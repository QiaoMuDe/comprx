package zip

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"gitee.com/MM-Q/comprx/config"
	"gitee.com/MM-Q/comprx/internal/utils"
)

// Unzip 解压缩 ZIP 文件到指定目录
//
// 参数:
//   - zipFilePath: 要解压缩的 ZIP 文件路径
//   - targetDir: 解压缩后的目标目录路径
//   - config: 解压缩配置
//
// 返回值:
//   - error: 解压缩过程中发生的错误
func Unzip(zipFilePath string, targetDir string, config *config.Config) error {
	// 确保路径为绝对路径
	var absErr error
	if zipFilePath, absErr = utils.EnsureAbsPath(zipFilePath, "ZIP文件路径"); absErr != nil {
		return absErr
	}
	if targetDir, absErr = utils.EnsureAbsPath(targetDir, "目标目录路径"); absErr != nil {
		return absErr
	}

	// 打开 ZIP 文件
	zipReader, err := zip.OpenReader(zipFilePath)
	if err != nil {
		return fmt.Errorf("打开 ZIP 文件失败: %w", err)
	}
	defer func() { _ = zipReader.Close() }()

	// 检查目标目录是否存在, 如果不存在, 则创建
	if err := utils.EnsureDir(targetDir); err != nil {
		return fmt.Errorf("创建目标目录失败: %w", err)
	}

	// 创建大小跟踪器用于解压过程
	tracker := utils.NewSizeTracker()

	// 遍历 ZIP 文件中的每个文件或目录
	for _, file := range zipReader.File {
		// 验证压缩比（防Zip Bomb攻击）
		if err := utils.ValidateCompressionRatio(config,
			int64(file.UncompressedSize64),
			int64(file.CompressedSize64)); err != nil {
			return fmt.Errorf("文件 %s 压缩比验证失败: %w", file.Name, err)
		}

		// 验证单个文件大小
		if err := utils.ValidateFileSize(config, file.Name, int64(file.UncompressedSize64)); err != nil {
			return fmt.Errorf("文件 %s 大小验证失败: %w", file.Name, err)
		}

		// 验证并更新累计大小
		if err := tracker.AddSize(config, int64(file.UncompressedSize64)); err != nil {
			return fmt.Errorf("累计大小验证失败: %w", err)
		}

		// 安全的路径验证和拼接
		targetPath, err := utils.ValidatePathSimple(targetDir, file.Name)
		if err != nil {
			return fmt.Errorf("处理文件 '%s' 时路径验证失败: %w", file.Name, err)
		}

		// 获取文件的模式
		mode := file.Mode()

		// 使用 switch 语句处理不同类型的文件
		switch {
		case mode.IsDir(): // 处理目录
			if err := extractDirectory(targetPath, file.Name); err != nil {
				return err
			}
		case mode&os.ModeSymlink != 0: // 处理软链接
			if err := extractSymlink(file, targetPath); err != nil {
				return err
			}
		default: // 处理普通文件
			if err := extractRegularFileWithValidation(file, targetPath, mode, config, tracker); err != nil {
				return err
			}
		}
	}

	return nil
}

// zipValidatingWriter 带验证功能的写入器包装器（用于ZIP解压）
type zipValidatingWriter struct {
	writer           io.Writer
	config           *config.Config
	compressedSize   int64
	uncompressedSize int64
	totalWritten     int64
	tracker          *utils.SizeTracker
}

// Write 实现io.Writer接口，在写入时进行安全验证
func (vw *zipValidatingWriter) Write(p []byte) (n int, err error) {
	// 写入数据
	n, err = vw.writer.Write(p)
	if err != nil {
		return n, err
	}

	// 更新总写入大小
	vw.totalWritten += int64(n)

	// 验证解压后的大小是否超过预期（基于ZIP文件头信息）
	if vw.totalWritten > vw.uncompressedSize {
		return n, fmt.Errorf("解压数据大小 %s 超过ZIP文件头声明的大小 %s",
			utils.FormatFileSize(vw.totalWritten), utils.FormatFileSize(vw.uncompressedSize))
	}

	// 验证解压后的大小是否超过单文件限制
	if vw.config.EnableSizeCheck && vw.totalWritten > vw.config.MaxFileSize {
		return n, fmt.Errorf("解压后文件大小 %s 超过单文件限制 %s",
			utils.FormatFileSize(vw.totalWritten), utils.FormatFileSize(vw.config.MaxFileSize))
	}

	// 验证压缩比（防止Zip Bomb攻击）
	if err := utils.ValidateCompressionRatio(vw.config, vw.totalWritten, vw.compressedSize); err != nil {
		return n, fmt.Errorf("压缩比验证失败: %w", err)
	}

	return n, nil
}

// extractDirectory 处理目录解压
//
// 参数:
//   - targetPath: 目标路径
//   - fileName: 文件名（用于错误信息）
//
// 返回值:
//   - error: 操作过程中遇到的错误
func extractDirectory(targetPath, fileName string) error {
	if err := utils.EnsureDir(targetPath); err != nil {
		return fmt.Errorf("处理目录 '%s' 时出错 - 创建目录失败: %w", fileName, err)
	}
	return nil
}

// extractSymlink 处理软链接解压
//
// 参数:
//   - file: ZIP文件条目
//   - targetPath: 目标路径
//
// 返回值:
//   - error: 操作过程中遇到的错误
func extractSymlink(file *zip.File, targetPath string) error {
	zipFileReader, err := file.Open()
	if err != nil {
		return fmt.Errorf("处理软链接 '%s' 时出错 - 打开 ZIP 文件中的软链接失败: %w", file.Name, err)
	}
	defer func() { _ = zipFileReader.Close() }()

	// 使用 io.ReadAll 读取完整的软链接目标路径
	targetBytes, err := io.ReadAll(zipFileReader)
	if err != nil {
		return fmt.Errorf("处理软链接 '%s' 时出错 - 读取软链接目标失败: %w", file.Name, err)
	}
	target := string(targetBytes) // 软链接的目标

	// 检查软链接的父目录是否存在，如果不存在，则创建
	parentDir := filepath.Dir(targetPath)
	if err := utils.EnsureDir(parentDir); err != nil {
		return fmt.Errorf("处理软链接 '%s' 时出错 - 创建软链接父目录失败: %w", file.Name, err)
	}

	// 创建软链接
	if err := os.Symlink(target, targetPath); err != nil {
		return fmt.Errorf("处理软链接 '%s' 时出错 - 创建软链接失败: %w", file.Name, err)
	}

	return nil
}

// extractRegularFileWithValidation 处理普通文件解压（带实时验证）
//
// 参数:
//   - file: ZIP文件条目
//   - targetPath: 目标路径
//   - mode: 文件模式
//   - config: 解压配置
//   - tracker: 大小跟踪器
//
// 返回值:
//   - error: 操作过程中遇到的错误
func extractRegularFileWithValidation(file *zip.File, targetPath string, mode os.FileMode, config *config.Config, tracker *utils.SizeTracker) error {
	// 创建带验证的写入器包装器
	return extractRegularFileWithWriter(file, targetPath, mode, config, func(fileWriter io.Writer) io.Writer {
		return &zipValidatingWriter{
			writer:           fileWriter,
			config:           config,
			compressedSize:   int64(file.CompressedSize64),
			uncompressedSize: int64(file.UncompressedSize64),
			tracker:          tracker,
		}
	})
}

// extractRegularFile 处理普通文件解压
//
// 参数:
//   - file: ZIP文件条目
//   - targetPath: 目标路径
//   - mode: 文件模式
//   - config: 解压配置
//
// 返回值:
//   - error: 操作过程中遇到的错误
func extractRegularFile(file *zip.File, targetPath string, mode os.FileMode, config *config.Config) error {
	return extractRegularFileWithWriter(file, targetPath, mode, config, func(fileWriter io.Writer) io.Writer {
		return fileWriter // 不使用验证包装器
	})
}

// extractRegularFileWithWriter 处理普通文件解压的通用实现
//
// 参数:
//   - file: ZIP文件条目
//   - targetPath: 目标路径
//   - mode: 文件模式
//   - config: 解压配置
//   - writerWrapper: 写入器包装函数
//
// 返回值:
//   - error: 操作过程中遇到的错误
func extractRegularFileWithWriter(file *zip.File, targetPath string, mode os.FileMode, config *config.Config, writerWrapper func(io.Writer) io.Writer) error {
	// 检查目标文件是否已存在
	if _, err := os.Stat(targetPath); err == nil {
		// 文件已存在，检查是否允许覆盖
		if !config.OverwriteExisting {
			return fmt.Errorf("目标文件已存在且不允许覆盖: %s", targetPath)
		}
	}

	// 检查file的父目录是否存在, 如果不存在, 则创建
	parentDir := filepath.Dir(targetPath)
	if err := utils.EnsureDir(parentDir); err != nil {
		return fmt.Errorf("处理文件 '%s' 时出错 - 创建文件父目录失败: %w", file.Name, err)
	}

	// 获取文件的大小
	fileSize := file.UncompressedSize64

	// 如果文件大小为0，只创建空文件，不进行读写操作
	if fileSize == 0 {
		// 创建空文件
		emptyFile, err := os.OpenFile(targetPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, mode.Perm())
		if err != nil {
			return fmt.Errorf("处理文件 '%s' 时出错 - 创建空文件失败: %w", file.Name, err)
		}
		defer func() { _ = emptyFile.Close() }()
		return nil
	}

	// 创建文件
	fileWriter, err := os.OpenFile(targetPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, mode.Perm())
	if err != nil {
		return fmt.Errorf("处理文件 '%s' 时出错 - 创建文件失败: %w", file.Name, err)
	}
	defer func() { _ = fileWriter.Close() }()

	// 使用包装器包装写入器
	wrappedWriter := writerWrapper(fileWriter)

	// 打开 ZIP 文件中的文件
	zipFileReader, err := file.Open()
	if err != nil {
		return fmt.Errorf("处理文件 '%s' 时出错 - 打开 zip 文件中的文件失败: %w", file.Name, err)
	}
	defer func() { _ = zipFileReader.Close() }()

	// 获取对应文件大小的缓冲区
	bufferSize := utils.GetBufferSize(int64(fileSize))

	// 创建缓冲区
	buffer := utils.GetBuffer(bufferSize)
	defer utils.PutBuffer(buffer)

	// 将文件内容写入目标文件（使用包装后的写入器）
	if _, err := io.CopyBuffer(wrappedWriter, zipFileReader, buffer); err != nil {
		return fmt.Errorf("处理文件 '%s' 时出错 - 写入文件失败: %w", file.Name, err)
	}

	return nil
}

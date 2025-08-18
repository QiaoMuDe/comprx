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

	// 遍历 ZIP 文件中的每个文件或目录
	for _, file := range zipReader.File {
		// 获取目标路径
		targetPath := filepath.Join(targetDir, file.Name)

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
			if err := extractRegularFile(file, targetPath, mode, config); err != nil {
				return err
			}
		}
	}

	return nil
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

	// 将文件内容写入目标文件
	if _, err := io.CopyBuffer(fileWriter, zipFileReader, buffer); err != nil {
		return fmt.Errorf("处理文件 '%s' 时出错 - 写入文件失败: %w", file.Name, err)
	}

	return nil
}

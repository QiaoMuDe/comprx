package tgz

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"gitee.com/MM-Q/comprx/config"
	"gitee.com/MM-Q/comprx/internal/utils"
)

// Untgz 解压缩 TGZ 文件到指定目录
//
// 参数:
//   - tgzFilePath: 要解压缩的 TGZ 文件路径
//   - targetDir: 解压缩后的目标目录路径
//   - config: 解压缩配置
//
// 返回值:
//   - error: 解压缩过程中发生的错误
func Untgz(tgzFilePath string, targetDir string, config *config.Config) error {
	// 打开 TGZ 文件
	tgzFile, err := os.Open(tgzFilePath)
	if err != nil {
		return fmt.Errorf("打开 TGZ 文件失败: %w", err)
	}
	defer func() { _ = tgzFile.Close() }()

	// 创建 GZIP 读取器
	gzipReader, err := gzip.NewReader(tgzFile)
	if err != nil {
		return fmt.Errorf("创建 GZIP 读取器失败: %w", err)
	}
	defer func() { _ = gzipReader.Close() }()

	// 创建 TAR 读取器
	tarReader := tar.NewReader(gzipReader)

	// 检查目标目录是否存在, 如果不存在, 则创建
	if err := utils.EnsureDir(targetDir); err != nil {
		return fmt.Errorf("创建目标目录失败: %w", err)
	}

	// 遍历 TAR 文件中的每个文件或目录
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break // 到达文件末尾
		}
		if err != nil {
			return fmt.Errorf("读取 TAR 文件头失败: %w", err)
		}

		// 安全的路径验证和拼接
		targetPath, err := utils.ValidatePathSimple(targetDir, header.Name, config.DisablePathValidation)
		if err != nil {
			return fmt.Errorf("处理文件 '%s' 时路径验证失败: %w", header.Name, err)
		}

		// 使用 switch 语句处理不同类型的文件
		switch header.Typeflag {
		case tar.TypeDir: // 处理目录
			if err := extractDirectory(targetPath, header.Name); err != nil {
				return err
			}
		case tar.TypeReg: // 处理普通文件
			if err := extractRegularFile(tarReader, targetPath, header, config); err != nil {
				return err
			}
		case tar.TypeSymlink: // 处理符号链接
			if err := extractSymlink(header, targetPath); err != nil {
				return err
			}
		case tar.TypeLink: // 处理硬链接
			if err := extractHardlink(header, targetPath, targetDir); err != nil {
				return err
			}
		default:
			// 对于其他类型的文件，我们跳过处理
			fmt.Printf("跳过不支持的文件类型: %s (类型: %c)\n", header.Name, header.Typeflag)
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

// extractRegularFile 处理普通文件解压
//
// 参数:
//   - tarReader: TAR读取器
//   - targetPath: 目标路径
//   - header: TAR文件头
//   - config: 解压配置
//
// 返回值:
//   - error: 操作过程中遇到的错误
func extractRegularFile(tarReader *tar.Reader, targetPath string, header *tar.Header, config *config.Config) error {
	// 检查目标文件是否已存在
	if _, err := os.Stat(targetPath); err == nil {
		// 文件已存在，检查是否允许覆盖
		if !config.OverwriteExisting {
			return fmt.Errorf("目标文件已存在且不允许覆盖: %s", targetPath)
		}
	}

	// 检查文件的父目录是否存在, 如果不存在, 则创建
	parentDir := filepath.Dir(targetPath)
	if err := utils.EnsureDir(parentDir); err != nil {
		return fmt.Errorf("处理文件 '%s' 时出错 - 创建文件父目录失败: %w", header.Name, err)
	}

	// 获取文件的大小
	fileSize := header.Size

	// 如果文件大小为0，只创建空文件，不进行读写操作
	if fileSize == 0 {
		// 创建空文件
		emptyFile, err := os.OpenFile(targetPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.FileMode(header.Mode))
		if err != nil {
			return fmt.Errorf("处理文件 '%s' 时出错 - 创建空文件失败: %w", header.Name, err)
		}
		defer func() { _ = emptyFile.Close() }()
		return nil
	}

	// 创建文件
	fileWriter, err := os.OpenFile(targetPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.FileMode(header.Mode))
	if err != nil {
		return fmt.Errorf("处理文件 '%s' 时出错 - 创建文件失败: %w", header.Name, err)
	}
	defer func() { _ = fileWriter.Close() }()

	// 获取对应文件大小的缓冲区
	bufferSize := utils.GetBufferSize(fileSize)

	// 创建缓冲区
	buffer := utils.GetBuffer(bufferSize)
	defer utils.PutBuffer(buffer)

	// 将文件内容写入目标文件
	if _, err := io.CopyBuffer(fileWriter, tarReader, buffer); err != nil {
		return fmt.Errorf("处理文件 '%s' 时出错 - 写入文件失败: %w", header.Name, err)
	}

	return nil
}

// extractSymlink 处理软链接解压
//
// 参数:
//   - header: TAR文件头
//   - targetPath: 目标路径
//
// 返回值:
//   - error: 操作过程中遇到的错误
func extractSymlink(header *tar.Header, targetPath string) error {
	// 检查软链接的父目录是否存在，如果不存在，则创建
	parentDir := filepath.Dir(targetPath)
	if err := utils.EnsureDir(parentDir); err != nil {
		return fmt.Errorf("处理软链接 '%s' 时出错 - 创建软链接父目录失败: %w", header.Name, err)
	}

	// 创建软链接
	if err := os.Symlink(header.Linkname, targetPath); err != nil {
		return fmt.Errorf("处理软链接 '%s' 时出错 - 创建软链接失败: %w", header.Name, err)
	}

	return nil
}

// extractHardlink 处理硬链接解压
//
// 参数:
//   - header: TAR文件头
//   - targetPath: 目标路径
//   - targetDir: 目标目录
//
// 返回值:
//   - error: 操作过程中遇到的错误
func extractHardlink(header *tar.Header, targetPath, targetDir string) error {
	// 检查硬链接的父目录是否存在，如果不存在，则创建
	parentDir := filepath.Dir(targetPath)
	if err := utils.EnsureDir(parentDir); err != nil {
		return fmt.Errorf("处理硬链接 '%s' 时出错 - 创建硬链接父目录失败: %w", header.Name, err)
	}

	// 获取硬链接的源文件路径
	linkSourcePath := filepath.Join(targetDir, header.Linkname)

	// 创建硬链接
	if err := os.Link(linkSourcePath, targetPath); err != nil {
		return fmt.Errorf("处理硬链接 '%s' 时出错 - 创建硬链接失败: %w", header.Name, err)
	}

	return nil
}

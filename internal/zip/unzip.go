package zip

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"gitee.com/MM-Q/comprx/internal/utils"
)

// Unzip 解压缩 ZIP 文件到指定目录
//
// 参数:
//   - zipFilePath: 要解压缩的 ZIP 文件路径
//   - targetDir: 解压缩后的目标目录路径
//
// 返回值:
//   - error: 解压缩过程中发生的错误
func Unzip(zipFilePath string, targetDir string) error {
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
	defer zipReader.Close()

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
		// 目录解压
		case mode.IsDir():
			if err := utils.EnsureDir(targetPath); err != nil {
				return fmt.Errorf("处理目录 '%s' 时出错 - 创建目录失败: %w", file.Name, err)
			}

		// 软链接解压
		case mode&os.ModeSymlink != 0:
			zipFileReader, err := file.Open()
			if err != nil {
				return fmt.Errorf("处理软链接 '%s' 时出错 - 打开 ZIP 文件中的软链接失败: %w", file.Name, err)
			}

			// 使用 io.ReadAll 读取完整的软链接目标路径
			targetBytes, err := io.ReadAll(zipFileReader)
			if err != nil {
				_ = zipFileReader.Close()
				return fmt.Errorf("处理软链接 '%s' 时出错 - 读取软链接目标失败: %w", file.Name, err)
			}
			target := string(targetBytes) // 软链接的目标

			// 检查软链接的父目录是否存在，如果不存在，则创建
			parentDir := filepath.Dir(targetPath)
			if err := utils.EnsureDir(parentDir); err != nil {
				_ = zipFileReader.Close()
				return fmt.Errorf("处理软链接 '%s' 时出错 - 创建软链接父目录失败: %w", file.Name, err)
			}

			// 创建软链接
			if err := os.Symlink(target, targetPath); err != nil {
				_ = zipFileReader.Close()
				return fmt.Errorf("处理软链接 '%s' 时出错 - 创建软链接失败: %w", file.Name, err)
			}

			// 关闭读取器并检查错误
			if err := zipFileReader.Close(); err != nil {
				return fmt.Errorf("处理软链接 '%s' 时出错 - 关闭 zip 读取器失败: %w", file.Name, err)
			}

		// 普通文件解压
		default:
			// 检查file的父目录是否存在, 如果不存在, 则创建
			parentDir := filepath.Dir(targetPath)
			if err := utils.EnsureDir(parentDir); err != nil {
				return fmt.Errorf("处理文件 '%s' 时出错 - 创建文件父目录失败: %w", file.Name, err)
			}

			// 创建文件
			fileWriter, err := os.OpenFile(targetPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, mode.Perm())
			if err != nil {
				return fmt.Errorf("处理文件 '%s' 时出错 - 创建文件失败: %w", file.Name, err)
			}

			// 打开 ZIP 文件中的文件
			zipFileReader, err := file.Open()
			if err != nil {
				return fmt.Errorf("处理文件 '%s' 时出错 - 打开 zip 文件中的文件失败: %w", file.Name, err)
			}

			// 获取文件的大小
			fileSize := file.UncompressedSize64

			// 获取对应文件大小的缓冲区
			bufferSize := utils.GetBufferSize(int64(fileSize))

			// 创建缓冲区
			buffer := utils.GetBuffer(bufferSize)

			// 将文件内容写入目标文件
			_, copyErr := io.CopyBuffer(fileWriter, zipFileReader, buffer)
			// 归还缓冲区
			utils.PutBuffer(buffer)
			if copyErr != nil {
				return fmt.Errorf("处理文件 '%s' 时出错 - 写入文件失败: %w", file.Name, copyErr)
			}

			// 关闭文件并检查错误
			if err := fileWriter.Close(); err != nil {
				return fmt.Errorf("处理文件 '%s' 时出错 - 关闭文件失败: %w", file.Name, err)
			}

			// 关闭读取器并检查错误
			if err := zipFileReader.Close(); err != nil {
				return fmt.Errorf("处理文件 '%s' 时出错 - 关闭 zip 读取器失败: %w", file.Name, err)
			}
		}
	}

	return nil
}

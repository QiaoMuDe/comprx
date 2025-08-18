package gzip

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"gitee.com/MM-Q/comprx/config"
	"gitee.com/MM-Q/comprx/internal/utils"
)

// Ungzip 解压缩 GZIP 文件
//
// 参数:
//   - gzipFilePath: 要解压缩的 GZIP 文件路径
//   - targetPath: 解压缩后的目标文件路径
//   - config: 解压缩配置
//
// 返回值:
//   - error: 解压缩过程中发生的错误
func Ungzip(gzipFilePath string, targetPath string, config *config.Config) error {
	// 确保路径为绝对路径
	var absErr error
	if gzipFilePath, absErr = utils.EnsureAbsPath(gzipFilePath, "GZIP文件路径"); absErr != nil {
		return absErr
	}
	if targetPath, absErr = utils.EnsureAbsPath(targetPath, "目标文件路径"); absErr != nil {
		return absErr
	}

	// 检查GZIP文件是否存在
	if _, err := os.Stat(gzipFilePath); err != nil {
		return fmt.Errorf("GZIP文件不存在: %w", err)
	}

	// 打开 GZIP 文件
	gzipFile, err := os.Open(gzipFilePath)
	if err != nil {
		return fmt.Errorf("打开 GZIP 文件失败: %w", err)
	}
	defer func() { _ = gzipFile.Close() }()

	// 创建 GZIP 读取器
	gzipReader, err := gzip.NewReader(gzipFile)
	if err != nil {
		return fmt.Errorf("创建 GZIP 读取器失败: %w", err)
	}
	defer func() { _ = gzipReader.Close() }()

	// 如果目标路径是目录，则使用GZIP文件头中的原始文件名
	if stat, err := os.Stat(targetPath); err == nil && stat.IsDir() {
		if gzipReader.Name != "" {
			targetPath = filepath.Join(targetPath, gzipReader.Name)
		} else {
			// 如果GZIP文件头中没有原始文件名，则去掉.gz扩展名
			baseName := filepath.Base(gzipFilePath)
			baseName = strings.TrimSuffix(baseName, ".gz")
			targetPath = filepath.Join(targetPath, baseName)
		}
	}

	// 检查目标文件的父目录是否存在，如果不存在则创建
	parentDir := filepath.Dir(targetPath)
	if err := utils.EnsureDir(parentDir); err != nil {
		return fmt.Errorf("创建目标文件父目录失败: %w", err)
	}

	// 创建目标文件
	targetFile, err := os.Create(targetPath)
	if err != nil {
		return fmt.Errorf("创建目标文件失败: %w", err)
	}
	defer func() { _ = targetFile.Close() }()

	// 获取GZIP文件大小来估算缓冲区大小
	gzipInfo, err := gzipFile.Stat()
	if err != nil {
		return fmt.Errorf("获取GZIP文件信息失败: %w", err)
	}

	// 获取缓冲区大小并创建缓冲区
	bufferSize := utils.GetBufferSize(gzipInfo.Size())
	buffer := utils.GetBuffer(bufferSize)
	defer utils.PutBuffer(buffer)

	// 解压缩文件内容
	if _, err := io.CopyBuffer(targetFile, gzipReader, buffer); err != nil {
		return fmt.Errorf("解压缩文件失败: %w", err)
	}

	// 如果GZIP文件头中有修改时间信息，则设置目标文件的修改时间
	if !gzipReader.ModTime.IsZero() {
		if err := os.Chtimes(targetPath, gzipReader.ModTime, gzipReader.ModTime); err != nil {
			// 设置时间失败不是致命错误，只记录警告
			fmt.Printf("警告: 设置文件修改时间失败: %v\n", err)
		}
	}

	return nil
}

package bzip2

import (
	"compress/bzip2"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"gitee.com/MM-Q/comprx/config"
	"gitee.com/MM-Q/comprx/internal/utils"
)

// Unbz2 解压缩 BZIP2 文件
//
// 参数:
//   - bz2FilePath: 要解压缩的 BZIP2 文件路径
//   - targetPath: 解压缩后的目标文件路径
//   - config: 解压缩配置
//
// 返回值:
//   - error: 解压缩过程中发生的错误
func Unbz2(bz2FilePath string, targetPath string, config *config.Config) error {
	// 确保路径为绝对路径
	var absErr error
	if bz2FilePath, absErr = utils.EnsureAbsPath(bz2FilePath, "BZIP2文件路径"); absErr != nil {
		return absErr
	}
	if targetPath, absErr = utils.EnsureAbsPath(targetPath, "目标文件路径"); absErr != nil {
		return absErr
	}

	// 检查BZIP2文件是否存在
	if _, err := os.Stat(bz2FilePath); err != nil {
		return fmt.Errorf("BZIP2文件不存在: %w", err)
	}

	// 打开 BZIP2 文件
	bz2File, err := os.Open(bz2FilePath)
	if err != nil {
		return fmt.Errorf("打开 BZIP2 文件失败: %w", err)
	}
	defer func() { _ = bz2File.Close() }()

	// 创建 BZIP2 读取器
	bz2Reader := bzip2.NewReader(bz2File)

	// 如果目标路径是目录，则使用去掉.bz2扩展名的文件名
	if stat, err := os.Stat(targetPath); err == nil && stat.IsDir() {
		baseName := filepath.Base(bz2FilePath)
		baseName = strings.TrimSuffix(baseName, ".bz2")
		targetPath = filepath.Join(targetPath, baseName)
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

	// 获取BZIP2文件大小来估算缓冲区大小
	bz2Info, err := bz2File.Stat()
	if err != nil {
		return fmt.Errorf("获取BZIP2文件信息失败: %w", err)
	}

	// 获取缓冲区大小并创建缓冲区
	bufferSize := utils.GetBufferSize(bz2Info.Size())
	buffer := utils.GetBuffer(bufferSize)
	defer utils.PutBuffer(buffer)

	// 解压缩文件内容
	if _, err := io.CopyBuffer(targetFile, bz2Reader, buffer); err != nil {
		return fmt.Errorf("解压缩文件失败: %w", err)
	}

	return nil
}

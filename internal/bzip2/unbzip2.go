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
//   - cfg: 解压缩配置
//
// 返回值:
//   - error: 解压缩过程中发生的错误
func Unbz2(bz2FilePath string, targetPath string, cfg *config.Config) error {
	// 打印压缩文件信息
	cfg.Progress.Compressing(bz2FilePath)

	// 打开 BZIP2 文件（同时检查文件是否存在）
	bz2File, err := os.Open(bz2FilePath)
	if err != nil {
		return fmt.Errorf("打开 BZIP2 文件失败: %w", err)
	}
	defer func() { _ = bz2File.Close() }()

	// 获取BZIP2文件信息
	bz2Info, err := bz2File.Stat()
	if err != nil {
		return fmt.Errorf("获取BZIP2文件信息失败: %w", err)
	}

	// 创建 BZIP2 读取器
	bz2Reader := bzip2.NewReader(bz2File)

	// 检查目标路径状态，处理目录情况和覆盖检查
	if targetStat, err := os.Stat(targetPath); err == nil {
		if targetStat.IsDir() {
			// 目标是目录，生成文件名
			baseName := filepath.Base(bz2FilePath)
			baseName = strings.TrimSuffix(baseName, ".bz2")
			baseName = strings.TrimSuffix(baseName, ".bzip2")

			// 添加安全验证
			validatedPath, err := utils.ValidatePathSimple(targetPath, baseName, cfg.DisablePathValidation)
			if err != nil {
				return fmt.Errorf("BZIP2文件名包含不安全的路径: %w", err)
			}
			targetPath = validatedPath

			// 重新检查生成的目标文件是否存在
			if _, err := os.Stat(targetPath); err == nil && !cfg.OverwriteExisting {
				return fmt.Errorf("目标文件已存在且不允许覆盖: %s", targetPath)
			}
		} else {
			// 目标是文件，检查是否允许覆盖
			if !cfg.OverwriteExisting {
				return fmt.Errorf("目标文件已存在且不允许覆盖: %s", targetPath)
			}
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

	// 获取缓冲区大小并创建缓冲区
	bufferSize := utils.GetBufferSize(bz2Info.Size())
	buffer := utils.GetBuffer(bufferSize)
	defer utils.PutBuffer(buffer)

	// 打印解压缩进度
	cfg.Progress.Inflating(targetPath)

	// 解压缩文件内容到目标文件
	if _, err := io.CopyBuffer(targetFile, bz2Reader, buffer); err != nil {
		return fmt.Errorf("解压缩文件失败: %w", err)
	}

	return nil
}

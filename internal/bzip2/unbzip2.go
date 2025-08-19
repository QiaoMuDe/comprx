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

	// 打开 BZIP2 文件（同时检查文件是否存在）
	bz2File, err := os.Open(bz2FilePath)
	if err != nil {
		return fmt.Errorf("打开 BZIP2 文件失败: %w", err)
	}
	defer func() { _ = bz2File.Close() }()

	// 获取BZIP2文件信息用于预验证
	bz2Info, err := bz2File.Stat()
	if err != nil {
		return fmt.Errorf("获取BZIP2文件信息失败: %w", err)
	}

	// 预验证BZIP2文件大小（检查输入文件合理性）
	if config.EnableSizeCheck && bz2Info.Size() > config.MaxTotalSize {
		return fmt.Errorf("BZIP2文件大小 %s 超过处理限制 %s",
			utils.FormatFileSize(bz2Info.Size()), utils.FormatFileSize(config.MaxTotalSize))
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
			targetPath = filepath.Join(targetPath, baseName)

			// 重新检查生成的目标文件是否存在
			if _, err := os.Stat(targetPath); err == nil && !config.OverwriteExisting {
				return fmt.Errorf("目标文件已存在且不允许覆盖: %s", targetPath)
			}
		} else {
			// 目标是文件，检查是否允许覆盖
			if !config.OverwriteExisting {
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

	// 使用之前获取的bz2Info来估算缓冲区大小

	// 获取缓冲区大小并创建缓冲区
	bufferSize := utils.GetBufferSize(bz2Info.Size())
	buffer := utils.GetBuffer(bufferSize)
	defer utils.PutBuffer(buffer)

	// 创建大小跟踪器（用于解压过程中的大小检查）
	tracker := utils.NewSizeTracker()

	// 创建带验证的写入器包装器
	validatingWriter := &validatingWriter{
		writer:         targetFile,
		config:         config,
		compressedSize: bz2Info.Size(),
		tracker:        tracker,
	}

	// 解压缩文件内容（使用带验证的写入器）
	if _, err := io.CopyBuffer(validatingWriter, bz2Reader, buffer); err != nil {
		return fmt.Errorf("解压缩文件失败: %w", err)
	}

	return nil
}

// validatingWriter 带验证功能的写入器包装器
type validatingWriter struct {
	writer         io.Writer
	config         *config.Config
	compressedSize int64
	totalWritten   int64
	tracker        *utils.SizeTracker
}

// Write 实现io.Writer接口，在写入时进行安全验证
func (vw *validatingWriter) Write(p []byte) (n int, err error) {
	// 写入数据
	n, err = vw.writer.Write(p)
	if err != nil {
		return n, err
	}

	// 更新总写入大小
	vw.totalWritten += int64(n)

	// 验证解压后的大小是否超过单文件限制
	if vw.config.EnableSizeCheck && vw.totalWritten > vw.config.MaxFileSize {
		return n, fmt.Errorf("解压后文件大小 %s 超过单文件限制 %s",
			utils.FormatFileSize(vw.totalWritten), utils.FormatFileSize(vw.config.MaxFileSize))
	}

	// 验证解压后的大小是否超过总大小限制
	if vw.config.EnableSizeCheck && vw.totalWritten > vw.config.MaxTotalSize {
		return n, fmt.Errorf("解压后文件大小 %s 超过总大小限制 %s",
			utils.FormatFileSize(vw.totalWritten), utils.FormatFileSize(vw.config.MaxTotalSize))
	}

	// 验证压缩比（防止Zip Bomb攻击）
	if err := utils.ValidateCompressionRatio(vw.config, vw.totalWritten, vw.compressedSize); err != nil {
		return n, fmt.Errorf("压缩比验证失败: %w", err)
	}

	return n, nil
}

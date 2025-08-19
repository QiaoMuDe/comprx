package utils

import (
	"fmt"
	"os"
	"path/filepath"

	"gitee.com/MM-Q/comprx/config"
)

// ValidateFileSize 验证单个文件大小
//
// 参数:
//   - cfg: 压缩器配置
//   - filePath: 文件路径
//   - size: 文件大小（字节）
//
// 返回:
//   - error: 如果文件大小超过限制，返回错误
func ValidateFileSize(cfg *config.Config, filePath string, size int64) error {
	if !cfg.EnableSizeCheck {
		return nil
	}

	if size > cfg.MaxFileSize {
		return fmt.Errorf("文件 %s 大小 %s 超过单文件限制 %s",
			filePath, FormatFileSize(size), FormatFileSize(cfg.MaxFileSize))
	}
	return nil
}

// ValidateTotalSize 验证累计处理大小
//
// 参数:
//   - cfg: 压缩器配置
//   - currentTotal: 当前累计大小（字节）
//   - additionalSize: 要添加的大小（字节）
//
// 返回:
//   - error: 如果累计大小超过限制，返回错误
func ValidateTotalSize(cfg *config.Config, currentTotal, additionalSize int64) error {
	if !cfg.EnableSizeCheck {
		return nil
	}

	newTotal := currentTotal + additionalSize
	if newTotal > cfg.MaxTotalSize {
		return fmt.Errorf("累计处理大小 %s 超过总大小限制 %s",
			FormatFileSize(newTotal), FormatFileSize(cfg.MaxTotalSize))
	}

	return nil
}

// ValidateCompressionRatio 验证压缩比（防Zip Bomb攻击）
//
// 参数:
//   - cfg: 压缩器配置
//   - originalSize: 原始文件大小（字节）
//   - compressedSize: 压缩后大小（字节）
//
// 返回:
//   - error: 如果压缩比异常，返回错误
func ValidateCompressionRatio(cfg *config.Config, originalSize, compressedSize int64) error {
	if !cfg.EnableSizeCheck || compressedSize == 0 {
		return nil
	}

	ratio := float64(originalSize) / float64(compressedSize)

	if ratio > cfg.MaxCompressionRatio {
		return fmt.Errorf("压缩比 %.2f:1 超过安全限制 %.0f:1，可能存在Zip Bomb攻击",
			ratio, cfg.MaxCompressionRatio)
	}
	return nil
}

// PreCheckDirectorySize 预检查目录总大小
//
// 参数:
//   - cfg: 压缩器配置
//   - dirPath: 目录路径
//
// 返回:
//   - int64: 目录总大小（字节）
//   - error: 检查过程中的错误
func PreCheckDirectorySize(cfg *config.Config, dirPath string) (int64, error) {
	if !cfg.EnableSizeCheck {
		return 0, nil
	}

	// 使用便捷函数检查目录是否存在
	if !Exists(dirPath) {
		return 0, fmt.Errorf("目录不存在: %s", dirPath)
	}

	var totalSize int64

	err := filepath.WalkDir(dirPath, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("遍历路径 %s 时出错: %w", path, err)
		}

		// 跳过目录
		if entry.IsDir() {
			return nil
		}

		info, err := entry.Info()
		if err != nil {
			return fmt.Errorf("获取文件信息失败 %s: %w", path, err)
		}

		fileSize := info.Size()

		// 验证单个文件大小
		if err := ValidateFileSize(cfg, path, fileSize); err != nil {
			return err
		}

		totalSize += fileSize

		// 验证累计大小
		if totalSize > cfg.MaxTotalSize {
			return fmt.Errorf("目录 %s 总大小 %s 超过限制 %s",
				dirPath, FormatFileSize(totalSize), FormatFileSize(cfg.MaxTotalSize))
		}

		return nil
	})

	return totalSize, err
}

// PreCheckSingleFile 预检查单个文件大小
//
// 参数:
//   - cfg: 压缩器配置
//   - filePath: 文件路径
//
// 返回:
//   - int64: 文件大小（字节）
//   - error: 检查过程中的错误
func PreCheckSingleFile(cfg *config.Config, filePath string) (int64, error) {
	if !cfg.EnableSizeCheck {
		return 0, nil
	}

	info, err := os.Stat(filePath)
	if err != nil {
		return 0, fmt.Errorf("获取文件信息失败 %s: %w", filePath, err)
	}

	fileSize := info.Size()

	// 验证单个文件大小
	if err := ValidateFileSize(cfg, filePath, fileSize); err != nil {
		return 0, err
	}

	// 验证是否超过总大小限制
	if fileSize > cfg.MaxTotalSize {
		return 0, fmt.Errorf("文件 %s 大小 %s 超过总大小限制 %s",
			filePath, FormatFileSize(fileSize), FormatFileSize(cfg.MaxTotalSize))
	}

	return fileSize, nil
}

// SizeTracker 大小跟踪器，用于在处理过程中跟踪累计大小
type SizeTracker struct {
	processedSize int64
}

// NewSizeTracker 创建大小跟踪器
//
// 返回:
//   - *SizeTracker: 大小跟踪器实例
func NewSizeTracker() *SizeTracker {
	return &SizeTracker{
		processedSize: 0,
	}
}

// AddSize 添加处理的大小并验证
//
// 参数:
//   - cfg: 压缩器配置
//   - size: 要添加的大小（字节）
//
// 返回:
//   - error: 如果累计大小超过限制，返回错误
func (st *SizeTracker) AddSize(cfg *config.Config, size int64) error {
	if err := ValidateTotalSize(cfg, st.processedSize, size); err != nil {
		return err
	}
	st.processedSize += size
	return nil
}

// GetProcessedSize 获取已处理的累计大小
//
// 返回:
//   - int64: 已处理的累计大小（字节）
func (st *SizeTracker) GetProcessedSize() int64 {
	return st.processedSize
}

// Reset 重置累计大小计数器
func (st *SizeTracker) Reset() {
	st.processedSize = 0
}

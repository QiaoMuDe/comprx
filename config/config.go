package config

import (
	"compress/gzip"
)

// CompressionLevel 压缩等级类型
type CompressionLevel int

const (
	// 压缩等级常量
	CompressionLevelDefault     CompressionLevel = -1 // 默认压缩等级
	CompressionLevelNone        CompressionLevel = 0  // 不压缩
	CompressionLevelFast        CompressionLevel = 1  // 快速压缩
	CompressionLevelBest        CompressionLevel = 9  // 最佳压缩
	CompressionLevelHuffmanOnly CompressionLevel = -2 // 仅使用Huffman编码
)

// Config 压缩器配置
type Config struct {
	CompressionLevel    CompressionLevel // 压缩等级
	OverwriteExisting   bool             // 是否覆盖已存在的文件
	MaxFileSize         int64            // 单个文件最大大小（字节）
	MaxTotalSize        int64            // 压缩包总大小限制（字节）
	EnableSizeCheck     bool             // 是否启用大小检查
	MaxCompressionRatio float64          // 最大压缩比限制（防Zip Bomb攻击）
}

// New 创建新的压缩器配置
func New() *Config {
	return &Config{
		CompressionLevel:    CompressionLevelDefault, // 默认压缩等级
		OverwriteExisting:   false,                   // 默认不覆盖已存在文件
		MaxFileSize:         100 * 1024 * 1024,       // 默认单文件最大100MB
		MaxTotalSize:        5120 * 1024 * 1024,      // 默认总大小最大5GB
		EnableSizeCheck:     true,                    // 默认启用大小检查
		MaxCompressionRatio: 500.0,                   // 默认最大压缩比500:1
	}
}

// GetCompressionLevel 根据配置返回对应的压缩等级
//
// 参数:
//   - cfg: *Config - 配置
//
// 返回值:
//   - int - 压缩等级
func GetCompressionLevel(cfg *Config) int {
	switch cfg.CompressionLevel {
	case CompressionLevelNone:
		return gzip.NoCompression // 不进行压缩
	case CompressionLevelFast:
		return gzip.BestSpeed // 快速压缩
	case CompressionLevelBest:
		return gzip.BestCompression // 最佳压缩
	case CompressionLevelHuffmanOnly:
		return gzip.HuffmanOnly // 只使用哈夫曼编码
	default:
		return gzip.DefaultCompression // 默认的压缩等级
	}
}

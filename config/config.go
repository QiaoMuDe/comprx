package config

import (
	"compress/gzip"

	"gitee.com/MM-Q/comprx/internal/progress"
	"gitee.com/MM-Q/comprx/types"
)

// Config 压缩器配置
type Config struct {
	CompressionLevel      types.CompressionLevel // 压缩等级
	OverwriteExisting     bool                   // 是否覆盖已存在的文件
	Progress              *progress.Progress     // 进度显示
	DisablePathValidation bool                   // 是否禁用路径验证
}

// New 创建新的压缩器配置
func New() *Config {
	return &Config{
		CompressionLevel:      types.CompressionLevelDefault, // 默认压缩等级
		OverwriteExisting:     false,                         // 默认不覆盖已存在文件
		Progress:              progress.New(),                // 创建进度显示
		DisablePathValidation: false,                         // 默认启用路径验证
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
	// 不进行压缩
	case types.CompressionLevelNone:
		return gzip.NoCompression

	// 快速压缩
	case types.CompressionLevelFast:
		return gzip.BestSpeed

	// 最佳压缩
	case types.CompressionLevelBest:
		return gzip.BestCompression

	// 只使用哈夫曼编码
	case types.CompressionLevelHuffmanOnly:
		return gzip.HuffmanOnly

	// 默认的压缩等级
	default:
		return gzip.DefaultCompression
	}
}

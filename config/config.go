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
	CompressionLevel  CompressionLevel // 压缩等级
	OverwriteExisting bool             // 是否覆盖已存在的文件
	ExcludePatterns   []string         // 排除的文件模式
}

// New 创建新的压缩器配置
func New() *Config {
	return &Config{
		CompressionLevel:  CompressionLevelDefault, // 默认压缩等级
		OverwriteExisting: false,                   // 默认不覆盖已存在文件
		ExcludePatterns:   []string{},              // 默认不排除任何文件
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

// // SetCompression 设置是否启用压缩
// func (c *Config) SetCompression(enabled bool) *Config {
// 	c.EnableCompression = enabled
// 	return c
// }

// // SetOverwrite 设置是否覆盖已存在文件
// func (c *Config) SetOverwrite(enabled bool) *Config {
// 	c.OverwriteExisting = enabled
// 	return c
// }

// // SetExcludePatterns 设置排除的文件模式
// func (c *Config) SetExcludePatterns(patterns ...string) *Config {
// 	c.ExcludePatterns = patterns
// 	return c
// }

// // AddExcludePattern 添加排除的文件模式
// func (c *Config) AddExcludePattern(pattern string) *Config {
// 	c.ExcludePatterns = append(c.ExcludePatterns, pattern)
// 	return c
// }

// // shouldExclude 检查文件是否应该被排除
// func (c *Config) shouldExclude(filePath string) bool {
// 	fileName := filepath.Base(filePath)
// 	for _, pattern := range c.ExcludePatterns {
// 		// 简单的通配符匹配
// 		if matched, _ := filepath.Match(pattern, fileName); matched {
// 			return true
// 		}
// 		// 检查是否包含指定字符串
// 		if strings.Contains(fileName, pattern) {
// 			return true
// 		}
// 	}
// 	return false
// }

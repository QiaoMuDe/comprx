package config

// Config 压缩器配置
type Config struct {
	EnableCompression bool     // 是否启用压缩
	OverwriteExisting bool     // 是否覆盖已存在的文件
	ExcludePatterns   []string // 排除的文件模式
}

// New 创建新的压缩器配置
func New() *Config {
	return &Config{
		EnableCompression: true,       // 默认启用压缩
		OverwriteExisting: false,      // 默认不覆盖已存在文件
		ExcludePatterns:   []string{}, // 默认不排除任何文件
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

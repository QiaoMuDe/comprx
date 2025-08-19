package progress

import (
	"fmt"
	"path/filepath"
)

// 进度条样式常量
const (
	// StyleText 文本样式进度条 - 使用文字描述进度
	StyleText = "text"

	// StyleUnicode Unicode样式进度条 - 使用Unicode字符绘制精美进度条
	// 示例: ████████████░░░░░░░░ 60%
	StyleUnicode = "unicode"

	// StyleASCII ASCII样式进度条 - 使用基础ASCII字符绘制兼容性最好的进度条
	// 示例: [##########          ] 50%
	StyleASCII = "ascii"
)

// IsSupportedStyle 判断是否是受支持的进度条样式
//
// 参数:
//   - style: 要检查的样式字符串
//
// 返回:
//   - bool: true表示支持该样式，false表示不支持
func IsSupportedStyle(style string) bool {
	switch style {
	case StyleText, StyleUnicode, StyleASCII:
		return true
	default:
		return false
	}
}

// ConsoleProgress 控制台进度显示器
type ConsoleProgress struct {
	enabled  bool   // 是否启用进度显示
	barStyle string // 进度条样式
}

// NewConsoleProgress 创建简单进度显示器
//
// 参数:
//   - enabled: 是否启用进度显示
//   - barStyle: 进度条样式
//
// 返回:
//   - *ConsoleProgress: 简单进度显示器
func NewConsoleProgress(enabled bool, barStyle string) *ConsoleProgress {
	// 检查样式是否支持,如果不支持则使用默认样式
	if !IsSupportedStyle(barStyle) {
		barStyle = StyleText
	}

	return &ConsoleProgress{
		enabled:  enabled,  // 是否启用进度显示
		barStyle: barStyle, // 进度条样式
	}
}

// Archive 显示压缩文件信息
//
// 参数:
//   - archivePath: 压缩文件路径
func (s *ConsoleProgress) Archive(archivePath string) {
	if s.enabled {
		fmt.Printf("Archive: %s\n", filepath.Base(archivePath))
	}
}

// Inflating 显示解压文件
//
// 参数:
//   - filePath: 文件路径
func (s *ConsoleProgress) Inflating(filePath string) {
	if s.enabled {
		fmt.Printf("\tinflating: %s\n", filePath)
	}
}

// Creating 显示创建目录
//
// 参数:
//   - dirPath: 目录路径
func (s *ConsoleProgress) Creating(dirPath string) {
	if s.enabled {
		fmt.Printf("\tcreating: %s\n", dirPath)
	}
}

// Extracting 显示提取文件（TAR）
//
// 参数:
//   - filePath: 文件路径
func (s *ConsoleProgress) Extracting(filePath string) {
	if s.enabled {
		fmt.Printf("\textracting: %s\n", filePath)
	}
}

// Adding 显示添加文件
//
// 参数:
//   - filePath: 文件路径
func (s *ConsoleProgress) Adding(filePath string) {
	if s.enabled {
		fmt.Printf("\tadding: %s\n", filePath)
	}
}

// Storing 显示存储目录
//
// 参数:
//   - dirPath: 目录路径
func (s *ConsoleProgress) Storing(dirPath string) {
	if s.enabled {
		fmt.Printf("\tstoring: %s\n", dirPath)
	}
}

// Compressing 显示压缩文件
//
// 参数:
//   - filePath: 文件路径
func (s *ConsoleProgress) Compressing(filePath string) {
	if s.enabled {
		fmt.Printf("compressing: %s\n", filePath)
	}
}

// File 通用文件处理显示（自动判断是解压还是创建目录）
//
// 参数:
//   - filePath: 文件路径
//   - isDirectory: 是否为目录
func (s *ConsoleProgress) File(filePath string, isDirectory bool) {
	if !s.enabled {
		return
	}

	if isDirectory {
		s.Creating(filePath)
	} else {
		s.Inflating(filePath)
	}
}

// IsEnabled 检查是否启用
//
// 返回:
//   - bool: 是否启用
func (s *ConsoleProgress) IsEnabled() bool {
	return s.enabled
}

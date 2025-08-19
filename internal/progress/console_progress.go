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

// 操作标签常量 - 确保冒号对齐
const (
	labelArchive     = "Archive:    " // 表示操作压缩包
	labelInflating   = "  inflating:" // 表示操作解压文件
	labelCreating    = "   creating:" // 表示操作创建目录
	labelExtracting  = " extracting:" // 表示操作解压文件(TAR)
	labelAdding      = "     adding:" // 表示操作添加文件
	labelStoring     = "    storing:" // 表示操作存储目录
	labelCompressing = "compressing:" // 表示操作压缩文件
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

// Progress 控制台进度显示器
type Progress struct {
	enabled  bool   // 是否启用进度显示
	barStyle string // 进度条样式
}

// NewProgress 创建简单进度显示器
//
// 参数:
//   - enabled: 是否启用进度显示
//   - barStyle: 进度条样式
//
// 返回:
//   - *Progress: 简单进度显示器
func NewProgress(enabled bool, barStyle string) *Progress {
	// 检查样式是否支持,如果不支持则使用默认样式
	if !IsSupportedStyle(barStyle) {
		barStyle = StyleText
	}

	return &Progress{
		enabled:  enabled,  // 是否启用进度显示
		barStyle: barStyle, // 进度条样式
	}
}

// IsEnabled 检查是否启用
//
// 返回:
//   - bool: 是否启用
func (s *Progress) IsEnabled() bool {
	return s.enabled
}

// Archive 显示压缩文件信息
//
// 参数:
//   - archivePath: 压缩文件路径
func (s *Progress) Archive(archivePath string) {
	if !s.enabled {
		return
	}
	fmt.Printf("%s %s\n", labelArchive, filepath.Base(archivePath))
}

// Compressing 显示压缩文件
//
// 参数:
//   - filePath: 文件路径
func (s *Progress) Compressing(filePath string) {
	if !s.enabled {
		return
	}
	fmt.Printf("%s %s\n", labelCompressing, filePath)
}

// Inflating 显示解压文件
//
// 参数:
//   - filePath: 文件路径
func (s *Progress) Inflating(filePath string) {
	if !s.enabled {
		return
	}
	fmt.Printf("%s %s\n", labelInflating, filePath)
}

// Creating 显示创建目录
//
// 参数:
//   - dirPath: 目录路径
func (s *Progress) Creating(dirPath string) {
	if !s.enabled {
		return
	}
	fmt.Printf("%s %s\n", labelCreating, dirPath)
}

// Extracting 显示提取文件（TAR）
//
// 参数:
//   - filePath: 文件路径
func (s *Progress) Extracting(filePath string) {
	if !s.enabled {
		return
	}
	fmt.Printf("%s %s\n", labelExtracting, filePath)
}

// Adding 显示添加文件
//
// 参数:
//   - filePath: 文件路径
func (s *Progress) Adding(filePath string) {
	if !s.enabled {
		return
	}
	fmt.Printf("%s %s\n", labelAdding, filePath)
}

// Storing 显示存储目录
//
// 参数:
//   - dirPath: 目录路径
func (s *Progress) Storing(dirPath string) {
	if !s.enabled {
		return
	}
	fmt.Printf("%s %s\n", labelStoring, dirPath)
}

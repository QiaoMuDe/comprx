package progress

import (
	"fmt"
	"path/filepath"

	"gitee.com/MM-Q/comprx/types"
	"github.com/schollz/progressbar/v3"
)

// 操作标签常量 - 确保冒号对齐
const (
	labelArchive     = "Archive:    " // 表示操作压缩包
	labelInflating   = "  inflating:" // 表示操作解压文件
	labelCreating    = "   creating:" // 表示操作创建目录
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
	case types.StyleText, types.StyleUnicode, types.StyleASCII:
		return true
	default:
		return false
	}
}

// Progress 控制台进度显示器
type Progress struct {
	Enabled  bool   // 是否启用进度显示
	BarStyle string // 进度条样式
}

// New 创建进度显示器
//
// 参数:
//   - enabled: 是否启用进度显示
//   - barStyle: 进度条样式
//
// 返回:
//   - *Progress: 简单进度显示器
func New() *Progress {
	return &Progress{
		Enabled:  true,            // 是否启用进度显示
		BarStyle: types.StyleText, // 进度条样式
	}
}

// NewProgressBar 创建一个进度条
//
// 参数:
//   - total: 进度条总大小
//   - description: 进度条描述信息
//
// 返回:
//   - *progressbar.ProgressBar: 进度条指针
//
// 进度条样式:
//   - types.StyleUnicode: Unicode样式进度条 - 使用Unicode字符绘制精美进度条
//   - types.StyleASCII: ASCII样式进度条 - 使用基础ASCII字符绘制兼容性最好的进度条
func (s *Progress) NewProgressBar(total int64, description string) *progressbar.ProgressBar {
	var theme progressbar.Theme
	// 如果设置样式为Unicode, 否则默认使用ASCII样式
	if s.BarStyle == types.StyleUnicode {
		theme = progressbar.ThemeUnicode
	} else {
		theme = progressbar.ThemeASCII
	}

	return progressbar.NewOptions64(
		total,                             // 进度条总大小
		progressbar.OptionClearOnFinish(), // 完成后清除进度条
		progressbar.OptionSetDescription(description), // 进度条描述信息
		progressbar.OptionSetElapsedTime(true),        // 显示已用时间
		progressbar.OptionSetPredictTime(true),        // 显示预计剩余时间
		progressbar.OptionSetRenderBlankState(true),   // 在进度条完成之前显示空白状态
		progressbar.OptionShowBytes(true),             // 显示进度条传输的字节
		progressbar.OptionShowCount(),                 // 显示当前进度的总和
		//progressbar.OptionShowElapsedTimeOnFinish(),        // 完成后显示已用时间
		progressbar.OptionSetTheme(theme), // ASCII 进度条主题(默认为 Unicode 进度条主题)
	)
}

// CloseBar 关闭进度条
//
// 参数:
//   - bar: 进度条指针
//
// 返回:
//   - error: 错误信息
func CloseBar(bar *progressbar.ProgressBar) error {
	// 如果进度条为空，则返回
	if bar == nil {
		return nil
	}

	// 完成进度条
	if err := bar.Finish(); err != nil {
		return err
	}

	// 关闭进度条
	if err := bar.Close(); err != nil {
		return err
	}

	return nil
}

// IsEnabled 检查是否启用
//
// 返回:
//   - bool: 是否启用
func (s *Progress) IsEnabled() bool {
	return s.Enabled
}

// Archive 显示压缩文件信息
//
// 参数:
//   - archivePath: 压缩文件路径
func (s *Progress) Archive(archivePath string) {
	// 如果不启用进度显示, 则直接返回
	// 如果进度条样式不是文本样式, 则直接返回
	if !s.Enabled || s.BarStyle != types.StyleText {
		return
	}
	fmt.Printf("%s %s\n", labelArchive, filepath.Base(archivePath))
}

// Compressing 显示压缩文件
//
// 参数:
//   - filePath: 文件路径
func (s *Progress) Compressing(filePath string) {
	if !s.Enabled || s.BarStyle != types.StyleText {
		return
	}
	fmt.Printf("%s %s\n", labelCompressing, filePath)
}

// ======================================================
// 解压进度
// ======================================================

// Inflating 显示解压文件
//
// 参数:
//   - filePath: 文件路径
func (s *Progress) Inflating(filePath string) {
	if !s.Enabled || s.BarStyle != types.StyleText {
		return
	}
	fmt.Printf("%s %s\n", labelInflating, filePath)
}

// Creating 显示创建目录
//
// 参数:
//   - dirPath: 目录路径
func (s *Progress) Creating(dirPath string) {
	if !s.Enabled || s.BarStyle != types.StyleText {
		return
	}
	fmt.Printf("%s %s\n", labelCreating, dirPath)
}

// ======================================================
// 压缩进度
// ======================================================

// Adding 显示添加文件
//
// 参数:
//   - filePath: 文件路径
func (s *Progress) Adding(filePath string) {
	if !s.Enabled || s.BarStyle != types.StyleText {
		return
	}
	fmt.Printf("%s %s\n", labelAdding, filePath)
}

// Storing 显示存储目录
//
// 参数:
//   - dirPath: 目录路径
func (s *Progress) Storing(dirPath string) {
	if !s.Enabled || s.BarStyle != types.StyleText {
		return
	}
	fmt.Printf("%s %s\n", labelStoring, dirPath)
}

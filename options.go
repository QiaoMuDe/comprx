package comprx

import (
	"gitee.com/MM-Q/comprx/config"
	"gitee.com/MM-Q/comprx/types"
)

// Options 压缩/解压配置选项
type Options struct {
	CompressionLevel      config.CompressionLevel // 压缩等级
	OverwriteExisting     bool                    // 是否覆盖已存在的文件
	ProgressEnabled       bool                    // 是否启用进度显示
	ProgressStyle         types.ProgressStyle     // 进度条样式
	DisablePathValidation bool                    // 是否禁用路径验证
}

// DefaultOptions 返回默认配置选项
//
// 返回:
//   - Options: 默认配置选项
//
// 默认配置:
//   - CompressionLevel: 默认压缩等级
//   - OverwriteExisting: false (不覆盖已存在文件)
//   - ProgressEnabled: false (不显示进度)
//   - ProgressStyle: 文本样式
//   - DisablePathValidation: false (启用路径验证)
func DefaultOptions() Options {
	return Options{
		CompressionLevel:      config.CompressionLevelDefault,
		OverwriteExisting:     false,
		ProgressEnabled:       false,
		ProgressStyle:         types.ProgressStyleText,
		DisablePathValidation: false,
	}
}

// ProgressOptions 返回带进度显示的配置选项
//
// 参数:
//   - style: 进度条样式
//
// 返回:
//   - Options: 带进度显示的配置选项
func ProgressOptions(style types.ProgressStyle) Options {
	opts := DefaultOptions()
	opts.ProgressEnabled = true
	opts.ProgressStyle = style
	return opts
}

// TextProgressOptions 返回文本样式进度条配置选项
//
// 返回:
//   - Options: 文本样式进度条配置选项
//
// 使用示例:
//
//	err := PackOptions("output.zip", "input_dir", TextProgressOptions())
func TextProgressOptions() Options {
	opts := DefaultOptions()
	opts.ProgressEnabled = true
	opts.ProgressStyle = types.ProgressStyleText
	return opts
}

// UnicodeProgressOptions 返回Unicode样式进度条配置选项
//
// 返回:
//   - Options: Unicode样式进度条配置选项
//
// 使用示例:
//
//	err := PackOptions("output.zip", "input_dir", UnicodeProgressOptions())
func UnicodeProgressOptions() Options {
	opts := DefaultOptions()
	opts.ProgressEnabled = true
	opts.ProgressStyle = types.ProgressStyleUnicode
	return opts
}

// ASCIIProgressOptions 返回ASCII样式进度条配置选项
//
// 返回:
//   - Options: ASCII样式进度条配置选项
//
// 使用示例:
//
//	err := PackOptions("output.zip", "input_dir", ASCIIProgressOptions())
func ASCIIProgressOptions() Options {
	opts := DefaultOptions()
	opts.ProgressEnabled = true
	opts.ProgressStyle = types.ProgressStyleASCII
	return opts
}

// ForceOptions 返回强制模式配置选项
//
// 返回:
//   - Options: 强制模式配置选项
//
// 配置特点:
//   - OverwriteExisting: true (覆盖已存在文件)
//   - DisablePathValidation: true (禁用路径验证)
//   - ProgressEnabled: false (关闭进度条)
//
// 使用示例:
//
//	err := PackOptions("output.zip", "input_dir", ForceOptions())
func ForceOptions() Options {
	opts := DefaultOptions()
	opts.OverwriteExisting = true
	opts.DisablePathValidation = true
	opts.ProgressEnabled = false
	return opts
}

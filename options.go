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

// FastOptions 返回快速压缩配置选项
//
// 返回:
//   - Options: 快速压缩配置选项
//
// 配置特点:
//   - 使用快速压缩等级
//   - 启用进度显示
//   - 允许覆盖已存在文件
func FastOptions() Options {
	return Options{
		CompressionLevel:      config.CompressionLevelFast, // 快速压缩
		OverwriteExisting:     true,                        // 允许覆盖已存在文件
		ProgressEnabled:       true,                        // 启用进度显示
		ProgressStyle:         types.ProgressStyleText,     // 文本样式
		DisablePathValidation: false,                       // 启用路径验证
	}
}

// BestOptions 返回最佳压缩配置选项
//
// 返回:
//   - Options: 最佳压缩配置选项
//
// 配置特点:
//   - 使用最佳压缩等级
//   - 启用进度显示
//   - 允许覆盖已存在文件
func BestOptions() Options {
	return Options{
		CompressionLevel:      config.CompressionLevelBest, // 最佳压缩
		OverwriteExisting:     true,                        // 允许覆盖已存在文件
		ProgressEnabled:       true,                        // 启用进度显示
		ProgressStyle:         types.ProgressStyleText,     // 文本样式
		DisablePathValidation: false,                       // 启用路径验证
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

// QuietOptions 返回静默模式配置选项
//
// 返回:
//   - Options: 静默模式配置选项
//
// 配置特点:
//   - 不显示进度
//   - 允许覆盖已存在文件
//   - 禁用路径验证（提高性能）
func QuietOptions() Options {
	return Options{
		CompressionLevel:      config.CompressionLevelDefault,
		OverwriteExisting:     true,
		ProgressEnabled:       false,
		ProgressStyle:         types.ProgressStyleText,
		DisablePathValidation: true,
	}
}

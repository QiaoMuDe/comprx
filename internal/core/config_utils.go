package core

import (
	"gitee.com/MM-Q/comprx/config"
	"gitee.com/MM-Q/comprx/types"
)

// ==============================================
// 配置方法
// ==============================================

// WithProgressAndStyle 设置进度条样式和是否启用进度条
//
// 参数:
//   - enabled: 是否启用进度条
//   - style: 进度条样式
//
// 返回:
//   - *Comprx: 压缩器实例
func (c *Comprx) WithProgressAndStyle(enabled bool, style types.ProgressStyle) *Comprx {
	c.SetProgressAndStyle(enabled, style)
	return c
}

// SetProgressAndStyle 设置进度条样式和是否启用进度条
//
// 参数:
//   - enabled: 是否启用进度条
//   - style: 进度条样式
func (c *Comprx) SetProgressAndStyle(enabled bool, style types.ProgressStyle) {
	if !style.IsValid() {
		style = types.ProgressStyleText
	}

	c.Config.Progress.Enabled = enabled
	c.Config.Progress.BarStyle = style
}

// WithOverwriteExisting 设置是否覆盖已存在的文件
//
// 参数:
//   - overwrite: 是否覆盖已存在文件
//
// 返回:
//   - *Comprx: 压缩器实例
func (c *Comprx) WithOverwriteExisting(overwrite bool) *Comprx {
	c.SetOverwriteExisting(overwrite)
	return c
}

// SetOverwriteExisting 设置是否覆盖已存在的文件
//
// 参数:
//   - overwrite: 是否覆盖已存在文件
func (c *Comprx) SetOverwriteExisting(overwrite bool) {
	c.Config.OverwriteExisting = overwrite
}

// WithCompressionLevel 设置压缩级别
//
// 参数:
//   - level: 压缩级别
//
// 返回:
//   - *Comprx: 压缩器实例
func (c *Comprx) WithCompressionLevel(level config.CompressionLevel) *Comprx {
	c.SetCompressionLevel(level)
	return c
}

// SetCompressionLevel 设置压缩级别
//
// 参数:
//   - level: 压缩级别
func (c *Comprx) SetCompressionLevel(level config.CompressionLevel) {
	c.Config.CompressionLevel = level
}

// WithDisablePathValidation 设置禁用路径验证
//
// 参数:
//   - disable: 是否禁用路径验证
//
// 返回:
//   - *Comprx: 压缩器实例
func (c *Comprx) WithDisablePathValidation(disable bool) *Comprx {
	c.SetDisablePathValidation(disable)
	return c
}

// SetDisablePathValidation 设置禁用路径验证
//
// 参数:
//   - disable: 是否禁用路径验证
func (c *Comprx) SetDisablePathValidation(disable bool) {
	c.Config.DisablePathValidation = disable
}

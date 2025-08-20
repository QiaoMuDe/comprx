package comprx

import (
	"gitee.com/MM-Q/comprx/internal/core"
)

// ==============================================
// 简单便捷函数 - 线程安全版本
// ==============================================

// Pack 压缩文件或目录(禁用进度条) - 线程安全
//
// 参数:
//   - dst: 目标文件路径
//   - src: 源文件路径
//
// 返回:
//   - error: 错误信息
//
// 使用示例:
//
//	err := Pack("output.zip", "input_dir")
func Pack(dst string, src string) error {
	return PackWithOptions(dst, src, DefaultOptions())
}

// Unpack 解压文件(禁用进度条) - 线程安全
//
// 参数:
//   - src: 源文件路径
//   - dst: 目标目录路径
//
// 返回:
//   - error: 错误信息
//
// 使用示例:
//
//	err := Unpack("archive.zip", "output_dir")
func Unpack(src string, dst string) error {
	return UnpackWithOptions(src, dst, DefaultOptions())
}

// PackWithProgress 压缩文件或目录(启用进度条) - 线程安全
//
// 参数:
//   - dst: 目标文件路径
//   - src: 源文件路径
//
// 返回:
//   - error: 错误信息
//
// 使用示例:
//
//	err := PackWithProgress("output.zip", "input_dir")
func PackWithProgress(dst string, src string) error {
	opts := DefaultOptions()
	opts.ProgressEnabled = true
	return PackWithOptions(dst, src, opts)
}

// UnpackWithProgress 解压文件(启用进度条) - 线程安全
//
// 参数:
//   - src: 源文件路径
//   - dst: 目标目录路径
//
// 返回:
//   - error: 错误信息
//
// 使用示例:
//
//	err := UnpackWithProgress("archive.zip", "output_dir")
func UnpackWithProgress(src string, dst string) error {
	opts := DefaultOptions()
	opts.ProgressEnabled = true
	return UnpackWithOptions(src, dst, opts)
}

// ==============================================
// 配置化便捷函数 - 线程安全版本
// ==============================================

// PackWithOptions 使用指定配置压缩文件或目录 - 线程安全
//
// 参数:
//   - dst: 目标文件路径
//   - src: 源文件路径
//   - opts: 配置选项
//
// 返回:
//   - error: 错误信息
//
// 使用示例:
//
//	opts := Options{
//	    CompressionLevel: config.CompressionLevelBest,
//	    OverwriteExisting: true,
//	    ProgressEnabled: true,
//	    ProgressStyle: types.ProgressStyleUnicode,
//	}
//	err := PackWithOptions("output.zip", "input_dir", opts)
func PackWithOptions(dst string, src string, opts Options) error {
	comprx := core.New().
		WithCompressionLevel(opts.CompressionLevel).
		WithOverwriteExisting(opts.OverwriteExisting).
		WithProgressAndStyle(opts.ProgressEnabled, opts.ProgressStyle).
		WithDisablePathValidation(opts.DisablePathValidation)

	return comprx.Pack(dst, src)
}

// UnpackWithOptions 使用指定配置解压文件 - 线程安全
//
// 参数:
//   - src: 源文件路径
//   - dst: 目标目录路径
//   - opts: 配置选项
//
// 返回:
//   - error: 错误信息
//
// 使用示例:
//
//	opts := Options{
//	    OverwriteExisting: true,
//	    ProgressEnabled: true,
//	    ProgressStyle: types.ProgressStyleASCII,
//	}
//	err := UnpackWithOptions("archive.zip", "output_dir", opts)
func UnpackWithOptions(src string, dst string, opts Options) error {
	comprx := core.New().
		WithOverwriteExisting(opts.OverwriteExisting).
		WithProgressAndStyle(opts.ProgressEnabled, opts.ProgressStyle).
		WithDisablePathValidation(opts.DisablePathValidation)

	return comprx.Unpack(src, dst)
}

// ==============================================
// 快捷配置函数 - 线程安全版本
// ==============================================

// PackFast 快速压缩文件或目录 - 线程安全
//
// 参数:
//   - dst: 目标文件路径
//   - src: 源文件路径
//
// 返回:
//   - error: 错误信息
//
// 特点:
//   - 使用快速压缩等级
//   - 启用进度显示
//   - 允许覆盖已存在文件
//
// 使用示例:
//
//	err := PackFast("output.zip", "input_dir")
func PackFast(dst string, src string) error {
	return PackWithOptions(dst, src, FastOptions())
}

// PackBest 最佳压缩文件或目录 - 线程安全
//
// 参数:
//   - dst: 目标文件路径
//   - src: 源文件路径
//
// 返回:
//   - error: 错误信息
//
// 特点:
//   - 使用最佳压缩等级
//   - 启用进度显示
//   - 允许覆盖已存在文件
//
// 使用示例:
//
//	err := PackBest("output.zip", "input_dir")
func PackBest(dst string, src string) error {
	return PackWithOptions(dst, src, BestOptions())
}

// PackQuiet 静默压缩文件或目录 - 线程安全
//
// 参数:
//   - dst: 目标文件路径
//   - src: 源文件路径
//
// 返回:
//   - error: 错误信息
//
// 特点:
//   - 不显示进度
//   - 允许覆盖已存在文件
//   - 禁用路径验证（提高性能）
//
// 使用示例:
//
//	err := PackQuiet("output.zip", "input_dir")
func PackQuiet(dst string, src string) error {
	return PackWithOptions(dst, src, QuietOptions())
}

// UnpackQuiet 静默解压文件 - 线程安全
//
// 参数:
//   - src: 源文件路径
//   - dst: 目标目录路径
//
// 返回:
//   - error: 错误信息
//
// 特点:
//   - 不显示进度
//   - 允许覆盖已存在文件
//   - 禁用路径验证（提高性能）
//
// 使用示例:
//
//	err := UnpackQuiet("archive.zip", "output_dir")
func UnpackQuiet(src string, dst string) error {
	return UnpackWithOptions(src, dst, QuietOptions())
}

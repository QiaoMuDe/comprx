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
	return PackOptions(dst, src, DefaultOptions())
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
	return UnpackOptions(src, dst, DefaultOptions())
}

// PackProgress 压缩文件或目录(启用进度条) - 线程安全
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
func PackProgress(dst string, src string) error {
	opts := DefaultOptions()
	opts.ProgressEnabled = true
	return PackOptions(dst, src, opts)
}

// UnpackProgress 解压文件(启用进度条) - 线程安全
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
func UnpackProgress(src string, dst string) error {
	opts := DefaultOptions()
	opts.ProgressEnabled = true
	return UnpackOptions(src, dst, opts)
}

// ==============================================
// 配置化便捷函数 - 线程安全版本
// ==============================================

// PackOptions 使用指定配置压缩文件或目录 - 线程安全
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
//	err := PackOptions("output.zip", "input_dir", opts)
func PackOptions(dst string, src string, opts Options) error {
	comprx := core.New().
		WithCompressionLevel(opts.CompressionLevel).
		WithOverwriteExisting(opts.OverwriteExisting).
		WithProgressAndStyle(opts.ProgressEnabled, opts.ProgressStyle).
		WithDisablePathValidation(opts.DisablePathValidation)

	return comprx.Pack(dst, src)
}

// UnpackOptions 使用指定配置解压文件 - 线程安全
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
//	err := UnpackOptions("archive.zip", "output_dir", opts)
func UnpackOptions(src string, dst string, opts Options) error {
	comprx := core.New().
		WithOverwriteExisting(opts.OverwriteExisting).
		WithProgressAndStyle(opts.ProgressEnabled, opts.ProgressStyle).
		WithDisablePathValidation(opts.DisablePathValidation)

	return comprx.Unpack(src, dst)
}

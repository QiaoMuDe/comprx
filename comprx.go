package comprx

import (
	"gitee.com/MM-Q/comprx/internal/core"
	"gitee.com/MM-Q/comprx/internal/gzip"
	"gitee.com/MM-Q/comprx/types"
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

// ==================== 内存压缩API ====================

// GzipBytes 压缩字节数据
//
// 参数:
//   - data: 要压缩的字节数据
//   - level: 压缩级别
//
// 返回:
//   - []byte: 压缩后的数据
//   - error: 错误信息
//
// 使用示例:
//
//	compressed, err := GzipBytes([]byte("hello world"), types.CompressionLevelDefault)
func GzipBytes(data []byte, level types.CompressionLevel) ([]byte, error) {
	return gzip.CompressBytes(data, level)
}

// UngzipBytes 解压字节数据
//
// 参数:
//   - compressedData: 压缩的字节数据
//
// 返回:
//   - []byte: 解压后的数据
//   - error: 错误信息
//
// 使用示例:
//
//	decompressed, err := UngzipBytes(compressedData)
func UngzipBytes(compressedData []byte) ([]byte, error) {
	return gzip.DecompressBytes(compressedData)
}

// GzipString 压缩字符串
//
// 参数:
//   - text: 要压缩的字符串
//   - level: 压缩级别
//
// 返回:
//   - []byte: 压缩后的数据
//   - error: 错误信息
//
// 使用示例:
//
//	compressed, err := GzipString("hello world", types.CompressionLevelBest)
func GzipString(text string, level types.CompressionLevel) ([]byte, error) {
	return gzip.CompressString(text, level)
}

// UngzipString 解压为字符串
//
// 参数:
//   - compressedData: 压缩的字节数据
//
// 返回:
//   - string: 解压后的字符串
//   - error: 错误信息
//
// 使用示例:
//
//	text, err := UngzipString(compressedData)
func UngzipString(compressedData []byte) (string, error) {
	return gzip.DecompressString(compressedData)
}

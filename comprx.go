package comprx

import (
	"gitee.com/MM-Q/comprx/internal/core"
	"gitee.com/MM-Q/comprx/types"
)

// ==============================================
// 全局便捷函数 - 线程安全版本
// ==============================================

// Pack 压缩文件或目录(禁用进度条) - 线程安全
//
// 参数:
//   - dst: 目标文件路径
//   - src: 源文件路径
//
// 返回:
//   - error: 错误信息
func Pack(dst string, src string) error {
	comprx := core.New().WithProgressAndStyle(false, types.ProgressStyleText)
	return comprx.Pack(dst, src)
}

// Unpack 解压文件(禁用进度条) - 线程安全
//
// 参数:
//   - src: 源文件路径
//   - dst: 目标目录路径
//
// 返回:
//   - error: 错误信息
func Unpack(src string, dst string) error {
	comprx := core.New().WithProgressAndStyle(false, types.ProgressStyleText)
	return comprx.Unpack(src, dst)
}

// PackWithProgress 压缩文件或目录(启用进度条) - 线程安全
//
// 参数:
//   - dst: 目标文件路径
//   - src: 源文件路径
//
// 返回:
//   - error: 错误信息
func PackWithProgress(dst string, src string) error {
	comprx := core.New().WithProgressAndStyle(true, types.ProgressStyleText)
	return comprx.Pack(dst, src)
}

// UnpackWithProgress 解压文件(启用进度条) - 线程安全
//
// 参数:
//   - src: 源文件路径
//   - dst: 目标目录路径
//
// 返回:
//   - error: 错误信息
func UnpackWithProgress(src string, dst string) error {
	comprx := core.New().WithProgressAndStyle(true, types.ProgressStyleText)
	return comprx.Unpack(src, dst)
}

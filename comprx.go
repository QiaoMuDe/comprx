package comprx

import (
	"fmt"
	"path/filepath"

	"gitee.com/MM-Q/comprx/config"
	"gitee.com/MM-Q/comprx/internal/tar"
	"gitee.com/MM-Q/comprx/internal/utils"
	"gitee.com/MM-Q/comprx/internal/zip"
	"gitee.com/MM-Q/comprx/types"
)

// Comprx 压缩器
type Comprx struct {
	config *config.Config // 压缩器配置
}

// New 创建压缩器实例
//
// 返回:
//   - *Comprx: 压缩器实例
func New() *Comprx {
	return &Comprx{
		config: config.New(),
	}
}

// // WithEnableCompression 设置是否启用压缩
// //
// // 参数:
// //   - enable: 是否启用压缩
// //
// // 返回:
// //   - *Comprx: 压缩器实例
// func (c *Comprx) WithEnableCompression(enable bool) *Comprx {
// 	c.config.EnableCompression = enable
// 	return c
// }

// // WithOverwriteExisting 设置是否覆盖已存在的文件
// //
// // 参数:
// //   - overwrite: 是否覆盖已存在文件
// //
// // 返回:
// //   - *Comprx: 压缩器实例
// func (c *Comprx) WithOverwriteExisting(overwrite bool) *Comprx {
// 	c.config.OverwriteExisting = overwrite
// 	return c
// }

// // WithExcludePatterns 设置排除的文件模式
// //
// // 参数:
// //   - excludePatterns: 排除的模式列表
// //
// // 返回:
// //   - *Comprx: 压缩器实例
// func (c *Comprx) WithExcludePatterns(excludePatterns ...string) *Comprx {
// 	c.config.ExcludePatterns = excludePatterns
// 	return c
// }

// Pack 压缩文件或目录
func (c *Comprx) Pack(dst string, src string) error {
	// 智能检测压缩文件格式
	compressType, err := types.DetectCompressFormat(dst)
	if err != nil {
		return fmt.Errorf("检测压缩格式失败: %v", err)
	}

	// 检查目标文件是否存在
	if utils.Exists(dst) {
		if !c.config.OverwriteExisting {
			return fmt.Errorf("文件 %s 已存在，如需覆盖请设置 OverwriteExisting 为 true", dst)
		}
	}

	// 检查目标目录是否存在, 不存在则创建
	targetDir := filepath.Dir(dst)
	if err := utils.EnsureDir(targetDir); err != nil {
		return fmt.Errorf("创建目标目录失败: %v", err)
	}

	// 根据压缩格式进行打包
	switch compressType {
	case types.CompressTypeZip:
		return zip.Zip(dst, src, c.config)
	case types.CompressTypeTar:
		return tar.Tar(dst, src, c.config)
	default:
		return fmt.Errorf("不支持的压缩格式: %s", compressType)
	}
}

// Unpack 解压文件
//
// 参数:
//   - src: 源文件路径
//   - dst: 目标目录路径
//
// 返回:
//   - error: 错误信息
func (c *Comprx) Unpack(src string, dst string) error {
	// 智能检测压缩文件格式
	compressType, err := types.DetectCompressFormat(src)
	if err != nil {
		return fmt.Errorf("检测压缩格式失败: %v", err)
	}

	// 检查源文件是否存在
	if !utils.Exists(src) {
		return fmt.Errorf("源文件 %s 不存在", src)
	}

	// 检查目标目录是否存在, 不存在则创建
	if err := utils.EnsureDir(dst); err != nil {
		return fmt.Errorf("创建目标目录失败: %v", err)
	}

	// 根据压缩格式进行解压
	switch compressType {
	case types.CompressTypeZip:
		return zip.Unzip(src, dst, c.config)
	case types.CompressTypeTar:
		return tar.Untar(src, dst, c.config)
	default:
		return fmt.Errorf("不支持的压缩格式: %s", compressType)
	}
}

// // 便捷函数，使用默认配置
// var defaultComprx = New()

// // Pack 使用默认配置压缩文件
// func Pack(dst string, src string) error {
// 	return defaultComprx.Pack(dst, src)
// }

// // Unpack 使用默认配置解压文件
// func Unpack(src string, dst string) error {
// 	return defaultComprx.Unpack(src, dst)
// }

//func (c *Comprx) List(archivePath string) error

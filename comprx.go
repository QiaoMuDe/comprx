package comprx

import (
	"fmt"
	"path/filepath"
	"strings"

	"gitee.com/MM-Q/comprx/config"
	"gitee.com/MM-Q/comprx/internal/bzip2"
	"gitee.com/MM-Q/comprx/internal/gzip"
	"gitee.com/MM-Q/comprx/internal/tar"
	"gitee.com/MM-Q/comprx/internal/tgz"
	"gitee.com/MM-Q/comprx/internal/utils"
	"gitee.com/MM-Q/comprx/internal/zip"
	"gitee.com/MM-Q/comprx/types"
)

// Comprx 压缩器
type Comprx struct {
	config *config.Config // 压缩器配置
}

// ==============================================
// 构造函数
// ==============================================

// New 创建压缩器实例(NewComprx的别名)
//
// 返回:
//   - *Comprx: 压缩器实例
var New = NewComprx

// NewComprx 创建压缩器实例
//
// 返回:
//   - *Comprx: 压缩器实例
func NewComprx() *Comprx {
	return &Comprx{
		config: config.New(),
	}
}

// ==============================================
// 配置方法
// ==============================================

// WithOverwriteExisting 设置是否覆盖已存在的文件
//
// 参数:
//   - overwrite: 是否覆盖已存在文件
//
// 返回:
//   - *Comprx: 压缩器实例
func (c *Comprx) WithOverwriteExisting(overwrite bool) *Comprx {
	c.config.OverwriteExisting = overwrite
	return c
}

// SetOverwriteExisting 设置是否覆盖已存在的文件
//
// 参数:
//   - overwrite: 是否覆盖已存在文件
func (c *Comprx) SetOverwriteExisting(overwrite bool) {
	c.config.OverwriteExisting = overwrite
}

// WithCompressionLevel 设置压缩级别
//
// 参数:
//   - level: 压缩级别
//
// 返回:
//   - *Comprx: 压缩器实例
func (c *Comprx) WithCompressionLevel(level config.CompressionLevel) *Comprx {
	c.config.CompressionLevel = level
	return c
}

// SetCompressionLevel 设置压缩级别
//
// 参数:
//   - level: 压缩级别
func (c *Comprx) SetCompressionLevel(level config.CompressionLevel) {
	c.config.CompressionLevel = level
}

// ==============================================
// 打包压缩方法
// ==============================================

// Pack 压缩文件或目录
func (c *Comprx) Pack(dst string, src string) error {
	// 检查参数
	if src == "" || dst == "" {
		return fmt.Errorf("源文件路径或目标文件路径不能为空")
	}

	// 智能检测压缩文件格式
	compressType, err := types.DetectCompressFormat(dst)
	if err != nil {
		return fmt.Errorf("检测压缩格式失败: %v", err)
	}

	// 检查是否为.bz2格式的压缩文件，暂不支持
	if compressType == types.CompressTypeBz2 || compressType == types.CompressTypeBzip2 {
		return fmt.Errorf("暂不支持 %s 和 %s 格式的压缩文件", types.CompressTypeBz2.String(), types.CompressTypeBzip2.String())
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
	case types.CompressTypeZip: // Zip
		return zip.Zip(dst, src, c.config)
	case types.CompressTypeTar: // Tar
		return tar.Tar(dst, src, c.config)
	case types.CompressTypeTgz, types.CompressTypeTarGz: // Tar.gz 或 .tgz
		return tgz.Tgz(dst, src, c.config)
	case types.CompressTypeGz: // Gz
		return gzip.Gzip(dst, src, c.config)
	default:
		return fmt.Errorf("不支持的压缩格式: %s", compressType)
	}
}

// ==============================================
// 解压方法
// ==============================================

// Unpack 解压文件
//
// 参数:
//   - src: 源文件路径
//   - dst: 目标目录路径
//
// 返回:
//   - error: 错误信息
func (c *Comprx) Unpack(src string, dst string) error {
	// 检查源文件路径是否为空
	if src == "" {
		return fmt.Errorf("源文件路径不能为空")
	}

	// 智能检测压缩文件格式
	compressType, err := types.DetectCompressFormat(src)
	if err != nil {
		return fmt.Errorf("检测压缩格式失败: %v", err)
	}

	// 检查源文件是否存在
	if !utils.Exists(src) {
		return fmt.Errorf("源文件 %s 不存在", src)
	}

	// 当目标目录为空时，自动生成目标目录, 如: /path/to/file.tar.gz -> /path/to/file
	if dst == "" {
		baseName := filepath.Base(src)
		baseName = strings.TrimSuffix(baseName, ".tar.gz")
		baseName = strings.TrimSuffix(baseName, filepath.Ext(baseName))
		dst = filepath.Join(filepath.Dir(src), baseName)
	}

	// 检查目标目录是否存在, 不存在则创建
	if err := utils.EnsureDir(dst); err != nil {
		return fmt.Errorf("创建目标目录失败: %v", err)
	}

	// 根据压缩格式进行解压
	switch compressType {
	case types.CompressTypeZip: // Zip
		return zip.Unzip(src, dst, c.config)
	case types.CompressTypeTar: // Tar
		return tar.Untar(src, dst, c.config)
	case types.CompressTypeTgz: // TarGz
		return tgz.Untgz(src, dst, c.config)
	case types.CompressTypeGz: // Gzip
		return gzip.Ungzip(src, dst, c.config)
	case types.CompressTypeBz2, types.CompressTypeBzip2: // Bz2
		return bzip2.Unbz2(src, dst, c.config)
	default:
		return fmt.Errorf("不支持的压缩格式: %s", compressType)
	}
}

// ==============================================
// 列表方法
// ==============================================

//func (c *Comprx) List(archivePath string) error

// ==============================================
// 便捷函数
// ==============================================

// 便捷函数，使用默认配置
var defaultComprx = New()

// Pack 使用默认配置压缩文件
func Pack(dst string, src string) error {
	return defaultComprx.Pack(dst, src)
}

// Unpack 使用默认配置解压文件
func Unpack(src string, dst string) error {
	return defaultComprx.Unpack(src, dst)
}

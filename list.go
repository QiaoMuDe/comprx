package comprx

import (
	"fmt"

	"gitee.com/MM-Q/comprx/internal/bzip2"
	"gitee.com/MM-Q/comprx/internal/gzip"
	"gitee.com/MM-Q/comprx/internal/tar"
	"gitee.com/MM-Q/comprx/internal/tgz"
	"gitee.com/MM-Q/comprx/internal/utils"
	"gitee.com/MM-Q/comprx/internal/zip"
	"gitee.com/MM-Q/comprx/types"
)

// List 列出压缩包的所有文件信息
//
// 参数:
//   - archivePath: 压缩包文件路径
//
// 返回:
//   - *types.ArchiveInfo: 压缩包信息
//   - error: 错误信息
func List(archivePath string) (*types.ArchiveInfo, error) {
	// 智能检测压缩文件格式
	compressType, err := types.DetectCompressFormat(archivePath)
	if err != nil {
		return nil, fmt.Errorf("检测压缩格式失败: %v", err)
	}

	// 检查源文件是否存在
	if !utils.Exists(archivePath) {
		return nil, fmt.Errorf("压缩包文件 %s 不存在", archivePath)
	}

	// 根据压缩格式调用对应的列表函数
	switch compressType {
	case types.CompressTypeZip: // Zip
		return zip.ListZip(archivePath)
	case types.CompressTypeTar: // Tar
		return tar.ListTar(archivePath)
	case types.CompressTypeTgz, types.CompressTypeTarGz: // Tar.gz 或 .tgz
		return tgz.ListTgz(archivePath)
	case types.CompressTypeGz: // Gz
		return gzip.ListGzip(archivePath)
	case types.CompressTypeBz2, types.CompressTypeBzip2: // Bz2
		return bzip2.ListBz2(archivePath)
	default:
		return nil, fmt.Errorf("不支持的压缩格式: %s", compressType)
	}
}

// ListLimit 列出指定数量的文件信息
//
// 参数:
//   - archivePath: 压缩包文件路径
//   - limit: 限制返回的文件数量
//
// 返回:
//   - *types.ArchiveInfo: 压缩包信息
//   - error: 错误信息
func ListLimit(archivePath string, limit int) (*types.ArchiveInfo, error) {
	// 智能检测压缩文件格式
	compressType, err := types.DetectCompressFormat(archivePath)
	if err != nil {
		return nil, fmt.Errorf("检测压缩格式失败: %v", err)
	}

	// 检查源文件是否存在
	if !utils.Exists(archivePath) {
		return nil, fmt.Errorf("压缩包文件 %s 不存在", archivePath)
	}

	// 根据压缩格式调用对应的列表函数
	switch compressType {
	case types.CompressTypeZip: // Zip
		return zip.ListZipLimit(archivePath, limit)
	case types.CompressTypeTar: // Tar
		return tar.ListTarLimit(archivePath, limit)
	case types.CompressTypeTgz, types.CompressTypeTarGz: // Tar.gz 或 .tgz
		return tgz.ListTgzLimit(archivePath, limit)
	case types.CompressTypeGz: // Gz
		return gzip.ListGzipLimit(archivePath, limit)
	case types.CompressTypeBz2, types.CompressTypeBzip2: // Bz2
		return bzip2.ListBz2Limit(archivePath, limit)
	default:
		return nil, fmt.Errorf("不支持的压缩格式: %s", compressType)
	}
}

// ListMatch 列出匹配指定模式的文件信息
//
// 参数:
//   - archivePath: 压缩包文件路径
//   - pattern: 文件名匹配模式 (支持通配符 * 和 ?)
//
// 返回:
//   - *types.ArchiveInfo: 压缩包信息
//   - error: 错误信息
func ListMatch(archivePath string, pattern string) (*types.ArchiveInfo, error) {
	// 智能检测压缩文件格式
	compressType, err := types.DetectCompressFormat(archivePath)
	if err != nil {
		return nil, fmt.Errorf("检测压缩格式失败: %v", err)
	}

	// 检查源文件是否存在
	if !utils.Exists(archivePath) {
		return nil, fmt.Errorf("压缩包文件 %s 不存在", archivePath)
	}

	// 根据压缩格式调用对应的列表函数
	switch compressType {
	case types.CompressTypeZip: // Zip
		return zip.ListZipMatch(archivePath, pattern)
	case types.CompressTypeTar: // Tar
		return tar.ListTarMatch(archivePath, pattern)
	case types.CompressTypeTgz, types.CompressTypeTarGz: // Tar.gz 或 .tgz
		return tgz.ListTgzMatch(archivePath, pattern)
	case types.CompressTypeGz: // Gz
		return gzip.ListGzipMatch(archivePath, pattern)
	case types.CompressTypeBz2, types.CompressTypeBzip2: // Bz2
		return bzip2.ListBz2Match(archivePath, pattern)
	default:
		return nil, fmt.Errorf("不支持的压缩格式: %s", compressType)
	}
}

// PrintList 打印压缩包的所有文件信息
//
// 参数:
//   - archivePath: 压缩包文件路径
//
// 返回:
//   - error: 错误信息
func PrintList(archivePath string) error {
	archiveInfo, err := List(archivePath)
	if err != nil {
		return err
	}

	utils.PrintArchiveInfo(archiveInfo, true)
	return nil
}

// PrintListLimit 打印指定数量的文件信息
//
// 参数:
//   - archivePath: 压缩包文件路径
//   - limit: 限制打印的文件数量
//
// 返回:
//   - error: 错误信息
func PrintListLimit(archivePath string, limit int) error {
	archiveInfo, err := ListLimit(archivePath, limit)
	if err != nil {
		return err
	}

	utils.PrintArchiveInfo(archiveInfo, true)
	return nil
}

// PrintListMatch 打印匹配指定模式的文件信息
//
// 参数:
//   - archivePath: 压缩包文件路径
//   - pattern: 文件名匹配模式 (支持通配符 * 和 ?)
//
// 返回:
//   - error: 错误信息
func PrintListMatch(archivePath string, pattern string) error {
	archiveInfo, err := ListMatch(archivePath, pattern)
	if err != nil {
		return err
	}

	utils.PrintArchiveInfo(archiveInfo, true)
	return nil
}

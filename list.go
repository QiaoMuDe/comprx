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

// ==============================================
// 压缩包信息获取方法
// ==============================================

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

// ==============================================
// 打印压缩包本身信息
// ==============================================

// PrintArchiveInfo 打印压缩包本身的基本信息
//
// 参数:
//   - archivePath: 压缩包文件路径
//
// 返回:
//   - error: 错误信息
func PrintArchiveInfo(archivePath string) error {
	archiveInfo, err := List(archivePath)
	if err != nil {
		return err
	}

	utils.PrintArchiveSummary(archiveInfo)
	return nil
}

// ==============================================
// 打印压缩包内文件信息
// ==============================================

// PrintFiles 打印压缩包内所有文件信息
//
// 参数:
//   - archivePath: 压缩包文件路径
//   - detailed: true=详细样式, false=简洁样式(默认)
//
// 返回:
//   - error: 错误信息
func PrintFiles(archivePath string, detailed bool) error {
	archiveInfo, err := List(archivePath)
	if err != nil {
		return err
	}

	utils.PrintFileList(archiveInfo.Files, detailed)
	return nil
}

// PrintFilesLimit 打印压缩包内指定数量的文件信息
//
// 参数:
//   - archivePath: 压缩包文件路径
//   - limit: 限制打印的文件数量
//   - detailed: true=详细样式, false=简洁样式(默认)
//
// 返回:
//   - error: 错误信息
func PrintFilesLimit(archivePath string, limit int, detailed bool) error {
	archiveInfo, err := ListLimit(archivePath, limit)
	if err != nil {
		return err
	}

	utils.PrintFileList(archiveInfo.Files, detailed)
	return nil
}

// PrintFilesMatch 打印压缩包内匹配指定模式的文件信息
//
// 参数:
//   - archivePath: 压缩包文件路径
//   - pattern: 文件名匹配模式 (支持通配符 * 和 ?)
//   - detailed: true=详细样式, false=简洁样式(默认)
//
// 返回:
//   - error: 错误信息
func PrintFilesMatch(archivePath string, pattern string, detailed bool) error {
	archiveInfo, err := ListMatch(archivePath, pattern)
	if err != nil {
		return err
	}

	utils.PrintFileList(archiveInfo.Files, detailed)
	return nil
}

// ==============================================
// 便捷函数 - 简洁样式
// ==============================================

// PrintLs 打印压缩包内所有文件信息（简洁样式）
//
// 参数:
//   - archivePath: 压缩包文件路径
//
// 返回:
//   - error: 错误信息
func PrintLs(archivePath string) error {
	return PrintFiles(archivePath, false)
}

// PrintLsLimit 打印压缩包内指定数量的文件信息（简洁样式）
//
// 参数:
//   - archivePath: 压缩包文件路径
//   - limit: 限制打印的文件数量
//
// 返回:
//   - error: 错误信息
func PrintLsLimit(archivePath string, limit int) error {
	return PrintFilesLimit(archivePath, limit, false)
}

// PrintLsMatch 打印压缩包内匹配指定模式的文件信息（简洁样式）
//
// 参数:
//   - archivePath: 压缩包文件路径
//   - pattern: 文件名匹配模式 (支持通配符 * 和 ?)
//
// 返回:
//   - error: 错误信息
func PrintLsMatch(archivePath string, pattern string) error {
	return PrintFilesMatch(archivePath, pattern, false)
}

// ==============================================
// 便捷函数 - 详细样式
// ==============================================

// PrintLl 打印压缩包内所有文件信息（详细样式）
//
// 参数:
//   - archivePath: 压缩包文件路径
//
// 返回:
//   - error: 错误信息
func PrintLl(archivePath string) error {
	return PrintFiles(archivePath, true)
}

// PrintLlLimit 打印压缩包内指定数量的文件信息（详细样式）
//
// 参数:
//   - archivePath: 压缩包文件路径
//   - limit: 限制打印的文件数量
//
// 返回:
//   - error: 错误信息
func PrintLlLimit(archivePath string, limit int) error {
	return PrintFilesLimit(archivePath, limit, true)
}

// PrintLlMatch 打印压缩包内匹配指定模式的文件信息（详细样式）
//
// 参数:
//   - archivePath: 压缩包文件路径
//   - pattern: 文件名匹配模式 (支持通配符 * 和 ?)
//
// 返回:
//   - error: 错误信息
func PrintLlMatch(archivePath string, pattern string) error {
	return PrintFilesMatch(archivePath, pattern, true)
}

// ==============================================
// 打印压缩包信息+文件信息
// ==============================================

// PrintArchiveAndFiles 打印压缩包信息和所有文件信息
//
// 参数:
//   - archivePath: 压缩包文件路径
//   - detailed: true=详细样式, false=简洁样式(默认)
//
// 返回:
//   - error: 错误信息
func PrintArchiveAndFiles(archivePath string, detailed bool) error {
	archiveInfo, err := List(archivePath)
	if err != nil {
		return err
	}

	utils.PrintArchiveSummary(archiveInfo)
	utils.PrintFileList(archiveInfo.Files, detailed)
	return nil
}

// PrintArchiveAndFilesLimit 打印压缩包信息和指定数量的文件信息
//
// 参数:
//   - archivePath: 压缩包文件路径
//   - limit: 限制打印的文件数量
//   - detailed: true=详细样式, false=简洁样式(默认)
//
// 返回:
//   - error: 错误信息
func PrintArchiveAndFilesLimit(archivePath string, limit int, detailed bool) error {
	archiveInfo, err := ListLimit(archivePath, limit)
	if err != nil {
		return err
	}

	utils.PrintArchiveSummary(archiveInfo)
	utils.PrintFileList(archiveInfo.Files, detailed)
	return nil
}

// PrintArchiveAndFilesMatch 打印压缩包信息和匹配指定模式的文件信息
//
// 参数:
//   - archivePath: 压缩包文件路径
//   - pattern: 文件名匹配模式 (支持通配符 * 和 ?)
//   - detailed: true=详细样式, false=简洁样式(默认)
//
// 返回:
//   - error: 错误信息
func PrintArchiveAndFilesMatch(archivePath string, pattern string, detailed bool) error {
	archiveInfo, err := ListMatch(archivePath, pattern)
	if err != nil {
		return err
	}

	utils.PrintArchiveSummary(archiveInfo)
	utils.PrintFileList(archiveInfo.Files, detailed)
	return nil
}

// ==============================================
// 便捷函数 - 压缩包信息+文件信息（简洁样式）
// ==============================================

// PrintInfo 打印压缩包信息和所有文件信息（简洁样式）
//
// 参数:
//   - archivePath: 压缩包文件路径
//
// 返回:
//   - error: 错误信息
func PrintInfo(archivePath string) error {
	return PrintArchiveAndFiles(archivePath, false)
}

// PrintInfoLimit 打印压缩包信息和指定数量的文件信息（简洁样式）
//
// 参数:
//   - archivePath: 压缩包文件路径
//   - limit: 限制打印的文件数量
//
// 返回:
//   - error: 错误信息
func PrintInfoLimit(archivePath string, limit int) error {
	return PrintArchiveAndFilesLimit(archivePath, limit, false)
}

// PrintInfoMatch 打印压缩包信息和匹配指定模式的文件信息（简洁样式）
//
// 参数:
//   - archivePath: 压缩包文件路径
//   - pattern: 文件名匹配模式 (支持通配符 * 和 ?)
//
// 返回:
//   - error: 错误信息
func PrintInfoMatch(archivePath string, pattern string) error {
	return PrintArchiveAndFilesMatch(archivePath, pattern, false)
}

// ==============================================
// 便捷函数 - 压缩包信息+文件信息（详细样式）
// ==============================================

// PrintInfoDetailed 打印压缩包信息和所有文件信息（详细样式）
//
// 参数:
//   - archivePath: 压缩包文件路径
//
// 返回:
//   - error: 错误信息
func PrintInfoDetailed(archivePath string) error {
	return PrintArchiveAndFiles(archivePath, true)
}

// PrintInfoDetailedLimit 打印压缩包信息和指定数量的文件信息（详细样式）
//
// 参数:
//   - archivePath: 压缩包文件路径
//   - limit: 限制打印的文件数量
//
// 返回:
//   - error: 错误信息
func PrintInfoDetailedLimit(archivePath string, limit int) error {
	return PrintArchiveAndFilesLimit(archivePath, limit, true)
}

// PrintInfoDetailedMatch 打印压缩包信息和匹配指定模式的文件信息（详细样式）
//
// 参数:
//   - archivePath: 压缩包文件路径
//   - pattern: 文件名匹配模式 (支持通配符 * 和 ?)
//
// 返回:
//   - error: 错误信息
func PrintInfoDetailedMatch(archivePath string, pattern string) error {
	return PrintArchiveAndFilesMatch(archivePath, pattern, true)
}

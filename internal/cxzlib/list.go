package cxzlib

import (
	"compress/zlib"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"gitee.com/MM-Q/comprx/internal/utils"
	"gitee.com/MM-Q/comprx/types"
)

// ListZlib 获取ZLIB压缩包的文件信息
func ListZlib(archivePath string) (*types.ArchiveInfo, error) {
	// 确保路径为绝对路径
	absPath, err := utils.EnsureAbsPath(archivePath, "ZLIB文件路径")
	if err != nil {
		return nil, err
	}

	// 根据文件名检测压缩格式类型
	compressType, err := types.DetectCompressFormat(absPath)
	if err != nil {
		return nil, fmt.Errorf("检测压缩格式失败: %w", err)
	}

	// 打开ZLIB文件
	file, err := os.Open(absPath)
	if err != nil {
		return nil, fmt.Errorf("打开ZLIB文件失败: %w", err)
	}
	defer func() { _ = file.Close() }()

	// 获取压缩包文件信息
	stat, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("获取ZLIB文件信息失败: %w", err)
	}

	// 创建ZLIB读取器
	zlibReader, err := zlib.NewReader(file)
	if err != nil {
		return nil, fmt.Errorf("创建ZLIB读取器失败: %w", err)
	}
	defer func() { _ = zlibReader.Close() }()

	// 获取原始文件名（ZLIB格式没有文件名信息，从压缩包文件名推导）
	baseName := filepath.Base(absPath)
	var originalName string
	if ext := filepath.Ext(baseName); ext == ".zlib" {
		// 去除.zlib后缀
		originalName = strings.TrimSuffix(baseName, ".zlib")
	} else {
		originalName = baseName + utils.DecompressedSuffix
	}

	// ZLIB是单文件压缩，需要读取整个文件来获取原始大小
	// 使用io.CopyBuffer配合io.Discard，既高效又准确
	buffer := utils.GetBuffer(utils.DefaultBufferSize)
	defer utils.PutBuffer(buffer)

	originalSize, err := io.CopyBuffer(io.Discard, zlibReader, buffer)
	if err != nil {
		// 如果读取失败，使用压缩文件大小作为估算
		originalSize = stat.Size()
	}

	// 创建FileInfo
	fileInfo := types.FileInfo{
		Name:           originalName,
		Size:           originalSize,
		CompressedSize: stat.Size(),
		ModTime:        stat.ModTime(),        // ZLIB不保存修改时间，使用压缩文件的修改时间
		Mode:           utils.DefaultFileMode, // ZLIB不保存文件权限，使用默认权限
		IsDir:          false,
		IsSymlink:      false,
	}

	// 创建ArchiveInfo
	archiveInfo := &types.ArchiveInfo{
		Type:           compressType,               // 类型
		TotalFiles:     1,                          // 文件数量
		TotalSize:      originalSize,               // 原始文件大小
		CompressedSize: stat.Size(),                // 压缩文件大小
		Files:          []types.FileInfo{fileInfo}, // 文件列表
	}

	return archiveInfo, nil
}

// ListZlibLimit 获取ZLIB压缩包指定数量的文件信息
func ListZlibLimit(archivePath string, limit int) (*types.ArchiveInfo, error) {
	archiveInfo, err := ListZlib(archivePath)
	if err != nil {
		return nil, err
	}

	// ZLIB只有一个文件，limit不影响结果
	return archiveInfo, nil
}

// ListZlibMatch 获取ZLIB压缩包中匹配指定模式的文件信息
func ListZlibMatch(archivePath string, pattern string) (*types.ArchiveInfo, error) {
	archiveInfo, err := ListZlib(archivePath)
	if err != nil {
		return nil, err
	}

	// 检查单个文件是否匹配模式
	if len(archiveInfo.Files) > 0 && utils.MatchPattern(archiveInfo.Files[0].Name, pattern) {
		return archiveInfo, nil
	}

	// 如果不匹配，返回空列表
	archiveInfo.Files = []types.FileInfo{}
	archiveInfo.TotalFiles = 0
	archiveInfo.TotalSize = 0

	return archiveInfo, nil
}

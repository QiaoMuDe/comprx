package cxbzip2

import (
	"compress/bzip2"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"gitee.com/MM-Q/comprx/internal/utils"
	"gitee.com/MM-Q/comprx/types"
)

// ListBz2 获取BZ2压缩包的文件信息
func ListBz2(archivePath string) (*types.ArchiveInfo, error) {
	// 确保路径为绝对路径
	absPath, err := utils.EnsureAbsPath(archivePath, "BZ2文件路径")
	if err != nil {
		return nil, err
	}

	// 打开BZ2文件
	file, err := os.Open(absPath)
	if err != nil {
		return nil, fmt.Errorf("打开BZ2文件失败: %w", err)
	}
	defer func() { _ = file.Close() }()

	// 获取压缩包文件信息
	stat, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("获取BZ2文件信息失败: %w", err)
	}

	// 创建BZ2读取器
	bz2Reader := bzip2.NewReader(file)

	// 获取原始文件名
	baseName := filepath.Base(absPath)
	var originalName string
	if ext := filepath.Ext(baseName); ext == ".bz2" {
		originalName = baseName[:len(baseName)-len(ext)]
	} else if ext == ".bzip2" {
		originalName = baseName[:len(baseName)-len(ext)]
	} else {
		originalName = baseName + utils.DecompressedSuffix
	}

	// BZ2是单文件压缩，需要读取整个文件来获取原始大小
	// 为了避免读取大文件，我们使用一个估算方法
	var originalSize int64
	buffer := make([]byte, utils.DefaultBufferSize)
	for {
		n, err := bz2Reader.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			}
			// 如果读取失败，使用压缩文件大小作为估算
			originalSize = stat.Size()
			break
		}
		originalSize += int64(n)
	}

	// 创建BZ2文件信息
	fileInfo := types.FileInfo{
		Name:           originalName,
		Size:           originalSize,
		CompressedSize: stat.Size(),
		ModTime:        stat.ModTime(),
		Mode:           utils.DefaultFileMode, // BZ2不保存文件权限，使用默认权限
		IsDir:          false,
		IsSymlink:      false,
	}

	// 根据文件名检测压缩格式类型
	compressType, err := types.DetectCompressFormat(absPath)
	if err != nil {
		return nil, fmt.Errorf("检测压缩格式失败: %w", err)
	}

	// 创建BZ2文件信息
	archiveInfo := &types.ArchiveInfo{
		Type:           compressType,
		TotalFiles:     1,
		TotalSize:      originalSize,
		CompressedSize: stat.Size(),
		Files:          []types.FileInfo{fileInfo},
	}

	return archiveInfo, nil
}

// ListBz2Limit 获取BZ2压缩包指定数量的文件信息
func ListBz2Limit(archivePath string, limit int) (*types.ArchiveInfo, error) {
	archiveInfo, err := ListBz2(archivePath)
	if err != nil {
		return nil, err
	}

	// BZ2只有一个文件，limit不影响结果
	return archiveInfo, nil
}

// ListBz2Match 获取BZ2压缩包中匹配指定模式的文件信息
func ListBz2Match(archivePath string, pattern string) (*types.ArchiveInfo, error) {
	archiveInfo, err := ListBz2(archivePath)
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

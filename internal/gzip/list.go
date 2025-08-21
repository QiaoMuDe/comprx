package gzip

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"gitee.com/MM-Q/comprx/internal/utils"
	"gitee.com/MM-Q/comprx/types"
)

// ListGzip 获取GZIP压缩包的文件信息
func ListGzip(archivePath string) (*types.ArchiveInfo, error) {
	// 确保路径为绝对路径
	absPath, err := utils.EnsureAbsPath(archivePath, "GZIP文件路径")
	if err != nil {
		return nil, err
	}

	// 打开GZIP文件
	file, err := os.Open(absPath)
	if err != nil {
		return nil, fmt.Errorf("打开GZIP文件失败: %w", err)
	}
	defer func() { _ = file.Close() }()

	// 获取压缩包文件信息
	stat, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("获取GZIP文件信息失败: %w", err)
	}

	// 创建GZIP读取器
	gzipReader, err := gzip.NewReader(file)
	if err != nil {
		return nil, fmt.Errorf("创建GZIP读取器失败: %w", err)
	}
	defer func() { _ = gzipReader.Close() }()

	// 获取原始文件名
	originalName := gzipReader.Name
	if originalName == "" {
		// 如果GZIP头中没有文件名，从压缩包文件名推导
		baseName := filepath.Base(absPath)
		if ext := filepath.Ext(baseName); ext == ".gz" {
			//originalName = baseName[:len(baseName)-len(ext)]
			// 去除.gz后缀
			originalName = strings.TrimSuffix(baseName, ".gz")
		} else {
			originalName = baseName + utils.DecompressedSuffix
		}
	}

	// GZIP是单文件压缩，需要读取整个文件来获取原始大小
	// 为了避免读取大文件，我们使用一个估算方法
	// 实际应用中可以考虑只读取部分数据或使用其他方法
	var originalSize int64
	buffer := make([]byte, utils.DefaultBufferSize)
	for {
		n, readErr := gzipReader.Read(buffer)
		if readErr != nil {
			if readErr == io.EOF {
				break
			}
			// 如果读取失败，使用压缩文件大小作为估算
			originalSize = stat.Size()
			break
		}
		originalSize += int64(n)
	}

	// 根据文件名检测压缩格式类型
	compressType, err := types.DetectCompressFormat(absPath)
	if err != nil {
		return nil, fmt.Errorf("检测压缩格式失败: %w", err)
	}

	// 创建FileInfo
	fileInfo := types.FileInfo{
		Name:           originalName,
		Size:           originalSize,
		CompressedSize: stat.Size(),
		ModTime:        gzipReader.ModTime,
		Mode:           utils.DefaultFileMode, // GZIP不保存文件权限，使用默认权限
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

// ListGzipLimit 获取GZIP压缩包指定数量的文件信息
func ListGzipLimit(archivePath string, limit int) (*types.ArchiveInfo, error) {
	archiveInfo, err := ListGzip(archivePath)
	if err != nil {
		return nil, err
	}

	// GZIP只有一个文件，limit不影响结果
	return archiveInfo, nil
}

// ListGzipMatch 获取GZIP压缩包中匹配指定模式的文件信息
func ListGzipMatch(archivePath string, pattern string) (*types.ArchiveInfo, error) {
	archiveInfo, err := ListGzip(archivePath)
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

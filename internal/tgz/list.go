package tgz

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"

	"gitee.com/MM-Q/comprx/internal/utils"
	"gitee.com/MM-Q/comprx/types"
)

// ListTgz 获取TGZ压缩包的所有文件信息
func ListTgz(archivePath string) (*types.ArchiveInfo, error) {
	// 确保路径为绝对路径
	absPath, err := utils.EnsureAbsPath(archivePath, "TGZ文件路径")
	if err != nil {
		return nil, err
	}

	// 打开TGZ文件
	file, err := os.Open(absPath)
	if err != nil {
		return nil, fmt.Errorf("打开TGZ文件失败: %w", err)
	}
	defer func() { _ = file.Close() }()

	// 获取压缩包文件信息
	stat, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("获取TGZ文件信息失败: %w", err)
	}

	// 创建GZIP读取器
	gzipReader, err := gzip.NewReader(file)
	if err != nil {
		return nil, fmt.Errorf("创建GZIP读取器失败: %w", err)
	}
	defer func() { _ = gzipReader.Close() }()

	// 创建TAR读取器
	tarReader := tar.NewReader(gzipReader)

	// 根据文件名检测压缩格式类型
	compressType, err := types.DetectCompressFormat(absPath)
	if err != nil {
		return nil, fmt.Errorf("检测压缩格式失败: %w", err)
	}

	// 创建 ArchiveInfo 结构体
	archiveInfo := &types.ArchiveInfo{
		Type:           compressType,
		CompressedSize: stat.Size(),
		Files:          make([]types.FileInfo, 0, utils.DefaultFileCapacity),
	}

	// 遍历TAR文件中的每个条目
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("读取TGZ条目失败: %w", err)
		}

		fileInfo := types.FileInfo{
			Name:           header.Name,
			Size:           header.Size,
			CompressedSize: 0, // TGZ整体压缩，单个文件压缩大小无法准确计算
			ModTime:        header.ModTime,
			Mode:           header.FileInfo().Mode(),
			IsDir:          header.FileInfo().IsDir(),
			IsSymlink:      header.Typeflag == tar.TypeSymlink || header.Typeflag == tar.TypeLink,
		}

		// 如果是符号链接，设置链接目标
		if fileInfo.IsSymlink {
			fileInfo.LinkTarget = header.Linkname
		}

		archiveInfo.Files = append(archiveInfo.Files, fileInfo)
		archiveInfo.TotalSize += fileInfo.Size
		archiveInfo.TotalFiles++
	}

	return archiveInfo, nil
}

// ListTgzLimit 获取TGZ压缩包指定数量的文件信息
func ListTgzLimit(archivePath string, limit int) (*types.ArchiveInfo, error) {
	// 确保路径为绝对路径
	absPath, err := utils.EnsureAbsPath(archivePath, "TGZ文件路径")
	if err != nil {
		return nil, err
	}

	// 打开TGZ文件
	file, err := os.Open(absPath)
	if err != nil {
		return nil, fmt.Errorf("打开TGZ文件失败: %w", err)
	}
	defer func() { _ = file.Close() }()

	// 获取压缩包文件信息
	stat, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("获取TGZ文件信息失败: %w", err)
	}

	// 创建GZIP读取器
	gzipReader, err := gzip.NewReader(file)
	if err != nil {
		return nil, fmt.Errorf("创建GZIP读取器失败: %w", err)
	}
	defer func() { _ = gzipReader.Close() }()

	// 创建TAR读取器
	tarReader := tar.NewReader(gzipReader)

	// 根据文件名检测压缩格式类型
	compressType, err := types.DetectCompressFormat(absPath)
	if err != nil {
		return nil, fmt.Errorf("检测压缩格式失败: %w", err)
	}

	// 创建 ArchiveInfo 结构体
	archiveInfo := &types.ArchiveInfo{
		Type:           compressType,
		CompressedSize: stat.Size(),
		Files:          make([]types.FileInfo, 0, utils.DefaultFileCapacity),
	}

	// 遍历TAR文件中的每个条目，但限制数量
	count := 0
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("读取TGZ条目失败: %w", err)
		}

		// 达到限制数量就提前退出
		if limit > 0 && count >= limit {
			break
		}

		fileInfo := types.FileInfo{
			Name:           header.Name,
			Size:           header.Size,
			CompressedSize: 0, // TGZ整体压缩，单个文件压缩大小无法准确计算
			ModTime:        header.ModTime,
			Mode:           header.FileInfo().Mode(),
			IsDir:          header.FileInfo().IsDir(),
			IsSymlink:      header.Typeflag == tar.TypeSymlink || header.Typeflag == tar.TypeLink,
		}

		// 如果是符号链接，设置链接目标
		if fileInfo.IsSymlink {
			fileInfo.LinkTarget = header.Linkname
		}

		archiveInfo.Files = append(archiveInfo.Files, fileInfo)
		archiveInfo.TotalSize += fileInfo.Size
		count++
	}

	archiveInfo.TotalFiles = count
	return archiveInfo, nil
}

// ListTgzMatch 获取TGZ压缩包中匹配指定模式的文件信息
func ListTgzMatch(archivePath string, pattern string) (*types.ArchiveInfo, error) {
	archiveInfo, err := ListTgz(archivePath)
	if err != nil {
		return nil, err
	}

	archiveInfo.Files = utils.FilterFilesByPattern(archiveInfo.Files, pattern)
	archiveInfo.TotalFiles = len(archiveInfo.Files)

	// 重新计算总大小
	var totalSize int64
	for _, file := range archiveInfo.Files {
		totalSize += file.Size
	}
	archiveInfo.TotalSize = totalSize

	return archiveInfo, nil
}

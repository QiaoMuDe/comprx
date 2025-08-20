package tar

import (
	"archive/tar"
	"fmt"
	"io"
	"os"

	"gitee.com/MM-Q/comprx/internal/utils"
	"gitee.com/MM-Q/comprx/types"
)

// ListTar 获取TAR压缩包的所有文件信息
func ListTar(archivePath string) (*types.ArchiveInfo, error) {
	// 确保路径为绝对路径
	absPath, err := utils.EnsureAbsPath(archivePath, "TAR文件路径")
	if err != nil {
		return nil, err
	}

	// 打开TAR文件
	file, err := os.Open(absPath)
	if err != nil {
		return nil, fmt.Errorf("打开TAR文件失败: %w", err)
	}
	defer func() { _ = file.Close() }()

	// 获取压缩包文件信息
	stat, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("获取TAR文件信息失败: %w", err)
	}

	// 创建TAR读取器
	tarReader := tar.NewReader(file)

	archiveInfo := &types.ArchiveInfo{
		Type:           types.CompressTypeTar,
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
			return nil, fmt.Errorf("读取TAR条目失败: %w", err)
		}

		fileInfo := types.FileInfo{
			Name:           header.Name,
			Size:           header.Size,
			CompressedSize: header.Size, // TAR不压缩，压缩大小等于原始大小
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

// ListTarLimit 获取TAR压缩包指定数量的文件信息
func ListTarLimit(archivePath string, limit int) (*types.ArchiveInfo, error) {
	archiveInfo, err := ListTar(archivePath)
	if err != nil {
		return nil, err
	}

	archiveInfo.Files = utils.LimitFiles(archiveInfo.Files, limit)
	return archiveInfo, nil
}

// ListTarMatch 获取TAR压缩包中匹配指定模式的文件信息
func ListTarMatch(archivePath string, pattern string) (*types.ArchiveInfo, error) {
	archiveInfo, err := ListTar(archivePath)
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

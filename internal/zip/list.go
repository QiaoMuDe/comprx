package zip

import (
	"archive/zip"
	"fmt"
	"os"

	"gitee.com/MM-Q/comprx/internal/utils"
	"gitee.com/MM-Q/comprx/types"
)

// ListZip 获取ZIP压缩包的所有文件信息
func ListZip(archivePath string) (*types.ArchiveInfo, error) {
	// 确保路径为绝对路径
	absPath, err := utils.EnsureAbsPath(archivePath, "ZIP文件路径")
	if err != nil {
		return nil, err
	}

	// 打开ZIP文件
	reader, err := zip.OpenReader(absPath)
	if err != nil {
		return nil, fmt.Errorf("打开ZIP文件失败: %w", err)
	}
	defer func() { _ = reader.Close() }()

	// 获取压缩包文件信息
	stat, err := os.Stat(absPath)
	if err != nil {
		return nil, fmt.Errorf("获取ZIP文件信息失败: %w", err)
	}

	archiveInfo := &types.ArchiveInfo{
		Type:           types.CompressTypeZip,
		TotalFiles:     len(reader.File),
		CompressedSize: stat.Size(),
		Files:          make([]types.FileInfo, 0, len(reader.File)),
	}

	// 遍历ZIP文件中的每个条目
	for _, file := range reader.File {
		fileInfo := types.FileInfo{
			Name:           file.Name,
			Size:           int64(file.UncompressedSize64),
			CompressedSize: int64(file.CompressedSize64),
			ModTime:        file.Modified,
			Mode:           file.Mode(),
			IsDir:          file.Mode().IsDir(),
			IsSymlink:      file.Mode()&os.ModeSymlink != 0,
		}

		// 如果是符号链接，读取链接目标
		if fileInfo.IsSymlink {
			if target, err := readSymlinkTarget(file); err == nil {
				fileInfo.LinkTarget = target
			}
		}

		archiveInfo.Files = append(archiveInfo.Files, fileInfo)
		archiveInfo.TotalSize += fileInfo.Size
	}

	return archiveInfo, nil
}

// ListZipLimit 获取ZIP压缩包指定数量的文件信息
func ListZipLimit(archivePath string, limit int) (*types.ArchiveInfo, error) {
	archiveInfo, err := ListZip(archivePath)
	if err != nil {
		return nil, err
	}

	archiveInfo.Files = utils.LimitFiles(archiveInfo.Files, limit)
	return archiveInfo, nil
}

// ListZipMatch 获取ZIP压缩包中匹配指定模式的文件信息
func ListZipMatch(archivePath string, pattern string) (*types.ArchiveInfo, error) {
	archiveInfo, err := ListZip(archivePath)
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

// readSymlinkTarget 读取ZIP文件中符号链接的目标
func readSymlinkTarget(file *zip.File) (string, error) {
	reader, err := file.Open()
	if err != nil {
		return "", err
	}
	defer func() { _ = reader.Close() }()

	// 读取符号链接目标
	target := make([]byte, file.UncompressedSize64)
	n, err := reader.Read(target)
	if err != nil {
		return "", err
	}

	return string(target[:n]), nil
}

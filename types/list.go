package types

import (
	"os"
	"time"
)

// FileInfo 压缩包内文件信息
type FileInfo struct {
	Name           string      // 文件名/路径
	Size           int64       // 原始大小
	CompressedSize int64       // 压缩后大小
	ModTime        time.Time   // 修改时间
	Mode           os.FileMode // 文件权限
	IsDir          bool        // 是否为目录
	IsSymlink      bool        // 是否为符号链接
	LinkTarget     string      // 符号链接目标(如果是符号链接)
}

// ArchiveInfo 压缩包整体信息
type ArchiveInfo struct {
	Type           CompressType // 压缩包类型
	TotalFiles     int          // 总文件数
	TotalSize      int64        // 总原始大小
	CompressedSize int64        // 总压缩大小
	Files          []FileInfo   // 文件列表
}

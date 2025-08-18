package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gitee.com/MM-Q/comprx/types"
)

// FormatFileSize 格式化文件大小显示
func FormatFileSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}

// FormatFileMode 格式化文件权限显示
func FormatFileMode(mode os.FileMode) string {
	return mode.String()
}

// MatchPattern 文件名模式匹配
// 支持简单的通配符匹配: * 和 ?
func MatchPattern(name, pattern string) bool {
	if pattern == "" {
		return true
	}

	// 使用 filepath.Match 进行通配符匹配
	matched, err := filepath.Match(pattern, name)
	if err != nil {
		// 如果模式无效，尝试简单的字符串包含匹配
		return strings.Contains(strings.ToLower(name), strings.ToLower(pattern))
	}

	// 如果没有匹配到文件名，尝试匹配路径的任意部分
	if !matched {
		parts := strings.Split(name, "/")
		for _, part := range parts {
			if matched, _ := filepath.Match(pattern, part); matched {
				return true
			}
		}
		// 最后尝试字符串包含匹配
		return strings.Contains(strings.ToLower(name), strings.ToLower(pattern))
	}

	return matched
}

// PrintFileInfo 格式化打印单个文件信息
func PrintFileInfo(info types.FileInfo, showDetails bool) {
	if showDetails {
		// 详细模式：显示权限、大小、时间等
		modeStr := FormatFileMode(info.Mode)
		sizeStr := FormatFileSize(info.Size)
		timeStr := info.ModTime.Format("2006-01-02 15:04:05")

		if info.IsSymlink {
			fmt.Printf("%s %8s %s %s -> %s\n", modeStr, sizeStr, timeStr, info.Name, info.LinkTarget)
		} else {
			fmt.Printf("%s %8s %s %s\n", modeStr, sizeStr, timeStr, info.Name)
		}
	} else {
		// 简单模式：只显示文件名
		if info.IsSymlink {
			fmt.Printf("%s -> %s\n", info.Name, info.LinkTarget)
		} else {
			fmt.Printf("%s\n", info.Name)
		}
	}
}

// PrintArchiveInfo 格式化打印压缩包信息
func PrintArchiveInfo(archiveInfo *types.ArchiveInfo, showSummary bool) {
	if showSummary {
		fmt.Printf("压缩包类型: %s\n", archiveInfo.Type)
		fmt.Printf("文件总数: %d\n", archiveInfo.TotalFiles)
		fmt.Printf("原始大小: %s\n", FormatFileSize(archiveInfo.TotalSize))
		if archiveInfo.CompressedSize > 0 {
			fmt.Printf("压缩大小: %s\n", FormatFileSize(archiveInfo.CompressedSize))
			ratio := float64(archiveInfo.CompressedSize) / float64(archiveInfo.TotalSize) * 100
			fmt.Printf("压缩率: %.1f%%\n", 100-ratio)
		}
		fmt.Println(strings.Repeat("-", 50))
	}

	// 打印文件列表
	for _, file := range archiveInfo.Files {
		PrintFileInfo(file, showSummary)
	}
}

// FilterFilesByPattern 根据模式过滤文件列表
func FilterFilesByPattern(files []types.FileInfo, pattern string) []types.FileInfo {
	if pattern == "" {
		return files
	}

	var filtered []types.FileInfo
	for _, file := range files {
		if MatchPattern(file.Name, pattern) {
			filtered = append(filtered, file)
		}
	}
	return filtered
}

// LimitFiles 限制文件列表数量
func LimitFiles(files []types.FileInfo, limit int) []types.FileInfo {
	if limit <= 0 || limit >= len(files) {
		return files
	}
	return files[:limit]
}

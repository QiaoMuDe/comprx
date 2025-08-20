package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gitee.com/MM-Q/comprx/types"
)

// FormatSize 格式化文件大小显示
const (
	unit    = 1024
	unitStr = "KMGTPE"
)

// FormatFileSize 格式化文件大小显示
//
// 参数:
//   - size: 文件大小
//
// 返回:
//   - string: 格式化后的文件大小字符串
func FormatFileSize(size int64) string {
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), unitStr[exp])
}

// FormatFileMode 格式化文件权限显示
//
// 参数:
//   - mode: 文件权限
//
// 返回:
//   - string: 格式化后的文件权限字符串
func FormatFileMode(mode os.FileMode) string {
	return mode.String()
}

// MatchPattern 文件名模式匹配
// 支持简单的通配符匹配: * 和 ?
//
// 参数:
//   - name: 文件名
//   - pattern: 模式字符串
//
// 返回:
//   - bool: 是否匹配成功
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
//
// 参数:
//   - info: 文件信息
//   - showDetails: 是否显示详细信息
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

// PrintArchiveSummary 打印压缩包摘要信息
//
// 参数:
//   - archiveInfo: 压缩包信息
func PrintArchiveSummary(archiveInfo *types.ArchiveInfo) {
	fmt.Println(strings.Repeat("-", 50))                            // 分隔线
	fmt.Printf("压缩包类型: %s\n", archiveInfo.Type)                     // 压缩包类型
	fmt.Printf("文件总数: %d\n", archiveInfo.TotalFiles)                // 文件总数
	fmt.Printf("原始大小: %s\n", FormatFileSize(archiveInfo.TotalSize)) // 原始大小
	if archiveInfo.CompressedSize > 0 {
		fmt.Printf("压缩大小: %s\n", FormatFileSize(archiveInfo.CompressedSize)) // 压缩大小
		ratio := float64(archiveInfo.CompressedSize) / float64(archiveInfo.TotalSize) * 100
		fmt.Printf("压缩率: %.1f%%\n", 100-ratio) // 压缩率
	}
	fmt.Println(strings.Repeat("-", 50)) // 分隔线
}

// PrintFileList 打印文件列表
//
// 参数:
//   - files: 文件列表
//   - showDetails: 是否显示详细信息
func PrintFileList(files []types.FileInfo, showDetails bool) {
	// 遍历文件列表并打印
	for _, file := range files {
		PrintFileInfo(file, showDetails)
	}
}

// FilterFilesByPattern 根据模式过滤文件列表
//
// 参数:
//   - files: 文件列表
//   - pattern: 模式字符串
//
// 返回:
//   - []types.FileInfo: 过滤后的文件列表
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
//
// 参数:
//   - files: 文件列表
//   - limit: 限制数量
//
// 返回:
//   - []types.FileInfo: 限制后的文件列表
func LimitFiles(files []types.FileInfo, limit int) []types.FileInfo {
	if limit <= 0 || limit >= len(files) {
		return files
	}
	return files[:limit]
}

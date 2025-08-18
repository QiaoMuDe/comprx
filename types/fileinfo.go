package types

// import (
// 	"fmt"
// 	"path/filepath"
// 	"strings"
// 	"time"
// )

// // FileInfo 压缩包中的文件信息
// type FileInfo struct {
// 	Name             string    // 文件名（包含路径）
// 	Size             int64     // 原始文件大小
// 	CompressedSize   int64     // 压缩后大小
// 	ModTime          time.Time // 修改时间
// 	IsDir            bool      // 是否为目录
// 	CompressionRatio float64   // 压缩率 (0-100)
// 	Mode             string    // 文件权限模式
// }

// // GetBaseName 获取文件基本名称（不含路径）
// func (f *FileInfo) GetBaseName() string {
// 	return filepath.Base(f.Name)
// }

// // GetDir 获取文件所在目录
// func (f *FileInfo) GetDir() string {
// 	return filepath.Dir(f.Name)
// }

// // GetExtension 获取文件扩展名
// func (f *FileInfo) GetExtension() string {
// 	return filepath.Ext(f.Name)
// }

// // FormatSize 格式化文件大小为人类可读格式
// func (f *FileInfo) FormatSize() string {
// 	return formatBytes(f.Size)
// }

// // FormatCompressedSize 格式化压缩后大小为人类可读格式
// func (f *FileInfo) FormatCompressedSize() string {
// 	return formatBytes(f.CompressedSize)
// }

// // String 返回文件信息的字符串表示
// func (f *FileInfo) String() string {
// 	typeStr := "文件"
// 	if f.IsDir {
// 		typeStr = "目录"
// 	}

// 	return fmt.Sprintf("%-50s %8s %8s %6.1f%% %s %s %s",
// 		f.Name,
// 		f.FormatSize(),
// 		f.FormatCompressedSize(),
// 		f.CompressionRatio,
// 		f.ModTime.Format("2006-01-02 15:04"),
// 		f.Mode,
// 		typeStr,
// 	)
// }

// // MatchesPattern 检查文件名是否匹配指定模式
// func (f *FileInfo) MatchesPattern(pattern string) bool {
// 	fileName := f.GetBaseName()

// 	// 支持通配符匹配
// 	if matched, _ := filepath.Match(pattern, fileName); matched {
// 		return true
// 	}

// 	// 支持字符串包含匹配
// 	if strings.Contains(strings.ToLower(fileName), strings.ToLower(pattern)) {
// 		return true
// 	}

// 	// 支持路径匹配
// 	if strings.Contains(strings.ToLower(f.Name), strings.ToLower(pattern)) {
// 		return true
// 	}

// 	return false
// }

// // formatBytes 将字节数格式化为人类可读的格式
// func formatBytes(bytes int64) string {
// 	const unit = 1024
// 	if bytes < unit {
// 		return fmt.Sprintf("%d B", bytes)
// 	}

// 	div, exp := int64(unit), 0
// 	for n := bytes / unit; n >= unit; n /= unit {
// 		div *= unit
// 		exp++
// 	}

// 	units := []string{"KB", "MB", "GB", "TB", "PB"}
// 	return fmt.Sprintf("%.1f %s", float64(bytes)/float64(div), units[exp])
// }

// // CalculateCompressionRatio 计算压缩率
// func CalculateCompressionRatio(originalSize, compressedSize int64) float64 {
// 	if originalSize == 0 {
// 		return 0
// 	}
// 	return (1.0 - float64(compressedSize)/float64(originalSize)) * 100
// }

// // ListOptions 列表显示选项
// type ListOptions struct {
// 	ShowHeader     bool   // 是否显示表头
// 	ShowSize       bool   // 是否显示文件大小
// 	ShowCompressed bool   // 是否显示压缩后大小
// 	ShowRatio      bool   // 是否显示压缩率
// 	ShowTime       bool   // 是否显示修改时间
// 	ShowMode       bool   // 是否显示文件权限
// 	HumanReadable  bool   // 是否使用人类可读的大小格式
// 	SortBy         string // 排序方式: "name", "size", "time"
// 	Reverse        bool   // 是否反向排序
// }

// // DefaultListOptions 返回默认的列表选项
// func DefaultListOptions() *ListOptions {
// 	return &ListOptions{
// 		ShowHeader:     true,
// 		ShowSize:       true,
// 		ShowCompressed: true,
// 		ShowRatio:      true,
// 		ShowTime:       true,
// 		ShowMode:       false,
// 		HumanReadable:  true,
// 		SortBy:         "name",
// 		Reverse:        false,
// 	}
// }

// // ArchiveInfo 压缩包信息
// type ArchiveInfo struct {
// 	Path             string    // 压缩包路径
// 	Format           string    // 压缩格式
// 	Size             int64     // 压缩包文件大小
// 	ModTime          time.Time // 修改时间
// 	FileCount        int       // 文件数量
// 	DirCount         int       // 目录数量
// 	TotalSize        int64     // 原始总大小
// 	CompressedSize   int64     // 压缩后总大小
// 	CompressionRatio float64   // 总体压缩率
// }

// // FormatSize 格式化压缩包大小
// func (a *ArchiveInfo) FormatSize() string {
// 	return formatBytes(a.Size)
// }

// // FormatTotalSize 格式化原始总大小
// func (a *ArchiveInfo) FormatTotalSize() string {
// 	return formatBytes(a.TotalSize)
// }

// // FormatCompressedSize 格式化压缩后总大小
// func (a *ArchiveInfo) FormatCompressedSize() string {
// 	return formatBytes(a.CompressedSize)
// }

// // String 返回压缩包信息的字符串表示
// func (a *ArchiveInfo) String() string {
// 	return fmt.Sprintf("压缩包: %s\n格式: %s\n文件大小: %s\n文件数: %d\n目录数: %d\n原始大小: %s\n压缩后: %s\n压缩率: %.1f%%\n修改时间: %s",
// 		a.Path,
// 		a.Format,
// 		a.FormatSize(),
// 		a.FileCount,
// 		a.DirCount,
// 		a.FormatTotalSize(),
// 		a.FormatCompressedSize(),
// 		a.CompressionRatio,
// 		a.ModTime.Format("2006-01-02 15:04:05"),
// 	)
// }

// // PrintHeader 打印表头
// func (opts *ListOptions) PrintHeader() {
// 	if !opts.ShowHeader {
// 		return
// 	}

// 	fmt.Printf("%-50s", "文件名")
// 	if opts.ShowSize {
// 		fmt.Printf(" %8s", "大小")
// 	}
// 	if opts.ShowCompressed {
// 		fmt.Printf(" %8s", "压缩后")
// 	}
// 	if opts.ShowRatio {
// 		fmt.Printf(" %6s", "压缩率")
// 	}
// 	if opts.ShowTime {
// 		fmt.Printf(" %16s", "修改时间")
// 	}
// 	if opts.ShowMode {
// 		fmt.Printf(" %10s", "权限")
// 	}
// 	fmt.Printf(" %4s", "类型")
// 	fmt.Println()
// 	fmt.Println(strings.Repeat("-", 100))
// }

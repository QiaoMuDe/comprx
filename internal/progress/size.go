package progress

import (
	"io/fs"
	"os"
	"path/filepath"

	"gitee.com/MM-Q/comprx/types"
)

// CalculateSourceTotalSizeWithProgress 计算源路径中所有普通文件的总大小并显示进度
//
// 参数:
//   - srcPath: 源路径（文件或目录）
//   - progress: 进度显示对象
//   - scanMessage: 扫描时显示的消息，如 "正在分析内容..."
//
// 返回值:
//   - int64: 普通文件的总大小（字节）
//
// 功能:
//   - 只在进度条模式下计算总大小，文本模式返回 0
//   - 显示扫描进度条并实时更新
//   - 支持单个文件和目录的大小计算
//   - 只计算普通文件，忽略目录、符号链接等特殊文件
func CalculateSourceTotalSizeWithProgress(srcPath string, progress *Progress, scanMessage string) int64 {
	// 只在进度条模式下计算总大小
	if !progress.Enabled || progress.BarStyle == types.ProgressStyleText {
		return 0
	}

	// 开始扫描进度显示
	bar := progress.StartScan(scanMessage)
	defer func() {
		_ = progress.CloseBar(bar)
	}()

	var totalSize int64

	// 检查是文件还是目录
	info, err := os.Stat(srcPath)
	if err != nil {
		return 0
	}

	if info.Mode().IsRegular() {
		// 单个文件
		totalSize = info.Size()
		_ = bar.Add64(totalSize)
	} else if info.IsDir() {
		// 目录，需要遍历所有文件
		_ = filepath.WalkDir(srcPath, func(path string, entry fs.DirEntry, err error) error {
			if err != nil {
				return nil // 忽略错误，继续遍历
			}

			// 只计算普通文件的大小
			if entry.Type().IsRegular() {
				if fileInfo, err := entry.Info(); err == nil {
					fileSize := fileInfo.Size()
					totalSize += fileSize
					_ = bar.Add64(fileSize) // 实时更新进度条
				}
			}

			return nil
		})
	}

	return totalSize
}
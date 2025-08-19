package comprx

import (
	"os"
	"path/filepath"

	"gitee.com/MM-Q/comprx/internal/utils"
)

// GetFileSize 获取单个文件的大小（字节）
//
// 参数:
//   - filePath: 文件路径
//
// 返回:
//   - int64: 文件大小（字节）
//
// 注意:
//   - 如果文件不存在，则返回 0
//   - 如果文件不是普通文件，则返回 0
//   - 如果获取文件信息时发生错误，则返回 0
func GetFileSize(filePath string) int64 {
	info, err := os.Stat(filePath)
	if err != nil {
		return 0
	}

	// 如果不是普通文件，返回 0
	if !info.Mode().IsRegular() {
		return 0
	}

	return info.Size()
}

// GetDirectorySize 获取目录中所有文件的总大小（字节）
//
// 参数:
//   - dirPath: 目录路径
//
// 返回:
//   - int64: 目录中所有文件的总大小（字节）
//
// 注意:
//   - 如果目录不存在，则返回 0
//   - 如果目录为空，返回 0
//   - 如果目录中包含非普通文件，则忽略这些文件
//   - 如果遍历目录时发生错误，则返回 0
func GetDirectorySize(dirPath string) (totalSize int64) {
	if !utils.Exists(dirPath) {
		return 0
	}

	err := filepath.WalkDir(dirPath, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			// 忽略错误，继续遍历
			return nil
		}

		// 只计算普通文件的大小
		if entry.Type().IsRegular() {
			if info, err := entry.Info(); err == nil {
				totalSize += info.Size()
			}
		}

		return nil
	})

	// 如果遍历失败，返回 0
	if err != nil {
		return 0
	}

	return totalSize
}

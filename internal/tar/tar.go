package tar

import (
	"archive/tar"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"gitee.com/MM-Q/comprx/config"
	"gitee.com/MM-Q/comprx/internal/utils"
)

// Tar 函数用于创建TAR归档文件
//
// 参数:
//   - dst: 生成的TAR文件路径
//   - src: 需要归档的源路径
//   - config: 压缩配置指针
//
// 返回值:
//   - error: 操作过程中遇到的错误
func Tar(dst string, src string, config *config.Config) error {
	// 确保路径为绝对路径
	var absErr error
	if dst, absErr = utils.EnsureAbsPath(dst, "TAR文件路径"); absErr != nil {
		return absErr
	}
	if src, absErr = utils.EnsureAbsPath(src, "源路径"); absErr != nil {
		return absErr
	}

	// 检查目标文件是否已存在
	if _, err := os.Stat(dst); err == nil {
		// 文件已存在，检查是否允许覆盖
		if !config.OverwriteExisting {
			return fmt.Errorf("目标文件已存在且不允许覆盖: %s", dst)
		}
	}

	// 确保目标目录存在
	if err := utils.EnsureDir(filepath.Dir(dst)); err != nil {
		return fmt.Errorf("创建目标目录失败: %w", err)
	}

	// 创建 TAR 文件
	tarFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("创建 TAR 文件失败: %w", err)
	}
	defer func() { _ = tarFile.Close() }()

	// 创建 TAR 写入器
	tarWriter := tar.NewWriter(tarFile)
	defer func() { _ = tarWriter.Close() }()

	// 检查源路径是文件还是目录
	srcInfo, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("获取源路径信息失败: %w", err)
	}

	// 根据源路径类型处理
	var tarErr error
	if srcInfo.IsDir() {
		// 遍历目录并添加文件到 TAR 包
		tarErr = walkDirectoryForTar(src, tarWriter)
	} else {
		// 单文件处理逻辑
		tarErr = processRegularFile(tarWriter, src, filepath.Base(src), srcInfo)
	}

	// 检查是否有错误发生
	if tarErr != nil {
		return fmt.Errorf("打包目录到 TAR 失败: %w", tarErr)
	}

	return nil
}

// processRegularFile 处理普通文件
//
// 参数:
//   - tarWriter: *tar.Writer - TAR 文件写入器
//   - path: string - 文件路径
//   - headerName: string - TAR 文件中的文件名
//   - info: os.FileInfo - 文件信息
//
// 返回值:
//   - error - 操作过程中遇到的错误
func processRegularFile(tarWriter *tar.Writer, path, headerName string, info os.FileInfo) error {
	// 创建文件头
	header, err := tar.FileInfoHeader(info, "")
	if err != nil {
		return fmt.Errorf("处理文件 '%s' 时出错 - 创建 TAR 文件头失败: %w", path, err)
	}
	header.Name = headerName // 设置文件名

	// 写入文件头
	if err := tarWriter.WriteHeader(header); err != nil {
		return fmt.Errorf("处理文件 '%s' 时出错 - 写入 TAR 文件头失败: %w", path, err)
	}

	// 打开文件
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("处理文件 '%s' 时出错 - 打开文件失败: %w", path, err)
	}
	defer func() { _ = file.Close() }()

	// 获取文件大小
	fileSize := info.Size()

	// 获取缓冲区大小并创建缓冲区
	bufferSize := utils.GetBufferSize(fileSize)
	buffer := utils.GetBuffer(bufferSize)
	defer utils.PutBuffer(buffer)

	// 复制文件内容到TAR写入器
	if _, err := io.CopyBuffer(tarWriter, file, buffer); err != nil {
		return fmt.Errorf("处理文件 '%s' 时出错 - 写入 TAR 文件失败: %w", path, err)
	}

	return nil
}

// processDirectory 处理目录
//
// 参数:
//   - tarWriter: *tar.Writer - TAR 文件写入器
//   - headerName: string - TAR 文件中的目录名
//   - info: os.FileInfo - 目录信息
//
// 返回值:
//   - error - 操作过程中遇到的错误
func processDirectory(tarWriter *tar.Writer, headerName string, info os.FileInfo) error {
	// 创建目录文件头
	header, err := tar.FileInfoHeader(info, "")
	if err != nil {
		return fmt.Errorf("处理目录 '%s' 时出错 - 创建 TAR 文件头失败: %w", headerName, err)
	}
	// 设置目录名
	header.Name = headerName + "/" // 目录名后添加斜杠

	// 写入目录文件头
	if err := tarWriter.WriteHeader(header); err != nil {
		return fmt.Errorf("处理目录 '%s' 时出错 - 写入 TAR 目录头失败: %w", headerName, err)
	}
	return nil
}

// processSymlink 处理软链接
//
// 参数:
//   - tarWriter: *tar.Writer - TAR 文件写入器
//   - path: string - 软链接路径
//   - headerName: string - TAR 文件中的软链接名
//   - info: os.FileInfo - 文件信息
//
// 返回值:
//   - error - 操作过程中遇到的错误
func processSymlink(tarWriter *tar.Writer, path, headerName string, info os.FileInfo) error {
	// 读取软链接目标
	target, err := os.Readlink(path)
	if err != nil {
		return fmt.Errorf("处理软链接 '%s' 时出错 - 读取软链接目标失败: %w", path, err)
	}

	// 创建软链接文件头
	header, err := tar.FileInfoHeader(info, target)
	if err != nil {
		return fmt.Errorf("处理软链接 '%s' 时出错 - 创建 TAR 文件头失败: %w", path, err)
	}
	header.Name = headerName

	// 写入软链接文件头
	if err := tarWriter.WriteHeader(header); err != nil {
		return fmt.Errorf("处理软链接 '%s' 时出错 - 写入 TAR 软链接头失败: %w", path, err)
	}
	return nil
}

// processSpecialFile 处理特殊文件类型
//
// 参数:
//   - tarWriter: *tar.Writer - TAR 文件写入器
//   - headerName: string - TAR 文件中的特殊文件名
//   - info: os.FileInfo - 文件信息
//
// 返回值:
//   - error - 操作过程中遇到的错误
func processSpecialFile(tarWriter *tar.Writer, headerName string, info os.FileInfo) error {
	// 创建 TAR 文件头
	header, err := tar.FileInfoHeader(info, "")
	if err != nil {
		return fmt.Errorf("处理特殊文件 '%s' 时出错 - 创建 TAR 文件头失败: %w", headerName, err)
	}
	header.Name = headerName

	// 写入 TAR 文件头
	if err := tarWriter.WriteHeader(header); err != nil {
		return fmt.Errorf("处理特殊文件 '%s' 时出错 - 写入 TAR 特殊文件头失败: %w", headerName, err)
	}
	return nil
}

// walkDirectoryForTar 遍历目录并处理文件到TAR包
//
// 参数:
//   - src: 源目录路径
//   - tarWriter: TAR写入器
//
// 返回值:
//   - error: 遍历过程中发生的错误
func walkDirectoryForTar(src string, tarWriter *tar.Writer) error {
	return filepath.WalkDir(src, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			// 如果不存在则忽略
			if os.IsNotExist(err) {
				return nil
			}
			// 其他错误
			return fmt.Errorf("遍历路径 '%s' 时出错: %w", path, err)
		}

		// 获取相对路径，保留顶层目录
		headerName, err := filepath.Rel(filepath.Dir(src), path)
		if err != nil {
			return fmt.Errorf("处理路径 '%s' 时出错 - 获取相对路径失败: %w", path, err)
		}

		// 替换路径分隔符为正斜杠(TAR 文件格式要求)
		headerName = filepath.ToSlash(headerName)

		// 根据文件类型处理
		switch {
		case entry.Type().IsRegular(): // 处理普通文件
			info, err := entry.Info()
			if err != nil {
				return fmt.Errorf("处理文件 '%s' 时出错 - 获取文件信息失败: %w", path, err)
			}
			return processRegularFile(tarWriter, path, headerName, info)
		case entry.IsDir(): // 处理目录
			info, err := entry.Info()
			if err != nil {
				return fmt.Errorf("处理目录 '%s' 时出错 - 获取目录信息失败: %w", path, err)
			}
			return processDirectory(tarWriter, headerName, info)
		case entry.Type()&os.ModeSymlink != 0: // 处理符号链接
			info, err := entry.Info()
			if err != nil {
				return fmt.Errorf("处理符号链接 '%s' 时出错 - 获取文件信息失败: %w", path, err)
			}
			return processSymlink(tarWriter, path, headerName, info)
		default: // 处理特殊文件
			info, err := entry.Info()
			if err != nil {
				return fmt.Errorf("处理特殊文件 '%s' 时出错 - 获取文件信息失败: %w", path, err)
			}
			return processSpecialFile(tarWriter, headerName, info)
		}
	})
}

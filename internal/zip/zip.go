package zip

import (
	"archive/zip"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"gitee.com/MM-Q/comprx/config"
	"gitee.com/MM-Q/comprx/internal/utils"
)

// Zip 函数用于创建ZIP压缩文件
//
// 参数:
//   - dst: 生成的ZIP文件路径
//   - src: 需要压缩的源路径
//   - config: 压缩配置指针
//
// 返回值:
//   - error: 操作过程中遇到的错误
func Zip(dst string, src string, config *config.Config) error {
	// 确保路径为绝对路径
	var absErr error
	if dst, absErr = utils.EnsureAbsPath(dst, "ZIP文件路径"); absErr != nil {
		return absErr
	}
	if src, absErr = utils.EnsureAbsPath(src, "源路径"); absErr != nil {
		return absErr
	}

	// 创建 ZIP 文件
	zipFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("创建 ZIP 文件失败: %w", err)
	}
	defer func() { _ = zipFile.Close() }()

	// 创建 ZIP 写入器
	zipWriter := zip.NewWriter(zipFile)
	defer func() { _ = zipWriter.Close() }()

	// 检查源路径是文件还是目录
	srcInfo, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("获取源路径信息失败: %w", err)
	}

	// 根据源路径类型处理
	var zipErr error
	if srcInfo.IsDir() {
		// 遍历目录并添加文件到 ZIP 包 (使用 WalkDir 提升性能)
		zipErr = filepath.WalkDir(src, func(path string, entry fs.DirEntry, err error) error {
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

			// 替换路径分隔符为正斜杠(ZIP 文件格式要求)
			headerName = filepath.ToSlash(headerName)

			// 根据文件类型处理
			switch {
			case entry.Type().IsRegular(): // 处理普通文件
				info, err := entry.Info()
				if err != nil {
					return fmt.Errorf("处理文件 '%s' 时出错 - 获取文件信息失败: %w", path, err)
				}
				return processRegularFile(zipWriter, path, headerName, info, config)
			case entry.IsDir(): // 处理目录
				info, err := entry.Info()
				if err != nil {
					return fmt.Errorf("处理目录 '%s' 时出错 - 获取目录信息失败: %w", path, err)
				}
				return processDirectory(zipWriter, headerName, info)
			case entry.Type()&fs.ModeSymlink != 0: // 处理符号链接
				return processSymlink(zipWriter, path, headerName, entry.Type())
			default: // 处理特殊文件
				return processSpecialFile(zipWriter, headerName, entry.Type())
			}
		})
	} else {
		// 新增的单文件处理逻辑
		zipErr = processRegularFile(zipWriter, src, filepath.Base(src), srcInfo, config)
	}

	// 检查是否有错误发生
	if zipErr != nil {
		return fmt.Errorf("打包目录到 ZIP 失败: %w", zipErr)
	}

	return nil
}

// processRegularFile 处理普通文件
//
// 参数:
//   - zipWriter: *zip.Writer - ZIP 文件写入器
//   - path: string - 文件路径
//   - headerName: string - ZIP 文件中的文件名
//   - info: os.FileInfo - 文件信息
//   - config: 压缩配置
//
// 返回值:
//   - error - 操作过程中遇到的错误
func processRegularFile(zipWriter *zip.Writer, path, headerName string, info os.FileInfo, config *config.Config) error {
	// 创建文件头
	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return fmt.Errorf("处理文件 '%s' 时出错 - 创建 ZIP 文件头失败: %w", path, err)
	}
	header.Name = headerName                     // 设置文件名
	header.Method = getCompressionMethod(config) // 使用配置的压缩方法

	// 创建 ZIP 写入器
	fileWriter, err := zipWriter.CreateHeader(header)
	if err != nil {
		return fmt.Errorf("处理文件 '%s' 时出错 - 创建 ZIP 写入器失败: %w", path, err)
	}

	// 打开文件
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("处理文件 '%s' 时出错 - 打开文件失败: %w", path, err)
	}

	// 获取文件大小
	fileSize := info.Size()

	// 获取缓冲区大小并创建缓冲区
	bufferSize := utils.GetBufferSize(fileSize)
	buffer := utils.GetBuffer(bufferSize)
	defer utils.PutBuffer(buffer)

	// 复制文件内容到ZIP写入器
	_, err = io.CopyBuffer(fileWriter, file, buffer)

	// 立即关闭文件并检查错误
	closeErr := file.Close()
	if err != nil {
		return fmt.Errorf("处理文件 '%s' 时出错 - 写入 ZIP 文件失败: %w", path, err)
	}
	if closeErr != nil {
		return fmt.Errorf("处理文件 '%s' 时出错 - 关闭文件失败: %w", path, closeErr)
	}

	return nil
}

// processDirectory 处理目录
//
// 参数:
//   - zipWriter: *zip.Writer - ZIP 文件写入器
//   - headerName: string - ZIP 文件中的目录名
//   - info: os.FileInfo - 目录信息
//
// 返回值:
//   - error - 操作过程中遇到的错误
func processDirectory(zipWriter *zip.Writer, headerName string, info os.FileInfo) error {
	// 创建目录文件头
	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return fmt.Errorf("处理目录 '%s' 时出错 - 创建 ZIP 文件头失败: %w", headerName, err)
	}
	// 设置目录名
	header.Name = headerName + "/" // 目录名后添加斜杠
	header.Method = zip.Store      // 使用不压缩的方法

	// 创建目录文件头
	if _, err := zipWriter.CreateHeader(header); err != nil {
		return fmt.Errorf("处理目录 '%s' 时出错 - 创建 ZIP 目录失败: %w", headerName, err)
	}
	return nil
}

// processSymlink 处理软链接
//
// 参数:
//   - zipWriter: *zip.Writer - ZIP 文件写入器
//   - path: string - 软链接路径
//   - headerName: string - ZIP 文件中的软链接名
//   - mode: fs.FileMode - 文件模式
//
// 返回值:
//   - error - 操作过程中遇到的错误
func processSymlink(zipWriter *zip.Writer, path, headerName string, mode fs.FileMode) error {
	// 读取软链接目标
	target, err := os.Readlink(path)
	if err != nil {
		return fmt.Errorf("处理软链接 '%s' 时出错 - 读取软链接目标失败: %w", path, err)
	}

	// 创建软链接文件头
	header := &zip.FileHeader{
		Name:   headerName,
		Method: zip.Store,
	}
	header.SetMode(mode)

	// 创建软链接文件
	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return fmt.Errorf("处理软链接 '%s' 时出错 - 创建 ZIP 软链接失败: %w", path, err)
	}
	if _, err := writer.Write([]byte(target)); err != nil {
		return fmt.Errorf("处理软链接 '%s' 时出错 - 写入软链接目标失败: %w", path, err)
	}
	return nil
}

// processSpecialFile 处理特殊文件类型
//
// 参数:
//   - zipWriter: *zip.Writer - ZIP 文件写入器
//   - headerName: string - ZIP 文件中的特殊文件名
//   - mode: fs.FileMode - 文件模式
//
// 返回值:
//   - error - 操作过程中遇到的错误
func processSpecialFile(zipWriter *zip.Writer, headerName string, mode fs.FileMode) error {
	// 创建 ZIP 文件头
	header := &zip.FileHeader{
		Name:   headerName,
		Method: zip.Store,
	}
	header.SetMode(mode)

	// 创建 ZIP 文件写入器
	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return fmt.Errorf("处理特殊文件 '%s' 时出错 - 创建 ZIP 特殊文件失败: %w", headerName, err)
	}
	if _, err := writer.Write([]byte{}); err != nil {
		return fmt.Errorf("处理特殊文件 '%s' 时出错 - 写入特殊文件失败: %w", headerName, err)
	}
	return nil
}

// getCompressionMethod 根据配置返回对应的压缩方法
func getCompressionMethod(config *config.Config) uint16 {
	if config.EnableCompression {
		return zip.Deflate // 启用压缩
	}
	return zip.Store // 不压缩，只存储
}

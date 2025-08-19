package gzip

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"gitee.com/MM-Q/comprx/config"
	"gitee.com/MM-Q/comprx/internal/utils"
)

// Ungzip 解压缩 GZIP 文件
//
// 参数:
//   - gzipFilePath: 要解压缩的 GZIP 文件路径
//   - targetPath: 解压缩后的目标文件路径
//   - config: 解压缩配置
//
// 返回值:
//   - error: 解压缩过程中发生的错误
func Ungzip(gzipFilePath string, targetPath string, config *config.Config) error {
	// 确保路径为绝对路径
	var absErr error
	if gzipFilePath, absErr = utils.EnsureAbsPath(gzipFilePath, "GZIP文件路径"); absErr != nil {
		return absErr
	}
	if targetPath, absErr = utils.EnsureAbsPath(targetPath, "目标文件路径"); absErr != nil {
		return absErr
	}

	// 打开 GZIP 文件（同时检查文件是否存在）
	gzipFile, err := os.Open(gzipFilePath)
	if err != nil {
		return fmt.Errorf("打开 GZIP 文件失败: %w", err)
	}
	defer func() { _ = gzipFile.Close() }()

	// 获取GZIP文件信息用于预验证
	gzipInfo, err := gzipFile.Stat()
	if err != nil {
		return fmt.Errorf("获取GZIP文件信息失败: %w", err)
	}

	// 预验证GZIP文件大小（检查输入文件合理性）
	if config.EnableSizeCheck && gzipInfo.Size() > config.MaxTotalSize {
		return fmt.Errorf("GZIP文件大小 %s 超过处理限制 %s",
			utils.FormatFileSize(gzipInfo.Size()), utils.FormatFileSize(config.MaxTotalSize))
	}

	// 创建 GZIP 读取器
	gzipReader, err := gzip.NewReader(gzipFile)
	if err != nil {
		return fmt.Errorf("创建 GZIP 读取器失败: %w", err)
	}
	defer func() { _ = gzipReader.Close() }()

	// 检查目标路径状态，处理目录情况和覆盖检查
	if targetStat, err := os.Stat(targetPath); err == nil {
		if targetStat.IsDir() {
			// 目标是目录，生成文件名
			if gzipReader.Name != "" {
				// 直接验证 GZIP 头部的文件名，并与目标目录合并
				validatedPath, err := utils.ValidatePathSimple(targetPath, gzipReader.Name)
				if err != nil {
					return fmt.Errorf("GZIP文件头包含不安全的文件名: %w", err)
				}
				targetPath = validatedPath
			} else {
				// 如果GZIP文件头中没有原始文件名，则去掉.gz扩展名
				baseName := filepath.Base(gzipFilePath)
				baseName = strings.TrimSuffix(baseName, ".gz")
				targetPath = filepath.Join(targetPath, baseName)
			}

			// 重新检查生成的目标文件是否存在
			if _, err := os.Stat(targetPath); err == nil && !config.OverwriteExisting {
				return fmt.Errorf("目标文件已存在且不允许覆盖: %s", targetPath)
			}
		} else {
			// 目标是文件，检查是否允许覆盖
			if !config.OverwriteExisting {
				return fmt.Errorf("目标文件已存在且不允许覆盖: %s", targetPath)
			}
		}
	}

	// 检查目标文件的父目录是否存在，如果不存在则创建
	parentDir := filepath.Dir(targetPath)
	if err := utils.EnsureDir(parentDir); err != nil {
		return fmt.Errorf("创建目标文件父目录失败: %w", err)
	}

	// 创建目标文件
	targetFile, err := os.Create(targetPath)
	if err != nil {
		return fmt.Errorf("创建目标文件失败: %w", err)
	}
	defer func() { _ = targetFile.Close() }()

	// 使用之前获取的gzipInfo来估算缓冲区大小

	// 获取缓冲区大小并创建缓冲区
	bufferSize := utils.GetBufferSize(gzipInfo.Size())
	buffer := utils.GetBuffer(bufferSize)
	defer utils.PutBuffer(buffer)

	// 创建大小跟踪器（用于解压过程中的大小检查）
	tracker := utils.NewSizeTracker()

	// 创建带验证的写入器包装器
	validatingWriter := &ungzipValidatingWriter{
		writer:         targetFile,
		config:         config,
		compressedSize: gzipInfo.Size(),
		tracker:        tracker,
	}

	// 解压缩文件内容（使用带验证的写入器）
	if _, err := io.CopyBuffer(validatingWriter, gzipReader, buffer); err != nil {
		return fmt.Errorf("解压缩文件失败: %w", err)
	}

	// 如果GZIP文件头中有修改时间信息，则设置目标文件的修改时间
	if !gzipReader.ModTime.IsZero() {
		if err := os.Chtimes(targetPath, gzipReader.ModTime, gzipReader.ModTime); err != nil {
			// 设置时间失败不是致命错误，只记录警告
			fmt.Printf("警告: 设置文件修改时间失败: %v\n", err)
		}
	}

	return nil
}

// ungzipValidatingWriter 带验证功能的写入器包装器（用于GZIP解压）
type ungzipValidatingWriter struct {
	writer         io.Writer
	config         *config.Config
	compressedSize int64
	totalWritten   int64
	tracker        *utils.SizeTracker
}

// Write 实现io.Writer接口，在写入时进行安全验证
func (vw *ungzipValidatingWriter) Write(p []byte) (n int, err error) {
	// 写入数据
	n, err = vw.writer.Write(p)
	if err != nil {
		return n, err
	}

	// 更新总写入大小
	vw.totalWritten += int64(n)

	// 验证解压后的大小是否超过单文件限制
	if vw.config.EnableSizeCheck && vw.totalWritten > vw.config.MaxFileSize {
		return n, fmt.Errorf("解压后文件大小 %s 超过单文件限制 %s",
			utils.FormatFileSize(vw.totalWritten), utils.FormatFileSize(vw.config.MaxFileSize))
	}

	// 验证解压后的大小是否超过总大小限制
	if vw.config.EnableSizeCheck && vw.totalWritten > vw.config.MaxTotalSize {
		return n, fmt.Errorf("解压后文件大小 %s 超过总大小限制 %s",
			utils.FormatFileSize(vw.totalWritten), utils.FormatFileSize(vw.config.MaxTotalSize))
	}

	// 验证压缩比（防止Zip Bomb攻击）
	if err := utils.ValidateCompressionRatio(vw.config, vw.totalWritten, vw.compressedSize); err != nil {
		return n, fmt.Errorf("压缩比验证失败: %w", err)
	}

	return n, nil
}

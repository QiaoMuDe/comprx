package gzip

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"gitee.com/MM-Q/comprx/config"
	"gitee.com/MM-Q/comprx/internal/utils"
)

// Gzip 函数用于压缩单个文件为GZIP格式
//
// 参数:
//   - dst: 生成的GZIP文件路径
//   - src: 需要压缩的源文件路径
//   - cfg: 压缩配置指针
//
// 返回值:
//   - error: 操作过程中遇到的错误
func Gzip(dst string, src string, cfg *config.Config) error {
	// 确保路径为绝对路径
	var absErr error
	if dst, absErr = utils.EnsureAbsPath(dst, "GZIP文件路径"); absErr != nil {
		return absErr
	}
	if src, absErr = utils.EnsureAbsPath(src, "源文件路径"); absErr != nil {
		return absErr
	}

	// 检查源路径是否为文件
	srcInfo, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("获取源文件信息失败: %w", err)
	}

	// 检查源路径是否为目录
	if srcInfo.IsDir() {
		return fmt.Errorf("GZIP 只支持单文件压缩，不支持目录压缩")
	}

	// 预检查源文件大小
	if _, err := utils.PreCheckSingleFile(cfg, src); err != nil {
		return fmt.Errorf("源文件大小检查失败: %w", err)
	}

	// 检查目标文件是否已存在
	if _, err := os.Stat(dst); err == nil {
		// 文件已存在，检查是否允许覆盖
		if !cfg.OverwriteExisting {
			return fmt.Errorf("目标文件已存在且不允许覆盖: %s", dst)
		}
	}

	// 确保目标目录存在
	if err := utils.EnsureDir(filepath.Dir(dst)); err != nil {
		return fmt.Errorf("创建目标目录失败: %w", err)
	}

	// 创建 GZIP 文件
	gzipFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("创建 GZIP 文件失败: %w", err)
	}
	defer func() { _ = gzipFile.Close() }()

	// 创建 GZIP 写入器
	gzipWriter, err := gzip.NewWriterLevel(gzipFile, config.GetCompressionLevel(cfg))
	if err != nil {
		return fmt.Errorf("创建 GZIP 写入器失败: %w", err)
	}
	defer func() { _ = gzipWriter.Close() }()

	// 设置 GZIP 文件头信息
	gzipWriter.Name = filepath.Base(src)
	gzipWriter.ModTime = srcInfo.ModTime()

	// 打开源文件
	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("打开源文件失败: %w", err)
	}
	defer func() { _ = srcFile.Close() }()

	// 获取文件大小
	fileSize := srcInfo.Size()

	// 获取缓冲区大小并创建缓冲区
	bufferSize := utils.GetBufferSize(fileSize)
	buffer := utils.GetBuffer(bufferSize)
	defer utils.PutBuffer(buffer)

	// 创建通用的压缩验证写入器包装器
	validatingWriter := utils.NewCompressionValidatingWriter(gzipWriter, cfg)

	// 复制文件内容到GZIP写入器（使用带验证的写入器）
	if _, err := io.CopyBuffer(validatingWriter, srcFile, buffer); err != nil {
		return fmt.Errorf("压缩文件失败: %w", err)
	}

	return nil
}

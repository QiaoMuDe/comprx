package zip

import (
	"fmt"
	"io"

	"gitee.com/MM-Q/comprx/config"
	"gitee.com/MM-Q/comprx/internal/utils"
)

// zipValidatingWriter ZIP解压验证写入器包装器
//
// 专门用于ZIP解压的验证写入器，与通用解压写入器的区别：
// 1. 支持预知解压大小的验证（ZIP文件头包含未压缩大小信息）
// 2. 针对ZIP格式的特殊验证需求
type zipValidatingWriter struct {
	writer           io.Writer          // 被包装的写入器
	config           *config.Config     // 解压配置
	compressedSize   int64              // 压缩文件大小
	uncompressedSize int64              // 预期的解压大小（来自ZIP文件头）
	totalWritten     int64              // 已写入的总大小
	tracker          *utils.SizeTracker // 大小跟踪器
}

// newZipValidatingWriter 创建新的ZIP验证写入器
//
// 参数:
//   - writer: 被包装的写入器
//   - config: 解压配置
//   - uncompressedSize: 预期的解压大小（来自ZIP文件头）
//   - compressedSize: 压缩文件大小
//   - tracker: 大小跟踪器
//
// 返回值:
//   - *zipValidatingWriter: ZIP验证写入器实例
func newZipValidatingWriter(writer io.Writer, config *config.Config, uncompressedSize, compressedSize int64, tracker *utils.SizeTracker) *zipValidatingWriter {
	return &zipValidatingWriter{
		writer:           writer,
		config:           config,
		compressedSize:   compressedSize,
		uncompressedSize: uncompressedSize,
		totalWritten:     0,
		tracker:          tracker,
	}
}

// Write 实现io.Writer接口，在写入时进行安全验证
//
// 参数:
//   - p: 要写入的数据
//
// 返回值:
//   - n: 实际写入的字节数
//   - err: 写入过程中的错误
func (vw *zipValidatingWriter) Write(p []byte) (n int, err error) {
	// 预检查：验证即将写入的数据是否会超过预期大小
	if vw.uncompressedSize > 0 && vw.totalWritten+int64(len(p)) > vw.uncompressedSize {
		return 0, fmt.Errorf("写入数据将超过ZIP文件头声明的大小 %s",
			utils.FormatFileSize(vw.uncompressedSize))
	}

	// 写入数据到底层写入器
	n, err = vw.writer.Write(p)
	if err != nil {
		return n, err
	}

	// 更新总写入大小
	vw.totalWritten += int64(n)

	// 验证解压后的大小是否超过预期（基于ZIP文件头信息）
	if vw.uncompressedSize > 0 && vw.totalWritten > vw.uncompressedSize {
		return n, fmt.Errorf("解压数据大小 %s 超过ZIP文件头声明的大小 %s",
			utils.FormatFileSize(vw.totalWritten), utils.FormatFileSize(vw.uncompressedSize))
	}

	// 验证解压后的大小是否超过单文件限制
	if vw.config.EnableSizeCheck && vw.totalWritten > vw.config.MaxFileSize {
		return n, fmt.Errorf("解压后文件大小 %s 超过单文件限制 %s",
			utils.FormatFileSize(vw.totalWritten), utils.FormatFileSize(vw.config.MaxFileSize))
	}

	return n, nil
}

// GetTotalWritten 获取已写入的总大小
//
// 返回值:
//   - int64: 已写入的总大小（字节）
func (vw *zipValidatingWriter) GetTotalWritten() int64 {
	return vw.totalWritten
}

// GetUncompressedSize 获取预期的解压大小
//
// 返回值:
//   - int64: 预期的解压大小（字节）
func (vw *zipValidatingWriter) GetUncompressedSize() int64 {
	return vw.uncompressedSize
}

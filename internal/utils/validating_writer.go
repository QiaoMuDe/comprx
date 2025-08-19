package utils

import (
	"fmt"
	"io"

	"gitee.com/MM-Q/comprx/config"
)

// ValidatingWriter 验证写入器包装器
//
// 在压缩/解压过程中进行安全验证：
// - 监测写入文件大小
// - 验证是否超过配置的大小限制
type ValidatingWriter struct {
	writer       io.Writer      // 被包装的写入器
	config       *config.Config // 配置
	totalWritten int64          // 已写入的总大小
	errorPrefix  string         // 错误消息前缀（如"压缩后"、"解压后"）
	tracker      *SizeTracker   // 可选的大小跟踪器
}

// NewCompressionValidatingWriter 创建压缩验证写入器
func NewCompressionValidatingWriter(writer io.Writer, config *config.Config) *ValidatingWriter {
	return &ValidatingWriter{
		writer:       writer,
		config:       config,
		totalWritten: 0,
		errorPrefix:  "压缩后",
	}
}

// NewDecompressionValidatingWriter 创建解压验证写入器
func NewDecompressionValidatingWriter(writer io.Writer, config *config.Config, compressedSize int64, tracker *SizeTracker) *ValidatingWriter {
	return &ValidatingWriter{
		writer:       writer,
		config:       config,
		totalWritten: 0,
		errorPrefix:  "解压后",
		tracker:      tracker,
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
func (vw *ValidatingWriter) Write(p []byte) (n int, err error) {
	// 写入数据到底层写入器
	n, err = vw.writer.Write(p)
	if err != nil {
		return n, err
	}

	// 更新总写入大小
	vw.totalWritten += int64(n)

	// 验证文件大小是否超过单文件限制
	if vw.config.EnableSizeCheck && vw.totalWritten > vw.config.MaxFileSize {
		return n, fmt.Errorf("%s文件大小 %s 超过单文件限制 %s",
			vw.errorPrefix, FormatFileSize(vw.totalWritten), FormatFileSize(vw.config.MaxFileSize))
	}

	// 验证文件大小是否超过总大小限制
	if vw.config.EnableSizeCheck && vw.totalWritten > vw.config.MaxTotalSize {
		return n, fmt.Errorf("%s文件大小 %s 超过总大小限制 %s",
			vw.errorPrefix, FormatFileSize(vw.totalWritten), FormatFileSize(vw.config.MaxTotalSize))
	}

	return n, nil
}

// GetTotalWritten 获取已写入的总大小
//
// 返回值:
//   - int64: 已写入的总大小（字节）
func (vw *ValidatingWriter) GetTotalWritten() int64 {
	return vw.totalWritten
}

// GetSizeTracker 获取大小跟踪器（如果有）
//
// 返回值:
//   - *SizeTracker: 大小跟踪器，可能为nil
func (vw *ValidatingWriter) GetSizeTracker() *SizeTracker {
	return vw.tracker
}

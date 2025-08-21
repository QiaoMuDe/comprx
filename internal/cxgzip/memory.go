package cxgzip

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"

	"gitee.com/MM-Q/comprx/config"
	"gitee.com/MM-Q/comprx/types"
)

// CompressBytes 压缩字节数据到内存
//
// 参数:
//   - data: 要压缩的字节数据
//   - level: 压缩级别
//
// 返回:
//   - []byte: 压缩后的数据
//   - error: 错误信息
func CompressBytes(data []byte, level types.CompressionLevel) (result []byte, err error) {
	// 参数验证 - 更精确的nil检查
	if data == nil {
		return nil, fmt.Errorf("输入数据不能为nil")
	}
	if len(data) == 0 {
		return nil, fmt.Errorf("输入数据不能为空")
	}

	// 创建内存缓冲区 - 预分配容量减少重分配
	// 预分配原大小的50%
	estimatedSize := len(data) / 2
	if estimatedSize < 64 {
		estimatedSize = 64 // 最小64字节
	}
	buf := bytes.NewBuffer(make([]byte, 0, estimatedSize))

	// 创建gzip写入器
	writer, err := gzip.NewWriterLevel(buf, config.GetCompressionLevel(level))
	if err != nil {
		return nil, fmt.Errorf("创建gzip写入器失败: %w", err)
	}

	// 直接写入数据，无需额外缓冲区
	if _, err = writer.Write(data); err != nil {
		_ = writer.Close() // 确保资源清理
		return nil, fmt.Errorf("压缩数据失败: %w", err)
	}

	// 关闭写入器确保数据完整写入
	if err = writer.Close(); err != nil {
		return nil, fmt.Errorf("完成压缩失败: %w", err)
	}

	return buf.Bytes(), nil
}

// DecompressBytes 从内存解压字节数据
//
// 参数:
//   - compressedData: 压缩的字节数据
//
// 返回:
//   - []byte: 解压后的数据
//   - error: 错误信息
func DecompressBytes(compressedData []byte) (result []byte, err error) {
	// 参数验证 - 更精确的nil检查
	if compressedData == nil {
		return nil, fmt.Errorf("压缩数据不能为nil")
	}
	if len(compressedData) == 0 {
		return nil, fmt.Errorf("压缩数据不能为空")
	}

	// 创建字节读取器
	reader := bytes.NewReader(compressedData)

	// 创建gzip读取器
	gzipReader, err := gzip.NewReader(reader)
	if err != nil {
		return nil, fmt.Errorf("创建gzip读取器失败: %w", err)
	}

	// 预分配解压缓冲区 - 解压通常是压缩数据的2-3倍
	estimatedSize := len(compressedData) * 2
	if estimatedSize < 128 {
		estimatedSize = 128 // 最小128字节
	}
	buf := bytes.NewBuffer(make([]byte, 0, estimatedSize))

	// 直接读取解压数据，无需额外缓冲区
	if _, err = io.Copy(buf, gzipReader); err != nil {
		_ = gzipReader.Close() // 确保资源清理
		return nil, fmt.Errorf("解压数据失败: %w", err)
	}

	// 关闭读取器
	if err = gzipReader.Close(); err != nil {
		return nil, fmt.Errorf("关闭gzip读取器失败: %w", err)
	}

	return buf.Bytes(), nil
}

// CompressString 压缩字符串到内存
//
// 参数:
//   - text: 要压缩的字符串
//   - level: 压缩级别
//
// 返回:
//   - []byte: 压缩后的数据
//   - error: 错误信息
func CompressString(text string, level types.CompressionLevel) ([]byte, error) {
	// 快速失败判断
	if text == "" {
		return nil, fmt.Errorf("输入字符串不能为空")
	}

	// 直接复用CompressBytes
	return CompressBytes([]byte(text), level)
}

// DecompressString 从内存解压为字符串
//
// 参数:
//   - compressedData: 压缩的字节数据
//
// 返回:
//   - string: 解压后的字符串
//   - error: 错误信息
func DecompressString(compressedData []byte) (string, error) {
	// 快速失败判断
	if compressedData == nil {
		return "", fmt.Errorf("压缩数据不能为nil")
	}
	if len(compressedData) == 0 {
		return "", fmt.Errorf("压缩数据不能为空")
	}

	// 先解压为字节
	decompressed, err := DecompressBytes(compressedData)
	if err != nil {
		return "", err
	}

	// 转换为字符串
	return string(decompressed), nil
}

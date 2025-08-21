package cxgzip

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"

	"gitee.com/MM-Q/comprx/config"
	"gitee.com/MM-Q/comprx/internal/utils"
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
func CompressBytes(data []byte, level types.CompressionLevel) ([]byte, error) {
	// 参数验证
	if data == nil {
		return nil, fmt.Errorf("输入数据不能为nil")
	}

	// 创建内存缓冲区
	var buf bytes.Buffer

	// 创建gzip写入器
	writer, err := gzip.NewWriterLevel(&buf, config.GetCompressionLevel(&config.Config{
		CompressionLevel: level,
	}))
	if err != nil {
		return nil, fmt.Errorf("创建gzip写入器失败: %w", err)
	}
	defer func() { _ = writer.Close() }()

	// 使用32KB缓冲区和CopyBuffer进行数据传输
	reader := bytes.NewReader(data)
	buffer := utils.GetBuffer(32 * 1024) // 使用32KB缓冲区
	defer utils.PutBuffer(buffer)

	if _, err := io.CopyBuffer(writer, reader, buffer); err != nil {
		return nil, fmt.Errorf("压缩数据失败: %w", err)
	}

	// 关闭写入器确保数据完整写入
	if err := writer.Close(); err != nil {
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
func DecompressBytes(compressedData []byte) ([]byte, error) {
	// 参数验证
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
	defer func() { _ = gzipReader.Close() }()

	// 使用32KB缓冲区和CopyBuffer读取解压数据
	var buf bytes.Buffer
	buffer := utils.GetBuffer(32 * 1024) // 使用32KB缓冲区
	defer utils.PutBuffer(buffer)

	if _, err := io.CopyBuffer(&buf, gzipReader, buffer); err != nil {
		return nil, fmt.Errorf("解压数据失败: %w", err)
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
	// 先解压为字节
	decompressed, err := DecompressBytes(compressedData)
	if err != nil {
		return "", err
	}

	// 转换为字符串
	return string(decompressed), nil
}

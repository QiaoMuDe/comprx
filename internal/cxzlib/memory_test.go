package cxzlib

import (
	"bytes"
	"strings"
	"testing"

	"gitee.com/MM-Q/comprx/types"
)

func TestCompressBytes(t *testing.T) {
	testData := []byte("Hello, ZLIB memory compression test!")

	// 测试压缩
	compressed, err := CompressBytes(testData, types.CompressionLevelDefault)
	if err != nil {
		t.Fatalf("压缩字节数据失败: %v", err)
	}

	if len(compressed) == 0 {
		t.Fatalf("压缩数据长度为0")
	}

	// 验证压缩数据不等于原始数据
	if bytes.Equal(compressed, testData) {
		t.Fatalf("压缩数据与原始数据相同")
	}

	t.Logf("原始数据: %d 字节, 压缩数据: %d 字节", len(testData), len(compressed))
}

func TestDecompressBytes(t *testing.T) {
	testData := []byte("Hello, ZLIB memory decompression test!")

	// 先压缩
	compressed, err := CompressBytes(testData, types.CompressionLevelDefault)
	if err != nil {
		t.Fatalf("压缩字节数据失败: %v", err)
	}

	// 再解压
	decompressed, err := DecompressBytes(compressed)
	if err != nil {
		t.Fatalf("解压字节数据失败: %v", err)
	}

	// 验证解压结果
	if !bytes.Equal(decompressed, testData) {
		t.Fatalf("解压数据与原始数据不匹配\n期望: %s\n实际: %s", string(testData), string(decompressed))
	}

	t.Logf("往返压缩解压成功")
}

func TestCompressString(t *testing.T) {
	testString := "Hello, ZLIB string compression test!"

	// 测试压缩
	compressed, err := CompressString(testString, types.CompressionLevelDefault)
	if err != nil {
		t.Fatalf("压缩字符串失败: %v", err)
	}

	if len(compressed) == 0 {
		t.Fatalf("压缩数据长度为0")
	}

	t.Logf("原始字符串: %d 字节, 压缩数据: %d 字节", len(testString), len(compressed))
}

func TestDecompressString(t *testing.T) {
	testString := "Hello, ZLIB string decompression test!"

	// 先压缩
	compressed, err := CompressString(testString, types.CompressionLevelDefault)
	if err != nil {
		t.Fatalf("压缩字符串失败: %v", err)
	}

	// 再解压
	decompressed, err := DecompressString(compressed)
	if err != nil {
		t.Fatalf("解压字符串失败: %v", err)
	}

	// 验证解压结果
	if decompressed != testString {
		t.Fatalf("解压字符串与原始字符串不匹配\n期望: %s\n实际: %s", testString, decompressed)
	}

	t.Logf("字符串往返压缩解压成功")
}

func TestCompressBytesNilInput(t *testing.T) {
	// 测试nil输入
	_, err := CompressBytes(nil, types.CompressionLevelDefault)
	if err == nil {
		t.Fatalf("压缩nil数据应该失败，但成功了")
	}

	expectedError := "输入数据不能为nil"
	if err.Error() != expectedError {
		t.Fatalf("期望错误: %s, 实际错误: %s", expectedError, err.Error())
	}
}

func TestCompressBytesEmptyInput(t *testing.T) {
	// 测试空输入
	_, err := CompressBytes([]byte{}, types.CompressionLevelDefault)
	if err == nil {
		t.Fatalf("压缩空数据应该失败，但成功了")
	}

	expectedError := "输入数据不能为空"
	if err.Error() != expectedError {
		t.Fatalf("期望错误: %s, 实际错误: %s", expectedError, err.Error())
	}
}

func TestDecompressBytesNilInput(t *testing.T) {
	// 测试nil输入
	_, err := DecompressBytes(nil)
	if err == nil {
		t.Fatalf("解压nil数据应该失败，但成功了")
	}

	expectedError := "压缩数据不能为nil"
	if err.Error() != expectedError {
		t.Fatalf("期望错误: %s, 实际错误: %s", expectedError, err.Error())
	}
}

func TestDecompressBytesEmptyInput(t *testing.T) {
	// 测试空输入
	_, err := DecompressBytes([]byte{})
	if err == nil {
		t.Fatalf("解压空数据应该失败，但成功了")
	}

	expectedError := "压缩数据不能为空"
	if err.Error() != expectedError {
		t.Fatalf("期望错误: %s, 实际错误: %s", expectedError, err.Error())
	}
}

func TestCompressStringEmptyInput(t *testing.T) {
	// 测试空字符串
	_, err := CompressString("", types.CompressionLevelDefault)
	if err == nil {
		t.Fatalf("压缩空字符串应该失败，但成功了")
	}

	expectedError := "输入字符串不能为空"
	if err.Error() != expectedError {
		t.Fatalf("期望错误: %s, 实际错误: %s", expectedError, err.Error())
	}
}

func TestDecompressBytesInvalidData(t *testing.T) {
	// 测试无效的压缩数据
	invalidData := []byte("this is not valid zlib data")
	_, err := DecompressBytes(invalidData)
	if err == nil {
		t.Fatalf("解压无效数据应该失败，但成功了")
	}
}

func TestDifferentCompressionLevels(t *testing.T) {
	testData := []byte(strings.Repeat("Hello, ZLIB compression level test! ", 100))

	levels := []types.CompressionLevel{
		types.CompressionLevelFast,
		types.CompressionLevelDefault,
		types.CompressionLevelBest,
	}

	for _, level := range levels {
		t.Run(level.String(), func(t *testing.T) {
			// 压缩
			compressed, err := CompressBytes(testData, level)
			if err != nil {
				t.Fatalf("压缩失败 (等级 %s): %v", level.String(), err)
			}

			// 解压
			decompressed, err := DecompressBytes(compressed)
			if err != nil {
				t.Fatalf("解压失败 (等级 %s): %v", level.String(), err)
			}

			// 验证
			if !bytes.Equal(decompressed, testData) {
				t.Fatalf("往返压缩解压数据不匹配 (等级 %s)", level.String())
			}

			t.Logf("等级 %s: 原始 %d 字节 -> 压缩 %d 字节, 压缩率 %.2f%%",
				level.String(),
				len(testData),
				len(compressed),
				float64(len(compressed))/float64(len(testData))*100)
		})
	}
}

func TestCompressStreamNilWriter(t *testing.T) {
	src := strings.NewReader("test data")
	err := CompressStream(nil, src, types.CompressionLevelDefault)
	if err == nil {
		t.Fatalf("nil写入器应该失败，但成功了")
	}

	expectedError := "目标写入器不能为nil"
	if err.Error() != expectedError {
		t.Fatalf("期望错误: %s, 实际错误: %s", expectedError, err.Error())
	}
}

func TestCompressStreamNilReader(t *testing.T) {
	var buf bytes.Buffer
	err := CompressStream(&buf, nil, types.CompressionLevelDefault)
	if err == nil {
		t.Fatalf("nil读取器应该失败，但成功了")
	}

	expectedError := "源读取器不能为nil"
	if err.Error() != expectedError {
		t.Fatalf("期望错误: %s, 实际错误: %s", expectedError, err.Error())
	}
}

func TestDecompressStreamNilWriter(t *testing.T) {
	src := strings.NewReader("test data")
	err := DecompressStream(nil, src)
	if err == nil {
		t.Fatalf("nil写入器应该失败，但成功了")
	}

	expectedError := "目标写入器不能为nil"
	if err.Error() != expectedError {
		t.Fatalf("期望错误: %s, 实际错误: %s", expectedError, err.Error())
	}
}

func TestDecompressStreamNilReader(t *testing.T) {
	var buf bytes.Buffer
	err := DecompressStream(&buf, nil)
	if err == nil {
		t.Fatalf("nil读取器应该失败，但成功了")
	}

	expectedError := "源读取器不能为nil"
	if err.Error() != expectedError {
		t.Fatalf("期望错误: %s, 实际错误: %s", expectedError, err.Error())
	}
}

func TestStreamCompression(t *testing.T) {
	testData := "Hello, ZLIB stream compression test! " + strings.Repeat("This is repeated content. ", 50)

	// 压缩流
	src := strings.NewReader(testData)
	var compressed bytes.Buffer
	err := CompressStream(&compressed, src, types.CompressionLevelDefault)
	if err != nil {
		t.Fatalf("流式压缩失败: %v", err)
	}

	// 解压流
	var decompressed bytes.Buffer
	err = DecompressStream(&decompressed, &compressed)
	if err != nil {
		t.Fatalf("流式解压失败: %v", err)
	}

	// 验证结果
	if decompressed.String() != testData {
		t.Fatalf("流式压缩解压数据不匹配")
	}

	t.Logf("流式压缩: 原始 %d 字节 -> 压缩 %d 字节 -> 解压 %d 字节",
		len(testData), compressed.Len(), decompressed.Len())
}

func BenchmarkCompressBytes(b *testing.B) {
	testData := []byte(strings.Repeat("Hello, ZLIB benchmark test! ", 100))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := CompressBytes(testData, types.CompressionLevelDefault)
		if err != nil {
			b.Fatalf("压缩失败: %v", err)
		}
	}
}

func BenchmarkDecompressBytes(b *testing.B) {
	testData := []byte(strings.Repeat("Hello, ZLIB benchmark test! ", 100))
	compressed, err := CompressBytes(testData, types.CompressionLevelDefault)
	if err != nil {
		b.Fatalf("压缩失败: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := DecompressBytes(compressed)
		if err != nil {
			b.Fatalf("解压失败: %v", err)
		}
	}
}

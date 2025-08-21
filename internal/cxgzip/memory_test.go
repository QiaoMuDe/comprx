package cxgzip

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"gitee.com/MM-Q/comprx/types"
)

// TestCompressBytes_NormalCases 测试正常压缩场景
func TestCompressBytes_NormalCases(t *testing.T) {
	tests := []struct {
		name  string
		data  []byte
		level types.CompressionLevel
	}{
		{
			name:  "小数据压缩",
			data:  []byte("Hello, World!"),
			level: types.CompressionLevelDefault,
		},
		{
			name:  "中等数据压缩",
			data:  bytes.Repeat([]byte("This is a test string for compression. "), 100),
			level: types.CompressionLevelBest,
		},
		{
			name:  "大数据压缩",
			data:  bytes.Repeat([]byte("Large data compression test. "), 10000),
			level: types.CompressionLevelFast,
		},
		{
			name:  "重复数据压缩",
			data:  bytes.Repeat([]byte("A"), 1000),
			level: types.CompressionLevelDefault,
		},
		{
			name:  "二进制数据压缩",
			data:  []byte{0x00, 0x01, 0x02, 0x03, 0xFF, 0xFE, 0xFD, 0xFC},
			level: types.CompressionLevelDefault,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 压缩数据
			compressed, err := CompressBytes(tt.data, tt.level)
			if err != nil {
				t.Fatalf("压缩失败: %v", err)
			}

			// 验证压缩结果不为空
			if len(compressed) == 0 {
				t.Fatal("压缩结果为空")
			}

			// 解压数据
			decompressed, err := DecompressBytes(compressed)
			if err != nil {
				t.Fatalf("解压失败: %v", err)
			}

			// 验证数据一致性
			if !bytes.Equal(tt.data, decompressed) {
				t.Fatalf("数据不一致:\n原始: %v\n解压: %v", tt.data, decompressed)
			}

			t.Logf("原始大小: %d, 压缩大小: %d, 压缩比: %.2f%%",
				len(tt.data), len(compressed), float64(len(compressed))/float64(len(tt.data))*100)
		})
	}
}

// TestCompressBytes_EdgeCases 测试边界情况
func TestCompressBytes_EdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		data        []byte
		level       types.CompressionLevel
		expectError bool
		errorMsg    string
	}{
		{
			name:        "nil数据",
			data:        nil,
			level:       types.CompressionLevelDefault,
			expectError: true,
			errorMsg:    "输入数据不能为nil",
		},
		{
			name:        "空数据",
			data:        []byte{},
			level:       types.CompressionLevelDefault,
			expectError: true,
			errorMsg:    "输入数据不能为空",
		},
		{
			name:        "单字节数据",
			data:        []byte{0x42},
			level:       types.CompressionLevelDefault,
			expectError: false,
		},
		{
			name:        "最小有效数据",
			data:        []byte("A"),
			level:       types.CompressionLevelDefault,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			compressed, err := CompressBytes(tt.data, tt.level)

			if tt.expectError {
				if err == nil {
					t.Fatal("期望出现错误，但没有错误")
				}
				if !strings.Contains(err.Error(), tt.errorMsg) {
					t.Fatalf("错误信息不匹配，期望包含: %s, 实际: %s", tt.errorMsg, err.Error())
				}
				return
			}

			if err != nil {
				t.Fatalf("不期望出现错误: %v", err)
			}

			// 验证可以正确解压
			decompressed, err := DecompressBytes(compressed)
			if err != nil {
				t.Fatalf("解压失败: %v", err)
			}

			if !bytes.Equal(tt.data, decompressed) {
				t.Fatalf("数据不一致")
			}
		})
	}
}

// TestDecompressBytes_EdgeCases 测试解压边界情况
func TestDecompressBytes_EdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		data        []byte
		expectError bool
		errorMsg    string
	}{
		{
			name:        "nil压缩数据",
			data:        nil,
			expectError: true,
			errorMsg:    "压缩数据不能为nil",
		},
		{
			name:        "空压缩数据",
			data:        []byte{},
			expectError: true,
			errorMsg:    "压缩数据不能为空",
		},
		{
			name:        "无效gzip数据",
			data:        []byte("invalid gzip data"),
			expectError: true,
			errorMsg:    "创建gzip读取器失败",
		},
		{
			name:        "损坏的gzip头",
			data:        []byte{0x1f, 0x8b, 0x08, 0x00}, // 不完整的gzip头
			expectError: true,
			errorMsg:    "创建gzip读取器失败",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			decompressed, err := DecompressBytes(tt.data)

			if tt.expectError {
				if err == nil {
					t.Fatal("期望出现错误，但没有错误")
				}
				if !strings.Contains(err.Error(), tt.errorMsg) {
					t.Fatalf("错误信息不匹配，期望包含: %s, 实际: %s", tt.errorMsg, err.Error())
				}
				return
			}

			if err != nil {
				t.Fatalf("不期望出现错误: %v", err)
			}

			if decompressed == nil {
				t.Fatal("解压结果为nil")
			}
		})
	}
}

// TestCompressString_NormalCases 测试字符串压缩正常场景
func TestCompressString_NormalCases(t *testing.T) {
	tests := []struct {
		name  string
		text  string
		level types.CompressionLevel
	}{
		{
			name:  "简单字符串",
			text:  "Hello, World!",
			level: types.CompressionLevelDefault,
		},
		{
			name:  "中文字符串",
			text:  "你好，世界！这是一个测试字符串。",
			level: types.CompressionLevelDefault,
		},
		{
			name:  "长字符串",
			text:  strings.Repeat("This is a long string for testing compression. ", 100),
			level: types.CompressionLevelBest,
		},
		{
			name:  "特殊字符",
			text:  "!@#$%^&*()_+-=[]{}|;':\",./<>?`~",
			level: types.CompressionLevelDefault,
		},
		{
			name:  "JSON格式字符串",
			text:  `{"name":"test","value":123,"array":[1,2,3],"nested":{"key":"value"}}`,
			level: types.CompressionLevelDefault,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 压缩字符串
			compressed, err := CompressString(tt.text, tt.level)
			if err != nil {
				t.Fatalf("压缩失败: %v", err)
			}

			// 解压字符串
			decompressed, err := DecompressString(compressed)
			if err != nil {
				t.Fatalf("解压失败: %v", err)
			}

			// 验证数据一致性
			if tt.text != decompressed {
				t.Fatalf("字符串不一致:\n原始: %s\n解压: %s", tt.text, decompressed)
			}

			t.Logf("原始大小: %d, 压缩大小: %d, 压缩比: %.2f%%",
				len(tt.text), len(compressed), float64(len(compressed))/float64(len(tt.text))*100)
		})
	}
}

// TestCompressString_EdgeCases 测试字符串压缩边界情况
func TestCompressString_EdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		text        string
		level       types.CompressionLevel
		expectError bool
		errorMsg    string
	}{
		{
			name:        "空字符串",
			text:        "",
			level:       types.CompressionLevelDefault,
			expectError: true,
			errorMsg:    "输入字符串不能为空",
		},
		{
			name:        "单字符",
			text:        "A",
			level:       types.CompressionLevelDefault,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			compressed, err := CompressString(tt.text, tt.level)

			if tt.expectError {
				if err == nil {
					t.Fatal("期望出现错误，但没有错误")
				}
				if !strings.Contains(err.Error(), tt.errorMsg) {
					t.Fatalf("错误信息不匹配，期望包含: %s, 实际: %s", tt.errorMsg, err.Error())
				}
				return
			}

			if err != nil {
				t.Fatalf("不期望出现错误: %v", err)
			}

			// 验证可以正确解压
			decompressed, err := DecompressString(compressed)
			if err != nil {
				t.Fatalf("解压失败: %v", err)
			}

			if tt.text != decompressed {
				t.Fatalf("字符串不一致")
			}
		})
	}
}

// TestDecompressString_EdgeCases 测试字符串解压边界情况
func TestDecompressString_EdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		data        []byte
		expectError bool
		errorMsg    string
	}{
		{
			name:        "nil压缩数据",
			data:        nil,
			expectError: true,
			errorMsg:    "压缩数据不能为nil",
		},
		{
			name:        "空压缩数据",
			data:        []byte{},
			expectError: true,
			errorMsg:    "压缩数据不能为空",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			decompressed, err := DecompressString(tt.data)

			if tt.expectError {
				if err == nil {
					t.Fatal("期望出现错误，但没有错误")
				}
				if !strings.Contains(err.Error(), tt.errorMsg) {
					t.Fatalf("错误信息不匹配，期望包含: %s, 实际: %s", tt.errorMsg, err.Error())
				}
				return
			}

			if err != nil {
				t.Fatalf("不期望出现错误: %v", err)
			}

			if decompressed == "" {
				t.Fatal("解压结果为空字符串")
			}
		})
	}
}

// TestCompressionLevels 测试不同压缩级别
func TestCompressionLevels(t *testing.T) {
	testData := bytes.Repeat([]byte("This is test data for compression level testing. "), 200)

	levels := []types.CompressionLevel{
		types.CompressionLevelFast,
		types.CompressionLevelDefault,
		types.CompressionLevelBest,
	}

	results := make(map[types.CompressionLevel]int)

	for _, level := range levels {
		levelName := fmt.Sprintf("Level_%v", level)
		t.Run(levelName, func(t *testing.T) {
			compressed, err := CompressBytes(testData, level)
			if err != nil {
				t.Fatalf("压缩失败: %v", err)
			}

			// 验证可以正确解压
			decompressed, err := DecompressBytes(compressed)
			if err != nil {
				t.Fatalf("解压失败: %v", err)
			}

			if !bytes.Equal(testData, decompressed) {
				t.Fatal("数据不一致")
			}

			results[level] = len(compressed)
			t.Logf("压缩级别 %v: 原始 %d -> 压缩 %d (%.2f%%)",
				level, len(testData), len(compressed),
				float64(len(compressed))/float64(len(testData))*100)
		})
	}

	// 验证压缩级别效果（一般情况下 Best < Default < Fast）
	t.Logf("压缩大小比较: Fast=%d, Default=%d, Best=%d",
		results[types.CompressionLevelFast],
		results[types.CompressionLevelDefault],
		results[types.CompressionLevelBest])
}

// TestLargeData 测试大数据压缩
func TestLargeData(t *testing.T) {
	// 创建1MB的测试数据
	largeData := make([]byte, 1024*1024)
	for i := range largeData {
		largeData[i] = byte(i % 256)
	}

	compressed, err := CompressBytes(largeData, types.CompressionLevelDefault)
	if err != nil {
		t.Fatalf("大数据压缩失败: %v", err)
	}

	decompressed, err := DecompressBytes(compressed)
	if err != nil {
		t.Fatalf("大数据解压失败: %v", err)
	}

	if !bytes.Equal(largeData, decompressed) {
		t.Fatal("大数据压缩解压后不一致")
	}

	t.Logf("大数据测试: 原始 %d -> 压缩 %d (%.2f%%)",
		len(largeData), len(compressed),
		float64(len(compressed))/float64(len(largeData))*100)
}

// BenchmarkCompressBytes 压缩性能基准测试
func BenchmarkCompressBytes(b *testing.B) {
	data := bytes.Repeat([]byte("benchmark test data "), 1000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := CompressBytes(data, types.CompressionLevelDefault)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkDecompressBytes 解压性能基准测试
func BenchmarkDecompressBytes(b *testing.B) {
	data := bytes.Repeat([]byte("benchmark test data "), 1000)
	compressed, err := CompressBytes(data, types.CompressionLevelDefault)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := DecompressBytes(compressed)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkCompressString 字符串压缩性能基准测试
func BenchmarkCompressString(b *testing.B) {
	text := strings.Repeat("benchmark test string ", 1000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := CompressString(text, types.CompressionLevelDefault)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkDecompressString 字符串解压性能基准测试
func BenchmarkDecompressString(b *testing.B) {
	text := strings.Repeat("benchmark test string ", 1000)
	compressed, err := CompressString(text, types.CompressionLevelDefault)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := DecompressString(compressed)
		if err != nil {
			b.Fatal(err)
		}
	}
}

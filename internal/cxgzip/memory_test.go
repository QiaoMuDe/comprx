package cxgzip

import (
	"bytes"
	"fmt"
	"io"
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

// ==================== 流式压缩API测试 ====================

func TestCompressStream(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		level   types.CompressionLevel
		wantErr bool
		errMsg  string
	}{
		{
			name:    "正常压缩_默认等级",
			input:   "Hello, World! This is a test for stream compression.",
			level:   types.CompressionLevelDefault,
			wantErr: false,
		},
		{
			name:    "正常压缩_最快等级",
			input:   "Fast compression test data.",
			level:   types.CompressionLevelFast,
			wantErr: false,
		},
		{
			name:    "正常压缩_最佳等级",
			input:   "Best compression test data for maximum compression ratio.",
			level:   types.CompressionLevelBest,
			wantErr: false,
		},
		{
			name:    "空数据压缩",
			input:   "",
			level:   types.CompressionLevelDefault,
			wantErr: false,
		},
		{
			name:    "大数据压缩",
			input:   strings.Repeat("Large data test. ", 1000),
			level:   types.CompressionLevelDefault,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 准备输入和输出
			src := strings.NewReader(tt.input)
			var dst bytes.Buffer

			// 执行压缩
			err := CompressStream(&dst, src, tt.level)

			// 检查错误
			if (err != nil) != tt.wantErr {
				t.Errorf("CompressStream() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("期望错误信息包含 %q, 实际错误: %v", tt.errMsg, err)
				}
				return
			}

			// 验证压缩结果
			compressedData := dst.Bytes()
			if len(compressedData) == 0 && len(tt.input) > 0 {
				t.Error("压缩后数据为空")
				return
			}

			// 验证可以正确解压
			decompressed, err := DecompressBytes(compressedData)
			if err != nil {
				t.Errorf("解压验证失败: %v", err)
				return
			}

			if string(decompressed) != tt.input {
				t.Errorf("解压后数据不匹配, 期望: %q, 实际: %q", tt.input, string(decompressed))
			}
		})
	}
}

func TestCompressStream_ErrorCases(t *testing.T) {
	tests := []struct {
		name    string
		dst     io.Writer
		src     io.Reader
		level   types.CompressionLevel
		wantErr string
	}{
		{
			name:    "目标写入器为nil",
			dst:     nil,
			src:     strings.NewReader("test"),
			level:   types.CompressionLevelDefault,
			wantErr: "目标写入器不能为nil",
		},
		{
			name:    "源读取器为nil",
			dst:     &bytes.Buffer{},
			src:     nil,
			level:   types.CompressionLevelDefault,
			wantErr: "源读取器不能为nil",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CompressStream(tt.dst, tt.src, tt.level)
			if err == nil {
				t.Error("期望出现错误，但没有错误")
				return
			}

			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("期望错误信息包含 %q, 实际错误: %v", tt.wantErr, err)
			}
		})
	}
}

func TestDecompressStream(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "正常解压",
			input:   "Hello, World! This is a test for stream decompression.",
			wantErr: false,
		},
		{
			name:    "单字符解压",
			input:   "A",
			wantErr: false,
		},
		{
			name:    "大数据解压",
			input:   strings.Repeat("Large data decompression test. ", 1000),
			wantErr: false,
		},
		{
			name:    "中文数据解压",
			input:   "这是一个中文测试数据，用于验证流式解压功能。",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 先压缩数据
			compressedData, err := CompressBytes([]byte(tt.input), types.CompressionLevelDefault)
			if err != nil {
				t.Fatalf("压缩测试数据失败: %v", err)
			}

			// 准备输入和输出
			src := bytes.NewReader(compressedData)
			var dst bytes.Buffer

			// 执行解压
			err = DecompressStream(&dst, src)

			// 检查错误
			if (err != nil) != tt.wantErr {
				t.Errorf("DecompressStream() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			// 验证解压结果
			result := dst.String()
			if result != tt.input {
				t.Errorf("解压后数据不匹配, 期望: %q, 实际: %q", tt.input, result)
			}
		})
	}
}

func TestDecompressStream_ErrorCases(t *testing.T) {
	tests := []struct {
		name    string
		dst     io.Writer
		src     io.Reader
		wantErr string
	}{
		{
			name:    "目标写入器为nil",
			dst:     nil,
			src:     strings.NewReader("test"),
			wantErr: "目标写入器不能为nil",
		},
		{
			name:    "源读取器为nil",
			dst:     &bytes.Buffer{},
			src:     nil,
			wantErr: "源读取器不能为nil",
		},
		{
			name:    "无效gzip数据",
			dst:     &bytes.Buffer{},
			src:     strings.NewReader("invalid gzip data"),
			wantErr: "创建gzip读取器失败",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := DecompressStream(tt.dst, tt.src)
			if err == nil {
				t.Error("期望出现错误，但没有错误")
				return
			}

			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("期望错误信息包含 %q, 实际错误: %v", tt.wantErr, err)
			}
		})
	}
}

func TestStreamAPI_RoundTrip(t *testing.T) {
	tests := []struct {
		name  string
		data  []byte
		level types.CompressionLevel
	}{
		{
			name:  "小数据往返",
			data:  []byte("Hello, World!"),
			level: types.CompressionLevelDefault,
		},
		{
			name:  "中等数据往返",
			data:  bytes.Repeat([]byte("Test data. "), 100),
			level: types.CompressionLevelFast,
		},
		{
			name:  "大数据往返",
			data:  bytes.Repeat([]byte("Large test data for round trip. "), 1000),
			level: types.CompressionLevelBest,
		},
		{
			name:  "二进制数据往返",
			data:  []byte{0x00, 0x01, 0x02, 0x03, 0xFF, 0xFE, 0xFD, 0xFC},
			level: types.CompressionLevelDefault,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 第一步：压缩
			src := bytes.NewReader(tt.data)
			var compressed bytes.Buffer

			err := CompressStream(&compressed, src, tt.level)
			if err != nil {
				t.Fatalf("压缩失败: %v", err)
			}

			// 第二步：解压
			compressedSrc := bytes.NewReader(compressed.Bytes())
			var decompressed bytes.Buffer

			err = DecompressStream(&decompressed, compressedSrc)
			if err != nil {
				t.Fatalf("解压失败: %v", err)
			}

			// 验证数据一致性
			if !bytes.Equal(tt.data, decompressed.Bytes()) {
				t.Errorf("往返后数据不一致")
				t.Errorf("原始数据长度: %d", len(tt.data))
				t.Errorf("解压数据长度: %d", decompressed.Len())
			}
		})
	}
}

func TestStreamAPI_LargeData(t *testing.T) {
	// 创建1MB的测试数据
	largeData := make([]byte, 1024*1024)
	for i := range largeData {
		largeData[i] = byte(i % 256)
	}

	// 压缩
	src := bytes.NewReader(largeData)
	var compressed bytes.Buffer

	err := CompressStream(&compressed, src, types.CompressionLevelDefault)
	if err != nil {
		t.Fatalf("压缩大数据失败: %v", err)
	}

	// 验证压缩效果
	compressionRatio := float64(compressed.Len()) / float64(len(largeData))
	t.Logf("压缩比: %.2f%% (原始: %d bytes, 压缩后: %d bytes)",
		compressionRatio*100, len(largeData), compressed.Len())

	// 解压
	compressedSrc := bytes.NewReader(compressed.Bytes())
	var decompressed bytes.Buffer

	err = DecompressStream(&decompressed, compressedSrc)
	if err != nil {
		t.Fatalf("解压大数据失败: %v", err)
	}

	// 验证数据完整性
	if !bytes.Equal(largeData, decompressed.Bytes()) {
		t.Error("大数据往返后不一致")
	}
}

// ==================== 流式API基准测试 ====================

func BenchmarkCompressStream(b *testing.B) {
	testData := bytes.Repeat([]byte("Benchmark test data for stream compression. "), 1000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		src := bytes.NewReader(testData)
		var dst bytes.Buffer

		err := CompressStream(&dst, src, types.CompressionLevelDefault)
		if err != nil {
			b.Fatalf("压缩失败: %v", err)
		}
	}
}

func BenchmarkDecompressStream(b *testing.B) {
	// 准备压缩数据
	testData := bytes.Repeat([]byte("Benchmark test data for stream decompression. "), 1000)
	compressedData, err := CompressBytes(testData, types.CompressionLevelDefault)
	if err != nil {
		b.Fatalf("准备测试数据失败: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		src := bytes.NewReader(compressedData)
		var dst bytes.Buffer

		err := DecompressStream(&dst, src)
		if err != nil {
			b.Fatalf("解压失败: %v", err)
		}
	}
}

func BenchmarkStreamAPI_RoundTrip(b *testing.B) {
	testData := bytes.Repeat([]byte("Round trip benchmark test data. "), 1000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// 压缩
		src := bytes.NewReader(testData)
		var compressed bytes.Buffer

		err := CompressStream(&compressed, src, types.CompressionLevelDefault)
		if err != nil {
			b.Fatalf("压缩失败: %v", err)
		}

		// 解压
		compressedSrc := bytes.NewReader(compressed.Bytes())
		var decompressed bytes.Buffer

		err = DecompressStream(&decompressed, compressedSrc)
		if err != nil {
			b.Fatalf("解压失败: %v", err)
		}
	}
}

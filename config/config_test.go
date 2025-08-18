package config

import (
	"compress/gzip"
	"testing"
)

func TestCompressionLevel_Constants(t *testing.T) {
	// 测试压缩等级常量的值
	tests := []struct {
		name     string
		level    CompressionLevel
		expected int
	}{
		{"默认压缩等级", CompressionLevelDefault, -1},
		{"不压缩", CompressionLevelNone, 0},
		{"快速压缩", CompressionLevelFast, 1},
		{"最佳压缩", CompressionLevelBest, 9},
		{"仅Huffman编码", CompressionLevelHuffmanOnly, -2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if int(tt.level) != tt.expected {
				t.Errorf("CompressionLevel %s = %v, want %v", tt.name, int(tt.level), tt.expected)
			}
		})
	}
}

func TestNew(t *testing.T) {
	// 测试创建新的配置
	config := New()

	if config == nil {
		t.Fatal("New() 返回 nil")
	}

	// 检查默认值
	if config.CompressionLevel != CompressionLevelDefault {
		t.Errorf("默认压缩等级 = %v, want %v", config.CompressionLevel, CompressionLevelDefault)
	}

	if config.OverwriteExisting != false {
		t.Errorf("默认覆盖设置 = %v, want %v", config.OverwriteExisting, false)
	}
}

func TestConfig_Fields(t *testing.T) {
	// 测试配置结构体字段的设置和获取
	config := &Config{
		CompressionLevel:  CompressionLevelBest,
		OverwriteExisting: true,
	}

	if config.CompressionLevel != CompressionLevelBest {
		t.Errorf("CompressionLevel = %v, want %v", config.CompressionLevel, CompressionLevelBest)
	}

	if config.OverwriteExisting != true {
		t.Errorf("OverwriteExisting = %v, want %v", config.OverwriteExisting, true)
	}
}

func TestGetCompressionLevel(t *testing.T) {
	tests := []struct {
		name              string
		configLevel       CompressionLevel
		expectedGzipLevel int
	}{
		{
			name:              "不压缩",
			configLevel:       CompressionLevelNone,
			expectedGzipLevel: gzip.NoCompression,
		},
		{
			name:              "快速压缩",
			configLevel:       CompressionLevelFast,
			expectedGzipLevel: gzip.BestSpeed,
		},
		{
			name:              "最佳压缩",
			configLevel:       CompressionLevelBest,
			expectedGzipLevel: gzip.BestCompression,
		},
		{
			name:              "仅Huffman编码",
			configLevel:       CompressionLevelHuffmanOnly,
			expectedGzipLevel: gzip.HuffmanOnly,
		},
		{
			name:              "默认压缩",
			configLevel:       CompressionLevelDefault,
			expectedGzipLevel: gzip.DefaultCompression,
		},
		{
			name:              "未知压缩等级",
			configLevel:       CompressionLevel(999),
			expectedGzipLevel: gzip.DefaultCompression,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &Config{
				CompressionLevel: tt.configLevel,
			}

			result := GetCompressionLevel(config)
			if result != tt.expectedGzipLevel {
				t.Errorf("GetCompressionLevel() = %v, want %v", result, tt.expectedGzipLevel)
			}
		})
	}
}

func TestGetCompressionLevel_NilConfig(t *testing.T) {
	// 测试传入 nil 配置的情况
	defer func() {
		if r := recover(); r == nil {
			t.Error("GetCompressionLevel(nil) 应该引发 panic")
		}
	}()

	GetCompressionLevel(nil)
}

func TestConfig_Modification(t *testing.T) {
	// 测试配置的修改
	config := New()

	// 修改压缩等级
	config.CompressionLevel = CompressionLevelBest
	if config.CompressionLevel != CompressionLevelBest {
		t.Errorf("修改后的压缩等级 = %v, want %v", config.CompressionLevel, CompressionLevelBest)
	}

	// 修改覆盖设置
	config.OverwriteExisting = true
	if config.OverwriteExisting != true {
		t.Errorf("修改后的覆盖设置 = %v, want %v", config.OverwriteExisting, true)
	}
}

func TestConfig_MultipleInstances(t *testing.T) {
	// 测试多个配置实例的独立性
	config1 := New()
	config2 := New()

	// 修改第一个配置
	config1.CompressionLevel = CompressionLevelBest
	config1.OverwriteExisting = true

	// 检查第二个配置是否受影响
	if config2.CompressionLevel != CompressionLevelDefault {
		t.Error("config2 的压缩等级不应该受 config1 影响")
	}

	if config2.OverwriteExisting != false {
		t.Error("config2 的覆盖设置不应该受 config1 影响")
	}

	// 检查两个配置是否为不同的实例
	if config1 == config2 {
		t.Error("New() 应该返回不同的实例")
	}
}

func TestCompressionLevel_AllValues(t *testing.T) {
	// 测试所有压缩等级值的有效性
	levels := []CompressionLevel{
		CompressionLevelDefault,
		CompressionLevelNone,
		CompressionLevelFast,
		CompressionLevelBest,
		CompressionLevelHuffmanOnly,
	}

	for _, level := range levels {
		t.Run(string(rune(int(level))), func(t *testing.T) {
			config := &Config{
				CompressionLevel: level,
			}

			// 确保 GetCompressionLevel 不会 panic
			result := GetCompressionLevel(config)

			// 检查返回值是否在有效范围内
			validGzipLevels := []int{
				gzip.NoCompression,
				gzip.BestSpeed,
				gzip.BestCompression,
				gzip.HuffmanOnly,
				gzip.DefaultCompression,
			}

			found := false
			for _, validLevel := range validGzipLevels {
				if result == validLevel {
					found = true
					break
				}
			}

			if !found {
				t.Errorf("GetCompressionLevel() 返回无效的 gzip 压缩等级: %v", result)
			}
		})
	}
}

func TestConfig_ZeroValue(t *testing.T) {
	// 测试零值配置
	var config Config

	if config.CompressionLevel != 0 {
		t.Errorf("零值配置的压缩等级 = %v, want %v", config.CompressionLevel, 0)
	}

	if config.OverwriteExisting != false {
		t.Errorf("零值配置的覆盖设置 = %v, want %v", config.OverwriteExisting, false)
	}

	// 测试零值配置的 GetCompressionLevel
	// 零值配置的压缩等级为 0，对应 CompressionLevelNone，应该返回 gzip.NoCompression
	result := GetCompressionLevel(&config)
	if result != gzip.NoCompression {
		t.Errorf("零值配置的 GetCompressionLevel() = %v, want %v", result, gzip.NoCompression)
	}
}

func TestConfig_EdgeCases(t *testing.T) {
	// 测试边界情况
	t.Run("极大压缩等级值", func(t *testing.T) {
		config := &Config{
			CompressionLevel: CompressionLevel(1000),
		}
		result := GetCompressionLevel(config)
		if result != gzip.DefaultCompression {
			t.Errorf("极大压缩等级应该返回默认压缩等级，got %v", result)
		}
	})

	t.Run("极小压缩等级值", func(t *testing.T) {
		config := &Config{
			CompressionLevel: CompressionLevel(-1000),
		}
		result := GetCompressionLevel(config)
		if result != gzip.DefaultCompression {
			t.Errorf("极小压缩等级应该返回默认压缩等级，got %v", result)
		}
	})
}

// 基准测试
func BenchmarkNew(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = New()
	}
}

func BenchmarkGetCompressionLevel(b *testing.B) {
	config := &Config{
		CompressionLevel: CompressionLevelBest,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = GetCompressionLevel(config)
	}
}

func BenchmarkGetCompressionLevel_AllLevels(b *testing.B) {
	configs := []*Config{
		{CompressionLevel: CompressionLevelDefault},
		{CompressionLevel: CompressionLevelNone},
		{CompressionLevel: CompressionLevelFast},
		{CompressionLevel: CompressionLevelBest},
		{CompressionLevel: CompressionLevelHuffmanOnly},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, config := range configs {
			_ = GetCompressionLevel(config)
		}
	}
}

// 示例测试
func ExampleNew() {
	config := New()
	config.CompressionLevel = CompressionLevelBest
	config.OverwriteExisting = true
	// Output:
}

func ExampleGetCompressionLevel() {
	config := &Config{
		CompressionLevel: CompressionLevelBest,
	}
	level := GetCompressionLevel(config)
	_ = level // level 现在包含 gzip.BestCompression
	// Output:
}

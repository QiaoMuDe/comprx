package config

import (
	"compress/gzip"
	"testing"

	"gitee.com/MM-Q/comprx/internal/progress"
	"gitee.com/MM-Q/comprx/types"
)

// TestNew 测试创建新的配置实例
func TestNew(t *testing.T) {
	config := New()

	// 验证配置实例不为nil
	if config == nil {
		t.Fatal("New() 返回的配置实例不应该为nil")
	}

	// 验证默认压缩等级
	if config.CompressionLevel != types.CompressionLevelDefault {
		t.Errorf("期望默认压缩等级为 %v, 实际为 %v", types.CompressionLevelDefault, config.CompressionLevel)
	}

	// 验证默认覆盖设置
	if config.OverwriteExisting != false {
		t.Errorf("期望默认覆盖设置为 false, 实际为 %v", config.OverwriteExisting)
	}

	// 验证进度显示实例
	if config.Progress == nil {
		t.Error("进度显示实例不应该为nil")
	}

	// 验证默认路径验证设置
	if config.DisablePathValidation != false {
		t.Errorf("期望默认路径验证设置为 false, 实际为 %v", config.DisablePathValidation)
	}
}

// TestGetCompressionLevel 测试获取压缩等级
func TestGetCompressionLevel(t *testing.T) {
	testCases := []struct {
		name     string
		level    types.CompressionLevel
		expected int
	}{
		{
			name:     "无压缩",
			level:    types.CompressionLevelNone,
			expected: gzip.NoCompression,
		},
		{
			name:     "快速压缩",
			level:    types.CompressionLevelFast,
			expected: gzip.BestSpeed,
		},
		{
			name:     "最佳压缩",
			level:    types.CompressionLevelBest,
			expected: gzip.BestCompression,
		},
		{
			name:     "哈夫曼编码",
			level:    types.CompressionLevelHuffmanOnly,
			expected: gzip.HuffmanOnly,
		},
		{
			name:     "默认压缩",
			level:    types.CompressionLevelDefault,
			expected: gzip.DefaultCompression,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := GetCompressionLevel(tc.level)
			if result != tc.expected {
				t.Errorf("期望压缩等级为 %d, 实际为 %d", tc.expected, result)
			}
		})
	}
}

// TestGetCompressionLevel_InvalidLevel 测试无效压缩等级
func TestGetCompressionLevel_InvalidLevel(t *testing.T) {
	// 测试未定义的压缩等级，应该返回默认压缩等级
	invalidLevel := types.CompressionLevel(999)
	result := GetCompressionLevel(invalidLevel)
	expected := gzip.DefaultCompression

	if result != expected {
		t.Errorf("无效压缩等级应该返回默认压缩等级 %d, 实际为 %d", expected, result)
	}
}

// TestConfig_FieldModification 测试配置字段修改
func TestConfig_FieldModification(t *testing.T) {
	config := New()

	// 测试修改压缩等级
	config.CompressionLevel = types.CompressionLevelBest
	if config.CompressionLevel != types.CompressionLevelBest {
		t.Errorf("期望压缩等级为 %v, 实际为 %v", types.CompressionLevelBest, config.CompressionLevel)
	}

	// 测试修改覆盖设置
	config.OverwriteExisting = true
	if config.OverwriteExisting != true {
		t.Errorf("期望覆盖设置为 true, 实际为 %v", config.OverwriteExisting)
	}

	// 测试修改路径验证设置
	config.DisablePathValidation = true
	if config.DisablePathValidation != true {
		t.Errorf("期望路径验证设置为 true, 实际为 %v", config.DisablePathValidation)
	}
}

// TestConfig_ProgressInstance 测试进度显示实例
func TestConfig_ProgressInstance(t *testing.T) {
	config := New()

	// 验证进度显示实例不为nil且类型正确
	if config.Progress == nil {
		t.Error("进度显示实例不应该为nil")
	}

	// 测试替换进度显示实例
	newProgress := progress.New()
	config.Progress = newProgress

	if config.Progress != newProgress {
		t.Error("进度显示实例替换失败")
	}
}

// TestConfig_AllCompressionLevels 测试所有压缩等级的映射
func TestConfig_AllCompressionLevels(t *testing.T) {
	// 定义所有已知的压缩等级及其期望的映射值
	testCases := []struct {
		level    types.CompressionLevel
		expected int
	}{
		{types.CompressionLevelNone, gzip.NoCompression},
		{types.CompressionLevelFast, gzip.BestSpeed},
		{types.CompressionLevelBest, gzip.BestCompression},
		{types.CompressionLevelHuffmanOnly, gzip.HuffmanOnly},
		{types.CompressionLevelDefault, gzip.DefaultCompression},
	}

	// 验证每个压缩等级都能正确映射
	for _, tc := range testCases {
		result := GetCompressionLevel(tc.level)
		if result != tc.expected {
			t.Errorf("压缩等级 %v 期望映射值为 %d, 实际为 %d", tc.level, tc.expected, result)
		}
	}
}

// TestConfig_DefaultValues 测试默认值的一致性
func TestConfig_DefaultValues(t *testing.T) {
	config1 := New()
	config2 := New()

	// 验证两个新创建的配置实例具有相同的默认值
	if config1.CompressionLevel != config2.CompressionLevel {
		t.Error("两个新配置实例的压缩等级应该相同")
	}

	if config1.OverwriteExisting != config2.OverwriteExisting {
		t.Error("两个新配置实例的覆盖设置应该相同")
	}

	if config1.DisablePathValidation != config2.DisablePathValidation {
		t.Error("两个新配置实例的路径验证设置应该相同")
	}

	// 验证进度显示实例是独立的
	if config1.Progress == config2.Progress {
		t.Error("两个新配置实例的进度显示实例应该是独立的")
	}
}

// TestConfig_ZeroValue 测试零值配置
func TestConfig_ZeroValue(t *testing.T) {
	var config Config

	// 验证零值配置的字段
	if config.CompressionLevel != types.CompressionLevel(0) {
		t.Errorf("零值配置的压缩等级应该为 0, 实际为 %v", config.CompressionLevel)
	}

	if config.OverwriteExisting != false {
		t.Errorf("零值配置的覆盖设置应该为 false, 实际为 %v", config.OverwriteExisting)
	}

	if config.Progress != nil {
		t.Errorf("零值配置的进度显示实例应该为 nil, 实际为 %v", config.Progress)
	}

	if config.DisablePathValidation != false {
		t.Errorf("零值配置的路径验证设置应该为 false, 实际为 %v", config.DisablePathValidation)
	}
}

// BenchmarkNew 基准测试：创建新配置
func BenchmarkNew(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = New()
	}
}

// BenchmarkGetCompressionLevel 基准测试：获取压缩等级
func BenchmarkGetCompressionLevel(b *testing.B) {
	levels := []types.CompressionLevel{
		types.CompressionLevelNone,
		types.CompressionLevelFast,
		types.CompressionLevelBest,
		types.CompressionLevelHuffmanOnly,
		types.CompressionLevelDefault,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		level := levels[i%len(levels)]
		_ = GetCompressionLevel(level)
	}
}

// TestConfig_Concurrent 测试并发安全性
func TestConfig_Concurrent(t *testing.T) {
	config := New()
	done := make(chan bool, 10)

	// 启动多个goroutine同时访问配置
	for i := 0; i < 10; i++ {
		go func(id int) {
			defer func() { done <- true }()

			// 读取配置
			_ = config.CompressionLevel
			_ = config.OverwriteExisting
			_ = config.DisablePathValidation

			// 修改配置
			config.CompressionLevel = types.CompressionLevelFast
			config.OverwriteExisting = true
			config.DisablePathValidation = true
		}(i)
	}

	// 等待所有goroutine完成
	for i := 0; i < 10; i++ {
		<-done
	}

	// 注意：这个测试主要是检查是否有竞态条件导致的panic
	// 实际的并发安全需要在使用时通过互斥锁等机制保证
}

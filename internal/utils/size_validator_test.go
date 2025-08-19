package utils

import (
	"os"
	"path/filepath"
	"testing"

	"gitee.com/MM-Q/comprx/config"
)

// TestValidateFileSize 测试单个文件大小验证
func TestValidateFileSize(t *testing.T) {
	tests := []struct {
		name        string
		config      *config.Config
		filePath    string
		fileSize    int64
		expectError bool
	}{
		{
			name: "正常文件大小",
			config: &config.Config{
				EnableSizeCheck: true,
				MaxFileSize:     100 * 1024 * 1024, // 100MB
			},
			filePath:    "test.txt",
			fileSize:    50 * 1024 * 1024, // 50MB
			expectError: false,
		},
		{
			name: "文件大小超限",
			config: &config.Config{
				EnableSizeCheck: true,
				MaxFileSize:     100 * 1024 * 1024, // 100MB
			},
			filePath:    "large.txt",
			fileSize:    150 * 1024 * 1024, // 150MB
			expectError: true,
		},
		{
			name: "禁用大小检查",
			config: &config.Config{
				EnableSizeCheck: false,
				MaxFileSize:     100 * 1024 * 1024,
			},
			filePath:    "any.txt",
			fileSize:    200 * 1024 * 1024, // 200MB
			expectError: false,
		},
		{
			name: "零大小文件",
			config: &config.Config{
				EnableSizeCheck: true,
				MaxFileSize:     100 * 1024 * 1024,
			},
			filePath:    "empty.txt",
			fileSize:    0,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateFileSize(tt.config, tt.filePath, tt.fileSize)
			if tt.expectError && err == nil {
				t.Errorf("期望出现错误，但没有错误")
			}
			if !tt.expectError && err != nil {
				t.Errorf("不期望出现错误，但出现了错误: %v", err)
			}
		})
	}
}

// TestValidateCompressionRatio 测试压缩比验证
func TestValidateCompressionRatio(t *testing.T) {
	tests := []struct {
		name           string
		config         *config.Config
		originalSize   int64
		compressedSize int64
		expectError    bool
	}{
		{
			name: "正常压缩比",
			config: &config.Config{
				EnableSizeCheck:     true,
				MaxCompressionRatio: 100.0,
			},
			originalSize:   1000,
			compressedSize: 100, // 10:1 压缩比
			expectError:    false,
		},
		{
			name: "压缩比超限",
			config: &config.Config{
				EnableSizeCheck:     true,
				MaxCompressionRatio: 50.0,
			},
			originalSize:   10000,
			compressedSize: 100, // 100:1 压缩比
			expectError:    true,
		},
		{
			name: "禁用大小检查",
			config: &config.Config{
				EnableSizeCheck:     false,
				MaxCompressionRatio: 10.0,
			},
			originalSize:   10000,
			compressedSize: 100, // 100:1 压缩比
			expectError:    false,
		},
		{
			name: "压缩后大小为零",
			config: &config.Config{
				EnableSizeCheck:     true,
				MaxCompressionRatio: 100.0,
			},
			originalSize:   1000,
			compressedSize: 0,
			expectError:    false, // 压缩后大小为0时跳过检查
		},
		{
			name: "原始大小为零",
			config: &config.Config{
				EnableSizeCheck:     true,
				MaxCompressionRatio: 100.0,
			},
			originalSize:   0,
			compressedSize: 100,
			expectError:    false, // 原始大小为0时压缩比为0
		},
		{
			name: "无压缩情况",
			config: &config.Config{
				EnableSizeCheck:     true,
				MaxCompressionRatio: 100.0,
			},
			originalSize:   1000,
			compressedSize: 1000, // 1:1 无压缩
			expectError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCompressionRatio(tt.config, tt.originalSize, tt.compressedSize)
			if tt.expectError && err == nil {
				t.Errorf("期望出现错误，但没有错误")
			}
			if !tt.expectError && err != nil {
				t.Errorf("不期望出现错误，但出现了错误: %v", err)
			}
		})
	}
}

// TestSizeTracker 测试大小跟踪器
func TestSizeTracker(t *testing.T) {
	cfg := &config.Config{
		EnableSizeCheck: true,
		MaxTotalSize:    1000,
	}

	tracker := NewSizeTracker()

	// 测试添加正常大小
	err := tracker.AddSize(cfg, 300)
	if err != nil {
		t.Errorf("添加300字节时不应该出错: %v", err)
	}
	if tracker.processedSize != 300 {
		t.Errorf("期望处理大小为300，实际为%d", tracker.processedSize)
	}

	// 测试继续添加
	err = tracker.AddSize(cfg, 400)
	if err != nil {
		t.Errorf("添加400字节时不应该出错: %v", err)
	}
	if tracker.processedSize != 700 {
		t.Errorf("期望处理大小为700，实际为%d", tracker.processedSize)
	}

	// 测试超限
	err = tracker.AddSize(cfg, 400) // 总计1100，超过1000限制
	if err == nil {
		t.Errorf("期望出现超限错误，但没有错误")
	}
	if tracker.processedSize != 700 { // 超限时不应该更新大小
		t.Errorf("超限时处理大小不应该更新，期望700，实际为%d", tracker.processedSize)
	}

	// 测试重置
	tracker.Reset()
	if tracker.processedSize != 0 {
		t.Errorf("重置后处理大小应该为0，实际为%d", tracker.processedSize)
	}

	// 测试禁用检查
	disabledCfg := &config.Config{
		EnableSizeCheck: false,
		MaxTotalSize:    100,
	}
	err = tracker.AddSize(disabledCfg, 200) // 超过限制但禁用检查
	if err != nil {
		t.Errorf("禁用检查时不应该出错: %v", err)
	}
}

// TestPreCheckSingleFile 测试单文件预检查
func TestPreCheckSingleFile(t *testing.T) {
	// 创建临时文件
	tempDir, err := os.MkdirTemp("", "size_validator_test")
	if err != nil {
		t.Fatalf("创建临时目录失败: %v", err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	// 创建测试文件
	testFile := filepath.Join(tempDir, "test.txt")
	testContent := make([]byte, 1024) // 1KB
	for i := range testContent {
		testContent[i] = 'A'
	}
	err = os.WriteFile(testFile, testContent, 0644)
	if err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}

	tests := []struct {
		name        string
		config      *config.Config
		filePath    string
		expectError bool
	}{
		{
			name: "正常文件",
			config: &config.Config{
				EnableSizeCheck: true,
				MaxFileSize:     2048, // 2KB
				MaxTotalSize:    4096, // 4KB 总大小限制
			},
			filePath:    testFile,
			expectError: false,
		},
		{
			name: "文件过大",
			config: &config.Config{
				EnableSizeCheck: true,
				MaxFileSize:     512, // 512B
				MaxTotalSize:    4096,
			},
			filePath:    testFile,
			expectError: true,
		},
		{
			name: "文件不存在",
			config: &config.Config{
				EnableSizeCheck: true,
				MaxFileSize:     2048,
				MaxTotalSize:    4096,
			},
			filePath:    filepath.Join(tempDir, "nonexistent.txt"),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := PreCheckSingleFile(tt.config, tt.filePath)
			if tt.expectError && err == nil {
				t.Errorf("期望出现错误，但没有错误")
			}
			if !tt.expectError && err != nil {
				t.Errorf("不期望出现错误，但出现了错误: %v", err)
			}
		})
	}
}

// TestPreCheckDirectorySize 测试目录预检查
func TestPreCheckDirectorySize(t *testing.T) {
	// 创建临时目录结构
	tempDir, err := os.MkdirTemp("", "size_validator_test")
	if err != nil {
		t.Fatalf("创建临时目录失败: %v", err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	// 创建测试文件结构
	files := map[string]int{
		"file1.txt":        500,
		"file2.txt":        300,
		"subdir/file3.txt": 200,
		"subdir/file4.txt": 400,
	}

	for filePath, size := range files {
		fullPath := filepath.Join(tempDir, filePath)
		dir := filepath.Dir(fullPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("创建目录失败: %v", err)
		}

		content := make([]byte, size)
		for i := range content {
			content[i] = 'A'
		}
		if err := os.WriteFile(fullPath, content, 0644); err != nil {
			t.Fatalf("创建文件失败: %v", err)
		}
	}

	tests := []struct {
		name        string
		config      *config.Config
		dirPath     string
		expectError bool
	}{
		{
			name: "正常目录",
			config: &config.Config{
				EnableSizeCheck: true,
				MaxFileSize:     1000, // 1KB 单文件限制
				MaxTotalSize:    2000, // 2KB 总大小限制
			},
			dirPath:     tempDir,
			expectError: false,
		},
		{
			name: "总大小超限",
			config: &config.Config{
				EnableSizeCheck: true,
				MaxFileSize:     1000, // 1KB 单文件限制
				MaxTotalSize:    1000, // 1KB 总大小限制（实际总大小1400）
			},
			dirPath:     tempDir,
			expectError: true,
		},
		{
			name: "单文件超限",
			config: &config.Config{
				EnableSizeCheck: true,
				MaxFileSize:     400, // 400B 单文件限制（file1.txt为500B）
				MaxTotalSize:    2000,
			},
			dirPath:     tempDir,
			expectError: true,
		},
		{
			name: "目录不存在",
			config: &config.Config{
				EnableSizeCheck: true,
				MaxFileSize:     1000,
				MaxTotalSize:    2000,
			},
			dirPath:     "C:\\absolutely\\nonexistent\\directory\\path\\that\\should\\not\\exist",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := PreCheckDirectorySize(tt.config, tt.dirPath)
			if tt.expectError && err == nil {
				t.Errorf("期望出现错误，但没有错误")
			}
			if !tt.expectError && err != nil {
				t.Errorf("不期望出现错误，但出现了错误: %v", err)
			}
		})
	}
}

// BenchmarkValidateFileSize 基准测试文件大小验证
func BenchmarkValidateFileSize(b *testing.B) {
	config := &config.Config{
		EnableSizeCheck: true,
		MaxFileSize:     100 * 1024 * 1024,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ValidateFileSize(config, "test.txt", 50*1024*1024)
	}
}

// BenchmarkValidateCompressionRatio 基准测试压缩比验证
func BenchmarkValidateCompressionRatio(b *testing.B) {
	config := &config.Config{
		EnableSizeCheck:     true,
		MaxCompressionRatio: 100.0,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ValidateCompressionRatio(config, 10000, 100)
	}
}

// BenchmarkSizeTracker 基准测试大小跟踪器
func BenchmarkSizeTracker(b *testing.B) {
	config := &config.Config{
		EnableSizeCheck: true,
		MaxTotalSize:    1024 * 1024 * 1024, // 1GB
	}

	tracker := NewSizeTracker()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tracker.Reset()
		for j := 0; j < 100; j++ {
			_ = tracker.AddSize(config, 1024) // 每次添加1KB
		}
	}
}

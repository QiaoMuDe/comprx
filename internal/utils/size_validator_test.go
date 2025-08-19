package utils

import (
	"os"
	"path/filepath"
	"testing"

	"gitee.com/MM-Q/comprx/config"
)

func TestValidateFileSize(t *testing.T) {
	tests := []struct {
		name        string
		config      *config.Config
		filePath    string
		size        int64
		expectError bool
	}{
		{
			name: "文件大小在限制内",
			config: &config.Config{
				EnableSizeCheck: true,
				MaxFileSize:     1024,
			},
			filePath:    "test.txt",
			size:        512,
			expectError: false,
		},
		{
			name: "文件大小超过限制",
			config: &config.Config{
				EnableSizeCheck: true,
				MaxFileSize:     1024,
			},
			filePath:    "large.txt",
			size:        2048,
			expectError: true,
		},
		{
			name: "禁用大小检查",
			config: &config.Config{
				EnableSizeCheck: false,
				MaxFileSize:     1024,
			},
			filePath:    "large.txt",
			size:        2048,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateFileSize(tt.config, tt.filePath, tt.size)
			if tt.expectError && err == nil {
				t.Errorf("期望出现错误，但没有错误发生")
			}
			if !tt.expectError && err != nil {
				t.Errorf("不期望出现错误，但发生了错误: %v", err)
			}
		})
	}
}

func TestValidateTotalSize(t *testing.T) {
	tests := []struct {
		name           string
		config         *config.Config
		currentTotal   int64
		additionalSize int64
		expectError    bool
	}{
		{
			name: "总大小在限制内",
			config: &config.Config{
				EnableSizeCheck: true,
				MaxTotalSize:    2048,
			},
			currentTotal:   1000,
			additionalSize: 500,
			expectError:    false,
		},
		{
			name: "总大小超过限制",
			config: &config.Config{
				EnableSizeCheck: true,
				MaxTotalSize:    2048,
			},
			currentTotal:   1500,
			additionalSize: 1000,
			expectError:    true,
		},
		{
			name: "禁用大小检查",
			config: &config.Config{
				EnableSizeCheck: false,
				MaxTotalSize:    2048,
			},
			currentTotal:   1500,
			additionalSize: 1000,
			expectError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTotalSize(tt.config, tt.currentTotal, tt.additionalSize)
			if tt.expectError && err == nil {
				t.Errorf("期望出现错误，但没有错误发生")
			}
			if !tt.expectError && err != nil {
				t.Errorf("不期望出现错误，但发生了错误: %v", err)
			}
		})
	}
}

func TestPreCheckDirectorySize(t *testing.T) {
	// 创建临时目录和文件进行测试
	tempDir := t.TempDir()

	// 创建测试文件
	testFile := filepath.Join(tempDir, "test.txt")
	content := []byte("hello world")
	if err := os.WriteFile(testFile, content, 0644); err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}

	tests := []struct {
		name        string
		config      *config.Config
		dirPath     string
		expectError bool
	}{
		{
			name: "目录大小在限制内",
			config: &config.Config{
				EnableSizeCheck: true,
				MaxFileSize:     1024,
				MaxTotalSize:    2048,
			},
			dirPath:     tempDir,
			expectError: false,
		},
		{
			name: "目录不存在",
			config: &config.Config{
				EnableSizeCheck: true,
				MaxFileSize:     1024,
				MaxTotalSize:    2048,
			},
			dirPath:     "/nonexistent/path",
			expectError: true,
		},
		{
			name: "禁用大小检查",
			config: &config.Config{
				EnableSizeCheck: false,
				MaxFileSize:     1,
				MaxTotalSize:    1,
			},
			dirPath:     tempDir,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			size, err := PreCheckDirectorySize(tt.config, tt.dirPath)
			if tt.expectError && err == nil {
				t.Errorf("期望出现错误，但没有错误发生")
			}
			if !tt.expectError && err != nil {
				t.Errorf("不期望出现错误，但发生了错误: %v", err)
			}
			if !tt.expectError && size < 0 {
				t.Errorf("返回的大小不应该为负数: %d", size)
			}
		})
	}
}

func TestPreCheckSingleFile(t *testing.T) {
	// 创建临时文件进行测试
	tempFile := filepath.Join(t.TempDir(), "test.txt")
	content := []byte("hello world")
	if err := os.WriteFile(tempFile, content, 0644); err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}

	tests := []struct {
		name        string
		config      *config.Config
		filePath    string
		expectError bool
	}{
		{
			name: "文件大小在限制内",
			config: &config.Config{
				EnableSizeCheck: true,
				MaxFileSize:     1024,
				MaxTotalSize:    2048,
			},
			filePath:    tempFile,
			expectError: false,
		},
		{
			name: "文件不存在",
			config: &config.Config{
				EnableSizeCheck: true,
				MaxFileSize:     1024,
				MaxTotalSize:    2048,
			},
			filePath:    "/nonexistent/file.txt",
			expectError: true,
		},
		{
			name: "禁用大小检查",
			config: &config.Config{
				EnableSizeCheck: false,
				MaxFileSize:     1,
				MaxTotalSize:    1,
			},
			filePath:    tempFile,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			size, err := PreCheckSingleFile(tt.config, tt.filePath)
			if tt.expectError && err == nil {
				t.Errorf("期望出现错误，但没有错误发生")
			}
			if !tt.expectError && err != nil {
				t.Errorf("不期望出现错误，但发生了错误: %v", err)
			}
			if !tt.expectError && size < 0 {
				t.Errorf("返回的大小不应该为负数: %d", size)
			}
		})
	}
}

func TestSizeTracker(t *testing.T) {
	tracker := NewSizeTracker()

	// 测试初始状态
	if tracker.GetProcessedSize() != 0 {
		t.Errorf("初始处理大小应该为0，实际为 %d", tracker.GetProcessedSize())
	}

	config := &config.Config{
		EnableSizeCheck: true,
		MaxTotalSize:    1000,
	}

	// 添加大小
	err := tracker.AddSize(config, 500)
	if err != nil {
		t.Errorf("添加大小失败: %v", err)
	}
	if tracker.GetProcessedSize() != 500 {
		t.Errorf("处理大小应该为500，实际为 %d", tracker.GetProcessedSize())
	}

	// 添加更多大小，但不超过限制
	err = tracker.AddSize(config, 400)
	if err != nil {
		t.Errorf("添加大小失败: %v", err)
	}
	if tracker.GetProcessedSize() != 900 {
		t.Errorf("处理大小应该为900，实际为 %d", tracker.GetProcessedSize())
	}

	// 尝试添加超过限制的大小
	err = tracker.AddSize(config, 200)
	if err == nil {
		t.Errorf("期望出现错误（超过总大小限制），但没有错误发生")
	}

	// 重置跟踪器
	tracker.Reset()
	if tracker.GetProcessedSize() != 0 {
		t.Errorf("重置后处理大小应该为0，实际为 %d", tracker.GetProcessedSize())
	}
}

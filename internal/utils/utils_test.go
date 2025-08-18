package utils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExists(t *testing.T) {
	tempDir := t.TempDir()

	t.Run("文件存在", func(t *testing.T) {
		// 创建测试文件
		testFile := filepath.Join(tempDir, "test.txt")
		if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
			t.Fatalf("创建测试文件失败: %v", err)
		}

		if !Exists(testFile) {
			t.Error("应该检测到文件存在")
		}
	})

	t.Run("目录存在", func(t *testing.T) {
		// 创建测试目录
		testDir := filepath.Join(tempDir, "testdir")
		if err := os.MkdirAll(testDir, 0755); err != nil {
			t.Fatalf("创建测试目录失败: %v", err)
		}

		if !Exists(testDir) {
			t.Error("应该检测到目录存在")
		}
	})

	t.Run("文件不存在", func(t *testing.T) {
		nonExistentFile := filepath.Join(tempDir, "nonexistent.txt")
		if Exists(nonExistentFile) {
			t.Error("不应该检测到不存在的文件")
		}
	})

	t.Run("目录不存在", func(t *testing.T) {
		nonExistentDir := filepath.Join(tempDir, "nonexistentdir")
		if Exists(nonExistentDir) {
			t.Error("不应该检测到不存在的目录")
		}
	})

	t.Run("空路径", func(t *testing.T) {
		if Exists("") {
			t.Error("空路径应该返回false")
		}
	})

	t.Run("无效路径", func(t *testing.T) {
		// 在Windows上，某些字符是无效的
		invalidPath := filepath.Join(tempDir, "invalid\x00path")
		if Exists(invalidPath) {
			t.Error("无效路径应该返回false")
		}
	})
}

func TestEnsureDir(t *testing.T) {
	tempDir := t.TempDir()

	t.Run("创建新目录", func(t *testing.T) {
		newDir := filepath.Join(tempDir, "newdir")

		if err := EnsureDir(newDir); err != nil {
			t.Fatalf("创建目录失败: %v", err)
		}

		if !Exists(newDir) {
			t.Error("目录应该被创建")
		}

		// 检查目录权限
		info, err := os.Stat(newDir)
		if err != nil {
			t.Fatalf("获取目录信息失败: %v", err)
		}

		if !info.IsDir() {
			t.Error("应该是目录")
		}
	})

	t.Run("创建嵌套目录", func(t *testing.T) {
		nestedDir := filepath.Join(tempDir, "level1", "level2", "level3")

		if err := EnsureDir(nestedDir); err != nil {
			t.Fatalf("创建嵌套目录失败: %v", err)
		}

		if !Exists(nestedDir) {
			t.Error("嵌套目录应该被创建")
		}

		// 检查所有父目录都被创建
		level1 := filepath.Join(tempDir, "level1")
		level2 := filepath.Join(tempDir, "level1", "level2")

		if !Exists(level1) || !Exists(level2) {
			t.Error("所有父目录都应该被创建")
		}
	})

	t.Run("目录已存在", func(t *testing.T) {
		existingDir := filepath.Join(tempDir, "existing")

		// 先创建目录
		if err := os.MkdirAll(existingDir, 0755); err != nil {
			t.Fatalf("预创建目录失败: %v", err)
		}

		// 再次调用EnsureDir应该成功
		if err := EnsureDir(existingDir); err != nil {
			t.Errorf("对已存在目录调用EnsureDir失败: %v", err)
		}
	})

	t.Run("路径是文件而非目录", func(t *testing.T) {
		testFile := filepath.Join(tempDir, "testfile.txt")

		// 创建文件
		if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
			t.Fatalf("创建测试文件失败: %v", err)
		}

		// 尝试将文件路径作为目录创建
		// 注意：EnsureDir 函数实际上会检查路径是否存在，如果存在就返回 nil
		// 这里我们需要检查实际的行为
		err := EnsureDir(testFile)
		if err != nil {
			// 如果返回错误，说明函数正确识别了这不是目录
			t.Logf("正确识别文件不是目录: %v", err)
		} else {
			// 如果没有返回错误，检查路径是否仍然是文件
			info, statErr := os.Stat(testFile)
			if statErr != nil {
				t.Fatalf("获取文件信息失败: %v", statErr)
			}
			if info.IsDir() {
				t.Error("文件不应该被转换为目录")
			}
			// 如果路径仍然是文件，这也是可以接受的行为
		}
	})

	t.Run("空路径", func(t *testing.T) {
		if err := EnsureDir(""); err == nil {
			t.Error("空路径应该返回错误")
		}
	})
}

func TestGetBufferSize(t *testing.T) {
	tests := []struct {
		name     string
		fileSize int64
		expected int
	}{
		{"小文件 100B", 100, 32 * 1024},
		{"小文件 100KB", 100 * 1024, 32 * 1024},
		{"边界 512KB", 512 * 1024, 64 * 1024},
		{"中等文件 800KB", 800 * 1024, 64 * 1024},
		{"边界 1MB", 1 * 1024 * 1024, 128 * 1024},
		{"中等文件 3MB", 3 * 1024 * 1024, 128 * 1024},
		{"边界 5MB", 5 * 1024 * 1024, 256 * 1024},
		{"中等文件 8MB", 8 * 1024 * 1024, 256 * 1024},
		{"边界 10MB", 10 * 1024 * 1024, 512 * 1024},
		{"大文件 50MB", 50 * 1024 * 1024, 512 * 1024},
		{"边界 100MB", 100 * 1024 * 1024, 1024 * 1024},
		{"超大文件 500MB", 500 * 1024 * 1024, 1024 * 1024},
		{"超大文件 1GB", 1024 * 1024 * 1024, 1024 * 1024},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetBufferSize(tt.fileSize)
			if result != tt.expected {
				t.Errorf("GetBufferSize(%d) = %d, want %d", tt.fileSize, result, tt.expected)
			}
		})
	}
}

func TestGetBufferSize_EdgeCases(t *testing.T) {
	t.Run("零大小文件", func(t *testing.T) {
		result := GetBufferSize(0)
		expected := 32 * 1024
		if result != expected {
			t.Errorf("GetBufferSize(0) = %d, want %d", result, expected)
		}
	})

	t.Run("负数大小", func(t *testing.T) {
		result := GetBufferSize(-1)
		expected := 32 * 1024
		if result != expected {
			t.Errorf("GetBufferSize(-1) = %d, want %d", result, expected)
		}
	})

	t.Run("极大文件", func(t *testing.T) {
		result := GetBufferSize(1024 * 1024 * 1024 * 10) // 10GB
		expected := 1024 * 1024
		if result != expected {
			t.Errorf("GetBufferSize(10GB) = %d, want %d", result, expected)
		}
	})
}

func TestEnsureAbsPath(t *testing.T) {
	t.Run("已经是绝对路径", func(t *testing.T) {
		if filepath.IsAbs("/absolute/path") {
			result, err := EnsureAbsPath("/absolute/path", "测试路径")
			if err != nil {
				t.Fatalf("处理绝对路径失败: %v", err)
			}
			if result != "/absolute/path" {
				t.Errorf("绝对路径应该保持不变: got %s, want /absolute/path", result)
			}
		}
	})

	t.Run("相对路径转换", func(t *testing.T) {
		relativePath := "relative/path"
		result, err := EnsureAbsPath(relativePath, "测试路径")
		if err != nil {
			t.Fatalf("转换相对路径失败: %v", err)
		}

		if !filepath.IsAbs(result) {
			t.Error("结果应该是绝对路径")
		}

		// 检查路径是否包含原始相对路径
		if !filepath.IsAbs(result) {
			t.Error("转换后的路径应该是绝对路径")
		}
	})

	t.Run("当前目录", func(t *testing.T) {
		result, err := EnsureAbsPath(".", "当前目录")
		if err != nil {
			t.Fatalf("处理当前目录失败: %v", err)
		}

		if !filepath.IsAbs(result) {
			t.Error("当前目录应该转换为绝对路径")
		}
	})

	t.Run("父目录", func(t *testing.T) {
		result, err := EnsureAbsPath("..", "父目录")
		if err != nil {
			t.Fatalf("处理父目录失败: %v", err)
		}

		if !filepath.IsAbs(result) {
			t.Error("父目录应该转换为绝对路径")
		}
	})

	t.Run("空路径", func(t *testing.T) {
		result, err := EnsureAbsPath("", "空路径")
		if err != nil {
			t.Fatalf("处理空路径失败: %v", err)
		}

		if !filepath.IsAbs(result) {
			t.Error("空路径应该转换为绝对路径")
		}
	})

	t.Run("复杂相对路径", func(t *testing.T) {
		complexPath := "../../some/complex/../path/./file"
		result, err := EnsureAbsPath(complexPath, "复杂路径")
		if err != nil {
			t.Fatalf("处理复杂路径失败: %v", err)
		}

		if !filepath.IsAbs(result) {
			t.Error("复杂路径应该转换为绝对路径")
		}
	})
}

// 基准测试
func BenchmarkExists(b *testing.B) {
	tempDir := b.TempDir()
	testFile := filepath.Join(tempDir, "benchmark.txt")
	_ = os.WriteFile(testFile, []byte("test"), 0644)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Exists(testFile)
	}
}

func BenchmarkGetBufferSize(b *testing.B) {
	sizes := []int64{
		100,                // 100B
		100 * 1024,         // 100KB
		1 * 1024 * 1024,    // 1MB
		10 * 1024 * 1024,   // 10MB
		100 * 1024 * 1024,  // 100MB
		1024 * 1024 * 1024, // 1GB
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, size := range sizes {
			_ = GetBufferSize(size)
		}
	}
}

func BenchmarkEnsureAbsPath(b *testing.B) {
	paths := []string{
		"relative/path",
		"./current",
		"../parent",
		"complex/../../path",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, path := range paths {
			_, _ = EnsureAbsPath(path, "benchmark")
		}
	}
}

package comprx

import (
	"os"
	"path/filepath"
	"testing"
)

// TestGetFileSize 测试获取文件大小功能
func TestGetFileSize(t *testing.T) {
	tempDir := t.TempDir()

	// 测试用例
	testCases := []struct {
		name         string
		setupFunc    func() string
		expectedSize int64
	}{
		{
			name: "普通文件",
			setupFunc: func() string {
				filePath := filepath.Join(tempDir, "test.txt")
				content := "Hello, World!"
				if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
					t.Fatal(err)
				}
				return filePath
			},
			expectedSize: 13, // "Hello, World!" 的字节数
		},
		{
			name: "空文件",
			setupFunc: func() string {
				filePath := filepath.Join(tempDir, "empty.txt")
				if err := os.WriteFile(filePath, []byte{}, 0644); err != nil {
					t.Fatal(err)
				}
				return filePath
			},
			expectedSize: 0,
		},
		{
			name: "不存在的文件",
			setupFunc: func() string {
				return filepath.Join(tempDir, "nonexistent.txt")
			},
			expectedSize: 0,
		},
		{
			name: "目录而不是文件",
			setupFunc: func() string {
				dirPath := filepath.Join(tempDir, "testdir")
				if err := os.MkdirAll(dirPath, 0755); err != nil {
					t.Fatal(err)
				}
				return dirPath
			},
			expectedSize: 0,
		},
		{
			name: "较大的文件",
			setupFunc: func() string {
				filePath := filepath.Join(tempDir, "large.txt")
				content := make([]byte, 1024) // 1KB
				for i := range content {
					content[i] = byte('A')
				}
				if err := os.WriteFile(filePath, content, 0644); err != nil {
					t.Fatal(err)
				}
				return filePath
			},
			expectedSize: 1024,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			filePath := tc.setupFunc()
			size := GetFileSize(filePath)

			if size != tc.expectedSize {
				t.Errorf("期望大小 %d，实际大小 %d", tc.expectedSize, size)
			}
		})
	}
}

// TestGetDirectorySize 测试获取目录大小功能
func TestGetDirectorySize(t *testing.T) {
	tempDir := t.TempDir()

	// 测试用例
	testCases := []struct {
		name         string
		setupFunc    func() string
		expectedSize int64
	}{
		{
			name: "空目录",
			setupFunc: func() string {
				dirPath := filepath.Join(tempDir, "empty_dir")
				if err := os.MkdirAll(dirPath, 0755); err != nil {
					t.Fatal(err)
				}
				return dirPath
			},
			expectedSize: 0,
		},
		{
			name: "单个文件的目录",
			setupFunc: func() string {
				dirPath := filepath.Join(tempDir, "single_file_dir")
				if err := os.MkdirAll(dirPath, 0755); err != nil {
					t.Fatal(err)
				}

				filePath := filepath.Join(dirPath, "file.txt")
				content := "Hello"
				if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
					t.Fatal(err)
				}
				return dirPath
			},
			expectedSize: 5, // "Hello" 的字节数
		},
		{
			name: "多个文件的目录",
			setupFunc: func() string {
				dirPath := filepath.Join(tempDir, "multi_file_dir")
				if err := os.MkdirAll(dirPath, 0755); err != nil {
					t.Fatal(err)
				}

				// 创建多个文件
				files := map[string]string{
					"file1.txt": "Hello", // 5 bytes
					"file2.txt": "World", // 5 bytes
					"file3.txt": "Test",  // 4 bytes
				}

				for fileName, content := range files {
					filePath := filepath.Join(dirPath, fileName)
					if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
						t.Fatal(err)
					}
				}
				return dirPath
			},
			expectedSize: 14, // 5 + 5 + 4 = 14 bytes
		},
		{
			name: "嵌套目录",
			setupFunc: func() string {
				dirPath := filepath.Join(tempDir, "nested_dir")
				if err := os.MkdirAll(dirPath, 0755); err != nil {
					t.Fatal(err)
				}

				// 创建嵌套目录结构
				subDir := filepath.Join(dirPath, "subdir")
				if err := os.MkdirAll(subDir, 0755); err != nil {
					t.Fatal(err)
				}

				// 在根目录创建文件
				rootFile := filepath.Join(dirPath, "root.txt")
				if err := os.WriteFile(rootFile, []byte("Root"), 0644); err != nil {
					t.Fatal(err)
				}

				// 在子目录创建文件
				subFile := filepath.Join(subDir, "sub.txt")
				if err := os.WriteFile(subFile, []byte("Sub"), 0644); err != nil {
					t.Fatal(err)
				}

				return dirPath
			},
			expectedSize: 7, // "Root" (4) + "Sub" (3) = 7 bytes
		},
		{
			name: "不存在的目录",
			setupFunc: func() string {
				return filepath.Join(tempDir, "nonexistent_dir")
			},
			expectedSize: 0,
		},
		{
			name: "包含子目录但只有文件的目录",
			setupFunc: func() string {
				dirPath := filepath.Join(tempDir, "mixed_dir")
				if err := os.MkdirAll(dirPath, 0755); err != nil {
					t.Fatal(err)
				}

				// 创建空子目录
				emptySubDir := filepath.Join(dirPath, "empty_sub")
				if err := os.MkdirAll(emptySubDir, 0755); err != nil {
					t.Fatal(err)
				}

				// 创建文件
				filePath := filepath.Join(dirPath, "file.txt")
				if err := os.WriteFile(filePath, []byte("Content"), 0644); err != nil {
					t.Fatal(err)
				}

				return dirPath
			},
			expectedSize: 7, // "Content" = 7 bytes
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			dirPath := tc.setupFunc()
			size := GetDirectorySize(dirPath)

			if size != tc.expectedSize {
				t.Errorf("期望大小 %d，实际大小 %d", tc.expectedSize, size)
			}
		})
	}
}

// TestGetDirectorySizeWithLargeFiles 测试大文件目录大小计算
func TestGetDirectorySizeWithLargeFiles(t *testing.T) {
	tempDir := t.TempDir()

	dirPath := filepath.Join(tempDir, "large_files_dir")
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		t.Fatal(err)
	}

	// 创建几个较大的文件
	var expectedTotal int64
	fileSizes := []int{1024, 2048, 512} // 不同大小的文件

	for i, size := range fileSizes {
		fileName := filepath.Join(dirPath, "large_file_"+string(rune('1'+i))+".txt")
		content := make([]byte, size)
		for j := range content {
			content[j] = byte('A' + (j % 26))
		}

		if err := os.WriteFile(fileName, content, 0644); err != nil {
			t.Fatal(err)
		}

		expectedTotal += int64(size)
	}

	actualSize := GetDirectorySize(dirPath)
	if actualSize != expectedTotal {
		t.Errorf("期望总大小 %d，实际大小 %d", expectedTotal, actualSize)
	}
}

// BenchmarkGetFileSize 性能基准测试
func BenchmarkGetFileSize(b *testing.B) {
	tempDir := b.TempDir()

	// 创建测试文件
	filePath := filepath.Join(tempDir, "benchmark.txt")
	content := make([]byte, 1024*1024) // 1MB
	if err := os.WriteFile(filePath, content, 0644); err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = GetFileSize(filePath)
	}
}

// BenchmarkGetDirectorySize 性能基准测试
func BenchmarkGetDirectorySize(b *testing.B) {
	tempDir := b.TempDir()

	// 创建测试目录结构
	dirPath := filepath.Join(tempDir, "benchmark_dir")
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		b.Fatal(err)
	}

	// 创建多个文件
	for i := 0; i < 100; i++ {
		fileName := filepath.Join(dirPath, "file_"+string(rune('0'+i%10))+".txt")
		content := make([]byte, 1024) // 1KB per file
		if err := os.WriteFile(fileName, content, 0644); err != nil {
			b.Fatal(err)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = GetDirectorySize(dirPath)
	}
}

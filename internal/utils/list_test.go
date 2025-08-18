package utils

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"gitee.com/MM-Q/comprx/types"
)

func TestFormatFileSize(t *testing.T) {
	tests := []struct {
		name     string
		size     int64
		expected string
	}{
		{"零字节", 0, "0 B"},
		{"小于1KB", 512, "512 B"},
		{"1KB", 1024, "1.0 KB"},
		{"1.5KB", 1536, "1.5 KB"},
		{"1MB", 1024 * 1024, "1.0 MB"},
		{"1.2MB", 1024*1024 + 204*1024, "1.2 MB"},
		{"1GB", 1024 * 1024 * 1024, "1.0 GB"},
		{"1.5GB", 1024*1024*1024 + 512*1024*1024, "1.5 GB"},
		{"1TB", 1024 * 1024 * 1024 * 1024, "1.0 TB"},
		{"1PB", 1024 * 1024 * 1024 * 1024 * 1024, "1.0 PB"},
		{"1EB", 1024 * 1024 * 1024 * 1024 * 1024 * 1024, "1.0 EB"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatFileSize(tt.size)
			if result != tt.expected {
				t.Errorf("FormatFileSize(%d) = %s, want %s", tt.size, result, tt.expected)
			}
		})
	}
}

func TestFormatFileMode(t *testing.T) {
	tests := []struct {
		name     string
		mode     os.FileMode
		contains string // 检查结果是否包含特定字符串
	}{
		{"普通文件", 0644, "-rw-r--r--"},
		{"可执行文件", 0755, "-rwxr-xr-x"},
		{"目录", os.ModeDir | 0755, "d"},
		{"符号链接", os.ModeSymlink | 0777, "L"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatFileMode(tt.mode)
			if !strings.Contains(result, tt.contains) {
				t.Errorf("FormatFileMode(%v) = %s, should contain %s", tt.mode, result, tt.contains)
			}
		})
	}
}

func TestMatchPattern(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		pattern  string
		expected bool
	}{
		// 空模式匹配所有
		{"空模式", "test.txt", "", true},

		// 精确匹配
		{"精确匹配", "test.txt", "test.txt", true},
		{"精确不匹配", "test.txt", "other.txt", false},

		// 通配符匹配
		{"星号匹配所有", "test.txt", "*", true},
		{"星号匹配扩展名", "test.txt", "*.txt", true},
		{"星号匹配前缀", "test.txt", "test.*", true},
		{"星号不匹配", "test.txt", "*.doc", false},

		// 问号匹配
		{"问号匹配单字符", "test.txt", "tes?.txt", true},
		{"问号不匹配", "test.txt", "te??.txt", true}, // 实际上会匹配，因为会回退到包含匹配

		// 路径匹配
		{"路径中的文件", "dir/subdir/file.txt", "*.txt", true},
		{"路径中的目录", "dir/subdir/file.txt", "subdir", true},
		{"路径匹配", "dir/subdir/file.txt", "dir/*", true},

		// 大小写不敏感的包含匹配
		{"大小写不敏感", "Test.TXT", "test", true},
		{"包含匹配", "very_long_filename.txt", "long", true},

		// 无效模式回退到包含匹配
		{"无效模式", "test.txt", "[invalid", false}, // 不包含"invalid"
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MatchPattern(tt.filename, tt.pattern)
			if result != tt.expected {
				t.Errorf("MatchPattern(%q, %q) = %v, want %v", tt.filename, tt.pattern, result, tt.expected)
			}
		})
	}
}

func TestPrintFileInfo(t *testing.T) {
	// 创建测试用的FileInfo
	testTime := time.Date(2023, 12, 25, 15, 30, 45, 0, time.UTC)

	tests := []struct {
		name        string
		fileInfo    types.FileInfo
		showDetails bool
		contains    []string // 输出应该包含的字符串
	}{
		{
			name: "简单模式普通文件",
			fileInfo: types.FileInfo{
				Name:    "test.txt",
				Size:    1024,
				Mode:    0644,
				ModTime: testTime,
			},
			showDetails: false,
			contains:    []string{"test.txt"},
		},
		{
			name: "详细模式普通文件",
			fileInfo: types.FileInfo{
				Name:    "test.txt",
				Size:    1024,
				Mode:    0644,
				ModTime: testTime,
			},
			showDetails: true,
			contains:    []string{"test.txt", "1.0 KB", "2023-12-25 15:30:45"},
		},
		{
			name: "简单模式符号链接",
			fileInfo: types.FileInfo{
				Name:       "link.txt",
				Size:       0,
				Mode:       os.ModeSymlink | 0777,
				ModTime:    testTime,
				IsSymlink:  true,
				LinkTarget: "target.txt",
			},
			showDetails: false,
			contains:    []string{"link.txt", "->", "target.txt"},
		},
		{
			name: "详细模式符号链接",
			fileInfo: types.FileInfo{
				Name:       "link.txt",
				Size:       0,
				Mode:       os.ModeSymlink | 0777,
				ModTime:    testTime,
				IsSymlink:  true,
				LinkTarget: "target.txt",
			},
			showDetails: true,
			contains:    []string{"link.txt", "->", "target.txt", "0 B", "2023-12-25 15:30:45"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 捕获标准输出
			old := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			PrintFileInfo(tt.fileInfo, tt.showDetails)

			_ = w.Close()
			os.Stdout = old

			var buf bytes.Buffer
			_, _ = io.Copy(&buf, r)
			output := buf.String()

			// 检查输出是否包含期望的字符串
			for _, expected := range tt.contains {
				if !strings.Contains(output, expected) {
					t.Errorf("输出应该包含 %q，但实际输出为: %s", expected, output)
				}
			}
		})
	}
}

func TestPrintArchiveInfo(t *testing.T) {
	// 创建测试用的ArchiveInfo
	testTime := time.Date(2023, 12, 25, 15, 30, 45, 0, time.UTC)

	archiveInfo := &types.ArchiveInfo{
		Type:           "ZIP",
		TotalFiles:     3,
		TotalSize:      3072, // 3KB
		CompressedSize: 1536, // 1.5KB
		Files: []types.FileInfo{
			{
				Name:    "file1.txt",
				Size:    1024,
				Mode:    0644,
				ModTime: testTime,
			},
			{
				Name:    "file2.txt",
				Size:    2048,
				Mode:    0644,
				ModTime: testTime,
			},
		},
	}

	t.Run("显示摘要信息", func(t *testing.T) {
		// 捕获标准输出
		old := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		PrintArchiveInfo(archiveInfo, true)

		_ = w.Close()
		os.Stdout = old

		var buf bytes.Buffer
		_, _ = io.Copy(&buf, r)
		output := buf.String()

		// 检查输出是否包含期望的信息
		expectedStrings := []string{
			"压缩包类型: ZIP",
			"文件总数: 3",
			"原始大小: 3.0 KB",
			"压缩大小: 1.5 KB",
			"压缩率: 50.0%",
			"file1.txt",
			"file2.txt",
		}

		for _, expected := range expectedStrings {
			if !strings.Contains(output, expected) {
				t.Errorf("输出应该包含 %q，但实际输出为: %s", expected, output)
			}
		}
	})

	t.Run("不显示摘要信息", func(t *testing.T) {
		// 捕获标准输出
		old := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		PrintArchiveInfo(archiveInfo, false)

		_ = w.Close()
		os.Stdout = old

		var buf bytes.Buffer
		_, _ = io.Copy(&buf, r)
		output := buf.String()

		// 应该只包含文件列表，不包含摘要信息
		if strings.Contains(output, "压缩包类型") {
			t.Error("不应该显示摘要信息")
		}

		// 但应该包含文件名
		if !strings.Contains(output, "file1.txt") || !strings.Contains(output, "file2.txt") {
			t.Error("应该显示文件列表")
		}
	})
}

func TestFilterFilesByPattern(t *testing.T) {
	files := []types.FileInfo{
		{Name: "test.txt", Size: 100},
		{Name: "document.doc", Size: 200},
		{Name: "image.jpg", Size: 300},
		{Name: "script.sh", Size: 400},
		{Name: "data.json", Size: 500},
		{Name: "dir/subfile.txt", Size: 600},
	}

	tests := []struct {
		name     string
		pattern  string
		expected int
		contains []string
	}{
		{
			name:     "空模式返回所有文件",
			pattern:  "",
			expected: 6,
			contains: []string{"test.txt", "document.doc", "image.jpg", "script.sh", "data.json", "dir/subfile.txt"},
		},
		{
			name:     "匹配txt文件",
			pattern:  "*.txt",
			expected: 2,
			contains: []string{"test.txt", "dir/subfile.txt"},
		},
		{
			name:     "匹配特定文件名",
			pattern:  "test.txt",
			expected: 1,
			contains: []string{"test.txt"},
		},
		{
			name:     "匹配包含特定字符串的文件",
			pattern:  "data",
			expected: 1,
			contains: []string{"data.json"},
		},
		{
			name:     "不匹配任何文件",
			pattern:  "*.xyz",
			expected: 0,
			contains: []string{},
		},
		{
			name:     "匹配目录中的文件",
			pattern:  "subfile*",
			expected: 1,
			contains: []string{"dir/subfile.txt"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FilterFilesByPattern(files, tt.pattern)

			if len(result) != tt.expected {
				t.Errorf("FilterFilesByPattern() 返回 %d 个文件, want %d", len(result), tt.expected)
			}

			// 检查结果是否包含期望的文件
			for _, expectedFile := range tt.contains {
				found := false
				for _, file := range result {
					if file.Name == expectedFile {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("结果应该包含文件 %q", expectedFile)
				}
			}
		})
	}
}

func TestLimitFiles(t *testing.T) {
	files := []types.FileInfo{
		{Name: "file1.txt", Size: 100},
		{Name: "file2.txt", Size: 200},
		{Name: "file3.txt", Size: 300},
		{Name: "file4.txt", Size: 400},
		{Name: "file5.txt", Size: 500},
	}

	tests := []struct {
		name     string
		limit    int
		expected int
	}{
		{"无限制", 0, 5},
		{"负数限制", -1, 5},
		{"限制为3", 3, 3},
		{"限制等于文件数", 5, 5},
		{"限制大于文件数", 10, 5},
		{"限制为1", 1, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := LimitFiles(files, tt.limit)

			if len(result) != tt.expected {
				t.Errorf("LimitFiles() 返回 %d 个文件, want %d", len(result), tt.expected)
			}

			// 检查返回的文件是否是原始列表的前N个
			for i, file := range result {
				if file.Name != files[i].Name {
					t.Errorf("第 %d 个文件应该是 %q, 但得到 %q", i, files[i].Name, file.Name)
				}
			}
		})
	}
}

func TestLimitFiles_EmptySlice(t *testing.T) {
	var files []types.FileInfo

	result := LimitFiles(files, 5)
	if len(result) != 0 {
		t.Errorf("空切片应该返回空切片，但得到 %d 个文件", len(result))
	}
}

// 基准测试
func BenchmarkFormatFileSize(b *testing.B) {
	sizes := []int64{
		100,                       // 100B
		1024,                      // 1KB
		1024 * 1024,               // 1MB
		1024 * 1024 * 1024,        // 1GB
		1024 * 1024 * 1024 * 1024, // 1TB
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, size := range sizes {
			_ = FormatFileSize(size)
		}
	}
}

func BenchmarkMatchPattern(b *testing.B) {
	patterns := []string{
		"*.txt",
		"test*",
		"*.doc",
		"*file*",
	}

	filenames := []string{
		"test.txt",
		"document.doc",
		"image.jpg",
		"script.sh",
		"datafile.json",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, pattern := range patterns {
			for _, filename := range filenames {
				_ = MatchPattern(filename, pattern)
			}
		}
	}
}

func BenchmarkFilterFilesByPattern(b *testing.B) {
	files := make([]types.FileInfo, 1000)
	for i := 0; i < 1000; i++ {
		files[i] = types.FileInfo{
			Name: fmt.Sprintf("file%d.txt", i),
			Size: int64(i * 100),
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = FilterFilesByPattern(files, "*.txt")
	}
}

func BenchmarkLimitFiles(b *testing.B) {
	files := make([]types.FileInfo, 10000)
	for i := 0; i < 10000; i++ {
		files[i] = types.FileInfo{
			Name: fmt.Sprintf("file%d.txt", i),
			Size: int64(i * 100),
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = LimitFiles(files, 100)
	}
}

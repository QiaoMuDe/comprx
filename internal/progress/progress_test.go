package progress

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gitee.com/MM-Q/comprx/types"
)

// captureOutput 捕获标准输出
func captureOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	_ = w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	return buf.String()
}

// TestNew 测试创建进度显示器
func TestNew(t *testing.T) {
	progress := New()

	if progress == nil {
		t.Fatal("New() 返回了 nil")
	}

	if !progress.Enabled {
		t.Error("期望 Enabled 为 true，实际为 false")
	}

	if progress.BarStyle != types.ProgressStyleText {
		t.Errorf("期望 BarStyle 为 %v，实际为 %v", types.ProgressStyleText, progress.BarStyle)
	}
}

// TestIsEnabled 测试是否启用检查
func TestIsEnabled(t *testing.T) {
	tests := []struct {
		name     string
		enabled  bool
		expected bool
	}{
		{
			name:     "启用状态",
			enabled:  true,
			expected: true,
		},
		{
			name:     "禁用状态",
			enabled:  false,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			progress := &Progress{
				Enabled: tt.enabled,
			}

			result := progress.IsEnabled()
			if result != tt.expected {
				t.Errorf("期望 %v，实际 %v", tt.expected, result)
			}
		})
	}
}

// TestArchive 测试显示压缩文件信息
func TestArchive(t *testing.T) {
	tests := []struct {
		name         string
		enabled      bool
		barStyle     types.ProgressStyle
		archivePath  string
		expectOutput bool
	}{
		{
			name:         "启用且文本样式",
			enabled:      true,
			barStyle:     types.ProgressStyleText,
			archivePath:  "/path/to/test.zip",
			expectOutput: true,
		},
		{
			name:         "禁用状态",
			enabled:      false,
			barStyle:     types.ProgressStyleText,
			archivePath:  "/path/to/test.zip",
			expectOutput: false,
		},
		{
			name:         "启用但非文本样式",
			enabled:      true,
			barStyle:     types.ProgressStyleASCII,
			archivePath:  "/path/to/test.zip",
			expectOutput: false,
		},
		{
			name:         "启用但Unicode样式",
			enabled:      true,
			barStyle:     types.ProgressStyleUnicode,
			archivePath:  "/path/to/test.zip",
			expectOutput: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			progress := &Progress{
				Enabled:  tt.enabled,
				BarStyle: tt.barStyle,
			}

			output := captureOutput(func() {
				progress.Archive(tt.archivePath)
			})

			if tt.expectOutput {
				expectedOutput := fmt.Sprintf("%s %s\n", labelArchive, filepath.Base(tt.archivePath))
				if output != expectedOutput {
					t.Errorf("期望输出 %q，实际输出 %q", expectedOutput, output)
				}
			} else {
				if output != "" {
					t.Errorf("期望无输出，实际输出 %q", output)
				}
			}
		})
	}
}

// TestCompressing 测试显示压缩文件信息
func TestCompressing(t *testing.T) {
	tests := []struct {
		name         string
		enabled      bool
		barStyle     types.ProgressStyle
		filePath     string
		expectOutput bool
	}{
		{
			name:         "启用且文本样式",
			enabled:      true,
			barStyle:     types.ProgressStyleText,
			filePath:     "/path/to/file.txt",
			expectOutput: true,
		},
		{
			name:         "禁用状态",
			enabled:      false,
			barStyle:     types.ProgressStyleText,
			filePath:     "/path/to/file.txt",
			expectOutput: false,
		},
		{
			name:         "启用但非文本样式",
			enabled:      true,
			barStyle:     types.ProgressStyleASCII,
			filePath:     "/path/to/file.txt",
			expectOutput: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			progress := &Progress{
				Enabled:  tt.enabled,
				BarStyle: tt.barStyle,
			}

			output := captureOutput(func() {
				progress.Compressing(tt.filePath)
			})

			if tt.expectOutput {
				expectedOutput := fmt.Sprintf("%s %s\n", labelCompressing, filepath.Base(tt.filePath))
				if output != expectedOutput {
					t.Errorf("期望输出 %q，实际输出 %q", expectedOutput, output)
				}
			} else {
				if output != "" {
					t.Errorf("期望无输出，实际输出 %q", output)
				}
			}
		})
	}
}

// TestInflating 测试显示解压文件
func TestInflating(t *testing.T) {
	tests := []struct {
		name         string
		enabled      bool
		barStyle     types.ProgressStyle
		filePath     string
		expectOutput bool
	}{
		{
			name:         "启用且文本样式",
			enabled:      true,
			barStyle:     types.ProgressStyleText,
			filePath:     "path/to/file.txt",
			expectOutput: true,
		},
		{
			name:         "禁用状态",
			enabled:      false,
			barStyle:     types.ProgressStyleText,
			filePath:     "path/to/file.txt",
			expectOutput: false,
		},
		{
			name:         "启用但非文本样式",
			enabled:      true,
			barStyle:     types.ProgressStyleASCII,
			filePath:     "path/to/file.txt",
			expectOutput: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			progress := &Progress{
				Enabled:  tt.enabled,
				BarStyle: tt.barStyle,
			}

			output := captureOutput(func() {
				progress.Inflating(tt.filePath)
			})

			if tt.expectOutput {
				expectedOutput := fmt.Sprintf("%s %s\n", labelInflating, tt.filePath)
				if output != expectedOutput {
					t.Errorf("期望输出 %q，实际输出 %q", expectedOutput, output)
				}
			} else {
				if output != "" {
					t.Errorf("期望无输出，实际输出 %q", output)
				}
			}
		})
	}
}

// TestCreating 测试显示创建目录
func TestCreating(t *testing.T) {
	tests := []struct {
		name         string
		enabled      bool
		barStyle     types.ProgressStyle
		dirPath      string
		expectOutput bool
	}{
		{
			name:         "启用且文本样式",
			enabled:      true,
			barStyle:     types.ProgressStyleText,
			dirPath:      "path/to/directory/",
			expectOutput: true,
		},
		{
			name:         "禁用状态",
			enabled:      false,
			barStyle:     types.ProgressStyleText,
			dirPath:      "path/to/directory/",
			expectOutput: false,
		},
		{
			name:         "启用但非文本样式",
			enabled:      true,
			barStyle:     types.ProgressStyleASCII,
			dirPath:      "path/to/directory/",
			expectOutput: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			progress := &Progress{
				Enabled:  tt.enabled,
				BarStyle: tt.barStyle,
			}

			output := captureOutput(func() {
				progress.Creating(tt.dirPath)
			})

			if tt.expectOutput {
				expectedOutput := fmt.Sprintf("%s %s\n", labelCreating, tt.dirPath)
				if output != expectedOutput {
					t.Errorf("期望输出 %q，实际输出 %q", expectedOutput, output)
				}
			} else {
				if output != "" {
					t.Errorf("期望无输出，实际输出 %q", output)
				}
			}
		})
	}
}

// TestAdding 测试显示添加文件
func TestAdding(t *testing.T) {
	tests := []struct {
		name         string
		enabled      bool
		barStyle     types.ProgressStyle
		filePath     string
		expectOutput bool
	}{
		{
			name:         "启用且文本样式",
			enabled:      true,
			barStyle:     types.ProgressStyleText,
			filePath:     "path/to/file.txt",
			expectOutput: true,
		},
		{
			name:         "禁用状态",
			enabled:      false,
			barStyle:     types.ProgressStyleText,
			filePath:     "path/to/file.txt",
			expectOutput: false,
		},
		{
			name:         "启用但非文本样式",
			enabled:      true,
			barStyle:     types.ProgressStyleASCII,
			filePath:     "path/to/file.txt",
			expectOutput: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			progress := &Progress{
				Enabled:  tt.enabled,
				BarStyle: tt.barStyle,
			}

			output := captureOutput(func() {
				progress.Adding(tt.filePath)
			})

			if tt.expectOutput {
				expectedOutput := fmt.Sprintf("%s %s\n", labelAdding, tt.filePath)
				if output != expectedOutput {
					t.Errorf("期望输出 %q，实际输出 %q", expectedOutput, output)
				}
			} else {
				if output != "" {
					t.Errorf("期望无输出，实际输出 %q", output)
				}
			}
		})
	}
}

// TestStoring 测试显示存储目录
func TestStoring(t *testing.T) {
	tests := []struct {
		name         string
		enabled      bool
		barStyle     types.ProgressStyle
		dirPath      string
		expectOutput bool
	}{
		{
			name:         "启用且文本样式",
			enabled:      true,
			barStyle:     types.ProgressStyleText,
			dirPath:      "path/to/directory/",
			expectOutput: true,
		},
		{
			name:         "禁用状态",
			enabled:      false,
			barStyle:     types.ProgressStyleText,
			dirPath:      "path/to/directory/",
			expectOutput: false,
		},
		{
			name:         "启用但非文本样式",
			enabled:      true,
			barStyle:     types.ProgressStyleASCII,
			dirPath:      "path/to/directory/",
			expectOutput: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			progress := &Progress{
				Enabled:  tt.enabled,
				BarStyle: tt.barStyle,
			}

			output := captureOutput(func() {
				progress.Storing(tt.dirPath)
			})

			if tt.expectOutput {
				expectedOutput := fmt.Sprintf("%s %s\n", labelStoring, tt.dirPath)
				if output != expectedOutput {
					t.Errorf("期望输出 %q，实际输出 %q", expectedOutput, output)
				}
			} else {
				if output != "" {
					t.Errorf("期望无输出，实际输出 %q", output)
				}
			}
		})
	}
}

// TestProgressLabels 测试进度标签常量
func TestProgressLabels(t *testing.T) {
	tests := []struct {
		name     string
		label    string
		expected string
	}{
		{
			name:     "Archive标签",
			label:    labelArchive,
			expected: "Archive:    ",
		},
		{
			name:     "Inflating标签",
			label:    labelInflating,
			expected: "  inflating:",
		},
		{
			name:     "Creating标签",
			label:    labelCreating,
			expected: "   creating:",
		},
		{
			name:     "Adding标签",
			label:    labelAdding,
			expected: "     adding:",
		},
		{
			name:     "Storing标签",
			label:    labelStoring,
			expected: "    storing:",
		},
		{
			name:     "Compressing标签",
			label:    labelCompressing,
			expected: "compressing:",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.label != tt.expected {
				t.Errorf("期望标签 %q，实际 %q", tt.expected, tt.label)
			}
		})
	}
}

// TestProgressWorkflow 测试完整的进度显示工作流
func TestProgressWorkflow(t *testing.T) {
	progress := New()

	// 测试压缩工作流
	output := captureOutput(func() {
		progress.Archive("test.zip")
		progress.Adding("file1.txt")
		progress.Adding("file2.txt")
		progress.Storing("dir1/")
		progress.Compressing("large_file.dat")
	})

	expectedLines := []string{
		"Archive:     test.zip",
		"     adding: file1.txt",
		"     adding: file2.txt",
		"    storing: dir1/",
		"compressing: large_file.dat",
	}

	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) != len(expectedLines) {
		t.Errorf("期望 %d 行输出，实际 %d 行", len(expectedLines), len(lines))
	}

	for i, expectedLine := range expectedLines {
		if i < len(lines) && lines[i] != expectedLine {
			t.Errorf("第 %d 行期望 %q，实际 %q", i+1, expectedLine, lines[i])
		}
	}
}

// TestProgressWithDifferentPaths 测试不同路径格式
func TestProgressWithDifferentPaths(t *testing.T) {
	progress := New()

	tests := []struct {
		name     string
		path     string
		method   func(string)
		expected string
	}{
		{
			name:     "绝对路径文件",
			path:     "/home/user/documents/file.txt",
			method:   progress.Adding,
			expected: "     adding: /home/user/documents/file.txt\n",
		},
		{
			name:     "相对路径文件",
			path:     "docs/readme.md",
			method:   progress.Adding,
			expected: "     adding: docs/readme.md\n",
		},
		{
			name:     "Windows路径",
			path:     "C:\\Users\\test\\file.txt",
			method:   progress.Adding,
			expected: "     adding: C:\\Users\\test\\file.txt\n",
		},
		{
			name:     "带空格的路径",
			path:     "my documents/test file.txt",
			method:   progress.Adding,
			expected: "     adding: my documents/test file.txt\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := captureOutput(func() {
				tt.method(tt.path)
			})

			if output != tt.expected {
				t.Errorf("期望输出 %q，实际输出 %q", tt.expected, output)
			}
		})
	}
}

// BenchmarkProgressOperations 基准测试进度操作
func BenchmarkProgressOperations(b *testing.B) {
	progress := New()
	testPath := "test/path/file.txt"

	b.Run("Archive", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			progress.Archive(testPath)
		}
	})

	b.Run("Adding", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			progress.Adding(testPath)
		}
	})

	b.Run("Inflating", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			progress.Inflating(testPath)
		}
	})

	b.Run("Creating", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			progress.Creating(testPath)
		}
	})

	b.Run("Storing", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			progress.Storing(testPath)
		}
	})

	b.Run("Compressing", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			progress.Compressing(testPath)
		}
	})
}

// BenchmarkProgressDisabled 基准测试禁用状态下的性能
func BenchmarkProgressDisabled(b *testing.B) {
	progress := &Progress{
		Enabled:  false,
		BarStyle: types.ProgressStyleText,
	}
	testPath := "test/path/file.txt"

	b.Run("DisabledAdding", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			progress.Adding(testPath)
		}
	})
}

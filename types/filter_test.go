package types

import (
	"fmt"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestFilterOptions_ShouldSkipByParams_NilFilter(t *testing.T) {
	var filter *FilterOptions = nil

	result := filter.ShouldSkipByParams("test.go", 1000, false)
	if result != false {
		t.Errorf("期望 nil 过滤器返回 false，实际返回 %v", result)
	}
}

func TestFilterOptions_ShouldSkipByParams_EmptyFilter(t *testing.T) {
	filter := &FilterOptions{}

	result := filter.ShouldSkipByParams("test.go", 1000, false)
	if result != false {
		t.Errorf("期望空过滤器返回 false，实际返回 %v", result)
	}
}

func TestFilterOptions_ShouldSkipByParams_SizeFilter(t *testing.T) {
	tests := []struct {
		name     string
		filter   FilterOptions
		fileName string
		fileSize int64
		isDir    bool
		expected bool
		desc     string
	}{
		{
			name:     "文件大小在范围内",
			filter:   FilterOptions{MinSize: 100, MaxSize: 2000},
			fileName: "test.go",
			fileSize: 1000,
			isDir:    false,
			expected: false,
			desc:     "文件大小符合要求，不应该跳过",
		},
		{
			name:     "文件太小",
			filter:   FilterOptions{MinSize: 100},
			fileName: "small.go",
			fileSize: 50,
			isDir:    false,
			expected: true,
			desc:     "文件太小，应该跳过",
		},
		{
			name:     "文件太大",
			filter:   FilterOptions{MaxSize: 1000},
			fileName: "large.go",
			fileSize: 2000,
			isDir:    false,
			expected: true,
			desc:     "文件太大，应该跳过",
		},
		{
			name:     "目录忽略大小限制",
			filter:   FilterOptions{MinSize: 100, MaxSize: 1000},
			fileName: "testdir",
			fileSize: 0,
			isDir:    true,
			expected: false,
			desc:     "目录应该忽略大小限制",
		},
		{
			name:     "边界值-最小大小",
			filter:   FilterOptions{MinSize: 100},
			fileName: "boundary.go",
			fileSize: 100,
			isDir:    false,
			expected: false,
			desc:     "文件大小等于最小值，不应该跳过",
		},
		{
			name:     "边界值-最大大小",
			filter:   FilterOptions{MaxSize: 1000},
			fileName: "boundary.go",
			fileSize: 1000,
			isDir:    false,
			expected: false,
			desc:     "文件大小等于最大值，不应该跳过",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.filter.ShouldSkipByParams(tt.fileName, tt.fileSize, tt.isDir)
			if result != tt.expected {
				t.Errorf("%s: 期望 %v，实际 %v", tt.desc, tt.expected, result)
			}
		})
	}
}

func TestFilterOptions_ShouldSkipByParams_IncludePattern(t *testing.T) {
	tests := []struct {
		name     string
		filter   FilterOptions
		filePath string
		expected bool
		desc     string
	}{
		{
			name:     "匹配包含模式-文件扩展名",
			filter:   FilterOptions{Include: []string{"*.go"}},
			filePath: "main.go",
			expected: false,
			desc:     "Go文件匹配包含模式，不应该跳过",
		},
		{
			name:     "不匹配包含模式-文件扩展名",
			filter:   FilterOptions{Include: []string{"*.go"}},
			filePath: "config.json",
			expected: true,
			desc:     "非Go文件不匹配包含模式，应该跳过",
		},
		{
			name:     "匹配包含模式-路径模式",
			filter:   FilterOptions{Include: []string{"src/*.go"}},
			filePath: "src/main.go",
			expected: false,
			desc:     "src目录下的Go文件匹配包含模式，不应该跳过",
		},
		{
			name:     "不匹配包含模式-路径模式",
			filter:   FilterOptions{Include: []string{"src/*.go"}},
			filePath: "test/main.go",
			expected: true,
			desc:     "test目录下的Go文件不匹配包含模式，应该跳过",
		},
		{
			name:     "多个包含模式-匹配其中一个",
			filter:   FilterOptions{Include: []string{"*.go", "*.md"}},
			filePath: "README.md",
			expected: false,
			desc:     "Markdown文件匹配包含模式之一，不应该跳过",
		},
		{
			name:     "多个包含模式-都不匹配",
			filter:   FilterOptions{Include: []string{"*.go", "*.md"}},
			filePath: "config.json",
			expected: true,
			desc:     "JSON文件不匹配任何包含模式，应该跳过",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.filter.ShouldSkipByParams(tt.filePath, 1000, false)
			if result != tt.expected {
				t.Errorf("%s: 期望 %v，实际 %v", tt.desc, tt.expected, result)
			}
		})
	}
}

func TestFilterOptions_ShouldSkipByParams_ExcludePattern(t *testing.T) {
	tests := []struct {
		name     string
		filter   FilterOptions
		filePath string
		expected bool
		desc     string
	}{
		{
			name:     "匹配排除模式",
			filter:   FilterOptions{Exclude: []string{"*_test.go"}},
			filePath: "main_test.go",
			expected: true,
			desc:     "测试文件匹配排除模式，应该跳过",
		},
		{
			name:     "不匹配排除模式",
			filter:   FilterOptions{Exclude: []string{"*_test.go"}},
			filePath: "main.go",
			expected: false,
			desc:     "普通Go文件不匹配排除模式，不应该跳过",
		},
		{
			name:     "匹配排除模式-目录",
			filter:   FilterOptions{Exclude: []string{"vendor/*"}},
			filePath: "vendor/package.go",
			expected: true,
			desc:     "vendor目录下的文件匹配排除模式，应该跳过",
		},
		{
			name:     "多个排除模式-匹配其中一个",
			filter:   FilterOptions{Exclude: []string{"*.tmp", "*.log"}},
			filePath: "debug.log",
			expected: true,
			desc:     "日志文件匹配排除模式之一，应该跳过",
		},
		{
			name:     "多个排除模式-都不匹配",
			filter:   FilterOptions{Exclude: []string{"*.tmp", "*.log"}},
			filePath: "main.go",
			expected: false,
			desc:     "Go文件不匹配任何排除模式，不应该跳过",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.filter.ShouldSkipByParams(tt.filePath, 1000, false)
			if result != tt.expected {
				t.Errorf("%s: 期望 %v，实际 %v", tt.desc, tt.expected, result)
			}
		})
	}
}

func TestFilterOptions_ShouldSkipByParams_CombinedFilters(t *testing.T) {
	tests := []struct {
		name     string
		filter   FilterOptions
		filePath string
		fileSize int64
		expected bool
		desc     string
	}{
		{
			name: "包含+排除-通过包含但被排除",
			filter: FilterOptions{
				Include: []string{"*.go"},
				Exclude: []string{"*_test.go"},
			},
			filePath: "main_test.go",
			fileSize: 1000,
			expected: true,
			desc:     "测试文件虽然匹配包含模式，但匹配排除模式，应该跳过",
		},
		{
			name: "包含+排除-通过包含且不被排除",
			filter: FilterOptions{
				Include: []string{"*.go"},
				Exclude: []string{"*_test.go"},
			},
			filePath: "main.go",
			fileSize: 1000,
			expected: false,
			desc:     "普通Go文件匹配包含模式且不匹配排除模式，不应该跳过",
		},
		{
			name: "包含+大小-匹配包含但文件太大",
			filter: FilterOptions{
				Include: []string{"*.go"},
				MaxSize: 500,
			},
			filePath: "main.go",
			fileSize: 1000,
			expected: true,
			desc:     "Go文件匹配包含模式但超过大小限制，应该跳过",
		},
		{
			name: "排除+大小-不匹配排除但文件太小",
			filter: FilterOptions{
				Exclude: []string{"*_test.go"},
				MinSize: 500,
			},
			filePath: "main.go",
			fileSize: 100,
			expected: true,
			desc:     "普通Go文件不匹配排除模式但小于最小大小，应该跳过",
		},
		{
			name: "全部条件-都通过",
			filter: FilterOptions{
				Include: []string{"*.go"},
				Exclude: []string{"*_test.go"},
				MinSize: 100,
				MaxSize: 2000,
			},
			filePath: "main.go",
			fileSize: 1000,
			expected: false,
			desc:     "普通Go文件通过所有过滤条件，不应该跳过",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.filter.ShouldSkipByParams(tt.filePath, tt.fileSize, false)
			if result != tt.expected {
				t.Errorf("%s: 期望 %v，实际 %v", tt.desc, tt.expected, result)
			}
		})
	}
}

func TestFilterOptions_matchPattern(t *testing.T) {
	filter := &FilterOptions{}

	tests := []struct {
		pattern  string
		path     string
		expected bool
		desc     string
	}{
		// 文件名匹配
		{"*.go", "main.go", true, "文件扩展名匹配"},
		{"*.go", "config.json", false, "文件扩展名不匹配"},
		{"main.*", "main.go", true, "文件名前缀匹配"},
		{"main.*", "test.go", false, "文件名前缀不匹配"},

		// 完整路径匹配
		{"src/*.go", "src/main.go", true, "路径模式匹配"},
		{"src/*.go", "test/main.go", false, "路径模式不匹配"},
		{"*/main.go", "src/main.go", true, "通配符路径匹配"},
		{"*/main.go", "src/test.go", false, "通配符路径不匹配"},

		// 目录匹配（以/结尾）
		{"vendor/", "vendor", true, "目录模式匹配"},
		{"vendor/", "src", false, "目录模式不匹配"},
		{"test\\", "test", true, "Windows风格目录模式匹配"},

		// 边界情况
		{"", "test.go", false, "空模式不匹配任何文件"},
		{"*", "test.go", true, "通配符匹配所有文件"},
		{"test.go", "test.go", true, "完全匹配"},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			// 预计算路径信息以匹配新的函数签名
			baseName := filepath.Base(tt.path)
			slashPath := filepath.ToSlash(tt.path)
			// 标准化模式以匹配新的函数签名
			normalizedPattern := strings.ReplaceAll(tt.pattern, "\\", "/")
			result := filter.matchPattern(normalizedPattern, tt.path, baseName, slashPath)
			if result != tt.expected {
				t.Errorf("模式 '%s' 匹配路径 '%s': 期望 %v，实际 %v",
					tt.pattern, tt.path, tt.expected, result)
			}
		})
	}
}

func TestFilterOptions_Validate(t *testing.T) {
	tests := []struct {
		name     string
		filter   FilterOptions
		hasError bool
		desc     string
	}{
		{
			name:     "有效的过滤器",
			filter:   FilterOptions{Include: []string{"*.go"}, Exclude: []string{"*_test.go"}, MinSize: 100, MaxSize: 2000},
			hasError: false,
			desc:     "所有参数都有效",
		},
		{
			name:     "最小大小为负数",
			filter:   FilterOptions{MinSize: -100},
			hasError: true,
			desc:     "最小大小不能为负数",
		},
		{
			name:     "最大大小为负数",
			filter:   FilterOptions{MaxSize: -1000},
			hasError: true,
			desc:     "最大大小不能为负数",
		},
		{
			name:     "最小大小大于最大大小",
			filter:   FilterOptions{MinSize: 2000, MaxSize: 1000},
			hasError: true,
			desc:     "最小大小不能大于最大大小",
		},
		{
			name:     "包含模式为空字符串",
			filter:   FilterOptions{Include: []string{"*.go", ""}},
			hasError: true,
			desc:     "包含模式不能为空字符串",
		},
		{
			name:     "排除模式为空字符串",
			filter:   FilterOptions{Exclude: []string{"*_test.go", ""}},
			hasError: true,
			desc:     "排除模式不能为空字符串",
		},
		{
			name:     "边界值-最小等于最大",
			filter:   FilterOptions{MinSize: 1000, MaxSize: 1000},
			hasError: false,
			desc:     "最小大小等于最大大小是有效的",
		},
		{
			name:     "零值",
			filter:   FilterOptions{MinSize: 0, MaxSize: 0},
			hasError: false,
			desc:     "零值是有效的",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.filter.Validate()
			hasError := err != nil
			if hasError != tt.hasError {
				t.Errorf("%s: 期望错误状态 %v，实际 %v，错误: %v",
					tt.desc, tt.hasError, hasError, err)
			}
		})
	}
}

func TestHasFilterConditions(t *testing.T) {
	tests := []struct {
		name     string
		filter   *FilterOptions
		expected bool
		desc     string
	}{
		{
			name:     "nil过滤器",
			filter:   nil,
			expected: false,
			desc:     "nil过滤器没有过滤条件",
		},
		{
			name:     "空过滤器",
			filter:   &FilterOptions{},
			expected: false,
			desc:     "空过滤器没有过滤条件",
		},
		{
			name:     "有包含模式",
			filter:   &FilterOptions{Include: []string{"*.go"}},
			expected: true,
			desc:     "有包含模式的过滤器有过滤条件",
		},
		{
			name:     "有排除模式",
			filter:   &FilterOptions{Exclude: []string{"*_test.go"}},
			expected: true,
			desc:     "有排除模式的过滤器有过滤条件",
		},
		{
			name:     "有最小大小",
			filter:   &FilterOptions{MinSize: 100},
			expected: true,
			desc:     "有最小大小限制的过滤器有过滤条件",
		},
		{
			name:     "有最大大小",
			filter:   &FilterOptions{MaxSize: 1000},
			expected: true,
			desc:     "有最大大小限制的过滤器有过滤条件",
		},
		{
			name:     "有多个条件",
			filter:   &FilterOptions{Include: []string{"*.go"}, MinSize: 100},
			expected: true,
			desc:     "有多个过滤条件的过滤器有过滤条件",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := HasFilterConditions(tt.filter)
			if result != tt.expected {
				t.Errorf("%s: 期望 %v，实际 %v", tt.desc, tt.expected, result)
			}
		})
	}
}

// ==============================================
// 路径分隔符兼容性测试
// ==============================================

func TestFilterOptions_PathSeparatorCompatibility(t *testing.T) {
	tests := []struct {
		name     string
		filter   FilterOptions
		filePath string
		expected bool
		desc     string
	}{
		{
			name:     "Windows路径分隔符-包含模式",
			filter:   FilterOptions{Include: []string{"src\\*.go"}},
			filePath: "src/main.go",
			expected: false,
			desc:     "Windows风格模式应该匹配Unix风格路径",
		},
		{
			name:     "Windows路径分隔符-排除模式",
			filter:   FilterOptions{Exclude: []string{"vendor\\node_modules\\*"}},
			filePath: "vendor/node_modules/package.json",
			expected: true,
			desc:     "Windows风格排除模式应该匹配Unix风格路径",
		},
		{
			name:     "混合路径分隔符",
			filter:   FilterOptions{Exclude: []string{"vendor\\*/node_modules/*"}},
			filePath: "vendor/lib/node_modules/package.json",
			expected: true,
			desc:     "混合路径分隔符应该正确处理",
		},
		{
			name:     "Unix路径分隔符",
			filter:   FilterOptions{Include: []string{"src/*.go"}},
			filePath: "src\\main.go",
			expected: false,
			desc:     "Unix风格模式应该匹配Windows风格路径",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.filter.ShouldSkipByParams(tt.filePath, 1000, false)
			if result != tt.expected {
				t.Errorf("%s: 期望 %v，实际 %v", tt.desc, tt.expected, result)
			}
		})
	}
}

// ==============================================
// Unicode和特殊字符测试
// ==============================================

func TestFilterOptions_UnicodeAndSpecialChars(t *testing.T) {
	tests := []struct {
		name     string
		filter   FilterOptions
		filePath string
		expected bool
		desc     string
	}{
		{
			name:     "中文文件名",
			filter:   FilterOptions{Include: []string{"*.txt"}},
			filePath: "测试文件.txt",
			expected: false,
			desc:     "应该支持Unicode文件名",
		},
		{
			name:     "日文文件名",
			filter:   FilterOptions{Include: []string{"*テスト*"}},
			filePath: "テストファイル.go",
			expected: false,
			desc:     "应该支持日文文件名匹配",
		},
		{
			name:     "空格文件名",
			filter:   FilterOptions{Include: []string{"*test*"}},
			filePath: "my test file.go",
			expected: false,
			desc:     "应该支持包含空格的文件名",
		},
		{
			name:     "特殊字符文件名-@符号",
			filter:   FilterOptions{Exclude: []string{"*@*"}},
			filePath: "file@version.txt",
			expected: true,
			desc:     "应该支持@符号匹配",
		},
		{
			name:     "特殊字符文件名-连字符",
			filter:   FilterOptions{Include: []string{"*-config.*"}},
			filePath: "app-config.json",
			expected: false,
			desc:     "应该支持连字符匹配",
		},
		{
			name:     "特殊字符文件名-下划线",
			filter:   FilterOptions{Exclude: []string{"*_backup_*"}},
			filePath: "data_backup_2023.sql",
			expected: true,
			desc:     "应该支持下划线匹配",
		},
		{
			name:     "括号文件名",
			filter:   FilterOptions{Include: []string{"*(*).*"}},
			filePath: "file(1).txt",
			expected: false,
			desc:     "应该支持括号文件名",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.filter.ShouldSkipByParams(tt.filePath, 1000, false)
			if result != tt.expected {
				t.Errorf("%s: 期望 %v，实际 %v", tt.desc, tt.expected, result)
			}
		})
	}
}

// ==============================================
// 深层嵌套路径测试
// ==============================================

func TestFilterOptions_DeepNestedPaths(t *testing.T) {
	tests := []struct {
		name     string
		filter   FilterOptions
		filePath string
		expected bool
		desc     string
	}{
		{
			name:     "深层嵌套匹配-test目录",
			filter:   FilterOptions{Include: []string{"*/test/*"}},
			filePath: "a/b/c/test/d/e/file.go",
			expected: false,
			desc:     "应该匹配深层嵌套的test目录",
		},
		{
			name:     "深层嵌套排除-node_modules",
			filter:   FilterOptions{Exclude: []string{"*node_modules*"}},
			filePath: "project/src/components/node_modules/package/index.js",
			expected: true,
			desc:     "应该排除任意深度的node_modules",
		},
		{
			name:     "超深层路径",
			filter:   FilterOptions{Include: []string{"*.go"}},
			filePath: "a/b/c/d/e/f/g/h/i/j/k/l/m/n/o/p/q/r/s/t/main.go",
			expected: false,
			desc:     "应该处理超深层路径",
		},
		{
			name:     "深层目录匹配",
			filter:   FilterOptions{Exclude: []string{"*/cache/*"}},
			filePath: "app/data/temp/cache/session/file.tmp",
			expected: true,
			desc:     "应该匹配深层目录中的cache",
		},
		{
			name:     "多级通配符匹配",
			filter:   FilterOptions{Include: []string{"src/*/components/*.vue"}},
			filePath: "src/admin/components/UserList.vue",
			expected: false,
			desc:     "应该匹配多级通配符路径",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.filter.ShouldSkipByParams(tt.filePath, 1000, false)
			if result != tt.expected {
				t.Errorf("%s: 期望 %v，实际 %v", tt.desc, tt.expected, result)
			}
		})
	}
}

// ==============================================
// 性能边界测试
// ==============================================

func TestFilterOptions_PerformanceBoundary(t *testing.T) {
	// 大量模式测试
	t.Run("大量包含模式", func(t *testing.T) {
		manyPatterns := make([]string, 1000)
		for i := 0; i < 1000; i++ {
			manyPatterns[i] = fmt.Sprintf("*pattern%d*", i)
		}
		
		filter := FilterOptions{Include: manyPatterns}
		
		// 测试大量模式不会导致性能问题
		start := time.Now()
		result := filter.ShouldSkipByParams("test.go", 1000, false)
		duration := time.Since(start)
		
		if duration > time.Millisecond*100 {
			t.Errorf("大量包含模式匹配耗时过长: %v", duration)
		}
		
		if !result {
			t.Error("应该跳过不匹配的文件")
		}
	})

	t.Run("大量排除模式", func(t *testing.T) {
		manyPatterns := make([]string, 500)
		for i := 0; i < 500; i++ {
			manyPatterns[i] = fmt.Sprintf("*exclude%d*", i)
		}
		
		filter := FilterOptions{Exclude: manyPatterns}
		
		start := time.Now()
		result := filter.ShouldSkipByParams("normal.go", 1000, false)
		duration := time.Since(start)
		
		if duration > time.Millisecond*50 {
			t.Errorf("大量排除模式匹配耗时过长: %v", duration)
		}
		
		if result {
			t.Error("不应该跳过不匹配排除模式的文件")
		}
	})

	t.Run("超长路径", func(t *testing.T) {
		// 创建超长路径
		longPath := strings.Repeat("verylongdirectoryname/", 100) + "file.go"
		
		filter := FilterOptions{Include: []string{"*.go"}}
		
		start := time.Now()
		result := filter.ShouldSkipByParams(longPath, 1000, false)
		duration := time.Since(start)
		
		if duration > time.Millisecond*10 {
			t.Errorf("超长路径匹配耗时过长: %v", duration)
		}
		
		if result {
			t.Error("不应该跳过匹配的.go文件")
		}
	})
}

// ==============================================
// 边界值和异常输入测试
// ==============================================

func TestFilterOptions_BoundaryAndEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		filter   FilterOptions
		filePath string
		fileSize int64
		isDir    bool
		expected bool
		desc     string
	}{
		{
			name:     "空路径",
			filter:   FilterOptions{Include: []string{"*.go"}},
			filePath: "",
			fileSize: 1000,
			isDir:    false,
			expected: true,
			desc:     "空路径应该被跳过",
		},
		{
			name:     "极大文件大小",
			filter:   FilterOptions{MaxSize: 1024 * 1024 * 1024}, // 1GB
			filePath: "huge.file",
			fileSize: 1024 * 1024 * 1024 * 10, // 10GB
			isDir:    false,
			expected: true,
			desc:     "超大文件应该被跳过",
		},
		{
			name:     "零大小文件",
			filter:   FilterOptions{MinSize: 1},
			filePath: "empty.txt",
			fileSize: 0,
			isDir:    false,
			expected: true,
			desc:     "零大小文件应该被跳过",
		},
		{
			name:     "负数文件大小",
			filter:   FilterOptions{MinSize: 100},
			filePath: "invalid.txt",
			fileSize: -1,
			isDir:    false,
			expected: true,
			desc:     "负数大小文件应该被跳过",
		},
		{
			name:     "单字符路径",
			filter:   FilterOptions{Include: []string{"*"}},
			filePath: "a",
			fileSize: 100,
			isDir:    false,
			expected: false,
			desc:     "单字符路径应该被匹配",
		},
		{
			name:     "只有扩展名的文件",
			filter:   FilterOptions{Include: []string{"*.go"}},
			filePath: ".go",
			fileSize: 100,
			isDir:    false,
			expected: false,
			desc:     "只有扩展名的文件应该被匹配",
		},
		{
			name:     "没有扩展名的文件",
			filter:   FilterOptions{Exclude: []string{"*.*"}},
			filePath: "README",
			fileSize: 100,
			isDir:    false,
			expected: false,
			desc:     "没有扩展名的文件不应该被排除",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.filter.ShouldSkipByParams(tt.filePath, tt.fileSize, tt.isDir)
			if result != tt.expected {
				t.Errorf("%s: 期望 %v，实际 %v", tt.desc, tt.expected, result)
			}
		})
	}
}

// ==============================================
// 复杂glob模式测试
// ==============================================

func TestFilterOptions_ComplexGlobPatterns(t *testing.T) {
	tests := []struct {
		name     string
		filter   FilterOptions
		filePath string
		expected bool
		desc     string
	}{
		{
			name:     "字符类匹配-c和h文件",
			filter:   FilterOptions{Include: []string{"*.[ch]"}},
			filePath: "main.c",
			expected: false,
			desc:     "应该匹配.c文件",
		},
		{
			name:     "字符类匹配-h文件",
			filter:   FilterOptions{Include: []string{"*.[ch]"}},
			filePath: "header.h",
			expected: false,
			desc:     "应该匹配.h文件",
		},
		{
			name:     "字符类不匹配",
			filter:   FilterOptions{Include: []string{"*.[ch]"}},
			filePath: "main.go",
			expected: true,
			desc:     "不应该匹配.go文件",
		},
		{
			name:     "字符范围匹配",
			filter:   FilterOptions{Include: []string{"test[0-9].txt"}},
			filePath: "test5.txt",
			expected: false,
			desc:     "应该匹配数字范围",
		},
		{
			name:     "字符范围不匹配",
			filter:   FilterOptions{Include: []string{"test[0-9].txt"}},
			filePath: "testa.txt",
			expected: true,
			desc:     "不应该匹配字母",
		},
		{
			name:     "问号通配符",
			filter:   FilterOptions{Include: []string{"test?.go"}},
			filePath: "test1.go",
			expected: false,
			desc:     "问号应该匹配单个字符",
		},
		{
			name:     "问号通配符不匹配多字符",
			filter:   FilterOptions{Include: []string{"test?.go"}},
			filePath: "test12.go",
			expected: true,
			desc:     "问号不应该匹配多个字符",
		},
		{
			name:     "复合模式",
			filter:   FilterOptions{Include: []string{"src/**/test*.go"}},
			filePath: "src/components/test_utils.go",
			expected: false,
			desc:     "应该匹配复合模式",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.filter.ShouldSkipByParams(tt.filePath, 1000, false)
			if result != tt.expected {
				t.Errorf("%s: 期望 %v，实际 %v", tt.desc, tt.expected, result)
			}
		})
	}
}

// ==============================================
// 并发安全测试
// ==============================================

func TestFilterOptions_ConcurrentSafety(t *testing.T) {
	filter := &FilterOptions{
		Include: []string{"*.go", "*.md", "*.json"},
		Exclude: []string{"*_test.go", "vendor/*", "node_modules/*"},
		MinSize: 100,
		MaxSize: 1024 * 1024,
	}

	// 并发测试
	var wg sync.WaitGroup
	const numGoroutines = 100
	const numOperations = 1000

	// 用于收集错误
	errors := make(chan error, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			
			for j := 0; j < numOperations; j++ {
				// 测试不同类型的文件路径
				testPaths := []string{
					fmt.Sprintf("test%d_%d.go", id, j),
					fmt.Sprintf("src/main%d.go", j),
					fmt.Sprintf("vendor/lib%d.js", j),
					fmt.Sprintf("README%d.md", id),
					fmt.Sprintf("config%d.json", j),
				}
				
				for _, path := range testPaths {
					result := filter.ShouldSkipByParams(path, int64(1000+j), false)
					_ = result // 使用结果避免编译器优化
				}
			}
		}(i)
	}

	// 等待所有goroutine完成
	go func() {
		wg.Wait()
		close(errors)
	}()

	// 检查是否有错误
	for err := range errors {
		if err != nil {
			t.Errorf("并发测试出现错误: %v", err)
		}
	}
}

// ==============================================
// 基准测试
// ==============================================

func BenchmarkFilterOptions_ShouldSkipByParams(b *testing.B) {
	filter := &FilterOptions{
		Include: []string{"*.go", "*.md"},
		Exclude: []string{"*_test.go", "vendor/*"},
		MinSize: 100,
		MaxSize: 1024 * 1024,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		filter.ShouldSkipByParams("src/main.go", 1000, false)
	}
}

func BenchmarkFilterOptions_matchPattern(b *testing.B) {
	filter := &FilterOptions{}

	// 预计算路径信息以匹配新的函数签名
	path := "main.go"
	baseName := filepath.Base(path)
	slashPath := filepath.ToSlash(path)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		filter.matchPattern("*.go", path, baseName, slashPath)
	}
}

func BenchmarkFilterOptions_FastMatch(b *testing.B) {
	filter := &FilterOptions{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		filter.fastMatch("*.go", "main.go", "main.go")
	}
}

func BenchmarkFilterOptions_ComplexMatch(b *testing.B) {
	filter := &FilterOptions{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		filter.complexMatch("src/[a-z]*.go", "src/main.go", "main.go", "src/main.go")
	}
}

func BenchmarkFilterOptions_ConcurrentAccess(b *testing.B) {
	filter := &FilterOptions{
		Include: []string{"*.go", "*.md"},
		Exclude: []string{"*_test.go"},
		MinSize: 100,
		MaxSize: 1024 * 1024,
	}

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			filter.ShouldSkipByParams("src/main.go", 1000, false)
		}
	})
}

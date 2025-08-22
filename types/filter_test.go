package types

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

// mockFileInfo 模拟文件信息，用于测试
type mockFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
	isDir   bool
}

func (m mockFileInfo) Name() string       { return m.name }
func (m mockFileInfo) Size() int64        { return m.size }
func (m mockFileInfo) Mode() os.FileMode  { return m.mode }
func (m mockFileInfo) ModTime() time.Time { return m.modTime }
func (m mockFileInfo) IsDir() bool        { return m.isDir }
func (m mockFileInfo) Sys() interface{}   { return nil }

// 创建文件信息的辅助函数
func createFileInfo(name string, size int64, isDir bool) os.FileInfo {
	return mockFileInfo{
		name:    name,
		size:    size,
		mode:    0644,
		modTime: time.Now(),
		isDir:   isDir,
	}
}

func TestFilterOptions_ShouldSkip_NilFilter(t *testing.T) {
	var filter *FilterOptions = nil
	fileInfo := createFileInfo("test.go", 1000, false)

	result := filter.ShouldSkip("test.go", fileInfo)
	if result != false {
		t.Errorf("期望 nil 过滤器返回 false，实际返回 %v", result)
	}
}

func TestFilterOptions_ShouldSkip_EmptyFilter(t *testing.T) {
	filter := &FilterOptions{}
	fileInfo := createFileInfo("test.go", 1000, false)

	result := filter.ShouldSkip("test.go", fileInfo)
	if result != false {
		t.Errorf("期望空过滤器返回 false，实际返回 %v", result)
	}
}

func TestFilterOptions_ShouldSkip_SizeFilter(t *testing.T) {
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
			fileInfo := createFileInfo(tt.fileName, tt.fileSize, tt.isDir)
			result := tt.filter.ShouldSkip(tt.fileName, fileInfo)
			if result != tt.expected {
				t.Errorf("%s: 期望 %v，实际 %v", tt.desc, tt.expected, result)
			}
		})
	}
}

func TestFilterOptions_ShouldSkip_IncludePattern(t *testing.T) {
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
			fileInfo := createFileInfo(filepath.Base(tt.filePath), 1000, false)
			result := tt.filter.ShouldSkip(tt.filePath, fileInfo)
			if result != tt.expected {
				t.Errorf("%s: 期望 %v，实际 %v", tt.desc, tt.expected, result)
			}
		})
	}
}

func TestFilterOptions_ShouldSkip_ExcludePattern(t *testing.T) {
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
			fileInfo := createFileInfo(filepath.Base(tt.filePath), 1000, false)
			result := tt.filter.ShouldSkip(tt.filePath, fileInfo)
			if result != tt.expected {
				t.Errorf("%s: 期望 %v，实际 %v", tt.desc, tt.expected, result)
			}
		})
	}
}

func TestFilterOptions_ShouldSkip_CombinedFilters(t *testing.T) {
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
			fileInfo := createFileInfo(filepath.Base(tt.filePath), tt.fileSize, false)
			result := tt.filter.ShouldSkip(tt.filePath, fileInfo)
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
			result := filter.matchPattern(tt.pattern, tt.path)
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

func TestLoadExcludeFromFile(t *testing.T) {
	// 创建临时测试文件
	tempDir := t.TempDir()

	// 测试正常的忽略文件
	normalIgnoreFile := filepath.Join(tempDir, ".testignore")
	normalContent := `# 这是注释
*.tmp
*.log

# 另一个注释
vendor/
node_modules/
*.exe`

	err := os.WriteFile(normalIgnoreFile, []byte(normalContent), 0644)
	if err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}

	// 测试加载正常文件
	patterns, err := LoadExcludeFromFile(normalIgnoreFile)
	if err != nil {
		t.Fatalf("加载忽略文件失败: %v", err)
	}

	expectedPatterns := []string{"*.tmp", "*.log", "vendor/", "node_modules/", "*.exe"}
	if len(patterns) != len(expectedPatterns) {
		t.Errorf("期望 %d 个模式，实际 %d 个", len(expectedPatterns), len(patterns))
	}

	for i, expected := range expectedPatterns {
		if i >= len(patterns) || patterns[i] != expected {
			t.Errorf("模式 %d: 期望 '%s'，实际 '%s'", i, expected, patterns[i])
		}
	}

	// 测试不存在的文件
	nonExistentFile := filepath.Join(tempDir, ".nonexistent")
	_, err = LoadExcludeFromFile(nonExistentFile)
	if err == nil {
		t.Error("期望加载不存在的文件时返回错误")
	}
}

func TestLoadExcludeFromFileOrEmpty(t *testing.T) {
	tempDir := t.TempDir()

	// 测试存在的文件
	existingFile := filepath.Join(tempDir, ".testignore")
	content := "*.tmp\n*.log"
	err := os.WriteFile(existingFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}

	patterns := LoadExcludeFromFileOrEmpty(existingFile)
	if len(patterns) != 2 {
		t.Errorf("期望 2 个模式，实际 %d 个", len(patterns))
	}

	// 测试不存在的文件
	nonExistentFile := filepath.Join(tempDir, ".nonexistent")
	patterns = LoadExcludeFromFileOrEmpty(nonExistentFile)
	if len(patterns) != 0 {
		t.Errorf("期望空列表，实际 %d 个模式", len(patterns))
	}
}

func TestParseIgnoreFileContent(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected []string
		desc     string
	}{
		{
			name:     "正常内容",
			content:  "*.tmp\n*.log\nvendor/",
			expected: []string{"*.tmp", "*.log", "vendor/"},
			desc:     "解析正常的忽略文件内容",
		},
		{
			name:     "包含注释和空行",
			content:  "# 注释\n*.tmp\n\n*.log\n# 另一个注释",
			expected: []string{"*.tmp", "*.log"},
			desc:     "跳过注释和空行",
		},
		{
			name:     "空内容",
			content:  "",
			expected: []string{},
			desc:     "空内容返回空列表",
		},
		{
			name:     "只有注释",
			content:  "# 注释1\n# 注释2",
			expected: []string{},
			desc:     "只有注释时返回空列表",
		},
		{
			name:     "只有空行",
			content:  "\n\n\n",
			expected: []string{},
			desc:     "只有空行时返回空列表",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseIgnoreFileContent(tt.content)
			if len(result) != len(tt.expected) {
				t.Errorf("%s: 期望 %d 个模式，实际 %d 个", tt.desc, len(tt.expected), len(result))
				return
			}

			for i, expected := range tt.expected {
				if result[i] != expected {
					t.Errorf("%s: 模式 %d 期望 '%s'，实际 '%s'", tt.desc, i, expected, result[i])
				}
			}
		})
	}
}

// 基准测试
func BenchmarkFilterOptions_ShouldSkip(b *testing.B) {
	filter := &FilterOptions{
		Include: []string{"*.go", "*.md"},
		Exclude: []string{"*_test.go", "vendor/*"},
		MinSize: 100,
		MaxSize: 1024 * 1024,
	}

	fileInfo := createFileInfo("main.go", 1000, false)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		filter.ShouldSkip("src/main.go", fileInfo)
	}
}

func BenchmarkFilterOptions_matchPattern(b *testing.B) {
	filter := &FilterOptions{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		filter.matchPattern("*.go", "main.go")
	}
}

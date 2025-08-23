package types

import (
	"testing"
)

// TestFilterOptions_matchAnyPattern 测试 matchAnyPattern 的模式和路径匹配情况
func TestFilterOptions_matchAnyPattern(t *testing.T) {
	tests := []struct {
		name     string
		patterns []string
		path     string
		expected bool
	}{
		// 基础测试用例
		{
			name:     "空模式列表",
			patterns: []string{},
			path:     "main.go",
			expected: false,
		},
		{
			name:     "空路径",
			patterns: []string{"*.go"},
			path:     "",
			expected: false,
		},
		{
			name:     "空模式和空路径",
			patterns: []string{},
			path:     "",
			expected: false,
		},

		// 文件扩展名匹配测试
		{
			name:     "Go文件扩展名匹配",
			patterns: []string{"*.go"},
			path:     "main.go",
			expected: true,
		},
		{
			name:     "多个扩展名匹配-匹配第一个",
			patterns: []string{"*.go", "*.md"},
			path:     "main.go",
			expected: true,
		},
		{
			name:     "多个扩展名匹配-匹配第二个",
			patterns: []string{"*.go", "*.md"},
			path:     "README.md",
			expected: true,
		},
		{
			name:     "扩展名不匹配",
			patterns: []string{"*.go"},
			path:     "main.js",
			expected: false,
		},
		{
			name:     "复杂扩展名匹配",
			patterns: []string{"*.test.go", "*.spec.js"},
			path:     "main.test.go",
			expected: true,
		},

		// 精确匹配测试
		{
			name:     "文件名精确匹配",
			patterns: []string{"main.go"},
			path:     "main.go",
			expected: true,
		},
		{
			name:     "文件名精确匹配-路径中的文件",
			patterns: []string{"main.go"},
			path:     "src/main.go",
			expected: true,
		},
		{
			name:     "文件名不匹配",
			patterns: []string{"main.go"},
			path:     "test.go",
			expected: false,
		},

		// 前缀匹配测试
		{
			name:     "文件名前缀匹配",
			patterns: []string{"test*"},
			path:     "test_main.go",
			expected: true,
		},
		{
			name:     "路径前缀匹配",
			patterns: []string{"src/*"},
			path:     "src/main.go",
			expected: true,
		},
		{
			name:     "深层路径前缀匹配",
			patterns: []string{"vendor/*"},
			path:     "project/vendor/package.go",
			expected: true,
		},
		{
			name:     "前缀不匹配",
			patterns: []string{"test*"},
			path:     "main.go",
			expected: false,
		},

		// 目录匹配测试
		{
			name:     "目录匹配-完整路径",
			patterns: []string{"vendor/"},
			path:     "vendor",
			expected: true,
		},
		{
			name:     "目录匹配-目录下文件",
			patterns: []string{"vendor/"},
			path:     "vendor/package.go",
			expected: true,
		},
		{
			name:     "目录匹配-嵌套目录",
			patterns: []string{"node_modules/"},
			path:     "src/node_modules/package.json",
			expected: true,
		},
		{
			name:     "目录不匹配",
			patterns: []string{"vendor/"},
			path:     "src/main.go",
			expected: false,
		},

		// 中间通配符匹配测试
		{
			name:     "中间通配符匹配文件名",
			patterns: []string{"*test*"},
			path:     "unit_test_main.go",
			expected: true,
		},
		{
			name:     "中间通配符匹配路径",
			patterns: []string{"*vendor*"},
			path:     "src/vendor_packages/main.go",
			expected: true,
		},
		{
			name:     "中间通配符不匹配",
			patterns: []string{"*test*"},
			path:     "main.go",
			expected: false,
		},

		// 复杂glob模式测试
		{
			name:     "问号通配符匹配",
			patterns: []string{"test?.go"},
			path:     "test1.go",
			expected: true,
		},
		{
			name:     "字符类匹配",
			patterns: []string{"test[0-9].go"},
			path:     "test5.go",
			expected: true,
		},
		{
			name:     "复杂路径模式",
			patterns: []string{"src/*/test/*.go"},
			path:     "src/main/test/unit.go",
			expected: true,
		},

		// 路径分隔符处理测试
		{
			name:     "Windows路径分隔符",
			patterns: []string{"src\\*.go"},
			path:     "src/main.go",
			expected: true,
		},
		{
			name:     "混合路径分隔符",
			patterns: []string{"src\\test/*.go"},
			path:     "src/test/main.go",
			expected: true,
		},
		{
			name:     "Unix路径分隔符",
			patterns: []string{"src/*.go"},
			path:     "src\\main.go",
			expected: true,
		},

		// 空模式处理测试
		{
			name:     "包含空模式的列表",
			patterns: []string{"", "*.go", ""},
			path:     "main.go",
			expected: true,
		},
		{
			name:     "全部为空模式",
			patterns: []string{"", "", ""},
			path:     "main.go",
			expected: false,
		},

		// 多模式组合测试
		{
			name:     "多种模式组合-匹配扩展名",
			patterns: []string{"*.js", "test*", "vendor/"},
			path:     "main.js",
			expected: true,
		},
		{
			name:     "多种模式组合-匹配前缀",
			patterns: []string{"*.js", "test*", "vendor/"},
			path:     "test_main.go",
			expected: true,
		},
		{
			name:     "多种模式组合-匹配目录",
			patterns: []string{"*.js", "test*", "vendor/"},
			path:     "vendor/package.go",
			expected: true,
		},
		{
			name:     "多种模式组合-都不匹配",
			patterns: []string{"*.js", "test*", "vendor/"},
			path:     "src/main.go",
			expected: false,
		},

		// 边界情况测试
		{
			name:     "单字符文件名",
			patterns: []string{"*"},
			path:     "a",
			expected: true,
		},
		{
			name:     "长路径匹配",
			patterns: []string{"*.go"},
			path:     "very/long/path/to/some/deep/directory/structure/main.go",
			expected: true,
		},
		{
			name:     "特殊字符文件名",
			patterns: []string{"test-*.go"},
			path:     "test-main.go",
			expected: true,
		},
		{
			name:     "数字文件名",
			patterns: []string{"*123*"},
			path:     "file123.txt",
			expected: true,
		},

		// 大小写敏感测试（Go的filepath.Match是大小写敏感的）
		{
			name:     "大小写敏感-匹配",
			patterns: []string{"*.GO"},
			path:     "main.GO",
			expected: true,
		},
		{
			name:     "大小写敏感-不匹配",
			patterns: []string{"*.GO"},
			path:     "main.go",
			expected: false,
		},

		// 相对路径和绝对路径测试
		{
			name:     "相对路径匹配",
			patterns: []string{"./src/*.go"},
			path:     "./src/main.go",
			expected: true,
		},
		{
			name:     "父目录路径匹配",
			patterns: []string{"../test/*.go"},
			path:     "../test/main.go",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &FilterOptions{}
			result := f.matchAnyPattern(tt.patterns, tt.path)
			if result != tt.expected {
				t.Errorf("matchAnyPattern(%v, %q) = %v, expected %v",
					tt.patterns, tt.path, result, tt.expected)
			}
		})
	}
}

// TestFilterOptions_matchAnyPattern_Performance 性能测试
func TestFilterOptions_matchAnyPattern_Performance(t *testing.T) {
	f := &FilterOptions{}
	patterns := []string{
		"*.go", "*.js", "*.ts", "*.py", "*.java",
		"test*", "spec*", "mock*",
		"vendor/", "node_modules/", ".git/",
		"*test*", "*spec*", "*mock*",
	}

	testPaths := []string{
		"main.go",
		"src/test/unit.go",
		"vendor/package.json",
		"node_modules/react/index.js",
		"test_helper.py",
		"spec_runner.js",
		"mock_data.json",
		"very/long/path/to/some/file.go",
	}

	// 预热
	for i := 0; i < 100; i++ {
		for _, path := range testPaths {
			f.matchAnyPattern(patterns, path)
		}
	}

	// 性能测试
	b := testing.B{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, path := range testPaths {
			f.matchAnyPattern(patterns, path)
		}
	}
}

// TestFilterOptions_matchAnyPattern_EdgeCases 边界情况测试
func TestFilterOptions_matchAnyPattern_EdgeCases(t *testing.T) {
	f := &FilterOptions{}

	// 测试nil patterns
	result := f.matchAnyPattern(nil, "main.go")
	if result != false {
		t.Errorf("matchAnyPattern(nil, \"main.go\") = %v, expected false", result)
	}

	// 测试超长模式
	longPattern := ""
	for i := 0; i < 1000; i++ {
		longPattern += "a"
	}
	longPattern += "*.go"
	result = f.matchAnyPattern([]string{longPattern}, "main.go")
	if result != false {
		t.Errorf("matchAnyPattern with very long pattern should return false")
	}

	// 测试超长路径
	longPath := ""
	for i := 0; i < 100; i++ {
		longPath += "dir/"
	}
	longPath += "main.go"
	result = f.matchAnyPattern([]string{"*.go"}, longPath)
	if result != true {
		t.Errorf("matchAnyPattern with very long path should match *.go pattern")
	}

	// 测试包含特殊字符的模式
	specialPatterns := []string{
		"[A-Z]*.go",    // 大写字母开头
		"*[0-9].go",    // 以数字结尾
		"test[abc].go", // 字符类
	}

	testCases := []struct {
		path     string
		expected bool
	}{
		{"Test1.go", true},  // 匹配 [A-Z]*.go (T是大写字母)
		{"main5.go", true},  // 匹配 *[0-9].go (以5结尾)
		{"testa.go", true},  // 匹配 test[abc].go (testa匹配)
		{"hello.go", false}, // 不匹配任何模式
	}

	for _, tc := range testCases {
		result = f.matchAnyPattern(specialPatterns, tc.path)
		if result != tc.expected {
			t.Errorf("matchAnyPattern(%v, %q) = %v, expected %v",
				specialPatterns, tc.path, result, tc.expected)
		}
	}
}

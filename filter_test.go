package comprx

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"
)

// TestLoadExcludeFromFile_Normal 测试正常功能
func TestLoadExcludeFromFile_Normal(t *testing.T) {
	// 创建临时测试文件
	content := `# 这是注释
*.log
*.tmp

# 另一个注释
temp/
*.bak
`
	tempFile := createTempFile(t, content)
	defer func() { _ = os.Remove(tempFile) }()

	patterns, err := LoadExcludeFromFile(tempFile)
	if err != nil {
		t.Fatalf("LoadExcludeFromFile() error = %v", err)
	}

	expected := []string{"*.log", "*.tmp", "temp/", "*.bak"}
	if !reflect.DeepEqual(patterns, expected) {
		t.Errorf("LoadExcludeFromFile() = %v, want %v", patterns, expected)
	}
}

// TestLoadExcludeFromFile_EmptyPath 测试空路径参数验证
func TestLoadExcludeFromFile_EmptyPath(t *testing.T) {
	_, err := LoadExcludeFromFile("")
	if err == nil {
		t.Error("LoadExcludeFromFile(\"\") should return error")
	}
	if !strings.Contains(err.Error(), "忽略文件路径不能为空") {
		t.Errorf("Expected empty path error, got: %v", err)
	}
}

// TestLoadExcludeFromFile_FileNotExist 测试文件不存在
func TestLoadExcludeFromFile_FileNotExist(t *testing.T) {
	_, err := LoadExcludeFromFile("nonexistent.ignore")
	if err == nil {
		t.Error("LoadExcludeFromFile() should return error for non-existent file")
	}
	if !strings.Contains(err.Error(), "忽略文件不存在") {
		t.Errorf("Expected file not exist error, got: %v", err)
	}
}

// TestLoadExcludeFromFile_EmptyFile 测试空文件
func TestLoadExcludeFromFile_EmptyFile(t *testing.T) {
	tempFile := createTempFile(t, "")
	defer func() { _ = os.Remove(tempFile) }()

	patterns, err := LoadExcludeFromFile(tempFile)
	if err != nil {
		t.Fatalf("LoadExcludeFromFile() error = %v", err)
	}

	if len(patterns) != 0 {
		t.Errorf("LoadExcludeFromFile() = %v, want empty slice", patterns)
	}
}

// TestLoadExcludeFromFile_OnlyComments 测试只有注释的文件
func TestLoadExcludeFromFile_OnlyComments(t *testing.T) {
	content := `# 注释1
# 注释2
# 注释3
`
	tempFile := createTempFile(t, content)
	defer func() { _ = os.Remove(tempFile) }()

	patterns, err := LoadExcludeFromFile(tempFile)
	if err != nil {
		t.Fatalf("LoadExcludeFromFile() error = %v", err)
	}

	if len(patterns) != 0 {
		t.Errorf("LoadExcludeFromFile() = %v, want empty slice", patterns)
	}
}

// TestLoadExcludeFromFile_OnlyEmptyLines 测试只有空行的文件
func TestLoadExcludeFromFile_OnlyEmptyLines(t *testing.T) {
	content := `

   
	
`
	tempFile := createTempFile(t, content)
	defer func() { _ = os.Remove(tempFile) }()

	patterns, err := LoadExcludeFromFile(tempFile)
	if err != nil {
		t.Fatalf("LoadExcludeFromFile() error = %v", err)
	}

	if len(patterns) != 0 {
		t.Errorf("LoadExcludeFromFile() = %v, want empty slice", patterns)
	}
}

// TestLoadExcludeFromFile_Deduplication 测试去重功能
func TestLoadExcludeFromFile_Deduplication(t *testing.T) {
	content := `*.log
*.tmp
*.log
temp/
*.tmp
*.bak
*.log
`
	tempFile := createTempFile(t, content)
	defer func() { _ = os.Remove(tempFile) }()

	patterns, err := LoadExcludeFromFile(tempFile)
	if err != nil {
		t.Fatalf("LoadExcludeFromFile() error = %v", err)
	}

	expected := []string{"*.log", "*.tmp", "temp/", "*.bak"}
	if !reflect.DeepEqual(patterns, expected) {
		t.Errorf("LoadExcludeFromFile() = %v, want %v", patterns, expected)
	}
}

// TestLoadExcludeFromFile_WhitespaceHandling 测试空白字符处理
func TestLoadExcludeFromFile_WhitespaceHandling(t *testing.T) {
	content := `  *.log  
	*.tmp	
 temp/ 
*.bak
`
	tempFile := createTempFile(t, content)
	defer func() { _ = os.Remove(tempFile) }()

	patterns, err := LoadExcludeFromFile(tempFile)
	if err != nil {
		t.Fatalf("LoadExcludeFromFile() error = %v", err)
	}

	expected := []string{"*.log", "*.tmp", "temp/", "*.bak"}
	if !reflect.DeepEqual(patterns, expected) {
		t.Errorf("LoadExcludeFromFile() = %v, want %v", patterns, expected)
	}
}

// TestLoadExcludeFromFile_InvalidGlobPattern 测试无效的glob模式
func TestLoadExcludeFromFile_InvalidGlobPattern(t *testing.T) {
	content := `*.log
[invalid
*.tmp
`
	tempFile := createTempFile(t, content)
	defer func() { _ = os.Remove(tempFile) }()

	_, err := LoadExcludeFromFile(tempFile)
	if err == nil {
		t.Error("LoadExcludeFromFile() should return error for invalid glob pattern")
	}
	if !strings.Contains(err.Error(), "无效的 glob 模式") {
		t.Errorf("Expected invalid glob pattern error, got: %v", err)
	}
	if !strings.Contains(err.Error(), "第 2 行") {
		t.Errorf("Expected line number in error, got: %v", err)
	}
}

// TestLoadExcludeFromFile_LargeFile 测试大文件处理
func TestLoadExcludeFromFile_LargeFile(t *testing.T) {
	// 创建包含1000行的大文件，每行都不同以测试容量预分配
	var lines []string
	for i := 0; i < 1000; i++ {
		// 确保每个模式都是唯一的：使用行号作为后缀，格式化为4位数字
		lines = append(lines, fmt.Sprintf("pattern_%04d_*.log", i))
	}
	content := strings.Join(lines, "\n")

	tempFile := createTempFile(t, content)
	defer func() { _ = os.Remove(tempFile) }()

	patterns, err := LoadExcludeFromFile(tempFile)
	if err != nil {
		t.Fatalf("LoadExcludeFromFile() error = %v", err)
	}

	if len(patterns) != 1000 {
		t.Errorf("LoadExcludeFromFile() returned %d patterns, want 1000", len(patterns))
	}
}

// TestLoadExcludeFromFile_MixedContent 测试混合内容
func TestLoadExcludeFromFile_MixedContent(t *testing.T) {
	content := `# 开头注释
*.log

# 中间注释
*.tmp
   # 带空格的注释
temp/

*.bak
# 结尾注释
`
	tempFile := createTempFile(t, content)
	defer func() { _ = os.Remove(tempFile) }()

	patterns, err := LoadExcludeFromFile(tempFile)
	if err != nil {
		t.Fatalf("LoadExcludeFromFile() error = %v", err)
	}

	expected := []string{"*.log", "*.tmp", "temp/", "*.bak"}
	if !reflect.DeepEqual(patterns, expected) {
		t.Errorf("LoadExcludeFromFile() = %v, want %v", patterns, expected)
	}
}

// TestLoadExcludeFromFileOrEmpty_Normal 测试OrEmpty版本正常功能
func TestLoadExcludeFromFileOrEmpty_Normal(t *testing.T) {
	content := `*.log
*.tmp
`
	tempFile := createTempFile(t, content)
	defer func() { _ = os.Remove(tempFile) }()

	patterns, err := LoadExcludeFromFileOrEmpty(tempFile)
	if err != nil {
		t.Fatalf("LoadExcludeFromFileOrEmpty() error = %v", err)
	}

	expected := []string{"*.log", "*.tmp"}
	if !reflect.DeepEqual(patterns, expected) {
		t.Errorf("LoadExcludeFromFileOrEmpty() = %v, want %v", patterns, expected)
	}
}

// TestLoadExcludeFromFileOrEmpty_FileNotExist 测试OrEmpty版本文件不存在
func TestLoadExcludeFromFileOrEmpty_FileNotExist(t *testing.T) {
	patterns, err := LoadExcludeFromFileOrEmpty("nonexistent.ignore")
	if err != nil {
		t.Fatalf("LoadExcludeFromFileOrEmpty() should not return error for non-existent file, got: %v", err)
	}

	if len(patterns) != 0 {
		t.Errorf("LoadExcludeFromFileOrEmpty() = %v, want empty slice", patterns)
	}
}

// TestLoadExcludeFromFileOrEmpty_EmptyPath 测试OrEmpty版本空路径
func TestLoadExcludeFromFileOrEmpty_EmptyPath(t *testing.T) {
	_, err := LoadExcludeFromFileOrEmpty("")
	if err == nil {
		t.Error("LoadExcludeFromFileOrEmpty(\"\") should return error")
	}
	if !strings.Contains(err.Error(), "忽略文件路径不能为空") {
		t.Errorf("Expected empty path error, got: %v", err)
	}
}

// TestLoadExcludeFromFileOrEmpty_InvalidGlob 测试OrEmpty版本无效glob
func TestLoadExcludeFromFileOrEmpty_InvalidGlob(t *testing.T) {
	content := `[invalid`
	tempFile := createTempFile(t, content)
	defer func() { _ = os.Remove(tempFile) }()

	_, err := LoadExcludeFromFileOrEmpty(tempFile)
	if err == nil {
		t.Error("LoadExcludeFromFileOrEmpty() should return error for invalid glob pattern")
	}
}

// BenchmarkLoadExcludeFromFile_Small 小文件性能测试
func BenchmarkLoadExcludeFromFile_Small(b *testing.B) {
	content := `*.log
*.tmp
temp/
*.bak
`
	tempFile := createTempFile(b, content)
	defer func() { _ = os.Remove(tempFile) }()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := LoadExcludeFromFile(tempFile)
		if err != nil {
			b.Fatalf("LoadExcludeFromFile() error = %v", err)
		}
	}
}

// BenchmarkLoadExcludeFromFile_Large 大文件性能测试
func BenchmarkLoadExcludeFromFile_Large(b *testing.B) {
	// 创建包含1000行的文件
	var lines []string
	for i := 0; i < 1000; i++ {
		lines = append(lines, "*.log"+string(rune('a'+i%26)))
	}
	content := strings.Join(lines, "\n")

	tempFile := createTempFile(b, content)
	defer func() { _ = os.Remove(tempFile) }()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := LoadExcludeFromFile(tempFile)
		if err != nil {
			b.Fatalf("LoadExcludeFromFile() error = %v", err)
		}
	}
}

// BenchmarkLoadExcludeFromFile_WithDuplicates 包含重复项的性能测试
func BenchmarkLoadExcludeFromFile_WithDuplicates(b *testing.B) {
	// 创建包含重复项的文件
	var lines []string
	patterns := []string{"*.log", "*.tmp", "temp/", "*.bak"}
	for i := 0; i < 1000; i++ {
		lines = append(lines, patterns[i%len(patterns)])
	}
	content := strings.Join(lines, "\n")

	tempFile := createTempFile(b, content)
	defer func() { _ = os.Remove(tempFile) }()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := LoadExcludeFromFile(tempFile)
		if err != nil {
			b.Fatalf("LoadExcludeFromFile() error = %v", err)
		}
	}
}

// 辅助函数：创建临时文件
func createTempFile(t testing.TB, content string) string {
	tempFile, err := os.CreateTemp("", "test_ignore_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	if _, err := tempFile.WriteString(content); err != nil {
		t.Fatalf("Failed to write temp file: %v", err)
	}

	if err := tempFile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	return tempFile.Name()
}

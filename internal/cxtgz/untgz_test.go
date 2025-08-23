package cxtgz

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gitee.com/MM-Q/comprx/internal/config"
)

func TestUntgz_SingleFile(t *testing.T) {
	tempDir := t.TempDir()

	// 创建测试文件
	testFile := filepath.Join(tempDir, "test.txt")
	testContent := "Hello, Untgz World!"
	if err := os.WriteFile(testFile, []byte(testContent), 0644); err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}

	// 压缩文件
	tgzFile := filepath.Join(tempDir, "test.tgz")
	cfg := config.New()

	if err := Tgz(tgzFile, testFile, cfg); err != nil {
		t.Fatalf("TGZ压缩失败: %v", err)
	}

	// 解压文件
	extractDir := filepath.Join(tempDir, "extract")
	if err := Untgz(tgzFile, extractDir, cfg); err != nil {
		t.Fatalf("TGZ解压失败: %v", err)
	}

	// 验证解压的文件
	extractedFile := filepath.Join(extractDir, "test.txt")
	content, err := os.ReadFile(extractedFile)
	if err != nil {
		t.Fatalf("读取解压文件失败: %v", err)
	}

	if string(content) != testContent {
		t.Errorf("解压内容 = %q, want %q", string(content), testContent)
	}
}

func TestUntgz_Directory(t *testing.T) {
	tempDir := t.TempDir()

	// 创建测试目录结构
	testDir := filepath.Join(tempDir, "testdir")
	if err := os.MkdirAll(testDir, 0755); err != nil {
		t.Fatalf("创建测试目录失败: %v", err)
	}

	// 创建测试文件
	files := map[string]string{
		"file1.txt":        "Content 1",
		"file2.txt":        "Content 2",
		"subdir/file3.txt": "Content 3",
	}

	for name, content := range files {
		filePath := filepath.Join(testDir, name)
		if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
			t.Fatalf("创建目录失败: %v", err)
		}
		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			t.Fatalf("创建文件失败: %v", err)
		}
	}

	// 压缩目录
	tgzFile := filepath.Join(tempDir, "test.tgz")
	cfg := config.New()

	if err := Tgz(tgzFile, testDir, cfg); err != nil {
		t.Fatalf("TGZ压缩失败: %v", err)
	}

	// 解压目录
	extractDir := filepath.Join(tempDir, "extract")
	if err := Untgz(tgzFile, extractDir, cfg); err != nil {
		t.Fatalf("TGZ解压失败: %v", err)
	}

	// 验证解压的文件
	for name, expectedContent := range files {
		extractedFile := filepath.Join(extractDir, "testdir", name)
		content, err := os.ReadFile(extractedFile)
		if err != nil {
			t.Fatalf("读取解压文件 %s 失败: %v", name, err)
		}

		if string(content) != expectedContent {
			t.Errorf("文件 %s 内容 = %q, want %q", name, string(content), expectedContent)
		}
	}
}

func TestUntgz_NonExistentFile(t *testing.T) {
	tempDir := t.TempDir()

	nonExistentFile := filepath.Join(tempDir, "nonexistent.tgz")
	extractDir := filepath.Join(tempDir, "extract")
	cfg := config.New()

	err := Untgz(nonExistentFile, extractDir, cfg)
	if err == nil {
		t.Errorf("期望解压不存在的文件时返回错误")
	}
}

func TestUntgz_InvalidTgzFile(t *testing.T) {
	tempDir := t.TempDir()

	// 创建无效的TGZ文件
	invalidFile := filepath.Join(tempDir, "invalid.tgz")
	if err := os.WriteFile(invalidFile, []byte("not a tgz file"), 0644); err != nil {
		t.Fatalf("创建无效文件失败: %v", err)
	}

	extractDir := filepath.Join(tempDir, "extract")
	cfg := config.New()

	err := Untgz(invalidFile, extractDir, cfg)
	if err == nil {
		t.Errorf("期望解压无效TGZ文件时返回错误")
	}
}

func TestUntgz_OverwriteExisting(t *testing.T) {
	tempDir := t.TempDir()

	// 创建测试文件
	testFile := filepath.Join(tempDir, "test.txt")
	testContent := "Original content"
	if err := os.WriteFile(testFile, []byte(testContent), 0644); err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}

	// 压缩文件
	tgzFile := filepath.Join(tempDir, "test.tgz")
	cfg := config.New()

	if err := Tgz(tgzFile, testFile, cfg); err != nil {
		t.Fatalf("TGZ压缩失败: %v", err)
	}

	// 创建解压目录和已存在的文件
	extractDir := filepath.Join(tempDir, "extract")
	if err := os.MkdirAll(extractDir, 0755); err != nil {
		t.Fatalf("创建解压目录失败: %v", err)
	}

	existingFile := filepath.Join(extractDir, "test.txt")
	existingContent := "Existing content"
	if err := os.WriteFile(existingFile, []byte(existingContent), 0644); err != nil {
		t.Fatalf("创建已存在文件失败: %v", err)
	}

	// 测试不覆盖
	cfg.OverwriteExisting = false

	err := Untgz(tgzFile, extractDir, cfg)
	if err == nil {
		t.Errorf("期望不覆盖已存在文件时返回错误")
	}

	// 验证文件未被覆盖
	content, err := os.ReadFile(existingFile)
	if err != nil {
		t.Fatalf("读取文件失败: %v", err)
	}
	if string(content) != existingContent {
		t.Errorf("文件被意外覆盖")
	}

	// 测试覆盖
	cfg.OverwriteExisting = true

	if err := Untgz(tgzFile, extractDir, cfg); err != nil {
		t.Fatalf("覆盖解压失败: %v", err)
	}

	// 验证文件被覆盖
	content, err = os.ReadFile(existingFile)
	if err != nil {
		t.Fatalf("读取覆盖文件失败: %v", err)
	}
	if string(content) != testContent {
		t.Errorf("文件内容 = %q, want %q", string(content), testContent)
	}
}

func TestUntgz_EmptyTgz(t *testing.T) {
	tempDir := t.TempDir()

	// 创建空目录
	emptyDir := filepath.Join(tempDir, "empty")
	if err := os.MkdirAll(emptyDir, 0755); err != nil {
		t.Fatalf("创建空目录失败: %v", err)
	}

	// 压缩空目录
	tgzFile := filepath.Join(tempDir, "empty.tgz")
	cfg := config.New()

	if err := Tgz(tgzFile, emptyDir, cfg); err != nil {
		t.Fatalf("压缩空目录失败: %v", err)
	}

	// 解压空TGZ
	extractDir := filepath.Join(tempDir, "extract")
	if err := Untgz(tgzFile, extractDir, cfg); err != nil {
		t.Fatalf("解压空TGZ失败: %v", err)
	}

	// 验证解压目录存在
	if _, err := os.Stat(extractDir); os.IsNotExist(err) {
		t.Errorf("解压目录未创建")
	}
}

func TestUntgz_LargeFile(t *testing.T) {
	tempDir := t.TempDir()

	// 创建大文件
	testFile := filepath.Join(tempDir, "large.txt")
	largeContent := strings.Repeat("Large file content for testing. ", 10000)
	if err := os.WriteFile(testFile, []byte(largeContent), 0644); err != nil {
		t.Fatalf("创建大文件失败: %v", err)
	}

	// 压缩大文件
	tgzFile := filepath.Join(tempDir, "large.tgz")
	cfg := config.New()

	if err := Tgz(tgzFile, testFile, cfg); err != nil {
		t.Fatalf("压缩大文件失败: %v", err)
	}

	// 解压大文件
	extractDir := filepath.Join(tempDir, "extract")
	if err := Untgz(tgzFile, extractDir, cfg); err != nil {
		t.Fatalf("解压大文件失败: %v", err)
	}

	// 验证解压的大文件
	extractedFile := filepath.Join(extractDir, "large.txt")
	content, err := os.ReadFile(extractedFile)
	if err != nil {
		t.Fatalf("读取解压大文件失败: %v", err)
	}

	if string(content) != largeContent {
		t.Errorf("大文件内容不匹配")
	}
}

// 基准测试
func BenchmarkUntgz_SmallFile(b *testing.B) {
	tempDir := b.TempDir()

	// 创建测试文件
	testFile := filepath.Join(tempDir, "test.txt")
	content := "Small file content for benchmarking."
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		b.Fatalf("创建测试文件失败: %v", err)
	}

	// 压缩文件
	tgzFile := filepath.Join(tempDir, "test.tgz")
	cfg := config.New()

	if err := Tgz(tgzFile, testFile, cfg); err != nil {
		b.Fatalf("压缩文件失败: %v", err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		extractDir := filepath.Join(tempDir, "extract_"+string(rune(i)))
		_ = Untgz(tgzFile, extractDir, cfg)
	}
}

func BenchmarkUntgz_LargeFile(b *testing.B) {
	tempDir := b.TempDir()

	// 创建大文件
	testFile := filepath.Join(tempDir, "large.txt")
	largeContent := strings.Repeat("Large file content for benchmarking. ", 5000)
	if err := os.WriteFile(testFile, []byte(largeContent), 0644); err != nil {
		b.Fatalf("创建大文件失败: %v", err)
	}

	// 压缩大文件
	tgzFile := filepath.Join(tempDir, "large.tgz")
	cfg := config.New()

	if err := Tgz(tgzFile, testFile, cfg); err != nil {
		b.Fatalf("压缩大文件失败: %v", err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		extractDir := filepath.Join(tempDir, "extract_"+string(rune(i)))
		_ = Untgz(tgzFile, extractDir, cfg)
	}
}

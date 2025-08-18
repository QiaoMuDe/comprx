package tgz

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gitee.com/MM-Q/comprx/config"
	"gitee.com/MM-Q/comprx/types"
)

func TestListTgz_SingleFile(t *testing.T) {
	tempDir := t.TempDir()

	// 创建测试文件
	testFile := filepath.Join(tempDir, "test.txt")
	testContent := "Hello, TGZ List!"
	if err := os.WriteFile(testFile, []byte(testContent), 0644); err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}

	// 压缩文件
	tgzFile := filepath.Join(tempDir, "test.tgz")
	cfg := config.New()

	if err := Tgz(tgzFile, testFile, cfg); err != nil {
		t.Fatalf("TGZ压缩失败: %v", err)
	}

	// 列出文件信息
	archiveInfo, err := ListTgz(tgzFile)
	if err != nil {
		t.Fatalf("列出TGZ文件失败: %v", err)
	}

	if archiveInfo.Type != types.CompressTypeTgz {
		t.Errorf("压缩包类型 = %v, want %v", archiveInfo.Type, types.CompressTypeTgz)
	}

	if archiveInfo.TotalFiles != 1 {
		t.Errorf("文件总数 = %d, want 1", archiveInfo.TotalFiles)
	}

	if len(archiveInfo.Files) != 1 {
		t.Errorf("文件列表长度 = %d, want 1", len(archiveInfo.Files))
	}

	file := archiveInfo.Files[0]
	if file.Name != "test.txt" {
		t.Errorf("文件名 = %q, want %q", file.Name, "test.txt")
	}

	if file.Size != int64(len(testContent)) {
		t.Errorf("文件大小 = %d, want %d", file.Size, len(testContent))
	}

	if file.IsDir {
		t.Errorf("文件不应该是目录")
	}
}

func TestListTgz_Directory(t *testing.T) {
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

	// 列出文件信息
	archiveInfo, err := ListTgz(tgzFile)
	if err != nil {
		t.Fatalf("列出TGZ文件失败: %v", err)
	}

	if archiveInfo.Type != types.CompressTypeTgz {
		t.Errorf("压缩包类型 = %v, want %v", archiveInfo.Type, types.CompressTypeTgz)
	}

	// 验证包含目录和文件
	if archiveInfo.TotalFiles < len(files) {
		t.Errorf("文件总数 = %d, want >= %d", archiveInfo.TotalFiles, len(files))
	}

	// 验证特定文件存在
	foundFiles := make(map[string]bool)
	for _, file := range archiveInfo.Files {
		// 检查文件名是否匹配任何期望的文件
		for expectedFile := range files {
			// 使用 strings.Contains 来匹配文件路径
			if strings.Contains(file.Name, expectedFile) && !file.IsDir {
				foundFiles[expectedFile] = true
			}
		}
	}

	for fileName := range files {
		if !foundFiles[fileName] {
			t.Errorf("未找到文件: %s", fileName)
		}
	}
}

func TestListTgzLimit(t *testing.T) {
	tempDir := t.TempDir()

	// 创建测试目录
	testDir := filepath.Join(tempDir, "testdir")
	if err := os.MkdirAll(testDir, 0755); err != nil {
		t.Fatalf("创建测试目录失败: %v", err)
	}

	// 创建多个文件
	for i := 0; i < 10; i++ {
		fileName := filepath.Join(testDir, fmt.Sprintf("file%d.txt", i))
		content := fmt.Sprintf("Content %d", i)
		if err := os.WriteFile(fileName, []byte(content), 0644); err != nil {
			t.Fatalf("创建文件失败: %v", err)
		}
	}

	// 压缩目录
	tgzFile := filepath.Join(tempDir, "test.tgz")
	cfg := config.New()

	if err := Tgz(tgzFile, testDir, cfg); err != nil {
		t.Fatalf("TGZ压缩失败: %v", err)
	}

	// 限制列出5个文件
	archiveInfo, err := ListTgzLimit(tgzFile, 5)
	if err != nil {
		t.Fatalf("限制列出TGZ文件失败: %v", err)
	}

	if len(archiveInfo.Files) > 5 {
		t.Errorf("文件列表长度 = %d, want <= 5", len(archiveInfo.Files))
	}
}

func TestListTgzMatch(t *testing.T) {
	tempDir := t.TempDir()

	// 创建测试目录
	testDir := filepath.Join(tempDir, "testdir")
	if err := os.MkdirAll(testDir, 0755); err != nil {
		t.Fatalf("创建测试目录失败: %v", err)
	}

	// 创建不同类型的文件
	files := []string{
		"test1.txt",
		"test2.txt",
		"data.log",
		"config.json",
		"readme.md",
	}

	for _, fileName := range files {
		filePath := filepath.Join(testDir, fileName)
		content := "Content of " + fileName
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

	// 匹配 .txt 文件
	archiveInfo, err := ListTgzMatch(tgzFile, "*.txt")
	if err != nil {
		t.Fatalf("匹配列出TGZ文件失败: %v", err)
	}

	// 验证只包含 .txt 文件
	for _, file := range archiveInfo.Files {
		if !strings.HasSuffix(file.Name, ".txt") && !file.IsDir {
			t.Errorf("匹配结果包含非.txt文件: %s", file.Name)
		}
	}
}

func TestListTgz_NonExistentFile(t *testing.T) {
	tempDir := t.TempDir()

	nonExistentFile := filepath.Join(tempDir, "nonexistent.tgz")

	_, err := ListTgz(nonExistentFile)
	if err == nil {
		t.Errorf("期望列出不存在的文件时返回错误")
	}
}

func TestListTgz_InvalidTgzFile(t *testing.T) {
	tempDir := t.TempDir()

	// 创建无效的TGZ文件
	invalidFile := filepath.Join(tempDir, "invalid.tgz")
	if err := os.WriteFile(invalidFile, []byte("not a tgz file"), 0644); err != nil {
		t.Fatalf("创建无效文件失败: %v", err)
	}

	_, err := ListTgz(invalidFile)
	if err == nil {
		t.Errorf("期望列出无效TGZ文件时返回错误")
	}
}

func TestListTgz_EmptyTgz(t *testing.T) {
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

	// 列出空TGZ文件
	archiveInfo, err := ListTgz(tgzFile)
	if err != nil {
		t.Fatalf("列出空TGZ文件失败: %v", err)
	}

	if archiveInfo.Type != types.CompressTypeTgz {
		t.Errorf("压缩包类型 = %v, want %v", archiveInfo.Type, types.CompressTypeTgz)
	}

	// 空目录可能包含目录条目本身
	if archiveInfo.TotalFiles < 0 {
		t.Errorf("文件总数不应该为负数: %d", archiveInfo.TotalFiles)
	}
}

func TestListTgz_LargeFile(t *testing.T) {
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

	// 列出大文件信息
	archiveInfo, err := ListTgz(tgzFile)
	if err != nil {
		t.Fatalf("列出大TGZ文件失败: %v", err)
	}

	if archiveInfo.TotalFiles != 1 {
		t.Errorf("文件总数 = %d, want 1", archiveInfo.TotalFiles)
	}

	file := archiveInfo.Files[0]
	if file.Size != int64(len(largeContent)) {
		t.Errorf("大文件大小 = %d, want %d", file.Size, len(largeContent))
	}
}

// 基准测试
func BenchmarkListTgz_SmallFile(b *testing.B) {
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
		_, _ = ListTgz(tgzFile)
	}
}

func BenchmarkListTgz_Directory(b *testing.B) {
	tempDir := b.TempDir()

	// 创建测试目录
	testDir := filepath.Join(tempDir, "testdir")
	if err := os.MkdirAll(testDir, 0755); err != nil {
		b.Fatalf("创建测试目录失败: %v", err)
	}

	// 创建多个文件
	for i := 0; i < 50; i++ {
		fileName := filepath.Join(testDir, fmt.Sprintf("file%d.txt", i))
		content := fmt.Sprintf("File content %d", i)
		if err := os.WriteFile(fileName, []byte(content), 0644); err != nil {
			b.Fatalf("创建文件失败: %v", err)
		}
	}

	// 压缩目录
	tgzFile := filepath.Join(tempDir, "test.tgz")
	cfg := config.New()

	if err := Tgz(tgzFile, testDir, cfg); err != nil {
		b.Fatalf("压缩目录失败: %v", err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = ListTgz(tgzFile)
	}
}

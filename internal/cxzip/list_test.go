package cxzip

import (
	"archive/zip"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gitee.com/MM-Q/comprx/types"
)

func TestListZip_SingleFile(t *testing.T) {
	tempDir := t.TempDir()
	zipFile := filepath.Join(tempDir, "test.zip")

	// 创建测试ZIP文件
	files := map[string]string{
		"test.txt": "Hello, World!",
	}
	createTestZip(t, zipFile, files)

	archiveInfo, err := ListZip(zipFile)
	if err != nil {
		t.Fatalf("列出ZIP文件失败: %v", err)
	}

	if archiveInfo.Type != types.CompressTypeZip {
		t.Errorf("压缩包类型 = %v, want %v", archiveInfo.Type, types.CompressTypeZip)
	}

	if archiveInfo.TotalFiles != 1 {
		t.Errorf("文件总数 = %d, want 1", archiveInfo.TotalFiles)
	}

	if len(archiveInfo.Files) != 1 {
		t.Errorf("文件列表长度 = %d, want 1", len(archiveInfo.Files))
	}

	file := archiveInfo.Files[0]
	if file.Name != "test.txt" {
		t.Errorf("文件名 = %s, want test.txt", file.Name)
	}

	if file.Size != int64(len("Hello, World!")) {
		t.Errorf("文件大小 = %d, want %d", file.Size, len("Hello, World!"))
	}
}

func TestListZip_MultipleFiles(t *testing.T) {
	tempDir := t.TempDir()
	zipFile := filepath.Join(tempDir, "test.zip")

	// 创建测试ZIP文件
	files := map[string]string{
		"file1.txt":            "Content 1",
		"file2.txt":            "Content 2",
		"dir/file3.txt":        "Content 3",
		"dir/subdir/file4.txt": "Content 4",
	}
	createTestZip(t, zipFile, files)

	archiveInfo, err := ListZip(zipFile)
	if err != nil {
		t.Fatalf("列出ZIP文件失败: %v", err)
	}

	if archiveInfo.TotalFiles != 4 {
		t.Errorf("文件总数 = %d, want 4", archiveInfo.TotalFiles)
	}

	if len(archiveInfo.Files) != 4 {
		t.Errorf("文件列表长度 = %d, want 4", len(archiveInfo.Files))
	}

	// 验证文件名
	expectedFiles := []string{"file1.txt", "file2.txt", "dir/file3.txt", "dir/subdir/file4.txt"}
	for _, expected := range expectedFiles {
		found := false
		for _, file := range archiveInfo.Files {
			if file.Name == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("未找到预期文件: %s", expected)
		}
	}
}

func TestListZip_WithDirectories(t *testing.T) {
	tempDir := t.TempDir()
	zipFile := filepath.Join(tempDir, "test.zip")

	// 创建包含目录的ZIP文件
	zipFileHandle, err := os.Create(zipFile)
	if err != nil {
		t.Fatalf("创建ZIP文件失败: %v", err)
	}
	defer func() { _ = zipFileHandle.Close() }()

	zipWriter := zip.NewWriter(zipFileHandle)
	defer func() { _ = zipWriter.Close() }()

	// 添加目录
	dirHeader := &zip.FileHeader{
		Name: "testdir/",
	}
	dirHeader.SetMode(os.ModeDir | 0755)
	if _, err := zipWriter.CreateHeader(dirHeader); err != nil {
		t.Fatalf("创建目录条目失败: %v", err)
	}

	// 添加文件
	fileWriter, err := zipWriter.Create("testdir/file.txt")
	if err != nil {
		t.Fatalf("创建文件条目失败: %v", err)
	}
	if _, err := fileWriter.Write([]byte("test content")); err != nil {
		t.Fatalf("写入文件内容失败: %v", err)
	}

	if err := zipWriter.Close(); err != nil {
		t.Fatalf("关闭 ZIP writer 失败: %v", err)
	}
	if err := zipFileHandle.Close(); err != nil {
		t.Fatalf("关闭 ZIP 文件失败: %v", err)
	}

	archiveInfo, err := ListZip(zipFile)
	if err != nil {
		t.Fatalf("列出ZIP文件失败: %v", err)
	}

	if archiveInfo.TotalFiles != 2 {
		t.Errorf("文件总数 = %d, want 2", archiveInfo.TotalFiles)
	}

	// 检查是否包含目录和文件
	hasDir := false
	hasFile := false
	for _, file := range archiveInfo.Files {
		if file.Name == "testdir/" && file.IsDir {
			hasDir = true
		}
		if file.Name == "testdir/file.txt" && !file.IsDir {
			hasFile = true
		}
	}

	if !hasDir {
		t.Error("未找到目录条目")
	}
	if !hasFile {
		t.Error("未找到文件条目")
	}
}

func TestListZipLimit(t *testing.T) {
	tempDir := t.TempDir()
	zipFile := filepath.Join(tempDir, "test.zip")

	// 创建包含多个文件的ZIP
	files := map[string]string{
		"file1.txt": "Content 1",
		"file2.txt": "Content 2",
		"file3.txt": "Content 3",
		"file4.txt": "Content 4",
		"file5.txt": "Content 5",
	}
	createTestZip(t, zipFile, files)

	// 测试限制为3个文件
	archiveInfo, err := ListZipLimit(zipFile, 3)
	if err != nil {
		t.Fatalf("列出限制文件失败: %v", err)
	}

	if len(archiveInfo.Files) != 3 {
		t.Errorf("限制文件列表长度 = %d, want 3", len(archiveInfo.Files))
	}

	// 测试限制为0（应该返回所有文件）
	archiveInfo, err = ListZipLimit(zipFile, 0)
	if err != nil {
		t.Fatalf("列出所有文件失败: %v", err)
	}

	if len(archiveInfo.Files) != 5 {
		t.Errorf("所有文件列表长度 = %d, want 5", len(archiveInfo.Files))
	}
}

func TestListZipMatch(t *testing.T) {
	tempDir := t.TempDir()
	zipFile := filepath.Join(tempDir, "test.zip")

	// 创建包含不同类型文件的ZIP
	files := map[string]string{
		"file1.txt":     "Content 1",
		"file2.go":      "package main",
		"file3.txt":     "Content 3",
		"dir/file4.go":  "package dir",
		"dir/file5.txt": "Content 5",
	}
	createTestZip(t, zipFile, files)

	// 测试匹配 .txt 文件
	archiveInfo, err := ListZipMatch(zipFile, "*.txt")
	if err != nil {
		t.Fatalf("匹配文件失败: %v", err)
	}

	expectedTxtFiles := 3
	if len(archiveInfo.Files) != expectedTxtFiles {
		t.Errorf("匹配的txt文件数量 = %d, want %d", len(archiveInfo.Files), expectedTxtFiles)
	}

	// 验证所有匹配的文件都是.txt文件
	for _, file := range archiveInfo.Files {
		if !strings.HasSuffix(file.Name, ".txt") {
			t.Errorf("匹配的文件不是txt文件: %s", file.Name)
		}
	}

	// 测试匹配 .go 文件
	archiveInfo, err = ListZipMatch(zipFile, "*.go")
	if err != nil {
		t.Fatalf("匹配go文件失败: %v", err)
	}

	expectedGoFiles := 2
	if len(archiveInfo.Files) != expectedGoFiles {
		t.Errorf("匹配的go文件数量 = %d, want %d", len(archiveInfo.Files), expectedGoFiles)
	}

	// 测试不匹配的模式
	archiveInfo, err = ListZipMatch(zipFile, "*.nonexistent")
	if err != nil {
		t.Fatalf("匹配不存在文件失败: %v", err)
	}

	if len(archiveInfo.Files) != 0 {
		t.Errorf("不匹配的文件数量 = %d, want 0", len(archiveInfo.Files))
	}
}

func TestListZip_NonExistentFile(t *testing.T) {
	tempDir := t.TempDir()
	nonExistentZip := filepath.Join(tempDir, "nonexistent.zip")

	_, err := ListZip(nonExistentZip)
	if err == nil {
		t.Error("应该返回错误，因为ZIP文件不存在")
	}
}

func TestListZip_InvalidZip(t *testing.T) {
	tempDir := t.TempDir()
	invalidZip := filepath.Join(tempDir, "invalid.zip")

	// 创建无效的ZIP文件
	if err := os.WriteFile(invalidZip, []byte("not a zip file"), 0644); err != nil {
		t.Fatalf("创建无效ZIP文件失败: %v", err)
	}

	_, err := ListZip(invalidZip)
	if err == nil {
		t.Error("应该返回错误，因为ZIP文件无效")
	}
}

func TestListZip_EmptyZip(t *testing.T) {
	tempDir := t.TempDir()
	zipFile := filepath.Join(tempDir, "empty.zip")

	// 创建空的ZIP文件
	zipFileHandle, err := os.Create(zipFile)
	if err != nil {
		t.Fatalf("创建ZIP文件失败: %v", err)
	}
	defer func() { _ = zipFileHandle.Close() }()

	zipWriter := zip.NewWriter(zipFileHandle)
	if err := zipWriter.Close(); err != nil {
		t.Fatalf("关闭 ZIP writer 失败: %v", err)
	}
	_ = zipFileHandle.Close()

	archiveInfo, err := ListZip(zipFile)
	if err != nil {
		t.Fatalf("列出空ZIP文件失败: %v", err)
	}

	if archiveInfo.TotalFiles != 0 {
		t.Errorf("空ZIP文件总数 = %d, want 0", archiveInfo.TotalFiles)
	}

	if len(archiveInfo.Files) != 0 {
		t.Errorf("空ZIP文件列表长度 = %d, want 0", len(archiveInfo.Files))
	}
}

// 基准测试
func BenchmarkListZip_SmallFiles(b *testing.B) {
	tempDir := b.TempDir()
	zipFile := filepath.Join(tempDir, "test.zip")

	// 创建测试ZIP文件
	files := map[string]string{
		"file1.txt": "Content 1",
		"file2.txt": "Content 2",
		"file3.txt": "Content 3",
	}
	createTestZipForBench(b, zipFile, files)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ListZip(zipFile)
	}
}

func BenchmarkListZip_ManyFiles(b *testing.B) {
	tempDir := b.TempDir()
	zipFile := filepath.Join(tempDir, "many.zip")

	// 创建包含很多文件的ZIP
	files := make(map[string]string)
	for i := 0; i < 100; i++ {
		files[filepath.Join("dir", "file"+string(rune(i))+".txt")] = "Content " + string(rune(i))
	}
	createTestZipForBench(b, zipFile, files)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ListZip(zipFile)
	}
}

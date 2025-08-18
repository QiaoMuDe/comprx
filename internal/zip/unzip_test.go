package zip

import (
	"archive/zip"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gitee.com/MM-Q/comprx/config"
	"gitee.com/MM-Q/comprx/internal/utils"
)

// 创建测试ZIP文件的辅助函数
func createTestZip(t *testing.T, zipPath string, files map[string]string) {
	zipFile, err := os.Create(zipPath)
	if err != nil {
		t.Fatalf("创建ZIP文件失败: %v", err)
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	for name, content := range files {
		writer, err := zipWriter.Create(name)
		if err != nil {
			t.Fatalf("创建ZIP条目失败: %v", err)
		}

		if _, err := writer.Write([]byte(content)); err != nil {
			t.Fatalf("写入ZIP条目失败: %v", err)
		}
	}
}

func TestUnzip_SingleFile(t *testing.T) {
	tempDir := t.TempDir()
	zipFile := filepath.Join(tempDir, "test.zip")
	extractDir := filepath.Join(tempDir, "extract")

	// 创建测试ZIP文件
	files := map[string]string{
		"test.txt": "Hello, World!",
	}
	createTestZip(t, zipFile, files)

	cfg := config.New()
	err := Unzip(zipFile, extractDir, cfg)
	if err != nil {
		t.Fatalf("解压失败: %v", err)
	}

	// 验证解压结果
	extractedFile := filepath.Join(extractDir, "test.txt")
	if !utils.Exists(extractedFile) {
		t.Error("解压的文件不存在")
	}

	content, err := os.ReadFile(extractedFile)
	if err != nil {
		t.Fatalf("读取解压文件失败: %v", err)
	}

	if string(content) != "Hello, World!" {
		t.Errorf("解压文件内容 = %s, want Hello, World!", string(content))
	}
}

func TestUnzip_MultipleFiles(t *testing.T) {
	tempDir := t.TempDir()
	zipFile := filepath.Join(tempDir, "test.zip")
	extractDir := filepath.Join(tempDir, "extract")

	// 创建测试ZIP文件
	files := map[string]string{
		"file1.txt":            "Content 1",
		"file2.txt":            "Content 2",
		"dir/file3.txt":        "Content 3",
		"dir/subdir/file4.txt": "Content 4",
	}
	createTestZip(t, zipFile, files)

	cfg := config.New()
	err := Unzip(zipFile, extractDir, cfg)
	if err != nil {
		t.Fatalf("解压失败: %v", err)
	}

	// 验证所有文件都被解压
	for fileName, expectedContent := range files {
		extractedFile := filepath.Join(extractDir, fileName)
		if !utils.Exists(extractedFile) {
			t.Errorf("解压的文件不存在: %s", fileName)
			continue
		}

		content, err := os.ReadFile(extractedFile)
		if err != nil {
			t.Errorf("读取解压文件失败 %s: %v", fileName, err)
			continue
		}

		if string(content) != expectedContent {
			t.Errorf("文件 %s 内容 = %s, want %s", fileName, string(content), expectedContent)
		}
	}
}

func TestUnzip_WithDirectories(t *testing.T) {
	tempDir := t.TempDir()
	zipFile := filepath.Join(tempDir, "test.zip")
	extractDir := filepath.Join(tempDir, "extract")

	// 创建包含目录的ZIP文件
	zipFileHandle, err := os.Create(zipFile)
	if err != nil {
		t.Fatalf("创建ZIP文件失败: %v", err)
	}
	defer zipFileHandle.Close()

	zipWriter := zip.NewWriter(zipFileHandle)
	defer zipWriter.Close()

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

	zipWriter.Close()
	zipFileHandle.Close()

	cfg := config.New()
	err = Unzip(zipFile, extractDir, cfg)
	if err != nil {
		t.Fatalf("解压失败: %v", err)
	}

	// 验证目录和文件都被创建
	extractedDir := filepath.Join(extractDir, "testdir")
	if !utils.Exists(extractedDir) {
		t.Error("解压的目录不存在")
	}

	extractedFile := filepath.Join(extractDir, "testdir", "file.txt")
	if !utils.Exists(extractedFile) {
		t.Error("解压的文件不存在")
	}
}

func TestUnzip_EmptyFiles(t *testing.T) {
	tempDir := t.TempDir()
	zipFile := filepath.Join(tempDir, "test.zip")
	extractDir := filepath.Join(tempDir, "extract")

	// 创建包含空文件的ZIP
	files := map[string]string{
		"empty.txt":    "",
		"nonempty.txt": "content",
	}
	createTestZip(t, zipFile, files)

	cfg := config.New()
	err := Unzip(zipFile, extractDir, cfg)
	if err != nil {
		t.Fatalf("解压失败: %v", err)
	}

	// 验证空文件被正确创建
	emptyFile := filepath.Join(extractDir, "empty.txt")
	if !utils.Exists(emptyFile) {
		t.Error("空文件未被创建")
	}

	stat, err := os.Stat(emptyFile)
	if err != nil {
		t.Fatalf("获取空文件信息失败: %v", err)
	}

	if stat.Size() != 0 {
		t.Errorf("空文件大小 = %d, want 0", stat.Size())
	}
}

func TestUnzip_NonExistentZip(t *testing.T) {
	tempDir := t.TempDir()
	nonExistentZip := filepath.Join(tempDir, "nonexistent.zip")
	extractDir := filepath.Join(tempDir, "extract")

	cfg := config.New()
	err := Unzip(nonExistentZip, extractDir, cfg)
	if err == nil {
		t.Error("应该返回错误，因为ZIP文件不存在")
	}
}

func TestUnzip_InvalidZip(t *testing.T) {
	tempDir := t.TempDir()
	invalidZip := filepath.Join(tempDir, "invalid.zip")
	extractDir := filepath.Join(tempDir, "extract")

	// 创建无效的ZIP文件
	if err := os.WriteFile(invalidZip, []byte("not a zip file"), 0644); err != nil {
		t.Fatalf("创建无效ZIP文件失败: %v", err)
	}

	cfg := config.New()
	err := Unzip(invalidZip, extractDir, cfg)
	if err == nil {
		t.Error("应该返回错误，因为ZIP文件无效")
	}
}

func TestUnzip_CreateTargetDirectory(t *testing.T) {
	tempDir := t.TempDir()
	zipFile := filepath.Join(tempDir, "test.zip")
	extractDir := filepath.Join(tempDir, "nonexistent", "extract")

	// 创建测试ZIP文件
	files := map[string]string{
		"test.txt": "Hello, World!",
	}
	createTestZip(t, zipFile, files)

	cfg := config.New()
	err := Unzip(zipFile, extractDir, cfg)
	if err != nil {
		t.Fatalf("解压失败: %v", err)
	}

	// 验证目标目录被创建
	if !utils.Exists(extractDir) {
		t.Error("目标目录未被创建")
	}

	// 验证文件被解压
	extractedFile := filepath.Join(extractDir, "test.txt")
	if !utils.Exists(extractedFile) {
		t.Error("解压的文件不存在")
	}
}

func TestUnzip_SymbolicLinks(t *testing.T) {
	// 在Windows上跳过符号链接测试
	if os.Getenv("GOOS") == "windows" {
		t.Skip("跳过Windows上的符号链接测试")
	}

	tempDir := t.TempDir()
	zipFile := filepath.Join(tempDir, "test.zip")
	extractDir := filepath.Join(tempDir, "extract")

	// 创建包含符号链接的ZIP文件
	zipFileHandle, err := os.Create(zipFile)
	if err != nil {
		t.Fatalf("创建ZIP文件失败: %v", err)
	}
	defer zipFileHandle.Close()

	zipWriter := zip.NewWriter(zipFileHandle)
	defer zipWriter.Close()

	// 添加普通文件
	fileWriter, err := zipWriter.Create("target.txt")
	if err != nil {
		t.Fatalf("创建文件条目失败: %v", err)
	}
	if _, err := fileWriter.Write([]byte("target content")); err != nil {
		t.Fatalf("写入文件内容失败: %v", err)
	}

	// 添加符号链接
	linkHeader := &zip.FileHeader{
		Name: "link.txt",
	}
	linkHeader.SetMode(os.ModeSymlink | 0777)
	linkWriter, err := zipWriter.CreateHeader(linkHeader)
	if err != nil {
		t.Fatalf("创建符号链接条目失败: %v", err)
	}
	if _, err := linkWriter.Write([]byte("target.txt")); err != nil {
		t.Fatalf("写入符号链接目标失败: %v", err)
	}

	zipWriter.Close()
	zipFileHandle.Close()

	cfg := config.New()
	err = Unzip(zipFile, extractDir, cfg)
	if err != nil {
		t.Fatalf("解压失败: %v", err)
	}

	// 验证符号链接被正确创建
	linkFile := filepath.Join(extractDir, "link.txt")
	if !utils.Exists(linkFile) {
		t.Error("符号链接未被创建")
	}

	// 验证符号链接指向正确的目标
	target, err := os.Readlink(linkFile)
	if err != nil {
		t.Fatalf("读取符号链接目标失败: %v", err)
	}

	if target != "target.txt" {
		t.Errorf("符号链接目标 = %s, want target.txt", target)
	}
}

func TestUnzip_LargeFile(t *testing.T) {
	tempDir := t.TempDir()
	zipFile := filepath.Join(tempDir, "large.zip")
	extractDir := filepath.Join(tempDir, "extract")

	// 创建包含大文件的ZIP
	largeContent := strings.Repeat("This is a large file content. ", 10000)
	files := map[string]string{
		"large.txt": largeContent,
	}
	createTestZip(t, zipFile, files)

	cfg := config.New()
	err := Unzip(zipFile, extractDir, cfg)
	if err != nil {
		t.Fatalf("解压大文件失败: %v", err)
	}

	// 验证大文件内容
	extractedFile := filepath.Join(extractDir, "large.txt")
	content, err := os.ReadFile(extractedFile)
	if err != nil {
		t.Fatalf("读取解压的大文件失败: %v", err)
	}

	if string(content) != largeContent {
		t.Error("解压的大文件内容不匹配")
	}
}

// 为基准测试创建ZIP文件的辅助函数
func createTestZipForBench(b *testing.B, zipPath string, files map[string]string) {
	zipFile, err := os.Create(zipPath)
	if err != nil {
		b.Fatalf("创建ZIP文件失败: %v", err)
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	for name, content := range files {
		writer, err := zipWriter.Create(name)
		if err != nil {
			b.Fatalf("创建ZIP条目失败: %v", err)
		}

		if _, err := writer.Write([]byte(content)); err != nil {
			b.Fatalf("写入ZIP条目失败: %v", err)
		}
	}
}

// 基准测试
func BenchmarkUnzip_SmallFiles(b *testing.B) {
	tempDir := b.TempDir()
	zipFile := filepath.Join(tempDir, "test.zip")

	// 创建测试ZIP文件
	files := map[string]string{
		"file1.txt": "Content 1",
		"file2.txt": "Content 2",
		"file3.txt": "Content 3",
	}
	createTestZipForBench(b, zipFile, files)

	cfg := config.New()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		extractDir := filepath.Join(tempDir, "extract_"+string(rune(i)))
		_ = Unzip(zipFile, extractDir, cfg)
	}
}

func BenchmarkUnzip_LargeFile(b *testing.B) {
	tempDir := b.TempDir()
	zipFile := filepath.Join(tempDir, "large.zip")

	// 创建包含大文件的ZIP
	largeContent := strings.Repeat("Large file content for benchmarking. ", 5000)
	files := map[string]string{
		"large.txt": largeContent,
	}
	createTestZipForBench(b, zipFile, files)

	cfg := config.New()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		extractDir := filepath.Join(tempDir, "extract_"+string(rune(i)))
		_ = Unzip(zipFile, extractDir, cfg)
	}
}

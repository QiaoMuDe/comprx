package tar

import (
	"archive/tar"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gitee.com/MM-Q/comprx/config"
)

func TestTar_SingleFile(t *testing.T) {
	tempDir := t.TempDir()

	// 创建测试文件
	testFile := filepath.Join(tempDir, "test.txt")
	testContent := "Hello, TAR World!"
	if err := os.WriteFile(testFile, []byte(testContent), 0644); err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}

	// 压缩文件
	tarFile := filepath.Join(tempDir, "test.tar")
	cfg := config.New()

	if err := Tar(tarFile, testFile, cfg); err != nil {
		t.Fatalf("TAR压缩失败: %v", err)
	}

	// 验证压缩文件存在
	if _, err := os.Stat(tarFile); os.IsNotExist(err) {
		t.Errorf("TAR文件未创建")
	}

	// 验证文件大小
	info, err := os.Stat(tarFile)
	if err != nil {
		t.Fatalf("获取TAR文件信息失败: %v", err)
	}

	if info.Size() == 0 {
		t.Errorf("TAR文件大小为0")
	}
}

func TestTar_Directory(t *testing.T) {
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
	tarFile := filepath.Join(tempDir, "test.tar")
	cfg := config.New()

	if err := Tar(tarFile, testDir, cfg); err != nil {
		t.Fatalf("TAR压缩失败: %v", err)
	}

	// 验证压缩文件存在
	if _, err := os.Stat(tarFile); os.IsNotExist(err) {
		t.Errorf("TAR文件未创建")
	}
}

func TestTar_EmptyDirectory(t *testing.T) {
	tempDir := t.TempDir()

	// 创建空目录
	emptyDir := filepath.Join(tempDir, "empty")
	if err := os.MkdirAll(emptyDir, 0755); err != nil {
		t.Fatalf("创建空目录失败: %v", err)
	}

	// 压缩空目录
	tarFile := filepath.Join(tempDir, "empty.tar")
	cfg := config.New()

	if err := Tar(tarFile, emptyDir, cfg); err != nil {
		t.Fatalf("TAR压缩空目录失败: %v", err)
	}

	// 验证压缩文件存在
	if _, err := os.Stat(tarFile); os.IsNotExist(err) {
		t.Errorf("TAR文件未创建")
	}
}

func TestTar_NonExistentSource(t *testing.T) {
	tempDir := t.TempDir()

	nonExistentPath := filepath.Join(tempDir, "nonexistent")
	tarFile := filepath.Join(tempDir, "test.tar")
	cfg := config.New()

	err := Tar(tarFile, nonExistentPath, cfg)
	if err == nil {
		t.Errorf("期望压缩不存在的文件时返回错误")
	}
}

func TestTar_InvalidDestination(t *testing.T) {
	tempDir := t.TempDir()

	// 创建测试文件
	testFile := filepath.Join(tempDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}

	// 使用无效的目标路径
	invalidPath := filepath.Join(tempDir, "nonexistent", "test.tar")
	cfg := config.New()

	err := Tar(invalidPath, testFile, cfg)
	if err != nil {
		// 应该自动创建目录，所以不应该出错
		// 但如果出错，检查是否是预期的错误类型
		if !strings.Contains(err.Error(), "创建目标目录失败") {
			t.Errorf("意外的错误类型: %v", err)
		}
	}
}

func TestTar_OverwriteExisting(t *testing.T) {
	tempDir := t.TempDir()

	// 创建测试文件
	testFile := filepath.Join(tempDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}

	tarFile := filepath.Join(tempDir, "test.tar")

	// 创建已存在的目标文件
	if err := os.WriteFile(tarFile, []byte("existing"), 0644); err != nil {
		t.Fatalf("创建已存在文件失败: %v", err)
	}

	// 测试不覆盖
	cfg := config.New()
	cfg.OverwriteExisting = false

	err := Tar(tarFile, testFile, cfg)
	if err == nil {
		t.Errorf("期望不覆盖已存在文件时返回错误")
	}

	// 测试覆盖
	cfg.OverwriteExisting = true

	if err := Tar(tarFile, testFile, cfg); err != nil {
		t.Fatalf("覆盖已存在文件失败: %v", err)
	}
}

func TestTar_LargeFile(t *testing.T) {
	tempDir := t.TempDir()

	// 创建大文件
	testFile := filepath.Join(tempDir, "large.txt")
	largeContent := strings.Repeat("Large file content for testing. ", 10000)
	if err := os.WriteFile(testFile, []byte(largeContent), 0644); err != nil {
		t.Fatalf("创建大文件失败: %v", err)
	}

	// 压缩大文件
	tarFile := filepath.Join(tempDir, "large.tar")
	cfg := config.New()

	if err := Tar(tarFile, testFile, cfg); err != nil {
		t.Fatalf("压缩大文件失败: %v", err)
	}

	// 验证文件存在
	if _, err := os.Stat(tarFile); os.IsNotExist(err) {
		t.Errorf("TAR文件未创建")
	}
}

func TestTar_SymbolicLinks(t *testing.T) {
	// 在Windows上跳过符号链接测试
	if os.Getenv("GOOS") == "windows" {
		t.Skip("跳过Windows上的符号链接测试")
	}

	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.txt")
	linkFile := filepath.Join(tempDir, "link.txt")
	tarFile := filepath.Join(tempDir, "test.tar")

	// 创建测试文件
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}

	// 创建符号链接
	if err := os.Symlink(testFile, linkFile); err != nil {
		t.Skipf("创建符号链接失败，跳过测试: %v", err)
	}

	cfg := config.New()
	err := Tar(tarFile, linkFile, cfg)
	if err != nil {
		t.Fatalf("压缩符号链接失败: %v", err)
	}

	// 验证TAR文件内容
	file, err := os.Open(tarFile)
	if err != nil {
		t.Fatalf("打开TAR文件失败: %v", err)
	}
	defer file.Close()

	tarReader := tar.NewReader(file)
	header, err := tarReader.Next()
	if err != nil {
		t.Fatalf("读取TAR文件头失败: %v", err)
	}

	if header.Name != "link.txt" {
		t.Errorf("TAR文件中的文件名 = %s, want link.txt", header.Name)
	}
}

// 基准测试
func BenchmarkTar_SmallFile(b *testing.B) {
	tempDir := b.TempDir()

	// 创建测试文件
	testFile := filepath.Join(tempDir, "test.txt")
	content := "Small file content for benchmarking."
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		b.Fatalf("创建测试文件失败: %v", err)
	}

	cfg := config.New()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		tarFile := filepath.Join(tempDir, fmt.Sprintf("bench_%d.tar", i))
		_ = Tar(tarFile, testFile, cfg)
	}
}

func BenchmarkTar_LargeFile(b *testing.B) {
	tempDir := b.TempDir()

	// 创建大文件
	testFile := filepath.Join(tempDir, "large.txt")
	largeContent := strings.Repeat("Large file content for benchmarking. ", 5000)
	if err := os.WriteFile(testFile, []byte(largeContent), 0644); err != nil {
		b.Fatalf("创建大文件失败: %v", err)
	}

	cfg := config.New()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		tarFile := filepath.Join(tempDir, fmt.Sprintf("bench_large_%d.tar", i))
		_ = Tar(tarFile, testFile, cfg)
	}
}

func BenchmarkTar_Directory(b *testing.B) {
	tempDir := b.TempDir()

	// 创建测试目录
	testDir := filepath.Join(tempDir, "testdir")
	if err := os.MkdirAll(testDir, 0755); err != nil {
		b.Fatalf("创建测试目录失败: %v", err)
	}

	// 创建多个文件
	for i := 0; i < 10; i++ {
		fileName := filepath.Join(testDir, fmt.Sprintf("file_%d.txt", i))
		content := fmt.Sprintf("File content %d", i)
		if err := os.WriteFile(fileName, []byte(content), 0644); err != nil {
			b.Fatalf("创建文件失败: %v", err)
		}
	}

	cfg := config.New()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		tarFile := filepath.Join(tempDir, fmt.Sprintf("bench_dir_%d.tar", i))
		_ = Tar(tarFile, testDir, cfg)
	}
}

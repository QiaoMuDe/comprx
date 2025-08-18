package tar

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gitee.com/MM-Q/comprx/config"
)

func TestUntar_SingleFile(t *testing.T) {
	tempDir := t.TempDir()

	// 创建测试文件
	testFile := filepath.Join(tempDir, "test.txt")
	testContent := "Hello, Untar World!"
	if err := os.WriteFile(testFile, []byte(testContent), 0644); err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}

	// 压缩文件
	tarFile := filepath.Join(tempDir, "test.tar")
	cfg := config.New()

	if err := Tar(tarFile, testFile, cfg); err != nil {
		t.Fatalf("TAR压缩失败: %v", err)
	}

	// 解压文件
	extractDir := filepath.Join(tempDir, "extract")
	if err := Untar(tarFile, extractDir, cfg); err != nil {
		t.Fatalf("TAR解压失败: %v", err)
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

func TestUntar_Directory(t *testing.T) {
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

	// 解压目录
	extractDir := filepath.Join(tempDir, "extract")
	if err := Untar(tarFile, extractDir, cfg); err != nil {
		t.Fatalf("TAR解压失败: %v", err)
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

func TestUntar_NonExistentFile(t *testing.T) {
	tempDir := t.TempDir()

	nonExistentFile := filepath.Join(tempDir, "nonexistent.tar")
	extractDir := filepath.Join(tempDir, "extract")
	cfg := config.New()

	err := Untar(nonExistentFile, extractDir, cfg)
	if err == nil {
		t.Errorf("期望解压不存在的文件时返回错误")
	}
}

func TestUntar_InvalidTarFile(t *testing.T) {
	tempDir := t.TempDir()

	// 创建无效的TAR文件
	invalidFile := filepath.Join(tempDir, "invalid.tar")
	if err := os.WriteFile(invalidFile, []byte("not a tar file"), 0644); err != nil {
		t.Fatalf("创建无效文件失败: %v", err)
	}

	extractDir := filepath.Join(tempDir, "extract")
	cfg := config.New()

	err := Untar(invalidFile, extractDir, cfg)
	if err == nil {
		t.Errorf("期望解压无效TAR文件时返回错误")
	}
}

func TestUntar_EmptyTar(t *testing.T) {
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
		t.Fatalf("压缩空目录失败: %v", err)
	}

	// 解压空TAR
	extractDir := filepath.Join(tempDir, "extract")
	if err := Untar(tarFile, extractDir, cfg); err != nil {
		t.Fatalf("解压空TAR失败: %v", err)
	}

	// 验证解压目录存在
	if _, err := os.Stat(extractDir); os.IsNotExist(err) {
		t.Errorf("解压目录未创建")
	}
}

func TestUntar_LargeFile(t *testing.T) {
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

	// 解压大文件
	extractDir := filepath.Join(tempDir, "extract")
	if err := Untar(tarFile, extractDir, cfg); err != nil {
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

func TestUntar_CreateTargetDirectory(t *testing.T) {
	tempDir := t.TempDir()

	// 创建测试文件
	testFile := filepath.Join(tempDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}

	// 压缩文件
	tarFile := filepath.Join(tempDir, "test.tar")
	cfg := config.New()

	if err := Tar(tarFile, testFile, cfg); err != nil {
		t.Fatalf("TAR压缩失败: %v", err)
	}

	// 解压到不存在的目录
	extractDir := filepath.Join(tempDir, "nonexistent", "extract")
	if err := Untar(tarFile, extractDir, cfg); err != nil {
		t.Fatalf("解压失败: %v", err)
	}

	// 验证目标目录被创建
	if _, err := os.Stat(extractDir); os.IsNotExist(err) {
		t.Errorf("目标目录未被创建")
	}

	// 验证文件被解压
	extractedFile := filepath.Join(extractDir, "test.txt")
	if _, err := os.Stat(extractedFile); os.IsNotExist(err) {
		t.Errorf("解压的文件不存在")
	}
}

func TestUntar_SymbolicLinks(t *testing.T) {
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

	// 压缩包含符号链接的目录
	testDir := filepath.Join(tempDir, "testdir")
	if err := os.MkdirAll(testDir, 0755); err != nil {
		t.Fatalf("创建测试目录失败: %v", err)
	}

	// 将文件复制到测试目录
	targetFile := filepath.Join(testDir, "target.txt")
	if err := os.WriteFile(targetFile, []byte("target content"), 0644); err != nil {
		t.Fatalf("创建目标文件失败: %v", err)
	}

	// 在测试目录中创建符号链接
	targetLink := filepath.Join(testDir, "link.txt")
	if err := os.Symlink("target.txt", targetLink); err != nil {
		t.Skipf("创建符号链接失败，跳过测试: %v", err)
	}

	cfg := config.New()
	if err := Tar(tarFile, testDir, cfg); err != nil {
		t.Fatalf("压缩符号链接失败: %v", err)
	}

	// 解压
	extractDir := filepath.Join(tempDir, "extract")
	if err := Untar(tarFile, extractDir, cfg); err != nil {
		t.Fatalf("解压符号链接失败: %v", err)
	}

	// 验证符号链接被正确创建
	extractedLink := filepath.Join(extractDir, "testdir", "link.txt")
	if _, err := os.Lstat(extractedLink); err != nil {
		t.Errorf("符号链接未被创建: %v", err)
	}
}

// 基准测试
func BenchmarkUntar_SmallFile(b *testing.B) {
	tempDir := b.TempDir()

	// 创建测试文件
	testFile := filepath.Join(tempDir, "test.txt")
	content := "Small file content for benchmarking."
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		b.Fatalf("创建测试文件失败: %v", err)
	}

	// 压缩文件
	tarFile := filepath.Join(tempDir, "test.tar")
	cfg := config.New()

	if err := Tar(tarFile, testFile, cfg); err != nil {
		b.Fatalf("压缩文件失败: %v", err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		extractDir := filepath.Join(tempDir, "extract_"+string(rune(i)))
		_ = Untar(tarFile, extractDir, cfg)
	}
}

func BenchmarkUntar_LargeFile(b *testing.B) {
	tempDir := b.TempDir()

	// 创建大文件
	testFile := filepath.Join(tempDir, "large.txt")
	largeContent := strings.Repeat("Large file content for benchmarking. ", 5000)
	if err := os.WriteFile(testFile, []byte(largeContent), 0644); err != nil {
		b.Fatalf("创建大文件失败: %v", err)
	}

	// 压缩大文件
	tarFile := filepath.Join(tempDir, "large.tar")
	cfg := config.New()

	if err := Tar(tarFile, testFile, cfg); err != nil {
		b.Fatalf("压缩大文件失败: %v", err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		extractDir := filepath.Join(tempDir, "extract_"+string(rune(i)))
		_ = Untar(tarFile, extractDir, cfg)
	}
}

func TestUntar_OverwriteExisting(t *testing.T) {
	tempDir := t.TempDir()

	// 创建测试文件
	testFile := filepath.Join(tempDir, "test.txt")
	originalContent := "Original content"
	if err := os.WriteFile(testFile, []byte(originalContent), 0644); err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}

	// 压缩文件
	tarFile := filepath.Join(tempDir, "test.tar")
	cfg := config.New()

	if err := Tar(tarFile, testFile, cfg); err != nil {
		t.Fatalf("TAR压缩失败: %v", err)
	}

	// 创建解压目录
	extractDir := filepath.Join(tempDir, "extracted")
	if err := os.MkdirAll(extractDir, 0755); err != nil {
		t.Fatalf("创建解压目录失败: %v", err)
	}

	// 在解压目录中创建同名文件
	existingFile := filepath.Join(extractDir, "test.txt")
	existingContent := "Existing content"
	if err := os.WriteFile(existingFile, []byte(existingContent), 0644); err != nil {
		t.Fatalf("创建已存在文件失败: %v", err)
	}

	t.Run("不覆盖已存在文件", func(t *testing.T) {
		// 设置不允许覆盖
		cfg := config.New()
		cfg.OverwriteExisting = false

		// 尝试解压，应该失败
		err := Untar(tarFile, extractDir, cfg)
		if err == nil {
			t.Errorf("应该返回错误，因为文件已存在且不允许覆盖")
		}

		// 验证原文件内容未被修改
		content, err := os.ReadFile(existingFile)
		if err != nil {
			t.Fatalf("读取文件失败: %v", err)
		}

		if string(content) != existingContent {
			t.Errorf("文件内容被意外修改: got %q, want %q", string(content), existingContent)
		}
	})

	t.Run("覆盖已存在文件", func(t *testing.T) {
		// 重新创建已存在文件
		if err := os.WriteFile(existingFile, []byte(existingContent), 0644); err != nil {
			t.Fatalf("重新创建已存在文件失败: %v", err)
		}

		// 设置允许覆盖
		cfg := config.New()
		cfg.OverwriteExisting = true

		// 解压，应该成功
		if err := Untar(tarFile, extractDir, cfg); err != nil {
			t.Fatalf("解压失败: %v", err)
		}

		// 验证文件内容被覆盖
		content, err := os.ReadFile(existingFile)
		if err != nil {
			t.Fatalf("读取文件失败: %v", err)
		}

		if string(content) != originalContent {
			t.Errorf("文件内容 = %q, want %q", string(content), originalContent)
		}
	})
}

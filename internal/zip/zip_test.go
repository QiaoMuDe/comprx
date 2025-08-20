package zip

import (
	"archive/zip"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gitee.com/MM-Q/comprx/config"
	"gitee.com/MM-Q/comprx/internal/utils"
)

// 创建测试目录和文件的辅助函数
func createTestFiles(t *testing.T, baseDir string) {
	// 创建测试目录结构
	dirs := []string{
		filepath.Join(baseDir, "dir1"),
		filepath.Join(baseDir, "dir2"),
		filepath.Join(baseDir, "dir1", "subdir"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("创建测试目录失败: %v", err)
		}
	}

	// 创建测试文件
	files := map[string]string{
		filepath.Join(baseDir, "file1.txt"):                   "Hello World 1",
		filepath.Join(baseDir, "file2.txt"):                   "Hello World 2",
		filepath.Join(baseDir, "dir1", "file3.txt"):           "Hello World 3",
		filepath.Join(baseDir, "dir1", "subdir", "file4.txt"): "Hello World 4",
		filepath.Join(baseDir, "dir2", "file5.txt"):           "Hello World 5",
	}

	for filePath, content := range files {
		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			t.Fatalf("创建测试文件失败: %v", err)
		}
	}
}

func TestZip_SingleFile(t *testing.T) {
	// 创建临时目录
	tempDir := t.TempDir()

	// 创建测试文件
	testFile := filepath.Join(tempDir, "test.txt")
	testContent := "Hello, World!"
	if err := os.WriteFile(testFile, []byte(testContent), 0644); err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}

	// 创建ZIP文件
	zipFile := filepath.Join(tempDir, "test.zip")
	cfg := config.New()

	err := Zip(zipFile, testFile, cfg)
	if err != nil {
		t.Fatalf("压缩单文件失败: %v", err)
	}

	// 验证ZIP文件是否创建
	if !utils.Exists(zipFile) {
		t.Error("ZIP文件未创建")
	}

	// 验证ZIP文件内容
	reader, err := zip.OpenReader(zipFile)
	if err != nil {
		t.Fatalf("打开ZIP文件失败: %v", err)
	}
	defer func() { _ = reader.Close() }()

	if len(reader.File) != 1 {
		t.Errorf("ZIP文件中的文件数量 = %d, want 1", len(reader.File))
	}

	file := reader.File[0]
	if file.Name != "test.txt" {
		t.Errorf("ZIP文件中的文件名 = %s, want test.txt", file.Name)
	}
}

func TestZip_Directory(t *testing.T) {
	// 创建临时目录
	tempDir := t.TempDir()
	testDir := filepath.Join(tempDir, "testdir")

	// 创建测试文件结构
	createTestFiles(t, testDir)

	// 创建ZIP文件
	zipFile := filepath.Join(tempDir, "test.zip")
	cfg := config.New()

	err := Zip(zipFile, testDir, cfg)
	if err != nil {
		t.Fatalf("压缩目录失败: %v", err)
	}

	// 验证ZIP文件是否创建
	if !utils.Exists(zipFile) {
		t.Error("ZIP文件未创建")
	}

	// 验证ZIP文件内容
	reader, err := zip.OpenReader(zipFile)
	if err != nil {
		t.Fatalf("打开ZIP文件失败: %v", err)
	}
	defer func() { _ = reader.Close() }()

	// 检查文件数量（包括目录）
	if len(reader.File) < 5 {
		t.Errorf("ZIP文件中的条目数量 = %d, want >= 5", len(reader.File))
	}

	// 检查是否包含预期的文件
	expectedFiles := []string{"testdir/file1.txt", "testdir/file2.txt", "testdir/dir1/file3.txt"}
	for _, expected := range expectedFiles {
		found := false
		for _, file := range reader.File {
			if strings.Contains(file.Name, filepath.Base(expected)) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("ZIP文件中未找到预期文件: %s", expected)
		}
	}
}

func TestZip_CompressionLevels(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.txt")

	// 创建一个较大的测试文件以便测试压缩效果
	content := strings.Repeat("Hello, World! This is a test content for compression. ", 1000)
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}

	compressionLevels := []config.CompressionLevel{
		config.CompressionLevelNone,
		config.CompressionLevelFast,
		config.CompressionLevelBest,
		config.CompressionLevelDefault,
	}

	for _, level := range compressionLevels {
		t.Run(fmt.Sprintf("level_%d", int(level)), func(t *testing.T) {
			zipFile := filepath.Join(tempDir, fmt.Sprintf("test_%d.zip", int(level)))
			cfg := config.New()
			cfg.CompressionLevel = level
			cfg.OverwriteExisting = true

			err := Zip(zipFile, testFile, cfg)
			if err != nil {
				t.Fatalf("压缩失败 (level %d): %v", level, err)
			}

			// 验证文件存在
			if !utils.Exists(zipFile) {
				t.Errorf("ZIP文件未创建 (level %d)", level)
			}
		})
	}
}

func TestZip_OverwriteExisting(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.txt")
	zipFile := filepath.Join(tempDir, "test.zip")

	// 创建测试文件
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}

	// 创建已存在的ZIP文件
	if err := os.WriteFile(zipFile, []byte("existing content"), 0644); err != nil {
		t.Fatalf("创建已存在的ZIP文件失败: %v", err)
	}

	t.Run("不覆盖已存在文件", func(t *testing.T) {
		cfg := config.New()
		cfg.OverwriteExisting = false

		err := Zip(zipFile, testFile, cfg)
		if err == nil {
			t.Error("应该返回错误，因为文件已存在且不允许覆盖")
		}
	})

	t.Run("覆盖已存在文件", func(t *testing.T) {
		cfg := config.New()
		cfg.OverwriteExisting = true

		err := Zip(zipFile, testFile, cfg)
		if err != nil {
			t.Errorf("覆盖已存在文件失败: %v", err)
		}
	})
}

func TestZip_NonExistentSource(t *testing.T) {
	tempDir := t.TempDir()
	nonExistentFile := filepath.Join(tempDir, "nonexistent.txt")
	zipFile := filepath.Join(tempDir, "test.zip")
	cfg := config.New()

	err := Zip(zipFile, nonExistentFile, cfg)
	if err == nil {
		t.Error("应该返回错误，因为源文件不存在")
	}
}

func TestZip_InvalidDestination(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.txt")

	// 创建测试文件
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}

	// 使用无效的目标路径（不存在的目录）
	invalidZipFile := filepath.Join(tempDir, "nonexistent", "test.zip")
	cfg := config.New()

	err := Zip(invalidZipFile, testFile, cfg)
	if err != nil {
		// 应该自动创建目录，所以不应该出错
		// 但如果出错，检查是否是预期的错误类型
		if !strings.Contains(err.Error(), "创建目标目录失败") {
			t.Errorf("意外的错误类型: %v", err)
		}
	}
}

func TestZip_EmptyDirectory(t *testing.T) {
	tempDir := t.TempDir()
	emptyDir := filepath.Join(tempDir, "empty")
	zipFile := filepath.Join(tempDir, "empty.zip")

	// 创建空目录
	if err := os.MkdirAll(emptyDir, 0755); err != nil {
		t.Fatalf("创建空目录失败: %v", err)
	}

	cfg := config.New()
	err := Zip(zipFile, emptyDir, cfg)
	if err != nil {
		t.Fatalf("压缩空目录失败: %v", err)
	}

	// 验证ZIP文件内容
	reader, err := zip.OpenReader(zipFile)
	if err != nil {
		t.Fatalf("打开ZIP文件失败: %v", err)
	}
	defer func() { _ = reader.Close() }()

	// 空目录应该至少包含目录本身
	if len(reader.File) == 0 {
		t.Error("空目录的ZIP文件应该包含目录条目")
	}
}

func TestZip_SymbolicLinks(t *testing.T) {
	// 在Windows上跳过符号链接测试
	if os.Getenv("GOOS") == "windows" {
		t.Skip("跳过Windows上的符号链接测试")
	}

	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.txt")
	linkFile := filepath.Join(tempDir, "link.txt")
	zipFile := filepath.Join(tempDir, "test.zip")

	// 创建测试文件
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}

	// 创建符号链接
	if err := os.Symlink(testFile, linkFile); err != nil {
		t.Skipf("创建符号链接失败，跳过测试: %v", err)
	}

	cfg := config.New()
	err := Zip(zipFile, linkFile, cfg)
	if err != nil {
		t.Fatalf("压缩符号链接失败: %v", err)
	}

	// 验证ZIP文件内容
	reader, err := zip.OpenReader(zipFile)
	if err != nil {
		t.Fatalf("打开ZIP文件失败: %v", err)
	}
	defer func() { _ = reader.Close() }()

	if len(reader.File) != 1 {
		t.Errorf("ZIP文件中的文件数量 = %d, want 1", len(reader.File))
	}
}

// 基准测试
func BenchmarkZip_SmallFile(b *testing.B) {
	tempDir := b.TempDir()
	testFile := filepath.Join(tempDir, "test.txt")
	content := "Hello, World!"

	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		b.Fatalf("创建测试文件失败: %v", err)
	}

	cfg := config.New()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		zipFile := filepath.Join(tempDir, "test_"+string(rune(i))+".zip")
		_ = Zip(zipFile, testFile, cfg)
	}
}

func BenchmarkZip_LargeFile(b *testing.B) {
	tempDir := b.TempDir()
	testFile := filepath.Join(tempDir, "large.txt")
	content := strings.Repeat("This is a large file content for benchmarking. ", 10000)

	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		b.Fatalf("创建大文件失败: %v", err)
	}

	cfg := config.New()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		zipFile := filepath.Join(tempDir, "large_"+string(rune(i))+".zip")
		_ = Zip(zipFile, testFile, cfg)
	}
}

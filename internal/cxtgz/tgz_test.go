package cxtgz

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gitee.com/MM-Q/comprx/internal/config"
	"gitee.com/MM-Q/comprx/types"
)

func TestTgz_SingleFile(t *testing.T) {
	tempDir := t.TempDir()

	// 创建测试文件
	testFile := filepath.Join(tempDir, "test.txt")
	testContent := "Hello, TGZ World!"
	if err := os.WriteFile(testFile, []byte(testContent), 0644); err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}

	// 压缩文件
	tgzFile := filepath.Join(tempDir, "test.tgz")
	cfg := config.New()

	if err := Tgz(tgzFile, testFile, cfg); err != nil {
		t.Fatalf("TGZ压缩失败: %v", err)
	}

	// 验证压缩文件存在
	if _, err := os.Stat(tgzFile); os.IsNotExist(err) {
		t.Errorf("TGZ文件未创建")
	}

	// 验证文件大小
	info, err := os.Stat(tgzFile)
	if err != nil {
		t.Fatalf("获取TGZ文件信息失败: %v", err)
	}

	if info.Size() == 0 {
		t.Errorf("TGZ文件大小为0")
	}
}

func TestTgz_Directory(t *testing.T) {
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

	// 验证压缩文件存在
	if _, err := os.Stat(tgzFile); os.IsNotExist(err) {
		t.Errorf("TGZ文件未创建")
	}
}

func TestTgz_EmptyDirectory(t *testing.T) {
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
		t.Fatalf("TGZ压缩空目录失败: %v", err)
	}

	// 验证压缩文件存在
	if _, err := os.Stat(tgzFile); os.IsNotExist(err) {
		t.Errorf("TGZ文件未创建")
	}
}

func TestTgz_NonExistentSource(t *testing.T) {
	tempDir := t.TempDir()

	nonExistentPath := filepath.Join(tempDir, "nonexistent")
	tgzFile := filepath.Join(tempDir, "test.tgz")
	cfg := config.New()

	err := Tgz(tgzFile, nonExistentPath, cfg)
	if err == nil {
		t.Errorf("期望压缩不存在的文件时返回错误")
	}
}

func TestTgz_InvalidDestination(t *testing.T) {
	tempDir := t.TempDir()

	// 创建测试文件
	testFile := filepath.Join(tempDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}

	// 使用无效的目标路径（空字符串）
	invalidPath := ""
	cfg := config.New()

	err := Tgz(invalidPath, testFile, cfg)
	if err == nil {
		t.Errorf("期望使用无效目标路径时返回错误")
	}
}

func TestTgz_OverwriteExisting(t *testing.T) {
	tempDir := t.TempDir()

	// 创建测试文件
	testFile := filepath.Join(tempDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}

	tgzFile := filepath.Join(tempDir, "test.tgz")

	// 创建已存在的目标文件
	if err := os.WriteFile(tgzFile, []byte("existing"), 0644); err != nil {
		t.Fatalf("创建已存在文件失败: %v", err)
	}

	// 测试不覆盖
	cfg := config.New()
	cfg.OverwriteExisting = false

	err := Tgz(tgzFile, testFile, cfg)
	if err == nil {
		t.Errorf("期望不覆盖已存在文件时返回错误")
	}

	// 测试覆盖
	cfg.OverwriteExisting = true

	if err := Tgz(tgzFile, testFile, cfg); err != nil {
		t.Fatalf("覆盖已存在文件失败: %v", err)
	}
}

func TestTgz_CompressionLevels(t *testing.T) {
	tempDir := t.TempDir()

	// 创建测试文件
	testFile := filepath.Join(tempDir, "test.txt")
	largeContent := strings.Repeat("This is test content for compression. ", 1000)
	if err := os.WriteFile(testFile, []byte(largeContent), 0644); err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}

	levels := []types.CompressionLevel{
		types.CompressionLevelNone,
		types.CompressionLevelFast,
		types.CompressionLevelDefault,
		types.CompressionLevelBest,
	}

	for _, level := range levels {
		t.Run(fmt.Sprintf("level_%d", int(level)), func(t *testing.T) {
			tgzFile := filepath.Join(tempDir, fmt.Sprintf("test_%d.tgz", int(level)))
			cfg := config.New()
			cfg.CompressionLevel = level

			if err := Tgz(tgzFile, testFile, cfg); err != nil {
				t.Fatalf("TGZ压缩失败 (level %d): %v", level, err)
			}

			// 验证文件存在
			if _, err := os.Stat(tgzFile); os.IsNotExist(err) {
				t.Errorf("TGZ文件未创建 (level %d)", level)
			}
		})
	}
}

// 基准测试
func BenchmarkTgz_SmallFile(b *testing.B) {
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
		tgzFile := filepath.Join(tempDir, "bench_"+string(rune(i))+".tgz")
		_ = Tgz(tgzFile, testFile, cfg)
	}
}

func BenchmarkTgz_LargeFile(b *testing.B) {
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
		tgzFile := filepath.Join(tempDir, "bench_large_"+string(rune(i))+".tgz")
		_ = Tgz(tgzFile, testFile, cfg)
	}
}

func BenchmarkTgz_Directory(b *testing.B) {
	tempDir := b.TempDir()

	// 创建测试目录
	testDir := filepath.Join(tempDir, "testdir")
	if err := os.MkdirAll(testDir, 0755); err != nil {
		b.Fatalf("创建测试目录失败: %v", err)
	}

	// 创建多个文件
	for i := 0; i < 10; i++ {
		fileName := filepath.Join(testDir, "file_"+string(rune(i))+".txt")
		content := "File content " + string(rune(i))
		if err := os.WriteFile(fileName, []byte(content), 0644); err != nil {
			b.Fatalf("创建文件失败: %v", err)
		}
	}

	cfg := config.New()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		tgzFile := filepath.Join(tempDir, "bench_dir_"+string(rune(i))+".tgz")
		_ = Tgz(tgzFile, testDir, cfg)
	}
}

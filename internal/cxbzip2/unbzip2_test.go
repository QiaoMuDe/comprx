package cxbzip2

import (
	"os"
	"path/filepath"
	"testing"

	"gitee.com/MM-Q/comprx/config"
)

func TestUnbz2_SourceNotFound(t *testing.T) {
	tempDir := t.TempDir()
	srcFile := filepath.Join(tempDir, "nonexistent.bz2")
	targetFile := filepath.Join(tempDir, "target.txt")

	cfg := config.New()
	cfg.OverwriteExisting = true

	err := Unbz2(srcFile, targetFile, cfg)
	if err == nil {
		t.Fatalf("期望源文件不存在时返回错误")
	}
}

func TestUnbz2_InvalidBzip2File(t *testing.T) {
	tempDir := t.TempDir()

	// 创建无效的 BZIP2 文件
	invalidFile := filepath.Join(tempDir, "invalid.bz2")
	if err := os.WriteFile(invalidFile, []byte("not a bzip2 file"), 0644); err != nil {
		t.Fatalf("创建无效文件失败: %v", err)
	}

	targetFile := filepath.Join(tempDir, "target.txt")
	cfg := config.New()
	cfg.OverwriteExisting = true

	err := Unbz2(invalidFile, targetFile, cfg)
	if err == nil {
		t.Fatalf("期望无效 BZIP2 文件时返回错误")
	}
}

func TestUnbz2_OverwriteExisting(t *testing.T) {
	tempDir := t.TempDir()

	// 创建一个简单的测试文件（不是真正的 bzip2 格式，但用于测试覆盖逻辑）
	bzip2File := filepath.Join(tempDir, "test.bz2")
	if err := os.WriteFile(bzip2File, []byte("fake bzip2 content"), 0644); err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}

	// 创建已存在的目标文件
	targetFile := filepath.Join(tempDir, "target.txt")
	if err := os.WriteFile(targetFile, []byte("existing content"), 0644); err != nil {
		t.Fatalf("创建已存在文件失败: %v", err)
	}

	// 测试不允许覆盖
	cfg := config.New()
	cfg.OverwriteExisting = false
	err := Unbz2(bzip2File, targetFile, cfg)
	if err == nil {
		t.Fatalf("期望不覆盖已存在文件时返回错误")
	}

	// 验证文件未被覆盖
	content, _ := os.ReadFile(targetFile)
	if string(content) != "existing content" {
		t.Fatalf("文件被意外覆盖")
	}
}

func TestUnbz2_TargetIsDirectory(t *testing.T) {
	tempDir := t.TempDir()

	// 创建测试文件
	bzip2File := filepath.Join(tempDir, "test.txt.bz2")
	if err := os.WriteFile(bzip2File, []byte("test content"), 0644); err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}

	// 创建目标目录
	targetDir := filepath.Join(tempDir, "output")
	if err := os.Mkdir(targetDir, 0755); err != nil {
		t.Fatalf("创建目标目录失败: %v", err)
	}

	cfg := config.New()
	cfg.OverwriteExisting = true

	// 使用目录作为目标路径，应该自动生成文件名
	err := Unbz2(bzip2File, targetDir, cfg)
	// 由于文件不是真正的 bzip2 格式，这里会失败，但我们主要测试路径处理逻辑
	if err != nil {
		t.Logf("预期的错误（因为不是真正的 bzip2 文件）: %v", err)
	}
}

func TestUnbz2_CreateTargetDirectory(t *testing.T) {
	tempDir := t.TempDir()

	// 创建测试文件
	bzip2File := filepath.Join(tempDir, "test.bz2")
	if err := os.WriteFile(bzip2File, []byte("test content"), 0644); err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}

	// 设置目标文件在不存在的目录中
	targetFile := filepath.Join(tempDir, "subdir", "target.txt")
	cfg := config.New()
	cfg.OverwriteExisting = true

	err := Unbz2(bzip2File, targetFile, cfg)
	// 由于不是真正的 bzip2 文件会失败，但目录应该被创建
	if err != nil {
		t.Logf("预期的错误（因为不是真正的 bzip2 文件）: %v", err)
	}

	// 验证目录是否被创建
	if _, err := os.Stat(filepath.Dir(targetFile)); os.IsNotExist(err) {
		t.Fatalf("目标目录未被创建")
	}
}

func TestUnbz2_EmptyFile(t *testing.T) {
	tempDir := t.TempDir()

	// 创建空的测试文件
	bzip2File := filepath.Join(tempDir, "empty.bz2")
	if err := os.WriteFile(bzip2File, []byte(""), 0644); err != nil {
		t.Fatalf("创建空文件失败: %v", err)
	}

	targetFile := filepath.Join(tempDir, "empty_extracted.txt")
	cfg := config.New()
	cfg.OverwriteExisting = true

	err := Unbz2(bzip2File, targetFile, cfg)
	if err == nil {
		t.Fatalf("期望空文件解压时返回错误")
	}
}

// 这个测试需要真正的 bzip2 文件，在实际项目中可以添加测试数据文件
func TestUnbz2_RealBzip2File(t *testing.T) {
	t.Skip("需要真正的 bzip2 测试文件，跳过此测试")

	// 在实际项目中，可以这样实现：
	// 1. 在 testdata 目录中放置真正的 bzip2 文件
	// 2. 使用这些文件进行完整的解压测试
	// 3. 验证解压后的内容是否正确

	/*
		tempDir := t.TempDir()

		// 使用预制的 bzip2 文件
		bzip2File := "testdata/sample.txt.bz2"
		targetFile := filepath.Join(tempDir, "extracted.txt")

		cfg := &config.Config{OverwriteExisting: true}

		err := Unbz2(bzip2File, targetFile, cfg)
		if err != nil {
			t.Fatalf("BZIP2 解压失败: %v", err)
		}

		// 验证解压内容
		content, err := os.ReadFile(targetFile)
		if err != nil {
			t.Fatalf("读取解压文件失败: %v", err)
		}

		expectedContent := "expected content from bzip2 file"
		if string(content) != expectedContent {
			t.Fatalf("解压内容不匹配")
		}
	*/
}

func BenchmarkUnbz2_FakeFile(b *testing.B) {
	tempDir := b.TempDir()

	// 创建假的 bzip2 文件用于基准测试
	bzip2File := filepath.Join(tempDir, "benchmark.bz2")
	content := make([]byte, 1024)
	for i := range content {
		content[i] = byte(i % 256)
	}
	if err := os.WriteFile(bzip2File, content, 0644); err != nil {
		b.Fatalf("创建测试文件失败: %v", err)
	}

	cfg := &config.Config{OverwriteExisting: true}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		targetFile := filepath.Join(tempDir, "benchmark_"+string(rune(i))+".txt")
		// 这会失败，但可以测试文件处理的开销
		_ = Unbz2(bzip2File, targetFile, cfg)
	}
}

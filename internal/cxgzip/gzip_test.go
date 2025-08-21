package cxgzip

import (
	"os"
	"path/filepath"
	"testing"

	"gitee.com/MM-Q/comprx/config"
	"gitee.com/MM-Q/comprx/types"
)

func TestGzip_Success(t *testing.T) {
	// 创建临时目录
	tempDir := t.TempDir()

	// 创建测试文件
	srcFile := filepath.Join(tempDir, "test.txt")
	content := "Hello, GZIP compression test!"
	if err := os.WriteFile(srcFile, []byte(content), 0644); err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}

	// 设置目标文件
	dstFile := filepath.Join(tempDir, "test.txt.gz")

	// 创建配置
	cfg := config.New()
	cfg.OverwriteExisting = true
	cfg.CompressionLevel = types.CompressionLevelDefault

	// 执行压缩
	err := Gzip(dstFile, srcFile, cfg)
	if err != nil {
		t.Fatalf("GZIP 压缩失败: %v", err)
	}

	// 验证压缩文件是否存在
	if _, err := os.Stat(dstFile); os.IsNotExist(err) {
		t.Fatalf("压缩文件未创建: %s", dstFile)
	}

	// 验证压缩文件大小
	info, err := os.Stat(dstFile)
	if err != nil {
		t.Fatalf("获取压缩文件信息失败: %v", err)
	}
	if info.Size() == 0 {
		t.Fatalf("压缩文件为空")
	}
}

func TestGzip_SourceNotFound(t *testing.T) {
	tempDir := t.TempDir()
	srcFile := filepath.Join(tempDir, "nonexistent.txt")
	dstFile := filepath.Join(tempDir, "test.txt.gz")

	cfg := config.New()
	cfg.OverwriteExisting = true

	err := Gzip(dstFile, srcFile, cfg)
	if err == nil {
		t.Fatalf("期望源文件不存在时返回错误")
	}
}

func TestGzip_SourceIsDirectory(t *testing.T) {
	tempDir := t.TempDir()
	srcDir := filepath.Join(tempDir, "testdir")
	if err := os.Mkdir(srcDir, 0755); err != nil {
		t.Fatalf("创建测试目录失败: %v", err)
	}

	dstFile := filepath.Join(tempDir, "test.txt.gz")
	cfg := config.New()
	cfg.OverwriteExisting = true

	err := Gzip(dstFile, srcDir, cfg)
	if err == nil {
		t.Fatalf("期望源为目录时返回错误")
	}
	if err.Error() != "GZIP 只支持单文件压缩，不支持目录压缩" {
		t.Fatalf("错误信息不匹配，得到: %v", err)
	}
}

func TestGzip_OverwriteExisting(t *testing.T) {
	tempDir := t.TempDir()

	// 创建源文件
	srcFile := filepath.Join(tempDir, "test.txt")
	if err := os.WriteFile(srcFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("创建源文件失败: %v", err)
	}

	// 创建已存在的目标文件
	dstFile := filepath.Join(tempDir, "test.txt.gz")
	if err := os.WriteFile(dstFile, []byte("existing"), 0644); err != nil {
		t.Fatalf("创建已存在文件失败: %v", err)
	}

	// 测试不允许覆盖
	cfg := config.New()
	cfg.OverwriteExisting = false
	err := Gzip(dstFile, srcFile, cfg)
	if err == nil {
		t.Fatalf("期望不覆盖已存在文件时返回错误")
	}

	// 测试允许覆盖
	cfg.OverwriteExisting = true
	err = Gzip(dstFile, srcFile, cfg)
	if err != nil {
		t.Fatalf("允许覆盖时不应返回错误: %v", err)
	}
}

func TestGzip_EmptyFile(t *testing.T) {
	tempDir := t.TempDir()

	// 创建空文件
	srcFile := filepath.Join(tempDir, "empty.txt")
	if err := os.WriteFile(srcFile, []byte(""), 0644); err != nil {
		t.Fatalf("创建空文件失败: %v", err)
	}

	dstFile := filepath.Join(tempDir, "empty.txt.gz")
	cfg := config.New()
	cfg.OverwriteExisting = true

	err := Gzip(dstFile, srcFile, cfg)
	if err != nil {
		t.Fatalf("压缩空文件失败: %v", err)
	}

	// 验证压缩文件存在
	if _, err := os.Stat(dstFile); os.IsNotExist(err) {
		t.Fatalf("压缩文件未创建")
	}
}

func TestGzip_LargeFile(t *testing.T) {
	tempDir := t.TempDir()

	// 创建较大的测试文件 (1MB)
	srcFile := filepath.Join(tempDir, "large.txt")
	content := make([]byte, 1024*1024)
	for i := range content {
		content[i] = byte(i % 256)
	}
	if err := os.WriteFile(srcFile, content, 0644); err != nil {
		t.Fatalf("创建大文件失败: %v", err)
	}

	dstFile := filepath.Join(tempDir, "large.txt.gz")
	cfg := config.New()
	cfg.OverwriteExisting = true
	cfg.CompressionLevel = 9 // 最高压缩级别

	err := Gzip(dstFile, srcFile, cfg)
	if err != nil {
		t.Fatalf("压缩大文件失败: %v", err)
	}

	// 验证压缩效果
	srcInfo, _ := os.Stat(srcFile)
	dstInfo, _ := os.Stat(dstFile)

	if dstInfo.Size() >= srcInfo.Size() {
		t.Logf("警告: 压缩后文件大小 (%d) >= 原文件大小 (%d)", dstInfo.Size(), srcInfo.Size())
	}
}

func BenchmarkGzip_SmallFile(b *testing.B) {
	tempDir := b.TempDir()
	srcFile := filepath.Join(tempDir, "small.txt")
	content := "Small file content for benchmarking"
	if err := os.WriteFile(srcFile, []byte(content), 0644); err != nil {
		b.Fatalf("创建测试文件失败: %v", err)
	}

	cfg := config.New()
	cfg.OverwriteExisting = true

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dstFile := filepath.Join(tempDir, "small_"+string(rune(i))+".txt.gz")
		if err := Gzip(dstFile, srcFile, cfg); err != nil {
			b.Fatalf("压缩失败: %v", err)
		}
	}
}

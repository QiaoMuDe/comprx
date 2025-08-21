package cxgzip

import (
	"compress/gzip"
	"os"
	"path/filepath"
	"testing"

	"gitee.com/MM-Q/comprx/config"
)

func TestUngzip_Success(t *testing.T) {
	tempDir := t.TempDir()

	// 创建测试内容
	originalContent := "Hello, GZIP decompression test!"

	// 先创建一个 GZIP 文件
	gzipFile := filepath.Join(tempDir, "test.txt.gz")
	file, err := os.Create(gzipFile)
	if err != nil {
		t.Fatalf("创建 GZIP 文件失败: %v", err)
	}
	defer func() {
		_ = file.Close()
		_ = os.Remove(gzipFile)
	}()

	writer := gzip.NewWriter(file)
	writer.Name = "test.txt"
	if _, err := writer.Write([]byte(originalContent)); err != nil {
		t.Fatalf("写入 GZIP 内容失败: %v", err)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("关闭 GZIP writer 失败: %v", err)
	}
	if err := file.Close(); err != nil {
		t.Fatalf("关闭文件失败: %v", err)
	}

	// 设置解压目标
	targetFile := filepath.Join(tempDir, "extracted.txt")
	cfg := config.New()
	cfg.OverwriteExisting = true

	// 执行解压
	err = Ungzip(gzipFile, targetFile, cfg)
	if err != nil {
		t.Fatalf("GZIP 解压失败: %v", err)
	}

	// 验证解压文件内容
	extractedContent, err := os.ReadFile(targetFile)
	if err != nil {
		t.Fatalf("读取解压文件失败: %v", err)
	}

	if string(extractedContent) != originalContent {
		t.Fatalf("解压内容不匹配，期望: %s, 得到: %s", originalContent, string(extractedContent))
	}
}

func TestUngzip_SourceNotFound(t *testing.T) {
	tempDir := t.TempDir()
	srcFile := filepath.Join(tempDir, "nonexistent.gz")
	targetFile := filepath.Join(tempDir, "target.txt")

	cfg := config.New()
	cfg.OverwriteExisting = true

	err := Ungzip(srcFile, targetFile, cfg)
	if err == nil {
		t.Fatalf("期望源文件不存在时返回错误")
	}
}

func TestUngzip_InvalidGzipFile(t *testing.T) {
	tempDir := t.TempDir()

	// 创建无效的 GZIP 文件
	invalidFile := filepath.Join(tempDir, "invalid.gz")
	if err := os.WriteFile(invalidFile, []byte("not a gzip file"), 0644); err != nil {
		t.Fatalf("创建无效文件失败: %v", err)
	}

	targetFile := filepath.Join(tempDir, "target.txt")
	cfg := config.New()
	cfg.OverwriteExisting = true

	err := Ungzip(invalidFile, targetFile, cfg)
	if err == nil {
		t.Fatalf("期望无效 GZIP 文件时返回错误")
	}
}

func TestUngzip_OverwriteExisting(t *testing.T) {
	tempDir := t.TempDir()

	// 创建有效的 GZIP 文件
	originalContent := "test content for overwrite"
	gzipFile := filepath.Join(tempDir, "test.txt.gz")
	file, err := os.Create(gzipFile)
	if err != nil {
		t.Fatalf("创建 GZIP 文件失败: %v", err)
	}

	writer := gzip.NewWriter(file)
	writer.Name = "test.txt"
	if _, err := writer.Write([]byte(originalContent)); err != nil {
		t.Fatalf("写入 GZIP 内容失败: %v", err)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("关闭 GZIP writer 失败: %v", err)
	}
	if err := file.Close(); err != nil {
		t.Fatalf("关闭文件失败: %v", err)
	}

	// 创建已存在的目标文件
	targetFile := filepath.Join(tempDir, "target.txt")
	if err := os.WriteFile(targetFile, []byte("existing content"), 0644); err != nil {
		t.Fatalf("创建已存在文件失败: %v", err)
	}

	// 测试不允许覆盖
	cfg := config.New()
	cfg.OverwriteExisting = false
	err = Ungzip(gzipFile, targetFile, cfg)
	if err == nil {
		t.Fatalf("期望不覆盖已存在文件时返回错误")
	}

	// 验证文件未被覆盖
	content, _ := os.ReadFile(targetFile)
	if string(content) != "existing content" {
		t.Fatalf("文件被意外覆盖")
	}

	// 测试允许覆盖
	cfg.OverwriteExisting = true
	err = Ungzip(gzipFile, targetFile, cfg)
	if err != nil {
		t.Fatalf("允许覆盖时不应返回错误: %v", err)
	}

	// 验证文件被正确覆盖
	content, _ = os.ReadFile(targetFile)
	if string(content) != originalContent {
		t.Fatalf("文件覆盖后内容不正确")
	}
}

func TestUngzip_EmptyGzipFile(t *testing.T) {
	tempDir := t.TempDir()

	// 创建空内容的 GZIP 文件
	gzipFile := filepath.Join(tempDir, "empty.txt.gz")
	file, err := os.Create(gzipFile)
	if err != nil {
		t.Fatalf("创建 GZIP 文件失败: %v", err)
	}

	writer := gzip.NewWriter(file)
	writer.Name = "empty.txt"
	if err := writer.Close(); err != nil {
		t.Fatalf("关闭 GZIP writer 失败: %v", err)
	}
	if err := file.Close(); err != nil {
		t.Fatalf("关闭文件失败: %v", err)
	}

	targetFile := filepath.Join(tempDir, "empty_extracted.txt")
	cfg := config.New()
	cfg.OverwriteExisting = true

	err = Ungzip(gzipFile, targetFile, cfg)
	if err != nil {
		t.Fatalf("解压空 GZIP 文件失败: %v", err)
	}

	// 验证解压后的文件为空
	content, err := os.ReadFile(targetFile)
	if err != nil {
		t.Fatalf("读取解压文件失败: %v", err)
	}
	if len(content) != 0 {
		t.Fatalf("期望空文件，但得到内容: %s", string(content))
	}
}

func TestUngzip_TargetIsDirectory(t *testing.T) {
	tempDir := t.TempDir()

	// 创建 GZIP 文件
	originalContent := "target is directory test"
	gzipFile := filepath.Join(tempDir, "test.txt.gz")
	file, err := os.Create(gzipFile)
	if err != nil {
		t.Fatalf("创建 GZIP 文件失败: %v", err)
	}
	defer func() {
		_ = file.Close()
		_ = os.Remove(gzipFile)
	}()

	writer := gzip.NewWriter(file)
	writer.Name = "test.txt"
	if _, err := writer.Write([]byte(originalContent)); err != nil {
		t.Fatalf("写入 GZIP 内容失败: %v", err)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("关闭 GZIP writer 失败: %v", err)
	}
	if err := file.Close(); err != nil {
		t.Fatalf("关闭文件失败: %v", err)
	}

	// 创建目标目录
	targetDir := filepath.Join(tempDir, "output")
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		t.Fatalf("创建目标目录失败: %v", err)
	}

	// 使用目录作为目标路径，应该自动生成文件名
	cfg := config.New()
	cfg.OverwriteExisting = true
	err = Ungzip(gzipFile, targetDir, cfg)
	if err != nil {
		t.Fatalf("目标为目录时解压失败: %v", err)
	}

	// 验证自动生成的文件
	expectedTarget := filepath.Join(targetDir, "test.txt")
	content, err := os.ReadFile(expectedTarget)
	if err != nil {
		t.Fatalf("读取自动生成的文件失败: %v", err)
	}

	if string(content) != originalContent {
		t.Fatalf("自动生成文件内容不匹配")
	}

	// 清理生成的文件
	_ = os.Remove(expectedTarget)
}

func BenchmarkUngzip_SmallFile(b *testing.B) {
	tempDir := b.TempDir()

	// 创建小的 GZIP 文件
	content := "Small file content for benchmarking ungzip"
	gzipFile := filepath.Join(tempDir, "small.txt.gz")
	file, _ := os.Create(gzipFile)
	writer := gzip.NewWriter(file)
	writer.Name = "small.txt"
	if _, err := writer.Write([]byte(content)); err != nil {
		b.Fatalf("写入 GZIP 内容失败: %v", err)
	}
	if err := writer.Close(); err != nil {
		b.Fatalf("关闭 GZIP writer 失败: %v", err)
	}
	if err := file.Close(); err != nil {
		b.Fatalf("关闭文件失败: %v", err)
	}

	cfg := config.New()
	cfg.OverwriteExisting = true

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		targetFile := filepath.Join(tempDir, "small_"+string(rune(i))+".txt")
		if err := Ungzip(gzipFile, targetFile, cfg); err != nil {
			b.Fatalf("解压失败: %v", err)
		}
	}
}

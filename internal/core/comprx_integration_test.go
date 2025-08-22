package core

import (
	"os"
	"path/filepath"
	"testing"

	"gitee.com/MM-Q/comprx/types"
)

// TestPackUnpackIntegration 集成测试：压缩后解压验证
func TestPackUnpackIntegration(t *testing.T) {
	tempDir := t.TempDir()

	// 创建测试目录结构
	srcDir := filepath.Join(tempDir, "source")
	if err := os.MkdirAll(srcDir, 0755); err != nil {
		t.Fatal(err)
	}

	// 创建测试文件
	testFiles := map[string]string{
		"file1.txt":        "这是第一个测试文件的内容",
		"file2.txt":        "这是第二个测试文件的内容",
		"subdir/file3.txt": "这是子目录中的文件内容",
	}

	for relPath, content := range testFiles {
		fullPath := filepath.Join(srcDir, relPath)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	c := New()
	c.Config.OverwriteExisting = true

	// 测试不同压缩格式的完整流程
	formats := []string{"zip", "tar", "tgz"}

	for _, format := range formats {
		t.Run(format, func(t *testing.T) {
			// 压缩
			archivePath := filepath.Join(tempDir, "test."+format)
			err := c.Pack(archivePath, srcDir)
			if err != nil {
				t.Fatalf("压缩失败: %v", err)
			}

			// 验证压缩文件存在
			if _, statErr := os.Stat(archivePath); os.IsNotExist(statErr) {
				t.Fatalf("压缩文件未创建 %s", statErr)
			}

			// 解压到新目录
			extractDir := filepath.Join(tempDir, "extracted_"+format)
			err = c.Unpack(archivePath, extractDir)
			if err != nil {
				t.Fatalf("解压失败: %v", err)
			}

			// 验证解压后的文件
			for relPath, expectedContent := range testFiles {
				extractedFile := filepath.Join(extractDir, "source", relPath)
				content, err := os.ReadFile(extractedFile)
				if err != nil {
					t.Errorf("读取解压文件失败 %s: %v", relPath, err)
					continue
				}
				if string(content) != expectedContent {
					t.Errorf("文件内容不匹配 %s: 期望 %q, 实际 %q", relPath, expectedContent, string(content))
				}
			}
		})
	}
}

// TestSingleFilePackUnpack 测试单文件压缩解压
func TestSingleFilePackUnpack(t *testing.T) {
	tempDir := t.TempDir()

	// 创建单个测试文件
	srcFile := filepath.Join(tempDir, "single.txt")
	testContent := "这是单个文件的测试内容\n包含多行\n和特殊字符: !@#$%^&*()"
	if err := os.WriteFile(srcFile, []byte(testContent), 0644); err != nil {
		t.Fatal(err)
	}

	c := New()
	c.Config.OverwriteExisting = true

	// 测试 gzip 格式（专门用于单文件）
	t.Run("gzip", func(t *testing.T) {
		gzFile := filepath.Join(tempDir, "single.txt.gz")

		// 压缩
		err := c.Pack(gzFile, srcFile)
		if err != nil {
			t.Fatalf("gzip压缩失败: %v", err)
		}

		// 解压
		extractDir := filepath.Join(tempDir, "gzip_extract")
		err = c.Unpack(gzFile, extractDir)
		if err != nil {
			t.Fatalf("gzip解压失败: %v", err)
		}

		// 验证解压后的文件
		extractedFile := filepath.Join(extractDir, "single.txt")
		content, err := os.ReadFile(extractedFile)
		if err != nil {
			t.Fatalf("读取解压文件失败: %v", err)
		}
		if string(content) != testContent {
			t.Errorf("文件内容不匹配: 期望 %q, 实际 %q", testContent, string(content))
		}
	})
}

// TestCompressionLevels 测试不同压缩级别
func TestCompressionLevels(t *testing.T) {
	tempDir := t.TempDir()

	// 创建较大的测试文件以便观察压缩效果
	srcFile := filepath.Join(tempDir, "large.txt")
	content := make([]byte, 10240) // 10KB
	for i := range content {
		content[i] = byte('A' + (i % 26)) // 重复字符模式，便于压缩
	}
	if err := os.WriteFile(srcFile, content, 0644); err != nil {
		t.Fatal(err)
	}

	levels := []types.CompressionLevel{
		types.CompressionLevelNone,
		types.CompressionLevelFast,
		types.CompressionLevelDefault,
		types.CompressionLevelBest,
	}

	levelNames := map[types.CompressionLevel]string{
		types.CompressionLevelNone:    "None",
		types.CompressionLevelFast:    "Fast",
		types.CompressionLevelDefault: "Default",
		types.CompressionLevelBest:    "Best",
	}

	for _, level := range levels {
		levelName := levelNames[level]
		t.Run(levelName, func(t *testing.T) {
			c := New()
			c.Config.OverwriteExisting = true
			c.Config.CompressionLevel = level

			zipFile := filepath.Join(tempDir, "test_"+levelName+".zip")
			err := c.Pack(zipFile, srcFile)
			if err != nil {
				t.Errorf("压缩级别 %v 失败: %v", level, err)
				return
			}

			// 验证文件存在
			if _, err := os.Stat(zipFile); os.IsNotExist(err) {
				t.Errorf("压缩文件未创建: %s", zipFile)
			}
		})
	}
}

// TestErrorHandling 测试错误处理场景
func TestErrorHandling(t *testing.T) {
	tempDir := t.TempDir()
	c := New()

	// 测试压缩不存在的文件
	t.Run("压缩不存在的文件", func(t *testing.T) {
		nonExistent := filepath.Join(tempDir, "nonexistent.txt")
		zipFile := filepath.Join(tempDir, "test.zip")

		err := c.Pack(zipFile, nonExistent)
		if err == nil {
			t.Error("期望返回错误，但没有错误")
		}
	})

	// 测试解压不存在的文件
	t.Run("解压不存在的文件", func(t *testing.T) {
		nonExistent := filepath.Join(tempDir, "nonexistent.zip")
		extractDir := filepath.Join(tempDir, "extract")

		err := c.Unpack(nonExistent, extractDir)
		if err == nil {
			t.Error("期望返回错误，但没有错误")
		}
	})

	// 测试无效的压缩格式
	t.Run("无效的压缩格式", func(t *testing.T) {
		srcFile := filepath.Join(tempDir, "test.txt")
		if err := os.WriteFile(srcFile, []byte("test"), 0644); err != nil {
			t.Fatal(err)
		}

		invalidFile := filepath.Join(tempDir, "test.invalid")
		err := c.Pack(invalidFile, srcFile)
		if err == nil {
			t.Error("期望返回错误，但没有错误")
		}
	})
}

// TestConcurrentAccess 测试并发访问
func TestConcurrentAccess(t *testing.T) {
	tempDir := t.TempDir()

	// 创建测试文件
	srcFile := filepath.Join(tempDir, "concurrent.txt")
	if err := os.WriteFile(srcFile, []byte("concurrent test content"), 0644); err != nil {
		t.Fatal(err)
	}

	// 创建多个压缩器实例并发执行
	const numGoroutines = 5
	done := make(chan error, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			c := New()
			c.Config.OverwriteExisting = true
			zipFile := filepath.Join(tempDir, "concurrent_"+string(rune('0'+id))+".zip")
			done <- c.Pack(zipFile, srcFile)
		}(i)
	}

	// 等待所有goroutine完成
	for i := 0; i < numGoroutines; i++ {
		if err := <-done; err != nil {
			t.Errorf("并发压缩失败: %v", err)
		}
	}
}

// TestLargeFile 测试大文件处理
func TestLargeFile(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过大文件测试（使用 -short 标志）")
	}

	tempDir := t.TempDir()

	// 创建1MB的测试文件
	srcFile := filepath.Join(tempDir, "large.txt")
	content := make([]byte, 1024*1024) // 1MB
	for i := range content {
		content[i] = byte(i % 256)
	}
	if err := os.WriteFile(srcFile, content, 0644); err != nil {
		t.Fatal(err)
	}

	c := New()
	c.Config.OverwriteExisting = true

	// 测试压缩大文件
	zipFile := filepath.Join(tempDir, "large.zip")
	err := c.Pack(zipFile, srcFile)
	if err != nil {
		t.Fatalf("压缩大文件失败: %v", err)
	}

	// 测试解压大文件
	extractDir := filepath.Join(tempDir, "large_extract")
	err = c.Unpack(zipFile, extractDir)
	if err != nil {
		t.Fatalf("解压大文件失败: %v", err)
	}

	// 验证文件完整性
	extractedFile := filepath.Join(extractDir, "large.txt")
	extractedContent, err := os.ReadFile(extractedFile)
	if err != nil {
		t.Fatalf("读取解压文件失败: %v", err)
	}

	if len(extractedContent) != len(content) {
		t.Errorf("文件大小不匹配: 期望 %d, 实际 %d", len(content), len(extractedContent))
	}

	// 检查前100字节和后100字节
	if len(extractedContent) >= 100 {
		for i := 0; i < 100; i++ {
			if extractedContent[i] != content[i] {
				t.Errorf("文件开头内容不匹配，位置 %d: 期望 %d, 实际 %d", i, content[i], extractedContent[i])
				break
			}
		}

		start := len(content) - 100
		for i := 0; i < 100; i++ {
			if extractedContent[start+i] != content[start+i] {
				t.Errorf("文件结尾内容不匹配，位置 %d: 期望 %d, 实际 %d", start+i, content[start+i], extractedContent[start+i])
				break
			}
		}
	}
}

// TestEmptyFiles 测试空文件处理
func TestEmptyFiles(t *testing.T) {
	tempDir := t.TempDir()

	// 创建空文件
	emptyFile := filepath.Join(tempDir, "empty.txt")
	if err := os.WriteFile(emptyFile, []byte{}, 0644); err != nil {
		t.Fatal(err)
	}

	c := New()
	c.Config.OverwriteExisting = true

	// 测试压缩空文件
	zipFile := filepath.Join(tempDir, "empty.zip")
	err := c.Pack(zipFile, emptyFile)
	if err != nil {
		t.Fatalf("压缩空文件失败: %v", err)
	}

	// 测试解压空文件
	extractDir := filepath.Join(tempDir, "empty_extract")
	err = c.Unpack(zipFile, extractDir)
	if err != nil {
		t.Fatalf("解压空文件失败: %v", err)
	}

	// 验证空文件
	extractedFile := filepath.Join(extractDir, "empty.txt")
	content, err := os.ReadFile(extractedFile)
	if err != nil {
		t.Fatalf("读取解压的空文件失败: %v", err)
	}
	if len(content) != 0 {
		t.Errorf("空文件应该为空，但包含 %d 字节", len(content))
	}
}

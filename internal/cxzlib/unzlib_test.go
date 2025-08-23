package cxzlib

import (
	"os"
	"path/filepath"
	"testing"

	"gitee.com/MM-Q/comprx/internal/config"
	"gitee.com/MM-Q/comprx/types"
)

func TestUnzlib(t *testing.T) {
	// 创建临时目录
	tempDir := t.TempDir()

	// 创建测试文件
	testFile := filepath.Join(tempDir, "test.txt")
	testContent := "Hello, ZLIB decompression test!"
	if err := os.WriteFile(testFile, []byte(testContent), 0644); err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}

	// 先压缩文件
	zlibFile := filepath.Join(tempDir, "test.zlib")
	cfg := config.New()
	cfg.CompressionLevel = types.CompressionLevelDefault

	err := Zlib(zlibFile, testFile, cfg)
	if err != nil {
		t.Fatalf("ZLIB压缩失败: %v", err)
	}

	// 解压文件
	outputFile := filepath.Join(tempDir, "output.txt")
	err = Unzlib(zlibFile, outputFile, cfg)
	if err != nil {
		t.Fatalf("ZLIB解压失败: %v", err)
	}

	// 验证解压文件内容
	decompressedContent, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("读取解压文件失败: %v", err)
	}

	if string(decompressedContent) != testContent {
		t.Fatalf("解压内容不匹配\n期望: %s\n实际: %s", testContent, string(decompressedContent))
	}

	t.Logf("解压成功，内容匹配")
}

func TestUnzlibToDirectory(t *testing.T) {
	// 创建临时目录
	tempDir := t.TempDir()

	// 创建测试文件
	testFile := filepath.Join(tempDir, "test.txt")
	testContent := "Hello, ZLIB decompression to directory test!"
	if err := os.WriteFile(testFile, []byte(testContent), 0644); err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}

	// 先压缩文件
	zlibFile := filepath.Join(tempDir, "test.zlib")
	cfg := config.New()

	err := Zlib(zlibFile, testFile, cfg)
	if err != nil {
		t.Fatalf("ZLIB压缩失败: %v", err)
	}

	// 创建输出目录
	outputDir := filepath.Join(tempDir, "output")
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		t.Fatalf("创建输出目录失败: %v", err)
	}

	// 解压到目录
	err = Unzlib(zlibFile, outputDir, cfg)
	if err != nil {
		t.Fatalf("ZLIB解压到目录失败: %v", err)
	}

	// 验证解压文件
	expectedFile := filepath.Join(outputDir, "test")
	decompressedContent, err := os.ReadFile(expectedFile)
	if err != nil {
		t.Fatalf("读取解压文件失败: %v", err)
	}

	if string(decompressedContent) != testContent {
		t.Fatalf("解压内容不匹配\n期望: %s\n实际: %s", testContent, string(decompressedContent))
	}

	t.Logf("解压到目录成功，内容匹配")
}

func TestUnzlibOverwrite(t *testing.T) {
	// 创建临时目录
	tempDir := t.TempDir()

	// 创建测试文件
	testFile := filepath.Join(tempDir, "test.txt")
	testContent := "Hello, ZLIB overwrite test!"
	if err := os.WriteFile(testFile, []byte(testContent), 0644); err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}

	// 先压缩文件
	zlibFile := filepath.Join(tempDir, "test.zlib")
	cfg := config.New()

	err := Zlib(zlibFile, testFile, cfg)
	if err != nil {
		t.Fatalf("ZLIB压缩失败: %v", err)
	}

	// 创建已存在的目标文件
	outputFile := filepath.Join(tempDir, "output.txt")
	if err := os.WriteFile(outputFile, []byte("existing content"), 0644); err != nil {
		t.Fatalf("创建已存在文件失败: %v", err)
	}

	// 测试不允许覆盖
	cfg.OverwriteExisting = false
	err = Unzlib(zlibFile, outputFile, cfg)
	if err == nil {
		t.Fatalf("不允许覆盖时解压应该失败，但成功了")
	}

	// 测试允许覆盖
	cfg.OverwriteExisting = true
	err = Unzlib(zlibFile, outputFile, cfg)
	if err != nil {
		t.Fatalf("允许覆盖时解压失败: %v", err)
	}

	// 验证内容被正确覆盖
	decompressedContent, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("读取解压文件失败: %v", err)
	}

	if string(decompressedContent) != testContent {
		t.Fatalf("解压内容不匹配\n期望: %s\n实际: %s", testContent, string(decompressedContent))
	}
}

func TestUnzlibNonExistentFile(t *testing.T) {
	// 创建临时目录
	tempDir := t.TempDir()

	// 不存在的ZLIB文件
	zlibFile := filepath.Join(tempDir, "nonexistent.zlib")
	outputFile := filepath.Join(tempDir, "output.txt")

	// 创建配置
	cfg := config.New()

	// 测试解压不存在的文件（应该失败）
	err := Unzlib(zlibFile, outputFile, cfg)
	if err == nil {
		t.Fatalf("解压不存在的文件应该失败，但成功了")
	}
}

func TestUnzlibInvalidFile(t *testing.T) {
	// 创建临时目录
	tempDir := t.TempDir()

	// 创建无效的ZLIB文件
	zlibFile := filepath.Join(tempDir, "invalid.zlib")
	if err := os.WriteFile(zlibFile, []byte("invalid zlib data"), 0644); err != nil {
		t.Fatalf("创建无效文件失败: %v", err)
	}

	outputFile := filepath.Join(tempDir, "output.txt")

	// 创建配置
	cfg := config.New()

	// 测试解压无效文件（应该失败）
	err := Unzlib(zlibFile, outputFile, cfg)
	if err == nil {
		t.Fatalf("解压无效文件应该失败，但成功了")
	}
}

func TestUnzlibRoundTrip(t *testing.T) {
	// 创建临时目录
	tempDir := t.TempDir()

	// 测试不同大小的文件
	testCases := []struct {
		name    string
		content string
	}{
		{"small", "Hello, World!"},
		{"medium", func() string {
			content := ""
			for i := 0; i < 100; i++ {
				content += "This is a medium-sized test content for ZLIB compression and decompression. "
			}
			return content
		}()},
		{"large", func() string {
			content := ""
			for i := 0; i < 1000; i++ {
				content += "This is a large test content for ZLIB compression and decompression testing. It contains repeated patterns that should compress well. "
			}
			return content
		}()},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 创建测试文件
			testFile := filepath.Join(tempDir, tc.name+".txt")
			if err := os.WriteFile(testFile, []byte(tc.content), 0644); err != nil {
				t.Fatalf("创建测试文件失败: %v", err)
			}

			// 压缩文件
			zlibFile := filepath.Join(tempDir, tc.name+".zlib")
			cfg := config.New()

			err := Zlib(zlibFile, testFile, cfg)
			if err != nil {
				t.Fatalf("ZLIB压缩失败: %v", err)
			}

			// 解压文件
			outputFile := filepath.Join(tempDir, tc.name+"_output.txt")
			err = Unzlib(zlibFile, outputFile, cfg)
			if err != nil {
				t.Fatalf("ZLIB解压失败: %v", err)
			}

			// 验证内容
			decompressedContent, err := os.ReadFile(outputFile)
			if err != nil {
				t.Fatalf("读取解压文件失败: %v", err)
			}

			if string(decompressedContent) != tc.content {
				t.Fatalf("往返压缩解压内容不匹配")
			}

			// 获取文件大小信息
			originalInfo, _ := os.Stat(testFile)
			compressedInfo, _ := os.Stat(zlibFile)
			decompressedInfo, _ := os.Stat(outputFile)

			t.Logf("%s: 原始 %d 字节 -> 压缩 %d 字节 -> 解压 %d 字节, 压缩率 %.2f%%",
				tc.name,
				originalInfo.Size(),
				compressedInfo.Size(),
				decompressedInfo.Size(),
				float64(compressedInfo.Size())/float64(originalInfo.Size())*100)
		})
	}
}

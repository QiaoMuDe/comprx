package cxzlib

import (
	"os"
	"path/filepath"
	"testing"

	"gitee.com/MM-Q/comprx/internal/config"
	"gitee.com/MM-Q/comprx/types"
)

func TestZlib(t *testing.T) {
	// 创建临时目录
	tempDir := t.TempDir()

	// 创建测试文件
	testFile := filepath.Join(tempDir, "test.txt")
	testContent := "Hello, ZLIB compression test!"
	if err := os.WriteFile(testFile, []byte(testContent), 0644); err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}

	// 创建输出文件路径
	zlibFile := filepath.Join(tempDir, "test.zlib")

	// 创建配置
	cfg := config.New()
	cfg.CompressionLevel = types.CompressionLevelDefault

	// 测试压缩
	err := Zlib(zlibFile, testFile, cfg)
	if err != nil {
		t.Fatalf("ZLIB压缩失败: %v", err)
	}

	// 验证压缩文件是否存在
	if _, err := os.Stat(zlibFile); os.IsNotExist(err) {
		t.Fatalf("压缩文件不存在: %s", zlibFile)
	}

	// 验证压缩文件大小
	zlibInfo, err := os.Stat(zlibFile)
	if err != nil {
		t.Fatalf("获取压缩文件信息失败: %v", err)
	}

	if zlibInfo.Size() == 0 {
		t.Fatalf("压缩文件大小为0")
	}

	t.Logf("原始文件大小: %d 字节", len(testContent))
	t.Logf("压缩文件大小: %d 字节", zlibInfo.Size())
}

func TestZlibDirectory(t *testing.T) {
	// 创建临时目录
	tempDir := t.TempDir()
	testDir := filepath.Join(tempDir, "testdir")
	if err := os.MkdirAll(testDir, 0755); err != nil {
		t.Fatalf("创建测试目录失败: %v", err)
	}

	// 创建输出文件路径
	zlibFile := filepath.Join(tempDir, "test.zlib")

	// 创建配置
	cfg := config.New()

	// 测试压缩目录（应该失败）
	err := Zlib(zlibFile, testDir, cfg)
	if err == nil {
		t.Fatalf("压缩目录应该失败，但成功了")
	}

	expectedError := "ZLIB 只支持单文件压缩，不支持目录压缩"
	if err.Error() != expectedError {
		t.Fatalf("期望错误: %s, 实际错误: %s", expectedError, err.Error())
	}
}

func TestZlibOverwrite(t *testing.T) {
	// 创建临时目录
	tempDir := t.TempDir()

	// 创建测试文件
	testFile := filepath.Join(tempDir, "test.txt")
	testContent := "Hello, ZLIB compression test!"
	if err := os.WriteFile(testFile, []byte(testContent), 0644); err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}

	// 创建输出文件路径
	zlibFile := filepath.Join(tempDir, "test.zlib")

	// 创建已存在的目标文件
	if err := os.WriteFile(zlibFile, []byte("existing"), 0644); err != nil {
		t.Fatalf("创建已存在文件失败: %v", err)
	}

	// 创建配置（不允许覆盖）
	cfg := config.New()
	cfg.OverwriteExisting = false

	// 测试压缩（应该失败）
	err := Zlib(zlibFile, testFile, cfg)
	if err == nil {
		t.Fatalf("不允许覆盖时压缩应该失败，但成功了")
	}

	// 测试允许覆盖
	cfg.OverwriteExisting = true
	err = Zlib(zlibFile, testFile, cfg)
	if err != nil {
		t.Fatalf("允许覆盖时压缩失败: %v", err)
	}
}

func TestZlibNonExistentFile(t *testing.T) {
	// 创建临时目录
	tempDir := t.TempDir()

	// 不存在的源文件
	testFile := filepath.Join(tempDir, "nonexistent.txt")
	zlibFile := filepath.Join(tempDir, "test.zlib")

	// 创建配置
	cfg := config.New()

	// 测试压缩不存在的文件（应该失败）
	err := Zlib(zlibFile, testFile, cfg)
	if err == nil {
		t.Fatalf("压缩不存在的文件应该失败，但成功了")
	}
}

func TestZlibDifferentCompressionLevels(t *testing.T) {
	// 创建临时目录
	tempDir := t.TempDir()

	// 创建测试文件（较大的内容以便看到压缩效果）
	testFile := filepath.Join(tempDir, "test.txt")
	testContent := ""
	for i := 0; i < 1000; i++ {
		testContent += "Hello, ZLIB compression test! This is a longer text for better compression testing. "
	}
	if err := os.WriteFile(testFile, []byte(testContent), 0644); err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}

	// 测试不同压缩等级
	levels := []types.CompressionLevel{
		types.CompressionLevelFast,
		types.CompressionLevelDefault,
		types.CompressionLevelBest,
	}

	for _, level := range levels {
		t.Run(level.String(), func(t *testing.T) {
			zlibFile := filepath.Join(tempDir, "test_"+level.String()+".zlib")

			// 创建配置
			cfg := config.New()
			cfg.CompressionLevel = level

			// 测试压缩
			err := Zlib(zlibFile, testFile, cfg)
			if err != nil {
				t.Fatalf("ZLIB压缩失败 (等级 %s): %v", level.String(), err)
			}

			// 验证压缩文件存在
			zlibInfo, err := os.Stat(zlibFile)
			if err != nil {
				t.Fatalf("获取压缩文件信息失败: %v", err)
			}

			if zlibInfo.Size() == 0 {
				t.Fatalf("压缩文件大小为0")
			}

			t.Logf("压缩等级 %s: 原始大小 %d 字节, 压缩大小 %d 字节, 压缩率 %.2f%%",
				level.String(),
				len(testContent),
				zlibInfo.Size(),
				float64(zlibInfo.Size())/float64(len(testContent))*100)
		})
	}
}

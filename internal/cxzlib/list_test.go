package cxzlib

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gitee.com/MM-Q/comprx/internal/config"
	"gitee.com/MM-Q/comprx/types"
)

func TestListZlib(t *testing.T) {
	// 创建临时目录
	tempDir := t.TempDir()

	// 创建测试文件
	testFile := filepath.Join(tempDir, "test.txt")
	testContent := "Hello, ZLIB list test!"
	if err := os.WriteFile(testFile, []byte(testContent), 0644); err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}

	// 压缩文件
	zlibFile := filepath.Join(tempDir, "test.zlib")
	cfg := config.New()
	err := Zlib(zlibFile, testFile, cfg)
	if err != nil {
		t.Fatalf("ZLIB压缩失败: %v", err)
	}

	// 测试列表功能
	archiveInfo, err := ListZlib(zlibFile)
	if err != nil {
		t.Fatalf("列出ZLIB文件信息失败: %v", err)
	}

	// 验证基本信息
	if archiveInfo.Type != types.CompressTypeZlib {
		t.Fatalf("压缩类型不匹配，期望: %s, 实际: %s", types.CompressTypeZlib, archiveInfo.Type)
	}

	if archiveInfo.TotalFiles != 1 {
		t.Fatalf("文件数量不匹配，期望: 1, 实际: %d", archiveInfo.TotalFiles)
	}

	if len(archiveInfo.Files) != 1 {
		t.Fatalf("文件列表长度不匹配，期望: 1, 实际: %d", len(archiveInfo.Files))
	}

	// 验证文件信息
	fileInfo := archiveInfo.Files[0]
	if fileInfo.Name != "test" {
		t.Fatalf("文件名不匹配，期望: test, 实际: %s", fileInfo.Name)
	}

	if fileInfo.Size != int64(len(testContent)) {
		t.Fatalf("文件大小不匹配，期望: %d, 实际: %d", len(testContent), fileInfo.Size)
	}

	if fileInfo.IsDir {
		t.Fatalf("文件不应该是目录")
	}

	if fileInfo.IsSymlink {
		t.Fatalf("文件不应该是符号链接")
	}

	// 验证压缩大小
	if archiveInfo.CompressedSize == 0 {
		t.Fatalf("压缩文件大小不应该为0")
	}

	if archiveInfo.TotalSize != int64(len(testContent)) {
		t.Fatalf("总大小不匹配，期望: %d, 实际: %d", len(testContent), archiveInfo.TotalSize)
	}

	t.Logf("列表功能测试成功: 文件名=%s, 原始大小=%d, 压缩大小=%d",
		fileInfo.Name, fileInfo.Size, fileInfo.CompressedSize)
}

func TestListZlibLimit(t *testing.T) {
	// 创建临时目录
	tempDir := t.TempDir()

	// 创建测试文件
	testFile := filepath.Join(tempDir, "test.txt")
	testContent := "Hello, ZLIB list limit test!"
	if err := os.WriteFile(testFile, []byte(testContent), 0644); err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}

	// 压缩文件
	zlibFile := filepath.Join(tempDir, "test.zlib")
	cfg := config.New()
	err := Zlib(zlibFile, testFile, cfg)
	if err != nil {
		t.Fatalf("ZLIB压缩失败: %v", err)
	}

	// 测试限制列表功能
	archiveInfo, err := ListZlibLimit(zlibFile, 5)
	if err != nil {
		t.Fatalf("限制列出ZLIB文件信息失败: %v", err)
	}

	// ZLIB只有一个文件，限制不应该影响结果
	if archiveInfo.TotalFiles != 1 {
		t.Fatalf("文件数量不匹配，期望: 1, 实际: %d", archiveInfo.TotalFiles)
	}

	if len(archiveInfo.Files) != 1 {
		t.Fatalf("文件列表长度不匹配，期望: 1, 实际: %d", len(archiveInfo.Files))
	}

	t.Logf("限制列表功能测试成功")
}

func TestListZlibMatch(t *testing.T) {
	// 创建临时目录
	tempDir := t.TempDir()

	// 创建测试文件
	testFile := filepath.Join(tempDir, "test.txt")
	testContent := "Hello, ZLIB list match test!"
	if err := os.WriteFile(testFile, []byte(testContent), 0644); err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}

	// 压缩文件
	zlibFile := filepath.Join(tempDir, "test.zlib")
	cfg := config.New()
	err := Zlib(zlibFile, testFile, cfg)
	if err != nil {
		t.Fatalf("ZLIB压缩失败: %v", err)
	}

	// 测试匹配模式 - 匹配的情况
	archiveInfo, err := ListZlibMatch(zlibFile, "test")
	if err != nil {
		t.Fatalf("匹配列出ZLIB文件信息失败: %v", err)
	}

	if archiveInfo.TotalFiles != 1 {
		t.Fatalf("匹配文件数量不匹配，期望: 1, 实际: %d", archiveInfo.TotalFiles)
	}

	if len(archiveInfo.Files) != 1 {
		t.Fatalf("匹配文件列表长度不匹配，期望: 1, 实际: %d", len(archiveInfo.Files))
	}

	// 测试匹配模式 - 不匹配的情况
	archiveInfo, err = ListZlibMatch(zlibFile, "nomatch")
	if err != nil {
		t.Fatalf("不匹配列出ZLIB文件信息失败: %v", err)
	}

	if archiveInfo.TotalFiles != 0 {
		t.Fatalf("不匹配文件数量应该为0，实际: %d", archiveInfo.TotalFiles)
	}

	if len(archiveInfo.Files) != 0 {
		t.Fatalf("不匹配文件列表长度应该为0，实际: %d", len(archiveInfo.Files))
	}

	if archiveInfo.TotalSize != 0 {
		t.Fatalf("不匹配总大小应该为0，实际: %d", archiveInfo.TotalSize)
	}

	t.Logf("匹配列表功能测试成功")
}

func TestListZlibNonExistentFile(t *testing.T) {
	// 创建临时目录
	tempDir := t.TempDir()

	// 不存在的ZLIB文件
	zlibFile := filepath.Join(tempDir, "nonexistent.zlib")

	// 测试列出不存在的文件（应该失败）
	_, err := ListZlib(zlibFile)
	if err == nil {
		t.Fatalf("列出不存在的文件应该失败，但成功了")
	}
}

func TestListZlibInvalidFile(t *testing.T) {
	// 创建临时目录
	tempDir := t.TempDir()

	// 创建无效的ZLIB文件
	zlibFile := filepath.Join(tempDir, "invalid.zlib")
	if err := os.WriteFile(zlibFile, []byte("invalid zlib data"), 0644); err != nil {
		t.Fatalf("创建无效文件失败: %v", err)
	}

	// 测试列出无效文件（应该失败）
	_, err := ListZlib(zlibFile)
	if err == nil {
		t.Fatalf("列出无效文件应该失败，但成功了")
	}
}

func TestListZlibDifferentExtensions(t *testing.T) {
	// 创建临时目录
	tempDir := t.TempDir()

	// 只测试支持的扩展名
	testCases := []struct {
		name         string
		zlibFileName string
		expectedName string
	}{
		{"standard_zlib", "test.zlib", "test"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 创建测试文件
			testFile := filepath.Join(tempDir, "source.txt")
			testContent := "Hello, ZLIB extension test!"
			if err := os.WriteFile(testFile, []byte(testContent), 0644); err != nil {
				t.Fatalf("创建测试文件失败: %v", err)
			}

			// 压缩文件
			zlibFile := filepath.Join(tempDir, tc.zlibFileName)
			cfg := config.New()
			err := Zlib(zlibFile, testFile, cfg)
			if err != nil {
				t.Fatalf("ZLIB压缩失败: %v", err)
			}

			// 测试列表功能
			archiveInfo, err := ListZlib(zlibFile)
			if err != nil {
				t.Fatalf("列出ZLIB文件信息失败: %v", err)
			}

			// 验证文件名
			if len(archiveInfo.Files) != 1 {
				t.Fatalf("文件列表长度不匹配，期望: 1, 实际: %d", len(archiveInfo.Files))
			}

			fileInfo := archiveInfo.Files[0]
			if fileInfo.Name != tc.expectedName {
				t.Fatalf("文件名不匹配，期望: %s, 实际: %s", tc.expectedName, fileInfo.Name)
			}

			t.Logf("扩展名测试成功: %s -> %s", tc.zlibFileName, tc.expectedName)
		})
	}
}

func TestListZlibUnsupportedExtension(t *testing.T) {
	// 创建临时目录
	tempDir := t.TempDir()

	// 创建测试文件
	testFile := filepath.Join(tempDir, "source.txt")
	testContent := "Hello, ZLIB unsupported extension test!"
	if err := os.WriteFile(testFile, []byte(testContent), 0644); err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}

	// 压缩文件到标准扩展名
	zlibFile := filepath.Join(tempDir, "test.zlib")
	cfg := config.New()
	err := Zlib(zlibFile, testFile, cfg)
	if err != nil {
		t.Fatalf("ZLIB压缩失败: %v", err)
	}

	// 重命名为不支持的扩展名
	unsupportedFile := filepath.Join(tempDir, "test.data")
	err = os.Rename(zlibFile, unsupportedFile)
	if err != nil {
		t.Fatalf("重命名文件失败: %v", err)
	}

	// 测试列表功能（应该失败）
	_, err = ListZlib(unsupportedFile)
	if err == nil {
		t.Fatalf("列出不支持扩展名的文件应该失败，但成功了")
	}

	// 验证错误信息包含格式检测失败
	if !strings.Contains(err.Error(), "检测压缩格式失败") {
		t.Fatalf("错误信息不正确，期望包含'检测压缩格式失败'，实际: %s", err.Error())
	}

	t.Logf("不支持扩展名测试成功，错误: %s", err.Error())
}

func TestListZlibWildcardMatch(t *testing.T) {
	// 创建临时目录
	tempDir := t.TempDir()

	// 创建测试文件
	testFile := filepath.Join(tempDir, "test.txt")
	testContent := "Hello, ZLIB wildcard test!"
	if err := os.WriteFile(testFile, []byte(testContent), 0644); err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}

	// 压缩文件
	zlibFile := filepath.Join(tempDir, "test.zlib")
	cfg := config.New()
	err := Zlib(zlibFile, testFile, cfg)
	if err != nil {
		t.Fatalf("ZLIB压缩失败: %v", err)
	}

	// 测试通配符匹配
	testCases := []struct {
		pattern     string
		shouldMatch bool
	}{
		{"*", true},
		{"test", true},
		{"te*", true},
		{"*st", true},
		{"t?st", true},
		{"nomatch", false},
		{"test*extra", false},
	}

	for _, tc := range testCases {
		t.Run(tc.pattern, func(t *testing.T) {
			archiveInfo, err := ListZlibMatch(zlibFile, tc.pattern)
			if err != nil {
				t.Fatalf("匹配模式 '%s' 失败: %v", tc.pattern, err)
			}

			if tc.shouldMatch {
				if archiveInfo.TotalFiles != 1 {
					t.Fatalf("模式 '%s' 应该匹配，但没有匹配到文件", tc.pattern)
				}
			} else {
				if archiveInfo.TotalFiles != 0 {
					t.Fatalf("模式 '%s' 不应该匹配，但匹配到了文件", tc.pattern)
				}
			}

			t.Logf("通配符测试成功: '%s' -> 匹配: %v", tc.pattern, tc.shouldMatch)
		})
	}
}

func BenchmarkListZlib(b *testing.B) {
	// 创建临时目录
	tempDir := b.TempDir()

	// 创建测试文件
	testFile := filepath.Join(tempDir, "test.txt")
	testContent := "Hello, ZLIB benchmark test!"
	if err := os.WriteFile(testFile, []byte(testContent), 0644); err != nil {
		b.Fatalf("创建测试文件失败: %v", err)
	}

	// 压缩文件
	zlibFile := filepath.Join(tempDir, "test.zlib")
	cfg := config.New()
	err := Zlib(zlibFile, testFile, cfg)
	if err != nil {
		b.Fatalf("ZLIB压缩失败: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := ListZlib(zlibFile)
		if err != nil {
			b.Fatalf("列表操作失败: %v", err)
		}
	}
}

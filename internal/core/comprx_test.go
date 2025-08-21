package core

import (
	"flag"
	"os"
	"path/filepath"
	"testing"

	"gitee.com/MM-Q/comprx/types"
)

// TestMain 全局测试入口，控制非verbose模式下的输出重定向
func TestMain(m *testing.M) {
	flag.Parse() // 解析命令行参数
	// 保存原始标准输出和错误输出
	originalStdout := os.Stdout
	originalStderr := os.Stderr
	var nullFile *os.File
	var err error

	// 非verbose模式下重定向到空设备
	if !testing.Verbose() {
		nullFile, err = os.OpenFile(os.DevNull, os.O_WRONLY, 0666)
		if err != nil {
			panic("无法打开空设备文件: " + err.Error())
		}
		os.Stdout = nullFile
		os.Stderr = nullFile
	}

	// 运行所有测试
	exitCode := m.Run()

	// 恢复原始输出
	if !testing.Verbose() {
		os.Stdout = originalStdout
		os.Stderr = originalStderr
		_ = nullFile.Close()
	}

	os.Exit(exitCode)
}

// TestNew 测试构造函数
func TestNew(t *testing.T) {
	c := New()
	if c == nil {
		t.Fatal("New() 返回 nil")
	}
	if c.Config == nil {
		t.Fatal("config 未初始化")
	}
	if c.Config.CompressionLevel != types.CompressionLevelDefault {
		t.Errorf("期望压缩级别为 %v, 实际为 %v", types.CompressionLevelDefault, c.Config.CompressionLevel)
	}
	if c.Config.OverwriteExisting != false {
		t.Errorf("期望 OverwriteExisting 为 false, 实际为 %v", c.Config.OverwriteExisting)
	}
}

// TestNewComprx 测试NewComprx构造函数
func TestNewComprx(t *testing.T) {
	c := NewComprx()
	if c == nil {
		t.Fatal("NewComprx() 返回 nil")
	}
	if c.Config == nil {
		t.Fatal("config 未初始化")
	}
}

// TestWithOverwriteExisting 测试链式设置覆盖选项
func TestWithOverwriteExisting(t *testing.T) {
	c := New()

	// 测试设置为 true
	result := c.WithOverwriteExisting(true)
	if result != c {
		t.Error("WithOverwriteExisting 应该返回同一个实例")
	}
	if !c.Config.OverwriteExisting {
		t.Error("OverwriteExisting 应该被设置为 true")
	}

	// 测试设置为 false
	c.WithOverwriteExisting(false)
	if c.Config.OverwriteExisting {
		t.Error("OverwriteExisting 应该被设置为 false")
	}
}

// TestSetOverwriteExisting 测试设置覆盖选项
func TestSetOverwriteExisting(t *testing.T) {
	c := New()

	// 测试设置为 true
	c.SetOverwriteExisting(true)
	if !c.Config.OverwriteExisting {
		t.Error("OverwriteExisting 应该被设置为 true")
	}

	// 测试设置为 false
	c.SetOverwriteExisting(false)
	if c.Config.OverwriteExisting {
		t.Error("OverwriteExisting 应该被设置为 false")
	}
}

// TestWithCompressionLevel 测试链式设置压缩级别
func TestWithCompressionLevel(t *testing.T) {
	c := New()

	testCases := []types.CompressionLevel{
		types.CompressionLevelNone,
		types.CompressionLevelFast,
		types.CompressionLevelBest,
		types.CompressionLevelHuffmanOnly,
	}

	for _, level := range testCases {
		result := c.WithCompressionLevel(level)
		if result != c {
			t.Error("WithCompressionLevel 应该返回同一个实例")
		}
		if c.Config.CompressionLevel != level {
			t.Errorf("期望压缩级别为 %v, 实际为 %v", level, c.Config.CompressionLevel)
		}
	}
}

// TestSetCompressionLevel 测试设置压缩级别
func TestSetCompressionLevel(t *testing.T) {
	c := New()

	testCases := []types.CompressionLevel{
		types.CompressionLevelNone,
		types.CompressionLevelFast,
		types.CompressionLevelBest,
		types.CompressionLevelHuffmanOnly,
	}

	for _, level := range testCases {
		c.SetCompressionLevel(level)
		if c.Config.CompressionLevel != level {
			t.Errorf("期望压缩级别为 %v, 实际为 %v", level, c.Config.CompressionLevel)
		}
	}
}

// TestPackEmptyPaths 测试空路径参数
func TestPackEmptyPaths(t *testing.T) {
	c := New()

	testCases := []struct {
		name string
		dst  string
		src  string
	}{
		{"空源路径", "test.zip", ""},
		{"空目标路径", "", "testfile.txt"},
		{"两个都为空", "", ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := c.Pack(tc.dst, tc.src)
			if err == nil {
				t.Error("期望返回错误，但没有错误")
			}
			if err.Error() != "源文件路径或目标文件路径不能为空" {
				t.Errorf("期望错误信息为 '源文件路径或目标文件路径不能为空', 实际为 '%s'", err.Error())
			}
		})
	}
}

// TestPackUnsupportedFormat 测试不支持的压缩格式
func TestPackUnsupportedFormat(t *testing.T) {
	c := New()

	// 创建临时测试文件
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatal(err)
	}

	testCases := []string{
		"test.rar",
		"test.7z",
		"test.unknown",
		"test",
	}

	for _, dst := range testCases {
		t.Run(dst, func(t *testing.T) {
			dstPath := filepath.Join(tempDir, dst)
			err := c.Pack(dstPath, testFile)
			if err == nil {
				t.Error("期望返回错误，但没有错误")
			}
		})
	}
}

// TestPackBz2Format 测试bz2格式（应该返回错误）
func TestPackBz2Format(t *testing.T) {
	c := New()

	// 创建临时测试文件
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatal(err)
	}

	testCases := []string{
		"test.bz2",
		"test.bzip2",
	}

	for _, dst := range testCases {
		t.Run(dst, func(t *testing.T) {
			dstPath := filepath.Join(tempDir, dst)
			err := c.Pack(dstPath, testFile)
			if err == nil {
				t.Error("期望返回错误，但没有错误")
			}
			expectedMsg := "暂不支持"
			if !contains(err.Error(), expectedMsg) {
				t.Errorf("期望错误信息包含 '%s', 实际为 '%s'", expectedMsg, err.Error())
			}
		})
	}
}

// TestPackNonExistentSource 测试不存在的源文件
func TestPackNonExistentSource(t *testing.T) {
	c := New()

	tempDir := t.TempDir()
	nonExistentFile := filepath.Join(tempDir, "nonexistent.txt")
	dstFile := filepath.Join(tempDir, "test.zip")

	err := c.Pack(dstFile, nonExistentFile)
	if err == nil {
		t.Error("期望返回错误，但没有错误")
	}
}

// TestPackOverwriteExisting 测试覆盖已存在文件
func TestPackOverwriteExisting(t *testing.T) {
	tempDir := t.TempDir()

	// 创建源文件
	srcFile := filepath.Join(tempDir, "source.txt")
	if err := os.WriteFile(srcFile, []byte("test content"), 0644); err != nil {
		t.Fatal(err)
	}

	// 创建目标文件
	dstFile := filepath.Join(tempDir, "test.zip")
	if err := os.WriteFile(dstFile, []byte("existing content"), 0644); err != nil {
		t.Fatal(err)
	}

	// 测试不允许覆盖
	c1 := New().WithOverwriteExisting(false)
	err := c1.Pack(dstFile, srcFile)
	if err == nil {
		t.Error("期望返回错误，但没有错误")
	}
	expectedMsg := "已存在"
	if !contains(err.Error(), expectedMsg) {
		t.Errorf("期望错误信息包含 '%s', 实际为 '%s'", expectedMsg, err.Error())
	}

	// 测试允许覆盖
	c2 := New().WithOverwriteExisting(true)
	err = c2.Pack(dstFile, srcFile)
	if err != nil {
		t.Errorf("不期望返回错误，但得到错误: %v", err)
	}
}

// TestUnpackEmptySource 测试空源路径
func TestUnpackEmptySource(t *testing.T) {
	c := New()

	err := c.Unpack("", "target")
	if err == nil {
		t.Error("期望返回错误，但没有错误")
	}
	if err.Error() != "源文件路径不能为空" {
		t.Errorf("期望错误信息为 '源文件路径不能为空', 实际为 '%s'", err.Error())
	}
}

// TestUnpackNonExistentFile 测试不存在的源文件
func TestUnpackNonExistentFile(t *testing.T) {
	c := New()

	tempDir := t.TempDir()
	nonExistentFile := filepath.Join(tempDir, "nonexistent.zip")
	targetDir := filepath.Join(tempDir, "target")

	err := c.Unpack(nonExistentFile, targetDir)
	if err == nil {
		t.Error("期望返回错误，但没有错误")
	}
	expectedMsg := "不存在"
	if !contains(err.Error(), expectedMsg) {
		t.Errorf("期望错误信息包含 '%s', 实际为 '%s'", expectedMsg, err.Error())
	}
}

// TestUnpackAutoGenerateTargetDir 测试自动生成目标目录
func TestUnpackAutoGenerateTargetDir(t *testing.T) {
	c := New()

	tempDir := t.TempDir()

	// 创建一个真实的压缩文件来测试目标目录自动生成
	srcFile := filepath.Join(tempDir, "source.txt")
	if err := os.WriteFile(srcFile, []byte("test content for auto target dir"), 0644); err != nil {
		t.Fatal(err)
	}

	// 先创建一个真实的压缩文件
	testArchive := filepath.Join(tempDir, "test.tar.gz")
	if err := c.Pack(testArchive, srcFile); err != nil {
		t.Fatal(err)
	}

	// 调用Unpack，目标目录为空字符串，应该自动生成目标目录
	err := c.Unpack(testArchive, "")
	if err != nil {
		t.Errorf("解压时不应该返回错误，但得到: %v", err)
	}

	// 检查自动生成的目标目录是否存在
	expectedTargetDir := filepath.Join(tempDir, "test")
	if _, err := os.Stat(expectedTargetDir); os.IsNotExist(err) {
		t.Errorf("自动生成的目标目录 %s 不存在", expectedTargetDir)
	}

	// 检查解压的文件是否存在
	extractedFile := filepath.Join(expectedTargetDir, "source.txt")
	if _, err := os.Stat(extractedFile); os.IsNotExist(err) {
		t.Errorf("解压的文件 %s 不存在", extractedFile)
	}
}

// TestChainedConfiguration 测试链式配置
func TestChainedConfiguration(t *testing.T) {
	c := New().
		WithOverwriteExisting(true).
		WithCompressionLevel(types.CompressionLevelBest)

	if !c.Config.OverwriteExisting {
		t.Error("OverwriteExisting 应该为 true")
	}
	if c.Config.CompressionLevel != types.CompressionLevelBest {
		t.Errorf("期望压缩级别为 %v, 实际为 %v", types.CompressionLevelBest, c.Config.CompressionLevel)
	}
}

// TestPackSupportedFormats 测试支持的压缩格式
func TestPackSupportedFormats(t *testing.T) {
	tempDir := t.TempDir()

	// 创建源文件
	srcFile := filepath.Join(tempDir, "source.txt")
	if err := os.WriteFile(srcFile, []byte("test content for compression"), 0644); err != nil {
		t.Fatal(err)
	}

	c := New().WithOverwriteExisting(true)

	supportedFormats := []string{
		"test.zip",
		"test.tar",
		"test.tgz",
		"test.tar.gz",
		"test.gz",
	}

	for _, format := range supportedFormats {
		t.Run(format, func(t *testing.T) {
			dstFile := filepath.Join(tempDir, format)
			err := c.Pack(dstFile, srcFile)
			if err != nil {
				t.Errorf("压缩格式 %s 应该被支持，但得到错误: %v", format, err)
			}

			// 检查文件是否创建
			if _, err := os.Stat(dstFile); os.IsNotExist(err) {
				t.Errorf("压缩文件 %s 未创建", dstFile)
			}
		})
	}
}

// TestPackDirectory 测试压缩目录
func TestPackDirectory(t *testing.T) {
	tempDir := t.TempDir()

	// 创建测试目录结构
	srcDir := filepath.Join(tempDir, "testdir")
	if err := os.MkdirAll(srcDir, 0755); err != nil {
		t.Fatal(err)
	}

	// 在目录中创建文件
	testFile := filepath.Join(srcDir, "file.txt")
	if err := os.WriteFile(testFile, []byte("directory content"), 0644); err != nil {
		t.Fatal(err)
	}

	c := New()
	dstFile := filepath.Join(tempDir, "directory.zip")

	err := c.Pack(dstFile, srcDir)
	if err != nil {
		t.Errorf("压缩目录时不应该返回错误，但得到: %v", err)
	}

	// 检查压缩文件是否创建
	if _, err := os.Stat(dstFile); os.IsNotExist(err) {
		t.Error("压缩文件未创建")
	}
}

// 辅助函数：检查字符串是否包含子字符串
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > len(substr) && containsHelper(s, substr)))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// BenchmarkPack 性能测试
func BenchmarkPack(b *testing.B) {
	tempDir := b.TempDir()

	// 创建测试文件
	srcFile := filepath.Join(tempDir, "benchmark.txt")
	content := make([]byte, 1024*1024) // 1MB 测试数据
	for i := range content {
		content[i] = byte(i % 256)
	}
	if err := os.WriteFile(srcFile, content, 0644); err != nil {
		b.Fatal(err)
	}

	c := New().WithOverwriteExisting(true)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dstFile := filepath.Join(tempDir, "benchmark.zip")
		if err := c.Pack(dstFile, srcFile); err != nil {
			b.Fatal(err)
		}
		// 清理文件以便下次测试
		_ = os.Remove(dstFile)
	}
}

// BenchmarkNew 构造函数性能测试
func BenchmarkNew(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = New()
	}
}

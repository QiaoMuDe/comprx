package utils

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestGetSize_SingleFile 测试获取单个文件大小
func TestGetSize_SingleFile(t *testing.T) {
	// 创建临时文件
	tmpFile, err := os.CreateTemp("", "test_file_*.txt")
	if err != nil {
		t.Fatalf("创建临时文件失败: %v", err)
	}
	defer func() {
		if err := os.Remove(tmpFile.Name()); err != nil {
			t.Logf("清理临时文件失败: %v", err)
		}
	}()

	// 写入测试数据
	testData := "Hello, World! This is a test file."
	if _, err := tmpFile.WriteString(testData); err != nil {
		t.Fatalf("写入测试数据失败: %v", err)
	}
	if err := tmpFile.Close(); err != nil {
		t.Fatalf("关闭临时文件失败: %v", err)
	}

	// 测试获取文件大小
	size, err := GetSize(tmpFile.Name())
	if err != nil {
		t.Errorf("GetSize() 返回错误: %v", err)
	}

	expectedSize := int64(len(testData))
	if size != expectedSize {
		t.Errorf("GetSize() = %d, 期望 %d", size, expectedSize)
	}
}

// TestGetSize_EmptyFile 测试获取空文件大小
func TestGetSize_EmptyFile(t *testing.T) {
	// 创建空文件
	tmpFile, err := os.CreateTemp("", "empty_file_*.txt")
	if err != nil {
		t.Fatalf("创建临时文件失败: %v", err)
	}
	defer func() {
		if err := os.Remove(tmpFile.Name()); err != nil {
			t.Logf("清理临时文件失败: %v", err)
		}
	}()
	if err := tmpFile.Close(); err != nil {
		t.Fatalf("关闭临时文件失败: %v", err)
	}

	// 测试获取空文件大小
	size, err := GetSize(tmpFile.Name())
	if err != nil {
		t.Errorf("GetSize() 返回错误: %v", err)
	}

	if size != 0 {
		t.Errorf("GetSize() = %d, 期望 0", size)
	}
}

// TestGetSize_Directory 测试获取目录大小
func TestGetSize_Directory(t *testing.T) {
	// 创建临时目录
	tmpDir, err := os.MkdirTemp("", "test_dir_*")
	if err != nil {
		t.Fatalf("创建临时目录失败: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			t.Logf("清理临时目录失败: %v", err)
		}
	}()

	// 创建测试文件
	testFiles := []struct {
		name    string
		content string
	}{
		{"file1.txt", "Hello"},
		{"file2.txt", "World"},
		{"subdir/file3.txt", "Test"},
	}

	var expectedSize int64
	for _, tf := range testFiles {
		filePath := filepath.Join(tmpDir, tf.name)

		// 创建子目录（如果需要）
		if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
			t.Fatalf("创建子目录失败: %v", err)
		}

		// 创建文件并写入内容
		if err := os.WriteFile(filePath, []byte(tf.content), 0644); err != nil {
			t.Fatalf("创建测试文件失败: %v", err)
		}

		expectedSize += int64(len(tf.content))
	}

	// 测试获取目录大小
	size, err := GetSize(tmpDir)
	if err != nil {
		t.Errorf("GetSize() 返回错误: %v", err)
	}

	if size != expectedSize {
		t.Errorf("GetSize() = %d, 期望 %d", size, expectedSize)
	}
}

// TestGetSize_EmptyDirectory 测试获取空目录大小
func TestGetSize_EmptyDirectory(t *testing.T) {
	// 创建空目录
	tmpDir, err := os.MkdirTemp("", "empty_dir_*")
	if err != nil {
		t.Fatalf("创建临时目录失败: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			t.Logf("清理临时目录失败: %v", err)
		}
	}()

	// 测试获取空目录大小
	size, err := GetSize(tmpDir)
	if err != nil {
		t.Errorf("GetSize() 返回错误: %v", err)
	}

	if size != 0 {
		t.Errorf("GetSize() = %d, 期望 0", size)
	}
}

// TestGetSize_NonExistentPath 测试不存在的路径
func TestGetSize_NonExistentPath(t *testing.T) {
	nonExistentPath := "/path/that/does/not/exist"

	size, err := GetSize(nonExistentPath)
	if err == nil {
		t.Error("GetSize() 应该返回错误，但没有")
	}

	if size != 0 {
		t.Errorf("GetSize() = %d, 期望 0", size)
	}

	// 检查错误消息
	if !strings.Contains(err.Error(), "路径不存在") {
		t.Errorf("错误消息不正确: %v", err)
	}
}

// TestGetSize_NestedDirectories 测试嵌套目录
func TestGetSize_NestedDirectories(t *testing.T) {
	// 创建临时目录
	tmpDir, err := os.MkdirTemp("", "nested_dir_*")
	if err != nil {
		t.Fatalf("创建临时目录失败: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			t.Logf("清理临时目录失败: %v", err)
		}
	}()

	// 创建嵌套目录结构
	testStructure := []struct {
		path    string
		content string
	}{
		{"level1/file1.txt", "content1"},
		{"level1/level2/file2.txt", "content2"},
		{"level1/level2/level3/file3.txt", "content3"},
		{"level1/level2/level3/file4.txt", "content4"},
	}

	var expectedSize int64
	for _, ts := range testStructure {
		fullPath := filepath.Join(tmpDir, ts.path)

		// 创建目录
		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			t.Fatalf("创建目录失败: %v", err)
		}

		// 创建文件
		if err := os.WriteFile(fullPath, []byte(ts.content), 0644); err != nil {
			t.Fatalf("创建文件失败: %v", err)
		}

		expectedSize += int64(len(ts.content))
	}

	// 测试获取嵌套目录大小
	size, err := GetSize(tmpDir)
	if err != nil {
		t.Errorf("GetSize() 返回错误: %v", err)
	}

	if size != expectedSize {
		t.Errorf("GetSize() = %d, 期望 %d", size, expectedSize)
	}
}

// TestGetSizeOrZero_Success 测试 GetSizeOrZero 成功情况
func TestGetSizeOrZero_Success(t *testing.T) {
	// 创建临时文件
	tmpFile, err := os.CreateTemp("", "test_file_*.txt")
	if err != nil {
		t.Fatalf("创建临时文件失败: %v", err)
	}
	defer func() {
		if err := os.Remove(tmpFile.Name()); err != nil {
			t.Logf("清理临时文件失败: %v", err)
		}
	}()

	// 写入测试数据
	testData := "Test content"
	if _, err := tmpFile.WriteString(testData); err != nil {
		t.Fatalf("写入测试数据失败: %v", err)
	}
	if err := tmpFile.Close(); err != nil {
		t.Fatalf("关闭临时文件失败: %v", err)
	}

	// 测试 GetSizeOrZero
	size := GetSizeOrZero(tmpFile.Name())
	expectedSize := int64(len(testData))

	if size != expectedSize {
		t.Errorf("GetSizeOrZero() = %d, 期望 %d", size, expectedSize)
	}
}

// TestGetSizeOrZero_Error 测试 GetSizeOrZero 错误情况
func TestGetSizeOrZero_Error(t *testing.T) {
	nonExistentPath := "/path/that/does/not/exist"

	size := GetSizeOrZero(nonExistentPath)
	if size != 0 {
		t.Errorf("GetSizeOrZero() = %d, 期望 0", size)
	}
}

// TestGetSize_LargeFile 测试大文件
func TestGetSize_LargeFile(t *testing.T) {
	// 创建临时文件
	tmpFile, err := os.CreateTemp("", "large_file_*.txt")
	if err != nil {
		t.Fatalf("创建临时文件失败: %v", err)
	}
	defer func() {
		if err := os.Remove(tmpFile.Name()); err != nil {
			t.Logf("清理临时文件失败: %v", err)
		}
	}()

	// 写入较大的数据（1MB）
	largeData := strings.Repeat("A", 1024*1024)
	if _, err := tmpFile.WriteString(largeData); err != nil {
		t.Fatalf("写入大文件数据失败: %v", err)
	}
	if err := tmpFile.Close(); err != nil {
		t.Fatalf("关闭临时文件失败: %v", err)
	}

	// 测试获取大文件大小
	size, err := GetSize(tmpFile.Name())
	if err != nil {
		t.Errorf("GetSize() 返回错误: %v", err)
	}

	expectedSize := int64(len(largeData))
	if size != expectedSize {
		t.Errorf("GetSize() = %d, 期望 %d", size, expectedSize)
	}
}

// TestGetSize_MixedContent 测试包含不同类型文件的目录
func TestGetSize_MixedContent(t *testing.T) {
	// 创建临时目录
	tmpDir, err := os.MkdirTemp("", "mixed_dir_*")
	if err != nil {
		t.Fatalf("创建临时目录失败: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			t.Logf("清理临时目录失败: %v", err)
		}
	}()

	// 创建普通文件
	file1Path := filepath.Join(tmpDir, "regular_file.txt")
	file1Content := "regular file content"
	if err := os.WriteFile(file1Path, []byte(file1Content), 0644); err != nil {
		t.Fatalf("创建普通文件失败: %v", err)
	}

	// 创建空文件
	file2Path := filepath.Join(tmpDir, "empty_file.txt")
	if err := os.WriteFile(file2Path, []byte{}, 0644); err != nil {
		t.Fatalf("创建空文件失败: %v", err)
	}

	// 创建子目录
	subDirPath := filepath.Join(tmpDir, "subdir")
	if err := os.Mkdir(subDirPath, 0755); err != nil {
		t.Fatalf("创建子目录失败: %v", err)
	}

	// 在子目录中创建文件
	file3Path := filepath.Join(subDirPath, "sub_file.txt")
	file3Content := "sub file content"
	if err := os.WriteFile(file3Path, []byte(file3Content), 0644); err != nil {
		t.Fatalf("创建子目录文件失败: %v", err)
	}

	// 测试获取混合内容目录大小
	size, err := GetSize(tmpDir)
	if err != nil {
		t.Errorf("GetSize() 返回错误: %v", err)
	}

	expectedSize := int64(len(file1Content) + len(file3Content))
	if size != expectedSize {
		t.Errorf("GetSize() = %d, 期望 %d", size, expectedSize)
	}
}

// BenchmarkGetSize_SingleFile 基准测试：单个文件
func BenchmarkGetSize_SingleFile(b *testing.B) {
	// 创建临时文件
	tmpFile, err := os.CreateTemp("", "bench_file_*.txt")
	if err != nil {
		b.Fatalf("创建临时文件失败: %v", err)
	}
	defer func() {
		if err := os.Remove(tmpFile.Name()); err != nil {
			b.Logf("清理临时文件失败: %v", err)
		}
	}()

	// 写入测试数据
	testData := "Benchmark test data"
	if _, err := tmpFile.WriteString(testData); err != nil {
		b.Fatalf("写入测试数据失败: %v", err)
	}
	if err := tmpFile.Close(); err != nil {
		b.Fatalf("关闭临时文件失败: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := GetSize(tmpFile.Name())
		if err != nil {
			b.Errorf("GetSize() 返回错误: %v", err)
		}
	}
}

// BenchmarkGetSize_Directory 基准测试：目录
func BenchmarkGetSize_Directory(b *testing.B) {
	// 创建临时目录
	tmpDir, err := os.MkdirTemp("", "bench_dir_*")
	if err != nil {
		b.Fatalf("创建临时目录失败: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			b.Logf("清理临时目录失败: %v", err)
		}
	}()

	// 创建多个测试文件
	for i := 0; i < 10; i++ {
		filePath := filepath.Join(tmpDir, "file"+string(rune('0'+i))+".txt")
		content := "Content of file " + string(rune('0'+i))
		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			b.Fatalf("创建测试文件失败: %v", err)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := GetSize(tmpDir)
		if err != nil {
			b.Errorf("GetSize() 返回错误: %v", err)
		}
	}
}

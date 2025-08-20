package comprx

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestListEmptyPaths 测试空路径参数
func TestListEmptyPaths(t *testing.T) {
	testCases := []struct {
		name string
		path string
	}{
		{"空路径", ""},
		{"空白字符串", "   "},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := List(tc.path)
			if err == nil {
				t.Error("期望返回错误，但没有错误")
			}
			expectedMsg := "检测压缩格式失败"
			if !strings.Contains(err.Error(), expectedMsg) {
				t.Errorf("期望错误信息包含 '%s', 实际为 '%s'", expectedMsg, err.Error())
			}
		})
	}
}

// TestListNonExistentFile 测试不存在的文件
func TestListNonExistentFile(t *testing.T) {
	tempDir := t.TempDir()
	nonExistentFile := filepath.Join(tempDir, "nonexistent.zip")

	_, err := List(nonExistentFile)
	if err == nil {
		t.Error("期望返回错误，但没有错误")
	}
	expectedMsg := "不存在"
	if !strings.Contains(err.Error(), expectedMsg) {
		t.Errorf("期望错误信息包含 '%s', 实际为 '%s'", expectedMsg, err.Error())
	}
}

// TestListUnsupportedFormat 测试不支持的压缩格式
func TestListUnsupportedFormat(t *testing.T) {
	tempDir := t.TempDir()

	// 创建一个不支持格式的文件
	unsupportedFile := filepath.Join(tempDir, "test.rar")
	if err := os.WriteFile(unsupportedFile, []byte("fake rar content"), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := List(unsupportedFile)
	if err == nil {
		t.Error("期望返回错误，但没有错误")
	}
	expectedMsg := "检测压缩格式失败"
	if !strings.Contains(err.Error(), expectedMsg) {
		t.Errorf("期望错误信息包含 '%s', 实际为 '%s'", expectedMsg, err.Error())
	}
}

// TestListZipFile 测试列出ZIP文件内容
func TestListZipFile(t *testing.T) {
	tempDir := t.TempDir()

	// 创建测试目录结构
	srcDir := filepath.Join(tempDir, "source")
	if err := os.MkdirAll(srcDir, 0755); err != nil {
		t.Fatal(err)
	}

	// 创建测试文件
	testFiles := []string{
		"file1.txt",
		"file2.txt",
		"subdir/file3.txt",
		"subdir/nested/file4.txt",
	}

	for _, relPath := range testFiles {
		fullPath := filepath.Join(srcDir, relPath)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			t.Fatal(err)
		}
		content := "测试内容: " + relPath
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	// 创建ZIP文件
	c := New().WithOverwriteExisting(true)
	zipFile := filepath.Join(tempDir, "test.zip")
	err := c.Pack(zipFile, srcDir)
	if err != nil {
		t.Fatalf("创建ZIP文件失败: %v", err)
	}

	// 测试列出ZIP文件内容
	entries, err := List(zipFile)
	if err != nil {
		t.Fatalf("列出ZIP文件内容失败: %v", err)
	}

	// 验证返回的条目数量（包括目录）
	if entries == nil || len(entries.Files) == 0 {
		t.Error("应该返回文件条目，但返回为空")
	}

	// 验证包含预期的文件
	foundFiles := make(map[string]bool)
	for _, entry := range entries.Files {
		foundFiles[entry.Name] = true
		t.Logf("找到条目: %s (大小: %d, 目录: %v)", entry.Name, entry.Size, entry.IsDir)
	}

	// 检查是否包含源目录
	sourceFound := false
	for name := range foundFiles {
		if strings.Contains(name, "source") {
			sourceFound = true
			break
		}
	}
	if !sourceFound {
		t.Error("应该包含源目录，但未找到")
	}
}

// TestListTarFile 测试列出TAR文件内容
func TestListTarFile(t *testing.T) {
	tempDir := t.TempDir()

	// 创建测试文件
	srcFile := filepath.Join(tempDir, "test.txt")
	testContent := "这是TAR测试文件的内容"
	if err := os.WriteFile(srcFile, []byte(testContent), 0644); err != nil {
		t.Fatal(err)
	}

	// 创建TAR文件
	c := New().WithOverwriteExisting(true)
	tarFile := filepath.Join(tempDir, "test.tar")
	err := c.Pack(tarFile, srcFile)
	if err != nil {
		t.Fatalf("创建TAR文件失败: %v", err)
	}

	// 测试列出TAR文件内容
	entries, err := List(tarFile)
	if err != nil {
		t.Fatalf("列出TAR文件内容失败: %v", err)
	}

	// 验证返回的条目
	if entries == nil || len(entries.Files) == 0 {
		t.Error("应该返回文件条目，但返回为空")
	}

	// 验证文件信息
	found := false
	for _, entry := range entries.Files {
		if strings.Contains(entry.Name, "test.txt") {
			found = true
			if entry.IsDir {
				t.Error("test.txt 应该是文件，不是目录")
			}
			if entry.Size != int64(len(testContent)) {
				t.Errorf("文件大小不匹配: 期望 %d, 实际 %d", len(testContent), entry.Size)
			}
			t.Logf("找到文件: %s (大小: %d)", entry.Name, entry.Size)
		}
	}

	if !found {
		t.Error("应该找到 test.txt 文件")
	}
}

// TestListTgzFile 测试列出TGZ文件内容
func TestListTgzFile(t *testing.T) {
	tempDir := t.TempDir()

	// 创建测试目录和文件
	srcDir := filepath.Join(tempDir, "tgz_source")
	if err := os.MkdirAll(srcDir, 0755); err != nil {
		t.Fatal(err)
	}

	testFile := filepath.Join(srcDir, "tgz_test.txt")
	if err := os.WriteFile(testFile, []byte("TGZ测试内容"), 0644); err != nil {
		t.Fatal(err)
	}

	// 创建TGZ文件
	c := New().WithOverwriteExisting(true)
	tgzFile := filepath.Join(tempDir, "test.tgz")
	err := c.Pack(tgzFile, srcDir)
	if err != nil {
		t.Fatalf("创建TGZ文件失败: %v", err)
	}

	// 测试列出TGZ文件内容
	entries, err := List(tgzFile)
	if err != nil {
		t.Fatalf("列出TGZ文件内容失败: %v", err)
	}

	// 验证返回的条目
	if entries == nil || len(entries.Files) == 0 {
		t.Error("应该返回文件条目，但返回为空")
	}

	// 验证包含预期的文件和目录
	foundDir := false
	foundFile := false
	for _, entry := range entries.Files {
		t.Logf("TGZ条目: %s (大小: %d, 目录: %v)", entry.Name, entry.Size, entry.IsDir)
		if strings.Contains(entry.Name, "tgz_source") && entry.IsDir {
			foundDir = true
		}
		if strings.Contains(entry.Name, "tgz_test.txt") && !entry.IsDir {
			foundFile = true
		}
	}

	if !foundDir {
		t.Error("应该找到源目录")
	}
	if !foundFile {
		t.Error("应该找到测试文件")
	}
}

// TestListGzipFile 测试列出GZIP文件内容
func TestListGzipFile(t *testing.T) {
	tempDir := t.TempDir()

	// 创建测试文件
	srcFile := filepath.Join(tempDir, "gzip_test.txt")
	testContent := "这是GZIP压缩的测试内容，包含中文字符"
	if err := os.WriteFile(srcFile, []byte(testContent), 0644); err != nil {
		t.Fatal(err)
	}

	// 创建GZIP文件
	c := New().WithOverwriteExisting(true)
	gzFile := filepath.Join(tempDir, "test.txt.gz")
	err := c.Pack(gzFile, srcFile)
	if err != nil {
		t.Fatalf("创建GZIP文件失败: %v", err)
	}

	// 测试列出GZIP文件内容
	entries, err := List(gzFile)
	if err != nil {
		t.Fatalf("列出GZIP文件内容失败: %v", err)
	}

	// 验证返回的条目
	if entries == nil || len(entries.Files) != 1 {
		t.Errorf("GZIP文件应该只包含一个条目，实际包含 %d 个", entries.TotalFiles)
	}

	if len(entries.Files) > 0 {
		entry := entries.Files[0]
		if entry.IsDir {
			t.Error("GZIP条目应该是文件，不是目录")
		}

		// GZIP文件的原始文件名应该是去掉.gz后缀的名称
		expectedName := "gzip_test.txt"
		if !strings.Contains(entry.Name, expectedName) {
			t.Errorf("期望文件名包含 '%s', 实际为 '%s'", expectedName, entry.Name)
		}

		t.Logf("GZIP条目: %s (大小: %d)", entry.Name, entry.Size)
	}
}

// TestListBz2File 测试列出BZ2文件内容
func TestListBz2File(t *testing.T) {
	tempDir := t.TempDir()

	// 创建一个假的BZ2文件来测试格式检测
	bz2File := filepath.Join(tempDir, "test.bz2")
	if err := os.WriteFile(bz2File, []byte("fake bz2 content"), 0644); err != nil {
		t.Fatal(err)
	}

	// 测试列出BZ2文件内容
	entries, err := List(bz2File)
	if err != nil {
		// BZ2格式应该被支持，但可能因为文件内容不是真正的BZ2而失败
		t.Logf("BZ2文件列出失败（预期）: %v", err)
		return
	}

	// 如果成功，验证返回的条目
	t.Logf("BZ2文件包含 %d 个条目", entries.TotalFiles)
	for _, entry := range entries.Files {
		t.Logf("BZ2条目: %s (大小: %d, 目录: %v)", entry.Name, entry.Size, entry.IsDir)
	}
}

// TestListWithGlobalFunction 测试使用全局List函数
func TestListWithGlobalFunction(t *testing.T) {
	tempDir := t.TempDir()

	// 创建测试文件
	srcFile := filepath.Join(tempDir, "global_test.txt")
	if err := os.WriteFile(srcFile, []byte("全局函数测试"), 0644); err != nil {
		t.Fatal(err)
	}

	// 创建压缩文件
	zipFile := filepath.Join(tempDir, "global_test.zip")
	err := Pack(zipFile, srcFile)
	if err != nil {
		t.Fatalf("压缩失败: %v", err)
	}

	// 使用全局函数列出内容
	entries, err := List(zipFile)
	if err != nil {
		t.Fatalf("使用全局函数列出内容失败: %v", err)
	}

	// 验证结果
	if entries == nil || len(entries.Files) == 0 {
		t.Error("应该返回文件条目，但返回为空")
	}

	found := false
	for _, entry := range entries.Files {
		if strings.Contains(entry.Name, "global_test.txt") {
			found = true
			t.Logf("全局函数找到文件: %s", entry.Name)
		}
	}

	if !found {
		t.Error("应该找到测试文件")
	}
}

// TestListComplexDirectory 测试复杂目录结构
func TestListComplexDirectory(t *testing.T) {
	tempDir := t.TempDir()

	// 创建复杂的目录结构
	structure := map[string]string{
		"root.txt":                    "根目录文件",
		"dir1/file1.txt":              "目录1文件1",
		"dir1/file2.txt":              "目录1文件2",
		"dir1/subdir1/file3.txt":      "子目录1文件3",
		"dir2/file4.txt":              "目录2文件4",
		"dir2/subdir2/file5.txt":      "子目录2文件5",
		"dir2/subdir2/deep/file6.txt": "深层目录文件6",
		"empty_dir/.gitkeep":          "", // 空目录标记文件
	}

	srcDir := filepath.Join(tempDir, "complex_source")
	for relPath, content := range structure {
		fullPath := filepath.Join(srcDir, relPath)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	// 创建ZIP文件
	c := New().WithOverwriteExisting(true)
	zipFile := filepath.Join(tempDir, "complex.zip")
	err := c.Pack(zipFile, srcDir)
	if err != nil {
		t.Fatalf("创建复杂ZIP文件失败: %v", err)
	}

	// 列出内容
	entries, err := List(zipFile)
	if err != nil {
		t.Fatalf("列出复杂ZIP文件内容失败: %v", err)
	}

	// 验证条目数量（应该包括目录和文件）
	if entries == nil || len(entries.Files) == 0 {
		t.Error("复杂目录结构应该包含多个条目")
	}

	// 统计文件和目录数量
	fileCount := 0
	dirCount := 0
	for _, entry := range entries.Files {
		if entry.IsDir {
			dirCount++
		} else {
			fileCount++
		}
		t.Logf("复杂结构条目: %s (大小: %d, 目录: %v)", entry.Name, entry.Size, entry.IsDir)
	}

	t.Logf("找到 %d 个文件和 %d 个目录", fileCount, dirCount)

	// 应该至少有一些文件和目录
	if fileCount == 0 {
		t.Error("应该包含文件")
	}
}

// TestListEmptyArchive 测试空压缩文件
func TestListEmptyArchive(t *testing.T) {
	tempDir := t.TempDir()

	// 创建空目录
	emptyDir := filepath.Join(tempDir, "empty")
	if err := os.MkdirAll(emptyDir, 0755); err != nil {
		t.Fatal(err)
	}

	// 压缩空目录
	c := New().WithOverwriteExisting(true)
	zipFile := filepath.Join(tempDir, "empty.zip")
	err := c.Pack(zipFile, emptyDir)
	if err != nil {
		t.Fatalf("压缩空目录失败: %v", err)
	}

	// 列出空压缩文件内容
	entries, err := List(zipFile)
	if err != nil {
		t.Fatalf("列出空压缩文件内容失败: %v", err)
	}

	// 空目录压缩后应该至少包含目录本身
	t.Logf("空压缩文件包含 %d 个条目", entries.TotalFiles)
	for _, entry := range entries.Files {
		t.Logf("空压缩文件条目: %s (目录: %v)", entry.Name, entry.IsDir)
	}
}

// TestListFileWithSpecialCharacters 测试包含特殊字符的文件名
func TestListFileWithSpecialCharacters(t *testing.T) {
	tempDir := t.TempDir()

	// 创建包含特殊字符的文件名（在Windows上可能有限制）
	specialFiles := []string{
		"normal.txt",
		"中文文件.txt",
		"file with spaces.txt",
		"file-with-dashes.txt",
		"file_with_underscores.txt",
		"file.with.dots.txt",
	}

	srcDir := filepath.Join(tempDir, "special_chars")
	if err := os.MkdirAll(srcDir, 0755); err != nil {
		t.Fatal(err)
	}

	for _, fileName := range specialFiles {
		filePath := filepath.Join(srcDir, fileName)
		content := "内容: " + fileName
		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			t.Logf("跳过文件 %s (可能包含系统不支持的字符): %v", fileName, err)
			continue
		}
	}

	// 创建ZIP文件
	c := New().WithOverwriteExisting(true)
	zipFile := filepath.Join(tempDir, "special_chars.zip")
	err := c.Pack(zipFile, srcDir)
	if err != nil {
		t.Fatalf("压缩特殊字符文件失败: %v", err)
	}

	// 列出内容
	entries, err := List(zipFile)
	if err != nil {
		t.Fatalf("列出特殊字符文件内容失败: %v", err)
	}

	// 验证特殊字符文件名
	for _, entry := range entries.Files {
		t.Logf("特殊字符条目: %s (大小: %d)", entry.Name, entry.Size)

		// 验证文件名编码正确
		if strings.Contains(entry.Name, "中文文件") {
			t.Logf("成功处理中文文件名: %s", entry.Name)
		}
		if strings.Contains(entry.Name, " ") {
			t.Logf("成功处理包含空格的文件名: %s", entry.Name)
		}
	}
}

// TestListLimit 测试限制返回条目数量
func TestListLimit(t *testing.T) {
	tempDir := t.TempDir()

	// 创建多个文件
	srcDir := filepath.Join(tempDir, "limit_test")
	if err := os.MkdirAll(srcDir, 0755); err != nil {
		t.Fatal(err)
	}

	// 创建10个文件
	for i := 0; i < 10; i++ {
		fileName := filepath.Join(srcDir, fmt.Sprintf("file%d.txt", i))
		content := fmt.Sprintf("文件 %d 的内容", i)
		if err := os.WriteFile(fileName, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	// 创建ZIP文件
	zipFile := filepath.Join(tempDir, "limit_test.zip")
	err := Pack(zipFile, srcDir)
	if err != nil {
		t.Fatalf("创建ZIP文件失败: %v", err)
	}

	// 测试限制返回5个条目
	limit := 5
	entries, err := ListLimit(zipFile, limit)
	if err != nil {
		t.Fatalf("ListLimit失败: %v", err)
	}

	// 注意：ListLimit函数可能返回所有文件的总数，但只包含限制数量的文件信息
	// 验证返回的文件信息数量不超过限制
	if len(entries.Files) > limit {
		t.Errorf("返回的文件信息数量应该不超过 %d，实际为 %d", limit, len(entries.Files))
	}

	t.Logf("限制为 %d 的情况下返回了 %d 个文件信息", limit, len(entries.Files))
}

// TestListMatch 测试匹配模式
func TestListMatch(t *testing.T) {
	tempDir := t.TempDir()

	// 创建多种类型的文件
	srcDir := filepath.Join(tempDir, "match_test")
	if err := os.MkdirAll(srcDir, 0755); err != nil {
		t.Fatal(err)
	}

	// 创建不同类型的文件
	fileTypes := []string{
		"doc1.txt", "doc2.txt", "doc3.txt",
		"image1.jpg", "image2.jpg",
		"data1.csv", "data2.csv",
	}

	for _, fileName := range fileTypes {
		filePath := filepath.Join(srcDir, fileName)
		content := "内容: " + fileName
		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	// 创建ZIP文件
	zipFile := filepath.Join(tempDir, "match_test.zip")
	err := Pack(zipFile, srcDir)
	if err != nil {
		t.Fatalf("创建ZIP文件失败: %v", err)
	}

	// 测试匹配模式
	testCases := []struct {
		pattern       string
		expectedCount int
	}{
		{"*.txt", 3},        // 应该匹配所有txt文件
		{"*.jpg", 2},        // 应该匹配所有jpg文件
		{"doc*.txt", 3},     // 应该匹配所有doc开头的txt文件
		{"image?.jpg", 2},   // 应该匹配image后跟单个字符的jpg文件
		{"nonexistent*", 0}, // 不应该匹配任何文件
	}

	for _, tc := range testCases {
		t.Run(tc.pattern, func(t *testing.T) {
			entries, err := ListMatch(zipFile, tc.pattern)
			if err != nil {
				t.Fatalf("ListMatch失败: %v", err)
			}

			// 验证匹配的文件数量
			matchCount := len(entries.Files)
			if matchCount != tc.expectedCount {
				t.Errorf("模式 '%s' 应该匹配 %d 个文件，实际匹配了 %d 个",
					tc.pattern, tc.expectedCount, matchCount)
			}

			t.Logf("模式 '%s' 匹配了 %d 个文件", tc.pattern, matchCount)
			for _, entry := range entries.Files {
				t.Logf("  - %s", entry.Name)
			}
		})
	}
}

// TestPrintArchiveInfo 测试打印压缩包信息
func TestPrintArchiveInfo(t *testing.T) {
	tempDir := t.TempDir()

	// 创建测试文件
	srcFile := filepath.Join(tempDir, "print_test.txt")
	if err := os.WriteFile(srcFile, []byte("打印测试内容"), 0644); err != nil {
		t.Fatal(err)
	}

	// 创建ZIP文件
	zipFile := filepath.Join(tempDir, "print_test.zip")
	err := Pack(zipFile, srcFile)
	if err != nil {
		t.Fatalf("创建ZIP文件失败: %v", err)
	}

	// 测试打印压缩包信息
	err = PrintArchiveInfo(zipFile)
	if err != nil {
		t.Fatalf("PrintArchiveInfo失败: %v", err)
	}
	t.Log("PrintArchiveInfo 执行成功")
}

// TestPrintFiles 测试打印文件信息
func TestPrintFiles(t *testing.T) {
	tempDir := t.TempDir()

	// 创建测试文件
	srcFile := filepath.Join(tempDir, "print_files_test.txt")
	if err := os.WriteFile(srcFile, []byte("打印文件测试内容"), 0644); err != nil {
		t.Fatal(err)
	}

	// 创建ZIP文件
	zipFile := filepath.Join(tempDir, "print_files_test.zip")
	err := Pack(zipFile, srcFile)
	if err != nil {
		t.Fatalf("创建ZIP文件失败: %v", err)
	}

	// 测试简洁样式
	t.Log("测试简洁样式:")
	err = PrintFiles(zipFile, false)
	if err != nil {
		t.Fatalf("PrintFiles(简洁样式)失败: %v", err)
	}

	// 测试详细样式
	t.Log("测试详细样式:")
	err = PrintFiles(zipFile, true)
	if err != nil {
		t.Fatalf("PrintFiles(详细样式)失败: %v", err)
	}
}

// TestPrintFilesLimit 测试限制打印文件数量
func TestPrintFilesLimit(t *testing.T) {
	tempDir := t.TempDir()

	// 创建多个文件
	srcDir := filepath.Join(tempDir, "print_limit_test")
	if err := os.MkdirAll(srcDir, 0755); err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 5; i++ {
		fileName := filepath.Join(srcDir, fmt.Sprintf("file%d.txt", i))
		content := fmt.Sprintf("文件 %d 的内容", i)
		if err := os.WriteFile(fileName, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	// 创建ZIP文件
	zipFile := filepath.Join(tempDir, "print_limit_test.zip")
	err := Pack(zipFile, srcDir)
	if err != nil {
		t.Fatalf("创建ZIP文件失败: %v", err)
	}

	// 测试限制打印3个文件（简洁样式）
	t.Log("测试限制打印3个文件（简洁样式）:")
	err = PrintFilesLimit(zipFile, 3, false)
	if err != nil {
		t.Fatalf("PrintFilesLimit失败: %v", err)
	}

	// 测试限制打印2个文件（详细样式）
	t.Log("测试限制打印2个文件（详细样式）:")
	err = PrintFilesLimit(zipFile, 2, true)
	if err != nil {
		t.Fatalf("PrintFilesLimit失败: %v", err)
	}
}

// TestPrintFilesMatch 测试匹配模式打印
func TestPrintFilesMatch(t *testing.T) {
	tempDir := t.TempDir()

	// 创建不同类型的文件
	srcDir := filepath.Join(tempDir, "print_match_test")
	if err := os.MkdirAll(srcDir, 0755); err != nil {
		t.Fatal(err)
	}

	fileTypes := []string{"doc1.txt", "doc2.txt", "image1.jpg", "data1.csv"}
	for _, fileName := range fileTypes {
		filePath := filepath.Join(srcDir, fileName)
		content := "内容: " + fileName
		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	// 创建ZIP文件
	zipFile := filepath.Join(tempDir, "print_match_test.zip")
	err := Pack(zipFile, srcDir)
	if err != nil {
		t.Fatalf("创建ZIP文件失败: %v", err)
	}

	// 测试匹配txt文件（简洁样式）
	t.Log("测试匹配*.txt文件（简洁样式）:")
	err = PrintFilesMatch(zipFile, "*.txt", false)
	if err != nil {
		t.Fatalf("PrintFilesMatch失败: %v", err)
	}

	// 测试匹配jpg文件（详细样式）
	t.Log("测试匹配*.jpg文件（详细样式）:")
	err = PrintFilesMatch(zipFile, "*.jpg", true)
	if err != nil {
		t.Fatalf("PrintFilesMatch失败: %v", err)
	}
}

// TestPrintArchiveAndFiles 测试打印压缩包信息和文件信息
func TestPrintArchiveAndFiles(t *testing.T) {
	tempDir := t.TempDir()

	// 创建测试文件
	srcFile := filepath.Join(tempDir, "archive_files_test.txt")
	if err := os.WriteFile(srcFile, []byte("压缩包和文件信息测试"), 0644); err != nil {
		t.Fatal(err)
	}

	// 创建ZIP文件
	zipFile := filepath.Join(tempDir, "archive_files_test.zip")
	err := Pack(zipFile, srcFile)
	if err != nil {
		t.Fatalf("创建ZIP文件失败: %v", err)
	}

	// 测试简洁样式
	t.Log("测试压缩包信息+文件信息（简洁样式）:")
	err = PrintArchiveAndFiles(zipFile, false)
	if err != nil {
		t.Fatalf("PrintArchiveAndFiles(简洁样式)失败: %v", err)
	}

	// 测试详细样式
	t.Log("测试压缩包信息+文件信息（详细样式）:")
	err = PrintArchiveAndFiles(zipFile, true)
	if err != nil {
		t.Fatalf("PrintArchiveAndFiles(详细样式)失败: %v", err)
	}
}

// TestConvenienceFunctions 测试便捷函数
func TestConvenienceFunctions(t *testing.T) {
	tempDir := t.TempDir()

	// 创建测试文件
	srcDir := filepath.Join(tempDir, "convenience_test")
	if err := os.MkdirAll(srcDir, 0755); err != nil {
		t.Fatal(err)
	}

	// 创建多个文件
	for i := 0; i < 5; i++ {
		fileName := filepath.Join(srcDir, fmt.Sprintf("file%d.txt", i))
		content := fmt.Sprintf("便捷函数测试文件 %d", i)
		if err := os.WriteFile(fileName, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	// 创建ZIP文件
	zipFile := filepath.Join(tempDir, "convenience_test.zip")
	err := Pack(zipFile, srcDir)
	if err != nil {
		t.Fatalf("创建ZIP文件失败: %v", err)
	}

	// 测试PrintLs（简洁样式）
	t.Log("测试PrintLs（简洁样式）:")
	err = PrintLs(zipFile)
	if err != nil {
		t.Fatalf("PrintLs失败: %v", err)
	}

	// 测试PrintLl（详细样式）
	t.Log("测试PrintLl（详细样式）:")
	err = PrintLl(zipFile)
	if err != nil {
		t.Fatalf("PrintLl失败: %v", err)
	}

	// 测试PrintLsLimit
	t.Log("测试PrintLsLimit（限制3个文件）:")
	err = PrintLsLimit(zipFile, 3)
	if err != nil {
		t.Fatalf("PrintLsLimit失败: %v", err)
	}

	// 测试PrintLlLimit
	t.Log("测试PrintLlLimit（限制2个文件）:")
	err = PrintLlLimit(zipFile, 2)
	if err != nil {
		t.Fatalf("PrintLlLimit失败: %v", err)
	}

	// 测试PrintLsMatch
	t.Log("测试PrintLsMatch（匹配*.txt）:")
	err = PrintLsMatch(zipFile, "*.txt")
	if err != nil {
		t.Fatalf("PrintLsMatch失败: %v", err)
	}

	// 测试PrintLlMatch
	t.Log("测试PrintLlMatch（匹配file?.txt）:")
	err = PrintLlMatch(zipFile, "file?.txt")
	if err != nil {
		t.Fatalf("PrintLlMatch失败: %v", err)
	}

	// 测试PrintInfo（简洁样式）
	t.Log("测试PrintInfo（压缩包信息+文件列表，简洁样式）:")
	err = PrintInfo(zipFile)
	if err != nil {
		t.Fatalf("PrintInfo失败: %v", err)
	}

	// 测试PrintInfoDetailed（详细样式）
	t.Log("测试PrintInfoDetailed（压缩包信息+文件列表，详细样式）:")
	err = PrintInfoDetailed(zipFile)
	if err != nil {
		t.Fatalf("PrintInfoDetailed失败: %v", err)
	}

	// 测试PrintInfoLimit
	t.Log("测试PrintInfoLimit（压缩包信息+限制3个文件）:")
	err = PrintInfoLimit(zipFile, 3)
	if err != nil {
		t.Fatalf("PrintInfoLimit失败: %v", err)
	}

	// 测试PrintInfoDetailedLimit
	t.Log("测试PrintInfoDetailedLimit（压缩包信息+限制2个文件，详细样式）:")
	err = PrintInfoDetailedLimit(zipFile, 2)
	if err != nil {
		t.Fatalf("PrintInfoDetailedLimit失败: %v", err)
	}

	// 测试PrintInfoMatch
	t.Log("测试PrintInfoMatch（压缩包信息+匹配*.txt）:")
	err = PrintInfoMatch(zipFile, "*.txt")
	if err != nil {
		t.Fatalf("PrintInfoMatch失败: %v", err)
	}

	// 测试PrintInfoDetailedMatch
	t.Log("测试PrintInfoDetailedMatch（压缩包信息+匹配file?.txt，详细样式）:")
	err = PrintInfoDetailedMatch(zipFile, "file?.txt")
	if err != nil {
		t.Fatalf("PrintInfoDetailedMatch失败: %v", err)
	}
}

// TestPrintFunctionsWithDifferentFormats 测试不同格式的打印函数
func TestPrintFunctionsWithDifferentFormats(t *testing.T) {
	tempDir := t.TempDir()

	// 创建测试文件
	srcFile := filepath.Join(tempDir, "format_test.txt")
	if err := os.WriteFile(srcFile, []byte("格式测试内容"), 0644); err != nil {
		t.Fatal(err)
	}

	// 测试不同格式
	formats := []struct {
		name string
		ext  string
	}{
		{"ZIP", ".zip"},
		{"TAR", ".tar"},
		{"TGZ", ".tgz"},
		{"GZ", ".gz"},
	}

	for _, format := range formats {
		t.Run(format.name, func(t *testing.T) {
			archiveFile := filepath.Join(tempDir, "format_test"+format.ext)

			// 创建压缩文件
			err := Pack(archiveFile, srcFile)
			if err != nil {
				t.Fatalf("创建%s文件失败: %v", format.name, err)
			}

			// 测试各种打印函数
			t.Logf("测试%s格式的打印函数:", format.name)

			// PrintArchiveInfo
			err = PrintArchiveInfo(archiveFile)
			if err != nil {
				t.Errorf("PrintArchiveInfo失败: %v", err)
			}

			// PrintLs
			err = PrintLs(archiveFile)
			if err != nil {
				t.Errorf("PrintLs失败: %v", err)
			}

			// PrintLl
			err = PrintLl(archiveFile)
			if err != nil {
				t.Errorf("PrintLl失败: %v", err)
			}

			// PrintInfo
			err = PrintInfo(archiveFile)
			if err != nil {
				t.Errorf("PrintInfo失败: %v", err)
			}

			// PrintInfoDetailed
			err = PrintInfoDetailed(archiveFile)
			if err != nil {
				t.Errorf("PrintInfoDetailed失败: %v", err)
			}
		})
	}
}

// BenchmarkList 性能基准测试
func BenchmarkList(b *testing.B) {
	tempDir := b.TempDir()

	// 创建测试文件
	srcFile := filepath.Join(tempDir, "benchmark.txt")
	content := make([]byte, 1024) // 1KB测试数据
	for i := range content {
		content[i] = byte('A' + (i % 26))
	}
	if err := os.WriteFile(srcFile, content, 0644); err != nil {
		b.Fatal(err)
	}

	// 创建ZIP文件
	c := New().WithOverwriteExisting(true)
	zipFile := filepath.Join(tempDir, "benchmark.zip")
	if err := c.Pack(zipFile, srcFile); err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := List(zipFile)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkListLargeArchive 大压缩文件性能测试
func BenchmarkListLargeArchive(b *testing.B) {
	if testing.Short() {
		b.Skip("跳过大文件基准测试（使用 -short 标志）")
	}

	tempDir := b.TempDir()

	// 创建多个文件的目录结构
	srcDir := filepath.Join(tempDir, "large_source")
	if err := os.MkdirAll(srcDir, 0755); err != nil {
		b.Fatal(err)
	}

	// 创建100个小文件
	for i := 0; i < 100; i++ {
		fileName := filepath.Join(srcDir, "file_"+string(rune('0'+i%10))+".txt")
		content := "文件内容 " + string(rune('0'+i%10))
		if err := os.WriteFile(fileName, []byte(content), 0644); err != nil {
			b.Fatal(err)
		}
	}

	// 创建ZIP文件
	c := New().WithOverwriteExisting(true)
	zipFile := filepath.Join(tempDir, "large_benchmark.zip")
	if err := c.Pack(zipFile, srcDir); err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		entries, err := List(zipFile)
		if err != nil {
			b.Fatal(err)
		}
		// 确保实际处理了数据
		if entries == nil || len(entries.Files) == 0 {
			b.Fatal("应该返回条目")
		}
	}
}

// BenchmarkPrintFunctions 打印函数性能测试
func BenchmarkPrintFunctions(b *testing.B) {
	tempDir := b.TempDir()

	// 创建测试文件
	srcFile := filepath.Join(tempDir, "print_benchmark.txt")
	if err := os.WriteFile(srcFile, []byte("打印函数性能测试"), 0644); err != nil {
		b.Fatal(err)
	}

	// 创建ZIP文件
	zipFile := filepath.Join(tempDir, "print_benchmark.zip")
	if err := Pack(zipFile, srcFile); err != nil {
		b.Fatal(err)
	}

	b.Run("PrintLs", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			err := PrintLs(zipFile)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("PrintLl", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			err := PrintLl(zipFile)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("PrintInfo", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			err := PrintInfo(zipFile)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

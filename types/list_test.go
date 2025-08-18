package types

import (
	"os"
	"testing"
	"time"
)

func TestFileInfo_Creation(t *testing.T) {
	// 测试 FileInfo 结构体的创建和字段设置
	modTime := time.Now()
	fileInfo := FileInfo{
		Name:           "test.txt",
		Size:           1024,
		CompressedSize: 512,
		ModTime:        modTime,
		Mode:           0644,
		IsDir:          false,
		IsSymlink:      false,
		LinkTarget:     "",
	}

	// 验证字段值
	if fileInfo.Name != "test.txt" {
		t.Errorf("FileInfo.Name = %v, want %v", fileInfo.Name, "test.txt")
	}
	if fileInfo.Size != 1024 {
		t.Errorf("FileInfo.Size = %v, want %v", fileInfo.Size, 1024)
	}
	if fileInfo.CompressedSize != 512 {
		t.Errorf("FileInfo.CompressedSize = %v, want %v", fileInfo.CompressedSize, 512)
	}
	if !fileInfo.ModTime.Equal(modTime) {
		t.Errorf("FileInfo.ModTime = %v, want %v", fileInfo.ModTime, modTime)
	}
	if fileInfo.Mode != 0644 {
		t.Errorf("FileInfo.Mode = %v, want %v", fileInfo.Mode, 0644)
	}
	if fileInfo.IsDir != false {
		t.Errorf("FileInfo.IsDir = %v, want %v", fileInfo.IsDir, false)
	}
	if fileInfo.IsSymlink != false {
		t.Errorf("FileInfo.IsSymlink = %v, want %v", fileInfo.IsSymlink, false)
	}
}

func TestFileInfo_Directory(t *testing.T) {
	// 测试目录类型的 FileInfo
	dirInfo := FileInfo{
		Name:           "testdir/",
		Size:           0,
		CompressedSize: 0,
		ModTime:        time.Now(),
		Mode:           os.ModeDir | 0755,
		IsDir:          true,
		IsSymlink:      false,
		LinkTarget:     "",
	}

	if !dirInfo.IsDir {
		t.Error("目录的 IsDir 应该为 true")
	}
	if dirInfo.IsSymlink {
		t.Error("目录的 IsSymlink 应该为 false")
	}
	if dirInfo.Size != 0 {
		t.Error("目录的 Size 应该为 0")
	}
}

func TestFileInfo_Symlink(t *testing.T) {
	// 测试符号链接类型的 FileInfo
	symlinkInfo := FileInfo{
		Name:           "link.txt",
		Size:           0,
		CompressedSize: 0,
		ModTime:        time.Now(),
		Mode:           os.ModeSymlink | 0777,
		IsDir:          false,
		IsSymlink:      true,
		LinkTarget:     "target.txt",
	}

	if symlinkInfo.IsDir {
		t.Error("符号链接的 IsDir 应该为 false")
	}
	if !symlinkInfo.IsSymlink {
		t.Error("符号链接的 IsSymlink 应该为 true")
	}
	if symlinkInfo.LinkTarget != "target.txt" {
		t.Errorf("符号链接的 LinkTarget = %v, want %v", symlinkInfo.LinkTarget, "target.txt")
	}
}

func TestArchiveInfo_Creation(t *testing.T) {
	// 测试 ArchiveInfo 结构体的创建
	files := []FileInfo{
		{
			Name:           "file1.txt",
			Size:           100,
			CompressedSize: 80,
			ModTime:        time.Now(),
			Mode:           0644,
			IsDir:          false,
			IsSymlink:      false,
		},
		{
			Name:           "file2.txt",
			Size:           200,
			CompressedSize: 150,
			ModTime:        time.Now(),
			Mode:           0644,
			IsDir:          false,
			IsSymlink:      false,
		},
	}

	archiveInfo := ArchiveInfo{
		Type:           CompressTypeZip,
		TotalFiles:     2,
		TotalSize:      300,
		CompressedSize: 230,
		Files:          files,
	}

	// 验证字段值
	if archiveInfo.Type != CompressTypeZip {
		t.Errorf("ArchiveInfo.Type = %v, want %v", archiveInfo.Type, CompressTypeZip)
	}
	if archiveInfo.TotalFiles != 2 {
		t.Errorf("ArchiveInfo.TotalFiles = %v, want %v", archiveInfo.TotalFiles, 2)
	}
	if archiveInfo.TotalSize != 300 {
		t.Errorf("ArchiveInfo.TotalSize = %v, want %v", archiveInfo.TotalSize, 300)
	}
	if archiveInfo.CompressedSize != 230 {
		t.Errorf("ArchiveInfo.CompressedSize = %v, want %v", archiveInfo.CompressedSize, 230)
	}
	if len(archiveInfo.Files) != 2 {
		t.Errorf("ArchiveInfo.Files length = %v, want %v", len(archiveInfo.Files), 2)
	}
}

func TestArchiveInfo_EmptyArchive(t *testing.T) {
	// 测试空压缩包
	emptyArchive := ArchiveInfo{
		Type:           CompressTypeTar,
		TotalFiles:     0,
		TotalSize:      0,
		CompressedSize: 0,
		Files:          []FileInfo{},
	}

	if emptyArchive.TotalFiles != 0 {
		t.Error("空压缩包的 TotalFiles 应该为 0")
	}
	if emptyArchive.TotalSize != 0 {
		t.Error("空压缩包的 TotalSize 应该为 0")
	}
	if len(emptyArchive.Files) != 0 {
		t.Error("空压缩包的 Files 应该为空切片")
	}
}

func TestArchiveInfo_DifferentTypes(t *testing.T) {
	// 测试不同压缩格式的 ArchiveInfo
	types := []CompressType{
		CompressTypeZip,
		CompressTypeTar,
		CompressTypeTgz,
		CompressTypeTarGz,
		CompressTypeGz,
		CompressTypeBz2,
		CompressTypeBzip2,
	}

	for _, ct := range types {
		t.Run(string(ct), func(t *testing.T) {
			archiveInfo := ArchiveInfo{
				Type:           ct,
				TotalFiles:     1,
				TotalSize:      100,
				CompressedSize: 80,
				Files: []FileInfo{
					{
						Name: "test.txt",
						Size: 100,
					},
				},
			}

			if archiveInfo.Type != ct {
				t.Errorf("ArchiveInfo.Type = %v, want %v", archiveInfo.Type, ct)
			}
		})
	}
}

func TestFileInfo_ZeroValues(t *testing.T) {
	// 测试零值情况
	var fileInfo FileInfo

	if fileInfo.Name != "" {
		t.Error("零值 FileInfo.Name 应该为空字符串")
	}
	if fileInfo.Size != 0 {
		t.Error("零值 FileInfo.Size 应该为 0")
	}
	if fileInfo.CompressedSize != 0 {
		t.Error("零值 FileInfo.CompressedSize 应该为 0")
	}
	if !fileInfo.ModTime.IsZero() {
		t.Error("零值 FileInfo.ModTime 应该为零时间")
	}
	if fileInfo.Mode != 0 {
		t.Error("零值 FileInfo.Mode 应该为 0")
	}
	if fileInfo.IsDir != false {
		t.Error("零值 FileInfo.IsDir 应该为 false")
	}
	if fileInfo.IsSymlink != false {
		t.Error("零值 FileInfo.IsSymlink 应该为 false")
	}
	if fileInfo.LinkTarget != "" {
		t.Error("零值 FileInfo.LinkTarget 应该为空字符串")
	}
}

func TestArchiveInfo_ZeroValues(t *testing.T) {
	// 测试零值情况
	var archiveInfo ArchiveInfo

	if archiveInfo.Type != "" {
		t.Error("零值 ArchiveInfo.Type 应该为空字符串")
	}
	if archiveInfo.TotalFiles != 0 {
		t.Error("零值 ArchiveInfo.TotalFiles 应该为 0")
	}
	if archiveInfo.TotalSize != 0 {
		t.Error("零值 ArchiveInfo.TotalSize 应该为 0")
	}
	if archiveInfo.CompressedSize != 0 {
		t.Error("零值 ArchiveInfo.CompressedSize 应该为 0")
	}
	if archiveInfo.Files != nil {
		t.Error("零值 ArchiveInfo.Files 应该为 nil")
	}
}

// 基准测试
func BenchmarkFileInfo_Creation(b *testing.B) {
	modTime := time.Now()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = FileInfo{
			Name:           "test.txt",
			Size:           1024,
			CompressedSize: 512,
			ModTime:        modTime,
			Mode:           0644,
			IsDir:          false,
			IsSymlink:      false,
			LinkTarget:     "",
		}
	}
}

func BenchmarkArchiveInfo_Creation(b *testing.B) {
	files := make([]FileInfo, 100)
	for i := range files {
		files[i] = FileInfo{
			Name: "test.txt",
			Size: 1024,
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ArchiveInfo{
			Type:           CompressTypeZip,
			TotalFiles:     100,
			TotalSize:      102400,
			CompressedSize: 51200,
			Files:          files,
		}
	}
}

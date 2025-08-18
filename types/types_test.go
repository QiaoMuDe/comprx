package types

import (
	"testing"
)

func TestCompressType_String(t *testing.T) {
	tests := []struct {
		name     string
		ct       CompressType
		expected string
	}{
		{"ZIP格式", CompressTypeZip, ".zip"},
		{"TAR格式", CompressTypeTar, ".tar"},
		{"TGZ格式", CompressTypeTgz, ".tgz"},
		{"TAR.GZ格式", CompressTypeTarGz, ".tar.gz"},
		{"GZ格式", CompressTypeGz, ".gz"},
		{"BZ2格式", CompressTypeBz2, ".bz2"},
		{"BZIP2格式", CompressTypeBzip2, ".bzip2"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ct.String(); got != tt.expected {
				t.Errorf("CompressType.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestIsSupportedCompressType(t *testing.T) {
	tests := []struct {
		name     string
		ct       string
		expected bool
	}{
		{"支持的ZIP格式", ".zip", true},
		{"支持的TAR格式", ".tar", true},
		{"支持的TGZ格式", ".tgz", true},
		{"支持的TAR.GZ格式", ".tar.gz", true},
		{"支持的GZ格式", ".gz", true},
		{"支持的BZ2格式", ".bz2", true},
		{"支持的BZIP2格式", ".bzip2", true},
		{"不支持的RAR格式", ".rar", false},
		{"不支持的7Z格式", ".7z", false},
		{"空字符串", "", false},
		{"无效格式", ".invalid", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsSupportedCompressType(tt.ct); got != tt.expected {
				t.Errorf("IsSupportedCompressType() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestSupportedCompressTypes(t *testing.T) {
	supportedTypes := SupportedCompressTypes()

	// 检查返回的切片不为空
	if len(supportedTypes) == 0 {
		t.Error("SupportedCompressTypes() 返回空切片")
	}

	// 检查是否包含所有预期的格式
	expectedTypes := []string{".zip", ".tar", ".tgz", ".tar.gz", ".gz", ".bz2", ".bzip2"}

	for _, expected := range expectedTypes {
		found := false
		for _, actual := range supportedTypes {
			if actual == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("SupportedCompressTypes() 缺少格式: %s", expected)
		}
	}

	// 检查返回的数量是否正确
	if len(supportedTypes) != len(expectedTypes) {
		t.Errorf("SupportedCompressTypes() 返回数量 = %d, 期望 %d", len(supportedTypes), len(expectedTypes))
	}
}

func TestDetectCompressFormat(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		expected CompressType
		wantErr  bool
	}{
		{"ZIP文件", "test.zip", CompressTypeZip, false},
		{"TAR文件", "test.tar", CompressTypeTar, false},
		{"TGZ文件", "test.tgz", CompressTypeTgz, false},
		{"TAR.GZ文件", "test.tar.gz", CompressTypeTarGz, false},
		{"GZ文件", "test.gz", CompressTypeGz, false},
		{"BZ2文件", "test.bz2", CompressTypeBz2, false},
		{"BZIP2文件", "test.bzip2", CompressTypeBzip2, false},
		{"大写扩展名", "TEST.ZIP", CompressTypeZip, false},
		{"混合大小写", "Test.Tar.Gz", CompressTypeTarGz, false},
		{"带路径的文件", "/path/to/file.zip", CompressTypeZip, false},
		{"Windows路径", "C:\\path\\to\\file.tar", CompressTypeTar, false},
		{"不支持的格式", "test.rar", "", true},
		{"无扩展名", "test", "", true},
		{"空文件名", "", "", true},
		{"只有扩展名", ".zip", CompressTypeZip, false},
		{"多个扩展名", "test.backup.zip", CompressTypeZip, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DetectCompressFormat(tt.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("DetectCompressFormat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.expected {
				t.Errorf("DetectCompressFormat() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestDetectCompressFormat_SpecialCases(t *testing.T) {
	// 测试特殊的 .tar.gz 处理
	t.Run("TAR.GZ优先级测试", func(t *testing.T) {
		// .tar.gz 应该被识别为 TAR.GZ 而不是 GZ
		got, err := DetectCompressFormat("test.tar.gz")
		if err != nil {
			t.Errorf("DetectCompressFormat() 意外错误 = %v", err)
		}
		if got != CompressTypeTarGz {
			t.Errorf("DetectCompressFormat() = %v, want %v", got, CompressTypeTarGz)
		}
	})

	t.Run("复杂文件名测试", func(t *testing.T) {
		testCases := []struct {
			filename string
			expected CompressType
		}{
			{"backup-2023-01-01.tar.gz", CompressTypeTarGz},
			{"data.backup.tar.gz", CompressTypeTarGz},
			{"file.TAR.GZ", CompressTypeTarGz},
			{"archive.Tar.Gz", CompressTypeTarGz},
		}

		for _, tc := range testCases {
			got, err := DetectCompressFormat(tc.filename)
			if err != nil {
				t.Errorf("DetectCompressFormat(%s) 意外错误 = %v", tc.filename, err)
			}
			if got != tc.expected {
				t.Errorf("DetectCompressFormat(%s) = %v, want %v", tc.filename, got, tc.expected)
			}
		}
	})
}

// 基准测试
func BenchmarkDetectCompressFormat(b *testing.B) {
	testFiles := []string{
		"test.zip",
		"test.tar",
		"test.tar.gz",
		"test.tgz",
		"test.gz",
		"test.bz2",
		"test.bzip2",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, file := range testFiles {
			_, _ = DetectCompressFormat(file)
		}
	}
}

func BenchmarkIsSupportedCompressType(b *testing.B) {
	testTypes := []string{".zip", ".tar", ".gz", ".rar", ".7z"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, ct := range testTypes {
			_ = IsSupportedCompressType(ct)
		}
	}
}

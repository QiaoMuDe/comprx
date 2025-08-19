package utils

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestValidatePathSimple_SafePaths(t *testing.T) {
	targetDir := "/tmp/extract"

	testCases := []struct {
		name     string
		filePath string
		expected string
	}{
		{
			name:     "普通文件名",
			filePath: "file.txt",
			expected: "/tmp/extract/file.txt",
		},
		{
			name:     "子目录中的文件",
			filePath: "subdir/file.txt",
			expected: "/tmp/extract/subdir/file.txt",
		},
		{
			name:     "多层目录",
			filePath: "dir1/dir2/dir3/file.txt",
			expected: "/tmp/extract/dir1/dir2/dir3/file.txt",
		},
		{
			name:     "带点的文件名",
			filePath: "file.name.txt",
			expected: "/tmp/extract/file.name.txt",
		},
		{
			name:     "空文件路径",
			filePath: "",
			expected: "/tmp/extract",
		},
		{
			name:     "当前目录标记（单独的点）",
			filePath: ".",
			expected: "/tmp/extract",
		},
		{
			name:     "包含空格的文件名",
			filePath: "my file.txt",
			expected: "/tmp/extract/my file.txt",
		},
		{
			name:     "包含特殊字符的文件名",
			filePath: "file-name_123.txt",
			expected: "/tmp/extract/file-name_123.txt",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := ValidatePathSimple(targetDir, tc.filePath)
			if err != nil {
				t.Errorf("安全路径 '%s' 被错误拒绝: %v", tc.filePath, err)
				return
			}

			// 标准化路径比较（处理不同操作系统的路径分隔符）
			expectedNorm := filepath.FromSlash(tc.expected)
			if result != expectedNorm {
				t.Errorf("路径结果不匹配:\n期望: %s\n实际: %s", expectedNorm, result)
			}
		})
	}
}

func TestValidatePathSimple_DangerousPaths(t *testing.T) {
	targetDir := "/tmp/extract"

	testCases := []struct {
		name     string
		filePath string
		reason   string
	}{
		{
			name:     "基本路径遍历",
			filePath: "../../../etc/passwd",
			reason:   "包含上级目录引用",
		},
		{
			name:     "混合路径遍历",
			filePath: "normal/../../../etc/passwd",
			reason:   "包含上级目录引用",
		},
		{
			name:     "Windows绝对路径",
			filePath: "C:\\Windows\\system32\\config\\sam",
			reason:   "绝对路径（Windows）",
		},
		{
			name:     "复杂路径遍历",
			filePath: "dir1/dir2/../../../etc/passwd",
			reason:   "包含上级目录引用",
		},
		{
			name:     "隐藏的路径遍历",
			filePath: "normal/./../../etc/passwd",
			reason:   "包含上级目录引用",
		},
		{
			name:     "多个连续的上级目录",
			filePath: "../../../../../../../../etc/passwd",
			reason:   "包含上级目录引用",
		},
		{
			name:     "路径末尾的上级目录",
			filePath: "somedir/..",
			reason:   "包含上级目录引用",
		},
		{
			name:     "UNC路径（Windows）",
			filePath: "\\\\server\\share\\file.txt",
			reason:   "UNC路径",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := ValidatePathSimple(targetDir, tc.filePath)
			if err == nil {
				t.Errorf("危险路径 '%s' 未被拒绝，返回结果: %s (原因: %s)", tc.filePath, result, tc.reason)
				return
			}

			// 验证错误消息包含预期内容
			if !strings.Contains(err.Error(), "不安全的路径") && !strings.Contains(err.Error(), "可疑路径") {
				t.Errorf("错误消息格式不正确: %v", err)
			}
		})
	}
}

func TestValidatePathSimple_SuspiciousPaths(t *testing.T) {
	targetDir := "/tmp/extract"

	// 这些路径可能被标记为可疑，取决于操作系统
	testCases := []struct {
		name     string
		filePath string
		reason   string
	}{
		{
			name:     "路径分隔符后跟点",
			filePath: "dir/.hidden",
			reason:   "可能的隐藏文件路径",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := ValidatePathSimple(targetDir, tc.filePath)

			// 在某些系统上这可能被拒绝，在其他系统上可能被接受
			if err != nil {
				t.Logf("路径 '%s' 被拒绝: %v (原因: %s)", tc.filePath, err, tc.reason)
			} else {
				t.Logf("路径 '%s' 被接受: %s (原因: %s)", tc.filePath, result, tc.reason)
			}
		})
	}
}

func TestValidatePathSimple_EdgeCases(t *testing.T) {
	testCases := []struct {
		name      string
		targetDir string
		filePath  string
		expectErr bool
		reason    string
	}{
		{
			name:      "空目标目录",
			targetDir: "",
			filePath:  "file.txt",
			expectErr: false,
			reason:    "应该能处理空目标目录",
		},
		{
			name:      "相对目标目录",
			targetDir: "relative/path",
			filePath:  "file.txt",
			expectErr: false,
			reason:    "应该能处理相对目标目录",
		},
		{
			name:      "目标目录和文件路径都为空",
			targetDir: "",
			filePath:  "",
			expectErr: false,
			reason:    "应该能处理都为空的情况",
		},
		{
			name:      "包含Unicode字符的文件名",
			targetDir: "/tmp/extract",
			filePath:  "文件名.txt",
			expectErr: false,
			reason:    "应该支持Unicode文件名",
		},
		{
			name:      "很长的文件名",
			targetDir: "/tmp/extract",
			filePath:  strings.Repeat("a", 255) + ".txt",
			expectErr: false,
			reason:    "应该能处理长文件名",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := ValidatePathSimple(tc.targetDir, tc.filePath)

			if tc.expectErr && err == nil {
				t.Errorf("期望返回错误但没有返回，结果: %s (原因: %s)", result, tc.reason)
			} else if !tc.expectErr && err != nil {
				t.Errorf("不期望返回错误但返回了: %v (原因: %s)", err, tc.reason)
			} else {
				t.Logf("测试通过: %s -> %s (原因: %s)", tc.filePath, result, tc.reason)
			}
		})
	}
}

func TestValidatePathSimple_CrossPlatform(t *testing.T) {
	targetDir := "/tmp/extract"

	// 测试不同操作系统的路径分隔符
	testCases := []struct {
		name      string
		filePath  string
		expectErr bool
	}{
		{
			name:      "Unix风格路径分隔符",
			filePath:  "dir/subdir/file.txt",
			expectErr: false,
		},
		{
			name:      "Windows风格路径分隔符",
			filePath:  "dir\\subdir\\file.txt",
			expectErr: false,
		},
		{
			name:      "混合路径分隔符",
			filePath:  "dir/subdir\\file.txt",
			expectErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := ValidatePathSimple(targetDir, tc.filePath)

			if tc.expectErr && err == nil {
				t.Errorf("期望返回错误但没有返回，结果: %s", result)
			} else if !tc.expectErr && err != nil {
				t.Errorf("不期望返回错误但返回了: %v", err)
			}

			// 验证结果路径使用了正确的分隔符
			if err == nil {
				// 标准化路径进行比较
				normalizedTargetDir := filepath.FromSlash(targetDir)
				if !strings.Contains(result, normalizedTargetDir) {
					t.Errorf("结果路径不包含目标目录: %s (期望包含: %s)", result, normalizedTargetDir)
				}
			}
		})
	}
}

func TestValidatePathSimple_RealWorldAttacks(t *testing.T) {
	targetDir := "/var/www/uploads"

	// 基于真实攻击案例的测试 - 专注于路径遍历技术而非特定系统路径
	realAttacks := []string{
		"../../../etc/passwd",
		"..\\..\\..\\windows\\system32\\config\\sam",
		"....//....//....//etc/passwd",
		"..%2F..%2F..%2Fetc%2Fpasswd",
		"../../../proc/self/environ",
		"../../../home/user/.ssh/id_rsa",
		"C:\\boot.ini",
		"\\\\?\\C:\\Windows\\system32\\drivers\\etc\\hosts",
		"file:///etc/passwd",
		"../../../var/log/auth.log",
	}

	for i, attack := range realAttacks {
		t.Run(fmt.Sprintf("RealAttack_%d", i+1), func(t *testing.T) {
			result, err := ValidatePathSimple(targetDir, attack)
			if err == nil {
				t.Errorf("真实攻击路径 '%s' 未被阻止，返回: %s", attack, result)
			} else {
				t.Logf("成功阻止攻击: %s -> %v", attack, err)
			}
		})
	}
}

func BenchmarkValidatePathSimple_SafePath(b *testing.B) {
	targetDir := "/tmp/extract"
	filePath := "normal/path/to/file.txt"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ValidatePathSimple(targetDir, filePath)
	}
}

func BenchmarkValidatePathSimple_DangerousPath(b *testing.B) {
	targetDir := "/tmp/extract"
	filePath := "../../../etc/passwd"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ValidatePathSimple(targetDir, filePath)
	}
}

func BenchmarkValidatePathSimple_ComplexPath(b *testing.B) {
	targetDir := "/tmp/extract"
	filePath := "very/deep/directory/structure/with/many/levels/file.txt"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ValidatePathSimple(targetDir, filePath)
	}
}

// TestValidatePathSimple_WindowsSpecific 测试Windows特定的路径问题
func TestValidatePathSimple_WindowsSpecific(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("跳过Windows特定测试")
	}

	targetDir := "C:\\temp\\extract"

	testCases := []struct {
		name      string
		filePath  string
		expectErr bool
	}{
		{
			name:      "Windows保留设备名",
			filePath:  "CON",
			expectErr: false, // 当前实现不检查保留名
		},
		{
			name:      "Windows保留设备名带扩展名",
			filePath:  "PRN.txt",
			expectErr: false, // 当前实现不检查保留名
		},
		{
			name:      "Windows UNC路径",
			filePath:  "\\\\server\\share\\file.txt",
			expectErr: true, // 应该被拒绝
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := ValidatePathSimple(targetDir, tc.filePath)

			if tc.expectErr && err == nil {
				t.Errorf("期望返回错误但没有返回，结果: %s", result)
			} else if !tc.expectErr && err != nil {
				t.Errorf("不期望返回错误但返回了: %v", err)
			}
		})
	}
}

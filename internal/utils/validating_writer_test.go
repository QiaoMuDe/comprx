package utils

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"gitee.com/MM-Q/comprx/config"
)

func TestCompressionValidatingWriter(t *testing.T) {
	tests := []struct {
		name          string
		config        *config.Config
		writeData     []string
		expectError   bool
		errorContains string
		expectedTotal int64
	}{
		{
			name: "正常写入小数据",
			config: &config.Config{
				EnableSizeCheck: true,
				MaxFileSize:     1024,
				MaxTotalSize:    2048,
			},
			writeData:     []string{"hello", "world"},
			expectError:   false,
			expectedTotal: 10,
		},
		{
			name: "超过单文件大小限制",
			config: &config.Config{
				EnableSizeCheck: true,
				MaxFileSize:     8,
				MaxTotalSize:    2048,
			},
			writeData:     []string{"hello", "world"},
			expectError:   true,
			errorContains: "压缩后文件大小",
		},
		{
			name: "超过总大小限制",
			config: &config.Config{
				EnableSizeCheck: true,
				MaxFileSize:     1024,
				MaxTotalSize:    8,
			},
			writeData:     []string{"hello", "world"},
			expectError:   true,
			errorContains: "压缩后文件大小",
		},
		{
			name: "禁用大小检查时不进行验证",
			config: &config.Config{
				EnableSizeCheck: false,
				MaxFileSize:     5,
				MaxTotalSize:    5,
			},
			writeData:     []string{"hello", "world", "this is long"},
			expectError:   false,
			expectedTotal: 22,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			writer := NewCompressionValidatingWriter(&buf, tt.config)

			var totalWritten int64
			var err error

			for _, data := range tt.writeData {
				n, writeErr := writer.Write([]byte(data))
				if writeErr != nil {
					err = writeErr
					break
				}
				totalWritten += int64(n)
			}

			if tt.expectError {
				if err == nil {
					t.Errorf("期望出现错误，但没有错误发生")
					return
				}
				if !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("错误信息不匹配，期望包含 '%s'，实际错误: %v", tt.errorContains, err)
				}
			} else {
				if err != nil {
					t.Errorf("不期望出现错误，但发生了错误: %v", err)
					return
				}
				if writer.GetTotalWritten() != tt.expectedTotal {
					t.Errorf("总写入大小不匹配，期望 %d，实际 %d", tt.expectedTotal, writer.GetTotalWritten())
				}
			}
		})
	}
}

func TestDecompressionValidatingWriter(t *testing.T) {
	tests := []struct {
		name          string
		config        *config.Config
		writeData     []string
		expectError   bool
		errorContains string
		expectedTotal int64
	}{
		{
			name: "正常解压写入",
			config: &config.Config{
				EnableSizeCheck: true,
				MaxFileSize:     1024,
				MaxTotalSize:    2048,
			},
			writeData:     []string{"hello", "world"},
			expectError:   false,
			expectedTotal: 10,
		},
		{
			name: "解压超过大小限制",
			config: &config.Config{
				EnableSizeCheck: true,
				MaxFileSize:     8,
				MaxTotalSize:    2048,
			},
			writeData:     []string{"hello", "world"},
			expectError:   true,
			errorContains: "解压后文件大小",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			tracker := NewSizeTracker()
			writer := NewDecompressionValidatingWriter(&buf, tt.config, 100, tracker)

			var err error
			for _, data := range tt.writeData {
				_, writeErr := writer.Write([]byte(data))
				if writeErr != nil {
					err = writeErr
					break
				}
			}

			if tt.expectError {
				if err == nil {
					t.Errorf("期望出现错误，但没有错误发生")
					return
				}
				if !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("错误信息不匹配，期望包含 '%s'，实际错误: %v", tt.errorContains, err)
				}
			} else {
				if err != nil {
					t.Errorf("不期望出现错误，但发生了错误: %v", err)
					return
				}
				if writer.GetTotalWritten() != tt.expectedTotal {
					t.Errorf("总写入大小不匹配，期望 %d，实际 %d", tt.expectedTotal, writer.GetTotalWritten())
				}
			}

			// 验证大小跟踪器是否正确设置
			if writer.GetSizeTracker() != tracker {
				t.Error("大小跟踪器未正确设置")
			}
		})
	}
}

func TestValidatingWriter_GetTotalWritten(t *testing.T) {
	var buf bytes.Buffer
	config := &config.Config{
		EnableSizeCheck: true,
		MaxFileSize:     1024,
		MaxTotalSize:    2048,
	}
	writer := NewCompressionValidatingWriter(&buf, config)

	// 初始状态应该为0
	if writer.GetTotalWritten() != 0 {
		t.Errorf("初始总写入大小应该为0，实际为 %d", writer.GetTotalWritten())
	}

	// 写入一些数据
	data1 := []byte("hello")
	n1, err := writer.Write(data1)
	if err != nil {
		t.Fatalf("写入数据失败: %v", err)
	}
	if writer.GetTotalWritten() != int64(n1) {
		t.Errorf("第一次写入后总大小不匹配，期望 %d，实际 %d", n1, writer.GetTotalWritten())
	}

	// 再写入一些数据
	data2 := []byte("world")
	n2, err := writer.Write(data2)
	if err != nil {
		t.Fatalf("写入数据失败: %v", err)
	}
	expectedTotal := int64(n1 + n2)
	if writer.GetTotalWritten() != expectedTotal {
		t.Errorf("第二次写入后总大小不匹配，期望 %d，实际 %d", expectedTotal, writer.GetTotalWritten())
	}
}

func TestValidatingWriter_DisabledSizeCheck(t *testing.T) {
	var buf bytes.Buffer
	config := &config.Config{
		EnableSizeCheck: false,
		MaxFileSize:     1,
		MaxTotalSize:    1,
	}
	writer := NewCompressionValidatingWriter(&buf, config)

	// 写入大量数据，应该不会触发大小检查
	largeData := strings.Repeat("a", 10000)
	n, err := writer.Write([]byte(largeData))
	if err != nil {
		t.Errorf("禁用大小检查时不应该出现大小相关错误，但发生了错误: %v", err)
	}
	if n != len(largeData) {
		t.Errorf("写入字节数不匹配，期望 %d，实际 %d", len(largeData), n)
	}
	if writer.GetTotalWritten() != int64(len(largeData)) {
		t.Errorf("总写入大小不匹配，期望 %d，实际 %d", len(largeData), writer.GetTotalWritten())
	}

	// 验证数据确实被写入到底层写入器
	if buf.Len() != len(largeData) {
		t.Errorf("底层缓冲区大小不匹配，期望 %d，实际 %d", len(largeData), buf.Len())
	}
}

func TestNewCompressionValidatingWriter(t *testing.T) {
	var buf bytes.Buffer
	config := &config.Config{
		EnableSizeCheck: true,
		MaxFileSize:     1024,
		MaxTotalSize:    2048,
	}

	writer := NewCompressionValidatingWriter(&buf, config)

	if writer == nil {
		t.Fatal("NewCompressionValidatingWriter 返回了 nil")
	}
	if writer.errorPrefix != "压缩后" {
		t.Errorf("压缩验证写入器错误前缀应该为'压缩后'，实际为 '%s'", writer.errorPrefix)
	}
	if writer.GetTotalWritten() != 0 {
		t.Errorf("初始总写入大小应该为0，实际为 %d", writer.GetTotalWritten())
	}
}

func TestNewDecompressionValidatingWriter(t *testing.T) {
	var buf bytes.Buffer
	config := &config.Config{
		EnableSizeCheck: true,
		MaxFileSize:     1024,
		MaxTotalSize:    2048,
	}
	tracker := NewSizeTracker()

	writer := NewDecompressionValidatingWriter(&buf, config, 100, tracker)

	if writer == nil {
		t.Fatal("NewDecompressionValidatingWriter 返回了 nil")
	}
	if writer.errorPrefix != "解压后" {
		t.Errorf("解压验证写入器错误前缀应该为'解压后'，实际为 '%s'", writer.errorPrefix)
	}
	if writer.GetSizeTracker() != tracker {
		t.Error("大小跟踪器未正确设置")
	}
}

// testErrorWriter 是一个总是返回错误的写入器，用于测试
type testErrorWriter struct{}

func (ew *testErrorWriter) Write(p []byte) (n int, err error) {
	return 0, errors.New("write error")
}

func TestValidatingWriter_WriteError(t *testing.T) {
	// 创建一个会产生错误的写入器
	errorWriter := &testErrorWriter{}
	config := &config.Config{
		EnableSizeCheck: true,
		MaxFileSize:     1024,
		MaxTotalSize:    2048,
	}
	writer := NewCompressionValidatingWriter(errorWriter, config)

	// 尝试写入数据，应该返回底层写入器的错误
	_, err := writer.Write([]byte("test"))
	if err == nil {
		t.Error("期望底层写入器错误被传播，但没有错误发生")
	}
	if err.Error() != "write error" {
		t.Errorf("期望错误信息为 'write error'，实际为 '%s'", err.Error())
	}
}

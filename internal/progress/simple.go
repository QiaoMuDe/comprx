package progress

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Simple 简单进度显示器
type Simple struct {
	writer  io.Writer
	enabled bool
}

// NewSimple 创建简单进度显示器
func NewSimple(enabled bool) *Simple {
	return &Simple{
		writer:  os.Stdout,
		enabled: enabled,
	}
}

// WithWriter 设置输出位置
func (s *Simple) WithWriter(writer io.Writer) *Simple {
	if writer != nil {
		s.writer = writer
	}
	return s
}

// Archive 显示压缩文件信息
func (s *Simple) Archive(archivePath string) {
	if s.enabled {
		fmt.Fprintf(s.writer, "Archive: %s\n", filepath.Base(archivePath))
	}
}

// Inflating 显示解压文件
func (s *Simple) Inflating(filePath string) {
	if s.enabled {
		fmt.Fprintf(s.writer, "  inflating: %s\n", filePath)
	}
}

// Creating 显示创建目录
func (s *Simple) Creating(dirPath string) {
	if s.enabled {
		fmt.Fprintf(s.writer, "   creating: %s\n", dirPath)
	}
}

// Extracting 显示提取文件（TAR）
func (s *Simple) Extracting(filePath string) {
	if s.enabled {
		fmt.Fprintf(s.writer, "extracting: %s\n", filePath)
	}
}

// Adding 显示添加文件
func (s *Simple) Adding(filePath string) {
	if s.enabled {
		fmt.Fprintf(s.writer, "     adding: %s\n", filePath)
	}
}

// Storing 显示存储目录
func (s *Simple) Storing(dirPath string) {
	if s.enabled {
		fmt.Fprintf(s.writer, "    storing: %s\n", dirPath)
	}
}

// Compressing 显示压缩文件
func (s *Simple) Compressing(filePath string) {
	if s.enabled {
		fmt.Fprintf(s.writer, "compressing: %s\n", filePath)
	}
}

// File 通用文件处理显示（自动判断是解压还是创建目录）
func (s *Simple) File(filePath string, isDirectory bool) {
	if !s.enabled {
		return
	}
	
	if isDirectory {
		s.Creating(filePath)
	} else {
		s.Inflating(filePath)
	}
}

// IsEnabled 检查是否启用
func (s *Simple) IsEnabled() bool {
	return s.enabled
}
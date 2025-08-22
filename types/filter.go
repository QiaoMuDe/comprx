package types

import (
	"fmt"
	"os"
	"path/filepath"
)

// FileFilter 文件过滤器接口
//
// 用于判断文件是否应该被包含在压缩/解压操作中
type FileFilter interface {
	// ShouldInclude 判断文件是否应该被包含
	//
	// 参数:
	//   - path: 文件路径
	//   - info: 文件信息
	//
	// 返回:
	//   - bool: true 表示应该包含，false 表示应该排除
	ShouldInclude(path string, info os.FileInfo) bool
}

// FilterOptions 过滤配置选项
//
// 用于指定压缩时的文件过滤条件：
//   - Include: 包含模式列表，支持 glob 语法
//   - Exclude: 排除模式列表，支持 glob 语法
//   - MaxSize: 最大文件大小限制（字节），0 表示无限制
//   - MinSize: 最小文件大小限制（字节），默认为 0
type FilterOptions struct {
	Include []string // 包含模式，支持 glob 语法
	Exclude []string // 排除模式，支持 glob 语法
	MaxSize int64    // 最大文件大小（字节），0 表示无限制
	MinSize int64    // 最小文件大小（字节），默认为 0
}

// ShouldInclude 判断文件是否应该被包含
//
// 过滤逻辑:
//  1. 检查文件大小是否符合要求
//  2. 如果指定了包含模式，检查文件是否匹配包含模式
//  3. 检查文件是否匹配排除模式
//
// 参数:
//   - path: 文件路径
//   - info: 文件信息
//
// 返回:
//   - bool: true 表示应该包含，false 表示应该排除
func (f *FilterOptions) ShouldInclude(path string, info os.FileInfo) bool {
	// 如果过滤器为空，包含所有文件
	if f == nil {
		return true
	}

	// 检查是否有任何过滤条件
	hasFilter := len(f.Include) > 0 || 
		len(f.Exclude) > 0 || 
		f.MinSize > 0 || 
		f.MaxSize > 0

	// 如果没有过滤条件，包含所有文件
	if !hasFilter {
		return true
	}

	// 1. 检查文件大小
	if !f.checkSizeFilter(info) {
		return false
	}

	// 2. 检查包含模式（如果指定了包含模式）
	if len(f.Include) > 0 {
		if !f.matchAnyPattern(f.Include, path) {
			return false
		}
	}

	// 3. 检查排除模式
	if len(f.Exclude) > 0 {
		if f.matchAnyPattern(f.Exclude, path) {
			return false
		}
	}

	return true
}

// checkSizeFilter 检查文件大小是否符合过滤条件
//
// 参数:
//   - info: 文件信息
//
// 返回:
//   - bool: true 表示通过大小过滤，false 表示不通过
func (f *FilterOptions) checkSizeFilter(info os.FileInfo) bool {
	if info.IsDir() {
		return true // 目录总是通过大小过滤
	}

	size := info.Size()

	// 检查最小大小
	if f.MinSize > 0 && size < f.MinSize {
		return false
	}

	// 检查最大大小
	if f.MaxSize > 0 && size > f.MaxSize {
		return false
	}

	return true
}

// matchAnyPattern 检查路径是否匹配任一模式
//
// 参数:
//   - patterns: 模式列表
//   - path: 文件路径
//
// 返回:
//   - bool: true 表示匹配任一模式，false 表示不匹配任何模式
func (f *FilterOptions) matchAnyPattern(patterns []string, path string) bool {
	for _, pattern := range patterns {
		if f.matchPattern(pattern, path) {
			return true
		}
	}
	return false
}

// matchPattern 检查路径是否匹配指定模式
//
// 支持多种匹配方式:
//  1. 文件名匹配（如 *.go 匹配 main.go）
//  2. 完整路径匹配（如 src/*.go 匹配 src/main.go）
//  3. 目录匹配（如 vendor/ 匹配 vendor 目录）
//
// 参数:
//   - pattern: glob 模式
//   - path: 文件路径
//
// 返回:
//   - bool: true 表示匹配，false 表示不匹配
func (f *FilterOptions) matchPattern(pattern, path string) bool {
	// 1. 尝试匹配文件名
	if matched, err := filepath.Match(pattern, filepath.Base(path)); err == nil && matched {
		return true
	}

	// 2. 尝试匹配完整路径
	if matched, err := filepath.Match(pattern, path); err == nil && matched {
		return true
	}

	// 3. 尝试匹配目录（处理以 / 结尾的模式）
	if len(pattern) > 0 && (pattern[len(pattern)-1] == '/' || pattern[len(pattern)-1] == '\\') {
		dirPattern := pattern[:len(pattern)-1]
		if matched, err := filepath.Match(dirPattern, filepath.Base(path)); err == nil && matched {
			return true
		}
		if matched, err := filepath.Match(dirPattern, path); err == nil && matched {
			return true
		}
	}

	// 4. 处理路径中包含模式的情况
	pathParts := filepath.SplitList(filepath.ToSlash(path))
	for _, part := range pathParts {
		if matched, err := filepath.Match(pattern, part); err == nil && matched {
			return true
		}
	}

	return false
}

// Validate 验证过滤器选项
//
// 返回:
//   - error: 验证错误，如果验证通过则返回 nil
func (f *FilterOptions) Validate() error {
	// 验证文件大小范围
	if f.MinSize < 0 {
		return fmt.Errorf("最小文件大小不能为负数: %d", f.MinSize)
	}
	
	if f.MaxSize < 0 {
		return fmt.Errorf("最大文件大小不能为负数: %d", f.MaxSize)
	}
	
	if f.MinSize > 0 && f.MaxSize > 0 && f.MinSize > f.MaxSize {
		return fmt.Errorf("最小文件大小 (%d) 不能大于最大文件大小 (%d)", f.MinSize, f.MaxSize)
	}

	// 验证包含模式
	for _, pattern := range f.Include {
		if pattern == "" {
			return fmt.Errorf("包含模式不能为空字符串")
		}
	}

	// 验证排除模式
	for _, pattern := range f.Exclude {
		if pattern == "" {
			return fmt.Errorf("排除模式不能为空字符串")
		}
	}

	return nil
}

// HasFilterConditions 检查过滤器是否有任何过滤条件
//
// 参数:
//   - filter: 过滤配置选项
//
// 返回:
//   - bool: true 表示有过滤条件，false 表示没有
func HasFilterConditions(filter *FilterOptions) bool {
	if filter == nil {
		return false
	}
	
	return len(filter.Include) > 0 || 
		len(filter.Exclude) > 0 || 
		filter.MinSize > 0 || 
		filter.MaxSize > 0
}

// LoadExcludeFromFile 从忽略文件加载排除模式
//
// 参数:
//   - ignoreFilePath: 忽略文件路径
//
// 返回:
//   - []string: 排除模式列表
//   - error: 错误信息
func LoadExcludeFromFile(ignoreFilePath string) ([]string, error) {
	content, err := os.ReadFile(ignoreFilePath)
	if err != nil {
		return nil, err
	}
	
	return parseIgnoreFileContent(string(content)), nil
}

// LoadExcludeFromFileOrEmpty 从忽略文件加载排除模式，文件不存在时返回空列表
//
// 参数:
//   - ignoreFilePath: 忽略文件路径
//
// 返回:
//   - []string: 排除模式列表
func LoadExcludeFromFileOrEmpty(ignoreFilePath string) []string {
	patterns, err := LoadExcludeFromFile(ignoreFilePath)
	if err != nil {
		return []string{}
	}
	return patterns
}

// parseIgnoreFileContent 解析忽略文件内容
//
// 参数:
//   - content: 文件内容
//
// 返回:
//   - []string: 排除模式列表
func parseIgnoreFileContent(content string) []string {
	var patterns []string
	
	// 按行分割
	lines := filepath.SplitList(content)
	
	for _, line := range lines {
		// 去除空白
		line = filepath.Clean(line)
		
		// 跳过空行和注释
		if line == "" || line[0] == '#' {
			continue
		}
		
		patterns = append(patterns, line)
	}
	
	return patterns
}
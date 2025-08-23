package types

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// FileFilter 文件过滤器接口
//
// 用于判断文件是否应该被跳过（不处理）
type FileFilter interface {
	// ShouldSkip 判断文件是否应该被跳过
	//
	// 参数:
	//   - path: 文件路径
	//   - info: 文件信息
	//
	// 返回:
	//   - bool: true 表示应该跳过，false 表示应该处理
	ShouldSkip(path string, info os.FileInfo) bool
}

// FilterOptions 过滤配置选项
//
// 用于指定压缩时的文件过滤条件：
//   - Include: 包含模式列表，支持 glob 语法，只有匹配的文件才会被处理
//   - Exclude: 排除模式列表，支持 glob 语法，匹配的文件会被跳过
//   - MaxSize: 最大文件大小限制（字节），0 表示无限制，超过此大小的文件会被跳过
//   - MinSize: 最小文件大小限制（字节），默认为 0，小于此大小的文件会被跳过
type FilterOptions struct {
	Include []string // 包含模式，支持 glob 语法，只处理匹配的文件
	Exclude []string // 排除模式，支持 glob 语法，跳过匹配的文件
	MaxSize int64    // 最大文件大小（字节），0 表示无限制
	MinSize int64    // 最小文件大小（字节），默认为 0
}

// NewFilterOptions 创建新的过滤器选项
//
// 返回:
//   - *FilterOptions: 新的过滤器选项实例
func NewFilterOptions() *FilterOptions {
	return &FilterOptions{}
}

// WithInclude 设置包含模式
//
// 参数:
//   - patterns: 包含模式列表
//
// 返回:
//   - *FilterOptions: 返回自身，支持链式调用
func (f *FilterOptions) WithInclude(patterns ...string) *FilterOptions {
	f.Include = append(f.Include, patterns...)
	return f
}

// WithExclude 设置排除模式
//
// 参数:
//   - patterns: 排除模式列表
//
// 返回:
//   - *FilterOptions: 返回自身，支持链式调用
func (f *FilterOptions) WithExclude(patterns ...string) *FilterOptions {
	f.Exclude = append(f.Exclude, patterns...)
	return f
}

// WithMaxSize 设置最大文件大小限制
//
// 参数:
//   - size: 最大文件大小（字节）
//
// 返回:
//   - *FilterOptions: 返回自身，支持链式调用
func (f *FilterOptions) WithMaxSize(size int64) *FilterOptions {
	f.MaxSize = size
	return f
}

// WithMinSize 设置最小文件大小限制
//
// 参数:
//   - size: 最小文件大小（字节）
//
// 返回:
//   - *FilterOptions: 返回自身，支持链式调用
func (f *FilterOptions) WithMinSize(size int64) *FilterOptions {
	f.MinSize = size
	return f
}

// WithSizeRange 设置文件大小范围
//
// 参数:
//   - minSize: 最小文件大小（字节）
//   - maxSize: 最大文件大小（字节）
//
// 返回:
//   - *FilterOptions: 返回自身，支持链式调用
func (f *FilterOptions) WithSizeRange(minSize, maxSize int64) *FilterOptions {
	f.MinSize = minSize
	f.MaxSize = maxSize
	return f
}

// WithIgnoreFile 从忽略文件加载排除模式
//
// 参数:
//   - ignoreFilePath: 忽略文件路径
//
// 返回:
//   - *FilterOptions: 返回自身，支持链式调用
func (f *FilterOptions) WithIgnoreFile(ignoreFilePath string) *FilterOptions {
	if ignoreFilePath == "" {
		return f
	}
	patterns := LoadExcludeFromFileOrEmpty(ignoreFilePath)
	f.Exclude = append(f.Exclude, patterns...)
	return f
}

// ShouldSkipByParams 判断文件是否应该被跳过(通用方法，用于压缩和解压)
//
// 过滤逻辑:
//  1. 检查文件大小是否符合要求
//  2. 如果指定了包含模式，检查文件是否匹配包含模式
//  3. 检查文件是否匹配排除模式
//
// 参数:
//   - path: 文件路径
//   - size: 文件大小（字节）
//   - isDir: 是否为目录
//
// 返回:
//   - bool: true 表示应该跳过，false 表示应该处理
func (f *FilterOptions) ShouldSkipByParams(path string, size int64, isDir bool) bool {
	// 如果过滤器为空或没有过滤条件，不跳过任何文件
	if !HasFilterConditions(f) {
		return false
	}

	// 1. 检查文件大小 - 不符合大小要求的文件应该被跳过
	if !f.checkSizeFilterByParams(size, isDir) {
		return true
	}

	// 2. 检查包含模式（如果指定了包含模式）
	// 不匹配包含模式的文件应该被跳过
	if len(f.Include) > 0 {
		if !f.matchAnyPattern(f.Include, path) {
			return true
		}
	}

	// 3. 检查排除模式
	// 匹配排除模式的文件应该被跳过
	if len(f.Exclude) > 0 {
		if f.matchAnyPattern(f.Exclude, path) {
			return true
		}
	}

	// 通过所有检查，不应该被跳过
	return false
}

// checkSizeFilterByParams 检查文件大小是否符合过滤条件（通用方法）
//
// 参数:
//   - size: 文件大小（字节）
//   - isDir: 是否为目录
//
// 返回:
//   - bool: true 表示符合大小要求，false 表示不符合大小要求
func (f *FilterOptions) checkSizeFilterByParams(size int64, isDir bool) bool {
	if isDir {
		return true // 目录总是符合大小要求
	}

	// 检查最小大小 - 文件太小不符合要求
	if f.MinSize > 0 && size < f.MinSize {
		return false
	}

	// 检查最大大小 - 文件太大不符合要求
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

	// 3. 尝试匹配目录（处理以路径分隔符结尾的模式）
	if len(pattern) > 0 && os.IsPathSeparator(pattern[len(pattern)-1]) {
		dirPattern := pattern[:len(pattern)-1]
		if matched, err := filepath.Match(dirPattern, filepath.Base(path)); err == nil && matched {
			return true
		}
		if matched, err := filepath.Match(dirPattern, path); err == nil && matched {
			return true
		}
	}

	// 4. 处理路径中包含模式的情况
	pathParts := strings.Split(filepath.ToSlash(path), "/")
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
	lines := strings.Split(content, "\n")

	for _, line := range lines {
		// 去除空白
		line = strings.TrimSpace(line)

		// 跳过空行和注释
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		patterns = append(patterns, line)
	}

	return patterns
}

package comprx

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gitee.com/MM-Q/comprx/types"
)

// LoadExcludeFromFile 从忽略文件加载排除模式
//
// 参数:
//   - ignoreFilePath: 忽略文件路径（如 ".comprxignore", ".gitignore"）
//
// 返回:
//   - []string: 排除模式列表
//   - error: 错误信息
//
// 支持的文件格式:
//   - 每行一个模式
//   - 支持 # 开头的注释行
//   - 自动忽略空行
//   - 支持 glob 模式匹配
//
// 使用示例:
//
//	patterns, err := comprx.LoadExcludeFromFile(".comprxignore")
func LoadExcludeFromFile(ignoreFilePath string) ([]string, error) {
	file, err := os.Open(ignoreFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("忽略文件不存在: %s", ignoreFilePath)
		}
		return nil, fmt.Errorf("打开忽略文件失败: %w", err)
	}
	defer file.Close()

	var patterns []string
	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		// 跳过空行和注释行
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// 验证模式是否有效
		if _, err := filepath.Match(line, "test"); err != nil {
			return nil, fmt.Errorf("第 %d 行包含无效的 glob 模式 '%s': %w", lineNum, line, err)
		}

		patterns = append(patterns, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("读取忽略文件失败: %w", err)
	}

	return patterns, nil
}

// LoadExcludeFromFileOrEmpty 从忽略文件加载排除模式，文件不存在时返回空列表
//
// 参数:
//   - ignoreFilePath: 忽略文件路径
//
// 返回:
//   - []string: 排除模式列表，文件不存在时返回空列表
//   - error: 错误信息（文件不存在不算错误）
//
// 使用示例:
//
//	patterns, err := comprx.LoadExcludeFromFileOrEmpty(".comprxignore")
func LoadExcludeFromFileOrEmpty(ignoreFilePath string) ([]string, error) {
	patterns, err := LoadExcludeFromFile(ignoreFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil // 文件不存在返回空列表，不是错误
		}
		return nil, err
	}
	return patterns, nil
}

// PackWithFilter 使用过滤选项压缩文件或目录
//
// 参数:
//   - dst: 目标文件路径
//   - src: 源文件路径
//   - filter: 过滤选项
//
// 返回:
//   - error: 错误信息
//
// 使用示例:
//
//	filter := types.FilterOptions{
//	    Include: []string{"*.go", "*.md"},
//	    Exclude: []string{"*_test.go", "vendor/*"},
//	    MaxSize: 10 * 1024 * 1024, // 10MB
//	}
//	err := comprx.PackWithFilter("output.zip", "src/", filter)
func PackWithFilter(dst string, src string, filter types.FilterOptions) error {
	opts := DefaultOptions()
	opts.Filter = filter
	return PackOptions(dst, src, opts)
}
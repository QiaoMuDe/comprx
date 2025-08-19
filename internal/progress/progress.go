package progress

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Operation 操作类型常量
type Operation string

const (
	// 解压操作
	OpInflating  Operation = "inflating"  // 解压文件
	OpCreating   Operation = "creating"   // 创建目录
	OpExtracting Operation = "extracting" // 提取文件（TAR）
	OpLinking    Operation = "linking"    // 创建符号链接

	// 压缩操作
	OpAdding      Operation = "adding"      // 添加文件
	OpStoring     Operation = "storing"     // 存储目录
	OpCompressing Operation = "compressing" // 压缩文件（GZIP）
	OpArchiving   Operation = "archiving"   // 归档文件（TAR）

	// 列表操作
	OpListing Operation = "listing" // 列出文件
)

// ProgressReporter 进度报告器
type ProgressReporter struct {
	writer      io.Writer // 输出位置
	enabled     bool      // 是否启用进度显示
	operation   Operation // 当前操作类型
	archivePath string    // 压缩文件路径
}

// New 创建新的进度报告器
//
// 参数:
//   - writer: 输出位置，如果为 nil 则使用 os.Stdout
//   - enabled: 是否启用进度显示
//
// 返回值:
//   - *ProgressReporter: 进度报告器实例
func New(writer io.Writer, enabled bool) *ProgressReporter {
	if writer == nil {
		writer = os.Stdout
	}

	return &ProgressReporter{
		writer:  writer,
		enabled: enabled,
	}
}

// StartArchive 开始处理压缩文件时调用，打印 Archive 信息
//
// 参数:
//   - archivePath: 压缩文件路径
//   - operation: 操作类型（压缩或解压相关的操作）
//
// 示例输出:
//   Archive: service.zip
func (pr *ProgressReporter) StartArchive(archivePath string, operation Operation) {
	if !pr.enabled {
		return
	}

	pr.archivePath = archivePath
	pr.operation = operation

	// 打印 Archive 信息，类似 Linux unzip
	fmt.Fprintf(pr.writer, "Archive: %s\n", filepath.Base(archivePath))
}

// ReportFile 报告文件处理进度
//
// 参数:
//   - filePath: 文件路径
//   - isDirectory: 是否为目录
//
// 示例输出:
//   creating: service/
//   inflating: service/file.txt
func (pr *ProgressReporter) ReportFile(filePath string, isDirectory bool) {
	if !pr.enabled {
		return
	}

	operation := pr.getFileOperation(isDirectory)
	pr.printFileProgress(operation, filePath)
}

// ReportFileWithOperation 使用指定操作报告文件处理进度
//
// 参数:
//   - filePath: 文件路径
//   - operation: 操作类型
//
// 示例输出:
//   linking: service/symlink -> target
func (pr *ProgressReporter) ReportFileWithOperation(filePath string, operation Operation) {
	if !pr.enabled {
		return
	}

	pr.printFileProgress(operation, filePath)
}

// SetOperation 设置当前操作类型
//
// 参数:
//   - operation: 操作类型
func (pr *ProgressReporter) SetOperation(operation Operation) {
	pr.operation = operation
}

// IsEnabled 检查是否启用进度显示
//
// 返回值:
//   - bool: 是否启用
func (pr *ProgressReporter) IsEnabled() bool {
	return pr.enabled
}

// SetEnabled 设置是否启用进度显示
//
// 参数:
//   - enabled: 是否启用
func (pr *ProgressReporter) SetEnabled(enabled bool) {
	pr.enabled = enabled
}

// SetWriter 设置输出位置
//
// 参数:
//   - writer: 输出位置
func (pr *ProgressReporter) SetWriter(writer io.Writer) {
	if writer == nil {
		writer = os.Stdout
	}
	pr.writer = writer
}

// getFileOperation 根据当前操作和文件类型获取具体的操作名称
func (pr *ProgressReporter) getFileOperation(isDirectory bool) Operation {
	switch pr.operation {
	case OpInflating, OpExtracting:
		if isDirectory {
			return OpCreating
		}
		return pr.operation

	case OpAdding, OpArchiving:
		if isDirectory {
			return OpStoring
		}
		return pr.operation

	case OpCompressing:
		// GZIP 只处理单文件，不区分目录
		return OpCompressing

	case OpListing:
		return OpListing

	default:
		// 默认行为
		if isDirectory {
			return OpCreating
		}
		return OpInflating
	}
}

// printFileProgress 打印文件处理进度
func (pr *ProgressReporter) printFileProgress(operation Operation, filePath string) {
	// 根据操作类型设置不同的格式和缩进
	switch operation {
	case OpInflating:
		fmt.Fprintf(pr.writer, "  inflating: %s\n", filePath)
	case OpCreating:
		fmt.Fprintf(pr.writer, "   creating: %s\n", filePath)
	case OpExtracting:
		fmt.Fprintf(pr.writer, "extracting: %s\n", filePath)
	case OpAdding:
		fmt.Fprintf(pr.writer, "     adding: %s\n", filePath)
	case OpStoring:
		fmt.Fprintf(pr.writer, "    storing: %s\n", filePath)
	case OpCompressing:
		fmt.Fprintf(pr.writer, "compressing: %s\n", filePath)
	case OpArchiving:
		fmt.Fprintf(pr.writer, " archiving: %s\n", filePath)
	case OpLinking:
		fmt.Fprintf(pr.writer, "   linking: %s\n", filePath)
	case OpListing:
		fmt.Fprintf(pr.writer, "   listing: %s\n", filePath)
	default:
		// 通用格式，右对齐操作名称
		fmt.Fprintf(pr.writer, "%11s: %s\n", string(operation), filePath)
	}
}
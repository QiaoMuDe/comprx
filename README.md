# ComprX - Go 压缩解压缩库

[![Go Version](https://img.shields.io/badge/Go-%3E%3D1.24.4-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

ComprX 是一个功能强大、易于使用的 Go 语言压缩解压缩库，支持多种压缩格式，提供线程安全的操作和丰富的配置选项。

## ✨ 特性

- 🗜️ **多格式支持**: ZIP、TAR、TGZ、TAR.GZ、GZ、BZ2/BZIP2
- 🔒 **线程安全**: 所有操作都是线程安全的
- 📊 **进度显示**: 支持多种样式的进度条（文本、Unicode、ASCII）
- 🎛️ **灵活配置**: 支持压缩级别、覆盖设置等多种配置选项
- 💾 **内存操作**: 支持字节数据和字符串的内存压缩/解压
- 🌊 **流式处理**: 支持流式压缩和解压缩
- 📝 **简单易用**: 提供简洁的 API 接口

## 📦 安装

```bash
go get gitee.com/MM-Q/comprx
```

## 🚀 快速开始

### 基本压缩和解压

```go
package main

import (
    "fmt"
    "gitee.com/MM-Q/comprx"
)

func main() {
    // 压缩文件或目录
    err := comprx.Pack("output.zip", "input_dir")
    if err != nil {
        fmt.Printf("压缩失败: %v\n", err)
        return
    }
    
    // 解压文件
    err = comprx.Unpack("output.zip", "output_dir")
    if err != nil {
        fmt.Printf("解压失败: %v\n", err)
        return
    }
    
    fmt.Println("操作完成!")
}
```

### 带进度条的压缩解压

```go
// 压缩时显示进度条
err := comprx.PackProgress("output.tar.gz", "large_directory")

// 解压时显示进度条
err := comprx.UnpackProgress("archive.zip", "output_dir")
```

### 自定义配置

```go
import (
    "gitee.com/MM-Q/comprx"
    "gitee.com/MM-Q/comprx/types"
)

// 创建自定义配置
opts := comprx.Options{
    CompressionLevel:  types.CompressionLevelBest,  // 最佳压缩
    OverwriteExisting: true,                        // 覆盖已存在文件
    ProgressEnabled:   true,                        // 启用进度条
    ProgressStyle:     types.ProgressStyleUnicode,  // Unicode 样式进度条
}

// 使用自定义配置压缩
err := comprx.PackOptions("output.zip", "input_dir", opts)

// 使用自定义配置解压
err := comprx.UnpackOptions("archive.zip", "output_dir", opts)
```

## 🧠 内存压缩 API

### 字节数据压缩

```go
// 压缩字节数据
data := []byte("Hello, World!")
compressed, err := comprx.GzipBytes(data, types.CompressionLevelDefault)

// 解压字节数据
decompressed, err := comprx.UngzipBytes(compressed)
```

### 字符串压缩

```go
// 压缩字符串
text := "这是一个测试字符串"
compressed, err := comprx.GzipString(text, types.CompressionLevelBest)

// 解压为字符串
decompressed, err := comprx.UngzipString(compressed)
```

## 🌊 流式压缩 API

```go
import (
    "os"
    "bytes"
)

// 流式压缩（默认压缩级别）
file, _ := os.Open("input.txt")
defer file.Close()

var buf bytes.Buffer
err := comprx.GzipStream(&buf, file)

// 流式压缩（指定压缩级别）
output, _ := os.Create("output.gz")
defer output.Close()

err := comprx.GzipStreamWithLevel(output, file, types.CompressionLevelBest)

// 流式解压
compressedFile, _ := os.Open("input.gz")
defer compressedFile.Close()

outputFile, _ := os.Create("output.txt")
defer outputFile.Close()

err := comprx.UngzipStream(outputFile, compressedFile)
```

## 📋 支持的格式

| 格式 | 扩展名 | 压缩 | 解压 | 说明 |
|------|--------|------|------|------|
| ZIP | `.zip` | ✅ | ✅ | 最常用的压缩格式 |
| TAR | `.tar` | ✅ | ✅ | Unix 标准归档格式 |
| TGZ | `.tgz` | ✅ | ✅ | TAR + GZIP 压缩 |
| TAR.GZ | `.tar.gz` | ✅ | ✅ | TAR + GZIP 压缩 |
| GZIP | `.gz` | ✅ | ✅ | 单文件 GZIP 压缩 |
| BZIP2 | `.bz2`, `.bzip2` | ❌ | ✅ | 仅支持解压 |

## ⚙️ 配置选项

### 压缩级别

```go
types.CompressionLevelDefault     // 默认压缩级别
types.CompressionLevelNone        // 不压缩
types.CompressionLevelFast        // 快速压缩
types.CompressionLevelBest        // 最佳压缩
types.CompressionLevelHuffmanOnly // 仅使用 Huffman 编码
```

### 进度条样式

```go
types.ProgressStyleText     // 文本样式
types.ProgressStyleDefault  // 默认样式
types.ProgressStyleUnicode  // Unicode 样式: ████████████░░░░░░░░ 60%
types.ProgressStyleASCII    // ASCII 样式: [##########          ] 50%
```

## 🏗️ 项目结构

```
comprx/
├── comprx.go              # 主要 API 接口
├── options.go             # 配置选项
├── types/                 # 类型定义
├── config/                # 配置管理
├── internal/
│   ├── core/             # 核心压缩逻辑
│   ├── cxzip/            # ZIP 格式处理
│   ├── cxtar/            # TAR 格式处理
│   ├── cxtgz/            # TGZ 格式处理
│   ├── cxgzip/           # GZIP 格式处理
│   ├── cxbzip2/          # BZIP2 格式处理
│   ├── progress/         # 进度条实现
│   └── utils/            # 工具函数
└── README.md
```

## 🧪 测试

运行所有测试：

```bash
go test ./...
```

运行特定模块测试：

```bash
go test ./internal/cxzip/
go test ./internal/cxgzip/
```

## 📄 许可证

本项目采用 MIT 许可证。详情请参阅 [LICENSE](LICENSE) 文件。

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

## 📞 联系

- 项目地址: [https://gitee.com/MM-Q/comprx](https://gitee.com/MM-Q/comprx)

---

**ComprX** - 让压缩解压变得简单高效！ 🚀
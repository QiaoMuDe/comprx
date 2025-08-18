package comprx

import (
	"fmt"

	"gitee.com/MM-Q/comprx/internal/utils"
	"gitee.com/MM-Q/comprx/internal/zip"
	"gitee.com/MM-Q/comprx/types"
)

// // 统一入口：src 可以是 string（路径/文件）、[]byte、io.Reader
// func Pack(dst string, src any, opts ...PackOption) error
// func Unpack(src string, dst string, opts ...UnpackOption) error

// // 只压缩内存数据，返回 []byte
// func PackBytes(src any, opts ...PackOption) ([]byte, error)

func Pack(dst string, src string) error {
	// 智能检测压缩文件格式
	compressType, err := types.DetectCompressFormat(dst)
	if err != nil {
		return fmt.Errorf("DetectCompressFormat failed: %v", err)
	}

	// 检查目标文件是否存在, 存在则返回错误
	if utils.Exists(dst) {
		return fmt.Errorf("file %s already exists", dst)
	}

	// 检查目标目录是否存在, 不存在则创建
	if err := utils.EnsureDir(dst); err != nil {
		return fmt.Errorf("EnsureDir failed: %v", err)
	}

	// 根据压缩格式进行打包
	switch compressType {
	case types.CompressTypeZip:
		if err := zip.Zip(dst, src); err != nil {
			return fmt.Errorf("Zip failed: %v", err)
		}
	}

	return nil
}

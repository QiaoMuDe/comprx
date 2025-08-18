package zip

// import (
// 	"archive/zip"
// 	"fmt"
// 	"os"

// 	"gitee.com/MM-Q/comprx/types"
// )

// // GetFileList 获取ZIP文件中的所有文件信息
// func GetFileList(zipPath string) ([]types.FileInfo, error) {
// 	reader, err := zip.OpenReader(zipPath)
// 	if err != nil {
// 		return nil, fmt.Errorf("打开ZIP文件失败: %w", err)
// 	}
// 	defer reader.Close()

// 	var fileInfos []types.FileInfo
// 	for _, file := range reader.File {
// 		info := types.FileInfo{
// 			Name:             file.Name,
// 			Size:             int64(file.UncompressedSize64),
// 			CompressedSize:   int64(file.CompressedSize64),
// 			ModTime:          file.Modified,
// 			IsDir:            file.Mode().IsDir(),
// 			CompressionRatio: types.CalculateCompressionRatio(int64(file.UncompressedSize64), int64(file.CompressedSize64)),
// 			Mode:             file.Mode().String(),
// 		}
// 		fileInfos = append(fileInfos, info)
// 	}

// 	return fileInfos, nil
// }

// // GetFileListLimit 获取ZIP文件中的文件信息（限制数量）
// func GetFileListLimit(zipPath string, limit int) ([]types.FileInfo, int, error) {
// 	reader, err := zip.OpenReader(zipPath)
// 	if err != nil {
// 		return nil, 0, fmt.Errorf("打开ZIP文件失败: %w", err)
// 	}
// 	defer reader.Close()

// 	totalCount := len(reader.File)
// 	actualLimit := limit
// 	if limit <= 0 || limit > totalCount {
// 		actualLimit = totalCount
// 	}

// 	var fileInfos []types.FileInfo
// 	for i := 0; i < actualLimit; i++ {
// 		file := reader.File[i]
// 		info := types.FileInfo{
// 			Name:             file.Name,
// 			Size:             int64(file.UncompressedSize64),
// 			CompressedSize:   int64(file.CompressedSize64),
// 			ModTime:          file.Modified,
// 			IsDir:            file.Mode().IsDir(),
// 			CompressionRatio: types.CalculateCompressionRatio(int64(file.UncompressedSize64), int64(file.CompressedSize64)),
// 			Mode:             file.Mode().String(),
// 		}
// 		fileInfos = append(fileInfos, info)
// 	}

// 	return fileInfos, totalCount, nil
// }

// // GetFileListWithFilter 获取ZIP文件中匹配模式的文件信息
// func GetFileListWithFilter(zipPath string, pattern string) ([]types.FileInfo, error) {
// 	allFiles, err := GetFileList(zipPath)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var filteredFiles []types.FileInfo
// 	for _, file := range allFiles {
// 		if file.MatchesPattern(pattern) {
// 			filteredFiles = append(filteredFiles, file)
// 		}
// 	}

// 	return filteredFiles, nil
// }

// // ListFiles 列出ZIP文件中的所有文件
// func ListFiles(zipPath string) error {
// 	files, err := GetFileList(zipPath)
// 	if err != nil {
// 		return err
// 	}

// 	return printFileList(files, len(files), types.DefaultListOptions())
// }

// // ListFilesLimit 列出ZIP文件中的文件（限制数量）
// func ListFilesLimit(zipPath string, limit int) error {
// 	files, total, err := GetFileListLimit(zipPath, limit)
// 	if err != nil {
// 		return err
// 	}

// 	opts := types.DefaultListOptions()
// 	err = printFileList(files, total, opts)
// 	if err != nil {
// 		return err
// 	}

// 	if limit > 0 && limit < total {
// 		fmt.Printf("\n显示了前 %d 个文件，总共 %d 个文件\n", len(files), total)
// 	}

// 	return nil
// }

// // ListFilesWithFilter 列出ZIP文件中匹配模式的文件
// func ListFilesWithFilter(zipPath string, pattern string) error {
// 	files, err := GetFileListWithFilter(zipPath, pattern)
// 	if err != nil {
// 		return err
// 	}

// 	opts := types.DefaultListOptions()
// 	err = printFileList(files, len(files), opts)
// 	if err != nil {
// 		return err
// 	}

// 	fmt.Printf("\n找到 %d 个匹配 '%s' 的文件\n", len(files), pattern)
// 	return nil
// }

// // printFileList 打印文件列表
// func printFileList(files []types.FileInfo, total int, opts *types.ListOptions) error {
// 	if len(files) == 0 {
// 		fmt.Println("没有找到文件")
// 		return nil
// 	}

// 	// 打印表头
// 	opts.PrintHeader()

// 	// 打印文件信息
// 	for _, file := range files {
// 		printFileInfo(file, opts)
// 	}

// 	// 打印统计信息
// 	fmt.Printf("\n总计: %d 个文件", total)

// 	// 计算总大小
// 	var totalSize, totalCompressed int64
// 	for _, file := range files {
// 		totalSize += file.Size
// 		totalCompressed += file.CompressedSize
// 	}

// 	if totalSize > 0 {
// 		overallRatio := types.CalculateCompressionRatio(totalSize, totalCompressed)
// 		totalSizeInfo := types.FileInfo{Size: totalSize}
// 		totalCompressedInfo := types.FileInfo{Size: totalCompressed}
// 		fmt.Printf("，原始大小: %s，压缩后: %s，压缩率: %.1f%%",
// 			totalSizeInfo.FormatSize(),
// 			totalCompressedInfo.FormatSize(),
// 			overallRatio)
// 	}
// 	fmt.Println()

// 	return nil
// }

// // printFileInfo 打印单个文件信息
// func printFileInfo(file types.FileInfo, opts *types.ListOptions) {
// 	fmt.Printf("%-50s", file.Name)

// 	if opts.ShowSize {
// 		if opts.HumanReadable {
// 			fmt.Printf(" %8s", file.FormatSize())
// 		} else {
// 			fmt.Printf(" %8d", file.Size)
// 		}
// 	}

// 	if opts.ShowCompressed {
// 		if opts.HumanReadable {
// 			fmt.Printf(" %8s", file.FormatCompressedSize())
// 		} else {
// 			fmt.Printf(" %8d", file.CompressedSize)
// 		}
// 	}

// 	if opts.ShowRatio {
// 		fmt.Printf(" %6.1f%%", file.CompressionRatio)
// 	}

// 	if opts.ShowTime {
// 		fmt.Printf(" %16s", file.ModTime.Format("2006-01-02 15:04"))
// 	}

// 	if opts.ShowMode {
// 		fmt.Printf(" %10s", file.Mode)
// 	}

// 	typeStr := "文件"
// 	if file.IsDir {
// 		typeStr = "目录"
// 	}
// 	fmt.Printf(" %4s", typeStr)

// 	fmt.Println()
// }

// // GetFileCount 获取ZIP文件中的文件总数
// func GetFileCount(zipPath string) (int, error) {
// 	reader, err := zip.OpenReader(zipPath)
// 	if err != nil {
// 		return 0, fmt.Errorf("打开ZIP文件失败: %w", err)
// 	}
// 	defer reader.Close()

// 	return len(reader.File), nil
// }

// // GetArchiveInfo 获取ZIP文件的基本信息
// func GetArchiveInfo(zipPath string) (*types.ArchiveInfo, error) {
// 	stat, err := os.Stat(zipPath)
// 	if err != nil {
// 		return nil, fmt.Errorf("获取文件信息失败: %w", err)
// 	}

// 	files, err := GetFileList(zipPath)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var totalSize, totalCompressed int64
// 	var fileCount, dirCount int

// 	for _, file := range files {
// 		if file.IsDir {
// 			dirCount++
// 		} else {
// 			fileCount++
// 		}
// 		totalSize += file.Size
// 		totalCompressed += file.CompressedSize
// 	}

// 	info := &types.ArchiveInfo{
// 		Path:             zipPath,
// 		Format:           "ZIP",
// 		Size:             stat.Size(),
// 		ModTime:          stat.ModTime(),
// 		FileCount:        fileCount,
// 		DirCount:         dirCount,
// 		TotalSize:        totalSize,
// 		CompressedSize:   totalCompressed,
// 		CompressionRatio: types.CalculateCompressionRatio(totalSize, totalCompressed),
// 	}

// 	return info, nil
// }

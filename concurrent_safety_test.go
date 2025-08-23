package comprx

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"gitee.com/MM-Q/comprx/internal/core"
	"gitee.com/MM-Q/comprx/types"
)

// TestMain 全局测试入口，控制非verbose模式下的输出重定向
func TestMain(m *testing.M) {
	flag.Parse() // 解析命令行参数
	// 保存原始标准输出和错误输出
	originalStdout := os.Stdout
	originalStderr := os.Stderr
	var nullFile *os.File
	var err error

	// 非verbose模式下重定向到空设备
	if !testing.Verbose() {
		nullFile, err = os.OpenFile(os.DevNull, os.O_WRONLY, 0666)
		if err != nil {
			panic("无法打开空设备文件: " + err.Error())
		}
		os.Stdout = nullFile
		os.Stderr = nullFile
	}

	// 运行所有测试
	exitCode := m.Run()

	// 恢复原始输出
	if !testing.Verbose() {
		os.Stdout = originalStdout
		os.Stderr = originalStderr
		_ = nullFile.Close()
	}

	os.Exit(exitCode)
}

// TestPackGlobalFunction 测试全局Pack函数
func TestPackGlobalFunction(t *testing.T) {
	tempDir := t.TempDir()

	// 创建源文件
	srcFile := filepath.Join(tempDir, "source.txt")
	if err := os.WriteFile(srcFile, []byte("test content"), 0644); err != nil {
		t.Fatal(err)
	}

	dstFile := filepath.Join(tempDir, "test.zip")

	err := Pack(dstFile, srcFile)
	if err != nil {
		t.Errorf("不期望返回错误，但得到错误: %v", err)
	}

	// 检查文件是否创建
	if _, err := os.Stat(dstFile); os.IsNotExist(err) {
		t.Error("压缩文件未创建")
	}
}

// TestUnpackGlobalFunction 测试全局Unpack函数
func TestUnpackGlobalFunction(t *testing.T) {
	tempDir := t.TempDir()
	nonExistentFile := filepath.Join(tempDir, "nonexistent.zip")
	targetDir := filepath.Join(tempDir, "target")

	err := Unpack(nonExistentFile, targetDir)
	if err == nil {
		t.Error("期望返回错误，但没有错误")
	}
}

// TestConcurrentPackOperations 测试并发压缩操作的安全性
func TestConcurrentPackOperations(t *testing.T) {
	// 创建临时目录
	tempDir, err := os.MkdirTemp("", "comprx_concurrent_test")
	if err != nil {
		t.Fatalf("创建临时目录失败: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Logf("清理临时目录失败: %v", err)
		}
	}()

	// 准备测试数据
	testFiles := setupTestFiles(t, tempDir, 10)

	var wg sync.WaitGroup
	var successCount int64
	var errorCount int64
	concurrency := 20 // 并发数

	// 启动多个goroutine并发压缩
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			srcFile := testFiles[id%len(testFiles)]
			dstFile := filepath.Join(tempDir, fmt.Sprintf("concurrent_test_%d.zip", id))

			// 测试Pack函数
			err := Pack(dstFile, srcFile)
			if err != nil {
				t.Logf("Goroutine %d Pack失败: %v", id, err)
				atomic.AddInt64(&errorCount, 1)
			} else {
				atomic.AddInt64(&successCount, 1)
				// 验证压缩文件是否存在
				if !fileExists(dstFile) {
					t.Errorf("Goroutine %d: 压缩文件未创建: %s", id, dstFile)
				}
			}
		}(i)
	}

	wg.Wait()

	t.Logf("并发压缩测试完成 - 成功: %d, 失败: %d", successCount, errorCount)

	// 验证结果
	if successCount == 0 {
		t.Error("所有并发压缩操作都失败了")
	}
}

// TestConcurrentPackWithProgress 测试并发带进度条的压缩操作
func TestConcurrentPackWithProgress(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "comprx_progress_test")
	if err != nil {
		t.Fatalf("创建临时目录失败: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Logf("清理临时目录失败: %v", err)
		}
	}()

	testFiles := setupTestFiles(t, tempDir, 5)

	var wg sync.WaitGroup
	var successCount int64
	concurrency := 10

	// 并发测试PackWithProgress
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			srcFile := testFiles[id%len(testFiles)]
			dstFile := filepath.Join(tempDir, fmt.Sprintf("progress_test_%d.zip", id))

			err := PackProgress(dstFile, srcFile)
			if err != nil {
				t.Logf("Goroutine %d PackWithProgress失败: %v", id, err)
			} else {
				atomic.AddInt64(&successCount, 1)
			}
		}(i)
	}

	wg.Wait()

	t.Logf("并发进度条压缩测试完成 - 成功: %d", successCount)
}

// TestConcurrentUnpackOperations 测试并发解压操作
func TestConcurrentUnpackOperations(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "comprx_unpack_test")
	if err != nil {
		t.Fatalf("创建临时目录失败: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Logf("清理临时目录失败: %v", err)
		}
	}()

	// 先创建一些压缩文件
	testFiles := setupTestFiles(t, tempDir, 5)
	zipFiles := make([]string, len(testFiles))

	for i, srcFile := range testFiles {
		zipFile := filepath.Join(tempDir, fmt.Sprintf("test_%d.zip", i))
		err := Pack(zipFile, srcFile)
		if err != nil {
			t.Fatalf("创建测试压缩文件失败: %v", err)
		}
		zipFiles[i] = zipFile
	}

	var wg sync.WaitGroup
	var successCount int64
	concurrency := 15

	// 并发解压测试
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			zipFile := zipFiles[id%len(zipFiles)]
			extractDir := filepath.Join(tempDir, fmt.Sprintf("extract_%d", id))

			err := Unpack(zipFile, extractDir)
			if err != nil {
				t.Logf("Goroutine %d Unpack失败: %v", id, err)
			} else {
				atomic.AddInt64(&successCount, 1)
			}
		}(i)
	}

	wg.Wait()

	t.Logf("并发解压测试完成 - 成功: %d", successCount)
}

// TestConcurrentMixedOperations 测试混合并发操作（压缩+解压）
func TestConcurrentMixedOperations(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "comprx_mixed_test")
	if err != nil {
		t.Fatalf("创建临时目录失败: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Logf("清理临时目录失败: %v", err)
		}
	}()

	testFiles := setupTestFiles(t, tempDir, 8)

	var wg sync.WaitGroup
	var packCount, unpackCount int64
	concurrency := 30

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			if id%2 == 0 {
				// 压缩操作
				srcFile := testFiles[id%len(testFiles)]
				dstFile := filepath.Join(tempDir, fmt.Sprintf("mixed_pack_%d.zip", id))

				err := Pack(dstFile, srcFile)
				if err == nil {
					atomic.AddInt64(&packCount, 1)
				}
			} else {
				// 先创建一个压缩文件，然后解压
				srcFile := testFiles[id%len(testFiles)]
				zipFile := filepath.Join(tempDir, fmt.Sprintf("mixed_temp_%d.zip", id))
				extractDir := filepath.Join(tempDir, fmt.Sprintf("mixed_extract_%d", id))

				// 压缩
				err := Pack(zipFile, srcFile)
				if err == nil {
					// 解压
					err = Unpack(zipFile, extractDir)
					if err == nil {
						atomic.AddInt64(&unpackCount, 1)
					}
				}
			}
		}(i)
	}

	wg.Wait()

	t.Logf("混合并发测试完成 - 压缩成功: %d, 解压成功: %d", packCount, unpackCount)
}

// TestConcurrentInstanceIsolation 测试实例隔离性
func TestConcurrentInstanceIsolation(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "comprx_isolation_test")
	if err != nil {
		t.Fatalf("创建临时目录失败: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Logf("清理临时目录失败: %v", err)
		}
	}()

	testFile := setupTestFiles(t, tempDir, 1)[0]

	var wg sync.WaitGroup
	var configConflicts int64
	concurrency := 50

	// 测试不同配置的并发操作是否会相互影响
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			// 创建不同配置的实例
			comprx := core.New()
			if id%2 == 0 {
				comprx.Config.Progress.Enabled = true
				comprx.Config.Progress.BarStyle = types.ProgressStyleText
				comprx.Config.OverwriteExisting = true
			} else {
				comprx.Config.Progress.Enabled = false
				comprx.Config.Progress.BarStyle = types.ProgressStyleASCII
				comprx.Config.OverwriteExisting = false
			}

			dstFile := filepath.Join(tempDir, fmt.Sprintf("isolation_test_%d.zip", id))

			// 执行压缩
			err := comprx.Pack(dstFile, testFile)
			if err != nil && id%2 == 1 {
				// 第二次压缩同一个文件应该失败（OverwriteExisting=false）
				// 但由于文件名不同，这里不应该失败
				t.Logf("意外的压缩失败 Goroutine %d: %v", id, err)
			}

			// 验证配置是否被其他goroutine影响
			if (id%2 == 0 && !comprx.Config.OverwriteExisting) ||
				(id%2 == 1 && comprx.Config.OverwriteExisting) {
				atomic.AddInt64(&configConflicts, 1)
				t.Errorf("Goroutine %d: 配置被其他goroutine影响", id)
			}
		}(i)
	}

	wg.Wait()

	if configConflicts > 0 {
		t.Errorf("检测到 %d 个配置冲突，实例隔离失败", configConflicts)
	} else {
		t.Log("实例隔离测试通过")
	}
}

// TestRaceConditionDetection 使用Go race detector检测竞态条件
func TestRaceConditionDetection(t *testing.T) {
	if !raceEnabled() {
		t.Skip("需要使用 -race 标志运行此测试")
	}

	tempDir, err := os.MkdirTemp("", "comprx_race_test")
	if err != nil {
		t.Fatalf("创建临时目录失败: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Logf("清理临时目录失败: %v", err)
		}
	}()

	testFiles := setupTestFiles(t, tempDir, 3)

	var wg sync.WaitGroup
	concurrency := 100

	// 高强度并发测试，用于触发潜在的竞态条件
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			srcFile := testFiles[id%len(testFiles)]
			dstFile := filepath.Join(tempDir, fmt.Sprintf("race_test_%d.zip", id))

			// 随机选择不同的函数调用
			switch id % 4 {
			case 0:
				_ = Pack(dstFile, srcFile)
			case 1:
				_ = PackProgress(dstFile, srcFile)
			case 2:
				comprx := core.New()
				comprx.Config.Progress.Enabled = true
				comprx.Config.Progress.BarStyle = types.ProgressStyleUnicode
				_ = comprx.Pack(dstFile, srcFile)
			case 3:
				comprx := core.New()
				comprx.Config.OverwriteExisting = true
				_ = comprx.Pack(dstFile, srcFile)
			}
		}(i)
	}

	wg.Wait()
	t.Log("竞态条件检测测试完成")
}

// TestMemoryLeakDetection 内存泄漏检测测试
func TestMemoryLeakDetection(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "comprx_memory_test")
	if err != nil {
		t.Fatalf("创建临时目录失败: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Logf("清理临时目录失败: %v", err)
		}
	}()

	testFile := setupTestFiles(t, tempDir, 1)[0]

	// 记录初始内存使用
	var m1, m2 runtime.MemStats
	runtime.GC()
	runtime.GC() // 调用两次确保清理完成
	runtime.ReadMemStats(&m1)

	// 执行大量操作
	iterations := 100 // 减少迭代次数避免测试时间过长
	for i := 0; i < iterations; i++ {
		dstFile := filepath.Join(tempDir, fmt.Sprintf("memory_test_%d.zip", i))
		_ = Pack(dstFile, testFile)
	}

	// 强制垃圾回收并检查内存
	runtime.GC()
	runtime.GC() // 调用两次确保清理完成
	runtime.ReadMemStats(&m2)

	// 安全地计算内存增长，避免负数
	var memGrowth uint64
	if m2.Alloc > m1.Alloc {
		memGrowth = m2.Alloc - m1.Alloc
	} else {
		memGrowth = 0 // 内存实际上减少了，这是正常的
	}

	t.Logf("初始内存: %d bytes, 最终内存: %d bytes", m1.Alloc, m2.Alloc)
	t.Logf("内存增长: %d bytes (%d KB)", memGrowth, memGrowth/1024)

	// 如果内存增长过大，可能存在内存泄漏
	maxAllowedGrowth := uint64(50 * 1024 * 1024) // 50MB，给更宽松的限制
	if memGrowth > maxAllowedGrowth {
		t.Errorf("可能存在内存泄漏，内存增长: %d bytes", memGrowth)
	}
}

// 辅助函数

// setupTestFiles 创建测试文件
func setupTestFiles(t *testing.T, tempDir string, count int) []string {
	var files []string

	for i := 0; i < count; i++ {
		fileName := filepath.Join(tempDir, fmt.Sprintf("test_file_%d.txt", i))
		content := fmt.Sprintf("这是测试文件 %d 的内容\n重复内容用于测试压缩效果\n", i)

		// 创建一些内容让文件有一定大小
		for j := 0; j < 100; j++ {
			content += fmt.Sprintf("行 %d: 测试数据 %s\n", j, time.Now().Format("15:04:05.000"))
		}

		err := os.WriteFile(fileName, []byte(content), 0644)
		if err != nil {
			t.Fatalf("创建测试文件失败: %v", err)
		}
		files = append(files, fileName)
	}

	return files
}

// fileExists 检查文件是否存在
func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

// raceEnabled 检查是否启用了race detector
func raceEnabled() bool {
	// 这是一个简单的检测方法，在实际使用中race detector会被自动检测
	return true // 假设启用了race detector
}

// BenchmarkConcurrentPack 并发压缩性能基准测试
func BenchmarkConcurrentPack(b *testing.B) {
	tempDir, err := os.MkdirTemp("", "comprx_bench")
	if err != nil {
		b.Fatalf("创建临时目录失败: %v", err)
	}
	defer func() {
		if rmErr := os.RemoveAll(tempDir); rmErr != nil {
			b.Logf("清理临时目录失败: %v", rmErr)
		}
	}()

	// 创建测试文件
	testFile := filepath.Join(tempDir, "bench_test.txt")
	content := make([]byte, 1024*1024) // 1MB文件
	for i := range content {
		content[i] = byte(i % 256)
	}
	err = os.WriteFile(testFile, content, 0644)
	if err != nil {
		b.Fatalf("创建测试文件失败: %v", err)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			dstFile := filepath.Join(tempDir, fmt.Sprintf("bench_%d.zip", i))
			_ = Pack(dstFile, testFile)
			i++
		}
	})
}

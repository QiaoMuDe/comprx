package comprx

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"gitee.com/MM-Q/comprx/internal/core"
)

// TestSimpleConcurrentSafety 简单的并发安全测试
func TestSimpleConcurrentSafety(t *testing.T) {
	// 创建临时目录
	tempDir, err := os.MkdirTemp("", "simple_concurrent_test")
	if err != nil {
		t.Fatalf("创建临时目录失败: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Logf("清理临时目录失败: %v", err)
		}
	}()

	// 创建测试文件
	testFile := filepath.Join(tempDir, "test.txt")
	content := "这是一个测试文件，用于验证并发安全性。\n"
	err = os.WriteFile(testFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}

	var wg sync.WaitGroup
	concurrency := 10
	successCount := 0
	var mu sync.Mutex

	// 启动多个goroutine并发压缩
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			zipFile := filepath.Join(tempDir, fmt.Sprintf("test_%d.zip", id))

			// 测试Pack函数的并发安全性
			err := Pack(zipFile, testFile)
			if err != nil {
				t.Logf("Goroutine %d 压缩失败: %v", id, err)
			} else {
				mu.Lock()
				successCount++
				mu.Unlock()
				t.Logf("Goroutine %d 压缩成功: %s", id, zipFile)
			}
		}(i)
	}

	wg.Wait()

	t.Logf("并发测试完成 - 成功: %d/%d", successCount, concurrency)

	if successCount == 0 {
		t.Error("所有并发操作都失败了")
	}
}

// TestProgressConcurrency 测试进度条的并发安全性
func TestProgressConcurrency(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "progress_concurrent_test")
	if err != nil {
		t.Fatalf("创建临时目录失败: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Logf("清理临时目录失败: %v", err)
		}
	}()

	// 创建测试文件
	testFile := filepath.Join(tempDir, "test.txt")
	content := "测试进度条并发安全性的文件内容。\n"
	err = os.WriteFile(testFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}

	var wg sync.WaitGroup
	concurrency := 5

	// 测试PackWithProgress的并发安全性
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			zipFile := filepath.Join(tempDir, fmt.Sprintf("progress_test_%d.zip", id))
			err := PackProgress(zipFile, testFile)
			if err != nil {
				t.Logf("Goroutine %d PackWithProgress失败: %v", id, err)
			} else {
				t.Logf("Goroutine %d PackWithProgress成功", id)
			}
		}(i)
	}

	wg.Wait()
	t.Log("进度条并发测试完成")
}

// TestInstanceIsolation 测试实例隔离
func TestInstanceIsolation(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "isolation_test")
	if err != nil {
		t.Fatalf("创建临时目录失败: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Logf("清理临时目录失败: %v", err)
		}
	}()

	// 创建测试文件
	testFile := filepath.Join(tempDir, "test.txt")
	content := "测试实例隔离的文件内容。\n"
	err = os.WriteFile(testFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}

	var wg sync.WaitGroup
	concurrency := 8
	var configErrors int
	var mu sync.Mutex

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			// 创建独立的实例并设置不同配置
			comprx := core.New()
			if id%2 == 0 {
				comprx.Config.OverwriteExisting = true
			} else {
				comprx.Config.OverwriteExisting = false
			}

			zipFile := filepath.Join(tempDir, fmt.Sprintf("isolation_%d.zip", id))
			err := comprx.Pack(zipFile, testFile)

			// 验证配置没有被其他goroutine影响
			expectedOverwrite := (id%2 == 0)
			actualOverwrite := comprx.Config.OverwriteExisting

			if expectedOverwrite != actualOverwrite {
				mu.Lock()
				configErrors++
				mu.Unlock()
				t.Errorf("Goroutine %d: 配置被污染，期望 %v，实际 %v",
					id, expectedOverwrite, actualOverwrite)
			}

			if err != nil {
				t.Logf("Goroutine %d 压缩失败: %v", id, err)
			}
		}(i)
	}

	wg.Wait()

	if configErrors > 0 {
		t.Errorf("检测到 %d 个配置污染问题", configErrors)
	} else {
		t.Log("实例隔离测试通过")
	}
}

package utils

import (
	"testing"
)

func TestGetBuffer(t *testing.T) {
	t.Run("获取默认大小缓冲区", func(t *testing.T) {
		size := 32 * 1024
		buffer := GetBuffer(size)

		if len(buffer) != size {
			t.Errorf("缓冲区长度 = %d, want %d", len(buffer), size)
		}

		if cap(buffer) < size {
			t.Errorf("缓冲区容量 = %d, want >= %d", cap(buffer), size)
		}
	})

	t.Run("获取小缓冲区", func(t *testing.T) {
		size := 1024
		buffer := GetBuffer(size)

		if len(buffer) != size {
			t.Errorf("缓冲区长度 = %d, want %d", len(buffer), size)
		}
	})

	t.Run("获取大缓冲区", func(t *testing.T) {
		size := 1024 * 1024 // 1MB
		buffer := GetBuffer(size)

		if len(buffer) != size {
			t.Errorf("缓冲区长度 = %d, want %d", len(buffer), size)
		}
	})

	t.Run("获取超大缓冲区", func(t *testing.T) {
		size := 10 * 1024 * 1024 // 10MB
		buffer := GetBuffer(size)

		if len(buffer) != size {
			t.Errorf("缓冲区长度 = %d, want %d", len(buffer), size)
		}
	})

	t.Run("获取零大小缓冲区", func(t *testing.T) {
		size := 0
		buffer := GetBuffer(size)

		if len(buffer) != size {
			t.Errorf("缓冲区长度 = %d, want %d", len(buffer), size)
		}
	})
}

func TestPutBuffer(t *testing.T) {
	t.Run("归还正常大小缓冲区", func(t *testing.T) {
		size := 64 * 1024
		buffer := GetBuffer(size)

		// 归还缓冲区不应该panic
		PutBuffer(buffer)

		// 再次获取缓冲区，应该能复用
		buffer2 := GetBuffer(size)
		if len(buffer2) != size {
			t.Errorf("复用缓冲区长度 = %d, want %d", len(buffer2), size)
		}
	})

	t.Run("归还大缓冲区", func(t *testing.T) {
		size := 512 * 1024
		buffer := GetBuffer(size)

		// 归还缓冲区不应该panic
		PutBuffer(buffer)
	})

	t.Run("归还超大缓冲区", func(t *testing.T) {
		size := 2 * 1024 * 1024 // 2MB，超过1MB限制
		buffer := make([]byte, size)

		// 归还超大缓冲区不应该panic，但不会被池回收
		PutBuffer(buffer)
	})

	t.Run("归还nil缓冲区", func(t *testing.T) {
		// 归还nil不应该panic
		PutBuffer(nil)
	})

	t.Run("归还空缓冲区", func(t *testing.T) {
		buffer := make([]byte, 0)
		// 归还空缓冲区不应该panic
		PutBuffer(buffer)
	})
}

func TestBufferPoolReuse(t *testing.T) {
	t.Run("缓冲区复用测试", func(t *testing.T) {
		size := 32 * 1024

		// 获取缓冲区
		buffer1 := GetBuffer(size)
		originalCap := cap(buffer1)

		// 修改缓冲区内容
		buffer1[0] = 0xFF
		buffer1[size-1] = 0xFF

		// 归还缓冲区
		PutBuffer(buffer1)

		// 再次获取缓冲区
		buffer2 := GetBuffer(size)

		// 检查是否复用了相同的底层数组
		// 注意：这个测试可能不稳定，因为对象池的行为不是确定的
		if cap(buffer2) == originalCap {
			// 如果容量相同，可能是复用了，检查内容是否被重置
			// 注意：GetBuffer不会清零内容，这是正常的
			t.Logf("可能复用了缓冲区，容量: %d", cap(buffer2))
		}

		if len(buffer2) != size {
			t.Errorf("复用缓冲区长度 = %d, want %d", len(buffer2), size)
		}
	})
}

func TestBufferPoolConcurrency(t *testing.T) {
	t.Run("并发获取和归还缓冲区", func(t *testing.T) {
		const numGoroutines = 100
		const numOperations = 100

		done := make(chan bool, numGoroutines)

		for i := 0; i < numGoroutines; i++ {
			go func() {
				defer func() { done <- true }()

				for j := 0; j < numOperations; j++ {
					size := 32 * 1024
					buffer := GetBuffer(size)

					// 简单使用缓冲区
					buffer[0] = byte(j)
					buffer[size-1] = byte(j)

					PutBuffer(buffer)
				}
			}()
		}

		// 等待所有goroutine完成
		for i := 0; i < numGoroutines; i++ {
			<-done
		}
	})
}

// 基准测试
func BenchmarkGetBuffer(b *testing.B) {
	sizes := []int{
		1024,        // 1KB
		32 * 1024,   // 32KB
		64 * 1024,   // 64KB
		256 * 1024,  // 256KB
		1024 * 1024, // 1MB
	}

	for _, size := range sizes {
		b.Run(formatSize(size), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				buffer := GetBuffer(size)
				PutBuffer(buffer)
			}
		})
	}
}

func BenchmarkGetBufferWithoutPool(b *testing.B) {
	sizes := []int{
		1024,        // 1KB
		32 * 1024,   // 32KB
		64 * 1024,   // 64KB
		256 * 1024,  // 256KB
		1024 * 1024, // 1MB
	}

	for _, size := range sizes {
		b.Run(formatSize(size), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = make([]byte, size)
			}
		})
	}
}

func BenchmarkBufferPoolConcurrent(b *testing.B) {
	size := 32 * 1024

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			buffer := GetBuffer(size)
			PutBuffer(buffer)
		}
	})
}

// 辅助函数：格式化大小显示
func formatSize(size int) string {
	if size < 1024 {
		return "1KB"
	}
	return "32KB"
}

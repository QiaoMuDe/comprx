package utils

import "sync"

// 缓冲区对象池，复用缓冲区减少内存分配
var bufferPool = sync.Pool{
	New: func() interface{} {
		return make([]byte, 32*1024) // 默认32KB缓冲区
	},
}

// GetBuffer 从对象池获取缓冲区
//
// 参数:
//   - size: 缓冲区大小
//
// 返回值:
//   - []byte: 获取到的缓冲区
func GetBuffer(size int) []byte {
	buffer, ok := bufferPool.Get().([]byte)
	if !ok || len(buffer) < size {
		// 如果类型断言失败或池中的缓冲区太小，创建新的
		return make([]byte, size)
	}
	return buffer[:size]
}

// PutBuffer 将缓冲区归还到对象池
//
// 参数:
//   - buffer: 要归还的缓冲区
//
// 说明:
//   - 该函数将缓冲区归还到对象池，以便后续复用。
//   - 只有容量不超过1MB的缓冲区才会被归还，以避免对象池占用过多内存。
func PutBuffer(buffer []byte) {
	if cap(buffer) <= 1024*1024 { // 只回收不超过1MB的缓冲区
		//nolint:staticcheck // SA6002: 忽略装箱警告，对象池的性能收益远大于装箱开销
		bufferPool.Put(buffer)
	}
}

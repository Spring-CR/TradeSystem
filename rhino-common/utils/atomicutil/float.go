package atomicutil

import (
	"math"
	"sync/atomic"
	"unsafe"
)

func LoadFloat64(addr *float64) float64 {
	uint64Val := atomic.LoadUint64((*uint64)(unsafe.Pointer(addr)))
	return math.Float64frombits(uint64Val)
}

func StoreFloat64(addr *float64, val float64) {
	atomic.StoreUint64((*uint64)(unsafe.Pointer(addr)), math.Float64bits(val))
}

func AddFloat64(addr *float64, delta float64) (new float64) {
	for {
		old := math.Float64bits(*addr)
		new = math.Float64frombits(old) + delta
		if atomic.CompareAndSwapUint64(
			(*uint64)(unsafe.Pointer(addr)),
			old,
			math.Float64bits(new),
		) {
			break
		}
	}
	return
}

// CompareAndSwapFloat64 原子地将 `val` 的值从 `old` 替换为 `new`。
// 如果 `*val` 的当前值等于 `old`，则将其设置为 `new` 并返回 true。
// 否则，返回 false。
func CompareAndSwapFloat64(val *float64, old, new float64) (swapped bool) {
    for {
        // 原子地获取当前值
        currentBits := atomic.LoadUint64((*uint64)(unsafe.Pointer(val)))
        currentVal := math.Float64frombits(currentBits)

        // 比较当前值是否等于期望的旧值
        // 注意：直接比较浮点数可能存在精度问题，在实际应用中需谨慎
        if currentVal != old {
            return false // 值已被其他 goroutine 修改，交换失败
        }

        // 准备新的位模式
        newBits := math.Float64bits(new)

        // 尝试进行原子比较并交换
        // 这里比较的是位模式，而不是浮点数值
        if atomic.CompareAndSwapUint64((*uint64)(unsafe.Pointer(val)), currentBits, newBits) {
            return true // 交换成功
        }
        // 如果 CAS 操作失败（说明 currentBits 在加载后已被修改），则循环重试
    }
}
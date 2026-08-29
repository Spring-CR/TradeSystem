package byteutils

import "unsafe"

// 正整数字符串字节数组转整形
func PositiveDecimalBytesToInt(data[]byte) int {
	result := 0
	// 使用 unsafe.Pointer 获取字节数组的底层数据地址
	for i := 0; i < len(data); i++ {
		// 通过 unsafe.Pointer 获取字节，并使用位移优化乘法
		result = (result << 3) + (result << 1) + int(*(*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(&data[0])) + uintptr(i)))) - '0' // result * 10 == (result << 3) + (result << 1)
	}
	return result
}
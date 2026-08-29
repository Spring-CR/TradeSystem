package byteutils

import (
	"reflect"
	"unsafe"
)

func GetZeroCopyString(bytes[]byte)string{
	var str string
	sliceHeader := (*reflect.SliceHeader)(unsafe.Pointer(&bytes))
	stringHeader := (*reflect.StringHeader)(unsafe.Pointer(&str))
	stringHeader.Len = sliceHeader.Len
	stringHeader.Data = sliceHeader.Data
	return str
}

func GetZeroCopyBytes(str string)[]byte{

	var bytes[]byte

	stringHeader := (*reflect.StringHeader)(unsafe.Pointer(&str))
	sliceHeader := (*reflect.SliceHeader)(unsafe.Pointer(&bytes))
	sliceHeader.Len = stringHeader.Len
	sliceHeader.Data = stringHeader.Data
	sliceHeader.Cap = stringHeader.Len

	return bytes
}
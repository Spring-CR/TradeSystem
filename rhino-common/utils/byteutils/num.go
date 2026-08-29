package byteutils

import (
	"math"
	"encoding/binary"
)

func Float32ToBytes(float float32, bytes[]byte) []byte {
	bits := math.Float32bits(float)
	//bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, bits)

	return bytes
}

func BytesToFloat32(bytes []byte) float32 {
	bits := binary.LittleEndian.Uint32(bytes)

	return math.Float32frombits(bits)
}

func Float64ToBytes(float float64, bytes[]byte) []byte {
	bits := math.Float64bits(float)
	//bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, bits)

	return bytes
}

func BytesToFloat64(bytes []byte) float64 {
	bits := binary.LittleEndian.Uint64(bytes)

	return math.Float64frombits(bits)
}

func Uint8ToBytes(v uint8, bytes[]byte){
	bytes[0]=byte(v)
}

func BytesToUint8(bytes[]byte)uint8{
	return uint8(bytes[0])
}

func Uint16ToBytes(v uint16, bytes[]byte){
	binary.LittleEndian.PutUint16(bytes,v)
}

func BytesToUint16(bytes[]byte)uint16{
	return binary.LittleEndian.Uint16(bytes)
}

func Uint32ToBytes(v uint32, bytes[]byte){
	binary.LittleEndian.PutUint32(bytes,v)
}

func BytesToUint32(bytes[]byte)uint32{
	return binary.LittleEndian.Uint32(bytes)
}

func Uint64ToBytes(v uint64, bytes[]byte){
	binary.LittleEndian.PutUint64(bytes,v)
}

func BytesToUint64(bytes[]byte)uint64{
	return binary.LittleEndian.Uint64(bytes)
}


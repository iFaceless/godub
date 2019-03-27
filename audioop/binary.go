package audioop

import (
	"encoding/binary"
)

func Int8LE(b []byte) int8 {
	_ = b[0] // bounds check hint to compiler
	return int8(b[0])
}

func Uint8LE(b []byte) uint8 {
	_ = b[0] // bounds check hint to compiler
	return b[0]
}

func Int16LE(b []byte) int16 {
	return int16(binary.LittleEndian.Uint16(b))
}

func Uint16LE(b []byte) uint16 {
	return binary.LittleEndian.Uint16(b)
}

func Int32LE(b []byte) int32 {
	return int32(binary.LittleEndian.Uint32(b))
}

func Uint32LE(b []byte) uint32 {
	return binary.LittleEndian.Uint32(b)
}

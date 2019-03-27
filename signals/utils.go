package signals

import (
	"encoding/binary"
	"io"
)

func maxBoundValue(bitDepth int) int {
	switch bitDepth {
	case 8:
		return 0x7F
	case 16:
		return 0x7FFF
	case 32:
		return 0x7FFFFFFF
	default:
		panic("invalid bit depth")
	}
}

func binaryWriteFunc(bitDepth int) func(w io.Writer, v int) error {
	return func(w io.Writer, v int) error {
		switch bitDepth {
		case 8:
			return binary.Write(w, binary.LittleEndian, int8(v))
		case 16:
			return binary.Write(w, binary.LittleEndian, int16(v))
		case 32:
			return binary.Write(w, binary.LittleEndian, int32(v))
		default:
			panic("invalid bit depth")
		}
	}
}

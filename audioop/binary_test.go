package audioop

import (
	"testing"

	"encoding/binary"

	"github.com/stretchr/testify/assert"
)

func TestInt8(t *testing.T) {
	buf := []byte("\x12\xA2")
	assert.Equal(t, int8(0x12), Int8LE(buf))
}

func TestUint8(t *testing.T) {
	buf := []byte("\x12\xA2")
	assert.Equal(t, uint8(0x12), Uint8LE(buf))
}

func TestInt16(t *testing.T) {
	buf := []byte("\x12\xA2")

	assert.Equal(t, int16(-24046), Int16LE(buf))
}

func TestUint16(t *testing.T) {
	buf := []byte("\x12\xA2")

	assert.Equal(t, uint16(0xa212), Uint16LE(buf))
}

func TestInt32(t *testing.T) {
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, 10248)

	assert.Equal(t, int32(10248), Int32LE(buf))
}

func TestUint32(t *testing.T) {
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, 10248)

	assert.Equal(t, uint32(0x2808), Uint32LE(buf))
}

func BenchmarkInt32(b *testing.B) {
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, 10248)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Int32LE(buf)
	}
}

// 10000000	       124 ns/op
func BenchmarkUint32(b *testing.B) {
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, 10248)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Uint32LE(buf)
	}
}

// 2000000000	         0.35 ns/op
func BenchmarkBuiltinUint32(b *testing.B) {
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, 10248)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		binary.LittleEndian.Uint32(buf)
	}
}

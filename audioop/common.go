package audioop

import (
	"bytes"
	"encoding/binary"
	"math"
)

func getSamples(cp []byte, size int) ([]int32, error) {
	read, err := getReadBinaryFunc(size)
	if err != nil {
		return nil, err
	}

	count := sampleCount(cp, size)
	samples := make([]int32, 0, count)

	for i := 0; i < count; i++ {
		start := i * size
		end := start + size
		samples = append(samples, read(cp[start:end]))
	}

	return samples, nil
}

func getReadBinaryFunc(size int) (func(b []byte) int32, error) {
	switch size {
	case 1:
		return func(b []byte) int32 {
			return int32(Int8LE(b))
		}, nil
	case 2:
		return func(b []byte) int32 {
			return int32(Int16LE(b))
		}, nil
	case 4:
		return func(b []byte) int32 {
			return Int32LE(b)
		}, nil
	default:
		return nil, NewError("failed to get sample, incorrect size: %d", size)
	}
}

func getSample(cp []byte, size int, offset int) (int32, error) {
	start := offset * size
	end := start + size

	switch size {
	case 1:
		return int32(Int8LE(cp[start:end])), nil
	case 2:
		return int32(Int16LE(cp[start:end])), nil
	case 4:
		return Int32LE(cp[start:end]), nil
	default:
		return 0, NewError("failed to get sample, incorrect size: %d", size)
	}
}

func sampleCount(cp []byte, size int) int {
	return len(cp) / size
}

func putSample(cp []byte, size int, offset int, value int32) error {
	start := offset * size
	end := start + size

	write := func(v interface{}) error {
		return binary.Write(bytes.NewBuffer(cp[start:end]), binary.LittleEndian, v)
	}

	switch size {
	case 1:
		return write(int8(value))
	case 2:
		return write(int16(value))
	case 4:
		return write(int32(value))
	default:
		return NewError("size should be 1, 2, or 4")
	}
}

func overflow(value int32, size int) int32 {
	minValue := getMinValue(size)
	maxValue := getMaxValue(size)
	if minValue <= value && value <= maxValue {
		return value
	}

	bits := size * 8
	offset := int(math.Pow(2, float64(bits-1)))
	result := (int(value) + offset) % (int(math.Pow(2, float64(bits))))
	return int32(result - offset)
}

func getClipFunc(size int) func(value int32) int32 {
	maxValue := getMaxValue(size)
	minValue := getMinValue(size)
	return func(value int32) int32 {
		return MaxInt32(MinInt32(value, maxValue), minValue)
	}
}

func getMaxValue(size int) int32 {
	switch size {
	case 1:
		return 0x7f
	case 2:
		return 0x7fff
	case 4:
		return 0x7fffffff
	default:
		return 0
	}
}

func getMinValue(size int) int32 {
	switch size {
	case 1:
		return -0x80
	case 2:
		return -0x8000
	case 4:
		return -0x80000000
	default:
		return 0
	}
}

func checkParameters(length int, size int) error {
	err := checkSize(size)
	if err != nil {
		return err
	}

	if length%size != 0 {
		return NewError("not a whole number of frames")
	}

	return nil
}

func checkSize(size int) error {
	if size != 1 && size != 2 && size != 4 {
		return NewError("size should be 1, 2, or 4")
	}
	return nil
}

type Int32Iterator struct {
	inner     []int32
	currIndex int
}

func NewInt32Interator(items []int32) *Int32Iterator {
	return &Int32Iterator{
		inner:     items,
		currIndex: 0,
	}
}

func (it *Int32Iterator) Next() int32 {
	r := it.inner[it.currIndex]
	it.currIndex += 1
	return r
}

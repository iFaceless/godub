package utils

import (
	"time"
)

func MapKeys(m map[string]interface{}) (keys []string) {
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func ConcatenateByteSlice(items ...[]byte) []byte {
	result := make([]byte, 0)
	for _, item := range items {
		result = append(result, item...)
	}

	return result
}

func Milliseconds(d time.Duration) int {
	return int(d / time.Millisecond)
}

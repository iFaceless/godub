package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConcatenateByteSlice(t *testing.T) {
	a := []byte{1, 2, 3}
	b := []byte{4, 5}
	c := []byte{6, 7, 8}

	expected := []byte{1, 2, 3, 4, 5, 6, 7, 8}
	assert.Equal(t, expected, ConcatenateByteSlice(a, b, c))
}

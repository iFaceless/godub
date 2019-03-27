package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStrSet(t *testing.T) {
	strSet := NewSet("hello", "world")
	assert.True(t, strSet.Has("hello"))
	assert.True(t, strSet.Has("world"))
	assert.False(t, strSet.Has("none"))
	assert.Equal(t, strSet.Len(), 2)
}

func TestIntSet(t *testing.T) {
	intSet := NewSet(1, 2, 3)
	assert.True(t, intSet.Has(1))
	assert.False(t, intSet.Has(0))
	assert.Equal(t, intSet.Len(), 3)
}

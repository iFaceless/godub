package audioop

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_checkParameters(t *testing.T) {
	assert.Nil(t, checkParameters(12, 1))
	assert.Nil(t, checkParameters(12, 2))
	assert.Nil(t, checkParameters(12, 4))

	assert.Error(t, checkParameters(0, 0))
	assert.Error(t, checkParameters(9, 2))
}

func Test_getMaxValue(t *testing.T) {
	assert.Equal(t, int32(0x7f), getMaxValue(1))
	assert.Equal(t, int32(0x7fff), getMaxValue(2))
	assert.Equal(t, int32(0x7fffffff), getMaxValue(4))
}

func Test_getMinValue(t *testing.T) {
	assert.Equal(t, int32(-0x80), getMinValue(1))
	assert.Equal(t, int32(-0x8000), getMinValue(2))
	assert.Equal(t, int32(-0x80000000), getMinValue(4))
}

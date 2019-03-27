package godub

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVolume(t *testing.T) {
	v := NewVolumeFromRatio(1, 0, true)
	assert.Equal(t, 0, int(v))

	v = NewVolumeFromRatio(10, 0, true)
	assert.Equal(t, 20, int(v))

	v = NewVolumeFromRatio(10, 0, false)
	assert.Equal(t, 10, int(v))

	assert.Equal(t, 10, int(Volume(20).ToRatio(true)))
	assert.Equal(t, 10, int(Volume(10).ToRatio(false)))

	assert.Equal(t, "2.120dBFS", Volume(2.12).String())
}

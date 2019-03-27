package audioop

import (
	"testing"

	"math/rand"

	"github.com/stretchr/testify/assert"
)

func TestAbsInt32(t *testing.T) {
	assert.Equal(t, int32(0), AbsInt32(0))
	assert.Equal(t, int32(1), AbsInt32(1))
	assert.Equal(t, int32(1), AbsInt32(-1))
}

func BenchmarkGCD(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GCD(rand.Int(), rand.Int())
	}
}

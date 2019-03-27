package signals

import "math/rand"

type WhiteNoiseSignal struct {
	signal
}

func NewWhiteNoiseSignal() *WhiteNoiseSignal {
	g := &WhiteNoiseSignal{}
	g.signal = newSignal(g)
	return g
}

func (g *WhiteNoiseSignal) Generate(count int) (samples []float64) {
	for i := 0; i < count; i++ {
		samples = append(samples, rand.Float64()*2-1.0)
	}
	return samples
}

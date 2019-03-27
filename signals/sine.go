package signals

import "math"

type SineSignal struct {
	signal
	frequency float64
}

func NewSineSignal(frequency float64) *SineSignal {
	g := &SineSignal{
		frequency: frequency,
	}
	g.signal = newSignal(g)
	return g
}

func (g *SineSignal) Generate(count int) (samples []float64) {
	sineOf := (g.frequency * 2 * math.Pi) / float64(g.sampleRate)

	for i := 0; i < count; i++ {
		samples = append(samples, math.Sin(sineOf*float64(i)))
	}

	return samples
}

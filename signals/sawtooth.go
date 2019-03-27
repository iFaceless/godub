package signals

type SawtoothSignal struct {
	signal
	frequency float64
	dutyCycle float64
}

func NewSawtoothSignal(freq float64) *SawtoothSignal {
	g := &SawtoothSignal{
		frequency: freq,
		dutyCycle: 1.0,
	}
	g.signal = newSignal(g)
	return g
}

func (g *SawtoothSignal) WithDutyCycle(v float64) *SawtoothSignal {
	g.dutyCycle = v
	return g
}

func (g *SawtoothSignal) Generate(count int) (samples []float64) {
	cycleLength := float64(g.sampleRate) / g.frequency
	midpoint := float64(cycleLength) * g.dutyCycle
	ascendLength := midpoint
	descendLength := cycleLength - ascendLength

	for i := 0; i < count; i++ {
		pos := i % int(cycleLength)
		if pos < int(midpoint) {
			samples = append(samples, 2*float64(pos)/ascendLength-1.0)
		} else {
			samples = append(samples, 1.0-(2*(float64(pos)-midpoint)/descendLength))
		}
	}

	return samples
}

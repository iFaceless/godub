package signals

type PulseSignal struct {
	signal
	frequency float64
	dutyCycle float64
}

func NewPulseSignal(freq float64) *PulseSignal {
	g := &PulseSignal{
		frequency: freq,
		dutyCycle: 0.5,
	}
	g.signal = newSignal(g)
	return g
}

func (g *PulseSignal) WithDutyCycle(v float64) *PulseSignal {
	g.dutyCycle = v
	return g
}

func (g *PulseSignal) Generate(count int) (samples []float64) {
	cycleLength := int(float64(g.sampleRate) / g.frequency)
	pulseLength := float64(cycleLength) * g.dutyCycle

	for i := 0; i < count; i++ {
		if i%cycleLength < int(pulseLength) {
			samples = append(samples, 1.0)
		} else {
			samples = append(samples, -1.0)
		}
	}

	return samples
}

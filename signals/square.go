package signals

type SquareSignal struct {
	PulseSignal
}

func NewSquareSignal(freq float64) *SquareSignal {
	return &SquareSignal{
		PulseSignal: *NewPulseSignal(freq).WithDutyCycle(0.5),
	}
}

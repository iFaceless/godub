package signals

type TriangleSignal struct {
	SawtoothSignal
}

func NewTriangleSignal(freq float64) *TriangleSignal {
	return &TriangleSignal{*NewSawtoothSignal(freq)}
}

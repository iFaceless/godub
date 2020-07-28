package signals

import (
	"time"

	"bytes"

	"github.com/yeoji/godub"
)

type SignalGenerator interface {
	GenerateAudioSegment(duration time.Duration, volume godub.Volume) (*godub.AudioSegment, error)
	Generate(count int) []float64
}

type signal struct {
	sampleRate int
	bitDepth   int

	// The real generator to generate signals.
	inner SignalGenerator
}

func newSignal(g SignalGenerator) signal {
	return signal{sampleRate: 44100, bitDepth: 16, inner: g}
}

func (g *signal) WithSampleRate(rate int) *signal {
	g.sampleRate = rate
	return g
}

func (g *signal) WithBitDepth(d int) *signal {
	g.bitDepth = d
	return g
}

// Generate generates an audio segment with given duration and volume.
func (g *signal) GenerateAudioSegment(duration time.Duration, volume godub.Volume) (*godub.AudioSegment, error) {
	data := make([]byte, 0)
	buf := bytes.NewBuffer(data)
	writeFunc := binaryWriteFunc(g.bitDepth)

	sampleCount := g.sampleRate * int(duration.Seconds())
	ratio := volume.ToRatio(true)
	maxBound := float64(maxBoundValue(g.bitDepth))

	for _, v := range g.inner.Generate(sampleCount) {
		writeFunc(buf, int(v*maxBound*ratio))
	}

	sampleWidth := g.bitDepth / 8
	return godub.NewAudioSegment(
		buf.Bytes(),
		godub.Channels(1),
		godub.SampleWidth(uint16(sampleWidth)),
		godub.FrameRate(uint32(g.sampleRate)),
		godub.FrameWidth(uint32(sampleWidth)),
	)
}

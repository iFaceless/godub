package godub

import (
	"fmt"

	"bytes"

	"io"

	"io/ioutil"

	"github.com/iFaceless/godub/converter"
	"github.com/iFaceless/godub/wav"
)

type Loader struct {
	converter *converter.Converter
	buf       io.Writer
}

func NewLoader() *Loader {
	var buf bytes.Buffer
	return &Loader{
		converter: converter.NewConverter(&buf),
		buf:       &buf,
	}
}

func (l *Loader) WithParams(params ...string) *Loader {
	l.converter.WithParams(params...)
	return l
}

func (l *Loader) Load(src interface{}) (*AudioSegment, error) {
	var buf []byte

	switch r := src.(type) {
	case io.Reader:
		result, err := ioutil.ReadAll(r)
		if err != nil {
			return nil, err
		}
		buf = result
	case string:
		result, err := ioutil.ReadFile(r)
		if err != nil {
			return nil, err
		}
		buf = result
	default:
		return nil, fmt.Errorf("expected `io.Reader` or file path to original audio")
	}

	// Try to decode it as wave audio
	waveAudio, err := wav.Decode(bytes.NewReader(buf))
	if err != nil {
		// Try to convert to wave audio, and decode it again!
		var tmpWavBuf bytes.Buffer
		conv := converter.NewConverter(&tmpWavBuf).WithDstFormat("wav")
		e := conv.Convert(bytes.NewReader(buf))
		if e != nil {
			return nil, e
		}

		waveAudio, e = wav.Decode(&tmpWavBuf)
		if e != nil {
			return nil, err
		}
	}
	return NewAudioSegmentFromWaveAudio(waveAudio)
}

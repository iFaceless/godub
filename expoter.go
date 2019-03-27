package godub

import (
	"io"
	"os"

	"bytes"

	"github.com/iFaceless/godub/converter"
	"github.com/iFaceless/godub/wav"
)

type Exporter struct {
	converter *converter.Converter
	dst       interface{}
}

func NewExporter(dst interface{}) *Exporter {
	return &Exporter{converter: converter.NewConverter(nil), dst: dst}
}

func (e *Exporter) Export(segment *AudioSegment) error {
	wavBuf := bytes.Buffer{}
	err := wav.Encode(&wavBuf, segment.AsWaveAudio())
	if err != nil {
		return err
	}

	var w io.Writer
	switch dst := e.dst.(type) {
	case io.Writer:
		w = dst
	case string:
		f, err := os.Create(dst)
		if err != nil {
			return err
		}
		w = f
		defer f.Close()
	}

	if e.converter.DstFormat() == "wav" {
		_, err := io.Copy(w, &wavBuf)
		return err
	} else {
		// Otherwise, convert it to the dst format using ffmpeg.
		return e.converter.WithWriter(w).Convert(&wavBuf)
	}
}

func (e *Exporter) WithCodec(c string) *Exporter {
	e.converter.WithCodec(c)
	return e
}

func (e *Exporter) WithCover(c string) *Exporter {
	e.converter.WithCover(c)
	return e
}

func (e *Exporter) WithBitRate(rate string) *Exporter {
	e.converter.WithBitRate(rate)
	return e
}

func (e *Exporter) WithSampleRate(rate int) *Exporter {
	e.converter.WithSampleRate(rate)
	return e
}

func (e *Exporter) WithTags(tags map[string]string) *Exporter {
	e.converter.WithTags(tags)
	return e
}

func (e *Exporter) WithDstFormat(f string) *Exporter {
	e.converter.WithDstFormat(f)
	return e
}

func (e *Exporter) WithID3TagVersion(v int) *Exporter {
	e.converter.WithID3TagVersion(v)
	return e
}

func (e *Exporter) WithParams(p ...string) *Exporter {
	e.converter.WithParams(p...)
	return e
}

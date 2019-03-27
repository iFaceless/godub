// Package wav defines WAV audio struct and Chunk struct. You can decode or encode WAV audio
// with this handy package.
//
// WAV audio format reference: http://www.topherlee.com/software/pcm-tut-wavformat.html.
//
// The header of WAV(RIFF) file is 44 bytes long and has following format:
//
// | Positions |  Sample Value  |                       Description                           |
// |-----------|----------------|-------------------------------------------------------------|
// | 1 - 4     | "RIFF"         | Marks the file as a RIFF file                               |
// | 5 - 8     | File size (u32)| Size of the overall file (fill this in after creation)      |
// | 9 - 12    | "WAVE"         | File type header                                            |
// | 13 -16    | "fmt "         | Format chunk marker                                         |
// | 17 - 20   | 16 (u32)       | Length of format data listed above                          |
// | 21 - 22   | 1 (u16)        | Type of format (1 is PCM)                                   |
// | 23 - 24   | 2 (u16)        | Number channels                                             |
// | 25 - 28   | 44100 (u32)    | Sample rate                                                 |
// | 29 - 32   | 176400 (u32)   | (Sample Rate * BitsPerSample * Channels) / 8                |
// | 33 - 34   | 4 (u16)        | Sound: 1 - mono, 2 - stereo: (BitsPerSample * Channels) / 8 |
// | 35 - 36   | 16 (u16)       | Bits per sample                                             |
// | 37 - 40   | "data"         | Data chunk header                                           |
// | 41 - 44   | File size (u32)| Size of the data section                                    |
//
package wav

const (
	AudioFormatPCM = 1
)

var (
	WaveHeader = []byte{'W', 'A', 'V', 'E'}
	RiffHeader = []byte{'R', 'I', 'F', 'F'}
	FmtHeader  = []byte{'f', 'm', 't', ' '}
	DataHeader = []byte{'d', 'a', 't', 'a'}
)

type Chunk struct {
	Header   []byte
	Position int
	Size     uint32
}

type WaveAudio struct {
	Format        uint16
	Channels      uint16
	SampleRate    uint32
	BitsPerSample uint16
	RawData       []byte
}

func (w *WaveAudio) DataSize() uint32 {
	return uint32(len(w.RawData))
}

func (w *WaveAudio) SampleFreq() uint32 {
	return (w.SampleRate * uint32(w.BitsPerSample) * uint32(w.Channels)) / 8
}

func (w *WaveAudio) Sound() uint16 {
	// mono: 8 bit, 1
	// stereo: 16 bit, 2
	return (w.BitsPerSample * w.Channels) / 8
}

package wav

import (
	"encoding/binary"
	"io"
)

// Encode encodes wave audio to a given writer.
// WAV file ref: http://www.topherlee.com/software/pcm-tut-wavformat.html
func Encode(w io.Writer, audio *WaveAudio) error {
	// Write RIFF header
	_, err := w.Write(RiffHeader)
	if err != nil {
		return err
	}
	riffSize := 4 + 8 + 16 + 8 + audio.DataSize()
	err = binary.Write(w, binary.LittleEndian, uint32(riffSize))
	if err != nil {
		return err
	}

	// Write file type
	_, err = w.Write(WaveHeader)
	if err != nil {
		return err
	}

	// Write format
	_, err = w.Write(FmtHeader)
	if err != nil {
		return err
	}
	// Write format length
	err = binary.Write(w, binary.LittleEndian, uint32(16))
	if err != nil {
		return err
	}

	err = binary.Write(w, binary.LittleEndian, audio.Format)
	if err != nil {
		return err
	}

	err = binary.Write(w, binary.LittleEndian, audio.Channels)
	if err != nil {
		return err
	}

	err = binary.Write(w, binary.LittleEndian, audio.SampleRate)
	if err != nil {
		return err
	}

	err = binary.Write(w, binary.LittleEndian, audio.SampleFreq())
	if err != nil {
		return err
	}

	err = binary.Write(w, binary.LittleEndian, audio.Sound())
	if err != nil {
		return err
	}

	err = binary.Write(w, binary.LittleEndian, audio.BitsPerSample)
	if err != nil {
		return err
	}

	// Write data
	_, err = w.Write(DataHeader)
	if err != nil {
		return err
	}

	err = binary.Write(w, binary.LittleEndian, audio.DataSize())
	if err != nil {
		return err
	}

	// Write raw data directly
	_, err = w.Write(audio.RawData)
	if err != nil {
		return err
	}

	return nil
}

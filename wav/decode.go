package wav

import (
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"

	"bytes"
)

func Decode(r io.Reader) (*WaveAudio, error) {
	d, err := NewDecoder(r)
	if err != nil {
		return nil, err
	}

	return d.Decode()
}

type Decoder struct {
	buffer []byte
	chunks []Chunk
}

func NewDecoder(r io.Reader) (*Decoder, error) {
	buf, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	d := &Decoder{buffer: buf}
	d.chunks = d.readChunks()

	return d, nil
}

func (d *Decoder) Decode() (*WaveAudio, error) {
	d.patchHeaders()

	fmtChunk := d.findChunk(FmtHeader)
	if fmtChunk == nil || fmtChunk.Size < 16 {
		return nil, DecodeError("Could not find fmt header in wav data")
	}

	pos := fmtChunk.Position + 8
	audioFormat := binary.LittleEndian.Uint16(d.buffer[pos : pos+2])
	if audioFormat != 1 && audioFormat != 0xFFFE {
		return nil, DecodeError(fmt.Sprintf("unknown audio format 0x%X in wav data", audioFormat))
	}

	channels := binary.LittleEndian.Uint16(d.buffer[pos+2 : pos+4])
	sampleRate := binary.LittleEndian.Uint32(d.buffer[pos+4 : pos+8])
	// bit depth
	bitsPerSample := binary.LittleEndian.Uint16(d.buffer[pos+14 : pos+16])

	dataChunk := d.findChunk(DataHeader)
	if dataChunk == nil {
		return nil, DecodeError("Could not find data header in wav data")
	}

	pos = dataChunk.Position + 8
	return &WaveAudio{
		Format:        audioFormat,
		Channels:      channels,
		SampleRate:    sampleRate,
		BitsPerSample: bitsPerSample,
		RawData:       d.buffer[pos : uint32(pos)+dataChunk.Size],
	}, nil
}

func (d *Decoder) readChunks() []Chunk {
	// The size of the RIFF chunk descriptors
	pos := 12
	subChunks := make([]Chunk, 0)

	for {
		if pos+8 >= len(d.buffer) || len(subChunks) >= 10 {
			break
		}

		header := d.buffer[pos : pos+4]
		subChunkSize := binary.LittleEndian.Uint32(d.buffer[pos+4 : pos+8])
		subChunks = append(
			subChunks,
			Chunk{Header: header, Position: pos, Size: subChunkSize},
		)

		if bytes.Equal(header, DataHeader) {
			// `data` is the last subchunk
			break
		}
		pos += int(subChunkSize) + 8
	}

	return subChunks
}

func (d *Decoder) findChunk(header []byte) *Chunk {
	for _, c := range d.chunks {
		if bytes.Equal(c.Header, header) {
			return &c
		}
	}

	return nil
}

func (d *Decoder) patchHeaders() {
	dataChunk := d.findChunk(DataHeader)
	if dataChunk == nil {
		return
	}

	// Set the file size in the RIFF chunk descriptor
	binary.LittleEndian.PutUint32(d.buffer[4:8], uint32(len(d.buffer)-8))

	// Set the data size in the subchunk
	pos := dataChunk.Position
	binary.LittleEndian.PutUint32(d.buffer[pos+4:pos+8], uint32(len(d.buffer)-pos-8))

	// Refresh chunks
	d.chunks = d.readChunks()
}

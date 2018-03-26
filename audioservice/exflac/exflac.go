package exflac

import (
	"errors"
	"fmt"
	"github.com/mewkiz/flac"
	"github.com/mewkiz/flac/frame"
	"io"
	"math"
)

type encFunction func(int32, int)
type sampleChan chan byte

// DecodedStreamBuffer represents a lazy byte buffer of decoded
// Samples.
type DecodedStreamBuffer struct {
	FileName  string
	Src       *flac.Stream
	BlockSize uint32
	position  int
	curFrame  *frame.Frame
	Samples   []sampleChan
	encode    encFunction
}

// Endianess generic type
type Endianess int

// Endianess types
const (
	LittleEndian Endianess = iota
	BigEndian
)

// NewBuffer returns a new instance of a byte buffer that contanins the
// decoded samples of a flac file encoded in little endian
func NewBuffer(filename string, blockSize uint32) (*DecodedStreamBuffer, error) {
	return NewBufferWithEndianess(filename, blockSize, LittleEndian)
}

// NewBufferWithEndianess returns a new instance of a byte buffer that contanins the
// decoded samples of a flac file with endianess e
func NewBufferWithEndianess(filename string, blockSize uint32, e Endianess) (*DecodedStreamBuffer, error) {
	stream, err := flac.Open(filename)
	if err != nil {
		return nil, err
	}
	ans := &DecodedStreamBuffer{
		FileName:  filename,
		Src:       stream,
		BlockSize: blockSize,
		position:  0,
	}
	var f encFunction
	if e == LittleEndian {
		switch bs := int(math.Ceil(float64(stream.Info.BitsPerSample) / 8)); bs {
		case 1:
			f = ans.PutInt8
		case 2:
			f = ans.lPutInt16
		case 3:
			f = ans.lPutInt24
		case 4:
			f = ans.lPutInt32
		}
	} else {
		switch bs := int(math.Ceil(float64(stream.Info.BitsPerSample) / 8)); bs {
		case 1:
			f = ans.PutInt8
		case 2:
			f = ans.bPutInt16
		case 3:
			f = ans.bPutInt24
		case 4:
			f = ans.bPutInt32
		}
	}
	ans.encode = f
	sChan := make([]sampleChan, stream.Info.NChannels)
	for i := range sChan {
		sChan[i] = make(sampleChan, blockSize)
	}
	ans.Samples = sChan
	return ans, nil
}

// Close termiates the file stream.
func (b *DecodedStreamBuffer) Close() error {
	for _, ch := range b.Samples {
		close(ch)
	}
	return b.Src.Close()
}

// PutInt8 puts a byte on the stream
func (b *DecodedStreamBuffer) PutInt8(v int32, nChan int) {
	b.Samples[nChan] <- byte(v & 0xFF)
}

// PutInt16 puts two bytes on the stream
func (b *DecodedStreamBuffer) lPutInt16(v int32, nChan int) {
	b.Samples[nChan] <- byte(v & 0xFF)
	b.Samples[nChan] <- byte(v >> 8 & 0xFF)
}

// PutInt24 puts two bytes on the stream
func (b *DecodedStreamBuffer) lPutInt24(v int32, nChan int) {
	b.Samples[nChan] <- byte(v & 0xFF)
	b.Samples[nChan] <- byte(v >> 8 & 0xFF)
	b.Samples[nChan] <- byte(v >> 16 & 0xFF)
}

// PutInt32 puts two bytes on the stream
func (b *DecodedStreamBuffer) lPutInt32(v int32, nChan int) {
	b.Samples[nChan] <- byte(v & 0xFF)
	b.Samples[nChan] <- byte(v >> 8 & 0xFF)
	b.Samples[nChan] <- byte(v >> 16 & 0xFF)
	b.Samples[nChan] <- byte(v >> 24 & 0xFF)
}

// PutInt16 puts two bytes on the stream
func (b *DecodedStreamBuffer) bPutInt16(v int32, nChan int) {
	b.Samples[nChan] <- byte(v >> 8 & 0xFF)
	b.Samples[nChan] <- byte(v & 0xFF)
}

// PutInt24 puts two bytes on the stream
func (b *DecodedStreamBuffer) bPutInt24(v int32, nChan int) {
	b.Samples[nChan] <- byte(v >> 16 & 0xFF)
	b.Samples[nChan] <- byte(v >> 8 & 0xFF)
	b.Samples[nChan] <- byte(v & 0xFF)
}

// PutInt32 puts two bytes on the stream
func (b *DecodedStreamBuffer) bPutInt32(v int32, nChan int) {
	b.Samples[nChan] <- byte(v >> 24 & 0xFF)
	b.Samples[nChan] <- byte(v >> 16 & 0xFF)
	b.Samples[nChan] <- byte(v >> 8 & 0xFF)
	b.Samples[nChan] <- byte(v & 0xFF)
}

func (b *DecodedStreamBuffer) Next() {
	var i int64
	end := int64(b.BlockSize)
	for i = 0; i < end; i++ {
		if b.position <= 0 {
			b.updateFrame()
		}
		for j := range b.Samples {
			b.encode(b.curFrame.Subframes[j].Samples[b.position], j)
		}
		b.position--
	}
}

func (b *DecodedStreamBuffer) updateFrame() {
	frm, err := b.Src.ParseNext()
	if err != nil {
		if err == io.EOF {

		}
		panic(fmt.Sprintln(err))
	}
	b.position = frm.Subframes[0].NSamples
	b.curFrame = frm
}

func (b *DecodedStreamBuffer) Seek(offset int) error {
	offset -= b.position
	var frm *frame.Frame
	var err error
	for offset > b.curFrame.Subframes[0].NSamples {
		offset -= b.curFrame.Subframes[0].NSamples
		frm, err = b.Src.Next()
		if err != nil {
			if err == io.EOF {
				return errors.New("Invalid offset, EOF reached")
			}
			panic(fmt.Sprintf("%v", err))
		}
	}
	b.curFrame = frm
	b.position = frm.Subframes[0].NSamples - offset
	return nil
}

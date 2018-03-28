package exflac

import (
	"errors"
	"fmt"
	"github.com/mewkiz/flac"
	"github.com/mewkiz/flac/frame"
	"io"
	"math"
	"os"
)

// ChannelState is the possible states of a stream channel generator
type ChannelState int

// Possible channel states
const (
	Open ChannelState = iota
	Close
)

type sampleChannel chan byte
type encFunction func(int32, sampleChannel)

// FlacBuffer contains a lazy byte buffer of a decoded flac with a go channel per audio channel
type FlacBuffer struct {
	FileName  string
	BlockSize uint32
	Channels  []*audioChannel
	encode    encFunction
}

// audioChannel is a lazy byte buffer of a single audio channel
type audioChannel struct {
	channelIdx int
	position   int
	Src        *flac.Stream
	curFrame   *frame.Frame
	Samples    sampleChannel
	State      ChannelState
	buffer     *FlacBuffer
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
func NewBuffer(filename string, blockSize uint32) (*FlacBuffer, error) {
	return NewBufferWithEndianess(filename, blockSize, LittleEndian)
}

// NewBufferWithEndianess returns a new instance of a byte buffer that contanins the
// decoded samples of a flac file with endianess e
func NewBufferWithEndianess(filename string, blockSize uint32, e Endianess) (*FlacBuffer, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	stream, err := flac.New(file)
	if err != nil {
		return nil, err
	}

	ans := &FlacBuffer{
		FileName:  filename,
		BlockSize: blockSize,
	}
	var f encFunction
	if e == LittleEndian {
		switch bs := int(math.Ceil(float64(stream.Info.BitsPerSample) / 8)); bs {
		case 1:
			f = PutInt8
		case 2:
			f = LPutInt16
		case 3:
			f = LPutInt24
		case 4:
			f = LPutInt32
		}
	} else {
		switch bs := int(math.Ceil(float64(stream.Info.BitsPerSample) / 8)); bs {
		case 1:
			f = PutInt8
		case 2:
			f = BPutInt16
		case 3:
			f = BPutInt24
		case 4:
			f = BPutInt32
		}
	}
	ans.encode = f
	channels := make([]*audioChannel, stream.Info.NChannels)
	for i := range channels {
		channels[i], err = newAudioChannel(file, i, blockSize, ans)
	}
	ans.Channels = channels
	return ans, nil
}

func newAudioChannel(file *os.File, id int, blockSize uint32, parent *FlacBuffer) (*audioChannel, error) {
	stream, err := flac.New(file)
	if err != nil {
		return nil, err
	}
	ans := &audioChannel{
		channelIdx: id,
		position:   0,
		Samples:    make(sampleChannel, blockSize),
		buffer:     parent,
		Src:        stream,
	}
	return ans, nil
}

// Close terminates the file stream.
func (f *FlacBuffer) Close() error {
	for _, ch := range f.Channels {
		close((*ch).Samples)
		err := (*ch).Src.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

// PutInt8 puts a byte on the stream
func PutInt8(val int32, s sampleChannel) {
	v := abs(val)
	packed := v & 0x7F
	if val < 0 {
		packed |= 0x80
	}
	s <- byte(packed)
}

// GetInt8 get a byte from the stream
func GetInt8(s sampleChannel) int32 {
	b := <-s
	ans := int32(b & 0x7F)
	if b&0x80 > 0 {
		ans *= -1
	}
	return ans
}

//LPutInt16 puts two bytes on the stream
func LPutInt16(val int32, s sampleChannel) {
	v := abs(val)
	s <- byte(v & 0xFF)
	s <- calcMSB(val, 2)
}

// LGetInt16 gets a two bytes int in little endian from the stream
func LGetInt16(s sampleChannel) int32 {
	b1 := <-s
	b2 := <-s
	var ans int32
	ans |= int32(b1)
	ans = unpackMSB(ans, b2, 2)
	return ans
}

//LPutInt24 puts two bytes on the stream
func LPutInt24(val int32, s sampleChannel) {
	v := abs(val)
	s <- byte(v & 0xFF)
	s <- byte(v >> 8 & 0xFF)
	s <- calcMSB(val, 3)
}

// LGetInt24 gets a two bytes int in little endian from the stream
func LGetInt24(s sampleChannel) int32 {
	b1 := <-s
	b2 := <-s
	b3 := <-s
	var ans int32
	ans |= int32(b2)
	ans <<= 8
	ans |= int32(b1)
	ans = unpackMSB(ans, b3, 3)
	return ans
}

//LPutInt32 puts two bytes on the stream
func LPutInt32(val int32, s sampleChannel) {
	v := uint32(val)
	s <- byte(v & 0xFF)
	s <- byte(v >> 8 & 0xFF)
	s <- byte(v >> 16 & 0xFF)
	s <- byte(v >> 24 & 0xFF)
}

// LGetInt32 gets a two bytes int in little endian from the stream
func LGetInt32(s sampleChannel) int32 {
	b1 := <-s
	b2 := <-s
	b3 := <-s
	b4 := <-s
	var ans uint32
	ans = uint32(b4)
	ans <<= 8
	ans |= uint32(b3)
	ans <<= 8
	ans |= uint32(b2)
	ans <<= 8
	ans |= uint32(b1)
	return int32(ans)
}

// BPutInt16 puts two bytes on the stream
func BPutInt16(val int32, s sampleChannel) {
	v := abs(val)
	s <- calcMSB(val, 2)
	s <- byte(v & 0xFF)
}

// BGetInt16 gets a two bytes int in big endian from the stream
func BGetInt16(s sampleChannel) int32 {
	b2 := <-s
	b1 := <-s
	var ans int32
	ans |= int32(b1)
	ans = unpackMSB(ans, b2, 2)
	return ans
}

// BPutInt24 puts two bytes on the stream
func BPutInt24(val int32, s sampleChannel) {
	v := abs(val)
	s <- calcMSB(val, 3)
	s <- byte(v >> 8 & 0xFF)
	s <- byte(v & 0xFF)
}

// BGetInt24 gets a two bytes int in big endian from the stream
func BGetInt24(s sampleChannel) int32 {
	b3 := <-s
	b2 := <-s
	b1 := <-s
	var ans int32
	ans |= int32(b2)
	ans <<= 8
	ans |= int32(b1)
	ans = unpackMSB(ans, b3, 3)
	return ans
}

// BPutInt32 puts two bytes on the stream
func BPutInt32(v int32, s sampleChannel) {
	s <- byte(v >> 24 & 0xFF)
	s <- byte(v >> 16 & 0xFF)
	s <- byte(v >> 8 & 0xFF)
	s <- byte(v & 0xFF)
}

// BGetInt32 gets a two bytes int in big endian from the stream
func BGetInt32(s sampleChannel) int32 {
	b4 := <-s
	b3 := <-s
	b2 := <-s
	b1 := <-s
	var ans int32
	ans = int32(b4)
	ans <<= 8
	ans |= int32(b3)
	ans <<= 8
	ans |= int32(b2)
	ans <<= 8
	ans |= int32(b1)
	return ans
}

// Decodes a sample and puts it into the channel
func (a audioChannel) Next() {
	var i int64
	end := int64(a.buffer.BlockSize)
	for i = 0; i < end; i++ {
		if a.position <= 0 {
			err := a.updateFrame()
			if err == io.EOF {
				close(a.Samples)
				a.Src.Close()
				a.State = Close
			}
		}
		a.buffer.encode(a.curFrame.Subframes[a.channelIdx].Samples[a.position], a.Samples)
		a.position--
	}
}

func (a audioChannel) updateFrame() error {
	frm, err := a.Src.ParseNext()
	if err != nil {
		if err == io.EOF {
			return err
		}
		panic(fmt.Sprintln(err))
	}
	a.position = frm.Subframes[a.channelIdx].NSamples
	a.curFrame = frm
	return nil
}

func (a *audioChannel) Seek(offset int) error {
	offset -= a.position
	var frm *frame.Frame
	var err error
	for offset > a.curFrame.Subframes[a.channelIdx].NSamples {
		offset -= a.curFrame.Subframes[a.channelIdx].NSamples
		frm, err = a.Src.Next()
		if err != nil {
			if err == io.EOF {
				return errors.New("Invalid offset, EOF reached")
			}
			panic(fmt.Sprintf("%v", err))
		}
	}
	a.curFrame = frm
	a.position = frm.Subframes[a.channelIdx].NSamples - offset
	return nil
}

func abs(v int32) int32 {
	if v >= 0 {
		return v
	}
	return int32(0x80000000 - (uint32(v) & 0x7FFFFFFF))

}

// calcMSB calculates the most significant octet of a signed integer
// of n byte size stored on a int32 type
func calcMSB(val int32, n uint) byte {
	v := abs(val)
	packed := v >> ((n - 1) * 8) & 0x7F
	if val < 0 {
		packed |= 0x80
	}
	return byte(packed)
}

// Unpacks signed msb of an integer of n bytes size
// into a int32 variable
func unpackMSB(val int32, msb byte, n uint) int32 {
	var unpack int32
	unpack = int32(msb & 0x7F)
	unpack <<= 8 * (n - 1)
	if msb&0x80 > 0 {
		return -1 * (val | unpack)
	}
	return val | unpack
}

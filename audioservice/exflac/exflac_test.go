package exflac_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/zenathark/astream/audioservice/exflac"
)

var _ = Describe("Exflac", func() {
	var (
		filename  string
		block     uint32
		bytetests = []int32{
			0x0,
			0x7,
			0x13,
			0x7F,
			-0x7F,
		}
		b2tests = []int32{
			0x0000,
			0x1007,
			0x13AC,
			0x7FFF,
			-0x7FFF,
		}
		b3tests = []int32{
			0x000000,
			0x011007,
			0x7513AC,
			0x7FFFFF,
			-0x7FFFFF,
		}
		b4tests = []int32{
			0x00000000,
			0x01391007,
			-0x7F5A13AC,
			-0x7FFFFFFF,
		}
	)

	BeforeEach(func() {
		filename = "../../database/m1.flac"
		block = 256
	})

	Describe("Initializing structures", func() {
		Context("Using the file m1.flac", func() {
			It("Should create a structure with two channels", func() {
				flc, _ := NewBuffer(filename, block)
				Ω(flc.FileName).Should(Equal(filename))
				Ω(flc.BlockSize).Should(Equal(block))
				Ω(len(flc.Channels)).Should(Equal(2))
			})
		})
	})

	Describe("Writing integers to channels", func() {
		Context("When writing on integer of a single byte", func() {
			It("Should copy the same value", func() {
				ch := make(chan byte, 1)
				for _, tt := range bytetests {
					PutInt8(tt, ch)
					actual := GetInt8(ch)
					Ω(actual).Should(Equal(tt))
				}
			})
		})
		Context("When writing on integers on little endian ", func() {
			It("Should be in right order using 2 bytes", func() {
				ch := make(chan byte, 2)
				for _, tt := range b2tests {
					LPutInt16(tt, ch)
					actual := LGetInt16(ch)
					Ω(actual).Should(Equal(tt))
				}
			})
			It("Should be in right order using 3 bytes", func() {
				ch := make(chan byte, 3)
				for _, tt := range b3tests {
					LPutInt24(tt, ch)
					actual := LGetInt24(ch)
					Ω(actual).Should(Equal(tt))
				}
			})
			It("Should be in right order using 4 bytes", func() {
				ch := make(chan byte, 4)
				for _, tt := range b4tests {
					LPutInt32(tt, ch)
					actual := LGetInt32(ch)
					Ω(actual).Should(Equal(tt))
				}
			})
		})

		Context("When writing on integers on big endian ", func() {
			It("Should be in right order using 2 bytes", func() {
				ch := make(chan byte, 2)
				for _, tt := range b2tests {
					BPutInt16(tt, ch)
					actual := BGetInt16(ch)
					Ω(actual).Should(Equal(tt))
				}
			})
			It("Should be in right order using 3 bytes", func() {
				ch := make(chan byte, 3)
				for _, tt := range b3tests {
					BPutInt24(tt, ch)
					actual := BGetInt24(ch)
					Ω(actual).Should(Equal(tt))
				}
			})
			It("Should be in right order using 4 bytes", func() {
				ch := make(chan byte, 4)
				for _, tt := range b4tests {
					BPutInt32(tt, ch)
					actual := BGetInt32(ch)
					Ω(actual).Should(Equal(tt))
				}
			})
		})
	})
})

package huffman

import (
	"fmt"
	"testing"
)

func TestBitStringWriteBytes(t *testing.T) {
	t.Run("large value", func(t *testing.T) {
		bs := &BitStringWriter{}
		bs.WriteBytes([]byte{0b1111_1111, 0b1111_1111, 0b0000_0001}, 17)

		Equal(t, []byte{0xff, 0xff, 0x80}, bs.buffer)
		Equal(t, 1, bs.offset)
	})

	t.Run("value less than 8 such that we have to ", func(t *testing.T) {
		bs := &BitStringWriter{}
		bs.WriteBytes([]byte{0b1111_1111}, 8)

		Equal(t, []byte{0xff}, bs.buffer)
		Equal(t, 8, bs.offset)
	})
}

func TestBitStringWrite(t *testing.T) {
	// the following tests will work upon the following example. Since the tree generation is non-deterministic (since I use a map to build the freq table and pulling k,v pairs from a map has no guarantee on order) I'm going to work with a static example (which happens to be a real case)
	//
	// example: "h" => 001  (Write(   0b1, 3))
	//          "e" => 011  (Write(  0b11, 3))
	//          "l" => 10   (Write(  0b10, 2))
	//          "l" => 10   (Write(  0b10, 2))
	//          "o" => 111  (Write( 0b111, 3))
	//          " " => 1100 (Write(0b1100, 4))
	//          "w" => 1101 (Write(0b1101, 4))
	//          "o" => 111  (Write( 0b111, 3))
	//          "r" => 000  (Write(   0b0, 3))
	//          "l" => 10   (Write(  0b10, 2))
	//          "d" => 010  (Write(  0b10, 3))
	//
	// "h"            => offset = 3 [[0010 0000]]
	//                                   ^
	// "e"            => offset = 6 [[0010 1100]]
	//                                       ^
	// "l"            => offset = 8 [[0010 1110]]
	//                                         ^
	//
	// taking a bit of time with this
	//  pre-write:
	// "l" (overflow) => offset = 0 [[0010 1110] [0000 0000]]
	//                                            ^
	// post-write:
	// "l" (overflow) => offset = 2 [[0010 1110] [1000 0000]]
	//                                              ^
	//
	// "o"            => offset = 5 [[0010 1110] [1011 1000]]
	//                                                  ^
	// " "            => offset = 1 [[0010 1110] [1011 1110] [0000 0000]]
	//                                                         ^
	t.Run("write empty", func(t *testing.T) {
		bs := &BitStringWriter{}
		bs.Write(0, 0)
		Equal(t, nil, bs.buffer)
	})

	t.Run("write a single byte", func(t *testing.T) {
		bs := &BitStringWriter{}

		// writing h
		bs.Write(bitPattern("h"))
		Equal(t, []byte{0b0010_0000}, bs.buffer)
		Equal(t, 3, bs.offset)
	})

	t.Run("write to the boundary of the first byte", func(t *testing.T) {
		bs := &BitStringWriter{}

		bs.Write(bitPattern("h"))
		bs.Write(bitPattern("e"))
		bs.Write(bitPattern("l"))

		Equal(t, []byte{0b0010_1110}, bs.buffer)
		Equal(t, 8, bs.offset)
	})

	t.Run("write just beyond the boundary of the first byte", func(t *testing.T) {
		bs := &BitStringWriter{}

		bs.Write(bitPattern("h"))
		bs.Write(bitPattern("e"))
		bs.Write(bitPattern("l"))
		bs.Write(bitPattern("l"))

		Equal(t, []byte{0b0010_1110, 0b1000_0000}, bs.buffer)
		Equal(t, 2, bs.offset)
	})

	t.Run("write a value that is split between bytes", func(t *testing.T) {
		bs := &BitStringWriter{}

		bs.Write(bitPattern("h"))
		bs.Write(bitPattern("e"))
		bs.Write(bitPattern("l"))
		bs.Write(bitPattern("l"))
		bs.Write(bitPattern("o"))
		bs.Write(bitPattern(" "))

		Equal(t, []byte{0b0010_1110, 0b1011_1110, 0b0000_0000}, bs.buffer)
		Equal(t, 1, bs.offset)
	})

	t.Run("write the whole thing", func(t *testing.T) {
		bs := &BitStringWriter{}

		bs.Write(bitPattern("h"))
		bs.Write(bitPattern("e"))
		bs.Write(bitPattern("l"))
		bs.Write(bitPattern("l"))
		bs.Write(bitPattern("o"))
		bs.Write(bitPattern(" "))
		bs.Write(bitPattern("w"))
		bs.Write(bitPattern("o"))
		bs.Write(bitPattern("r"))
		bs.Write(bitPattern("l"))
		bs.Write(bitPattern("d"))

		Equal(t, []byte{0b0010_1110, 0b1011_1110, 0b0110_1111, 0b0001_0010}, bs.buffer)
		Equal(t, 8, bs.offset)
	})
}

func bitPattern(s string) (byte, int) {
	switch s {
	case "h":
		return 0b0001, 3
	case "e":
		return 0b0011, 3
	case "l":
		return 0b0010, 2
	case "o":
		return 0b0111, 3
	case " ":
		return 0b1100, 4
	case "w":
		return 0b1101, 4
	case "r":
		return 0b0000, 3
	case "d":
		return 0b0010, 3
	}
	panic(fmt.Sprintf("unsupported character passed into convertCharToBitPattern: %s", s))
}

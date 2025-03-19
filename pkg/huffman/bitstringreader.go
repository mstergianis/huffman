package huffman

type BitStringReader struct {
	buffer      []byte
	offset      int
	currentByte int
}

func NewBitStringReader(input []byte) *BitStringReader {
	if input == nil || len(input) < 1 {
		return nil
	}
	return &BitStringReader{buffer: input, offset: 0, currentByte: 0}
}

func (bs *BitStringReader) Read(w int) byte {
	if w > 8 {
		panic("error: cannot read more than 8 bits at a time from BitStringReader")
	}
	// [[0101 1100] [1000 0000]]
	//                      ^

	var output byte
	leftBitsRemaining := 8 - bs.offset
	if w > leftBitsRemaining {
		// compute left side
		// TODO trailing 0
		output = (bs.buffer[bs.currentByte] & onesMask(leftBitsRemaining)) << bs.offset

		// compute right side
		rightBits := w - leftBitsRemaining
		rightMask := onesMask(rightBits) << (8 - rightBits)
		output = output | (bs.buffer[bs.currentByte+1] & rightMask >> (8 - rightBits))
	} else {
		mask := onesMask(w) << (8 - (w + bs.offset))
		output = bs.buffer[bs.currentByte] & mask >> (8 - (w + bs.offset))
	}

	bs.addOffset(w)

	return output
}

func (bs *BitStringReader) addOffset(w int) {
	if w+bs.offset >= 8 {
		bs.currentByte++
		bs.offset += w - 8
		return
	}

	bs.offset += w
}

func (bs *BitStringReader) InputRemaining() bool {
	return bs.currentByte < (len(bs.buffer)-1) && bs.offset < 8
}

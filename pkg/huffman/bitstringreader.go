package huffman

import "fmt"

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

func (bs *BitStringReader) Read(w int) (byte, error) {
	if w > 8 {
		return 0, fmt.Errorf("error: cannot read more than 8 bits at a time from BitStringReader")
	}

	var output byte
	leftBitsRemaining := 8 - bs.offset
	if w > leftBitsRemaining {
		// compute left side
		output = (bs.buffer[bs.currentByte] & onesMask(leftBitsRemaining)) << bs.offset

		// compute right side
		rightBits := w - leftBitsRemaining
		rightMaskShift := 8 - rightBits
		rightMask := onesMask(rightBits) << rightMaskShift
		if (bs.currentByte + 1) >= len(bs.buffer) {
			return 0, fmt.Errorf("error: attempting to read byte %d from a buffer with len %d", bs.currentByte+1, len(bs.buffer))
		}
		output = output | (bs.buffer[bs.currentByte+1] & rightMask >> rightMaskShift)
		bs.addOffset(w)
		return output, nil
	}

	// we can just take from the left byte
	mask := onesMask(w) << (8 - (w + bs.offset))
	output = bs.buffer[bs.currentByte] & mask >> (8 - (w + bs.offset))
	bs.addOffset(w)

	return output, nil
}

func (bs *BitStringReader) addOffset(w int) {
	if w+bs.offset >= 8 {
		bs.currentByte++
		bs.offset += w - 8
		return
	}

	bs.offset += w
}

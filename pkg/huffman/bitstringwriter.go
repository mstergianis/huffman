package huffman

type BitStringWriter struct {
	// 0000 0000
	// ^
	buffer []byte
	offset int
}

func (bs *BitStringWriter) WriteBytes(bytes []byte, w int) {
	currentByte := 0
	widthRemaining := w
	for widthRemaining >= 8 {
		bs.Write(bytes[currentByte], 8)
		currentByte++
		widthRemaining -= 8
	}
	if widthRemaining > 0 {
		bs.Write(bytes[currentByte], widthRemaining)
	}
}

// Write takes a byte and the width of the bits within that byte and writes to
// an internal buffer that the object maintains.
//
// Unfortunately we can't fulfill the writer interface. For one thing, we need
// the bitwidth to be passed in otherwise we can't distinguish leading 0s from
// good data. Additionally, outputting the number of bytes written isn't useful,
// nor are errors since this is buffered and shouldn't fail.
func (bs *BitStringWriter) Write(b byte, w int) {
	if w < 1 {
		return
	}
	// write the whole byte
	if bs.offset == 0 || bs.offset >= 8 {
		bs.addByte()
	}

	// do we have enough space for the whole "partial-byte"?
	overflow := bs.offset + w
	if overflow > 8 {
		// 1. write part to existing byte
		numBitsLeft := 8 - bs.offset
		left := b >> (w - numBitsLeft)
		bs.writeToLastByte(left, numBitsLeft)

		// 2. add new byte
		bs.addByte()

		// 3. write overflow to new byte
		numBitsRight := w - numBitsLeft
		right := computeRightByte(b, numBitsRight)
		bs.writeToLastByte(right, numBitsRight)

		return
	}

	bs.writeToLastByte(b, w)
}

func (bs *BitStringWriter) String() string {
	return string(bs.buffer)
}

func (bs *BitStringWriter) Bytes() []byte {
	return bs.buffer
}

func (bs *BitStringWriter) addByte() {
	bs.buffer = append(bs.buffer, 0)
	bs.offset = 0
}

func (bs *BitStringWriter) writeToLastByte(b byte, w int) {
	bs.buffer[len(bs.buffer)-1] = bs.buffer[len(bs.buffer)-1] | (b << (7 - bs.offset - w + 1))
	bs.offset += w
}

package huffman

type BitStringWriter struct {
	buffer []byte
	offset int
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

// WriteBytes takes bytes from its input from right to left writing the minimum number of bits at each phase.
//
// For example, let's say your input is [[0100 0000] [0000 0001]] (w = 9)
// Then WriteBytes will take the rightmost byte 1 and write the sole bit from that
// Write(1, 1)
// Then it will write the entire last byte
// Write(0b0100_0000, 8)
func (bs *BitStringWriter) WriteBytes(bytes []byte, w int) {
	currentByte := len(bytes) - 1
	widthRemaining := w
	for widthRemaining >= 8 {
		if widthRemaining%8 != 0 {
			bs.Write(bytes[currentByte], widthRemaining%8)
			widthRemaining -= widthRemaining % 8
		} else {
			bs.Write(bytes[currentByte], 8)
			widthRemaining -= 8
		}
		currentByte--
	}
	if widthRemaining > 0 {
		bs.Write(bytes[currentByte], widthRemaining)
	}
}

// Writes the trailing 30 bits (ignores the leading 2 bits) of contentLength
// representing the uncompressed content length in bytes.
func (bs *BitStringWriter) WriteContentLength(contentLength uint32) {
	bs.Write(byte(CONTROL_BIT_CONTENT_LENGTH), 2)
	// skip 2 bits
	bits := byte((contentLength & (uint32(onesMask(6)) << 24)) >> 24)
	bs.Write(bits, 6)

	bits = byte((contentLength & (uint32(onesMask(8)) << 16)) >> 16)
	bs.Write(bits, 8)

	bits = byte((contentLength & (uint32(onesMask(8)) << 8)) >> 8)
	bs.Write(bits, 8)

	bits = byte(contentLength & uint32(onesMask(8)))
	bs.Write(bits, 8)
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

func computeRightByte(b byte, w int) byte {
	return b & onesMask(w)
}

const (
	F8 byte = 0b1111_1111 >> iota
)

func onesMask(w int) byte {
	if w > 8 {
		panic("error: onesMask encountered a width greater than 8")
	}

	return F8 >> (8 - w)
}

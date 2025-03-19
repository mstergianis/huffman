package huffman

import "fmt"

func Decode(input []byte) []byte {
	// read in tree
	// read in content
	// decompress
	return nil
}

func NewNodeFromBytes(input []byte) *Node {
	if input == nil || len(input) < 1 {
		return nil
	}
	// read control bits
	bs := NewBitStringReader(input)
	var tree *Node = &Node{}
	newNodeFromBytes(bs, tree)

	return tree
}

func newNodeFromBytes(bs *BitStringReader, n *Node) {
	if n == nil {
		return
	}

	controlBits := ControlBit(bs.Read(2))
	if controlBits == CONTROL_BIT_FREQ_PAIR { // freq pair
		char := bs.Read(8)
		n.freqPair = &freqPair{b: char}
		return
	}

	var (
		remainingNode       **Node
		remainingControlBit ControlBit
	)
	switch controlBits {
	case CONTROL_BIT_LEFT:
		n.left = &Node{}
		newNodeFromBytes(bs, n.left)
		remainingNode = &n.right
		remainingControlBit = CONTROL_BIT_RIGHT
	case CONTROL_BIT_RIGHT:
		n.right = &Node{}
		newNodeFromBytes(bs, n.right)
		remainingNode = &n.left
		remainingControlBit = CONTROL_BIT_LEFT
	}

	if bs.InputRemaining() {
		controlBits := bs.Read(2)
		if controlBits != byte(remainingControlBit) {
			panic(fmt.Sprintf("error: newNodeFromBytes did not get the expected control bits %02b (%s) received %02b", byte(remainingControlBit), remainingControlBit, controlBits))
		}
		*remainingNode = &Node{}
		newNodeFromBytes(bs, *remainingNode)
		return
	}

	panic("error: newNodeFromBytes reached an unreachable line, you have a malformed tree input")
}

type ControlBit byte

const (
	CONTROL_BIT_FREQ_PAIR ControlBit = 0b01
	CONTROL_BIT_LEFT      ControlBit = 0b11
	CONTROL_BIT_RIGHT     ControlBit = 0b10
)

func (cb ControlBit) String() string {
	switch cb {
	case CONTROL_BIT_FREQ_PAIR:
		return "freqPair"
	case CONTROL_BIT_LEFT:
		return "left"
	case CONTROL_BIT_RIGHT:
		return "right"
	}
	panic(fmt.Sprintf("error: ControlBit.String unimplemented %02b", cb))
}

type BitStringReader struct {
	buffer      []byte
	offset      int
	currentByte int
}

func NewBitStringReader(input []byte) *BitStringReader {
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

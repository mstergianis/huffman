package huffman

import "fmt"

func Decode(input []byte) ([]byte, error) {
	if input == nil || len(input) < 1 {
		return nil, fmt.Errorf("error: while decoding the input was empty")
	}
	bs := NewBitStringReader(input)
	// read in tree
	NewNodeFromBytes(bs)

	// read in content

	// decompress
	return nil, nil
}

func NewNodeFromBytes(bs *BitStringReader) *Node {
	if bs == nil {
		return nil
	}
	// read control bits
	var tree *Node = &Node{}
	newNodeFromBytes(bs, tree)

	return tree
}

func newNodeFromBytes(bs *BitStringReader, n *Node) {
	if n == nil || bs == nil {
		return
	}

	controlBits := ControlBit(bs.Read(2))
	if controlBits == CONTROL_BIT_FREQ_PAIR { // freq pair
		char := bs.Read(8)
		n.freqPair = &freqPair{char: char}
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

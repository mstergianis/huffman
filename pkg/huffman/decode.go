package huffman

import (
	"bytes"
	"fmt"
)

func Decode(input []byte) ([]byte, error) {
	if input == nil || len(input) < 1 {
		return nil, fmt.Errorf("error: while decoding the input was empty")
	}
	bs := NewBitStringReader(input)

	// read in the content length
	contentLength, err := bs.ReadContentLength()
	if err != nil {
		return nil, err
	}

	// read in tree
	tree := NewNodeFromBytes(bs)

	// read in content
	contents, err := ReadContent(bs, tree, contentLength)
	if err != nil {
		return nil, err
	}

	return contents, nil
}

func (bs *BitStringReader) ReadContentLength() (ret uint32, err error) {
	bits, err := bs.Read(2)
	if err != nil {
		return 0, err
	}
	if ControlBit(bits) != CONTROL_BIT_CONTENT_LENGTH {
		return 0, fmt.Errorf("error: decoding expected first 2 bits to be the ControlBit header %02b", bits)
	}

	bits, err = bs.Read(6)
	if err != nil {
		return 0, err
	}
	ret = ret | (uint32(bits) << 24)

	bits, err = bs.Read(8)
	if err != nil {
		return 0, err
	}
	ret = ret | (uint32(bits) << 16)

	bits, err = bs.Read(8)
	if err != nil {
		return 0, err
	}
	ret = ret | (uint32(bits) << 8)

	bits, err = bs.Read(8)
	if err != nil {
		return 0, err
	}
	ret = ret | uint32(bits)

	return ret, nil
}

func ReadContent(bs *BitStringReader, tree *Node, contentLength uint32) ([]byte, error) {
	buf := &bytes.Buffer{}
	// read one bit at a time until you reach a leaf node, then write that byte.
	var readBytes uint32 = 0
	for readBytes < contentLength {
		n := tree
		for n.freqPair == nil {
			bit, err := bs.Read(1)
			if err != nil {
				return nil, err
			}
			switch bit {
			case LEFT:
				n = n.left
			case RIGHT:
				n = n.right
			}
		}
		readBytes++
		err := buf.WriteByte(n.freqPair.char)
		if err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}

const (
	LEFT  byte = 0
	RIGHT byte = 1
)

func NewNodeFromBytes(bs *BitStringReader) *Node {
	if bs == nil {
		return nil
	}
	var tree *Node = &Node{}
	newNodeFromBytes(bs, tree)

	return tree
}

func newNodeFromBytes(bs *BitStringReader, n *Node) error {
	if n == nil || bs == nil {
		return nil
	}

	bits, err := bs.Read(2)
	if err != nil {
		return err
	}
	controlBits := ControlBit(bits)
	if controlBits == CONTROL_BIT_FREQ_PAIR {
		char, err := bs.Read(8)
		if err != nil {
			return err
		}
		n.freqPair = &freqPair{char: char}
		return nil
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

	bits, err = bs.Read(2)
	if err != nil {
		return err
	}
	controlBits = ControlBit(bits)
	if controlBits != remainingControlBit {
		panic(fmt.Sprintf("error: newNodeFromBytes did not get the expected control bits %02b (%s) received %02b", byte(remainingControlBit), remainingControlBit, controlBits))
	}
	*remainingNode = &Node{}
	newNodeFromBytes(bs, *remainingNode)
	return nil
}

type ControlBit byte

const (
	CONTROL_BIT_CONTENT_LENGTH ControlBit = 0b00
	CONTROL_BIT_FREQ_PAIR      ControlBit = 0b01
	CONTROL_BIT_LEFT           ControlBit = 0b11
	CONTROL_BIT_RIGHT          ControlBit = 0b10
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

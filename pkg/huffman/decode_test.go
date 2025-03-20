package huffman

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadContent(t *testing.T) {
	tree := &Node{
		left: &Node{
			left: &Node{
				left:  &Node{freqPair: &freqPair{char: 'r'}},
				right: &Node{freqPair: &freqPair{char: 'h'}},
			},
			right: &Node{
				left:  &Node{freqPair: &freqPair{char: 'd'}},
				right: &Node{freqPair: &freqPair{char: 'e'}},
			},
		},
		right: &Node{
			left: &Node{
				freqPair: &freqPair{char: 'l'},
			},
			right: &Node{
				left: &Node{
					left:  &Node{freqPair: &freqPair{char: ' '}},
					right: &Node{freqPair: &freqPair{char: 'w'}},
				},
				right: &Node{freqPair: &freqPair{char: 'o'}},
			},
		},
	}
	bs := NewBitStringReader([]byte{0b0010_1110, 0b1011_1110, 0b0110_1111, 0b0001_0010})
	expected := []byte("hello world")
	contents, err := ReadContent(bs, tree, uint32(len(expected)))
	assert.NoError(t, err)
	Equal(t, expected, contents)
}

func TestNewNodeFromBytes(t *testing.T) {
	t.Run("empty input", func(t *testing.T) {
		bs := NewBitStringReader([]byte{})
		n := NewNodeFromBytes(bs)
		Equal(t, nil, n)
	})

	t.Run("single node tree", func(t *testing.T) {
		input := []byte{0b0101_1100, 0b1000_0000}
		bs := NewBitStringReader(input)
		n := NewNodeFromBytes(bs)

		expected := &Node{freqPair: &freqPair{char: 'r'}}
		Equal(t, expected, n, "expected: %08b actual: %08b", expected.freqPair.char, n.freqPair.char)
	})

	t.Run("multi node tree", func(t *testing.T) {
		input := []byte{
			0b1101_0110,
			0b1100_1001,
			0b0111_0010,
		}
		bs := NewBitStringReader(input)
		n := NewNodeFromBytes(bs)
		expected := &Node{
			left: &Node{
				freqPair: &freqPair{char: 'l'},
			},
			right: &Node{
				freqPair: &freqPair{char: 'r'},
			},
		}
		Equal(t, expected, n)
	})

	t.Run("multi node tree but mirrored (with right node first, then left node)", func(t *testing.T) {
		input := []byte{
			0b1001_0111,
			0b0010_1101,
			0b0110_1100,
		}
		bs := NewBitStringReader(input)
		n := NewNodeFromBytes(bs)
		expected := &Node{
			left: &Node{
				freqPair: &freqPair{char: 'l'},
			},
			right: &Node{
				freqPair: &freqPair{char: 'r'},
			},
		}
		Equal(t, expected, n)
	})

	t.Run("hello world tree", func(t *testing.T) {
		input := []byte{
			0b1111_1101, // l, l, l, f
			0b0111_0010, // r
			0b1001_0110, // r, f, 4 bits of h
			0b1000_1011,
			0b0101_1001,
			0b0010_0101, // up to d
			0b1001_0110, // up to e, then right control bits
			0b1101_0110, // first 4 bits of l
			0b1100_1011, // last 4 bits of l, then right left
			0b1101_0010, // left again then space
			0b0000_1001, // last 4 bits of space, right, then freq control
			0b0111_0111, // w
			0b1001_0110, // right, freq, first 4 bits of o
			0b1111_0000, // last 4 bits of o
		}

		// left
		//          "r" => 000  (Write(   0b0, 3))
		//          "h" => 001  (Write(   0b1, 3))
		//          "d" => 010  (Write(  0b10, 3))
		//          "e" => 011  (Write(  0b11, 3))

		// right
		//          "l" => 10   (Write(  0b10, 2))
		//          "o" => 111  (Write( 0b111, 3))
		//          " " => 1100 (Write(0b1100, 4))
		//          "w" => 1101 (Write(0b1101, 4))
		expected := &Node{
			left: &Node{
				left: &Node{
					left:  &Node{freqPair: &freqPair{char: 'r'}},
					right: &Node{freqPair: &freqPair{char: 'h'}},
				},
				right: &Node{
					left:  &Node{freqPair: &freqPair{char: 'd'}},
					right: &Node{freqPair: &freqPair{char: 'e'}},
				},
			},
			right: &Node{
				left: &Node{freqPair: &freqPair{char: 'l'}},
				right: &Node{
					left: &Node{
						left:  &Node{freqPair: &freqPair{char: ' '}},
						right: &Node{freqPair: &freqPair{char: 'w'}},
					},
					right: &Node{freqPair: &freqPair{char: 'o'}},
				},
			},
		}
		bs := NewBitStringReader(input)
		n := NewNodeFromBytes(bs)
		Equal(t, expected, n)
	})
}

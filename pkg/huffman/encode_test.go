package huffman

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncodeTree(t *testing.T) {
	t.Run("empty input", func(t *testing.T) {
		var n *Node
		bs := &BitStringWriter{}
		n.WriteBytes(bs)
		Equal(t, []byte(nil), bs.Bytes())
	})

	t.Run("single node tree", func(t *testing.T) {
		n := &Node{freqPair: &freqPair{char: 'r', freq: 1}}
		bs := &BitStringWriter{}
		n.WriteBytes(bs)
		Equal(t, []byte{0b0101_1100, 0b1000_0000}, bs.Bytes())
	})

	t.Run("multi node tree", func(t *testing.T) {
		n := &Node{
			freq: 2,
			left: &Node{
				freqPair: &freqPair{char: 'l', freq: 1},
			},
			right: &Node{
				freqPair: &freqPair{char: 'r', freq: 1},
			},
		}
		bs := &BitStringWriter{}
		n.WriteBytes(bs)
		expected := []byte{
			0b1101_0110,
			0b1100_1001,
			0b0111_0010,
		}
		Equal(t, expected, bs.Bytes())
	})

	t.Run("hello world tree", func(t *testing.T) {
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
		n := &Node{
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
		bs := &BitStringWriter{}
		n.WriteBytes(bs)
		expected := []byte{
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
		Equal(t, expected, bs.Bytes())
	})
}

func TestComputeFreqTable(t *testing.T) {
	type testCase struct {
		input    []byte
		expected []freqPair
		name     string
	}

	testCases := []testCase{
		{
			input:    []byte("aaaaaaaa"),
			expected: []freqPair{{char: 'a', freq: 8}},
		},
		{
			input: []byte("hello world"),
			expected: []freqPair{
				{char: ' ', freq: 1},
				{char: 'r', freq: 1},
				{char: 'h', freq: 1},
				{char: 'e', freq: 1},
				{char: 'l', freq: 3},
				{char: 'w', freq: 1},
				{char: 'd', freq: 1},
				{char: 'o', freq: 2},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			output := computeFreqTable(tc.input)

			for _, element := range tc.expected {
				assert.Contains(t, output, element, "output did not contain element: %s", element)
			}
		})
	}
}

func Equal[E any](t assert.TestingT, expected, actual E, msgAndArgs ...any) bool {
	return assert.Equal(t, expected, actual, msgAndArgs...)
}

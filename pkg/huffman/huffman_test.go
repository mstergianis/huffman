package huffman

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHuffman(t *testing.T) {

}

func TestEncodeTree(t *testing.T) {
	t.Run("empty input", func(t *testing.T) {
		var n *Node
		bs := n.Bytes()
		Equal(t, nil, bs)
	})

	t.Run("single node tree", func(t *testing.T) {
		n := &Node{freqPair: &freqPair{b: 'r', freq: 1}}
		bs := n.Bytes()
		Equal(t, []byte{0b0101_1100, 0b1000_0000}, bs)
	})

	t.Run("multi node tree", func(t *testing.T) {
		n := &Node{
			freq: 2,
			left: &Node{
				freqPair: &freqPair{b: 'l', freq: 1},
			},
			right: &Node{
				freqPair: &freqPair{b: 'r', freq: 1},
			},
		}
		bs := n.Bytes()
		expected := []byte{
			0b1101_0110,
			0b1100_1001,
			0b0111_0010,
		}
		Equal(t, expected, bs)
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
					left:  &Node{freqPair: &freqPair{b: 'r'}},
					right: &Node{freqPair: &freqPair{b: 'h'}},
				},
				right: &Node{
					left:  &Node{freqPair: &freqPair{b: 'd'}},
					right: &Node{freqPair: &freqPair{b: 'e'}},
				},
			},
			right: &Node{
				left: &Node{freqPair: &freqPair{b: 'l'}},
				right: &Node{
					left: &Node{
						left:  &Node{freqPair: &freqPair{b: ' '}},
						right: &Node{freqPair: &freqPair{b: 'w'}},
					},
					right: &Node{freqPair: &freqPair{b: 'o'}},
				},
			},
		}
		bs := n.Bytes()
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
		Equal(t, expected, bs)
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
			expected: []freqPair{{b: 'a', freq: 8}},
		},
		{
			input: []byte("hello world"),
			expected: []freqPair{
				{b: ' ', freq: 1},
				{b: 'r', freq: 1},
				{b: 'h', freq: 1},
				{b: 'e', freq: 1},
				{b: 'l', freq: 3},
				{b: 'w', freq: 1},
				{b: 'd', freq: 1},
				{b: 'o', freq: 2},
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

func TestSearch(t *testing.T) {
	type searchReturn struct {
		b []byte
		w int
	}
	type testCase struct {
		n        *Node
		query    byte
		expected searchReturn
		name     string
	}
	tree := &Node{
		right: &Node{
			right: &Node{
				right: makeHighlyRightNestedNode(13, 'c'),
				left: &Node{
					freqPair: &freqPair{
						b:    'r',
						freq: 2,
					},
				},
			},
			left: &Node{
				freqPair: &freqPair{
					b:    'z',
					freq: 2,
				},
			},
		},
		left: &Node{
			freqPair: &freqPair{
				b:    'o',
				freq: 4,
			},
		},
	}

	testCases := []testCase{
		// {
		// 	name:  "find r",
		// 	n:     tree,
		// 	query: 'r',
		// 	expected: searchReturn{
		// 		b: []byte{0b110},
		// 		w: 3,
		// 	},
		// },
		// {
		// 	name:  "find o",
		// 	n:     tree,
		// 	query: 'o',
		// 	expected: searchReturn{
		// 		b: []byte{0b0},
		// 		w: 1,
		// 	},
		// },
		// {
		// 	name:  "find z",
		// 	n:     tree,
		// 	query: 'z',
		// 	expected: searchReturn{
		// 		b: []byte{0b10},
		// 		w: 2,
		// 	},
		// },
		// {
		// 	name:  "find x, which is not present",
		// 	n:     tree,
		// 	query: 'x',
		// 	expected: searchReturn{
		// 		b: nil,
		// 		w: -1,
		// 	},
		// },
		{
			name:  "find an element which is wider than a byte",
			n:     tree,
			query: 'c',
			expected: searchReturn{
				b: []byte{0b1111_1111, 0b1111_1111, 0b0000_0001},
				w: 17,
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actualByte, actualWidth := tc.n.Search(tc.query)

			Equal(t, tc.expected.b, actualByte)
			Equal(t, tc.expected.w, actualWidth)
		})
	}
}

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
		//                                   ^       ^
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

func Equal[E any](t assert.TestingT, expected, actual E, msgAndArgs ...any) bool {
	return assert.Equal(t, expected, actual, msgAndArgs...)
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

func makeHighlyRightNestedNode(depth int, char byte) *Node {
	var head *Node = &Node{}
	var n = head
	for range depth {
		n.right = &Node{}
		n = n.right
	}
	n.right = &Node{freqPair: &freqPair{b: char}}
	return head
}

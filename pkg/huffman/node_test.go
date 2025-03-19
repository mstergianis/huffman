package huffman

import "testing"

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
						char: 'r',
						freq: 2,
					},
				},
			},
			left: &Node{
				freqPair: &freqPair{
					char: 'z',
					freq: 2,
				},
			},
		},
		left: &Node{
			freqPair: &freqPair{
				char: 'o',
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

func makeHighlyRightNestedNode(depth int, char byte) *Node {
	var head *Node = &Node{}
	var n = head
	for range depth {
		n.right = &Node{}
		n = n.right
	}
	n.right = &Node{freqPair: &freqPair{char: char}}
	return head
}

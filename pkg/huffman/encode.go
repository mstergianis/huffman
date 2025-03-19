package huffman

import (
	"fmt"
	"sort"
)

func Encode(input []byte) ([]byte, []byte, error) {
	for i, r := range input {
		if r > 127 || r < 0 {
			return nil, nil, fmt.Errorf("error: encountered a non-ascii character %c at position %d", r, i)
		}
	}
	ordered := computeFreqTable(input)

	tree := NewNode(ordered)
	// printTree(tree)

	bs := &BitStringWriter{}

	for _, b := range []byte(input) {
		bytes, bitWidth := tree.Search(b)
		if bitWidth == -1 {
			return nil, nil, fmt.Errorf("error: cannot find the byte %c in the tree", b)
		}
		bs.WriteBytes(bytes, bitWidth)
	}

	return tree.Bytes(), bs.Bytes(), nil
}

func computeRightByte(b byte, w int) byte {
	return b & onesMask(w)
}

func onesMask(w int) byte {
	switch w {
	case 0:
		return 0
	case 1:
		return 0b0000_0001
	case 2:
		return 0b0000_0011
	case 3:
		return 0b0000_0111
	case 4:
		return 0b0000_1111
	case 5:
		return 0b0001_1111
	case 6:
		return 0b0011_1111
	case 7:
		return 0b0111_1111
	case 8:
		return 0b1111_1111
	}
	panic("error: onesMask encountered a width greater than 8")
}

type freqPair struct {
	char byte
	freq int
}

func (f freqPair) Freq() int {
	return f.freq
}

func (f freqPair) String() string {
	return fmt.Sprintf("('%s', %d)", string(f.char), f.freq)
}

func computeFreqTable(input []byte) (ordered []freqPair) {
	freqTable := make(map[byte]int)
	ordered = make([]freqPair, 0, 64)
	for _, r := range []byte(input) {
		if _, ok := freqTable[r]; !ok {
			freqTable[r] = 0
		}
		freqTable[r]++
	}

	for k, v := range freqTable {
		ordered = append(ordered, freqPair{char: k, freq: v})
	}

	sort.SliceStable(ordered, func(i int, j int) bool {
		return ordered[i].freq < ordered[j].freq
	})

	return
}

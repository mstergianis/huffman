package huffman

import (
	"fmt"
	"sort"
)

func Encode(input []byte) ([]byte, error) {
	for i, r := range input {
		if r > 127 || r < 0 {
			return nil, fmt.Errorf("error: encountered a non-ascii character %c at position %d", r, i)
		}
	}
	ordered := computeFreqTable(input)

	tree := NewNode(ordered)

	bs := &BitStringWriter{}
	bs.WriteContentLength(uint32(len(input)))
	tree.WriteBytes(bs)

	for _, b := range []byte(input) {
		bytes, bitWidth := tree.Search(b)
		if bitWidth == -1 {
			return nil, fmt.Errorf("error: cannot find the byte %c in the tree", b)
		}
		bs.WriteBytes(bytes, bitWidth)
	}

	return bs.Bytes(), nil
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

type freqPair struct {
	char byte
	freq int
}

func (f freqPair) Freq() int {
	return f.freq
}

func (f freqPair) String() string {
	return fmt.Sprintf("(%q, %d)", string(f.char), f.freq)
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

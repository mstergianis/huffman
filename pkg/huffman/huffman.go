package huffman

import (
	"fmt"
	"sort"
	"strings"
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

type Node struct {
	freq     int
	freqPair *freqPair
	left     *Node
	right    *Node
}

type Frequentable interface {
	Freq() int
}

func NewNode(ordered []freqPair) *Node {
	nodes := make([]Frequentable, 0, 64)

	for _, o := range ordered {
		nodes = append(nodes, o)
	}

	for len(nodes) > 1 {
		var newNode = Node{}

		switch left := nodes[0].(type) {
		case freqPair:
			{
				newNode.left = &Node{
					freq:     left.freq,
					freqPair: &left,
				}
			}
		case *Node:
			{
				newNode.left = left
			}
		default:
			panic(fmt.Sprintf("unimplemented type %T: %v", left, left))
		}

		switch right := nodes[1].(type) {
		case freqPair:
			{
				newNode.right = &Node{
					freq:     right.freq,
					freqPair: &right,
				}
			}
		case *Node:
			{
				newNode.right = right
			}
		default:
			panic(fmt.Sprintf("unimplemented type %T: %v", right, right))
		}

		newNode.freq = newNode.left.Freq() + newNode.right.Freq()
		nodes = nodes[1:]
		nodes[0] = &newNode
		sort.SliceStable(nodes, func(i, j int) bool {
			return nodes[i].Freq() < nodes[j].Freq()
		})
	}

	head := nodes[0].(*Node)
	return head
}

func (n *Node) Search(b byte) ([]byte, int) {
	if n == nil {
		return nil, -1
	}
	if n.freqPair != nil {
		if n.freqPair.b == b {
			return []byte{0}, 0
		}
		return nil, -1
	}

	if n.right != nil {
		rightBytes, rightBitWidth := n.right.Search(b)
		if rightBitWidth >= 0 {
			adjustedBitWidth := rightBitWidth - (len(rightBytes)-1)*8
			if adjustedBitWidth >= 8 {
				rightBytes = append(rightBytes, 0)
				adjustedBitWidth -= 8
			}

			modifyingByte := &rightBytes[len(rightBytes)-1]
			*modifyingByte = *modifyingByte | (1 << adjustedBitWidth)
			return rightBytes, rightBitWidth + 1
		}
	}

	if n.left != nil {
		leftBytes, leftBitWidth := n.left.Search(b)
		if (leftBitWidth - (len(leftBytes)-1)*8) >= 8 {
			leftBytes = append(leftBytes, 0)
		}
		if leftBitWidth >= 0 {
			return leftBytes, leftBitWidth + 1
		}
	}

	return nil, -1
}

func (n *Node) String() string {
	s := &strings.Builder{}
	fmt.Fprintf(s, "(%d =>", n.freq)
	if n.freqPair != nil {
		fmt.Fprintf(s, " f%s)", n.freqPair)
		return s.String()
	}

	if n.left != nil {
		fmt.Fprintf(s, " *l[%s]", n.left)
	}

	if n.right != nil {
		fmt.Fprintf(s, " *r[%s]", n.right)
	}

	fmt.Fprint(s, ")")
	return s.String()
}

func (n *Node) printTreeString() string {
	s := &strings.Builder{}
	fmt.Fprintf(s, "(%d =>", n.freq)

	if n.freqPair != nil {
		fmt.Fprintf(s, " f%s)", n.freqPair)
		return s.String()
	}

	return s.String()
}

func (n *Node) Freq() int {
	return n.freq
}

// Bytes encodes the tree as an array of bytes
//
// Frequencies are omitted to save data and because they are not essential for
// recovering the compressed text.
//
// The tree is encoded depth-first, to make deserializing easier.
func (n *Node) Bytes() []byte {
	// grammar:
	//   node                       = (leftChild rightChild) | freqPair .
	//   leftChild         (2 bits) = leftBitString node .
	//   rightChild        (2 bits) = rightBitString node .
	//   freqPair         (10 bits) = freqPairBitstring byte .
	//   freqPairBitString (2 bits) = "01" .
	//   leftBitString     (2 bits) = "11" .
	//   rightBitString    (2 bits) = "10" .
	//   binaryDigit                = "1" | "0" .
	//
	// Notice that the grammar does not prescribe that all elements occupy a
	// full byte. But that all valid inputs would start with the bit width used.
	// Thus you should be able to write a program that correctly allocates a
	// tree based on the input.
	//
	// One could probably analyze the character set and eke out a bit more data
	// savings. But I'm going to be saving the whole byte (rune... char...
	// whatever).

	bs := &BitStringWriter{}
	nodeToBytes(n, bs)

	return bs.Bytes()
}

func nodeToBytes(n *Node, bs *BitStringWriter) {
	if n == nil {
		return
	}

	if n.freqPair != nil {
		bs.Write(byte(CONTROL_BIT_FREQ_PAIR), 2)
		bs.Write(n.freqPair.b, 8)
	}

	if n.left != nil {
		bs.Write(byte(CONTROL_BIT_LEFT), 2)
		nodeToBytes(n.left, bs)
	}

	if n.right != nil {
		bs.Write(byte(CONTROL_BIT_RIGHT), 2)
		nodeToBytes(n.right, bs)
	}

	return
}

// printTree is a debugging tool
func printTree(tree *Node) {
	printTreeWithDepth(0, tree)
}

// printTreeWithDepth is a debugging tool
func printTreeWithDepth(depth int, tree *Node) {
	if tree == nil {
		return
	}
	// print the depth, then the node
	fmt.Printf("%s%s\n", strings.Repeat(" ", depth*2), tree.printTreeString())
	printTreeWithDepth(depth+1, tree.left)
	printTreeWithDepth(depth+1, tree.right)
}

func maxDepth(depth int, tree *Node) int {
	//            a               => 1 (spaces = nil)
	//       /         \
	//      b           c      => 2 (spaces = 11)
	//    /   \       /   \
	//   d     e     f     g   => 3 (spaces = 5)
	//  / \   / \   / \   / \
	// h   i j   k l   m n   o => 4 (spaces = 3, 1)

	//             b          c
	//         / /            c
	//        / /             c
	//       / /         /         \
	//      d           e           f
	//    /   \       /   \       /   \
	//   h     i     j     k     l     m
	//  / \   / \   / \   / \   / \   / \
	// p   q r   s t   u v   w x   y z   a
	if tree == nil {
		return depth
	}

	leftDepth := 0
	if tree.left != nil {
		leftDepth = maxDepth(leftDepth, tree.left)
	}

	rightDepth := 0
	if tree.right != nil {
		rightDepth = maxDepth(rightDepth, tree.right)
	}

	if leftDepth > rightDepth {
		return leftDepth + 1
	}
	return rightDepth + 1
}

type BitStringWriter struct {
	// 0000 0000
	// ^
	buffer []byte
	offset int
}

func (bs *BitStringWriter) WriteBytes(bytes []byte, w int) {
	onByte := 0
	widthRemaining := w
	for widthRemaining >= 8 {
		bs.Write(bytes[onByte], 8)
		onByte++
		widthRemaining -= 8
	}
	if widthRemaining > 0 {
		bs.Write(bytes[onByte], widthRemaining)
	}
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
	// write the whole byte
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

func (bs *BitStringWriter) String() string {
	return string(bs.buffer)
}

func (bs *BitStringWriter) Bytes() []byte {
	return bs.buffer
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

func (bs *BitStringWriter) addByte() {
	bs.buffer = append(bs.buffer, 0)
	bs.offset = 0
}

func (bs *BitStringWriter) writeToLastByte(b byte, w int) {
	bs.buffer[len(bs.buffer)-1] = bs.buffer[len(bs.buffer)-1] | (b << (7 - bs.offset - w + 1))
	bs.offset += w
}

type freqPair struct {
	b    byte
	freq int
}

func (f freqPair) Freq() int {
	return f.freq
}

func (f freqPair) String() string {
	return fmt.Sprintf("('%s', %d)", string(f.b), f.freq)
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
		ordered = append(ordered, freqPair{b: k, freq: v})
	}

	sort.SliceStable(ordered, func(i int, j int) bool {
		return ordered[i].freq < ordered[j].freq
	})

	return
}

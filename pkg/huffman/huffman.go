package huffman

import (
	"fmt"
	"log"
	"sort"
	"strings"
)

func Huffman(input string) (*Node, []byte) {
	ordered := computeFreqTable(input)

	tree := NewNode(ordered)
	printTree(tree)

	bs := &BitStringWriter{}

	for _, r := range input {
		byte, bitWidth := tree.Search(r)
		if bitWidth == -1 {
			log.Fatalf("we can't find the rune %s in the tree", string(r))
		}
		fmt.Printf("%q => ", string(r))
		bs.Write(byte, bitWidth)
	}

	return tree, bs.Bytes()
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

func (n *Node) Search(r rune) (byte, int) {
	if n.freqPair != nil {
		if n.freqPair.r == r {
			return 0, 0
		}
		return 0, -1
	}

	if n.left == nil || n.right == nil {
		return 0, -1
	}

	leftByte, leftBitWidth := n.left.Search(r)
	if leftBitWidth >= 0 {
		return leftByte, leftBitWidth + 1
	}

	rightByte, rightBitWidth := n.right.Search(r)
	if rightBitWidth >= 0 {
		return rightByte | (1 << rightBitWidth), rightBitWidth + 1
	}

	return 0, -1
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
// The tree is encoded depth-first, to make deserializing easier.
func (n *Node) Bytes() []byte {
	// serialize the first node
	return nil
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
	switch w {
	case 0:
		return 0
	case 1:
		return b & 0b0000_0001
	case 2:
		return b & 0b0000_0011
	case 3:
		return b & 0b0000_0111
	case 4:
		return b & 0b0000_1111
	case 5:
		return b & 0b0001_1111
	case 6:
		return b & 0b0011_1111
	case 7:
		return b & 0b0111_1111
	}
	panic("computeRightByte: encountered a width greater than 7")
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
	r    rune
	freq int
}

func (f freqPair) Freq() int {
	return f.freq
}

func (f freqPair) String() string {
	return fmt.Sprintf("('%s', %d)", string(f.r), f.freq)
}

func computeFreqTable(input string) (ordered []freqPair) {
	freqTable := make(map[rune]int)
	ordered = make([]freqPair, 0, 64)
	for _, r := range input {
		if _, ok := freqTable[r]; !ok {
			freqTable[r] = 0
		}
		freqTable[r]++
	}

	for k, v := range freqTable {
		ordered = append(ordered, freqPair{r: k, freq: v})
	}

	sort.SliceStable(ordered, func(i int, j int) bool {
		return ordered[i].freq < ordered[j].freq
	})

	return
}

package huffman

import (
	"fmt"
	"sort"
	"strings"
)

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
		if n.freqPair.char == b {
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
		bs.Write(n.freqPair.char, 8)
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

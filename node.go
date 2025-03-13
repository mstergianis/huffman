package main

import (
	"fmt"
	"sort"
	"strings"
)

type node struct {
	freq     int
	freqPair *freqPair
	left     *node
	right    *node
}

type Frequentable interface {
	Freq() int
}

func makeTree(ordered []freqPair) *node {
	nodes := make([]Frequentable, 0, 64)

	for _, o := range ordered {
		nodes = append(nodes, o)
	}

	for len(nodes) > 1 {
		var newNode = node{}

		switch left := nodes[0].(type) {
		case freqPair:
			{
				newNode.left = &node{
					freq:     left.freq,
					freqPair: &left,
				}
			}
		case *node:
			{
				newNode.left = left
			}
		default:
			panic(fmt.Sprintf("unimplemented type %T: %v", left, left))
		}

		switch right := nodes[1].(type) {
		case freqPair:
			{
				newNode.right = &node{
					freq:     right.freq,
					freqPair: &right,
				}
			}
		case *node:
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

	head := nodes[0].(*node)
	return head
}

func (n *node) Search(r rune) (string, bool) {
	if n.freqPair != nil {
		if n.freqPair.r == r {
			return "", true
		}
		return "", false
	}

	if n.left == nil || n.right == nil {
		return "", false
	}

	leftS, leftB := n.left.Search(r)
	if leftB {
		return "0" + leftS, leftB
	}

	rightS, rightB := n.right.Search(r)
	if rightB {
		return "1" + rightS, rightB
	}

	return "", false
}

func (n *node) String() string {
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

func (n *node) PrintTreeString() string {
	s := &strings.Builder{}
	fmt.Fprintf(s, "(%d =>", n.freq)
	if n.freqPair != nil {
		fmt.Fprintf(s, " f%s)", n.freqPair)
		return s.String()
	}

	// if n.left != nil {
	// 	fmt.Fprintf(s, " *l[%s]", n.left)
	// }
	// if n.right != nil {
	// 	fmt.Fprintf(s, " *r[%s]", n.right)
	// }
	// fmt.Fprint(s, ")")
	return s.String()
}

func (n *node) Freq() int {
	return n.freq
}

func printTree(tree *node) {
	printTreeWithDepth(0, tree)
}

func printTreeWithDepth(depth int, tree *node) {
	if tree == nil {
		return
	}
	// print the depth, then the node
	fmt.Printf("%s%s\n", strings.Repeat(" ", depth*2), tree.PrintTreeString())
	printTreeWithDepth(depth+1, tree.left)
	printTreeWithDepth(depth+1, tree.right)
}

func maxDepth(depth int, tree *node) int {
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

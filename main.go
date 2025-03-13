package main

import (
	"bytes"
	"fmt"
	"log"
	"sort"
	"strings"
)

func main() {
	input := "hello world"
	huffman(input)
	return
}

type bitstringwriter struct {
	b        bytes.Buffer
	overflow *byte
}

func (b *bitstringwriter) Write(bitstring []byte) (int, error) {
	return 0, nil
}

type Frequentable interface {
	Freq() int
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

func huffman(input string) (mapping map[rune]int, output []byte) {

	ordered := computeFreqTable(input)

	tree := makeTree(ordered)
	printTree(tree)

	for _, r := range input {
		bitstring, ok := tree.Search(r)
		if !ok {
			log.Fatalf("we can't find the rune %s in the tree", string(r))
		}
		fmt.Printf("%s => %s\n", string(r), bitstring)
	}

	return nil, nil
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

type node struct {
	freq     int
	freqPair *freqPair
	left     *node
	right    *node
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

const WILD_WILD_WEST = `
(Alright man, fine)
("Wa-wa")
Uhh (doo-doo-doo-doo)
Wicki-wild wild (doo-doo-doo-doo-doo)
Wicki-wicki-wild
Wicki-wild
Wicki-wicki wild wild West
Jim West, desperado
Rough rider, no you don't want nada
None of this, six-gunnin' this, brother runnin' this
Buffalo soldier, look, it's like I told ya
Any damsel that's in distress
Be outta that dress when she meet Jim West
Rough neck so go check the law and abide
Watch your step or flex and get a hole in your side
Swallow your pride, don't let your lip react
You don't wanna see my hand where my hip be at
With Artemus, from the start of this, runnin' the game
James West, tamin' the West, so remember the name
Now who ya gonna call?
Not the GB's
Now who you gonna call?
J Dub and A.G
If you ever riff with either one of us
Break out, before you get bum-rushed, at the
Wild wild west (when I roll into the)
Wild wild west (when I stroll into the)
Wild wild west (when I bounce into the)
Wild wild west
We're goin' straight to the wild wild west (wild wild west, the wild wild west)
We're goin' straight to the wild wild west (wild wild west)
Now, now, now
Now once upon a time in the West
Mad man lost his damn mind in the West
Loveless, kidnap a dime, nothin' less
Now I must, put his behind to the test (can you feel me?)
Then through the shadows, in the saddle, ready for battle
Bring all your boys in, here come the poison
Behind my back, all that riffin' ya did
Front and center, now where your lip at kid?
Who dat is? A mean brother, bad for your health
Lookin' damn good, though, if I could say it myself
Told me Loveless is a mad man, but I don't fear that
He got mad weapons too? Ain't tryna hear that
Tryin' to bring down me, the champion?
When y'all clowns gon' see that it can't be done
Understand me, son, I'm the slickest they is
I'm the quickest they is (yeah)
Did I say I'm the slickest they is?
So if you barkin' up the wrong tree, we comin'
Don't be startin' nothin', me and my partner gonna
Test your chest, Loveless
Can't stand the heat, then get out the wild, wild, West
We're goin' straight to the wild wild West (when I roll into the)
We're goin' straight (when I stroll into the)
To (when I bounce into the)
The wild wild west
We're goin' straight (straight) to (to)
The wild wild West (the wild wild West)
We're goin' straight (straight) to (to)
The wild wild West (the wild wild West)
(The wild wild West)
Yeah
Can you feel it? C'mon, c'mon
Yeah (breakdown)
Keep it moving, keep it moving (breakdown)
Ooh, yeah
To any outlaw tryin' to draw, thinkin' you're bad
Any drawin' on West, best with a pen and a pad
Don't even think about it, six gun, weighin a ton
Ten paces and turn (one... two... three...), just for fun, son
Up 'til sundown, rollin' around
See where the bad guys are to be found and make 'em lay down
They're defenders of the West
Crushin' all pretenders in the West
Don't mess with us, 'cause we in the
Wild wild west (when I roll into the)
Wild wild west (when I stroll into the)
Wild wild west (when I bounce into the)
Wild wild west
We're goin' straight (straight) to (to)
The wild wild West (the wild wild West)
We're goin' straight (straight) to (to)
The wild wild West (the wild wild West)
We're goin' straight (straight) to (to)
The wild wild West (the wild wild West)
We're goin' straight (straight) to (to)
The wild wild West (the wild wild West)
Come on
The wild wild west (when I roll into the) come on
Wild wild west (when I stroll into the) we're goin' straight to (to)
The wild wild west (the wild wild west)
The wild wild west (whoo! Uhh)
The wild wild west (ha hah, ha hah)
The wild wild west (uhh)
The wild wild west (I done done it again y'all, done done it again)
The wild wild west (ha hah, ha hah)
The wild wild west (big Will, Dru Hill, uh)
The wild wild west (big Will, Dru Hill)
The wild wild west (ha hah, ha hah)
The wild wild west
The wild wild west (uhh)
The wild wild west (one time)
The wild wild west (uhh)
The wild wild west (bring in the heat, bring in the heat)
The wild wild west (what? Ha hah, ha hah)
The wild wild west (whoo, wild wild, wicki-wild)
The wild wild west (wick wild wild wild, wa-wicki wild wild)
The wild wild west (wickidy-wick wild wild wild)
The wild wild west
The wild wild west
The wild wild west (uhh, uhh)
The wild wild west (can't stop the bumrush)
The wild wild...
The wild wild west
`

package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/mstergianis/huffman/pkg/huffman"
)

var programName string

func main() {
	args := os.Args
	var err error
	programName, err = shift(&args)
	check(err)
	operatingMode, err := shift(&args)
	check(err)
	var (
		inputFile  string
		outputFile string
	)
	for len(args) > 0 {
		arg, err := shift(&args)
		check(err)
		switch arg {
		case "-i":
			inputFile, err = shift(&args)
			check(err)
		case "-o":
			outputFile, err = shift(&args)
			check(err)
		default:
			usage()
		}
	}

	input, err := os.ReadFile(inputFile)
	check(err)

	var writeToOutput func(w io.Writer) error
	switch operatingMode {
	case "encode":
		{
			treeAsBytes, contentAsBytes, err := huffman.Encode(input)
			check(err)

			writeToOutput = func(w io.Writer) error {
				n, err := w.Write(treeAsBytes)
				if err != nil {
					return err
				}
				if n < len(treeAsBytes) {
					return fmt.Errorf("error: wrote fewer bytes (%d) than the tree's length (%d)", n, len(treeAsBytes))
				}

				n, err = w.Write(contentAsBytes)
				if err != nil {
					return err
				}
				if n < len(contentAsBytes) {
					return fmt.Errorf("error: wrote fewer bytes (%d) than the content's length (%d)", n, len(contentAsBytes))
				}

				return nil
			}
		}
	case "decode":
		{
			panic("unimplemented")
		}
	}

	f, err := os.OpenFile(outputFile, os.O_CREATE|os.O_RDWR, 0644)
	check(err)
	defer f.Close()

	err = writeToOutput(f)
	check(err)

	return
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s COMMAND\n", programName)
	fmt.Fprintf(os.Stderr, "Available commands:\n")
	fmt.Fprintf(os.Stderr, "    encode -i INPUT-FILE -o OUTPUT-FILE\n")
	fmt.Fprintf(os.Stderr, "    decode -i INPUT-FILE -o OUTPUT-FILE\n")
	os.Exit(1)
}

func check(err error) {
	if errors.Is(err, shiftErr) {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		usage()
	}
	if err != nil {
		panic(err.Error())
	}
}

var (
	shiftErr = errors.New("error: expected an argument but did not find one")
)

func shift(ss *[]string) (string, error) {
	if len(*ss) < 1 {
		return "", shiftErr
	}
	s := (*ss)[0]
	*ss = (*ss)[1:]
	return s, nil
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

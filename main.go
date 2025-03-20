package main

import (
	"errors"
	"fmt"
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

	var modeFunc func([]byte) ([]byte, error)
	switch operatingMode {
	case "encode":
		{
			modeFunc = huffman.Encode
		}
	case "decode":
		{
			modeFunc = huffman.Decode
		}
	}

	f, err := os.OpenFile(outputFile, os.O_CREATE|os.O_RDWR, 0644)
	check(err)
	defer f.Close()

	contents, err := modeFunc(input)
	check(err)

	n, err := f.Write(contents)
	check(err)
	if n < len(contents) {
		panic(fmt.Sprintf("error: wrote fewer bytes (%d) than the content's length (%d)", n, len(contents)))
	}

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

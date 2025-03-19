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
			contents, err := huffman.Encode(input)
			check(err)

			writeToOutput = func(w io.Writer) error {
				n, err := w.Write(contents)
				if err != nil {
					return err
				}
				if n < len(contents) {
					return fmt.Errorf("error: wrote fewer bytes (%d) than the content's length (%d)", n, len(contents))
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

package main

import (
	"flag"
	"fmt"
	"io"
)

type LineLengthCommand struct {
	cmdOptions *flag.FlagSet
	delim      *string
}

func (LineLengthCommand) Name() string {
	return "llen"
}

func (LineLengthCommand) Desc() string {
	return "Prints length of each input line."
}

func (llc *LineLengthCommand) Init() {
	llc.cmdOptions = flag.NewFlagSet(llc.Name(), flag.ExitOnError)
	llc.delim = llc.cmdOptions.String("delim", "\n", "record delimiter.")
}

func (llc *LineLengthCommand) Interact(args []string, stdin io.Reader, stdout io.Writer, stderr io.Writer) error {
	if llc.cmdOptions.Parse(args) != nil {
		return CommandNotInitialized
	}

	bInOut, _ := ToBuffered(stdin, stdout, stderr)
	for {
		line, err := bInOut.Reader.ReadString([]byte(*(llc.delim))[0])
		if err != nil {
			break
		} else {
			line = line[:len(line)-1]
			WriteAndFlush(bInOut.Writer, fmt.Sprintf("%s:%d\n", line, len(line)))
		}
	}

	return nil
}

func (llc *LineLengthCommand) Help(stderr io.Writer) {
	llc.cmdOptions.SetOutput(stderr)
	llc.cmdOptions.Usage()
}

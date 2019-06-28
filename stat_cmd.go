package main

import (
	"flag"
	"fmt"
	"io"
	"strconv"
)

// Stat command -  Integer statistics for sidin integer streams
type StatCommand struct {
	delim      *string
	cmdOptions *flag.FlagSet
}

func (sc *StatCommand) Name() string {
	return "stat"
}

func (sc *StatCommand) Desc() string {
	return "Integer statistics for STDIN integer streams"
}

func (sc *StatCommand) Init() {
	sc.cmdOptions = flag.NewFlagSet(sc.Name(), flag.ExitOnError)
	sc.delim = sc.cmdOptions.String("delim", "\n", "record delimiter.")
}

func (sc *StatCommand) Help(stderr io.Writer) {
	sc.cmdOptions.SetOutput(stderr)
	sc.cmdOptions.Usage()
}

func (sc *StatCommand) Interact(args []string, stdin io.Reader, stdout io.Writer, stderr io.Writer) error {
	if sc.cmdOptions.Parse(args) != nil {
		return CommandNotInitialized
	}

	bInOut, _ := ToBuffered(stdin, stdout, stderr)
	sum := 0
	size := 0
	for {
		line, err := bInOut.ReadString([]byte(*(sc.delim))[0])
		if err != nil {
			break
		} else {
			line = line[:len(line)-1]
			value, err := strconv.Atoi(line)
			if err != nil {
				continue
			}
			sum += value
			size++
		}
	}
	WriteAndFlush(bInOut.Writer, fmt.Sprintf("sum: %d, size: %d, avg: %f\n", sum, size, float32(sum)/float32(size)))
	return nil
}

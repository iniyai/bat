package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"strconv"
)

// Stat command -  Integer statistics for sidin integer streams
type StatCommand struct {
	delim *string
}

func (sc *StatCommand) Desc() string {
	return "Integer statistics for STDIN integer streams"
}

func (sc *StatCommand) Init(args []string) error {
	flagSet := flag.NewFlagSet("stat", flag.ExitOnError)
	sc.delim = flagSet.String("delim", "\n", "record delimiter.")
	return flagSet.Parse(args)
}

func (sc *StatCommand) Help(stderr io.Writer) {
	bOut := bufio.NewWriter(stderr)
	WriteAndFlush(bOut, "Stat command help")
}

func (sc *StatCommand) Interact(stdin io.Reader, stdout io.Writer, stderr io.Writer) error {
	bInOut := bufio.NewReadWriter(bufio.NewReader(stdin), bufio.NewWriter(stdout))
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

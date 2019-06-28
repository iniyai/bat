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
	floatMode  *bool
	base       *int
	cmdOptions *flag.FlagSet
}

func (sc *StatCommand) Name() string {
	return "stat"
}

func (sc *StatCommand) Desc() string {
	return "Numerical statistics for STDIN reals"
}

func (sc *StatCommand) Init() {
	sc.cmdOptions = flag.NewFlagSet(sc.Name(), flag.ExitOnError)
	sc.delim = sc.cmdOptions.String("delim", "\n", "record delimiter.")
	sc.base = sc.cmdOptions.Int("base", 10, "numerical base of input")
	sc.floatMode = sc.cmdOptions.Bool("float", false, "float mode")
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
	var sumF float64 = 0.0
	var sumI int64 = 0
	size := 0
	for {
		line, err := bInOut.ReadString([]byte(*(sc.delim))[0])
		if err != nil {
			break
		} else {
			line = line[:len(line)-1]
			if *sc.floatMode {
				value, err := strconv.ParseFloat(line, 64)
				if err != nil {
					continue
				}
				sumF += value
			} else {
				value, err := strconv.ParseInt(line, *sc.base, 64)
				if err != nil {
					continue
				}
				sumI += value
			}
			size++
		}
	}
	if *sc.floatMode {
		WriteAndFlush(bInOut.Writer, fmt.Sprintf("sum: %f, size: %d, avg: %f\n", sumF, size, sumF/float64(size)))
	} else {
		WriteAndFlush(bInOut.Writer, fmt.Sprintf("sum: %d, size: %d, avg: %f\n", sumI, size, float64(sumI)/float64(size)))
	}
	return nil
}

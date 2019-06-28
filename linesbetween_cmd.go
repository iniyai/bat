package main

import (
	"flag"
	"io"
	"strconv"
)

type LinesBetweenCommand struct {
	start, end *int
	delim      *string
	cmdOptions *flag.FlagSet
}

func (*LinesBetweenCommand) Name() string {
	return "lbw"
}

func (*LinesBetweenCommand) Desc() string {
	return "Extracts Lines between given start-index end end-index."
}

func (lbw *LinesBetweenCommand) Init() {
	lbw.cmdOptions = flag.NewFlagSet(lbw.Name(), flag.ExitOnError)
	lbw.start = lbw.cmdOptions.Int("start", 0, "start index")
	lbw.end = lbw.cmdOptions.Int("end", -1, "end index")
	lbw.delim = lbw.cmdOptions.String("delim", "\n", "record delimiter")
}

func (lbw *LinesBetweenCommand) Interact(args []string, stdin io.Reader, stdout io.Writer, stderr io.Writer) error {
	if lbw.cmdOptions.Parse(args) != nil {
		return CommandNotInitialized
	}

	bInOut, bErr := ToBuffered(stdin, stdout, stderr)
	if *lbw.start < 0 {
		WriteAndFlush(bErr, "start index is negative: "+strconv.Itoa(*lbw.start))
		return IllegalCommandArguments
	}
	count := 0
	for {
		line, err := bInOut.ReadString([]byte(*(lbw.delim))[0])
		if err != nil {
			break
		} else {
			count++
			// in Range - print
			if *lbw.start <= count && (*lbw.end < 0 || *lbw.end >= count) {
				WriteAndFlush(bInOut.Writer, line)
			}

			// If we don't need to care about next lines, then break
			if *lbw.end > 0 && count > *lbw.end {
				break
			}
		}
	}
	return nil
}

func (lbw *LinesBetweenCommand) Help(stderr io.Writer) {
	lbw.cmdOptions.SetOutput(stderr)
	lbw.cmdOptions.Usage()
}

package main

import (
	"flag"
	"io"
)

type HistogramCommand struct {
	span            *int
	horizontalChart *bool
	maxHeight *int
//	floatMode *bool
	cmdOptions      *flag.FlagSet
}

func (HistogramCommand) Name() string {
	return "hist"
}

func (HistogramCommand) Desc() string {
	return "Prints histogram of input values"
}

func (hc *HistogramCommand) Init() {
	hc.cmdOptions = flag.NewFlagSet(hc.Name(), flag.ExitOnError)
	hc.span = hc.cmdOptions.Int("span", 5, "Histogram bucket size")
	hc.horizontalChart = hc.cmdOptions.Bool("hc", true, "Horizontal chart orientation")
	hc.maxHeight = hc.cmdOptions.Int("maxHeight", 50, "Maximum height of histograms")
}

func (HistogramCommand) Interact(args []string, stdin io.Reader, stdout io.Writer, stderr io.Writer) error {
	panic("implement me")
}

func (hc *HistogramCommand) Help(stderr io.Writer) {
	hc.cmdOptions.SetOutput(stderr)
	hc.cmdOptions.Usage()
}

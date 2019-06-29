package main

import (
	"flag"
	"fmt"
	"io"
	tmap "github.com/emirpasic/gods/maps/treemap"
	"strconv"
)

type HistogramCommand struct {
	bins            *int
	horizontalChart *bool
	maxHeight *int
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
	hc.bins = hc.cmdOptions.Int("bins", 10, "Number of Histogram buckets")
	hc.horizontalChart = hc.cmdOptions.Bool("hc", true, "Horizontal chart orientation")
	hc.maxHeight = hc.cmdOptions.Int("maxHeight", 50, "Maximum height of histograms")
}

func (hc *HistogramCommand) Interact(args []string, stdin io.Reader, stdout io.Writer, stderr io.Writer) error {
	bInOut, _ := ToBuffered(stdin, stdout, stderr)
	fmap := tmap.NewWithIntComparator()
	for {
		line, err := bInOut.Reader.ReadString('\n')
		if err != nil {
			break
		} else {
			line = line[:len(line)-1]
			value, err := strconv.Atoi(line)
			if err != nil {
				continue
			}
			freq, found := fmap.Get(value)
			if found {
				fmap.Put(value, freq.(int) + 1)
			} else {
				fmap.Put(value, 1)
			}
		}
	}

	if fmap.Empty() {
		return nil
	}

	min, _ := fmap.Min()
	max, _ := fmap.Max()

	// Size and Span
	hGram := Histogram{bins:*hc.bins, max:max.(int), min:min.(int), hMap:make(map[int]int)}
	fmap.Each(hGram.add)
	fmt.Println(hGram.hMap)
	fmt.Println(fmap)

	return nil
}

func (hc *HistogramCommand) Help(stderr io.Writer) {
	hc.cmdOptions.SetOutput(stderr)
	hc.cmdOptions.Usage()
}

type Histogram struct {
	bins, min, max int
	hMap map[int]int
}

func (hGram *Histogram) size()  int {
	return hGram.max - hGram.min + 1
}

func (hGram *Histogram) span() int {
	return hGram.size() / hGram.bins
}

func (hGram *Histogram) findBucket(score int) int {
	return score /  hGram.span()
}

func (hGram *Histogram) add(score, freq interface{})  {
	hGram.hMap[hGram.findBucket(score.(int))] += freq.(int)
}




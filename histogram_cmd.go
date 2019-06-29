package main

import (
	"bufio"
	"flag"
	"fmt"
	tmap "github.com/emirpasic/gods/maps/treemap"
	"github.com/emirpasic/gods/utils"
	"io"
	"math"
	"strconv"
	"strings"
)

type HistogramCommand struct {
	bins            *int64
	horizontalChart *bool
	maxHeight       *int64
	cmdOptions      *flag.FlagSet
}

func (HistogramCommand) Name() string {
	return "hist"
}

func (HistogramCommand) Desc() string {
	return "Prints histogram input stream of positive integers"
}

func (hc *HistogramCommand) Init() {
	hc.cmdOptions = flag.NewFlagSet(hc.Name(), flag.ExitOnError)
	hc.bins = hc.cmdOptions.Int64("bins", 10, "Number of Histogram buckets")
	hc.horizontalChart = hc.cmdOptions.Bool("hc", true, "Horizontal chart orientation")
	hc.maxHeight = hc.cmdOptions.Int64("maxHeight", 50, "Maximum height of histograms")
}

func (hc *HistogramCommand) Interact(args []string, stdin io.Reader, stdout io.Writer, stderr io.Writer) error {
	if hc.cmdOptions.Parse(args) != nil {
		return CommandNotInitialized
	}

	bInOut, _ := ToBuffered(stdin, stdout, stderr)

	fmap := make(map[int64]int64)
	var min int64 = math.MaxInt64
	var max int64 = math.MinInt64
	for {
		line, err := bInOut.Reader.ReadString('\n')
		if err != nil {
			break
		} else {
			line = line[:len(line)-1]
			value, err := strconv.ParseInt(line, 10, 32)
			if err != nil {
				continue
			}
			if value < min {
				min = value
			}

			if value > max {
				max = value
			}
			fmap[value]++
		}
	}

	if len(fmap) == 0 {
		return nil
	}

	// Size and Span
	hGram := New(min, max, *hc.bins)
	for score, freq := range fmap {
		hGram.add(score, freq)
	}
	hGram.plot(bInOut.Writer, *hc.horizontalChart, *hc.maxHeight)
	return nil
}

func (hc *HistogramCommand) Help(stderr io.Writer) {
	hc.cmdOptions.SetOutput(stderr)
	hc.cmdOptions.Usage()
}

// Histogram type
type Histogram struct {
	bins, min, max, size, span, maxFreq, totalFreq int64
	bucketMap                                      tmap.Map
}

func New(min, max, bins int64) Histogram {
	hist := Histogram{min: min, max: max + 1, bins: bins}
	hist.size = hist.max - hist.min + 1
	if hist.size/hist.bins == 0 {
		hist.span = 1
	} else {
		hist.span = hist.size / hist.bins
	}
	hist.bucketMap = *tmap.NewWith(utils.Int64Comparator)
	bs := hist.min
	var i int64 = 0
	for ; i < hist.bins; i++ {
		hist.bucketMap.Put(bs, int64(0))
		bs += hist.span
	}
	hist.maxFreq = int64(math.MinInt64)
	hist.totalFreq = int64(0)
	return hist
}

func (hGram *Histogram) add(score, freq int64) {
	bucketKey, bucketValue := hGram.bucketMap.Floor(score)
	if bucketKey != nil {
		hGram.totalFreq += freq
		newFreq := bucketValue.(int64) + freq
		hGram.bucketMap.Put(bucketKey, newFreq)
		if newFreq > hGram.maxFreq {
			hGram.maxFreq = newFreq
		}
	} else {
		panic("unknown histogram key!")
	}
}

func (hGram *Histogram) plot(writer *bufio.Writer, horizontal bool, height int64) {
	if horizontal {
		lastBucket, _ := hGram.bucketMap.Max()
		horizontalBucketStartWidth := int(math.Ceil(math.Log10(float64(lastBucket.(int64))))) + 1
		horizontalBucketEndWidth := int(math.Ceil(math.Log10(float64(lastBucket.(int64)+hGram.span)))) + 1
		fmtStr := fmt.Sprintf("[%%%dd -%%%dd) | %%-%ds | (%%d - %%.2f%%%%)\n", horizontalBucketStartWidth,
			horizontalBucketEndWidth, height)
		hGram.bucketMap.Each(func(key interface{}, value interface{}) {
			histPercent := float32(value.(int64)) / float32(hGram.totalFreq) * 100
			histLen := (value.(int64) * height) / hGram.maxFreq
			WriteAndFlush(writer,
				fmt.Sprintf(fmtStr, key, key.(int64)+hGram.span, strings.Repeat("*", int(histLen)), value, histPercent))
		})
	} else {
		panic("vertical map not yet implemented")
	}
}

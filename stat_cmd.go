package main

import (
	"flag"
	"fmt"
	"github.com/emirpasic/gods/maps/treemap"
	"github.com/emirpasic/gods/utils"
	"io"
	"math"
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
	return "Numerical statistics for STDIN reals."
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
	iStats := &intStats{}
	rStats := &realStats{}
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

				rStats.update(value)
			} else {
				value, err := strconv.ParseInt(line, *sc.base, 64)
				if err != nil {
					continue
				}

				iStats.update(value)
			}
		}
	}

	if *sc.floatMode && rStats.count() > 0 {
		WriteAndFlush(bInOut.Writer, fmt.Sprintf("sum: %f, size: %d, mean: %f, min: %f, max: %f, mode: %f\n", rStats.sum(),
			rStats.count(), rStats.mean(), rStats.min(), rStats.max(), rStats.mode()))
	} else if iStats.count() > 0 {
		WriteAndFlush(bInOut.Writer, fmt.Sprintf("sum: %d, size: %d, mean: %f, min: %d, max: %d, mode: %d\n", iStats.sum(),
			iStats.count(), iStats.mean(), iStats.min(), iStats.max(), iStats.mode()))
	}
	return nil
}

// Stat utilities
type stats interface {
	update(values ...interface{})

	min() interface{}

	max() interface{}

	sum() interface{}

	median() interface{}

	mode() interface{}

	mean() float64

	stddev() float64

	count() uint64
}

// Int stats
type intStats struct {
	tmap   *treemap.Map
	sumE   int64
	countE uint64
}

func (is *intStats) update(values ...interface{}) {
	if is.tmap == nil {
		is.tmap = treemap.NewWith(utils.Int64Comparator)
	}
	for _, value := range values {
		is.sumE += value.(int64)
		is.countE++
		oldValue, ok := is.tmap.Get(value.(int64))
		if ok {
			is.tmap.Put(value.(int64), oldValue.(uint64)+uint64(1))
		} else {
			is.tmap.Put(value.(int64), uint64(1))
		}
	}
}

func (is *intStats) min() interface{} {
	if is.countE == 0 {
		return math.NaN()
	}

	value, _ := is.tmap.Min()
	return value.(int64)
}

func (is *intStats) max() interface{} {
	if is.countE == 0 {
		return math.NaN()
	}

	value, _ := is.tmap.Max()
	return value.(int64)
}

func (is *intStats) sum() interface{} {
	if is.countE == 0 {
		return math.NaN()
	}

	return is.sumE
}

func (is *intStats) median() interface{} {
	if is.countE == 0 {
		return math.NaN()
	}

	return extractMedianFromFreqencyMap(is.tmap, is.countE)
}

func (is *intStats) mode() interface{} {
	if is.countE == 0 {
		return math.NaN()
	}

	return extractModeFromFreqencyMap(is.tmap)
}

func (is *intStats) mean() float64 {
	if is.countE == 0 {
		return math.NaN()
	}

	return float64(is.sumE) / float64(is.countE)
}

func (is *intStats) stddev() float64 {
	if is.countE == 0 {
		return math.NaN()
	}

	panic("implement me")
}

func (is *intStats) count() uint64 {
	return is.countE
}

// Real stats

type realStats struct {
	tmap   *treemap.Map
	sumE   float64
	countE uint64
}

func (rs *realStats) update(values ...interface{}) {
	if rs.tmap == nil {
		rs.tmap = treemap.NewWith(utils.Float64Comparator)
	}
	for _, value := range values {
		rs.sumE += value.(float64)
		rs.countE++
		oldValue, ok := rs.tmap.Get(value.(float64))
		if ok {
			rs.tmap.Put(value.(float64), oldValue.(uint64)+uint64(1))
		} else {
			rs.tmap.Put(value.(float64), uint64(1))
		}
	}
}

func (rs *realStats) min() interface{} {
	if rs.countE == 0 {
		return math.NaN()
	}

	value, _ := rs.tmap.Min()
	return value.(float64)
}

func (rs *realStats) max() interface{} {
	if rs.countE == 0 {
		return math.NaN()
	}

	value, _ := rs.tmap.Max()
	return value.(float64)
}

func (rs *realStats) sum() interface{} {
	if rs.countE == 0 {
		return math.NaN()
	}

	return rs.sumE
}

func (rs *realStats) median() interface{} {
	if rs.countE == 0 {
		return math.NaN()
	}

	return extractMedianFromFreqencyMap(rs.tmap, rs.countE)
}

func (rs *realStats) mode() interface{} {
	if rs.countE == 0 {
		return math.NaN()
	}

	return extractModeFromFreqencyMap(rs.tmap)
}

func (rs *realStats) mean() float64 {
	if rs.countE == 0 {
		return math.NaN()
	}

	return rs.sumE / float64(rs.countE)
}

func (rs *realStats) stddev() float64 {
	if rs.countE == 0 {
		return math.NaN()
	}

	panic("implement me")
}

func (rs *realStats) count() uint64 {
	return rs.countE
}

// Common util methods
func extractMedianFromFreqencyMap(tmap *treemap.Map, totalCount uint64) interface{} {
	medianIndex := uint64(totalCount) / 2
	for _, ele := range tmap.Keys() {
		freq, _ := tmap.Get(ele)
		if medianIndex <= freq.(uint64) {
			return ele
		}
		medianIndex -= freq.(uint64)
	}
	return nil
}

func extractModeFromFreqencyMap(tmap *treemap.Map) interface{} {
	var modeEle interface{} = nil
	var maxFreq uint64 = 0
	tmap.Each(func(key interface{}, value interface{}) {
		if maxFreq < value.(uint64) {
			maxFreq = value.(uint64)
			modeEle = key
		}
	})
	return modeEle
}

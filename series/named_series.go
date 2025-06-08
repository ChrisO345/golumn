package series

import "fmt"

// NewRangedSeries creates a new Series defined for a range of integers.
func NewRangedSeries(start, end int, t Type, name string) Series {
	numRange := make([]int, end-start)
	for i := start; i < end; i++ {
		numRange[i-start] = i
	}

	return New(numRange, t, name)
}

// NewEmptySeries creates a new Series with no values.
func NewEmptySeries(t Type, size int, name string) Series {
	switch t {
	case Int:
		return New(make([]int, size), t, name)
	case Float:
		return New(make([]float64, size), t, name)
	case Boolean:
		return New(make([]bool, size), t, name)
	case String:
		return New(make([]string, size), t, name)
	case Runic:
		panic("not implemented")
	default:
		panic(fmt.Errorf("type %v not supported", t))
	}
}

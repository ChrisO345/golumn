package series

import (
	"fmt"
	"math"
)

type intElement struct {
	e   int
	nan bool
}

// force implementation of Element interface
var _ Element = (*intElement)(nil)

func (i *intElement) Set(value any) {
	i.nan = false

	switch v := value.(type) {
	case int:
		i.e = v
	case bool:
		if v {
			i.e = 1
		} else {
			i.e = 0
		}
	case float64:
		if math.IsNaN(v) || math.IsInf(v, 0) {
			i.nan = true
			return
		}
		i.e = int(v)
	default:
		i.nan = true
		return
	}
}

func (i intElement) Get() any {
	return i.e
}

func (i intElement) IsNA() bool {
	return i.nan
}

func (i intElement) Type() Type {
	return Int
}

func (i intElement) IsNumeric() bool {
	return true
}

type floatElement struct {
	e   float64
	nan bool
}

// force implementation of Element interface
var _ Element = (*floatElement)(nil)

func (f *floatElement) Set(value any) {
	f.nan = false

	switch v := value.(type) {
	case float64:
		if math.IsNaN(v) || math.IsInf(v, 0) {
			f.nan = true
			return
		}
		f.e = v
	case int:
		f.e = float64(v)
	case bool:
		if v {
			f.e = 1.0
		} else {
			f.e = 0.0
		}
	default:
		f.nan = true
		return
	}
}

func (f floatElement) Get() any {
	return f.e
}

func (f floatElement) IsNA() bool {
	return f.nan
}

func (f floatElement) Type() Type {
	return Float
}

func (f floatElement) IsNumeric() bool {
	return true
}

type booleanElement struct {
	e   bool
	nan bool
}

// force implementation of Element interface
var _ Element = (*booleanElement)(nil)

func (b *booleanElement) Set(value any) {
	b.nan = false

	switch v := value.(type) {
	case int:
		b.e = v != 0
	case bool:
		b.e = v
	case float64:
		if math.IsNaN(v) || math.IsInf(v, 0) {
			b.nan = true
			return
		}
		b.e = v != 0.0
	default:
		b.nan = true
		return
	}
}

func (b booleanElement) Get() any {
	return b.e
}

func (b booleanElement) IsNA() bool {
	return b.nan
}

func (b booleanElement) Type() Type {
	return Boolean
}

func (b booleanElement) IsNumeric() bool {
	return true
}

type stringElement struct {
	e   string
	nan bool
}

// force implementation of Element interface
var _ Element = (*stringElement)(nil)

func (s *stringElement) Set(value any) {
	s.nan = false

	switch v := value.(type) {
	case int:
		s.e = fmt.Sprintf("%d", v)
	case bool:
		s.e = fmt.Sprintf("%t", v)
	case float64:
		s.e = fmt.Sprintf("%f", v)
	case string:
		s.e = v
	case rune:
		s.e = string(v)
	default:
		s.nan = true
		return
	}
}

func (s stringElement) Get() any {
	return s.e
}

func (s stringElement) IsNA() bool {
	return s.nan
}

func (s stringElement) Type() Type {
	return String
}

func (s stringElement) IsNumeric() bool {
	return false
}

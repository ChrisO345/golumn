package series

import (
	"fmt"
)

// Series is a collection of elements of the same type and
// is the basic building block of a DataFrame
type Series struct {
	Name     string
	elements Elements
	t        Type
}

// Elements is an interface that defines the methods that a collection of elements must implement
type Elements interface {
	Elem(int) Element
	Len() int
	Values() []any
}

// Element is an interface that defines the methods that an element must implement
type Element interface {
	Set(any)
	Get() any

	IsNA() bool
	IsNumeric() bool
	Type() Type
}

// intElement is the implementation of the Element interface for int types
type intElements []intElement

func (i intElements) Len() int           { return len(i) }
func (i intElements) Elem(j int) Element { return &i[j] }
func (i intElements) Values() []any {
	v := make([]any, len(i))
	for j, e := range i {
		v[j] = e.e
	}
	return v
}

// AsInt converts an Element to an int, returning false if the conversion is not possible
func AsInt(e Element) (int, bool) {
	if e.IsNA() {
		return 0, false
	}

	switch v := e.Get().(type) {
	case int:
		return v, true
	case float64:
		if v == float64(int(v)) {
			return int(v), true
		}
		return 0, false
	case bool:
		if v {
			return 1, true
		}
		return 0, true
	default:
		return 0, false
	}
}

// floatElement is the implementation of the Element interface for float types
type floatElements []floatElement

func (f floatElements) Len() int           { return len(f) }
func (f floatElements) Elem(j int) Element { return &f[j] }
func (f floatElements) Values() []any {
	v := make([]any, len(f))
	for j, e := range f {
		v[j] = e.e
	}
	return v
}

// AsFloat converts an Element to a float64, returning false if the conversion is not possible
func AsFloat(e Element) (float64, bool) {
	if e.IsNA() {
		return 0, false
	}

	switch v := e.Get().(type) {
	case float64:
		return v, true
	case int:
		return float64(v), true
	default:
		return 0, false
	}
}

// booleanElements is the implementation of the Element interface for float types
type booleanElements []booleanElement

func (b booleanElements) Len() int           { return len(b) }
func (b booleanElements) Elem(j int) Element { return &b[j] }
func (b booleanElements) Values() []any {
	v := make([]any, len(b))
	for j, e := range b {
		v[j] = e.e
	}
	return v
}

// AsBool converts an Element to a bool, returning false if the conversion is not possible
func AsBool(e Element) (bool, bool) {
	if e.IsNA() {
		return false, false
	}

	switch v := e.Get().(type) {
	case bool:
		return v, true
	case int:
		return v != 0, true
	case float64:
		if v == 0.0 {
			return false, true
		}
		return true, true
	default:
		return false, false
	}
}

// stringElements is the implementation of the Element interface for string types
type stringElements []stringElement

func (s stringElements) Len() int           { return len(s) }
func (s stringElements) Elem(j int) Element { return &s[j] }
func (s stringElements) Values() []any {
	v := make([]any, len(s))
	for j, e := range s {
		v[j] = e.e
	}
	return v
}

// AsString converts an Element to a string, returning false if the conversion is not possible
func AsString(e Element) (string, bool) {
	if e.IsNA() {
		return "", false
	}
	switch v := e.Get().(type) {
	case string:
		return v, true
	case int:
		return fmt.Sprint(v), true
	case float64:
		return fmt.Sprintf("%f", v), true
	case bool:
		return fmt.Sprintf("%t", v), true
	default:
		return "", false
	}
}

// Type defines the type of the series
type Type string

const (
	Int     Type = "int"
	Float   Type = "float"
	Boolean Type = "bool"
	String  Type = "string"
	Runic   Type = "rune"
)

// New creates a new series from a slice of values of type t, and a name
func New(v any, t Type, name string) Series {
	s := Series{Name: name, t: t}

	allocMemory := func(n int) {
		switch t {
		case Int:
			s.elements = make(intElements, n)
		case Float:
			s.elements = make(floatElements, n)
		case Boolean:
			s.elements = make(booleanElements, n)
		case String:
			s.elements = make(stringElements, n)
		case Runic:
			panic("not implemented")
		}
	}

	if v == nil {
		allocMemory(1)
		s.elements.Elem(0).Set(nil)
		return s
	}

	switch v_ := v.(type) {
	case []string:
		l := len(v_)
		allocMemory(l)
		for i, e := range v_ {
			s.elements.Elem(i).Set(e)
		}
	case []int:
		l := len(v_)
		allocMemory(l)
		for i, e := range v_ {
			s.elements.Elem(i).Set(e)
		}
	case []float64:
		l := len(v_)
		allocMemory(l)
		for i, e := range v_ {
			s.elements.Elem(i).Set(e)
		}
	case []bool:
		l := len(v_)
		allocMemory(l)
		for i, e := range v_ {
			s.elements.Elem(i).Set(e)
		}
	case []rune:
		panic("not implemented")
	default:
		panic(fmt.Sprintf("unsupported type, %T", v_))
	}

	return s
}

// Copy returns a memory copy of the series
func (s Series) Copy() Series {
	name := s.Name
	t := s.t

	var elements Elements
	switch s.t {
	case Int:
		elements = make(intElements, s.elements.Len())
		copy(elements.(intElements), s.elements.(intElements))
	case Float:
		elements = make(floatElements, s.elements.Len())
		copy(elements.(floatElements), s.elements.(floatElements))
	case Boolean:
		elements = make(booleanElements, s.elements.Len())
		copy(elements.(booleanElements), s.elements.(booleanElements))
	case String:
		elements = make(stringElements, s.elements.Len())
		copy(elements.(stringElements), s.elements.(stringElements))
	case Runic:
		panic("not implemented")
	}

	return Series{
		Name:     name,
		elements: elements,
		t:        t,
	}
}

// Len returns the number of elements in the series
func (s Series) Len() int {
	return s.elements.Len()
}

// Append appends a value to the series
func (s *Series) Append(v any) {
	switch s.t {
	case Int:
		s.elements = append(s.elements.(intElements), intElement{e: v.(int)})
	case Float:
		s.elements = append(s.elements.(floatElements), floatElement{e: v.(float64)})
	case Boolean:
		s.elements = append(s.elements.(booleanElements), booleanElement{e: v.(bool)})
	case String:
		s.elements = append(s.elements.(stringElements), stringElement{e: v.(string)})
	case Runic:
		panic("not implemented")
	}
}

// String returns the Stringer implementation of the series
func (s Series) String() string {
	return fmt.Sprintf("{%v %v %v}", s.Name, s.elements.Values(), s.t)
}

// Val returns the value of the element at index i
func (s Series) Val(i int) any {
	return s.elements.Elem(i).Get()
}

// Elem returns the element at index i
func (s Series) Elem(i int) Element {
	return s.elements.Elem(i)
}

// HasNa returns true if the series has any NA values
func (s Series) HasNa() bool {
	for i := range s.Len() {
		if s.Elem(i).IsNA() {
			return true
		}
	}
	return false
}

// Slice returns a copy of the series from index a to index b
func (s Series) Slice(a, b int) Series {
	if a < 0 {
		panic(fmt.Errorf("a index %v out of range", a))
	}

	if b > s.Len() || a > b {
		panic(fmt.Errorf("b index %v out of range", b))
	}

	se := Series{Name: s.Name, t: s.t}
	n := b - a

	allocMemory := func(n int) {
		switch s.t {
		case Int:
			se.elements = make(intElements, n)
		case Float:
			se.elements = make(floatElements, n)
		case Boolean:
			se.elements = make(booleanElements, n)
		case String:
			se.elements = make(stringElements, n)
		case Runic:
			panic("not implemented")
		default:
			panic("unsupported type")
		}
	}
	allocMemory(n)

	for i := a; i < b; i++ {
		se.Elem(i - a).Set(s.Val(i))
	}
	return se
}

// Head returns a slice of the first n elements of the series
func (s Series) Head(n int) Series {
	return s.Slice(0, n)
}

// Tail returns a slice of the last n elements of the series
func (s Series) Tail(n int) Series {
	return s.Slice(s.Len()-n, s.Len())
}

// Sort sorts the series in place using merge sort for better performance.
func (s Series) Sort() {
	idx := s.SortedIndex()
	// reorder according to sorted indices
	s.Order(idx...)
}

// SortedIndex returns the indices of the series sorted in ascending order using merge sort
func (s Series) SortedIndex() []int {
	n := s.Len()
	index := make([]int, n)
	for i := range n {
		index[i] = i
	}
	if n <= 1 {
		return index
	}

	// comparator for two element positions
	less := func(a, b int) bool {
		switch s.t {
		case Int:
			return s.Val(a).(int) < s.Val(b).(int)
		case Float:
			return s.Val(a).(float64) < s.Val(b).(float64)
		case Boolean:
			// false < true
			return !s.Val(a).(bool) && s.Val(b).(bool)
		case String:
			panic("not implemented")
		case Runic:
			panic("not implemented")
		default:
			panic("unsupported type")
		}
	}

	tmp := make([]int, n)
	var mergeSort func(lo, hi int)
	mergeSort = func(lo, hi int) {
		if hi-lo <= 1 {
			return
		}
		mid := (lo + hi) / 2
		mergeSort(lo, mid)
		mergeSort(mid, hi)
		i, j, k := lo, mid, lo
		for i < mid && j < hi {
			if less(index[i], index[j]) {
				tmp[k] = index[i]
				i++
			} else {
				tmp[k] = index[j]
				j++
			}
			k++
		}
		for i < mid {
			tmp[k] = index[i]
			i++
			k++
		}
		for j < hi {
			tmp[k] = index[j]
			j++
			k++
		}
		for p := lo; p < hi; p++ {
			index[p] = tmp[p]
		}
	}

	mergeSort(0, n)
	return index
}

// Order returns the series with the elements ordered according to the positions slice
func (s Series) Order(positions ...int) Series {
	if len(positions) != s.Len() {
		panic(fmt.Errorf("series and new positions must be the same length"))
	}

	// Need to copy otherwise positions collection will mutate
	newPositions := make([]int, s.Len())
	copy(newPositions, positions)

	for newPos, oldPos := range newPositions {
		if oldPos == newPos {
			continue
		}

		temp := s.Val(oldPos)
		s.Elem(oldPos).Set(s.Val(newPos))
		s.Elem(newPos).Set(temp)

		for i, pos := range newPositions {
			if pos == newPos {
				newPositions[i] = oldPos
				newPositions[newPos] = newPos
				break
			}
		}
	}

	return s
}

// Count returns the number of occurrences of the value v in the series
func (s Series) Count(v any) int {
	count := 0
	for i := range s.Len() {
		if s.Val(i) == v {
			count++
		}
	}
	return count
}

// Unique returns the true if there are no duplicates in the series
func (s Series) Unique() bool {
	seen := make(map[any]struct{})
	for i := range s.Len() {
		if _, ok := seen[s.Val(i)]; ok {
			return false
		}
		seen[s.Val(i)] = struct{}{}
	}
	return true
}

// Homogeneous returns true if there is only one value in the series
func (s Series) Homogeneous() bool {
	if s.Len() == 0 {
		panic(fmt.Errorf("cannot check homogeneity of an empty series"))
	}

	first := s.Val(0)
	for i := 1; i < s.Len(); i++ {
		if s.Val(i) != first {
			return false
		}
	}
	return true
}

// NUnique returns the number of unique values in the series
func (s Series) NUnique() int {
	seen := make(map[any]struct{})
	for i := range s.Len() {
		seen[s.Val(i)] = struct{}{}
	}
	return len(seen)
}

// ValueCounts returns a slice of the unique values in the series
func (s Series) ValueCounts() map[any]int {
	seen := make(map[any]int)
	for i := range s.Len() {
		seen[s.Val(i)] = seen[s.Val(i)] + 1
	}

	return seen
}

// Type returns the type of the series
func (s Series) Type() Type {
	return s.t
}

// InferType infers the type of the value v and returns the corresponding Type
func InferType(v any) Type {
	switch v.(type) {
	case int:
		return Int
	case float64:
		return Float
	case bool:
		return Boolean
	case string:
		return String
	case rune:
		return Runic
	default:
		panic(fmt.Errorf("unsupported type %T", v))
	}
}

// IsNumeric returns true if the series is of a numeric type (int, float, bool)
func (s Series) IsNumeric() bool {
	return s.t == Int || s.t == Float || s.t == Boolean
}

// IsObject returns true if the series is of a non-numeric type (string, rune, object)
func (s Series) IsObject() bool {
	return s.t == String || s.t == Runic
}

// Mode returns the most frequent value in the series
func (s Series) Mode() any {
	// TODO: mode only returns the first mode, need to return all modes
	counts := s.ValueCounts()
	max := 0
	var mode any
	for k, v := range counts {
		if v > max {
			max = v
			mode = k
		}
	}
	return mode
}

// Mean returns the mean of the series
func (s Series) Mean() float64 {
	if !s.IsNumeric() {
		panic(fmt.Errorf("mean is only supported for numeric types"))
	}

	var sum float64
	for i := range s.Len() {
		switch s.t {
		case Int:
			sum += float64(s.Val(i).(int))
		case Float:
			sum += s.Val(i).(float64)
		case Boolean:
			if s.Val(i).(bool) {
				sum++
			}
		}
	}

	return sum / float64(s.Len())
}

// Quantile returns the specified quantile of the series
func (s Series) Quantile(q float64) any {
	if !s.IsNumeric() {
		panic(fmt.Errorf("quantile is only supported for numeric types"))
	}

	if q < 0 || q > 1 {
		panic(fmt.Errorf("quantile must be between 0 and 1, but got %v", q))
	}

	se := s.Copy()
	se.Sort()
	index := int(float64(s.Len()) * q)
	return se.Val(index)
}

// Median returns the median of the series
func (s Series) Median() any {
	return s.Quantile(0.5)
}

func (s Series) Empty() bool {
	return s.Len() == 0
}

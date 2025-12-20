package series

import (
	"fmt"
)

// Series is a collection of elements of the same type and
// is the basic building block of a DataFrame
type Series struct {
	Name     string
	elements Elements
	// valid is a compact validity bitset for nullable values. nil means all values are valid.
	valid *Bitset
	t     Type
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
		if e.IsNA() {
			v[j] = nil
		} else {
			v[j] = e.e
		}
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
		if e.IsNA() {
			v[j] = nil
		} else {
			v[j] = e.e
		}
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
		if e.IsNA() {
			v[j] = nil
		} else {
			v[j] = e.e
		}
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
		if e.IsNA() {
			v[j] = nil
		} else {
			v[j] = e.e
		}
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
	return NewWithValidity(v, nil, t, name)
}

// NewWithValidity creates a Series and applies an optional validity mask. If mask is nil,
// validity is inferred from element-level NA indicators; mask entries set to false mark nulls.
func NewWithValidity(v any, mask []bool, t Type, name string) Series {
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
		// Create validity bitset only if element is NA
		if s.Elem(0).IsNA() {
			b := NewBitset(1)
			b.Clear(0)
			s.valid = b
		}
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

	// apply mask if provided
	if mask != nil {
		if len(mask) != s.Len() {
			panic(fmt.Errorf("validity mask length %v does not match values length %v", len(mask), s.Len()))
		}
		b := NewBitset(s.Len())
		for i := 0; i < s.Len(); i++ {
			if !mask[i] {
				// mark element as NA
				s.Elem(i).Set(nil)
				b.Clear(i)
			}
		}
		s.valid = b
		return s
	}

	// Only create a validity bitset if any element is NA; keep nil for fast-path when all valid
	anyNA := false
	for i := 0; i < s.Len(); i++ {
		if s.Elem(i).IsNA() {
			anyNA = true
			break
		}
	}
	if anyNA {
		b := NewBitset(s.Len())
		for i := 0; i < s.Len(); i++ {
			if s.Elem(i).IsNA() {
				b.Clear(i)
			}
		}
		s.valid = b
	}

	return s
}

// Copy returns a memory copy of the series
func (s Series) Copy() Series { // unchanged header, kept for edit context
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

	var valid *Bitset
	if s.valid != nil {
		valid = s.valid.Clone()
	}

	return Series{
		Name:     name,
		elements: elements,
		valid:    valid,
		t:        t,
	}
}

// NewWithValidity already defined above; add helpers: FillNA (non-mutating) and MutFillNA (mutating)

// FillNA returns a copy of the series with NA values replaced by value.
func (s Series) FillNA(value any) Series {
	res := s.Copy()
	for i := 0; i < res.Len(); i++ {
		if res.IsNull(i) {
			res.Elem(i).Set(value)
			if res.valid == nil {
				res.valid = NewBitset(res.Len())
				// set all existing to 1 except NA
				for j := 0; j < res.Len(); j++ {
					if !res.Elem(j).IsNA() {
						res.valid.Set(j)
					}
				}
			}
			res.valid.Set(i)
		}
	}
	return res
}

// MutFillNA mutates the series by replacing NA values with value.
func (s *Series) MutFillNA(value any) {
	for i := 0; i < s.Len(); i++ {
		if s.IsNull(i) {
			s.Elem(i).Set(value)
			if s.valid == nil {
				s.valid = NewBitset(s.Len())
				for j := 0; j < s.Len(); j++ {
					if !s.Elem(j).IsNA() {
						s.valid.Set(j)
					}
				}
			}
			s.valid.Set(i)
		}
	}
}

// DropNA returns a new Series with NA values removed.
func (s Series) DropNA() Series {
	n := s.Len() - s.CountNulls()
	res := Series{Name: s.Name, t: s.t}
	// allocate
	switch s.t {
	case Int:
		res.elements = make(intElements, n)
	case Float:
		res.elements = make(floatElements, n)
	case Boolean:
		res.elements = make(booleanElements, n)
	case String:
		res.elements = make(stringElements, n)
	default:
		panic("unsupported type")
	}
	idx := 0
	for i := 0; i < s.Len(); i++ {
		if s.IsNull(i) {
			continue
		}
		res.Elem(idx).Set(s.Val(i))
		idx++
	}
	return res
}

// CopyWithValidity returns a copy; if copyValues is false, values are zeroed but validity mask is preserved.
func (s Series) CopyWithValidity(copyValues bool) Series {
	res := s.Copy()
	if !copyValues {
		// zero values but clone validity
		for i := 0; i < res.Len(); i++ {
			res.Elem(i).Set(zeroForType(res.t))
		}
		if s.valid != nil {
			res.valid = s.valid.Clone()
		}
	}
	return res
}

// Len returns the number of elements in the series
func (s Series) Len() int {
	return s.elements.Len()
}

// zeroForType returns the zero placeholder value for a given series.Type
func zeroForType(t Type) any {
	switch t {
	case Int:
		return 0
	case Float:
		return 0.0
	case Boolean:
		return false
	case String:
		return ""
	default:
		return nil
	}
}

// Append appends a value to the series
func (s *Series) Append(v any) {
	// append using existing element Set semantics so NA handling is preserved
	swt := s.t
	// grow elements
	switch swt {
	case Int:
		el := intElement{}
		el.Set(v)
		s.elements = append(s.elements.(intElements), el)
	case Float:
		el := floatElement{}
		el.Set(v)
		s.elements = append(s.elements.(floatElements), el)
	case Boolean:
		el := booleanElement{}
		el.Set(v)
		s.elements = append(s.elements.(booleanElements), el)
	case String:
		el := stringElement{}
		el.Set(v)
		s.elements = append(s.elements.(stringElements), el)
	case Runic:
		panic("not implemented")
	}

	// maintain validity bitset: if s.valid is nil and no NA present, keep nil (fast path).
	// if any NA present, create or rebuild bitset to reflect current elements.
	if s.valid == nil {
		// detect if any NA exists now
		for i := 0; i < s.Len(); i++ {
			if s.Elem(i).IsNA() {
				b := NewBitset(s.Len())
				for j := 0; j < s.Len(); j++ {
					if s.Elem(j).IsNA() {
						b.Clear(j)
					}
				}
				s.valid = b
				break
			}
		}
	} else {
		// s.valid exists; ensure capacity and set/clear last bit accordingly
		s.valid.EnsureCapacity(s.Len())
		if s.Elem(s.Len() - 1).IsNA() {
			s.valid.Clear(s.Len() - 1)
		} else {
			s.valid.Set(s.Len() - 1)
		}
	}
}

// String returns the Stringer implementation of the series
func (s Series) String() string {
	return fmt.Sprintf("{%v %v %v}", s.Name, s.elements.Values(), s.t)
}

// Val returns the value of the element at index i; returns nil for NA/null values.
func (s Series) Val(i int) any {
	if s.IsNull(i) {
		return nil
	}
	return s.elements.Elem(i).Get()
}

// Elem returns the element at index i
func (s Series) Elem(i int) Element {
	return s.elements.Elem(i)
}

// HasNa returns true if the series has any NA values
func (s Series) HasNa() bool {
	for i := 0; i < s.Len(); i++ {
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

	// slice validity
	if s.valid == nil {
		// if original had no bitset but some elements might be NA flagged at element level
		anyNA := false
		for i := a; i < b; i++ {
			if s.Elem(i).IsNA() {
				anyNA = true
				break
			}
		}
		if anyNA {
			bset := NewBitset(n)
			for i := range n {
				if se.Elem(i).IsNA() {
					bset.Clear(i)
				}
			}
			se.valid = bset
		}
	} else {
		// clone relevant range
		bset := NewBitset(n)
		for i := range n {
			if s.valid.Get(a + i) {
				bset.Set(i)
			} else {
				bset.Clear(i)
			}
		}
		se.valid = bset
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

	// comparator for two element positions; NULLs are considered greater (sorted to end)
	less := func(a, b int) bool {
		// handle nulls
		if s.IsNull(a) {
			if s.IsNull(b) {
				return false
			}
			return false // a is null, b not -> a > b
		}
		if s.IsNull(b) {
			return true // a not null, b null -> a < b
		}
		switch s.t {
		case Int:
			return s.Val(a).(int) < s.Val(b).(int)
		case Float:
			return s.Val(a).(float64) < s.Val(b).(float64)
		case Boolean:
			// false < true
			return !s.Val(a).(bool) && s.Val(b).(bool)
		case String:
			return fmt.Sprint(s.Val(a)) < fmt.Sprint(s.Val(b))
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
	for i := 0; i < s.Len(); i++ {
		if s.Val(i) == v {
			count++
		}
	}
	return count
}

// Unique returns the true if there are no duplicates in the series
func (s Series) Unique() bool {
	seen := make(map[any]struct{})
	for i := 0; i < s.Len(); i++ {
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
	for i := 0; i < s.Len(); i++ {
		seen[s.Val(i)] = struct{}{}
	}
	return len(seen)
}

// ValueCounts returns a slice of the unique values in the series
func (s Series) ValueCounts() map[any]int {
	seen := make(map[any]int)
	for i := 0; i < s.Len(); i++ {
		seen[s.Val(i)] = seen[s.Val(i)] + 1
	}

	return seen
}

// Type returns the type of the series
func (s Series) Type() Type {
	return s.t
}

// IsNull returns true if the element at index i is NA/null.
func (s Series) IsNull(i int) bool {
	if s.valid != nil {
		return !s.valid.Get(i)
	}
	return s.Elem(i).IsNA()
}

// IsValid returns true if the element at index i is not NA/null.
func (s Series) IsValid(i int) bool {
	return !s.IsNull(i)
}

// CountNulls returns the number of NA/null values in the series.
func (s Series) CountNulls() int {
	if s.valid != nil {
		return s.Len() - s.valid.Count()
	}
	count := 0
	for i := 0; i < s.Len(); i++ {
		if s.Elem(i).IsNA() {
			count++
		}
	}
	return count
}

// AnyNull returns true if any element in the series is null.
func (s Series) AnyNull() bool {
	return s.CountNulls() > 0
}

// NullMask returns a slice of booleans where true indicates a null value.
func (s Series) NullMask() []bool {
	m := make([]bool, s.Len())
	for i := 0; i < s.Len(); i++ {
		m[i] = s.IsNull(i)
	}
	return m
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
	for i := 0; i < s.Len(); i++ {
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

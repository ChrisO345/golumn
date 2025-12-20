package series

import (
	"math"
	"testing"

	"github.com/chriso345/gore/assert"
)

func TestNewSeriesInt(t *testing.T) {
	s := New([]int{1, 2, 3}, Int, "Integers")
	assert.Equal(t, s.String(), "{Integers [1 2 3] int}")
}

func TestNewSeriesFloat(t *testing.T) {
	s := New([]float64{1.1, 2.2, 3.3}, Float, "Floats")
	assert.Equal(t, s.String(), "{Floats [1.1 2.2 3.3] float}")
}

func TestNewSeriesString(t *testing.T) {
	s := New([]string{"a", "b", "c"}, String, "Strings")
	assert.Equal(t, s.String(), "{Strings [a b c] string}")
}

func TestNewSeriesBool(t *testing.T) {
	s := New([]bool{true, false, true}, Boolean, "Booleans")
	assert.Equal(t, s.String(), "{Booleans [true false true] bool}")
}

func TestSeries_Copy(t *testing.T) {
	s := New([]int{1, 2, 3}, Int, "Integers")
	se := s.Copy()
	assert.Equal(t, se.String(), s.String())
}

func TestSeries_ValueCounts(t *testing.T) {
	expected := map[any]int{1: 2, 2: 2, 3: 1}
	s := New([]int{1, 2, 3, 2, 1}, Int, "Integers")
	unique := s.ValueCounts()

	for k, v := range expected {
		assert.Equal(t, unique[k], v)
	}
}

func TestSeries_Unique(t *testing.T) {
	s := New([]int{1, 4, 5, 2}, Int, "Integers")
	assert.Equal(t, s.Unique(), true)

	s = New([]int{1, 2, 3, 2, 1, 1}, Int, "Integers")
	assert.Equal(t, s.Unique(), false)
}

func TestSeries_Name(t *testing.T) {
	s := New([]int{1, 2, 3}, Int, "Integers")
	assert.Equal(t, s.Name, "Integers")
}

func TestSeries_EmptyAndLenType(t *testing.T) {
	s := New([]int{}, Int, "Empty")
	assert.Equal(t, s.Empty(), true)

	s = New([]int{1}, Int, "NotEmpty")
	assert.Equal(t, s.Empty(), false)

	s = New([]int{1, 2, 3}, Int, "Integers")
	assert.Equal(t, s.Len(), 3)
	assert.Equal(t, s.Type(), Int)
}

func TestSeries_ValuesAndElemSliceAppend(t *testing.T) {
	values := []int{1, 2, 3}
	s := New(values, Int, "Integers")
	assert.Equal(t, len(s.elements.Values()), len(values))
	for i, v := range s.elements.Values() {
		assert.Equal(t, v.(int), values[i])
	}

	e := s.Elem(1)
	assert.Equal(t, e.Get(), 2)

	slice := s.Slice(1, 3)
	assert.Equal(t, slice.String(), "{Integers [2 3] int}")

	s.Append(4)
	assert.Equal(t, s.String(), "{Integers [1 2 3 4] int}")
}

func TestSeries_AsAndTypeInterfaceCoverage(t *testing.T) {
	// As conversions
	i := intElement{e: 5}
	v, ok := AsInt(&i)
	assert.Equal(t, ok, true)
	assert.Equal(t, v, 5)

	f := floatElement{e: 3.0}
	fv, fok := AsFloat(&f)
	assert.Equal(t, fok, true)
	assert.Equal(t, fv, 3.0)

	b := booleanElement{e: true}
	bv, bok := AsBool(&b)
	assert.Equal(t, bok, true)
	assert.Equal(t, bv, true)

	sr := stringElement{e: "hi"}
	sv, sok := AsString(&sr)
	assert.Equal(t, sok, true)
	assert.Equal(t, sv, "hi")
}

func TestElementsBehaviorAndEdgeCases(t *testing.T) {
	// element set/get variations
	ie := intElement{}
	ie.Set(7)
	assert.Equal(t, ie.Get(), 7)
	ie.Set(true)
	assert.Equal(t, ie.Get(), 1)
	ie.Set(2.0)
	assert.Equal(t, ie.Get(), 2)
	ie.Set(math.Inf(1))
	assert.Equal(t, ie.IsNA(), true)

	fe := floatElement{}
	fe.Set(3.14)
	assert.Equal(t, fe.Get(), 3.14)
	fe.Set(5)
	assert.Equal(t, fe.Get(), 5.0)
	fe.Set(true)
	assert.Equal(t, fe.Get(), 1.0)
	fe.Set(math.NaN())
	assert.Equal(t, fe.IsNA(), true)

	be := booleanElement{}
	be.Set(0)
	assert.Equal(t, be.Get(), false)
	be.Set(2.0)
	assert.Equal(t, be.Get(), true)
	be.Set(false)
	assert.Equal(t, be.Get(), false)

	se := stringElement{}
	se.Set(42)
	assert.Equal(t, se.Get(), "42")
	se.Set(1.5)
	assert.Equal(t, true, len(se.Get().(string)) > 0)
	se.Set('z')
	assert.Equal(t, se.Get(), "z")
}

func TestSortSortedIndexOrderCountMisc(t *testing.T) {
	si := New([]int{3, 1, 2}, Int, "A")
	idx := si.SortedIndex()
	assert.Equal(t, idx[0], 1)

	c := New([]int{1, 2, 1, 3}, Int, "C")
	assert.Equal(t, c.Count(1), 2)
	assert.Equal(t, c.NUnique(), 3)

	h := New([]int{5, 5, 5}, Int, "H")
	assert.Equal(t, h.Homogeneous(), true)

	q := New([]int{1, 2, 3, 4}, Int, "Q")
	assert.Equal(t, q.Median(), 3)
}

func TestSortCopyAppendOthers(t *testing.T) {
	si := New([]int{5, 1, 3}, Int, "I")
	si.Sort()
	assert.Equal(t, si.Val(0), 1)

	sf := New([]float64{2.5, 0.5, 1.5}, Float, "F")
	sf.Sort()
	assert.Equal(t, sf.Val(0), 0.5)
	cf := sf.Copy()
	assert.Equal(t, cf.String(), sf.String())

	ss := New([]string{"a"}, String, "S")
	ss.Append("b")
	assert.Equal(t, ss.Len(), 2)

	sb := New([]bool{true}, Boolean, "B")
	sb.Append(false)
	assert.Equal(t, sb.Len(), 2)
}

func TestQuantileModeMeanValueCountsHasNaEmpty(t *testing.T) {
	q := New([]int{1, 2, 3, 4}, Int, "Q")
	assert.Equal(t, q.Quantile(0.0), 1)
	assert.Equal(t, q.Median(), q.Quantile(0.5))

	m := New([]int{1, 2, 2, 3}, Int, "M")
	assert.Equal(t, m.Mode(), 2)

	i := New([]int{1, 2, 3}, Int, "I")
	assert.Equal(t, i.Mean(), 2.0)

	f := New([]float64{1, 2, 3}, Float, "F")
	assert.Equal(t, f.Mean(), 2.0)

	b := New([]bool{true, false, true}, Boolean, "B")
	assert.Equal(t, b.Mean(), 2.0/3.0)

	s := New([]int{1, 2, 3}, Int, "A")
	s.Elem(1).Set(nil)
	assert.Equal(t, s.HasNa(), true)

	c := New([]string{"a", "b", "a"}, String, "C")
	vc := c.ValueCounts()
	assert.Equal(t, vc["a"], 2)
}

func TestInferTypeAndPanic(t *testing.T) {
	assert.Equal(t, InferType(1), Int)
	assert.Equal(t, InferType(1.0), Float)
	assert.Equal(t, InferType(true), Boolean)
	assert.Equal(t, InferType("x"), String)
	assert.Equal(t, InferType('x'), Runic)

	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("expected panic for unsupported type")
		}
	}()
	_ = InferType(struct{}{})
}

func TestSeriesNullHelpers(t *testing.T) {
	s := New([]int{1, 2, 3}, Int, "A")
	// set middle to NA
	s.Elem(1).Set(nil)
	// ensure bitset built
	_ = s.Copy()
	assert.Equal(t, s.IsNull(1), true)
	assert.Equal(t, s.CountNulls(), 1)
	b := s.FillNA(9)
	assert.Equal(t, b.IsNull(1), false)
	assert.Equal(t, b.Val(1), 9)

	// DropNA
	s2 := New([]int{1, 2, 3}, Int, "B")
	s2.Elem(0).Set(nil)
	d := s2.DropNA()
	assert.Equal(t, d.Len(), 2)
	assert.Equal(t, d.Val(0), 2)
}

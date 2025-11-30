package series

import (
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

func TestSeries_Empty(t *testing.T) {
	s := New([]int{}, Int, "Empty")
	assert.Equal(t, s.Empty(), true)

	s = New([]int{1}, Int, "NotEmpty")
	assert.Equal(t, s.Empty(), false)
}

func TestSeries_Len(t *testing.T) {
	s := New([]int{1, 2, 3}, Int, "Integers")
	assert.Equal(t, s.Len(), 3)
}

func TestSeries_Type(t *testing.T) {
	s := New([]int{1, 2, 3}, Int, "Integers")
	assert.Equal(t, s.Type(), Int)
}

func TestSeries_Values(t *testing.T) {
	values := []int{1, 2, 3}
	s := New(values, Int, "Integers")
	assert.Equal(t, len(s.elements.Values()), len(values))
	for i, v := range s.elements.Values() {
		assert.Equal(t, v.(int), values[i])
	}
}

func TestSeries_Elem(t *testing.T) {
	s := New([]int{1, 2, 3}, Int, "Integers")
	e := s.Elem(1)
	assert.Equal(t, e.Get(), 2)

	s = New([]bool{true, false, true}, Boolean, "Booleans")
	e = s.Elem(1)
	assert.Equal(t, e.Get(), false)
}

func TestSeries_Slice(t *testing.T) {
	s := New([]int{1, 2, 3, 4}, Int, "Integers")
	slice := s.Slice(1, 3)
	assert.Equal(t, slice.String(), "{Integers [2 3] int}")
}

func TestSeries_Append(t *testing.T) {
	s := New([]int{1, 2, 3}, Int, "Integers")
	s.Append(4)
	assert.Equal(t, s.String(), "{Integers [1 2 3 4] int}")
}

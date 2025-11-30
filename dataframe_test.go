package golumn

import (
	"testing"

	"github.com/chriso345/gore/assert"

	"github.com/chriso345/golumn/series"
)

func TestDataFrame_New_Int(t *testing.T) {
	expected := "   Integers1  Integers2\n0          1          4\n1          2          5\n2          3          6"

	df := New(
		series.New([]int{1, 2, 3}, series.Int, "Integers1"),
		series.New([]int{4, 5, 6}, series.Int, "Integers2"),
	)

	assert.Equal(t, df.String(), expected)
}

func TestDataFrame_New_Float(t *testing.T) {
	expected := "   Floats1  Floats2\n0      1.1      4.4\n1      2.2      5.5\n2      3.3      6.6"

	df := New(
		series.New([]float64{1.1, 2.2, 3.3}, series.Float, "Floats1"),
		series.New([]float64{4.4, 5.5, 6.6}, series.Float, "Floats2"),
	)

	assert.Equal(t, df.String(), expected)
}

func TestDataFrame_New_Mixed(t *testing.T) {
	expected := "   Integers  Floats\n0         1     4.4\n1         2     5.5\n2         3     6.6"

	df := New(
		series.New([]int{1, 2, 3}, series.Int, "Integers"),
		series.New([]float64{4.4, 5.5, 6.6}, series.Float, "Floats"),
	)

	assert.Equal(t, df.String(), expected)
}

func TestDataFrame_Shape(t *testing.T) {
	df := New(
		series.New([]int{1, 2, 3}, series.Int, "Integers"),
		series.New([]float64{4.4, 5.5, 6.6}, series.Float, "Floats"),
	)

	rows, cols := df.Shape()

	assert.Equal(t, rows, 3)
	assert.Equal(t, cols, 2)
}

func TestDataFrame_Slice(t *testing.T) {
	expected := "   Integers  Floats\n1         2     5.5\n2         3     6.6"

	df := New(
		series.New([]int{1, 2, 3, 4}, series.Int, "Integers"),
		series.New([]float64{4.4, 5.5, 6.6, 7.7}, series.Float, "Floats"),
	)
	result := df.Slice(1, 3).String()

	assert.Equal(t, result, expected)
}

func TestDataFrame_Head(t *testing.T) {
	expected := "   Integers  Floats  Integers2\n0         1     4.4          7\n1         2     5.5          8"

	df := New(
		series.New([]int{1, 2, 3}, series.Int, "Integers"),
		series.New([]float64{4.4, 5.5, 6.6}, series.Float, "Floats"),
		series.New([]int{7, 8, 9}, series.Int, "Integers2"),
	)
	result := df.Head(2).String()

	assert.Equal(t, result, expected)
}

func TestDataFrame_Tail(t *testing.T) {
	expected := "   Integers  Floats  Integers2\n1         2     5.5          8\n2         3     6.6          9"

	df := New(
		series.New([]int{1, 2, 3}, series.Int, "Integers"),
		series.New([]float64{4.4, 5.5, 6.6}, series.Float, "Floats"),
		series.New([]int{7, 8, 9}, series.Int, "Integers2"),
	)
	result := df.Tail(2).String()

	assert.Equal(t, result, expected)
}

func TestDataFrame_SetIndex(t *testing.T) {
	df := New(
		series.New([]int{1, 2, 3}, series.Int, "Integers"),
		series.New([]float64{4.4, 5.5, 6.6}, series.Float, "Floats"),
	)

	df = df.SetIndex(series.New([]int{7, 8, 9}, series.Int, "Integers2"))

	assert.Equal(t, df.Index().String(), "{Integers2 [7 8 9] int}")
}

func TestDataFrame_ResetIndex(t *testing.T) {
	df := New(
		series.New([]int{1, 2, 3}, series.Int, "Integers"),
		series.New([]float64{4.4, 5.5, 6.6}, series.Float, "Floats"),
	)

	df = df.SetIndex(series.New([]int{7, 8, 9}, series.Int, "Integers2"))
	df = df.ResetIndex()

	assert.Equal(t, df.Index().String(), "{Index [0 1 2] int}")
}

func TestDataFrame_Columns(t *testing.T) {
	a := series.New([]int{1, 2, 3}, series.Int, "Integers")
	b := series.New([]float64{4.4, 5.5, 6.6}, series.Float, "Floats")

	df := New(
		a,
		b,
	)

	cols := df.Columns()

	assert.Equal(t, len(cols), 2)

	assert.Equal(t, cols[0].String(), a.String())
	assert.Equal(t, cols[1].String(), b.String())
}

func TestDataFrame_Index(t *testing.T) {
	df := New(
		series.New([]int{1, 2, 3}, series.Int, "Integers"),
		series.New([]float64{4.4, 5.5, 6.6}, series.Float, "Floats"),
	)

	assert.Equal(t, df.Index().String(), "{Index [0 1 2] int}")
}

func TestDataFrame_Sort(t *testing.T) {
	expected := "   Integers  Floats\n1         1     6.6\n2         2     5.5\n0         3     4.4"

	df := New(
		series.New([]int{3, 1, 2}, series.Int, "Integers"),
		series.New([]float64{4.4, 6.6, 5.5}, series.Float, "Floats"),
	)
	df.Sort("Integers")

	assert.Equal(t, df.String(), expected)

	expected = "   Integers  Floats\n0         3     4.4\n2         2     5.5\n1         1     6.6"

	df.Sort("Floats")

	assert.Equal(t, df.String(), expected)
}

func TestDataFrame_Order(t *testing.T) {
	expected := "   Integers  Floats\n2         3     6.6\n1         2     5.5\n0         1     4.4"

	df := New(
		series.New([]int{1, 2, 3}, series.Int, "Integers"),
		series.New([]float64{4.4, 5.5, 6.6}, series.Float, "Floats"),
	)
	df = df.Order(2, 1, 0)

	assert.Equal(t, df.String(), expected)

	expected = "   Integers  Floats\n0         1     4.4\n1         2     5.5\n2         3     6.6"

	df = df.Order(2, 1, 0)

	assert.Equal(t, df.String(), expected)

	expected = "   Integers  Floats\n2         3     6.6\n0         1     4.4\n1         2     5.5"

	df = df.Order(2, 0, 1)

	assert.Equal(t, df.String(), expected)
}

func TestDataFrame_Append(t *testing.T) {
	expected := "   Integers  Floats\n0         1     4.4\n1         2     5.5\n2         3     6.6"

	df1 := New(
		series.New([]int{1, 2, 3}, series.Int, "Integers"),
	)
	s := series.New([]float64{4.4, 5.5, 6.6}, series.Float, "Floats")

	df1.Append(s)

	assert.Equal(t, df1.String(), expected)
}

func TestDataFrame_Drop(t *testing.T) {
	expected := "   Floats\n0     4.4\n1     5.5\n2     6.6"
	seriesExpected := "{Integers [1 2 3] int}"

	df := New(
		series.New([]int{1, 2, 3}, series.Int, "Integers"),
		series.New([]float64{4.4, 5.5, 6.6}, series.Float, "Floats"),
	)
	s := df.Drop("Integers")

	assert.Equal(t, df.String(), expected)
	assert.Equal(t, s.String(), seriesExpected)
}

func TestDataFrame_Filter(t *testing.T) {
	expected := "   Integers  Floats\n0         1     4.4\n1         2     5.5"

	df := New(
		series.New([]int{1, 2, 3}, series.Int, "Integers"),
		series.New([]float64{4.4, 5.5, 6.6}, series.Float, "Floats"),
	)

	filtered := df.Filter(func(row Row) bool {
		keep := row.At(0).(int)
		return keep < 3
	})

	assert.Equal(t, filtered.String(), expected)
}

func TestDataFrame_Apply(t *testing.T) {
	expected := "   Integers  Floats\n0         2     8.8\n1         4      11\n2         6    13.2"

	df := New(
		series.New([]int{1, 2, 3}, series.Int, "Integers"),
		series.New([]float64{4.4, 5.5, 6.6}, series.Float, "Floats"),
	)

	applied := df.Apply(func(row *Row) {
		row.Set("Integers", row.At(0).(int)*2)
		row.Set("Floats", row.At(1).(float64)*2)
	})

	assert.Equal(t, applied.String(), expected)
}

func TestDataFrame_Apply_Dynamic(t *testing.T) {
	expected := "   Integers  Floats  Strings\n0         2     8.8    Hello\n1         4      11    World\n2         6    13.2   Golumn"

	df := New(
		series.New([]int{1, 2, 3}, series.Int, "Integers"),
		series.New([]float64{4.4, 5.5, 6.6}, series.Float, "Floats"),
		series.New([]string{"A", "B", "C"}, series.String, "Strings"),
	)

	applied := df.Apply(func(row *Row) {
		row.Set("Integers", row.At(0).(int)*2)
		row.Set("Floats", row.At(1).(float64)*2)
		row.Set("Strings", []string{"Hello", "World", "Golumn"}[row.Position()])
	})

	assert.Equal(t, applied.String(), expected)
}

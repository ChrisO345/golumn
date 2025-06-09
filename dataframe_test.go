package golumn

import (
	"testing"

	"github.com/chriso345/golumn/series"
)

func TestDataFrame_New_Int(t *testing.T) {
	expected := "   Integers1  Integers2\n0          1          4\n1          2          5\n2          3          6"

	df := New(
		series.New([]int{1, 2, 3}, series.Int, "Integers1"),
		series.New([]int{4, 5, 6}, series.Int, "Integers2"),
	)

	if df.String() != expected {
		t.Errorf("Expected:\n%v\nGot:\n%v", expected, df.String())
	}
}

func TestDataFrame_New_Float(t *testing.T) {
	expected := "   Floats1  Floats2\n0      1.1      4.4\n1      2.2      5.5\n2      3.3      6.6"

	df := New(
		series.New([]float64{1.1, 2.2, 3.3}, series.Float, "Floats1"),
		series.New([]float64{4.4, 5.5, 6.6}, series.Float, "Floats2"),
	)

	if df.String() != expected {
		t.Errorf("Expected:\n%v\nGot:\n%v", expected, df.String())
	}
}

func TestDataFrame_New_Mixed(t *testing.T) {
	expected := "   Integers  Floats\n0         1     4.4\n1         2     5.5\n2         3     6.6"

	df := New(
		series.New([]int{1, 2, 3}, series.Int, "Integers"),
		series.New([]float64{4.4, 5.5, 6.6}, series.Float, "Floats"),
	)

	if df.String() != expected {
		t.Errorf("Expected:\n%v\nGot:\n%v", expected, df.String())
	}
}

func TestDataFrame_Shape(t *testing.T) {
	df := New(
		series.New([]int{1, 2, 3}, series.Int, "Integers"),
		series.New([]float64{4.4, 5.5, 6.6}, series.Float, "Floats"),
	)

	rows, cols := df.Shape()

	if rows != 3 || cols != 2 {
		t.Errorf("Expected (3, 2), got (%v, %v)", rows, cols)
	}
}

func TestDataFrame_Slice(t *testing.T) {
	expected := "   Integers  Floats\n1         2     5.5\n2         3     6.6"

	df := New(
		series.New([]int{1, 2, 3, 4}, series.Int, "Integers"),
		series.New([]float64{4.4, 5.5, 6.6, 7.7}, series.Float, "Floats"),
	)
	result := df.Slice(1, 3).String()

	if result != expected {
		t.Errorf("Expected:\n%v\nGot:\n%v", expected, result)
	}
}

func TestDataFrame_Head(t *testing.T) {
	expected := "   Integers  Floats  Integers2\n0         1     4.4          7\n1         2     5.5          8"

	df := New(
		series.New([]int{1, 2, 3}, series.Int, "Integers"),
		series.New([]float64{4.4, 5.5, 6.6}, series.Float, "Floats"),
		series.New([]int{7, 8, 9}, series.Int, "Integers2"),
	)
	result := df.Head(2).String()

	if result != expected {
		t.Errorf("Expected:\n%v\nGot:\n%v", expected, result)
	}
}

func TestDataFrame_Tail(t *testing.T) {
	expected := "   Integers  Floats  Integers2\n1         2     5.5          8\n2         3     6.6          9"

	df := New(
		series.New([]int{1, 2, 3}, series.Int, "Integers"),
		series.New([]float64{4.4, 5.5, 6.6}, series.Float, "Floats"),
		series.New([]int{7, 8, 9}, series.Int, "Integers2"),
	)
	result := df.Tail(2).String()

	if result != expected {
		t.Errorf("Expected:\n%v\nGot:\n%v", expected, result)
	}
}

func TestDataFrame_SetIndex(t *testing.T) {
	df := New(
		series.New([]int{1, 2, 3}, series.Int, "Integers"),
		series.New([]float64{4.4, 5.5, 6.6}, series.Float, "Floats"),
	)

	df = df.SetIndex(series.New([]int{7, 8, 9}, series.Int, "Integers2"))

	if df.Index().String() != "{Integers2 [7 8 9] int}" {
		t.Errorf("Expected index to be [7, 8, 9], got %v", df.Index().String())
	}
}

func TestDataFrame_ResetIndex(t *testing.T) {
	df := New(
		series.New([]int{1, 2, 3}, series.Int, "Integers"),
		series.New([]float64{4.4, 5.5, 6.6}, series.Float, "Floats"),
	)

	df = df.SetIndex(series.New([]int{7, 8, 9}, series.Int, "Integers2"))
	df = df.ResetIndex()

	if df.Index().String() != "{Index [0 1 2] int}" {
		t.Errorf("Expected index to be [0, 1, 2], got %v", df.Index().String())
	}
}

func TestDataFrame_Columns(t *testing.T) {
	a := series.New([]int{1, 2, 3}, series.Int, "Integers")
	b := series.New([]float64{4.4, 5.5, 6.6}, series.Float, "Floats")

	df := New(
		a,
		b,
	)

	cols := df.Columns()

	if len(cols) != 2 {
		t.Errorf("Expected 2 columns, got %v", len(cols))
	}

	if cols[0].String() != a.String() {
		t.Errorf("Expected first column to be %v, got %v", a.String(), cols[0].String())
	}
	if cols[1].String() != b.String() {
		t.Errorf("Expected second column to be %v, got %v", b.String(), cols[1].String())
	}
}

func TestDataFrame_Index(t *testing.T) {
	df := New(
		series.New([]int{1, 2, 3}, series.Int, "Integers"),
		series.New([]float64{4.4, 5.5, 6.6}, series.Float, "Floats"),
	)

	if df.Index().String() != "{Index [0 1 2] int}" {
		t.Errorf("Expected index to be nil, got %v", df.Index())
	}
}

func TestDataFrame_Sort(t *testing.T) {
	expected := "   Integers  Floats\n1         1     6.6\n2         2     5.5\n0         3     4.4"

	df := New(
		series.New([]int{3, 1, 2}, series.Int, "Integers"),
		series.New([]float64{4.4, 6.6, 5.5}, series.Float, "Floats"),
	)
	df.Sort("Integers")

	if df.String() != expected {
		t.Errorf("Expected:\n%v\nGot:\n%v", expected, df.String())
	}

	expected = "   Integers  Floats\n0         3     4.4\n2         2     5.5\n1         1     6.6"

	df.Sort("Floats")

	if df.String() != expected {
		t.Errorf("Expected:\n%v\nGot:\n%v", expected, df.String())
	}
}

func TestDataFrame_Order(t *testing.T) {
	expected := "   Integers  Floats\n2         3     6.6\n1         2     5.5\n0         1     4.4"

	df := New(
		series.New([]int{1, 2, 3}, series.Int, "Integers"),
		series.New([]float64{4.4, 5.5, 6.6}, series.Float, "Floats"),
	)
	df = df.Order(2, 1, 0)

	if df.String() != expected {
		t.Errorf("Expected:\n%v\nGot:\n%v", expected, df.String())
	}

	expected = "   Integers  Floats\n0         1     4.4\n1         2     5.5\n2         3     6.6"

	df = df.Order(2, 1, 0)

	if df.String() != expected {
		t.Errorf("Expected:\n%v\nGot:\n%v", expected, df.String())
	}

	expected = "   Integers  Floats\n2         3     6.6\n0         1     4.4\n1         2     5.5"
	df = df.Order(2, 0, 1)

	if df.String() != expected {
		t.Errorf("Expected:\n%v\nGot:\n%v", expected, df.String())
	}
}

func TestDataFrame_Append(t *testing.T) {
	expected := "   Integers  Floats\n0         1     4.4\n1         2     5.5\n2         3     6.6"

	df1 := New(
		series.New([]int{1, 2, 3}, series.Int, "Integers"),
	)
	s := series.New([]float64{4.4, 5.5, 6.6}, series.Float, "Floats")

	df1.Append(s)

	if df1.String() != expected {
		t.Errorf("Expected:\n%v\nGot:\n%v", expected, df1.String())
	}
}

func TestDataFrame_Drop(t *testing.T) {
	expected := "   Floats\n0     4.4\n1     5.5\n2     6.6"
	seriesExpected := "{Integers [1 2 3] int}"

	df := New(
		series.New([]int{1, 2, 3}, series.Int, "Integers"),
		series.New([]float64{4.4, 5.5, 6.6}, series.Float, "Floats"),
	)
	s := df.Drop("Integers")

	if df.String() != expected {
		t.Errorf("Expected:\n%v\nGot:\n%v", expected, df.String())
	}

	if s.String() != seriesExpected {
		t.Errorf("Expected:\n%v\nGot:\n%v", seriesExpected, s.String())
	}
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

	if filtered.String() != expected {
		t.Errorf("Expected:\n%v\nGot:\n%v", expected, filtered.String())
	}
}

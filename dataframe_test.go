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

func TestJoinInner(t *testing.T) {
	left := New(
		series.New([]int{1, 2}, series.Int, "id"),
		series.New([]string{"A", "B"}, series.String, "name"),
	)
	right := New(
		series.New([]int{2, 3}, series.Int, "id"),
		series.New([]string{"X", "Y"}, series.String, "val"),
	)
	j := left.Join(right, "id")
	r, c := j.Shape()
	assert.Equal(t, r, 1)
	// columns: id, name, val
	assert.Equal(t, c, 3)
	assert.Equal(t, j.At(0, 0), 2)
	assert.Equal(t, j.At(0, 1), "B")
	assert.Equal(t, j.At(0, 2), "X")
}

func colIndex(df DataFrame, name string) int {
	for i, n := range df.Names() {
		if n == name {
			return i
		}
	}
	return -1
}

func findRowByVal(df DataFrame, colName string, val any) int {
	ci := colIndex(df, colName)
	if ci == -1 {
		panic("column not found")
	}
	for i := 0; i < df.nrows; i++ {
		if df.At(i, ci) == val {
			return i
		}
	}
	return -1
}

func TestJoinLeftRightFull(t *testing.T) {
	left := New(
		series.New([]int{1, 2}, series.Int, "id"),
		series.New([]string{"A", "B"}, series.String, "name"),
	)
	right := New(
		series.New([]int{2, 3}, series.Int, "id"),
		series.New([]string{"X", "Y"}, series.String, "val"),
	)

	// Left join
	jl := left.JoinLeft(right, "id")
	rl, cl := jl.Shape()
	assert.Equal(t, rl, 2)
	assert.Equal(t, cl, 3)
	row1 := findRowByVal(jl, "id", 1)
	if row1 == -1 {
		t.Fatalf("expected to find id=1 in left join result; got names=%v", jl.Names())
	}
	nameIdx := colIndex(jl, "name")
	valIdx := colIndex(jl, "val")
	if nameIdx == -1 {
		t.Fatalf("expected column 'name' in join result, got %v", jl.Names())
	}
	assert.Equal(t, jl.At(row1, nameIdx), "A")
	// missing right val should be empty string
	if valIdx == -1 {
		// try val_y
		valIdx = colIndex(jl, "val_y")
		if valIdx == -1 {
			t.Fatalf("expected column 'val' or 'val_y' in join result, got %v", jl.Names())
		}
	}
	assert.Equal(t, jl.At(row1, valIdx), "")

	// Right join
	jr := left.JoinRight(right, "id")
	rr, cr := jr.Shape()
	assert.Equal(t, rr, 2)
	assert.Equal(t, cr, 3)
	row3 := findRowByVal(jr, "id", 3)
	if row3 == -1 {
		t.Fatalf("expected id=3 in right join result; got names=%v", jr.Names())
	}
	valIdxR := colIndex(jr, "val")
	if valIdxR == -1 {
		valIdxR = colIndex(jr, "val_y")
	}
	if valIdxR == -1 {
		t.Fatalf("expected val column in right join result, got %v", jr.Names())
	}
	assert.Equal(t, jr.At(row3, valIdxR), "Y")
	nameIdxR := colIndex(jr, "name")
	if nameIdxR == -1 {
		nameIdxR = colIndex(jr, "name_y")
	}
	if nameIdxR == -1 {
		t.Fatalf("expected name column in right join result, got %v", jr.Names())
	}
	assert.Equal(t, jr.At(row3, nameIdxR), "")

	// Full join
	jf := left.JoinFull(right, "id")
	rf, cf := jf.Shape()
	assert.Equal(t, rf, 3)
	assert.Equal(t, cf, 3)
	// ensure all ids 1,2,3 present
	assert.Equal(t, findRowByVal(jf, "id", 1) != -1, true)
	assert.Equal(t, findRowByVal(jf, "id", 2) != -1, true)
	assert.Equal(t, findRowByVal(jf, "id", 3) != -1, true)
}

func TestJoinIndexVariants(t *testing.T) {
	left := New(
		series.New([]string{"A", "B", "C"}, series.String, "name"),
		series.New([]int{10, 20, 30}, series.Int, "val1"),
	)
	// left index: 0,1,2
	right := New(
		series.New([]int{100, 200}, series.Int, "val2"),
	)
	// right index: 0,1

	// inner join on index should match indices 0 and 1
	ji := left.JoinIndex(right)
	ri, _ := ji.Shape()
	assert.Equal(t, ri, 2)
	// left join on index should preserve left rows (3)
	jl := left.JoinLeftIndex(right)
	rl, _ := jl.Shape()
	assert.Equal(t, rl, 3)
	// right join on index should preserve right rows (2)
	jr := left.JoinRightIndex(right)
	rr, _ := jr.Shape()
	assert.Equal(t, rr, 2)
	// full join index should have union (3)
	jf := left.JoinFullIndex(right)
	rf, _ := jf.Shape()
	assert.Equal(t, rf, 3)
}

func TestPivotAndUnpivot(t *testing.T) {
	df := New(
		series.New([]int{1, 1, 2}, series.Int, "id"),
		series.New([]string{"a", "b", "a"}, series.String, "key"),
		series.New([]int{10, 20, 30}, series.Int, "val"),
	)
	p := df.Pivot("id", "key", "val")
	r, c := p.Shape()
	assert.Equal(t, r, 2)
	// columns should be id, a, b
	assert.Equal(t, c, 3)
	// check some cells
	// id 1 should have a=10 b=20
	// find row with id==1
	var row1 int
	for i := range r {
		if p.At(i, 0) == 1 {
			row1 = i
			break
		}
	}
	assert.Equal(t, p.At(row1, 1), 10)
	assert.Equal(t, p.At(row1, 2), 20)

	// now unpivot back
	u := p.Unpivot([]string{"id"}, "key", "val")
	r2, c2 := u.Shape()
	// should have 3 rows (original combinations)
	assert.Equal(t, r2, 3)
	// columns id,key,val
	assert.Equal(t, c2, 3)
}

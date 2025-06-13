package golumn

import (
	"testing"

	"github.com/chriso345/golumn/series"
)

func TestRow_RowView(t *testing.T) {
	df := New(
		series.New([]int{1, 2, 3}, series.Int, "Integers"),
		series.New([]string{"A", "B", "C"}, series.String, "Strings"),
	)

	row := df.RowView(1)

	if row[0].Get() != 2 || row[1].Get() != "B" {
		t.Errorf("Expected row values [2, B], got [%v, %v]", row[0].Get(), row[1].Get())
	}
}

func TestRow_At(t *testing.T) {
	df := New(
		series.New([]int{1, 2, 3}, series.Int, "Integers"),
		series.New([]string{"A", "B", "C"}, series.String, "Strings"),
	)

	row := Row{parent: &df, index: 1}

	if row.At(0) != 2 || row.At(1) != "B" {
		t.Errorf("Expected row values [2, B], got [%v, %v]", row.At(0), row.At(1))
	}
}

func TestRow_Get(t *testing.T) {
	df := New(
		series.New([]int{1, 2, 3}, series.Int, "Integers"),
		series.New([]string{"A", "B", "C"}, series.String, "Strings"),
	)

	row := Row{parent: &df, index: 1}

	if row.Get("Integers") != 2 || row.Get("Strings") != "B" {
		t.Errorf("Expected row values [2, B], got [%v, %v]", row.Get("Integers"), row.Get("Strings"))
	}
}

func TestRow_Set(t *testing.T) {
	df := New(
		series.New([]int{1, 2, 3}, series.Int, "Integers"),
		series.New([]string{"A", "B", "C"}, series.String, "Strings"),
	)

	row := Row{parent: &df, index: 1}

	row.Set("Integers", 20)
	row.Set("Strings", "Z")

	if df.At(1, 0) != 20 || df.At(1, 1) != "Z" {
		t.Errorf("Expected updated row values [20, Z], got [%v, %v]", df.At(1, 0), df.At(1, 1))
	}
}

func TestJoinRows(t *testing.T) {
	df1 := New(
		series.New([]int{1, 2, 3}, series.Int, "Integers"),
		series.New([]string{"A", "B", "C"}, series.String, "Strings"),
	)

	row1 := Row{parent: &df1, index: 0}
	row2 := Row{parent: &df1, index: 1}
	row3 := Row{parent: &df1, index: 2}

	rows := []Row{row1, row2, row3}
	df2 := JoinRows(rows)

	if df2.nrows != 3 || df2.ncols != 2 {
		t.Errorf("Expected DataFrame with shape (3, 2), got (%v, %v)", df2.nrows, df2.ncols)
	}

	if df2.At(0, 0) != 1 || df2.At(0, 1) != "A" {
		t.Errorf("Expected first row values [1, A], got [%v, %v]", df2.At(0, 0), df2.At(0, 1))
	}
}

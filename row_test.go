package golumn

import (
	"testing"

	"github.com/chriso345/gore/assert"

	"github.com/chriso345/golumn/series"
)

func TestRow_RowView(t *testing.T) {
	df := New(
		series.New([]int{1, 2, 3}, series.Int, "Integers"),
		series.New([]string{"A", "B", "C"}, series.String, "Strings"),
	)

	row := df.RowView(1)

	assert.Equal(t, row[0].Get(), 2)
	assert.Equal(t, row[1].Get(), "B")
}

func TestRow_At(t *testing.T) {
	df := New(
		series.New([]int{1, 2, 3}, series.Int, "Integers"),
		series.New([]string{"A", "B", "C"}, series.String, "Strings"),
	)

	row := Row{parent: &df, index: 1}

	assert.Equal(t, row.At(0), 2)
	assert.Equal(t, row.At(1), "B")
}

func TestRow_Get(t *testing.T) {
	df := New(
		series.New([]int{1, 2, 3}, series.Int, "Integers"),
		series.New([]string{"A", "B", "C"}, series.String, "Strings"),
	)

	row := Row{parent: &df, index: 1}

	assert.Equal(t, row.Get("Integers"), 2)
	assert.Equal(t, row.Get("Strings"), "B")
}

func TestRow_Set(t *testing.T) {
	df := New(
		series.New([]int{1, 2, 3}, series.Int, "Integers"),
		series.New([]string{"A", "B", "C"}, series.String, "Strings"),
	)

	row := Row{parent: &df, index: 1}

	row.Set("Integers", 20)
	row.Set("Strings", "Z")

	assert.Equal(t, df.At(1, 0), 20)
	assert.Equal(t, df.At(1, 1), "Z")
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

	assert.Equal(t, df2.nrows, 3)
	assert.Equal(t, df2.ncols, 2)

	assert.Equal(t, df2.At(0, 0), 1)
	assert.Equal(t, df2.At(0, 1), "A")
}

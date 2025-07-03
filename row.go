package golumn

import (
	"fmt"

	"github.com/chriso345/golumn/series"
)

type Row struct {
	parent *DataFrame
	index  int
}

// RowView returns a Row view of the DataFrame at index i.
func (df *DataFrame) RowView(i int) []series.Element {
	row := make([]series.Element, df.ncols)
	for j, s := range df.columns {
		row[j] = s.Elem(i)
	}
	return row
}

// Position returns the index of the row in the DataFrame.
func (row Row) Position() int {
	return row.index
}

// At returns the value at the specified column.
func (row Row) At(i int) any {
	return row.parent.At(row.index, i)
}

// Get returns the value at the specified column name.
func (row Row) Get(name string) any {
	col := row.parent.Column(name)
	if col == nil {
		panic(fmt.Errorf("column %s not found", name))
	}
	return col.Val(row.index)
}

// Set sets the value at the specified column name
func (row Row) Set(name string, value any) {
	col := row.parent.Column(name)
	if col == nil {
		panic(fmt.Errorf("column %s not found", name))
	}
	if col.Type() != series.InferType(value) {
		panic(fmt.Errorf("type mismatch: expected %v, got %T", col.Type(), value))
	}
	col.Elem(row.index).Set(value)
}

// JoinRows creates a DataFrame from a slice of Row.
func JoinRows(rows []Row) DataFrame {
	if len(rows) == 0 {
		return DataFrame{}
	}

	ncols := len(rows[0].parent.columns)
	nrows := len(rows)

	cols := make([]series.Series, ncols)
	for i := range ncols {
		switch rows[0].parent.columns[i].Type() {
		case series.Int:
			cols[i] = series.New(make([]int, nrows), series.Int, rows[0].parent.columns[i].Name)
		case series.Float:
			cols[i] = series.New(make([]float64, nrows), series.Float, rows[0].parent.columns[i].Name)
		case series.Boolean:
			cols[i] = series.New(make([]bool, nrows), series.Boolean, rows[0].parent.columns[i].Name)
		case series.String:
			cols[i] = series.New(make([]string, nrows), series.String, rows[0].parent.columns[i].Name)
		default:
			panic(fmt.Errorf("unsupported series type: %v", rows[0].parent.columns[i].Type()))
		}

		for j, row := range rows {
			cols[i].Elem(j).Set(row.At(i))
		}
	}

	df := New(cols...)

	return df
}

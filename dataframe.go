package golumn

import (
	"fmt"
	"slices"
	"strings"

	"github.com/chriso345/golumn/series"
)

// DataFrame is a collection of series.Series with a shared index.
// It is similar to a table in a relational database, and is implemented
// similar to a dataframe in R or Python (pandas).
type DataFrame struct {
	index   series.Series
	columns []series.Series
	ncols   int
	nrows   int
}

// New creates a new DataFrame from a collection of series.Series.
// It has a shared index which defaults to a range of integers.
func New(se ...series.Series) DataFrame {
	if len(se) == 0 {
		panic("empty Series")
	}

	// Create index
	indices := make([]int, se[0].Len())
	for i := range se[0].Len() {
		indices[i] = i
	}

	index := series.New(indices, series.Int, "Index")

	columns := make([]series.Series, len(se))
	for i, s := range se {
		columns[i] = s.Copy()
	}
	ncols, nrows, err := checkColumnDimensions(columns...)
	if err != nil {
		panic(err)
	}

	df := DataFrame{
		index:   index,
		columns: columns,
		ncols:   ncols,
		nrows:   nrows,
	}

	// TODO: Currently assuming that column names are unique
	return df
}

// checkColumnDimensions checks that all series.Series have the same length.
func checkColumnDimensions(se ...series.Series) (ncols int, nrows int, err error) {
	ncols = len(se)
	nrows = -1
	if se == nil || ncols == 0 {
		err = fmt.Errorf("empty Series")
		return
	}

	for i, s := range se {
		if nrows == -1 {
			nrows = s.Len()
		} else if nrows != s.Len() {
			err = fmt.Errorf("series %v has length %v, expected %v", i, s.Len(), nrows)
			return
		}
	}
	return
}

// String implements the fmt.Stringer interface for DataFrame.
func (df DataFrame) String() string {
	var sb strings.Builder

	maxIndexWidth := 0
	for i := range df.nrows {
		width := len(fmt.Sprint(df.index.Val(i)))
		if width > maxIndexWidth {
			maxIndexWidth = width
		}
	}

	colWidths := make([]int, df.ncols)
	for j, col := range df.columns {
		maxWidth := len(col.Name)
		for i := range df.nrows {
			valWidth := len(fmt.Sprint(col.Val(i)))
			if valWidth > maxWidth {
				maxWidth = valWidth
			}
		}
		colWidths[j] = maxWidth
	}

	sb.WriteString(strings.Repeat(" ", maxIndexWidth))
	sb.WriteString("  ")
	for j, col := range df.columns {
		sb.WriteString(padLeft(col.Name, colWidths[j]))
		if j < df.ncols-1 {
			sb.WriteString("  ")
		}
	}
	sb.WriteString("\n")

	for i := range df.nrows {
		indexStr := fmt.Sprint(df.index.Val(i))
		sb.WriteString(padLeft(indexStr, maxIndexWidth))
		sb.WriteString("  ")

		for j, col := range df.columns {
			val := fmt.Sprint(col.Val(i))
			sb.WriteString(padLeft(val, colWidths[j]))
			if j < df.ncols-1 {
				sb.WriteString("  ")
			}
		}
		if i < df.nrows-1 {
			sb.WriteString("\n")
		}
	}

	return sb.String()
}

// padLeft adds spaces to the left of a string to fit the desired width.
func padLeft(s string, width int) string {
	return strings.Repeat(" ", width-len(s)) + s
}

// Shape returns the dimensions of the DataFrame in the form (nrows, ncols).
func (df DataFrame) Shape() (int, int) {
	return df.nrows, df.ncols
}

// Columns returns a collection of the series.Series of the DataFrame.
func (df DataFrame) Columns() []series.Series {
	return df.columns
}

// Column returns a series.Series of the DataFrame by name.
func (df DataFrame) Column(name string) *series.Series {
	for i := range df.columns {
		if df.columns[i].Name == name {
			return &df.columns[i]
		}
	}
	panic(fmt.Errorf("column %v not found", name))
}

// Names returns a collection of the names of the series.Series of the DataFrame.
func (df DataFrame) Names() []string {
	names := make([]string, df.ncols)
	for i, s := range df.columns {
		names[i] = s.Name
	}
	return names
}

// SetIndex sets the index of the DataFrame to a specified series.Series.
func (df DataFrame) SetIndex(s series.Series) DataFrame {
	if df.nrows != s.Len() {
		panic(fmt.Errorf("index length %v does not match DataFrame length %v", s.Len(), df.nrows))
	}

	df.index = s.Copy()
	return df
}

// ResetIndex resets the index of the DataFrame to a range of integers.
func (df DataFrame) ResetIndex() DataFrame {
	indices := make([]int, df.nrows)
	for i := range df.nrows {
		indices[i] = i
	}

	df.index = series.New(indices, series.Int, "Index")
	return df
}

// Index returns the index of the DataFrame.
func (df DataFrame) Index() series.Series {
	return df.index
}

// Slice returns a new DataFrame with rows from a to b
func (df DataFrame) Slice(a, b int) DataFrame {
	if a < 0 || b > df.nrows {
		panic(fmt.Errorf("b index %v out of range", b))
	}

	if a > b {
		panic(fmt.Errorf("a index %v greater than b index %v", a, b))
	}

	var s []series.Series
	for _, se := range df.columns {
		s = append(s, se.Slice(a, b))
	}

	dfNew := New(s...)
	dfNew.index = df.index.Slice(a, b)
	return dfNew
}

// Head returns a slice of the last n elements of the DataFrame. If n is not specified, it defaults to 5.
func (df DataFrame) Head(n ...int) DataFrame {
	if len(n) > 1 {
		panic("only one argument allowed")
	}
	if len(n) == 0 {
		n = []int{5}
	}

	return df.Slice(0, n[0])
}

// Tail returns a slice of the last n elements of the DataFrame. If n is not specified, it defaults to 5.
func (df DataFrame) Tail(n ...int) DataFrame {
	if len(n) > 1 {
		panic("only one argument allowed")
	}
	if len(n) == 0 {
		n = []int{5}
	}

	return df.Slice(df.nrows-n[0], df.nrows)
}

// At returns the value at the specified row and column of the DataFrame.
func (df DataFrame) At(i, j int) any {
	if i < 0 || i >= df.nrows {
		panic(fmt.Errorf("index %v out of range", i))
	}
	if j < 0 || j >= df.ncols {
		panic(fmt.Errorf("column %v out of range", j))
	}

	return df.columns[j].Val(i)
}

// Swap swaps the rows at index row1 and row2 of the DataFrame inplace.
func (df DataFrame) Swap(row1, row2 int) {
	// Swap index
	temp := df.index.Val(row1)
	df.index.Elem(row1).Set(df.index.Val(row2))
	df.index.Elem(row2).Set(temp)

	// Swap rows
	for k := range df.ncols {
		temp := df.columns[k].Val(row1)
		df.columns[k].Elem(row1).Set(df.columns[k].Val(row2))
		df.columns[k].Elem(row2).Set(temp)
	}
}

// Sort sorts the DataFrame inplace according to the specified columns.
func (df DataFrame) Sort(columns ...string) {
	if len(columns) == 0 {
		panic("no columns specified")
	}

	if len(columns) > 1 {
		panic("> 1 column not yet implemented")
	}

	column := df.Column(columns[0])

	// Sort via bubble sort according to specified column
	for i := range df.nrows {
		for j := range df.nrows - i - 1 {
			switch column.Type() {
			case series.Int:
				if column.Val(j).(int) > column.Val(j+1).(int) {
					df.Swap(j, j+1)
				}
			case series.Float:
				if column.Val(j).(float64) > column.Val(j+1).(float64) {
					df.Swap(j, j+1)
				}
			}
		}
	}
}

// Order orders the DataFrame inplace according to the specified positions.
func (df DataFrame) Order(positions ...int) DataFrame {
	if len(positions) != df.nrows {
		panic("positions must be the same length as the DataFrame")
	}

	// Need to copy otherwise positions collection will mutate
	newPositions := make([]int, df.nrows)
	copy(newPositions, positions)

	for newPos, oldPos := range newPositions {
		if oldPos == newPos {
			continue
		}

		df.Swap(oldPos, newPos)

		for i, pos := range newPositions {
			if pos == newPos {
				newPositions[i] = oldPos
				newPositions[newPos] = newPos
				break
			}
		}
	}

	return df
}

// Append appends a series.Series to right of the DataFrame.
func (df *DataFrame) Append(s series.Series) {
	if s.Len() != df.nrows {
		panic(fmt.Errorf("series length %v does not match DataFrame length %v", s.Len(), df.nrows))
	}

	df.columns = append(df.columns, s)
	df.ncols++
}

// Copy returns a deep copy of the DataFrame.
func (df DataFrame) Copy() DataFrame {
	var s []series.Series
	for _, se := range df.columns {
		s = append(s, se.Copy())
	}

	dfNew := New(s...)
	dfNew.index = df.index.Copy()
	return dfNew
}

// SelectObjectNames returns a collection of the names of object columns of the DataFrame.
func (df DataFrame) SelectObjectNames() []string {
	var objects []string
	for _, s := range df.columns {
		if s.IsObject() {
			objects = append(objects, s.Name)
		}
	}
	return objects
}

// SelectNumericNames returns a collection of the names of numeric columns of the DataFrame.
func (df DataFrame) SelectNumericNames() []string {
	var objects []string
	for _, s := range df.columns {
		if s.IsNumeric() {
			objects = append(objects, s.Name)
		}
	}
	return objects
}

// Drop removes the specified column from the DataFrame and returns it as a series.Series.
func (df *DataFrame) Drop(name string) series.Series {
	for i, s := range df.columns {
		if s.Name == name {
			df.columns = slices.Delete(df.columns, i, i+1)
			df.ncols--
			return s
		}
	}
	panic(fmt.Errorf("column %v not found", name))
}

// Filter returns a new DataFrame with rows that match the specified condition.
func (df DataFrame) Filter(condition func(row Row) bool) DataFrame {
	var s []Row
	for i := range df.nrows {
		row := Row{parent: &df, index: i}
		if condition(row) {
			s = append(s, row)
		}
	}

	return JoinRows(s)
}

// Apply applies a function to each row of the DataFrame in-place and returns the modified DataFrame.
func (df DataFrame) Apply(fn func(row *Row)) DataFrame {
	for i := range df.nrows {
		row := Row{parent: &df, index: i}
		fn(&row)
	}

	return df
}

// internal join implementation supporting modes: "inner", "left", "right", "full".
// onL/onR may be column names or "__index__" to join on the DataFrame index.
func (df DataFrame) joinInternal(other DataFrame, onL, onR, mode string) DataFrame {
	// getters for join keys
	var getLeftVal func(int) any
	var getRightVal func(int) any
	lIsIndex := onL == "__index__"
	rIsIndex := onR == "__index__"
	if lIsIndex {
		getLeftVal = func(i int) any { return df.index.Val(i) }
	} else {
		lcol := df.Column(onL)
		getLeftVal = func(i int) any { return lcol.Val(i) }
	}
	if rIsIndex {
		getRightVal = func(i int) any { return other.index.Val(i) }
	} else {
		rcol := other.Column(onR)
		getRightVal = func(i int) any { return rcol.Val(i) }
	}

	leftIdx := make(map[any][]int)
	rightIdx := make(map[any][]int)
	for i := 0; i < df.nrows; i++ {
		leftIdx[getLeftVal(i)] = append(leftIdx[getLeftVal(i)], i)
	}
	for i := 0; i < other.nrows; i++ {
		rightIdx[getRightVal(i)] = append(rightIdx[getRightVal(i)], i)
	}

	// prepare output columns: left columns then right (skip right join column when it's a column)
	outCols := make([]series.Series, 0, df.ncols+other.ncols)
	for _, s := range df.columns {
		outCols = append(outCols, series.NewEmptySeries(s.Type(), 0, s.Name))
	}
	// find right join index if needed
	rJoinIdx := -1
	if !rIsIndex {
		for idx, s := range other.columns {
			if s.Name == onR {
				rJoinIdx = idx
				break
			}
		}
	}
	// find left join index if needed
	lJoinIdx := -1
	if !lIsIndex {
		for idx, s := range df.columns {
			if s.Name == onL {
				lJoinIdx = idx
				break
			}
		}
	}
	for _, s := range other.columns {
		if !rIsIndex && s.Name == onR {
			continue
		}
		name := s.Name
		for _, ls := range df.columns {
			if ls.Name == name {
				name = name + "_y"
				break
			}
		}
		outCols = append(outCols, series.NewEmptySeries(s.Type(), 0, name))
	}

	// helper to append a combined row; li or ri may be nil to indicate missing side
	appendRow := func(li *int, ri *int) {
		// left values
		for j, s := range df.columns {
			if li == nil {
				// if left side missing but join key exists on right and left has the join column, populate it
				if lJoinIdx == j && ri != nil && !rIsIndex {
					outCols[j].Append(getRightVal(*ri))
				} else {
					outCols[j].Append(zeroForType(s.Type()))
				}
			} else {
				outCols[j].Append(s.Val(*li))
			}
		}
		// right values
		outOffset := df.ncols
		for j, s := range other.columns {
			if !rIsIndex && j == rJoinIdx {
				continue
			}
			if ri == nil {
				outCols[outOffset].Append(zeroForType(s.Type()))
			} else {
				outCols[outOffset].Append(s.Val(*ri))
			}
			outOffset++
		}
	}

	seen := make(map[any]bool)
	// process left-side keys
	for k, lrows := range leftIdx {
		rrows, ok := rightIdx[k]
		if ok {
			if mode == "right" {
				// in right mode we will handle matches only when iterating right-only later
				// but still emit matches for completeness
			}
			for _, li := range lrows {
				for _, ri := range rrows {
					tmpLi := li
					tmpRi := ri
					appendRow(&tmpLi, &tmpRi)
				}
			}
		} else {
			if mode == "left" || mode == "full" {
				for _, li := range lrows {
					tmpLi := li
					appendRow(&tmpLi, nil)
				}
			}
		}
		seen[k] = true
	}

	// right-only keys for right or full
	if mode == "right" || mode == "full" {
		for k, rrows := range rightIdx {
			if seen[k] {
				continue
			}
			for _, ri := range rrows {
				tmpRi := ri
				appendRow(nil, &tmpRi)
			}
		}
	}

	return New(outCols...)
}

// Join performs an inner join with another DataFrame on the specified column name.
// If column names collide (other than the join column), suffix "_y" is added to the right-hand columns.
func (df DataFrame) Join(other DataFrame, on string) DataFrame {
	return df.joinInternal(other, on, on, "inner")
}

// JoinLeft performs a left outer join on column 'on'.
func (df DataFrame) JoinLeft(other DataFrame, on string) DataFrame {
	return df.joinInternal(other, on, on, "left")
}

// JoinRight performs a right outer join on column 'on'.
func (df DataFrame) JoinRight(other DataFrame, on string) DataFrame {
	return df.joinInternal(other, on, on, "right")
}

// JoinFull performs a full outer join on column 'on'.
func (df DataFrame) JoinFull(other DataFrame, on string) DataFrame {
	return df.joinInternal(other, on, on, "full")
}

// JoinIndex performs an inner join on the DataFrame indices.
func (df DataFrame) JoinIndex(other DataFrame) DataFrame {
	return df.joinInternal(other, "__index__", "__index__", "inner")
}

// JoinLeftIndex performs a left outer join on the DataFrame indices.
func (df DataFrame) JoinLeftIndex(other DataFrame) DataFrame {
	return df.joinInternal(other, "__index__", "__index__", "left")
}

// JoinRightIndex performs a right outer join on the DataFrame indices.
func (df DataFrame) JoinRightIndex(other DataFrame) DataFrame {
	return df.joinInternal(other, "__index__", "__index__", "right")
}

// JoinFullIndex performs a full outer join on the DataFrame indices.
func (df DataFrame) JoinFullIndex(other DataFrame) DataFrame {
	return df.joinInternal(other, "__index__", "__index__", "full")
}

// Pivot creates a pivot table with indexCol as rows, columnsCol as columns and valuesCol as cell values.
// Aggregation is done by taking the first encountered value for the index/column pair.
func (df DataFrame) Pivot(indexCol, columnsCol, valuesCol string) DataFrame {
	idxSeries := df.Column(indexCol)
	colSeries := df.Column(columnsCol)
	valSeries := df.Column(valuesCol)

	// collect unique indices and columns in order
	indexOrder := []any{}
	indexSeen := make(map[any]bool)
	colOrder := []any{}
	colSeen := make(map[any]bool)

	for i := 0; i < df.nrows; i++ {
		iv := idxSeries.Val(i)
		cv := colSeries.Val(i)
		if !indexSeen[iv] {
			indexOrder = append(indexOrder, iv)
			indexSeen[iv] = true
		}
		if !colSeen[cv] {
			colOrder = append(colOrder, cv)
			colSeen[cv] = true
		}
	}

	// map of index->col->value
	cells := make(map[any]map[any]any)
	for i := 0; i < df.nrows; i++ {
		iv := idxSeries.Val(i)
		cv := colSeries.Val(i)
		if cells[iv] == nil {
			cells[iv] = make(map[any]any)
		}
		if _, ok := cells[iv][cv]; !ok {
			cells[iv][cv] = valSeries.Val(i)
		}
	}

	// prepare output series: index column plus one column per unique column value
	outCols := make([]series.Series, 0, 1+len(colOrder))
	// index series: reuse type
	outCols = append(outCols, series.NewEmptySeries(idxSeries.Type(), 0, indexCol))
	for _, c := range colOrder {
		name := fmt.Sprint(c)
		outCols = append(outCols, series.NewEmptySeries(valSeries.Type(), 0, name))
	}

	// populate rows
	for _, iv := range indexOrder {
		outCols[0].Append(iv)
		for j, cv := range colOrder {
			vmap := cells[iv]
			if vmap == nil {
				outCols[j+1].Append(zeroForType(valSeries.Type()))
				continue
			}
			if v, ok := vmap[cv]; ok {
				outCols[j+1].Append(v)
			} else {
				outCols[j+1].Append(zeroForType(valSeries.Type()))
			}
		}
	}

	return New(outCols...)
}

// Unpivot melts the DataFrame from wide to long format. idVars are kept as identifier columns;
// all other columns become variable/value pairs with names varName and valueName.
func (df DataFrame) Unpivot(idVars []string, varName, valueName string) DataFrame {
	idSet := make(map[string]bool)
	for _, v := range idVars {
		idSet[v] = true
	}

	// determine value columns
	valCols := []series.Series{}
	for _, s := range df.columns {
		if idSet[s.Name] {
			continue
		}
		valCols = append(valCols, s)
	}

	// prepare output: idVars..., varName (string), valueName (type depends on first valCol)
	outCols := []series.Series{}
	for _, id := range idVars {
		outCols = append(outCols, series.NewEmptySeries(df.Column(id).Type(), 0, id))
	}
	outCols = append(outCols, series.NewEmptySeries(series.String, 0, varName))
	// value type: use first valCols if exists else string
	valType := series.String
	if len(valCols) > 0 {
		valType = valCols[0].Type()
	}
	outCols = append(outCols, series.NewEmptySeries(valType, 0, valueName))

	for i := 0; i < df.nrows; i++ {
		// id values
		idVals := make([]any, len(idVars))
		for j, id := range idVars {
			idVals[j] = df.Column(id).Val(i)
		}
		for _, vc := range valCols {
			v := vc.Val(i)
			// skip placeholder zero values to avoid emitting missing combinations
			if v == zeroForType(vc.Type()) {
				continue
			}
			// append id values
			for j := range idVars {
				outCols[j].Append(idVals[j])
			}
			// var name
			outCols[len(idVars)].Append(vc.Name)
			// value
			outCols[len(idVars)+1].Append(v)
		}
	}

	return New(outCols...)
}

func zeroForType(t series.Type) any {
	switch t {
	case series.Int:
		return 0
	case series.Float:
		return 0.0
	case series.Boolean:
		return false
	case series.String:
		return ""
	default:
		return nil
	}
}

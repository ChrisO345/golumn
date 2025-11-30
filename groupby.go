package golumn

import (
	"fmt"
	"strings"

	"github.com/chriso345/golumn/series"
)

// GroupBy represents a grouping on one or more columns.
type GroupBy struct {
	parent DataFrame
	keys   []string
	groups map[string][]int // map serialized key -> row positions
	order  []string         // ordered list of group keys in first-seen order
}

// GroupBy creates a GroupBy object by specified key columns.
func (df DataFrame) GroupBy(keys ...string) GroupBy {
	if len(keys) == 0 {
		panic("no group by keys specified")
	}
	g := GroupBy{parent: df, keys: keys, groups: make(map[string][]int), order: make([]string, 0)}
	for i := 0; i < df.nrows; i++ {
		parts := make([]string, len(keys))
		for j, k := range keys {
			parts[j] = fmt.Sprint(df.Column(k).Val(i))
		}
		key := strings.Join(parts, "|~|")
		if _, ok := g.groups[key]; !ok {
			g.order = append(g.order, key)
		}
		g.groups[key] = append(g.groups[key], i)
	}
	return g
}

// String implements fmt.Stringer for GroupBy in a pandas-like way: prints rows with group key
// values shown once per consecutive group and hidden (empty) for duplicate consecutive rows.
func (g GroupBy) String() string {
	if g.parent.nrows == 0 || len(g.parent.columns) == 0 {
		return ""
	}

	// Build ordered columns: keys first (in given order), then remaining columns
	keySet := make(map[string]struct{})
	for _, k := range g.keys {
		keySet[k] = struct{}{}
	}
	var cols []series.Series
	// append key columns in order
	for _, k := range g.keys {
		for _, c := range g.parent.columns {
			if c.Name == k {
				cols = append(cols, c)
				break
			}
		}
	}
	// append other columns
	for _, c := range g.parent.columns {
		if _, ok := keySet[c.Name]; ok {
			continue
		}
		cols = append(cols, c)
	}

	// compute widths
	maxIndexWidth := 0
	for i := 0; i < g.parent.nrows; i++ {
		w := len(fmt.Sprint(g.parent.index.Val(i)))
		if w > maxIndexWidth {
			maxIndexWidth = w
		}
	}

	colWidths := make([]int, len(cols))
	for j, col := range cols {
		maxW := len(col.Name)
		for i := 0; i < g.parent.nrows; i++ {
			v := fmt.Sprint(col.Val(i))
			if lw := len(v); lw > maxW {
				maxW = lw
			}
		}
		colWidths[j] = maxW
	}

	var sb strings.Builder
	// header
	sb.WriteString(strings.Repeat(" ", maxIndexWidth))
	sb.WriteString("  ")
	for j, col := range cols {
		sb.WriteString(padLeft(col.Name, colWidths[j]))
		if j < len(cols)-1 {
			sb.WriteString("  ")
		}
	}
	sb.WriteString("\n")

	// print rows grouped by group order; hide duplicate group key values within each group
	firstRow := true
	for _, key := range g.order {
		positions := g.groups[key]
		for gi, pos := range positions {
			if !firstRow {
				sb.WriteString("\n")
			}
			firstRow = false
			indexStr := fmt.Sprint(g.parent.index.Val(pos))
			sb.WriteString(padLeft(indexStr, maxIndexWidth))
			sb.WriteString("  ")

			for j := range cols {
				val := fmt.Sprint(cols[j].Val(pos))
				if j < len(g.keys) && gi > 0 {
					// hide duplicate key values for subsequent rows in the same group
					sb.WriteString(padLeft("", colWidths[j]))
				} else {
					sb.WriteString(padLeft(val, colWidths[j]))
				}
				if j < len(cols)-1 {
					sb.WriteString("  ")
				}
			}
		}
	}

	return sb.String()
}

// Groups returns a map of serialized group key -> DataFrame for that group.
func (g GroupBy) Groups() map[string]DataFrame {
	out := make(map[string]DataFrame)
	for _, key := range g.order {
		positions := g.groups[key]
		subs := make([]series.Series, len(g.parent.columns))
		for ci, col := range g.parent.columns {
			switch col.Type() {
			case series.Int:
				subs[ci] = series.New(make([]int, len(positions)), series.Int, col.Name)
			case series.Float:
				subs[ci] = series.New(make([]float64, len(positions)), series.Float, col.Name)
			case series.Boolean:
				subs[ci] = series.New(make([]bool, len(positions)), series.Boolean, col.Name)
			case series.String:
				subs[ci] = series.New(make([]string, len(positions)), series.String, col.Name)
			default:
				panic(fmt.Errorf("unsupported series type: %v", col.Type()))
			}
		}
		for ri, pos := range positions {
			for ci := range subs {
				subs[ci].Elem(ri).Set(g.parent.columns[ci].Val(pos))
			}
		}
		out[key] = New(subs...)
	}
	return out
}

// Aggregate applies aggregation functions to each group. aggs is a map from output name to func(rows DataFrame) series.Series
func (g GroupBy) Aggregate(agg func(df DataFrame) DataFrame) DataFrame {
	first := true
	var out DataFrame
	for _, key := range g.order {
		positions := g.groups[key]
		// build a sub-DataFrame for the group
		subs := make([]series.Series, len(g.parent.columns))
		for ci, col := range g.parent.columns {
			switch col.Type() {
			case series.Int:
				subs[ci] = series.New(make([]int, len(positions)), series.Int, col.Name)
			case series.Float:
				subs[ci] = series.New(make([]float64, len(positions)), series.Float, col.Name)
			case series.Boolean:
				subs[ci] = series.New(make([]bool, len(positions)), series.Boolean, col.Name)
			case series.String:
				subs[ci] = series.New(make([]string, len(positions)), series.String, col.Name)
			default:
				panic(fmt.Errorf("unsupported series type: %v", col.Type()))
			}
		}
		for ri, pos := range positions {
			for ci := range subs {
				subs[ci].Elem(ri).Set(g.parent.columns[ci].Val(pos))
			}
		}
		subdf := New(subs...)
		// apply agg which should return a DataFrame with a single row
		res := agg(subdf)
		if res.nrows == 0 {
			continue
		}
		if first {
			// initialize out with res copy
			out = res.Copy()
			first = false
			continue
		}
		// append row-wise: res is expected to be single-row DataFrame
		if res.nrows != 1 {
			panic("aggregate function must return a single-row DataFrame per group")
		}
		// append each value to the corresponding column of out
		for ci := range res.columns {
			out.columns[ci].Append(res.columns[ci].Val(0))
		}
		out.nrows = out.columns[0].Len()
	}
	if first {
		return DataFrame{}
	}
	return out
}

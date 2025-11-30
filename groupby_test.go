package golumn

import (
	"testing"

	"github.com/chriso345/gore/assert"

	"github.com/chriso345/golumn/series"
)

func TestGroupBy_StringAndGroups(t *testing.T) {
	df := New(
		series.New([]string{"a", "b", "a"}, series.String, "Key"),
		series.New([]int{1, 2, 3}, series.Int, "Val"),
	)

	gb := df.GroupBy("Key")

	// Stringer: grouped/pseudo-sorted display where groups are printed together
	expected := "   Key  Val\n0    a    1\n2         3\n1    b    2"
	assert.Equal(t, gb.String(), expected)

	groups := gb.Groups()
	// two groups: "a" and "b"
	assert.Equal(t, len(groups), 2)

	gA := groups["a"]
	rows, cols := gA.Shape()
	assert.Equal(t, rows, 2)
	assert.Equal(t, cols, 2)
	// values preserved in order
	valCol := gA.Column("Val")
	assert.Equal(t, valCol.Val(0), 1)
	assert.Equal(t, valCol.Val(1), 3)
}

func TestGroupBy_Aggregate(t *testing.T) {
	// group and aggregate: compute sum per group
	df := New(
		series.New([]string{"x", "y", "x", "y", "x"}, series.String, "Key"),
		series.New([]int{1, 2, 3, 4, 5}, series.Int, "Val"),
	)
	gb := df.GroupBy("Key")
	agg := func(d DataFrame) DataFrame {
		// return single-row DataFrame with sum of Val
		s := d.Column("Val")
		sv := 0
		for i := 0; i < s.Len(); i++ {
			sv += s.Val(i).(int)
		}
		return New(series.New([]string{"sum"}, series.String, "Metric"), series.New([]int{sv}, series.Int, "Val"))
	}

	res := gb.Aggregate(func(d DataFrame) DataFrame {
		// produce a df with single row: Metric, Val
		return agg(d)
	})

	// Expect two rows in result for x and y groups
	rows, cols := res.Shape()
	assert.Equal(t, rows, 2)
	assert.Equal(t, cols, 2)

	// Because groups preserve order of first appearance, x then y
	// check Val sums
	valCol := res.Column("Val")
	assert.Equal(t, valCol.Val(0), 9) // 1+3+5
	assert.Equal(t, valCol.Val(1), 6) // 2+4
}

func TestGroupBy_AggregateEmptyAndPanic(t *testing.T) {
	// aggregate function returning >1 row should panic
	df := New(
		series.New([]string{"a", "b"}, series.String, "Key"),
		series.New([]int{1, 2}, series.Int, "Val"),
	)
	gb := df.GroupBy("Key")
	badAgg := func(d DataFrame) DataFrame {
		// return two-row DataFrame to trigger panic
		return New(series.New([]string{"r1", "r2"}, series.String, "Metric"), series.New([]int{1, 2}, series.Int, "Val"))
	}

	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("expected panic from Aggregate when agg returns multiple rows")
		}
	}()
	_ = gb.Aggregate(badAgg)
}

func TestGroupBy_MultiKeyGroups(t *testing.T) {
	// group by two keys and verify serialized group keys and group DataFrames
	df := New(
		series.New([]string{"a", "a", "b", "b"}, series.String, "K1"),
		series.New([]int{1, 1, 2, 2}, series.Int, "K2"),
		series.New([]int{10, 20, 30, 40}, series.Int, "Val"),
	)
	gb := df.GroupBy("K1", "K2")
	groups := gb.Groups()
	// expect two groups: a|~|1 and b|~|2
	if len(groups) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(groups))
	}
	// validate groups
	if g, ok := groups["a|~|1"]; ok {
		if rows, _ := g.Shape(); rows != 2 {
			t.Fatalf("expected 2 rows in group a|~|1, got %d", rows)
		}
		val := g.Column("Val")
		if val.Val(0) != 10 || val.Val(1) != 20 {
			t.Fatalf("unexpected vals for group a|~|1: %v,%v", val.Val(0), val.Val(1))
		}
	} else {
		t.Fatalf("group a|~|1 missing")
	}
	if g, ok := groups["b|~|2"]; ok {
		if rows, _ := g.Shape(); rows != 2 {
			t.Fatalf("expected 2 rows in group b|~|2, got %d", rows)
		}
		val := g.Column("Val")
		if val.Val(0) != 30 || val.Val(1) != 40 {
			t.Fatalf("unexpected vals for group b|~|2: %v,%v", val.Val(0), val.Val(1))
		}
	} else {
		t.Fatalf("group b|~|2 missing")
	}
}

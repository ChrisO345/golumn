package series_test

import (
	"testing"

	"github.com/chriso345/golumn/internal/testutils/assert"
	"github.com/chriso345/golumn/series"
)

func TestNewRangedSeries(t *testing.T) {
	expected := "{Integers [1 2 3 4] int}"
	s := series.NewRangedSeries(1, 5, series.Int, "Integers")

	assert.AssertEqual(t, s.Name, "Integers")
	assert.AssertEqual(t, s.Type(), series.Int)
	assert.AssertEqual(t, s.Len(), 4)
	assert.AssertEqual(t, s.String(), expected)
}

func TestNewEmptySeries_Int(t *testing.T) {
	expected := "{Integers [0 0 0 0] int}"
	s := series.NewEmptySeries(series.Int, 4, "Integers")

	assert.AssertEqual(t, s.Name, "Integers")
	assert.AssertEqual(t, s.Type(), series.Int)
	assert.AssertEqual(t, s.Len(), 4)
	assert.AssertEqual(t, s.String(), expected)

	for i := range s.Len() {
		assert.AssertEqual(t, s.Val(i), 0)
	}
}

func TestNewEmptySeries_String(t *testing.T) {
	expected := "{Strings [ ] string}"
	s := series.NewEmptySeries(series.String, 2, "Strings")

	assert.AssertEqual(t, s.Name, "Strings")
	assert.AssertEqual(t, s.Type(), series.String)
	assert.AssertEqual(t, s.Len(), 2)
	assert.AssertEqual(t, s.String(), expected)

	for i := range s.Len() {
		assert.AssertEqual(t, s.Val(i), "")
	}
}

func TestNewEmptySeries_Float(t *testing.T) {
	expected := "{Floats [0 0] float}"
	s := series.NewEmptySeries(series.Float, 2, "Floats")

	assert.AssertEqual(t, s.Name, "Floats")
	assert.AssertEqual(t, s.Type(), series.Float)
	assert.AssertEqual(t, s.Len(), 2)
	assert.AssertEqual(t, s.String(), expected)

	for i := range s.Len() {
		assert.AssertEqual(t, s.Val(i), 0.0)
	}
}

func TestNewEmptySeries_Boolean(t *testing.T) {
	expected := "{Bools [false false] bool}"
	s := series.NewEmptySeries(series.Boolean, 2, "Bools")

	assert.AssertEqual(t, s.Name, "Bools")
	assert.AssertEqual(t, s.Type(), series.Boolean)
	assert.AssertEqual(t, s.Len(), 2)
	assert.AssertEqual(t, s.String(), expected)

	for i := range s.Len() {
		assert.AssertEqual(t, s.Val(i), false)
	}
}

package series

import "testing"

func TestNewRangedSeries(t *testing.T) {
	expected := "{Integers [1 2 3 4] int}"
	s := NewRangedSeries(1, 5, Int, "Integers")

	if s.String() != expected {
		t.Errorf("Expected:\n%v\nGot:\n%v", expected, s.String())
	}
	// Output: {Integers [1 2 3 4] int}
}

func TestNewEmptySeries(t *testing.T) {
	expected := "{Integers [0 0 0 0] int}"
	s := NewEmptySeries(Int, 4, "Integers")

	if s.String() != expected {
		t.Errorf("Expected:\n%v\nGot:\n%v", expected, s.String())
	}
}


package series

import "testing"

func TestBitsetBasic(t *testing.T) {
	b := NewBitset(130)
	if b.n != 130 {
		t.Fatalf("expected n=130, got %d", b.n)
	}
	if !b.Get(0) || !b.Get(64) || !b.Get(129) {
		t.Fatalf("expected bits set by default")
	}
	b.Clear(64)
	if b.Get(64) {
		t.Fatalf("expected bit 64 cleared")
	}
	if b.Count() != 129 {
		t.Fatalf("expected count 129, got %d", b.Count())
	}
	b2 := b.Clone()
	b2.Set(64)
	if b.Get(64) || !b2.Get(64) {
		t.Fatalf("clone semantics broken")
	}
	idx := b2.ToIndices()
	if len(idx) != 130 {
		t.Fatalf("expected 130 indices, got %d", len(idx))
	}
}

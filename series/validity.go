package series

import (
	"math/bits"
)

// Bitset stores validity bits for n elements. A set bit means the value is valid (not null).
type Bitset struct {
	words []uint64
	n     int
}

// NewBitset allocates a Bitset for n elements, defaulting all bits to 1 (valid).
func NewBitset(n int) *Bitset {
	if n <= 0 {
		return &Bitset{words: nil, n: 0}
	}
	w := (n + 63) / 64
	b := &Bitset{words: make([]uint64, w), n: n}
	for i := range b.words {
		b.words[i] = ^uint64(0)
	}
	// clear bits that are out of range in the last word
	if n%64 != 0 {
		b.words[len(b.words)-1] &= (uint64(1)<<(uint(n%64)) - 1)
	}
	return b
}

// EnsureCapacity ensures the bitset can hold at least n bits.
func (b *Bitset) EnsureCapacity(n int) {
	if b == nil {
		return
	}
	if n <= b.n {
		return
	}
	w := (n + 63) / 64
	if len(b.words) < w {
		newWords := make([]uint64, w)
		copy(newWords, b.words)
		for i := len(b.words); i < w; i++ {
			newWords[i] = ^uint64(0)
		}
		b.words = newWords
	}
	b.n = n
	if n%64 != 0 {
		b.words[len(b.words)-1] &= (uint64(1)<<(uint(n%64)) - 1)
	}
}

// Clone returns a deep copy of the bitset.
func (b *Bitset) Clone() *Bitset {
	if b == nil {
		return nil
	}
	w := make([]uint64, len(b.words))
	copy(w, b.words)
	return &Bitset{words: w, n: b.n}
}

// Set marks index i as valid. Panics on out of range.
func (b *Bitset) Set(i int) {
	if b == nil {
		return
	}
	if i < 0 || i >= b.n {
		panic("index out of range")
	}
	word := i / 64
	bit := uint(i % 64)
	b.words[word] |= 1 << bit
}

// Clear marks index i as null. Panics on out of range.
func (b *Bitset) Clear(i int) {
	if b == nil {
		return
	}
	if i < 0 || i >= b.n {
		panic("index out of range")
	}
	word := i / 64
	bit := uint(i % 64)
	b.words[word] &^= 1 << bit
}

// Get returns true if index i is valid (not null).
func (b *Bitset) Get(i int) bool {
	if b == nil {
		return true
	}
	if i < 0 || i >= b.n {
		return false
	}
	word := i / 64
	bit := uint(i % 64)
	return (b.words[word]&(1<<bit) != 0)
}

// Count returns the number of set bits (valid values).
func (b *Bitset) Count() int {
	if b == nil {
		return 0
	}
	c := 0
	for _, w := range b.words {
		c += bits.OnesCount64(w)
	}
	return c
}

// ToIndices returns the indices of bits that are valid.
func (b *Bitset) ToIndices() []int {
	if b == nil {
		return nil
	}
	res := make([]int, 0, b.Count())
	for wi, w := range b.words {
		for w != 0 {
			tz := bits.TrailingZeros64(w)
			idx := wi*64 + tz
			if idx >= b.n {
				break
			}
			res = append(res, idx)
			w &^= 1 << uint(tz)
		}
	}
	return res
}

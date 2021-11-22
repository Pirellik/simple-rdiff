package rollsum

import (
	"testing"
)

func TestInit(t *testing.T) {
	giveData := []byte{1, 2, 3}
	wantSum := uint32(655366)

	rSum := New()
	rSum.Init(giveData)
	gotSum := rSum.Sum()

	if gotSum != wantSum {
		t.Errorf("sums do not match, got = %d; want = %d", gotSum, wantSum)
	}
}

// Init() + Roll() should give the same result as only Init() with all the data
func TestRoll(t *testing.T) {
	data := []byte{'h', 'e', 'l', 'l', 'o', ' ', 'w', 'o', 'r', 'l', 'd'}

	r1 := New()
	r1.Init(data[6:11])
	r2 := New()
	r2.Init(data[:5])
	for i, b := range data[5:] {
		r2.Roll(data[i], b)
	}
	r1Sum := r1.Sum()
	r2Sum := r2.Sum()

	if r1Sum != r2Sum {
		t.Errorf("sums do not match, r1Sum = %d; r2Sum = %d", r1Sum, r2Sum)
	}
}

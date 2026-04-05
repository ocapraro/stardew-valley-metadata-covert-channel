package inventory

import (
	"fmt"
	"maps"
	"math/big"
	"slices"
	"testing"
)

func TestPartitionCounts(t *testing.T) {
	c := cache{partitionCounts: map[[2]uint8]uint64{}}
	if got := partitionCounts(5, 2, c); got != 3 {
		t.Fatalf("partitionCounts(5,2) expected 3, got %d", got)
	}
	if got := partitionCounts(4, 3, c); got != 4 {
		t.Fatalf("partitionCounts(4,3) expected 4, got %d", got)
	}
}

func TestSumFactorization(t *testing.T) {
	c := cache{sumFactorizations: map[[2]uint16][]map[uint16]uint16{}}
	got := sumFactorization(5, 2, c)
	if len(got) != 2 {
		t.Fatalf("expected 2 factorizations, got %d", len(got))
	}

	enc := make([]string, 0, len(got))
	for _, f := range got {
		keys := slices.Collect(maps.Keys(f))
		slices.Sort(keys)
		s := ""
		for _, k := range keys {
			s += fmt.Sprintf("%d:%d;", k, f[k])
		}
		enc = append(enc, s)
	}
	slices.Sort(enc)
	if enc[0] != "1:1;4:1;" || enc[1] != "2:1;3:1;" {
		t.Fatalf("unexpected sumFactorization output: %v", enc)
	}
}

func TestGetDuplicateWeight(t *testing.T) {
	c := cache{smallFactorial: map[int]int{}, duplicateWeights: map[[3]uint16]*big.Rat{}}
	got := getDuplicateWeight(4, 2, 4, c)
	want := big.NewRat(3, 2) // 1+3 and 2+2 (weighted by 1/2!)
	if got.Cmp(want) != 0 {
		t.Fatalf("expected %s, got %s", want.RatString(), got.RatString())
	}
}

func TestCombinatoricsHelpers(t *testing.T) {
	if got := nCr(5, 2); got != 10 {
		t.Fatalf("nCr(5,2) expected 10, got %d", got)
	}
	if got := factorialBig(6).Int64(); got != 720 {
		t.Fatalf("factorialBig(6) expected 720, got %d", got)
	}

	c := cache{smallFactorial: map[int]int{}}
	if got := factorial(6, c); got != 720 {
		t.Fatalf("factorial(6) expected 720, got %d", got)
	}
}

func TestSpreadHelpers(t *testing.T) {
	spread := calculateSpread(2, 2)
	if len(spread) != 3 {
		t.Fatalf("calculateSpread(2,2) expected 3 rows, got %d", len(spread))
	}

	spares := spreadSpares(2, 2)
	if len(spares) != 2 {
		t.Fatalf("spreadSpares(2,2) expected 2 rows, got %d", len(spares))
	}
}

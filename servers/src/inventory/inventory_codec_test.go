package inventory

import (
	"math/big"
	"testing"
)

func TestTextNumberRoundTrip(t *testing.T) {
	cases := []string{"", "Lily", "Stardew123", "A"}
	for _, input := range cases {
		n := TextToNumber(input)
		got := NumberToText(n)
		if got != input {
			t.Fatalf("roundtrip mismatch for %q: got %q", input, got)
		}
	}
}

func TestNumberToTextPadsToFullBytes(t *testing.T) {
	// Binary 1000001 should decode to "A" after left-padding to 8 bits.
	n := big.NewRat(65, 1)
	got := NumberToText(n)
	if got != "A" {
		t.Fatalf("expected A, got %q", got)
	}
}

func TestTotalInventoryBits(t *testing.T) {
	inv := Inventory{Items: []Item{{Name: "A", Stack: 3}, {Name: "B", Stack: 2}, {Name: "Blank", Stack: 0}}}
	inv.NewCache()
	bits := totalInventoryBits(inv)
	if bits != 4 {
		t.Fatalf("expected 4 bits for 15 combinations, got %d", bits)
	}
}

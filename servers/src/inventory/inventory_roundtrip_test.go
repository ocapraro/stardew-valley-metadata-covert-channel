package inventory

import (
	"fmt"
	"math/big"
	"testing"
)

func totalCombinations(inv Inventory) *big.Int {
	counts := inv.CalculateSpreadPermutationCounts()
	total := big.NewRat(0, 1)
	for _, c := range counts {
		total.Add(total, c)
	}
	if !total.IsInt() {
		return big.NewInt(0)
	}
	return new(big.Int).Quo(total.Num(), total.Denom())
}

func invSig(inv Inventory) string {
	return fmt.Sprint(inv.Items)
}

func TestInventoryEncodeDecodeHelpers(t *testing.T) {
	inv := Inventory{Items: []Item{{Name: "A", Stack: 1}, {Name: "B", Stack: 2}, {Name: "Blank", Stack: 0}}}
	encoded := inv.encode()
	decoded := decodeInventory(encoded)
	if len(decoded.Items) != len(inv.Items) {
		t.Fatalf("decoded length mismatch: got %d want %d", len(decoded.Items), len(inv.Items))
	}
	for i := range inv.Items {
		if decoded.Items[i] != inv.Items[i] {
			t.Fatalf("decoded item mismatch at %d: got %+v want %+v", i, decoded.Items[i], inv.Items[i])
		}
	}
}

func TestSmallInventoryFullBijection(t *testing.T) {
	inv := Inventory{Items: []Item{{Name: "A", Stack: 3}, {Name: "B", Stack: 2}, {Name: "Blank", Stack: 0}}}
	inv.NewCache()

	limit := totalCombinations(inv).Int64()
	if limit != 15 {
		t.Fatalf("unexpected total combinations: got %d want 15", limit)
	}

	seen := map[string]int64{}
	for i := int64(1); i <= limit; i++ {
		variation := inv.GetVariation(big.NewRat(i, 1))
		variation.Cache = inv.Cache
		recovered := variation.GetIndex()
		if recovered.Cmp(big.NewRat(i, 1)) != 0 {
			t.Fatalf("roundtrip mismatch at %d: recovered %s", i, recovered.RatString())
		}

		sig := invSig(variation)
		if prev, ok := seen[sig]; ok {
			t.Fatalf("collision: index %d and %d produce same variation %s", prev, i, sig)
		}
		seen[sig] = i
	}
}

func TestKnownInventoryRoundTrip(t *testing.T) {
	inv := Inventory{Items: []Item{{"Pickaxe", 1}, {"Scythe", 1}, {"Axe", 1}, {"Coal", 3}, {"Wood", 123}, {"Stone", 11}, {"Parsnip Seeds", 13}, {"Sap", 54}, {"Fiber", 14}, {"Blank", 0}, {"Blank", 0}, {"Blank", 0}}}
	inv.NewCache()

	target := big.NewRat(1281977465, 1)
	variation := inv.GetVariation(target)
	variation.Cache = inv.Cache
	recovered := variation.GetIndex()

	if recovered.Cmp(target) != 0 {
		t.Fatalf("round-trip mismatch: target=%s recovered=%s variation=%v", target.RatString(), recovered.RatString(), variation.Items)
	}
}

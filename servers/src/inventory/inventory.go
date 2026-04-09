package inventory

import (
	"fmt"
	"maps"
	"math/big"
	"slices"
	"sort"
	"strings"
)

type cache struct {
	smallFactorial     map[int]int
	duplicateWeights   map[[3]uint16]*big.Rat
	partitionCounts    map[[2]uint8]uint64
	sumFactorizations  map[[2]uint16][]map[uint16]uint16
	perms              map[[2]uint8][][]uint8
	weightedPartCounts map[[3]uint16]*big.Rat
}

type Item struct {
	Name  string
	Stack uint16
}

func (i Item) encode() string {
	return fmt.Sprintf("[%s:%d]", i.Name, i.Stack)
}

type Inventory struct {
	Items []Item
	Cache cache
}

func (i *Inventory) NewCache() {
	i.Cache = cache{
		smallFactorial:     make(map[int]int),
		duplicateWeights:   make(map[[3]uint16]*big.Rat),
		partitionCounts:    make(map[[2]uint8]uint64),
		sumFactorizations:  make(map[[2]uint16][]map[uint16]uint16),
		perms:              make(map[[2]uint8][][]uint8),
		weightedPartCounts: make(map[[3]uint16]*big.Rat),
	}
}

// Contains checks to see if an Item is in an Inventory
func (i Inventory) Contains(target Item) bool {
	return slices.Contains(i.Items, target)
}

// collapse combines all stacks of the same item
func (i *Inventory) collapse() {
	length := len(i.Items)
	counts := make(map[string]uint16)
	for _, item := range i.Items {
		if item.Stack < 1 {
			continue
		}
		counts[item.Name] += item.Stack
	}

	i.Items = []Item{}

	for name, stack := range counts {
		item := Item{
			Name:  name,
			Stack: stack,
		}
		i.Items = append(i.Items, item)
	}

	// Fill the rest of the inventory with blanks
	for len(i.Items) < length {
		i.Items = append(i.Items, Item{
			Name:  "Blank",
			Stack: 0,
		})
	}
}

func (i *Inventory) removeBlanks() {
	var deblanked []Item
	for _, item := range i.Items {
		if item.Stack < 1 {
			continue
		}
		deblanked = append(deblanked, item)
	}
	i.Items = deblanked
}

func (i *Inventory) addBlanks(length int) {
	for len(i.Items) < length {
		i.Items = append(i.Items, Item{
			Name:  "Blank",
			Stack: 0,
		})
	}
}

func (i Inventory) encode() string {
	var b strings.Builder
	b.WriteString("[")
	for index, item := range i.Items {
		if index > 0 {
			b.WriteString(",")
		}
		b.WriteString(item.Name)
		b.WriteString(":")
		fmt.Fprint(&b, item.Stack)
	}
	b.WriteString("]")
	return b.String()
}

func decodeInventory(encodedInventory string) Inventory {
	i := Inventory{}
	encodedInventory = strings.Trim(encodedInventory, "[]")
	itemPairs := strings.Split(encodedInventory, ",")

	for _, pair := range itemPairs {
		parts := strings.Split(pair, ":")
		if len(parts) != 2 {
			continue
		}

		name := parts[0]
		var stack uint16
		fmt.Sscanf(parts[1], "%d", &stack)

		i.Items = append(i.Items, Item{
			Name:  name,
			Stack: stack,
		})
	}

	return i
}

// partitionCounts calculates the number of ways x can be made via a sum of y positive numbers
func partitionCounts(x uint8, y uint8, c cache) uint64 {
	key := [2]uint8{x, y}

	if cached, ok := c.partitionCounts[key]; ok {
		return cached
	}
	if x == 0 || y == 1 {
		return uint64(1)
	}
	if y == 0 {
		return 0
	}

	// The general logic is that to calc the number of partitions, you can split
	// them up into either partitions where at least one of the numbers being added
	// is 0, meaning you are effectively conducting p(x, y-1), or partitions where
	// none of the numbers added are zero, meaning that you are figuring out how to
	// add y numbers together to get (x-y) since every number needs to be at least 1
	// so that is simplified to p(x-y,y) so long as x>=y
	withSumOfZero := partitionCounts(x, y-1, c)
	withoutSumOfZero := uint64(0)
	if x >= y {
		withoutSumOfZero = partitionCounts(x-y, y, c)
	}

	c.partitionCounts[key] = withSumOfZero + withoutSumOfZero
	return withSumOfZero + withoutSumOfZero
}

// sumFactorization returns a list of all possible numbers that can sum up to x
func sumFactorization(x, length uint16, c cache) []map[uint16]uint16 {
	key := [2]uint16{x, length}
	if cached, ok := c.sumFactorizations[key]; ok {
		out := make([]map[uint16]uint16, len(cached))
		for index, factorization := range cached {
			out[index] = maps.Clone(factorization)
		}
		return out
	}

	if length == 0 {
		if x == 0 {
			return []map[uint16]uint16{{}}
		}
		return nil
	}
	if x < length {
		return nil
	}
	if length == 1 {
		result := []map[uint16]uint16{{x: 1}}
		c.sumFactorizations[key] = []map[uint16]uint16{maps.Clone(result[0])}
		return result
	}

	result := []map[uint16]uint16{}
	for i := uint16(1); i <= x; i++ {
		if x-i < length-1 {
			break
		}
		factors := sumFactorization(x-i, length-1, c)
		for _, factor := range factors {
			if len(factor) > 0 && slices.Min(slices.Collect(maps.Keys(factor))) < i {
				continue
			}
			next := maps.Clone(factor)
			next[i]++
			result = append(result, next)
		}
	}
	c.sumFactorizations[key] = make([]map[uint16]uint16, len(result))
	for index, factorization := range result {
		c.sumFactorizations[key][index] = maps.Clone(factorization)
	}
	return result
}

// weightedPartitionCount computes the sum of (1/∏ci!) over all unordered partitions
// of x into exactly k positive integers, each ≥ minVal.
// This allows navigating the partition space in O(x*k) time without enumerating
// all partitions explicitly.
func weightedPartitionCount(x, k, minVal uint16, c cache) *big.Rat {
	if k == 0 {
		if x == 0 {
			return big.NewRat(1, 1)
		}
		return big.NewRat(0, 1)
	}
	if x < k*minVal {
		return big.NewRat(0, 1)
	}
	if k == 1 {
		if x >= minVal {
			return big.NewRat(1, 1)
		}
		return big.NewRat(0, 1)
	}

	key := [3]uint16{x, k, minVal}
	if cached, ok := c.weightedPartCounts[key]; ok {
		return new(big.Rat).Set(cached)
	}

	result := new(big.Rat)
	// For each possible smallest distinct value v starting at minVal
	for v := minVal; uint32(v)*uint32(k) <= uint32(x); v++ {
		// For each count c of v (1 to k, subject to c*v ≤ x)
		for count := uint16(1); count <= k && uint32(count)*uint32(v) <= uint32(x); count++ {
			newX := x - count*v
			newK := k - count
			var subCount *big.Rat
			if newK == 0 {
				if newX == 0 {
					subCount = big.NewRat(1, 1)
				} else {
					continue
				}
			} else {
				subCount = weightedPartitionCount(newX, newK, v+1, c)
				if subCount.Sign() == 0 {
					continue
				}
			}
			countFact := factorialBig(int64(count))
			term := new(big.Rat).SetFrac(big.NewInt(1), countFact)
			term.Mul(term, subCount)
			result.Add(result, term)
		}
	}

	c.weightedPartCounts[key] = new(big.Rat).Set(result)
	return result
}

func getDuplicateWeight(total, split, max uint16, c cache) *big.Rat {
	key := [3]uint16{total, split, max}

	if cached, ok := c.duplicateWeights[key]; ok {
		return new(big.Rat).Set(cached)
	}

	if total == 0 && split == 0 {
		return big.NewRat(1, 1)
	}
	if max == 0 || total == 0 || split == 0 {
		return big.NewRat(0, 1)
	}

	limit := min(split, total/max)

	result := new(big.Rat)
	for i := uint16(0); i <= limit; i++ {
		term := big.NewRat(1, int64(factorial(int(i), c)))
		if split >= i && total >= (i*max) {
			term = term.Mul(term, getDuplicateWeight(total-i*max, split-i, max-1, c))
		} else {
			term = big.NewRat(0, 1)
		}

		result.Add(result, term)
	}
	c.duplicateWeights[key] = new(big.Rat).Set(result)
	return result
}

// waysForStackCount calculates the number of ways an item can be stacked
// @param stop: number of stacks to stop at
func (item Item) waysForStackCount(stop uint8, c cache) map[uint8]*big.Rat {
	stackCounts := make(map[uint8]*big.Rat)
	stackCounts[0] = big.NewRat(1, 1)
	for i := uint8(1); uint16(i) <= item.Stack && i <= stop; i++ {
		stackCounts[i-1] = getDuplicateWeight(item.Stack, uint16(i), item.Stack, c)
	}
	return stackCounts
}

func nCr(n, r uint64) uint64 {
	if r > n {
		return 0
	}
	if r > n/2 {
		r = n - r // Use symmetry: nCr = nC(n-r)
	}
	res := uint64(1)
	for i := uint64(1); i <= r; i++ {
		res = res * (n - i + 1) / i
	}
	return res
}

// TODO: Cache this
func factorialBig(n int64) *big.Int {
	result := big.NewInt(1)
	for i := int64(2); i <= n; i++ {
		result.Mul(result, big.NewInt(i))
	}
	return result
}

// combineSpreads returns a combined mapping of 2 spreads, combining their weights
func combineSpreads(s1, s2 map[uint8]*big.Rat, blankCount uint8) map[uint8]*big.Rat {
	combinedSpreads := map[uint8]*big.Rat{}
	s1Keys := slices.Collect(maps.Keys(s1))
	s2Keys := slices.Collect(maps.Keys(s2))
	slices.Sort(s1Keys)
	slices.Sort(s2Keys)

	for _, extra1 := range s1Keys {
		weight1 := new(big.Rat).Set(s1[extra1])
		for _, extra2 := range s2Keys {
			weight2 := new(big.Rat).Set(s2[extra2])
			combinedWeight := new(big.Rat)
			combinedWeight.Mul(weight1, weight2)
			newExtra := extra1 + extra2
			if newExtra > blankCount {
				break
			}
			if combinedSpreads[newExtra] == nil {
				combinedSpreads[newExtra] = new(big.Rat)
			}
			combinedSpreads[newExtra].Add(combinedSpreads[newExtra], combinedWeight)
		}
	}
	return combinedSpreads
}

// TODO: Cache This
// getWeights gets the combined weights of all items in the inventory, starting at `startingFrom`
func (i Inventory) getWeights(blankCount, startingFrom uint8) map[uint8]*big.Rat {
	combinedSpreads := map[uint8]*big.Rat{
		0: big.NewRat(1, 1),
	}

	for index, item := range i.Items {
		if index < int(startingFrom) {
			continue
		}
		if item.Stack > 1 {
			spread := item.waysForStackCount(blankCount+1, i.Cache)
			combinedSpreads = combineSpreads(combinedSpreads, spread, blankCount)
		}
	}

	return combinedSpreads
}

// TODO: Cache This
// calculateSpreadPermutationCounts calculates the number of possible inventory perms for each number of blanks
func (i Inventory) CalculateSpreadPermutationCounts() map[uint8]*big.Rat {
	deblankedInventory := i.Copy()
	deblankedInventory.collapse()
	deblankedInventory.removeBlanks()
	deblankedInventory.Sort()
	blankCount := uint8(len(i.Items) - len(deblankedInventory.Items))
	inventoryLengthFactorial := factorialBig(int64(len(i.Items)))

	result := map[uint8]*big.Rat{}

	combinedSpreads := deblankedInventory.getWeights(blankCount, 0)

	for extra := uint8(0); extra <= blankCount; extra++ {
		if combinedSpreads[extra] == nil {
			continue
		}
		blankCountFactorial := factorialBig(int64(blankCount - extra))
		permCount := new(big.Rat).SetFrac(inventoryLengthFactorial, blankCountFactorial)
		permCount.Mul(permCount, combinedSpreads[extra])

		result[blankCount-extra] = permCount
	}

	return result
}

// GetVariation gets the variation of a collapsed inventory at a given index
func (i Inventory) GetVariation(target *big.Rat) Inventory {
	spreadPermutationCounts := i.CalculateSpreadPermutationCounts()

	deblankedInventory := i.Copy()
	deblankedInventory.collapse()
	deblankedInventory.removeBlanks()
	deblankedInventory.Sort()
	blankCount := uint8(len(i.Items) - len(deblankedInventory.Items))

	// Convert to 0-based residual rank and consume branch sizes bucket-by-bucket.
	rank := new(big.Rat).Set(target)
	rank.Sub(rank, big.NewRat(1, 1))
	if rank.Sign() < 0 {
		rank = big.NewRat(0, 1)
	}

	blanks := slices.Collect(maps.Keys(spreadPermutationCounts))
	slices.Sort(blanks)
	targetBlankCount := uint8(0)

	// Bucket 1
	selectedBlankCount := false
	for _, blank := range blanks {
		branchCount := spreadPermutationCounts[blank]
		if branchCount == nil {
			continue
		}

		if rank.Cmp(branchCount) < 0 {
			targetBlankCount = blank
			selectedBlankCount = true
			break
		}
		rank.Sub(rank, branchCount)
	}
	if !selectedBlankCount {
		panic("bucket 1 failed to select blank count")
	}
	targetFillCount := blankCount - targetBlankCount

	// Bucket 2
	inventoryLengthFactorial := factorialBig(int64(len(i.Items)))
	blankCountFactorial := factorialBig(int64(targetBlankCount))
	rawPermCount := new(big.Rat).SetFrac(inventoryLengthFactorial, blankCountFactorial)
	targetSplits := []uint8{}
	prefixWeight := big.NewRat(1, 1)
	for itemIndex, item := range deblankedInventory.Items {
		weights := map[uint8]*big.Rat{0: big.NewRat(1, 1)}
		if item.Stack > 1 {
			weights = item.waysForStackCount(targetFillCount+1, i.Cache)
		}

		weightSplits := slices.Collect(maps.Keys(weights))
		slices.Sort(weightSplits)
		slices.Reverse(weightSplits)

		chosen := false

		for _, split := range weightSplits {
			if split > targetFillCount {
				continue
			}

			weight := weights[split]
			suffixWeight := deblankedInventory.getWeights(targetFillCount-split, uint8(itemIndex)+1)
			suffix := suffixWeight[targetFillCount-split]
			if suffix == nil {
				suffix = big.NewRat(0, 1)
			}

			branchCount := new(big.Rat).Set(rawPermCount)
			branchCount.Mul(branchCount, prefixWeight)
			branchCount.Mul(branchCount, weight)
			branchCount.Mul(branchCount, suffix)

			if rank.Cmp(branchCount) < 0 {
				targetSplits = append(targetSplits, split)
				targetFillCount -= split
				prefixWeight = new(big.Rat).Mul(prefixWeight, weight)
				chosen = true
				break
			}
			rank.Sub(rank, branchCount)
		}
		if !chosen {
			panic(fmt.Sprintf("bucket 2 failed at item %d (%s)", itemIndex, item.Name))
		}
	}

	// Bucket 3
	exactSuffixes := make([]*big.Rat, len(deblankedInventory.Items)+1)
	exactSuffixes[len(deblankedInventory.Items)] = big.NewRat(1, 1)
	for itemIndex := len(deblankedInventory.Items) - 1; itemIndex >= 0; itemIndex-- {
		parts := targetSplits[itemIndex] + 1
		itemContribution := big.NewRat(1, 1)
		if parts > 1 {
			itemStack := deblankedInventory.Items[itemIndex].Stack
			itemContribution = getDuplicateWeight(itemStack, uint16(parts), itemStack, i.Cache)
		}
		exactSuffixes[itemIndex] = new(big.Rat).Set(itemContribution)
		exactSuffixes[itemIndex] = exactSuffixes[itemIndex].Mul(exactSuffixes[itemIndex], exactSuffixes[itemIndex+1])
	}
	finalInventoryCombination := Inventory{}
	exactPrefix := big.NewRat(1, 1)
	for itemIndex, split := range targetSplits {
		item := deblankedInventory.Items[itemIndex]
		split++
		if split == 1 {
			finalInventoryCombination.Items = append(finalInventoryCombination.Items, item)
			continue
		}
		suffix := exactSuffixes[itemIndex+1]
		// Tree-walk the partition space: choose (value, count) pairs in ascending order
		// instead of enumerating all partitions via sumFactorization.
		remaining := item.Stack
		remParts := uint16(split)
		minVal := uint16(1)
		factMap := map[uint16]uint16{}
		for remParts > 0 {
			found := false
			for v := minVal; !found && uint32(v)*uint32(remParts) <= uint32(remaining); v++ {
				for c := uint16(1); c <= remParts && uint32(c)*uint32(v) <= uint32(remaining); c++ {
					newX := remaining - c*v
					newK := remParts - c
					var subCount *big.Rat
					if newK == 0 {
						if newX == 0 {
							subCount = big.NewRat(1, 1)
						} else {
							continue
						}
					} else {
						subCount = weightedPartitionCount(newX, newK, v+1, i.Cache)
						if subCount.Sign() == 0 {
							continue
						}
					}
					cFact := factorialBig(int64(c))
					subW := new(big.Rat).SetFrac(big.NewInt(1), cFact)
					branchCount := new(big.Rat).Mul(rawPermCount, exactPrefix)
					branchCount.Mul(branchCount, subW)
					branchCount.Mul(branchCount, subCount)
					branchCount.Mul(branchCount, suffix)
					if rank.Cmp(branchCount) < 0 {
						factMap[v] = c
						exactPrefix.Mul(exactPrefix, subW)
						remaining = newX
						remParts = newK
						minVal = v + 1
						found = true
						break
					}
					rank.Sub(rank, branchCount)
				}
			}
			if !found {
				panic(fmt.Sprintf("bucket 3 failed at item %d (%s)", itemIndex, item.Name))
			}
		}
		factVals := slices.Collect(maps.Keys(factMap))
		slices.Sort(factVals)
		for _, v := range factVals {
			newItem := Item{Name: item.Name, Stack: v}
			for range factMap[v] {
				finalInventoryCombination.Items = append(finalInventoryCombination.Items, newItem)
			}
		}
	}
	finalInventoryCombination.addBlanks(len(i.Items))
	finalInventoryCombination.Sort()

	// Bucket 4
	finalInventory := Inventory{}

	for len(finalInventoryCombination.Items) > 0 {
		chosen := false
		checkedItems := Inventory{}
		for index, item := range finalInventoryCombination.Items {
			if checkedItems.Contains(item) {
				continue
			}
			checkedItems.Items = append(checkedItems.Items, item)
			checkingInventory := finalInventoryCombination.Copy()
			checkingInventory.Items = slices.Delete(checkingInventory.Items, index, index+1)
			branchCount := checkingInventory.getPermutationCount()

			if rank.Cmp(branchCount) < 0 {
				finalInventoryCombination.Items = slices.Delete(finalInventoryCombination.Items, index, index+1)
				finalInventory.Items = append(finalInventory.Items, item)
				chosen = true
				break
			}
			rank.Sub(rank, branchCount)
		}
		if !chosen {
			panic("bucket 4 failed to select next item")
		}
	}

	return finalInventory
}

func (i Inventory) GetIndex() *big.Rat {
	deblankedFullInventory := i.Copy()
	deblankedFullInventory.removeBlanks()
	deblankedFullInventory.Sort()
	deblankedInventory := i.Copy()
	deblankedInventory.collapse()
	deblankedInventory.removeBlanks()
	deblankedInventory.Sort()

	blankCount := uint8(len(i.Items) - len(deblankedFullInventory.Items))
	fillCount := uint8(len(deblankedFullInventory.Items) - len(deblankedInventory.Items))

	// Bucket 1
	index := new(big.Rat)
	inventoryCounts := i.CalculateSpreadPermutationCounts()
	blankKeys := slices.Collect(maps.Keys(inventoryCounts))
	slices.Sort(blankKeys)
	for _, blanks := range blankKeys {
		if blanks >= blankCount {
			break
		}
		if inventoryCounts[blanks] == nil {
			continue
		}
		index.Add(index, inventoryCounts[blanks])
	}

	// Bucket 2
	splits := make([]uint8, len(deblankedInventory.Items))
	for itemIndex, item := range deblankedInventory.Items {
		for _, splitItem := range deblankedFullInventory.Items {
			if splitItem.Name == item.Name {
				splits[itemIndex]++
			}
		}
		// convert "number of stacks used" into "extra"
		if splits[itemIndex] > 0 {
			splits[itemIndex]--
		}
	}

	inventoryLengthFactorial := factorialBig(int64(len(i.Items)))
	blankCountFactorial := factorialBig(int64(blankCount))
	rawPermCount := new(big.Rat).SetFrac(inventoryLengthFactorial, blankCountFactorial)

	prefixWeight := big.NewRat(1, 1)
	targetFillCount := fillCount

	for itemIndex, item := range deblankedInventory.Items {
		actualSplit := splits[itemIndex]

		weights := map[uint8]*big.Rat{0: big.NewRat(1, 1)}
		if item.Stack > 1 {
			weights = item.waysForStackCount(targetFillCount+1, i.Cache)
		}

		weightSplits := slices.Collect(maps.Keys(weights))
		slices.Sort(weightSplits)
		slices.Reverse(weightSplits)

		for _, split := range weightSplits {
			if split > targetFillCount {
				continue
			}

			weight := weights[split]

			suffixWeight := deblankedInventory.getWeights(targetFillCount-split, uint8(itemIndex)+1)
			suffix := suffixWeight[targetFillCount-split]
			if suffix == nil {
				suffix = big.NewRat(0, 1)
			}

			permCount := new(big.Rat).Set(rawPermCount)
			permCount.Mul(permCount, prefixWeight)
			permCount.Mul(permCount, weight)
			permCount.Mul(permCount, suffix)

			if split == actualSplit {
				prefixWeight = new(big.Rat).Mul(prefixWeight, weight)
				targetFillCount -= split
				break
			}

			// Earlier branch: add it to the rank.
			index.Add(index, permCount)
		}
	}

	// Bucket 3
	exactSuffixes := make([]*big.Rat, len(deblankedInventory.Items)+1)
	exactSuffixes[len(deblankedInventory.Items)] = big.NewRat(1, 1)
	for itemIndex := len(deblankedInventory.Items) - 1; itemIndex >= 0; itemIndex-- {
		parts := splits[itemIndex] + 1
		itemContribution := big.NewRat(1, 1)
		if parts > 1 {
			itemStack := deblankedInventory.Items[itemIndex].Stack
			itemContribution = getDuplicateWeight(itemStack, uint16(parts), itemStack, i.Cache)
		}
		exactSuffixes[itemIndex] = new(big.Rat).Mul(itemContribution, exactSuffixes[itemIndex+1])
	}

	workingInventory := deblankedFullInventory.Copy()
	exactPrefix := big.NewRat(1, 1)

	for itemIndex, extra := range splits {
		item := deblankedInventory.Items[itemIndex]
		parts := extra + 1
		if parts == 1 {
			continue
		}
		suffix := exactSuffixes[itemIndex+1]

		// Build actual factorization for this item from workingInventory
		actualStacks := map[uint16]uint16{}
		for _, invItem := range workingInventory.Items {
			if invItem.Name == item.Name {
				actualStacks[invItem.Stack]++
			}
		}

		// Tree-walk: count branches that precede the actual factorization
		remaining := item.Stack
		remParts := uint16(parts)
		minVal := uint16(1)
		for remParts > 0 {
			// Find actual (v, c) at this level: smallest key in actualStacks >= minVal
			actualV := uint16(0)
			for sv := range actualStacks {
				if sv >= minVal && (actualV == 0 || sv < actualV) {
					actualV = sv
				}
			}
			actualC := actualStacks[actualV]

			done := false
			for v := minVal; !done && uint32(v)*uint32(remParts) <= uint32(remaining); v++ {
				for c := uint16(1); c <= remParts && uint32(c)*uint32(v) <= uint32(remaining); c++ {
					newX := remaining - c*v
					newK := remParts - c

					if v == actualV && c == actualC {
						// Matching branch: enter it
						cFact := factorialBig(int64(c))
						subW := new(big.Rat).SetFrac(big.NewInt(1), cFact)
						exactPrefix.Mul(exactPrefix, subW)
						remaining = newX
						remParts = newK
						minVal = v + 1
						delete(actualStacks, actualV)
						for r := uint16(0); r < c; r++ {
							toRemove := Item{Name: item.Name, Stack: v}
							for idx, invItem := range workingInventory.Items {
								if invItem == toRemove {
									workingInventory.Items = slices.Delete(workingInventory.Items, idx, idx+1)
									break
								}
							}
						}
						done = true
						break
					}

					// Count this branch (it precedes the actual)
					var subCount *big.Rat
					if newK == 0 {
						if newX == 0 {
							subCount = big.NewRat(1, 1)
						} else {
							continue
						}
					} else {
						subCount = weightedPartitionCount(newX, newK, v+1, i.Cache)
						if subCount.Sign() == 0 {
							continue
						}
					}
					cFact := factorialBig(int64(c))
					subW := new(big.Rat).SetFrac(big.NewInt(1), cFact)
					permCount := new(big.Rat).Mul(rawPermCount, exactPrefix)
					permCount.Mul(permCount, subW)
					permCount.Mul(permCount, subCount)
					permCount.Mul(permCount, suffix)
					index.Add(index, permCount)
				}
			}
		}
	}

	// Bucket 4
	inventoryCombination := i.Copy()
	inventoryCombination.Sort()

	for len(inventoryCombination.Items) > 0 {
		itemBeingChecked := len(i.Items) - len(inventoryCombination.Items)
		actualItem := i.Items[itemBeingChecked]

		checkedItems := Inventory{}

		for itemIndex, item := range inventoryCombination.Items {
			if checkedItems.Contains(item) {
				continue
			}
			checkedItems.Items = append(checkedItems.Items, item)

			checkingInventory := inventoryCombination.Copy()
			checkingInventory.Items = slices.Delete(checkingInventory.Items, itemIndex, itemIndex+1)
			permCount := checkingInventory.getPermutationCount()

			if item == actualItem {
				inventoryCombination.Items = slices.Delete(inventoryCombination.Items, itemIndex, itemIndex+1)
				break
			}
			index.Add(index, permCount)
		}
	}

	index.Add(index, big.NewRat(1, 1))
	return index
}

func (i Inventory) getPermutationCount() *big.Rat {
	abstract := i.getAbstract()
	numerator := new(big.Int)
	denominator := big.NewInt(1)
	for _, itemCount := range abstract {
		bigItemCount := big.NewInt(int64(itemCount))
		numerator.Add(numerator, bigItemCount)
		denominator.Mul(denominator, factorialBig(int64(itemCount)))
	}
	numerator = factorialBig(numerator.Int64())
	return new(big.Rat).SetFrac(numerator, denominator)
}

// Copy returns a copy of the Inventory
func (i Inventory) Copy() Inventory {
	copiedInventory := Inventory{
		Items: make([]Item, len(i.Items)),
		Cache: i.Cache,
	}
	copy(copiedInventory.Items, i.Items)
	return copiedInventory
}

func (i *Inventory) Sort() {
	sort.Slice(i.Items, func(j, k int) bool {
		jBlank := i.Items[j].Stack < 1
		kBlank := i.Items[k].Stack < 1
		if jBlank != kBlank {
			return !jBlank // blanks go to the end
		}
		if i.Items[j].Name != i.Items[k].Name {
			return i.Items[j].Name < i.Items[k].Name
		}
		return i.Items[j].Stack < i.Items[k].Stack
	})
}

func (i Inventory) Print() {
	for _, item := range i.Items {
		print(item.Name, " x", item.Stack, " ")
	}
	println()
}

// getAbstract turns the inventory into just a slice of Item stacks
// It's used for getting standardized permutation counts
func (i Inventory) getAbstract() []uint16 {
	var abstractRepresentation []uint16
	abstractMap := make(map[Item]uint16)
	for _, item := range i.Items {
		abstractMap[item]++
	}
	for _, count := range abstractMap {
		abstractRepresentation = append(abstractRepresentation, count)
	}
	slices.Sort(abstractRepresentation)
	return abstractRepresentation
}

// calculateSpread generates all the different combinations of the numbers
// of blank spaces to spread each stack by
func calculateSpread(blanks uint8, stacks uint8) [][]uint8 {
	var spread [][]uint8
	if blanks < 1 {
		return [][]uint8{make([]uint8, stacks)}
	}
	if stacks < 2 {
		return [][]uint8{{blanks}}
	}
	blank := blanks
	for {
		nextBlanks := blanks - blank
		for _, row := range calculateSpread(nextBlanks, stacks-1) {
			spread = append(spread, append([]uint8{blank}, row...))
		}
		if blank < 1 {
			return spread
		}
		blank--
	}
}

func spreadSpares(spares uint8, blanks uint8) [][]uint8 {
	var spread [][]uint8
	if spares < 1 {
		return [][]uint8{make([]uint8, blanks)}
	}
	if blanks < 2 {
		return [][]uint8{{spares}}
	}

	nextStack := spares
	for nextStack > 0 {
		for _, row := range spreadSpares(spares-nextStack, blanks-1) {
			if row[0] <= nextStack {
				spread = append(spread, append([]uint8{nextStack}, row...))
			}
		}

		nextStack--
	}
	return spread
}

func factorial(n int, c cache) int {
	if c.smallFactorial[n] > 0 {
		return c.smallFactorial[n]
	}
	if n <= 1 {
		return 1
	}
	result := 1
	for i := 2; i <= n; i++ {
		result *= i
	}
	c.smallFactorial[n] = result
	return result
}

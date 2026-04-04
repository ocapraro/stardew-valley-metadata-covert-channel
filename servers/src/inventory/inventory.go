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
	smallFactorial    map[int]int
	duplicateWeights  map[[3]uint16]*big.Rat
	partitionCounts   map[[2]uint8]uint64
	sumFactorizations map[uint8][][]uint8
	perms             map[[2]uint8][][]uint8
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
		smallFactorial:    make(map[int]int),
		duplicateWeights:  make(map[[3]uint16]*big.Rat),
		partitionCounts:   make(map[[2]uint8]uint64),
		sumFactorizations: make(map[uint8][][]uint8),
		perms:             make(map[[2]uint8][][]uint8),
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
func sumFactorization(x uint8, c cache) [][]uint8 {
	if cached, ok := c.sumFactorizations[x]; ok {
		out := make([][]uint8, len(cached))
		copy(out, cached)
		return out
	}
	result := [][]uint8{{x}}
	for i := uint8(1); i < x; i++ {
		for _, subResult := range sumFactorization(x-i, c) {
			if slices.Min(subResult) >= i {
				result = append(result, append([]uint8{i}, subResult...))
			}
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return len(result[i]) < len(result[j])
	})
	c.sumFactorizations[x] = make([][]uint8, len(result))
	copy(c.sumFactorizations[x], result)
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
		term.Mul(term, getDuplicateWeight(total-i*max, split-i, max-1, c))
		result.Add(result, term)
	}
	c.duplicateWeights[key] = new(big.Rat).Set(result)
	return result
}

// waysForStackCount calculates the number of ways an item can be stacked
func (item Item) waysForStackCount(stop uint8, c cache) map[uint8]*big.Rat {
	stackCounts := make(map[uint8]*big.Rat)
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

// calculateSpreadPermutationCounts calculates the number of possible inventory perms for each number of blanks
func (i Inventory) CalculateSpreadPermutationCounts() map[uint8]*big.Rat {
	deblankedInventory := i.Copy()
	deblankedInventory.collapse()
	deblankedInventory.removeBlanks()
	deblankedInventory.Sort()
	blankCount := uint8(len(i.Items) - len(deblankedInventory.Items))
	inventoryLengthFactorial := factorialBig(int64(len(i.Items)))

	result := map[uint8]*big.Rat{}

	combinedSpreads := map[uint8]*big.Rat{
		0: big.NewRat(1, 1),
	}

	for _, item := range deblankedInventory.Items {
		if item.Stack > 1 {
			spread := item.waysForStackCount(blankCount+1, i.Cache)
			combinedSpreads = combineSpreads(combinedSpreads, spread, blankCount)
		}
	}

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

type bounds struct {
	Upper big.Int
	Lower big.Int
}

// calculateCombinations gets each way the stacks in an inventory can be split
func (i Inventory) calculateCombinations() []Inventory {
	println("Calculating combinations")
	deblankedInventory := i.Copy()
	deblankedInventory.collapse()
	deblankedInventory.removeBlanks()
	deblankedInventory.Sort()
	// blankCount := len(i.Items) - len(deblankedInventory.Items)

	inventoryCombinations := []Inventory{
		deblankedInventory.Copy(),
	}

	// var stacks []Item
	// for _, item := range deblankedInventory.Items {
	// 	if item.Stack < 2 {
	// 		continue
	// 	}
	// 	stacks = append(stacks, item)
	// }

	// for blanks := 1; blanks <= blankCount; blanks++ {
	// 	spreads := calculateSpread(uint8(blanks), uint8(len(stacks)))
	// 	for _, spread := range spreads {
	// 		var spreadCombinations [][][]uint16
	// 		for stackIndex, fillsNeeded := range spread {
	// 			stack := stacks[stackIndex]
	// 			// Make sure the stack can actually fill the number of slots assigned to it here
	// 			if stack.Stack <= fillsNeeded {
	// 				continue
	// 			}
	// 			// Add the initial slot
	// 			fillsNeeded++

	// 			spareItems := stack.Stack - fillsNeeded
	// 			spareSpread := spreadSpares(spareItems, fillsNeeded)
	// 			for row := range spareSpread {
	// 				for s := range spareSpread[row] {
	// 					spareSpread[row][s]++
	// 				}
	// 			}
	// 			spreadCombinations = append(spreadCombinations, spareSpread)

	// 		}
	// 		stackCombinations := []Inventory{{}}
	// 		stackIndex := 0
	// 		for len(spreadCombinations) > 0 {
	// 			previousStackCombinations := stackCombinations
	// 			stackCombinations = []Inventory{}

	// 			spreads := spreadCombinations[0]
	// 			currentStack := stacks[stackIndex]
	// 			spreadCombinations = spreadCombinations[1:]
	// 			stackIndex++
	// 			for _, currentSpread := range spreads {
	// 				var stackCombinationAddition Inventory
	// 				for _, quantity := range currentSpread {
	// 					stackCombinationAddition.Items = append(stackCombinationAddition.Items, Item{
	// 						Name:  currentStack.Name,
	// 						Stack: quantity,
	// 					})
	// 				}
	// 				for _, stackCombination := range previousStackCombinations {
	// 					newStackCombination := stackCombination.Copy()
	// 					newStackCombination.Items = append(newStackCombination.Items, stackCombinationAddition.Copy().Items...)
	// 					stackCombinations = append(stackCombinations, newStackCombination)
	// 				}
	// 			}

	// 		}
	// 		for _, stackCombination := range stackCombinations {
	// 			for _, item := range deblankedInventory.Items {
	// 				if item.Stack > 1 {
	// 					continue
	// 				}
	// 				stackCombination.Items = append(stackCombination.Items, item)
	// 			}
	// 			stackCombination = stackCombination.Copy()
	// 			stackCombination.Sort()
	// 			inventoryCombinations = append(inventoryCombinations, stackCombination)
	// 		}
	// 	}
	// }

	return inventoryCombinations
}

func (i Inventory) getPermutationCount() uint64 {
	abstract := i.getAbstract()
	numerator := 0
	denominator := 1
	for _, itemCount := range abstract {
		numerator += int(itemCount)
		denominator *= factorial(int(itemCount), i.Cache)
	}
	numerator = factorial(numerator, i.Cache)
	return uint64(numerator / denominator)
}

func (i Inventory) getCounts() map[string]uint64 {
	println("Getting counts")
	inventories := i.calculateCombinations()
	println("Calculated combos")
	sort.Slice(inventories, func(j, k int) bool {
		return fmt.Sprint(inventories[j].Items) < fmt.Sprint(inventories[k].Items)
	})
	invetoryCounts := make(map[string]uint64)
	abstractCounts := make(map[string]uint64)

	for _, inventory := range inventories {
		inventory.addBlanks(len(i.Items))
		abstract := inventory.getAbstract()
		if abstractCounts[string(abstract)] > 0 {
			invetoryCounts[inventory.encode()] = abstractCounts[string(abstract)]
			continue
		}
		perms := inventory.getPermutationCount()
		invetoryCounts[inventory.encode()] = perms
		abstractCounts[string(abstract)] = perms
	}
	return invetoryCounts
}

// GetVariation gets the variation of a collapsed inventory at a given index
func (i Inventory) GetVariation(target uint64) Inventory {
	println("Getting variation")
	inventoryCounts := i.getCounts()
	println("Counts gotten")
	var encodedInventories []string
	for inventory := range inventoryCounts {
		encodedInventories = append(encodedInventories, inventory)
	}
	slices.Sort(encodedInventories)

	count := uint64(0)
	encodedInventory := ""
	for _, inventory := range encodedInventories {
		currentInventoryCount := inventoryCounts[inventory]
		if count+currentInventoryCount > target {
			encodedInventory = inventory
			break
		}
		count += currentInventoryCount
	}
	if len(encodedInventory) < 1 {
		return Inventory{}
	}

	finalInventory := Inventory{}
	workingInventory := decodeInventory(encodedInventory)
	workingInventory.Sort()

	for len(workingInventory.Items) > 0 {
		checkedItems := Inventory{}
		for index, item := range workingInventory.Items {
			if checkedItems.Contains(item) {
				continue
			}
			checkedItems.Items = append(checkedItems.Items, item)
			checkingInventory := workingInventory.Copy()
			checkingInventory.Items = slices.Delete(checkingInventory.Items, index, index+1)
			permCount := checkingInventory.getPermutationCount()
			if count+permCount > target {
				workingInventory.Items = slices.Delete(workingInventory.Items, index, index+1)
				finalInventory.Items = append(finalInventory.Items, item)
				break
			}
			count += permCount
		}
	}
	return finalInventory
}

func (i Inventory) GetIndex() uint64 {
	collapsedInventory := i.Copy()
	collapsedInventory.collapse()
	collapsedInventory.Sort()

	sortedInventory := i.Copy()
	sortedInventory.Sort()
	encodedInventory := sortedInventory.encode()

	inventoryCounts := collapsedInventory.getCounts()
	var encodedInventories []string
	for inventory := range inventoryCounts {
		encodedInventories = append(encodedInventories, inventory)
	}
	slices.Sort(encodedInventories)
	if !slices.Contains(encodedInventories, encodedInventory) {
		return 0
	}

	count := uint64(0)
	for _, inventory := range encodedInventories {
		if inventory == encodedInventory {
			break
		}
		count += inventoryCounts[inventory]
	}

	finalInventory := Inventory{}
	workingInventory := decodeInventory(encodedInventory)
	workingInventory.Sort()

	for len(workingInventory.Items) > 0 {
		checkedItems := Inventory{}
		for index, item := range workingInventory.Items {
			if checkedItems.Contains(item) {
				continue
			}
			checkedItems.Items = append(checkedItems.Items, item)
			checkingInventory := workingInventory.Copy()
			checkingInventory.Items = slices.Delete(checkingInventory.Items, index, index+1)
			permCount := checkingInventory.getPermutationCount()
			if slices.Equal(i.Items[:len(finalInventory.Items)+1], append(finalInventory.Items, item)) {
				workingInventory.Items = slices.Delete(workingInventory.Items, index, index+1)
				finalInventory.Items = append(finalInventory.Items, item)
				break
			}
			count += permCount
		}
	}
	return count
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
func (i Inventory) getAbstract() []uint8 {
	var abstractRepresentation []uint8
	abstractMap := make(map[Item]uint8)
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

package inventory

import (
	"slices"
	"sort"
)

type Item struct {
	Name  string
	Stack uint8
}

type Inventory struct {
	Items []Item
}

// Contains checks to see if an Item is in an Inventory
func (i Inventory) Contains(target Item) bool {
	return slices.Contains(i.Items, target)
}

// collapse combines all stacks of the same item
func (i *Inventory) collapse() {
	length := len(i.Items)
	counts := make(map[string]uint8)
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

// calculateCombinations gets each way the stacks in an inventory can be split
func (i Inventory) calculateCombinations() []Inventory {
	deblankedInventory := i.Copy()
	deblankedInventory.collapse()
	deblankedInventory.removeBlanks()
	deblankedInventory.Sort()
	blankCount := len(i.Items) - len(deblankedInventory.Items)

	inventoryCombinations := []Inventory{
		deblankedInventory.Copy(),
	}

	var stacks []Item
	for _, item := range deblankedInventory.Items {
		if item.Stack < 2 {
			continue
		}
		stacks = append(stacks, item)
	}

	for blanks := 1; blanks <= blankCount; blanks++ {
		spreads := calculateSpread(uint8(blanks), uint8(len(stacks)))
		for _, spread := range spreads {
			var spreadCombinations [][][]uint8
			for stackIndex, fillsNeeded := range spread {
				stack := stacks[stackIndex]
				// Make sure the stack can actually fill the number of slots assigned to it here
				if stack.Stack <= fillsNeeded {
					continue
				}
				// Add the initial slot
				fillsNeeded++

				spareItems := stack.Stack - fillsNeeded
				spareSpread := spreadSpares(spareItems, fillsNeeded)
				for row := range spareSpread {
					for s := range spareSpread[row] {
						spareSpread[row][s]++
					}
				}
				spreadCombinations = append(spreadCombinations, spareSpread)

			}
			stackCombinations := []Inventory{{}}
			stackIndex := 0
			for len(spreadCombinations) > 0 {
				previousStackCombinations := stackCombinations
				stackCombinations = []Inventory{}

				spreads := spreadCombinations[0]
				currentStack := stacks[stackIndex]
				spreadCombinations = spreadCombinations[1:]
				stackIndex++
				for _, currentSpread := range spreads {
					var stackCombinationAddition Inventory
					for _, quantity := range currentSpread {
						stackCombinationAddition.Items = append(stackCombinationAddition.Items, Item{
							Name:  currentStack.Name,
							Stack: quantity,
						})
					}
					for _, stackCombination := range previousStackCombinations {
						newStackCombination := stackCombination.Copy()
						newStackCombination.Items = append(newStackCombination.Items, stackCombinationAddition.Copy().Items...)
						stackCombinations = append(stackCombinations, newStackCombination)
					}
				}

			}
			for _, stackCombination := range stackCombinations {
				for _, item := range deblankedInventory.Items {
					if item.Stack > 1 {
						continue
					}
					stackCombination.Items = append(stackCombination.Items, item)
				}
				stackCombination = stackCombination.Copy()
				stackCombination.Sort()
				inventoryCombinations = append(inventoryCombinations, stackCombination)
			}
		}
	}

	return inventoryCombinations
}

// GetVariation gets the variation of a collapsed inventory at a given index
// func (i Inventory) GetVariation(index uint64) Inventory {

// }

// Copy returns a copy of the Inventory
func (i Inventory) Copy() Inventory {
	copiedInventory := Inventory{
		Items: make([]Item, len(i.Items)),
	}
	copy(copiedInventory.Items, i.Items)
	return copiedInventory
}

func (i *Inventory) Sort() {
	sort.Slice(i.Items, func(j, k int) bool {
		return i.Items[j].Name < i.Items[k].Name
	})
}

func (i Inventory) Print() {
	for _, item := range i.Items {
		print(item.Name, " x", item.Stack, " ")
	}
	println()
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

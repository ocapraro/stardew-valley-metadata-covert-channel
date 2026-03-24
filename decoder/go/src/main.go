package main

import (
	"decoder/inventory"
)

func main() {
	currentInventory := inventory.Inventory{
		Items: []inventory.Item{
			{
				Name:  "Axe",
				Stack: 1,
			},
			{
				Name:  "Hoe",
				Stack: 1,
			},
			{
				Name:  "WateringCan",
				Stack: 1,
			},
			{
				Name:  "Pickaxe",
				Stack: 1,
			},
			{
				Name:  "Parsnip Seeds",
				Stack: 10,
			},
			{
				Name:  "MeleeWeapon",
				Stack: 1,
			},
			{
				Name:  "Parsnip Seeds",
				Stack: 5,
			},
			{
				Name:  "Blank",
				Stack: 0,
			},
			{
				Name:  "Blank",
				Stack: 0,
			},
			{
				Name:  "Blank",
				Stack: 0,
			},
			{
				Name:  "Blank",
				Stack: 0,
			},
			{
				Name:  "Blank",
				Stack: 0,
			},
		},
	}

	// msg := "Lily"
	// numberMsg := inventory.TextToNumber(msg)
	// println(numberMsg)

	combos := currentInventory.CalculateCombinations()
	// for _, i := range combos {
	// 	i.Print()
	// }
	println(len(combos))
}

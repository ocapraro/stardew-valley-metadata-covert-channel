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
	println(len(currentInventory.Items))

	counts := currentInventory.GetCounts()
	total := 0
	for key, count := range counts {
		total += int(count)
		println(key, count)
		println("\n")
		// println(len(strings.Split(key, ",")))
	}
	println(total)
}

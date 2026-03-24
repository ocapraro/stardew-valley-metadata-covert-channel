package main

import (
	"decoder/inventory"
	"fmt"
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

	msg := "Lily"
	numberMsg := inventory.TextToNumber(msg)
	fmt.Printf("%s : %d\n", msg, numberMsg)
	// println(len(currentInventory.Items))
	variation := currentInventory.GetVariation(numberMsg)
	variation.Print()
	msgNumber := variation.GetIndex()
	newMsg := inventory.NumberToText(msgNumber)
	fmt.Printf("%d : %s\n", msgNumber, newMsg)

	// inventory.StartEncoder()
}

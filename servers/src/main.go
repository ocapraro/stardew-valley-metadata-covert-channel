package main

import (
	"fmt"
	"stardewChannel/inventory"
)

func main() {
	currentInventory := inventory.Inventory{
		Items: []inventory.Item{
			{
				Name:  "Parsnip Seeds",
				Stack: 1,
			},
			{
				Name:  "Axe",
				Stack: 1,
			},
			{
				Name:  "Pickaxe",
				Stack: 1,
			},
			{
				Name:  "Watering Can",
				Stack: 1,
			},
			{
				Name:  "Parsnip Seeds",
				Stack: 2,
			},
			{
				Name:  "Parsnip Seeds",
				Stack: 1,
			},
			{
				Name:  "Blank",
				Stack: 0,
			},
			{
				Name:  "Scythe",
				Stack: 1,
			},
			{
				Name:  "Blank",
				Stack: 0,
			},
			{
				Name:  "Parsnip Seeds",
				Stack: 5,
			},
			{
				Name:  "Parsnip Seeds",
				Stack: 6,
			},
			{
				Name:  "Hoe",
				Stack: 1,
			},
		},
	}

	msg := "Test"
	numberMsg := inventory.TextToNumber(msg)
	fmt.Printf("%s : %d\n", msg, numberMsg)
	// println(len(currentInventory.Items))
	variation := currentInventory.GetVariation(numberMsg)
	variation.Print()
	msgNumber := currentInventory.GetIndex()
	newMsg := inventory.NumberToText(msgNumber)
	fmt.Printf("%d : %s\n", msgNumber, newMsg)

	// inventory.StartEncoder()
	// inventory.StartDecoder()
}

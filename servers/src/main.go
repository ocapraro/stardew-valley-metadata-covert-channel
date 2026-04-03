package main

import (
	"fmt"
	"math"
	"math/big"
	"stardewChannel/inventory"
)

func logBigInt(n *big.Int) float64 {
	f, _ := new(big.Float).SetInt(n).Float64()
	return math.Log(f)
}

func main() {
	currentInventory := inventory.Inventory{
		Items: []inventory.Item{
			{
				Name:  "Parsnip Seeds",
				Stack: 85,
			},
			{
				Name:  "Apple",
				Stack: 99,
			},
			{
				Name:  "Pickaxe",
				Stack: 99,
			},
			{
				Name:  "Watering Can",
				Stack: 99,
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
				Stack: 99,
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
				Stack: 2,
			},
			{
				Name:  "Parsnip Seeds2",
				Stack: 85,
			},
			{
				Name:  "Apple2",
				Stack: 99,
			},
			{
				Name:  "Pickaxe2",
				Stack: 99,
			},
			{
				Name:  "Watering Can2",
				Stack: 99,
			},
			{
				Name:  "Parsnip Seeds2",
				Stack: 2,
			},
			{
				Name:  "Parsnip Seeds2",
				Stack: 1,
			},
			{
				Name:  "ahaaa",
				Stack: 99,
			},
			{
				Name:  "Scythe2",
				Stack: 99,
			},
			{
				Name:  "aaaga",
				Stack: 99,
			},
			{
				Name:  "Parsnip Seeds2",
				Stack: 5,
			},
			{
				Name:  "Parsnip Seeds2",
				Stack: 6,
			},
			{
				Name:  "Hoe2",
				Stack: 99,
			},
			{
				Name:  "Parsnip Seeds3",
				Stack: 85,
			},
			{
				Name:  "Apple3",
				Stack: 99,
			},
			{
				Name:  "Pickaxe3",
				Stack: 99,
			},
			{
				Name:  "Watering Can3",
				Stack: 99,
			},
			{
				Name:  "Parsnip Seeds3",
				Stack: 3,
			},
			{
				Name:  "Parsnip Seeds3",
				Stack: 1,
			},
			{
				Name:  "aaaa",
				Stack: 99,
			},
			{
				Name:  "Scythe3",
				Stack: 99,
			},
			{
				Name:  "test",
				Stack: 99,
			},
			{
				Name:  "Parsnip Seeds3",
				Stack: 5,
			},
			{
				Name:  "Parsnip Seeds3",
				Stack: 6,
			},
			{
				Name:  "Hoe3",
				Stack: 99,
			},
		},
	}

	// msg := "Test"
	// numberMsg := inventory.TextToNumber(msg)
	// fmt.Printf("%s : %d\n", msg, numberMsg)
	fmt.Println(len(currentInventory.Items))
	// variation := currentInventory.GetVariation(numberMsg)
	// variation.Print()
	// msgNumber := currentInventory.GetIndex()
	// newMsg := inventory.NumberToText(msgNumber)
	// fmt.Printf("%d : %s\n", msgNumber, newMsg)
	// fmt.Println(currentInventory.CalculateSpreadCombinationCounts())

	// blankBounds := currentInventory.CalculateSpreadPermutationBounds()
	// lowerCount := new(big.Int)
	// upperCount := new(big.Int)
	// for _, bounds := range blankBounds {
	// 	lowerCount.Add(lowerCount, &bounds.Lower)
	// 	upperCount.Add(upperCount, &bounds.Upper)
	// }
	// lowerBitCount := logBigInt(lowerCount) / logBigInt(big.NewInt(2))
	// upperBitCount := logBigInt(upperCount) / logBigInt(big.NewInt(2))
	// fmt.Println(lowerBitCount)
	// fmt.Println(upperBitCount)

	// inventory.StartEncoder()
	// inventory.StartDecoder()
}

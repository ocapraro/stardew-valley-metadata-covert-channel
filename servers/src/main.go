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
				Stack: 99,
			},
			{
				Name:  "Cauliflower Seeds",
				Stack: 99,
			},
			{
				Name:  "Potato Seeds",
				Stack: 99,
			},
			{
				Name:  "Strawberry Seeds",
				Stack: 99,
			},
			{
				Name:  "Melon Seeds",
				Stack: 99,
			},
			{
				Name:  "Tomato Seeds",
				Stack: 99,
			},
			{
				Name:  "Blueberry Seeds",
				Stack: 99,
			},
			{
				Name:  "Pumpkin Seeds",
				Stack: 99,
			},
			{
				Name:  "Bok Choy Seeds",
				Stack: 99,
			},
			{
				Name:  "Yam Seeds",
				Stack: 99,
			},
			{
				Name:  "Cranberry Seeds",
				Stack: 99,
			},
			{
				Name:  "Artichoke Seeds",
				Stack: 99,
			},
			{
				Name:  "Parsnip",
				Stack: 99,
			},
			{
				Name:  "Cauliflower",
				Stack: 99,
			},
			{
				Name:  "Potato",
				Stack: 99,
			},
			{
				Name:  "Strawberry",
				Stack: 99,
			},
			{
				Name:  "Melon",
				Stack: 99,
			},
			{
				Name:  "Tomato",
				Stack: 99,
			},
			{
				Name:  "Pumpkin",
				Stack: 99,
			},
			{
				Name:  "Hoe",
				Stack: 99,
			},
			{
				Name:  "Pickaxe",
				Stack: 99,
			},
			{
				Name:  "Axe",
				Stack: 99,
			},
			{
				Name:  "Watering Can",
				Stack: 99,
			},
			{
				Name:  "Scythe",
				Stack: 99,
			},
			{
				Name:  "aBlank",
				Stack: 99,
			},
			{
				Name:  "bBlank",
				Stack: 99,
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

	currentInventory.NewCache()

	// msg := "Test"
	// numberMsg := inventory.TextToNumber(msg)
	// fmt.Printf("%s : %d\n", msg, numberMsg)
	// fmt.Println(len(currentInventory.Items))
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

	// fmt.Println(inventory.GetDuplicateWeight(4, 1, 4)) // 1
	// fmt.Println(inventory.GetDuplicateWeight(4, 2, 4)) // 3/2
	// fmt.Println(inventory.GetDuplicateWeight(4, 3, 4)) // 1/2
	// fmt.Println(inventory.GetDuplicateWeight(4, 4, 4)) // 1/24

	// fmt.Println(currentInventory.Items[0].WaysForStackCount(12, currentInventory.Cache))
	// fmt.Println(currentInventory.Cache)

	// fmt.Println(inventory.SumFactorization(4))

	permCounts := currentInventory.CalculateSpreadPermutationCounts()
	total := new(big.Rat)
	for _, permCount := range permCounts {
		total.Add(total, permCount)
	}
	totalApproxBits := total.Num().BitLen() - total.Denom().BitLen()
	fmt.Println(total, totalApproxBits)

	// ========== COMBINING STACKS TEST ==============
	// items := []inventory.Item{
	// 	{
	// 		Name:  "Apples",
	// 		Stack: 3,
	// 	},
	// 	{
	// 		Name:  "Oranges",
	// 		Stack: 4,
	// 	},
	// }

	// combinedSpreads := map[uint8]*big.Rat{
	// 	0: big.NewRat(1, 1),
	// }

	// for _, item := range items {
	// 	combinedSpreads = inventory.CombineSpreads(combinedSpreads, item.WaysForStackCount(5, currentInventory.Cache))
	// }

	// fmt.Println(combinedSpreads)

	// inventory.StartEncoder()
	// inventory.StartDecoder()
}

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

func totalInventoryBits(current inventory.Inventory) int {
	permCounts := current.CalculateSpreadPermutationCounts()
	total := new(big.Rat)
	for _, permCount := range permCounts {
		total.Add(total, permCount)
	}
	if total.Sign() <= 0 {
		return 0
	}
	if total.IsInt() {
		count := new(big.Int).Set(total.Num())
		if count.Cmp(big.NewInt(1)) <= 0 {
			return 0
		}
		count.Sub(count, big.NewInt(1))
		return count.BitLen()
	}

	f, _ := total.Float64()
	if f <= 1 {
		return 0
	}
	return int(math.Ceil(math.Log2(f)))
}

func fillInventoryAndPrintBits(current *inventory.Inventory) {
	step := 0
	for {
		fmt.Printf("Step %02d: %d bits\n", step, totalInventoryBits(*current))

		nextBlank := -1
		for index, item := range current.Items {
			if item.Stack == 0 {
				nextBlank = index
				break
			}
		}
		if nextBlank < 0 {
			return
		}

		step++
		current.Items[nextBlank] = inventory.Item{
			Name:  fmt.Sprintf("Unique Item %02d", step),
			Stack: 999,
		}
	}
}

func main() {
	currentInventory := inventory.Inventory{
		Items: []inventory.Item{
			{
				Name:  "Potato Seeds",
				Stack: 99,
			},
			{
				Name:  "Bean Starter",
				Stack: 99,
			},
			{
				Name:  "Sugar",
				Stack: 99,
			},
			{
				Name:  "Rice",
				Stack: 99,
			},
			{
				Name:  "Vinegar",
				Stack: 99,
			},
			{
				Name:  "Kale Seeds",
				Stack: 99,
			},
			{
				Name:  "Parsnip Seeds",
				Stack: 99,
			},
			{
				Name:  "Tulip Bulb",
				Stack: 99,
			},
			{
				Name:  "Jazz Seeds",
				Stack: 99,
			},
			{
				Name:  "Cauliflower Seeds",
				Stack: 99,
			},
			{
				Name:  "Oil",
				Stack: 99,
			},
			{
				Name:  "Wood",
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

	fmt.Println("Starting...")
	// fillInventoryAndPrintBits(&currentInventory)
	// fmt.Println(inventory.TextToNumber("Lily"))
	// message := "Lily"
	message := "a9951ce549e42b03146871b65e314929"
	numberMessage := inventory.TextToNumber(message)
	variation := currentInventory.GetVariation(numberMessage)
	fmt.Println("Message:", numberMessage, message)
	fmt.Println(variation)
	// fmt.Println(currentInventory.CalculateSpreadPermutationCounts())
	// variation := currentInventory.GetVariation(big.NewRat(204223865, 1))
	// variation.Print()
	// variation.NewCache()
	// fmt.Println("Generating Index")
	// fmt.Println(variation.GetIndex())

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

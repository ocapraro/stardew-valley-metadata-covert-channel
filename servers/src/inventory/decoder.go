package inventory

import (
	"encoding/xml"
	"fmt"
	"math/big"
	"os"
	"slices"
	"strings"
	"time"
)

// NumberToText changes a binary into its string equivilent
func NumberToText(number *big.Rat) string {
	if number == nil || number.Sign() == 0 {
		return ""
	}

	integerPart := new(big.Int).Quo(number.Num(), number.Denom())
	binaryString := integerPart.Text(2)
	if rem := len(binaryString) % 8; rem != 0 {
		binaryString = strings.Repeat("0", 8-rem) + binaryString
	}

	var message strings.Builder
	for i := 0; i < len(binaryString); i += 8 {
		charBits := binaryString[i : i+8]
		charValue := uint8(0)
		for _, bit := range charBits {
			charValue <<= 1
			if bit == '1' {
				charValue |= 1
			}
		}
		message.WriteByte(charValue)
	}
	return message.String()
}

type xmlItem struct {
	Name  *string `xml:"name"`
	Stack *uint16 `xml:"stack"`
}

type xmlRoot struct {
	Items []xmlItem `xml:"items>Item"`
}

func GetCurrentInventory(savePath string, c cache) (Inventory, error) {
	data, err := os.ReadFile(savePath)
	if err != nil {
		return Inventory{}, err
	}

	var root xmlRoot
	err = xml.Unmarshal(data, &root)
	if err != nil {
		return Inventory{}, err
	}

	inventory := Inventory{}
	inventory.Cache = c
	for _, item := range root.Items {
		if item.Name == nil || item.Stack == nil {
			inventory.Items = append(inventory.Items, Item{
				Name:  "Blank",
				Stack: 0,
			})
		} else {
			inventory.Items = append(inventory.Items, Item{
				Name:  *item.Name,
				Stack: *item.Stack,
			})
		}
	}

	return inventory, nil
}

func StartDecoder() {
	println("Decoder Started!")
	path, ok := os.LookupEnv("SAVE_PATH")
	if !ok || strings.TrimSpace(path) == "" {
		panic("SAVE_PATH environment variable is not set")
	}
	inventory := Inventory{}
	inventory.NewCache()
	cache := inventory.Cache
	for {
		time.Sleep(time.Second)
		newInventory, _ := GetCurrentInventory(path, cache)
		if slices.Equal(newInventory.Items, inventory.Items) {
			continue
		}
		inventory = newInventory
		inventory.Print()
		msgNumber := inventory.GetIndex()
		msg := NumberToText(msgNumber)
		fmt.Println("Message:", msgNumber, msg)
	}
}

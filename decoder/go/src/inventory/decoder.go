package inventory

import (
	"encoding/xml"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"
)

// NumberToText changes a binary into its string equivilent
func NumberToText(number uint64) string {
	binaryString := fmt.Sprintf("%064b", number)
	var message strings.Builder
	for i := 0; i < len(binaryString); i += 8 {
		charCode, _ := strconv.ParseUint(binaryString[i:i+8], 2, 8)
		message.WriteRune(rune(charCode))
	}
	return message.String()
}

type xmlItem struct {
	Name  *string `xml:"name"`
	Stack *uint8  `xml:"stack"`
}

type xmlRoot struct {
	Items []xmlItem `xml:"items>Item"`
}

func GetCurrentInventory(savePath string) (Inventory, error) {
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
	for _, item := range root.Items {
		if item.Name == nil || item.Stack == nil {
			inventory.Items = append(inventory.Items, Item{
				Name:  "Blank",
				Stack: 1,
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
	const path = "/Users/ocapraro/.config/StardewValley/Saves/CHANNEL_431325361/SaveGameInfo"
	inventory, _ := GetCurrentInventory(path)
	for {
		time.Sleep(time.Second)
		newInventory, _ := GetCurrentInventory(path)
		if slices.Equal(newInventory.Items, inventory.Items) {
			continue
		}
		inventory = newInventory
		inventory.Print()
	}
}

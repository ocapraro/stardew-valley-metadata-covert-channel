package inventory

import (
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"net/http"
	"strings"
)

// TextToNumber changes a string into its binary equivilent
func TextToNumber(message string) *big.Rat {
	var binaryString strings.Builder
	for _, char := range message {
		fmt.Fprintf(&binaryString, "%08b", char)
	}
	binaryMessage := new(big.Int)
	if binaryString.Len() == 0 {
		return new(big.Rat)
	}
	binaryMessage.SetString(binaryString.String(), 2)
	return new(big.Rat).SetInt(binaryMessage)
}

func totalInventoryBits(current Inventory) int {
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

type Message struct {
	Message string `json:"message"`
}

func StartEncoder() {
	currentInventory := Inventory{
		Items: []Item{
			{
				Name:  "Axe",
				Stack: 1,
			},
			{
				Name:  "Hoe",
				Stack: 1,
			},
			{
				Name:  "Watering Can",
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
				Name:  "Scythe",
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
	currentInventory.NewCache()

	message := "Lily"

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Message{Message: "ok"})
	})

	mux.HandleFunc("GET /encode", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		numberMessage := TextToNumber(message)
		variation := currentInventory.GetVariation(numberMessage)
		fmt.Println("Message:", numberMessage, message)
		json.NewEncoder(w).Encode(variation)
	})

	mux.HandleFunc("POST /setInventory", func(w http.ResponseWriter, r *http.Request) {
		inventory := Inventory{}
		inventory.Cache = currentInventory.Cache
		json.NewDecoder(r.Body).Decode(&inventory)
		currentInventory = inventory
		count := totalInventoryBits(currentInventory)
		currentInventory.Print()
		fmt.Printf("New inventory combinations: %d bits\n", count)
		json.NewEncoder(w).Encode(Message{Message: "ok"})
	})

	mux.HandleFunc("POST /setMessage", func(w http.ResponseWriter, r *http.Request) {
		newMessage := Message{}
		json.NewDecoder(r.Body).Decode(&newMessage)
		message = newMessage.Message
		json.NewEncoder(w).Encode(Message{Message: fmt.Sprintf("Message set to: '%s'", message)})
	})

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	println("Encoder Started!")
	print(server.ListenAndServe())
}

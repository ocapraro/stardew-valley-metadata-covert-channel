package inventory

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// TextToNumber changes a string into its binary equivilent
func TextToNumber(message string) uint64 {
	var binaryString strings.Builder
	for _, char := range message {
		fmt.Fprintf(&binaryString, "%08b", char)
	}
	binaryMessage, _ := strconv.ParseUint(binaryString.String(), 2, 64)
	return binaryMessage
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

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Message{Message: "ok"})
	})

	mux.HandleFunc("GET /encode/{message}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		message := r.PathValue("message")
		numberMessage := TextToNumber(message)
		variation := currentInventory.GetVariation(numberMessage)
		json.NewEncoder(w).Encode(variation)
	})

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	println("Encoder Started!")
	print(server.ListenAndServe())
}

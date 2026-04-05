package inventory

import (
	// "encoding/json"
	"fmt"
	// "net/http"
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
	// currentInventory := Inventory{
	// 	Items: []Item{
	// 		{
	// 			Name:  "Axe",
	// 			Stack: 1,
	// 		},
	// 		{
	// 			Name:  "Hoe",
	// 			Stack: 1,
	// 		},
	// 		{
	// 			Name:  "Watering Can",
	// 			Stack: 1,
	// 		},
	// 		{
	// 			Name:  "Pickaxe",
	// 			Stack: 1,
	// 		},
	// 		{
	// 			Name:  "Parsnip Seeds",
	// 			Stack: 10,
	// 		},
	// 		{
	// 			Name:  "Scythe",
	// 			Stack: 1,
	// 		},
	// 		{
	// 			Name:  "Parsnip Seeds",
	// 			Stack: 5,
	// 		},
	// 		{
	// 			Name:  "Blank",
	// 			Stack: 0,
	// 		},
	// 		{
	// 			Name:  "Blank",
	// 			Stack: 0,
	// 		},
	// 		{
	// 			Name:  "Blank",
	// 			Stack: 0,
	// 		},
	// 		{
	// 			Name:  "Blank",
	// 			Stack: 0,
	// 		},
	// 		{
	// 			Name:  "Blank",
	// 			Stack: 0,
	// 		},
	// 	},
	// }

	// message := "Lily"

	// mux := http.NewServeMux()
	// mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
	// 	w.Header().Set("Content-Type", "application/json")
	// 	json.NewEncoder(w).Encode(Message{Message: "ok"})
	// })

	// mux.HandleFunc("GET /encode", func(w http.ResponseWriter, r *http.Request) {
	// 	w.Header().Set("Content-Type", "application/json")
	// 	numberMessage := TextToNumber(message)
	// 	variation := currentInventory.GetVariationOLD(numberMessage)
	// 	json.NewEncoder(w).Encode(variation)
	// })

	// mux.HandleFunc("POST /setInventory", func(w http.ResponseWriter, r *http.Request) {
	// 	inventory := Inventory{}
	// 	json.NewDecoder(r.Body).Decode(&inventory)
	// 	currentInventory = inventory
	// 	count := 0
	// 	for _, permCount := range inventory.getCounts() {
	// 		count += int(permCount)
	// 	}
	// 	fmt.Printf("New inventory combinations: %d\n", count)
	// 	json.NewEncoder(w).Encode(Message{Message: "ok"})
	// })

	// mux.HandleFunc("POST /setMessage", func(w http.ResponseWriter, r *http.Request) {
	// 	newMessage := Message{}
	// 	json.NewDecoder(r.Body).Decode(&newMessage)
	// 	message = newMessage.Message
	// 	json.NewEncoder(w).Encode(Message{Message: fmt.Sprintf("Message set to: '%s'", message)})
	// })

	// server := &http.Server{
	// 	Addr:    ":8080",
	// 	Handler: mux,
	// }
	println("Encoder Started!")
	// print(server.ListenAndServe())
}

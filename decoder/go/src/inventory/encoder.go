package inventory

import (
	"fmt"
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

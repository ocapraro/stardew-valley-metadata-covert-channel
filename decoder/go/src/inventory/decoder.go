package inventory

import (
	"fmt"
	"strconv"
	"strings"
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

package challenges

import (
	"fmt"
	"strings"
)

func StringValueHashing(value string) {

	lengthOfVal := len(value)

	fmt.Printf("\nLength := %d, value := %s\n", lengthOfVal, value)
}

func MaskTheNames(fullName string) {

	names := strings.Split(fullName, " ")
	var maskedFirstName, maskedLastName string

	// Mask the first name
	for i, char := range names[0] {
		if i == 0 {
			maskedFirstName += string(char) //keep the first letter
		} else {
			maskedFirstName += "*" //mask the rest of the letters
		}
	}

	// Mask the last name
	for i, char := range names[1] {
		if i == 0 {
			maskedLastName += string(char) //keep the first letter
		} else {
			maskedLastName += "*" //mask the rest of the letters
		}
	}

	fmt.Printf("maskedFirstName : %s\nmaskedLastName : %s\n", maskedFirstName, maskedLastName)
}

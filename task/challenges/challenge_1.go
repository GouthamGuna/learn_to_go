package challenges

import (
	"fmt"
	"regexp"
	"strconv"
)

func StringArrToInt() {

	pro := []string{"get fee list of 456", "fee balance of 80", "pending fee of 74", "total fee of 70", "10"}

	var i []int
	re := regexp.MustCompile(`\d+`)

	for _, str := range pro {
		match := re.FindString(str)

		if match != "" {
			num, err := strconv.Atoi(match)

			if err == nil {
				i = append(i, num)
			}
		}
	}

	fmt.Printf("\n%d", i)
}

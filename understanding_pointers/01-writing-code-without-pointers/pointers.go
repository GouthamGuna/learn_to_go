package main

import "fmt"

func main() {
	age := 18 // regular variable

	fmt.Printf("The age is: %d\n", age)
	fmt.Println(validateAge(age))
}

func validateAge(age int) string {
	if age < 18 {
		return "You are not old enough to drive."
	}
	return "You are old enough to drive."
}

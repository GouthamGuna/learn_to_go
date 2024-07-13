package main

import "fmt"

/*
=> `&` References the use pointer the variable and get the memoey address.
=> `*` DeReference to pointer address and get the value.
*/
func main() {
	age := 18 // regular variable

	var agePointer *int = &age

	// agePointer = &age // should merge variable declaration with assignment

	fmt.Printf("The age is: %d\n", *agePointer)
	fmt.Println(validateAge(age))
}

func validateAge(age int) string {
	if age < 18 {
		return "You are not old enough to drive."
	}
	return "You are old enough to drive."
}

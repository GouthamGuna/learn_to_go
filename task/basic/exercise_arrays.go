package basic

import "fmt"

// Arrays

var cars = [4]string{"Volvo", "BMW", "Ford", "Audi"}

func ExerciseOne() {
	fmt.Printf("\nprint array of string ::: %s\n", cars)
}

func ExerciseTwo() {
	fmt.Printf("Print the value of the second element in the cars array ::: %s\n", cars[1])
}

func ExerciseThree() {
	cars[0] = "Maruti Suzuki"
	fmt.Printf("Change the value from Volvo to Maruti Suzuki, in the cars array ::: %s\n", cars)
}

// Operators

func ExerciseFour() {
	fmt.Printf("Multiply 10 with 5, and print the result ::: %d\n", 10*5)
}

func ExerciseFive() {
	x := 1
	x++
	fmt.Printf("arithmetic operator to increase the value of the variable x by 1 ::: %d\n", x)
}

func ExerciseSix() {
	x := 1
	x += 5
	fmt.Printf("Use assignment operators to add 5 to x ::: %d\n", x)
}

// Conditions

func ExerciseSeven() {
	x := 50
	y := 25

	if x > y {
		fmt.Printf("Print Hello World if x is greater than y ::: %s\n", "Hello World")
	}
}

func ExerciseEight() {
	x := 50
	y := 50

	if x == y {
		fmt.Printf("Print Hello World if x is equal to y ::: %s\n", "Hello World")
	}
}

func ExerciseNine() {
	var x = 50
	var y = 500

	if x == y {
		fmt.Printf("Print Yes if x is equal to y, otherwise print No ::: %s\n", "Yes")
	} else {
		fmt.Printf("Print Yes if x is equal to y, otherwise print No ::: %s\n", "No")
	}
}

// Print "1" if x is equal to y, print "2" if x is greater than y, otherwise print "3".

func ExerciseTen() {
	var x = 50
	var y = 50

	if x == y {
		fmt.Print("If Loop ::: 1\n")
	} else if x > y {
		fmt.Print("Else If Loop ::: 2\n")
	} else {
		fmt.Print("Else Loop ::: 3\n")
	}
}

// Switch

func ExerciseEleven() {
	var day = 2
	switch day {
	case (1):
		fmt.Print("Saturday ::: from Switch case\n")
	case (2):
		fmt.Print("Sunday ::: from Switch case\n")
	}
}

func ExerciseTwelve() {
	var day = 4

	switch day {

	case (1):
		fmt.Print("From Switch case ::: Saturday\n")

	case (2):
		fmt.Print("From Switch case ::: Sunday\n")
	default:
		fmt.Print("From Switch case ::: Weekday\n")
	}
}

// Loops

func ExerciseThirteen() {

	for i := 0; i < 6; i++ {
		fmt.Printf("From Thirteen for loop ::: %d\n", i)
	}
}

func ExerciseFourteen() {

	for i := 0; i < 5; i++ {
		fmt.Println("---------------------------------")
		fmt.Printf("Use a for loop to print Yes 5 times ::: %d\n", i)
	}
}

func ExerciseFifteen() {

	for i := 0; i < 10; i++ {
		if i == 5 {
			fmt.Printf("Stop the loop if for is 5 ::: %d\n", i)
			break
		}
	}
}

func ExerciseSixteen() {

	for i := 0; i < 10; i++ {
		if i == 4 {
			fmt.Println("In the following loop, when the value is 4, jump directly to the next value.")
			continue
		}
	}
}

// Functions

func myFunction(fname string) {
	fmt.Printf("Author ::: %v Gunasekaran\n", fname)
}

func myFunctionReturnInt(x int) int {
	return 6 + x
}

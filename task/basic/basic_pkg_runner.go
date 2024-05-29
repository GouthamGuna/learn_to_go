package basic

import "fmt"

func Basic_PKG_Runner() {

	// external pkg
	QuotePrint()

	// Arrays
	ExerciseOne()
	ExerciseTwo()
	ExerciseThree()

	// Operators
	ExerciseFour()
	ExerciseFive()
	ExerciseSix()

	// Conditions
	ExerciseSeven()
	ExerciseEight()
	ExerciseNine()
	ExerciseTen()

	// Switch
	ExerciseEleven()
	ExerciseTwelve()

	// Loops
	ExerciseThirteen()
	ExerciseFourteen()
	ExerciseFifteen()
	ExerciseSixteen()

	// Functions
	fmt.Printf("Func return stmt ::: %d\n", myFunctionReturnInt(18))
	myFunction("Gowtham Sankar")

	fmt.Println("-----------------------------------------------------")
	fmt.Println(" Go - Variables Starting Here !")
	fmt.Println("-----------------------------------------------------")

	// Variables
	VariablesExample_1()
	VariablesExample_2()
	VariablesExample_3()
	VariablesExample_4()
	VariablesExample_5()
	VariablesExample_6()
	VariablesExample_7()
}

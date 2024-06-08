package basic

import "fmt"

/*
Variables are containers for storing data values.

var variablename type = value

Note: You always have to specify either type or value (or both).

2. With the := sign:
Use the := sign, followed by the variable value:

Syntax
variablename := value

*/

func VariablesExample_1() {
	var student1 string = "Jay Kumar" //type is string
	var student2 = "Manoj"            //type is inferred
	x := 2                            //type is inferred

	fmt.Println(student1)
	fmt.Println(student2)
	fmt.Println(x)
}

// Variable Declaration With Initial Value
func VariablesExample_2() {
	var mobile string = "android" // type is string
	var brand = "Xiaomi"          // type is inferred
	makeYear := 2024              // type is inferred

	fmt.Printf("Mobile type is %s made by %s year of release %d\n", mobile, brand, makeYear)
}

func VariablesExample_3() {
	var a string
	var b int
	var c bool

	fmt.Println(a)
	fmt.Println(b)
	fmt.Println(c)
}

// Value Assignment After Declaration
func VariablesExample_4() {
	var printer string
	printer = "Canon"

	// printer := "testing" // error ::: no new variables on left side of :=

	fmt.Printf("\nMy Printer is %s.\n", printer)
}

/**
 a := 1

Since := is used outside of a function, running the program results in an error.

syntax error: non-declaration statement outside function body

*/

// Go Multiple Variable Declaration

func VariablesExample_5() {

	var a, b, c, d int = 7, 77, 777, 7777

	fmt.Println(a)
	fmt.Println(b)
	fmt.Println(c)
	fmt.Println(d)

	// Note: If you use the type keyword, it is only possible to declare one type of variable per line.
}

func VariablesExample_6() {

	// If the type keyword is not specified, you can declare different types of variables in the same line:

	var a, b = 24, "Hello world"
	c, d := 27, "Hello Lunar"

	fmt.Println(a)
	fmt.Println(b)
	fmt.Println(c)
	fmt.Println(d)
}

// Go Variable Declaration in a Block
func VariablesExample_7() {

	fmt.Println("\nGo Variable Declaration in a Block")

	var (
		a int
		b = "Hello world"
		c = 2027
		d = "Hello Lunar"
	)

	a = 2024

	fmt.Println(a)
	fmt.Println(b)
	fmt.Println(c)
	fmt.Println(d)
}

// Go Variable Naming Rules
// -> The variable name cannot be any Go keywords

/**
Camel Case -> myVariableName = "John"

Pascal Case -> MyVariableName = "John"

Snake Case -> my_variable_name = "John"
*/

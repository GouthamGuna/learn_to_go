package basic

import (
	"fmt"
)

/*
	Go Constants

	If a variable should have a fixed value, that cannot be changed, you can use the `const` keyword

	The `const` keyword declares the variable as a constant it`s means a constant

	Syntax : const CONST_NAME = value

*/

const PI = 3.14

func ConstantsExample_1() {
	fmt.Println(PI)
}

const A int = 7

func ConstantsExample_2() {
	fmt.Printf("Type constants are declared with a defined type : %d\n", A)
}

const B = 77

func ConstantsExample_3() {
    fmt.Printf("Untyped constants are declared without a defined type : %d\n", B)
}

/*

    func ConstantsExample_4() {
		const A = 777

		A = 27

		fmt.Println(A) // error : cannot assign to A
	}
*/

// Multiple constants are declared

func ConstantsExample_5() {
    const (
        A = 777
        B = 888
        C = "Hi there!"
    )

    fmt.Println(A)
    fmt.Println(B)
    fmt.Println(C)
}
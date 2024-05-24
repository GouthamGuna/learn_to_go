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

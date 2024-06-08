package basic

import "fmt"

/**
General Formatting Verbs

Verb	Description
%v	Prints the value in the default format
%#v	Prints the value in Go-syntax format
%T	Prints the type of the value
%%	Prints the % sign
*/

func Formatting_Example_1() {

	var i = 15.5
	var txt = "Hello World!"

	fmt.Printf("%v Prints the value in the default format\n", i)
	fmt.Printf("%#v Prints the value in Go-syntax format\n", i)
	fmt.Printf("%T Prints the type of the value\n", i)
	fmt.Printf("%v%%\n", i)

	fmt.Println("-------------------")

	fmt.Printf("%v Prints the value in the default format\n", txt)
	fmt.Printf("%#v Prints the value in Go-syntax format\n", txt)
	fmt.Printf("%T Prints the type of the value\n", txt)
	fmt.Printf("%X Base 16, uppercase \n", txt)
}

/**
Integer Formatting Verbs

The following verbs can be used with the integer data type:

Verb	Description
%b	Base 2
%d	Base 10
%+d	Base 10 and always show sign
%o	Base 8
%O	Base 8, with leading 0o
%x	Base 16, lowercase
%X	Base 16, uppercase
%#x	Base 16, with leading 0x
%4d	Pad with spaces (width 4, right justified)
%-4d	Pad with spaces (width 4, left justified)
%04d	Pad with zeroes (width 4

*/

func Formatting_Example_2() {
	var i = 15

	fmt.Printf("%b\n", i)
	fmt.Printf("%d\n", i)
	fmt.Printf("%+d\n", i)
	fmt.Printf("%o\n", i)
	fmt.Printf("%O\n", i)
	fmt.Printf("%x\n", i)
	fmt.Printf("%X\n", i)
	fmt.Printf("%#x\n", i)
	fmt.Printf("%4d\n", i)
	fmt.Printf("%-4d\n", i)
	fmt.Printf("%04d\n", i)
}

/**
String Formatting Verbs

Verb	Description
%s	Prints the value as plain string
%q	Prints the value as a double-quoted string
%8s	Prints the value as plain string (width 8, right justified)
%-8s	Prints the value as plain string (width 8, left justified)
%x	Prints the value as hex dump of byte values
% x	Prints the value as hex dump with spaces
*/

func Formatting_Example_3() {
	var txt = "Hello"

	fmt.Printf("%s\n", txt)
	fmt.Printf("%q\n", txt)
	fmt.Printf("%8s\n", txt)
	fmt.Printf("%-8s\n", txt)
	fmt.Printf("%x\n", txt)
	fmt.Printf("% x\n", txt)
}

/**
Boolean Formatting Verbs

Verb	Description
%t	Value of the boolean operator in true or false format (same as using %v)
*/

func Formatting_Example_4() {
	var i = true
	var j = false

	fmt.Printf("%t\n", i)
	fmt.Printf("%t\n", j)
}

/**
Float Formatting Verbs

Verb	Description
%e	Scientific notation with 'e' as exponent
%f	Decimal point, no exponent
%.2f	Default width, precision 2
%6.2f	Width 6, precision 2
%g	Exponent as needed, only necessary digits
*/

func Formatting_Example_5() {
	var i = 3.141

	fmt.Printf("%e\n", i)
	fmt.Printf("%f\n", i)
	fmt.Printf("%.2f\n", i)
	fmt.Printf("%6.2f\n", i)
	fmt.Printf("%g\n", i)
}

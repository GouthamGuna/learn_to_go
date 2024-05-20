# Go Data Type`s

Data types in Go refer to the way you can define and use variables. Go is a statically typed language, which means you must declare the data type of a variable when you declare it.
Here are the basic data types in Go:

**Numeric Types**

 * `Integers`: int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64
 * int is the default integer type and is 32 bits long on 32-bit systems and 64 bits long on 64-bit systems.
 * The other integer types are similar to int but are 8, 16, or 32 bits long.

 * `Floating-point numbers`: float32, float64
 * float32 is a 32-bit floating-point number.
 * float64 is a 64-bit floating-point number, which is the default type for floating-point numbers.

 * `Complex numbers`: complex64, complex128
 * complex64 is a complex number with a 32-bit real part and a 32-bit imaginary part.
 * complex128 is a complex number with a 64-bit real part and a 64-bit imaginary part.

**String Type**

 * String: `string`
 * Strings are sequences of characters. You can use the string type to declare a variable that holds a string.

**Boolean Type**
 
 * Boolean:`bool`
 * The bool type is used to represent boolean values, which can be either true or false.

**Other Types**

 * Array: `array`
 * An array is a collection of values of the same type stored in a single variable.

 * Slice: `slice`
 * A slice is a dynamic array, which means it can grow or shrink in size.

 * Struct: `struct`
 * A struct is a collection of fields, which are individual variables of the same type.

 * Pointer: `pointer`
 * A pointer is a variable that holds the memory address of another variable.

 Here is an example of how to use some of these data types in Go:

```
package main

import "fmt"

func main() {
    // Integer
    var x int = 5
    fmt.Println(x)

    // Floating-point number
    var y float64 = 3.14
    fmt.Println(y)

    // String
    var name string = "John"
    fmt.Println(name)

    // Boolean
    var isAdmin bool = true
    fmt.Println(isAdmin)

    // Array
    var scores [5]int = [5]int{1, 2, 3, 4, 5}
    fmt.Println(scores)

    // Slice
    var ages []int = []int{12, 15, 18, 20}
    fmt.Println(ages)

    // Struct
    type person struct {
        name  string
        age   int
        email string
    }
    var p person = person{"John", 30, "john@example.com"}
    fmt.Println(p)

    // Pointer
    var p1 *int = &x
    fmt.Println(*p1)
}
```

This program demonstrates the use of various data types in Go, including integers, floating-point numbers, strings, booleans, arrays, slices, structs, and pointers.
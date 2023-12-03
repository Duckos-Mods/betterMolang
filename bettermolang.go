package main

import (
	bm "github.com/duckos-mods/betterMolang/bettermolang"
)

const printMessage string = `
Welcome to better molang compiler!
This program compiles a custom molang extension into a MCBE valid molang expression.
I hope you enjoy it!

Usage:
	- F = file path to the file you want to compile
	- O = file path to the file you want to output to
	- O = is optional and will run the optimizer if present

Featuring:
	- If statements (how the hell is this not a thing yet)
	- If else statements (how the hell is this not a thing yet)
	- Arrays (again how the hell is this not a thing yet)
	- Functions (hopefully)
	- Bypass the 1024 loop limit using nested loops

Basic Explanation of arrays:
	- Arrays are indexed with the [] operator (revolitionary I know)
	- Arrays are declared with the array keyword
	example:
		myArray = (3)<1, 2, 3>
	- Arrays are 0 indexed
	- Arrays size is fixed at declaration and cannot be changed (yet)
`

func main() {
	println(printMessage)
	bm.RunTest()
}

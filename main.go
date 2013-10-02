package main

import "fmt"

func main() {
	//Interpreter()

	toks := Tokenize("y + 6 = 5x - (4 + 3 * 2)")
	for _,v := range toks {
		fmt.Println(v)
	}
}

package main

import "fmt"

func main() {
	eq,_ := Parse("(119 + 21i) - j(15j / 3)")
	fmt.Println(eq.Print())
	fmt.Println(Vars)
	Vars['i'].val = 7
	Vars['j'].val = 13
	fmt.Println(eq.Value())
}


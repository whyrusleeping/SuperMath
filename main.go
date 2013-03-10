package main

import "fmt"

func main() {
	//eq,_ := Parse("(119 + 21i) - j(15j / 3)")
	eq2,_ := Parse("15")
	fmt.Println(eq2.Print())

	eq,_ := Parse("(12 + i) * (2 - j)")
	fmt.Println(eq.Print())
	Vars['i'].val = 2
	Vars['j'].val = 10
	le := Equality{eq,eq2}
	fmt.Println(le.SolveFor('j'))
}


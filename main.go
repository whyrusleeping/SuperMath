package main

import "fmt"

func main() {
	//eq,_ := ParseExpression("(119 + 21i) - j(15j / 3)")
	eq2,_ := ParseExpression("15 - 3(j ^ 3)")
	fmt.Println(eq2.Print())

	eq,_ := ParseExpression("(12 + i) * (2 - j)")
	fmt.Println(eq.Print())
	Vars['i'].val = 2
	Vars['j'].val = 10
	fmt.Println(eq2.Value())
	le := Equality{eq,eq2}
	fmt.Println(le.SolveFor('j'))
}


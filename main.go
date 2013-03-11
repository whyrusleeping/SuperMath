package main

import "fmt"

func main() {
	eq,_ := ParseEquation("15 * j = 2 + (j ^ 2) * 6")
	fmt.Printf("j = %f\n", eq.SolveFor('j'))
}


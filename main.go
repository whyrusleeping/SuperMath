package main

import "fmt"

func main() {
	//Interpreter()
	s := "atan(4+x)=y-3*5"
	toks := Tokenize(s)
	for _,v := range toks {
		fmt.Println(v)
	}

	eq,err := ParseEquation(s)

	if err != nil {
		panic(err)
	}
	fmt.Println(eq.Print())
}

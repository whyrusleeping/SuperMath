package main

import "testing"

func TestEquation(t *testing.T) {
	eq,err := ParseEquation("y ^ 3 = 4x + (8 - (3x^2 -6))")
	if err != nil {
		panic(err)
	}

	Vars['x'].val = 3

	//y = -1

	res,off := eq.SolveFor('y')
	if res != -1 || off != 0{
		t.Fatalf("solved to %f",res)
	}
}

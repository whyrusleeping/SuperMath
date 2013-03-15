package main

import (
	"testing"
	"math"
)

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

func TestFunction(t* testing.T) {
	eq, err := ParseEquation("tan(x + 4)=y -7")
	if err != nil {
		panic(err)
	}
	Vars['x'].val = 1.9
	actual := math.Tan(5.9) + 7
	solve,_ := eq.SolveFor('y')
	if actual != solve {
		t.Fatalf("%f %f",actual,solve)
	}
}

func BenchmarkAsyncEquate(b *testing.B) {
	eq,_ := ParseEquation("y^(x - 2) - x^4 * y= x^2.5 - 6(y - x^0.5)")
	for i := 0; i < b.N; i++ {
		eq.Difference()
	}
}

func BenchmarkEquate(b *testing.B) {
	eq,_ := ParseEquation("y^(x - 2) - x^4 * y= x^2.5 - 6(y - x^0.5)")
	r := 0.0
	for i := 0; i < b.N; i++ {
		r = eq.left.Value() - eq.right.Value()
	}
	b.Log(r)
}

package main

import (
	"fmt"
	"math"
)

type CalcTerm struct {
	factor float64
	vr *Variable
	power Equatable
}

func (ct *CalcTerm) Value() float64 {
	return ct.factor * math.Pow(ct.vr.Value(),ct.power.Value())
}

func (ct *CalcTerm) Print() string {
	return fmt.Sprintf("%f%s^%s",ct.factor, ct.vr.Print(), ct.power.Print())
}

func (ct *CalcTerm) simple() bool {
	return false
}

func (ct *CalcTerm) Integrate() *CalcTerm {
	nct := CalcTerm{}
	nct.factor = ct.factor / ct.power.Value()
	fmt.Println("Caution, integrate doesnt quite work yet")
	nct.power = ct.power //+ 1
	return &nct
}

func (ct *CalcTerm) ContainsVar(v string) bool {
	return ct.vr.C == v || ct.power.ContainsVar(v)
}

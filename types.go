package main

import (
	"math"
	"strconv"
	"fmt"
)

//Equation Forms
const (
	FSimple = iota
	FCalculus
)
const (
	OAdd = iota
	OSub
	OMul
	ODiv
	OPow
	OLog
)

type Equatable interface {
	Value() float64
	Print() string
}

type Constant struct {
	val float64
}

func NewConstant(val string) *Constant {
	n, _ := strconv.Atoi(val)
	return &(Constant{float64(n)})
}

func (c *Constant) Value() float64{
	return c.val
}

func (c *Constant) Print() string{
	return fmt.Sprint(c.val)
}

type Equation struct {
	vars []*Variable
}

type Variable struct {
	C uint8
	val float64
}

func NewVariable(C string) *Variable {
	v := Variable{C[0], 0}
	return &v
}

func (v *Variable) Print() string{
	s := make([]byte,1)
	s[0] = v.C
	return string(s)
}

type Term struct {
	left, right Equatable
	operator int
}

type CalcTerm struct {
	factor float64
	vr *Variable
	power float64
}

func (ct *CalcTerm) Value() float64 {
	return ct.factor * math.Pow(ct.vr.Value(),ct.power)
}

func (ct *CalcTerm) Integrate() *CalcTerm {
	nct := CalcTerm{}
	nct.factor = ct.factor / ct.power
	nct.power = ct.power + 1
	return &nct
}

func (t *Term) Value() float64 {
	switch t.operator {
		case OAdd: return t.left.Value() + t.right.Value()
		case OSub: return t.left.Value() - t.right.Value()
		case OMul: return t.left.Value() * t.right.Value()
		case ODiv: return t.left.Value() / t.right.Value()
	}
	return 0
}

func (t *Term) Print() string{
	ops := ""
	switch t.operator {
		case OAdd:
			ops = "+"
		case OSub:
			ops = "-"
		case OMul:
			ops = "*"
		case ODiv:
			ops = "/"
		case OPow:
			ops = "^"
	}
	return fmt.Sprintf("(%s %s %s)",t.left.Print(), ops, t.right.Print())
}

func (v *Variable) Value() float64 {
	return v.val
}


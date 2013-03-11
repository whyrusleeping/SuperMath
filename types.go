package main

import (
	"math"
	"strconv"
	"fmt"
)

//Map for keeping track of variables
var Vars = make(map[uint8]*Variable)

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

type Equality struct {
	left, right Equatable
}

func (e *Equality) SolveFor(v uint8) float64 {
	tolerance := 0.00000000000000001
	vr := Vars[v]
	if math.Abs(vr.val) < 1 {
		vr.val = 5 //5 is sufficiently random, selected by a random dice roll
	}
	difference := e.left.Value() - e.right.Value()
	delta := vr.Value() / 2
	for i:= uint64(0); math.Abs(difference) > tolerance && i < 10e6; i++ {
		vr.val += delta
		ndiff := e.left.Value() - e.right.Value()
		if (ndiff > difference && ndiff > 0) || (ndiff < difference && ndiff < 0) {
			delta = delta * -0.9
		}
		if (ndiff < 0  && difference > 0) || (ndiff > 0 && difference < 0) {
			vr.val -= delta
			delta /= 4
		}
		difference = ndiff
	}
	return vr.val
}

func (e *Equality) Print() string {
	return fmt.Sprintf("%s = %s",e.left.Print(), e.right.Print())
}

type Variable struct {
	C uint8
	val float64
}

func NewVariable(C string) *Variable {
	v,ok := Vars[C[0]]
	if !ok {
		v = &Variable{C[0], 0}
		Vars[C[0]] = v
	}
	return v
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
		case OPow: return math.Pow(t.left.Value(), t.right.Value())
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


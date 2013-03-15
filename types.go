package main

import (
	"math"
	"strconv"
	"errors"
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
	OFun
)

type Equatable interface {
	Value() float64
	Print() string
	simple() bool
}

type Constant struct {
	val float64
}

func NewConstant(val string) *Constant {
	n,_ := strconv.ParseFloat(val, 64)
	return &(Constant{float64(n)})
}

func (c *Constant) Value() float64{
	return c.val
}

func (c *Constant) Print() string{
	return fmt.Sprint(c.val)
}

func (c *Constant) simple() bool {
	return true
}

type Equality struct {
	left, right Equatable
}

func (e *Equality) Difference() float64 {
	return e.left.Value() - e.right.Value()
}

func (e *Equality) SolveFor(v uint8) (float64, float64) {
	tolerance := 0.00000000000000001
	vr := Vars[v]
	if math.Abs(vr.val) < 1 {
		vr.val = 5 //5 is sufficiently random, selected by a random dice roll
	}
	difference := e.Difference()
	h := 0.0000001
	for i:= uint64(0); math.Abs(difference) > tolerance && i < 10e6; i++ {
		difference = e.Difference()
		tmp := vr.val
		vr.val += h
		pos := e.Difference()
		vr.val = tmp - h
		neg := e.Difference()
		vr.val = tmp
		vr.val -= difference / ((pos - neg) / (2 * h))
	}
	return vr.val, difference
}

//Note, this doesnt actually do anything because im not that smart
func (e *Equality) Differentiate(of, to uint8) (*Equality, error) {
	_,okl := e.left.(*Variable)
	_,okr := e.right.(*Variable)
	//First, get a single variable on the left side, if we cant, exit with a failure
	if !okl && !okr {
		return nil, errors.New("Equation must have single variable on one side for now")
	}
	return nil,nil
}

func (e *Equality) Print() string {
	return fmt.Sprintf("%s = %s",e.left.Print(), e.right.Print())
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

func (t *Term) simple() bool {
	return t.left.simple() && t.right.simple()
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

func Simplify(e Equatable) Equatable {
	if e.simple() {
		return &Constant{e.Value()}
	} else {
		t,ok := e.(*Term)
		if ok {
			t.left = Simplify(t.left)
			t.right = Simplify(t.right)
			vl, okl := t.left.(*Variable)
			vr, okr := t.right.(*Variable)
			if okl && okr {
				if vl.C == vr.C {
					switch t.operator {
					case OAdd:
						return &Term{&Constant{2.0}, vl, OMul}
					case OSub:
						return &Constant{0.0}
					case OMul:
						return &Term{vl, &Constant{2.0}, OPow}
					case ODiv:
						return &Constant{1.0}
					}
				}
			} else if okl {
				r, cr := simplifyVars(vl, t.right, t.operator)
				if cr {
					return r
				}
			} else if okr {
				l, cl := simplifyVars(vr, t.left, t.operator)
				if cl {
					return l
				}
			}
			//insert further logic here.
		}
	}
	//Next up: check for cancelling, ie ((6 * x^2) / 2) = (3 * x^2)
	return e
}

//Checks for cases like:
// X * (X ^3) -> X ^ 4
// X + (X - 6) -> ((2 * X) - 6)
// X - (X + 5) -> 5
//Also try to do:
// (1 / x) -> (x ^ -1)
// (3 / x) -> (3 * (x ^ -1))
func simplifyVars(v *Variable, e Equatable, op int) (Equatable, bool) {
	//check if 'e' is a term, otherwise we have nothing to do
	t, ok := e.(*Term)
	if !ok {
		return nil, false
	}

	ttlv, ttokl := t.left.(*Variable)
	//For this case only work where left side of nested statement is the Variable in question
	if ttokl && ttlv.C == v.C {
		if t.operator == OPow {
			if op == OMul {
				if t.right.simple() {
					return &Term{v, &Constant{t.right.Value() + 1}, OPow}, true
				} else {
					returnTerm := &Term{v, nil, OPow}
					returnTerm.right = &Term{t.right, &Constant{1.0}, OAdd}
					return returnTerm
				}
			} else if op == ODiv {
				
			}
		}
	}

	// Now attempt combination
	// General case : X  * (  X     ^     3)
	//               [v][op][subv][t.op][par]
	return nil, false
}

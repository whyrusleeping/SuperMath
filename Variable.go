package main

type Variable struct {
	C string
	val float64
}

func NewVariable(C string) *Variable {
	v,ok := Vars[C]
	if !ok {
		v = &Variable{C, 0}
		Vars[C] = v
	}
	return v
}

func (v *Variable) Print() string {
	return v.C
}

func (v *Variable) simple() bool {
	return false
}

func (v *Variable) Value() float64 {
	return v.val
}

func (v *Variable) ContainsVar(vc string) bool {
	return v.C == vc
}

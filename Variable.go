package main

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

func (v *Variable) Print() string {
	s := make([]byte,1)
	s[0] = v.C
	return string(s)
}

func (v *Variable) simple() bool {
	return false
}

func (v *Variable) Value() float64 {
	return v.val
}

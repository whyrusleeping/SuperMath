package main

import (
	"strings"
	"fmt"
	"math"
)

const (
	FSin = iota
	FCos
	FTan
	FASin
	FACos
	FATan
)

type Function struct {
	ftype int
	arg Equatable
}

func FStrToConst(f string) int {
	f = strings.ToLower(f)
	switch f {
		case "sin":
			return FSin
		case "cos":
			return FCos
		case "tan":
			return FTan
		case "asin":
			return FASin
		case "acos":
			return FACos
		case "atan":
			return FATan
	}
	return -1
}

func ConstToFStr(f int) string {
	switch f {
		case FASin:
			return "asin"
		case FACos:
			return "acos"
		case FATan:
			return "atan"
		case FSin:
			return "sin"
		case FCos:
			return "cos"
		case FTan:
			return "tan"
	}
	return "UNKNOWN"
}

func (f *Function) Value() float64 {
	switch f.ftype {
		case FASin:
			return math.Asin(f.arg.Value())
		case FACos:
			return math.Acos(f.arg.Value())
		case FATan:
			return math.Atan(f.arg.Value())
		case FSin:
			return math.Sin(f.arg.Value())
		case FCos:
			return math.Cos(f.arg.Value())
		case FTan:
			return math.Tan(f.arg.Value())
	}
	return 0.0
}


func (f *Function) Print() string {
	var fmtstr string
	if _,ok := f.arg.(*Term); ok {
		fmtstr = "%s%s"
	} else {
		fmtstr = "%s(%s)"
	}
	return fmt.Sprintf(fmtstr, ConstToFStr(f.ftype), f.arg.Print())
}

func (f *Function) simple() bool {
	return f.arg.simple()
}

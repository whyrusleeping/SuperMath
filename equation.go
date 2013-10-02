package main

import (
	"errors"
	"strings"
	"fmt"
	"strconv"
	"os"
	"bufio"
	"bytes"
)

func IsOperator(c uint8) bool {
	if c == '+' ||
	c == '-' ||
	c == '*' ||
	c == '/' ||
	c == '^' {
		return true
	}
	return false
}

func IsAlpha(c uint8) bool {
	if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') {
		return true
	}
	return false
}

func IsNum(c uint8) bool {
	if c >= '0' && c <= '9' {
		return true
	}
	return false
}

const (
	TLparen = byte(iota)
	TRparen
	TNumber
	TVariable
	TOperator
	TFunction
	TUnknown
)

type Token struct {
	kind byte
	val string
}

func IsFunction(segment string) (bool, string) {
	for i := 0; i < len(segment); i++ {
		if segment[i] == ' ' || IsOperator(segment[i]) {
			return false,""
		}
		if segment[i] == '(' {
			return true, segment[:i]
		}
	}
	return false,""
	/*
	if len(segment) > 3 {
		tr := segment[:3]
		switch tr {
		case "sin", "cos", "tan":
			return true, tr
		}
		tr = segment[:4]
		switch tr {
		case "asin", "acos", "atan":
			return true, tr
		}
	}
	return false, ""
	*/
}

func Tokenize(input string) []*Token {
	tokens := make([]*Token, 0, len(input))
	buf := new(bytes.Buffer)

	for i := 0; i < len(input); i++ {
		c := input[i]
		if IsNum(c) || c == '.' {
			buf.WriteByte(input[i])
		} else {
			if buf.Len() > 0 {
				tokens = append(tokens,&Token{TNumber, buf.String()})
				buf.Reset()
			}
			t := TUnknown
			if c == '(' {
				t = TLparen
			} else if c == ')' {
				t = TRparen
			} else if IsOperator(c) {
				t = TOperator
			} else if IsAlpha(c) {
				fi, str := IsFunction(input[i:])
				if fi {
					i += len(str) - 1
					t = TFunction
					tokens = append(tokens,&Token{t, str})
					tokens = append(tokens,&Token{TOperator, "F"})
					continue
				} else {
					t = TVariable
				}
			}

			//This is wrong as it assumes that each token is only one character long
			if t != TUnknown {
				tokens = append(tokens,&Token{t,input[i:i+1]})
			}
		}
	}
	if buf.Len() > 0 {
		tokens = append(tokens,&Token{TNumber, buf.String()})
	}
	return tokens
}

//ParseExpression and validate syntax, also expand any 'shortcuts'
func Validate(tokens []*Token) ([]*Token , error) {
	lt := TUnknown
	passtwo := make([]*Token, len(tokens)*2)
	tokc := 0
	for	i := 0; i < len(tokens); i++ {
		switch tokens[i].kind {
		case TLparen:
			if lt == TVariable || lt == TNumber || lt == TRparen {
				// Implicit multiplication
				passtwo[tokc] = &Token{TOperator, "*"}
				tokc++
			}
			passtwo[tokc] = tokens[i]
			tokc++
		case TRparen:
			if lt == TOperator {
				return nil, errors.New("Invalid syntax, Closing Paren cannot follow operator")
			}
			passtwo[tokc] = tokens[i]
			tokc++
		case TOperator:
			if lt == TOperator || lt == TLparen {
				return nil, errors.New("Invalid syntax, improper operator placement.")
			}
			passtwo[tokc] = tokens[i]
			tokc++
		case TVariable, TNumber, TFunction:
			if lt == TVariable || lt == TNumber || lt == TRparen {
				passtwo[tokc] = &Token{TOperator, "*"}
				tokc++
			}
			passtwo[tokc] = tokens[i]
			tokc++
		}
		lt = tokens[i].kind
	}
	return passtwo[:tokc], nil
}

//returns true if operator a has higher precedence than b
func comparePrecedence(a, b int) bool {
	if a == b {
		return false
	}

	if a == OFun {
		return false
	}

	if a == OPow {
		return true
	}

	if (a == OMul || a == ODiv) && (b == OAdd || b == OSub) {
		return true
	}

	return false
}

func OpSignToConst(op string) (rop int) {
	switch op {
	case "+":
		rop = OAdd
	case "-":
		rop = OSub
	case "*":
		rop = OMul
	case "/":
		rop = ODiv
	case "^":
		rop = OPow
	case "F":
		rop = OFun
	}
	return rop
}

func build(tokens []*Token) Equatable {
	stack := NewTokStack(len(tokens))
	postfix := NewTokStack(len(tokens))
	for _,t := range tokens {
		switch t.kind {
		case TNumber, TVariable, TFunction:
			postfix.Push(t)
		case TLparen:
			stack.Push(t)
		case TOperator:
			for stack.Size() > 0 && stack.Peek().kind != TLparen {
				if comparePrecedence(OpSignToConst(stack.Peek().val),
						OpSignToConst(t.val)) {
					postfix.Push(stack.Pop())
				} else {
					break
				}
			}
			stack.Push(t)
		case TRparen:
			for stack.Size() > 0 && stack.Peek().kind != TLparen {
				postfix.Push(stack.Pop())
			}
			if stack.Size() > 0 {
				stack.Pop()
			}
		}
	}
	for stack.Size() > 0 {
		postfix.Push(stack.Pop())
	}
	eqs := make([]Equatable, len(postfix.GetSlice()))
	eqsc := 0
	for _,t :=  range postfix.GetSlice() {
		if t.kind == TVariable {
			eqs[eqsc] = NewVariable(t.val)
			eqsc++
		} else if t.kind == TNumber {
			eqs[eqsc] = NewConstant(t.val)
			eqsc++
		} else if t.kind == TFunction {
			eqs[eqsc] = &Function{FStrToConst(t.val),nil}
			eqsc++
		} else if t.kind == TOperator {
			op := OpSignToConst(t.val)
			if op == OFun {
				tpar := eqs[eqsc - 1]
				tfun,_ := eqs[eqsc - 2].(*Function)
				tfun.arg = tpar
				eqsc--
			} else {
				neq := &Term{eqs[eqsc - 2] ,eqs[eqsc - 1], op}
				eqsc--
				eqs[eqsc - 1] = neq
			}
		}
	}
	return eqs[0]
}

func ParseEquation(input string) (*Equality, error) {
	if !strings.Contains(input,"=") {
		return nil, errors.New("Not a valid equality, must contain '='.")
	}
	spl := strings.Split(input, "=")
	l,lerr := ParseExpression(spl[0])
	r,rerr := ParseExpression(spl[1])
	if lerr != nil {
		return nil, lerr
	}
	if rerr != nil {
		return nil, rerr
	}
	l = Simplify(l)
	r = Simplify(r)
	return &Equality{l,r}, nil
}

func SimilarOp(a, b int) bool {
	if (a == OAdd || a == OSub) && (b == OAdd || b == OSub) {
		return true
	} else if (a == OMul || a == ODiv) && (b == OMul || b == ODiv) {
		return true
	}
	return false
}

func ParseExpression(input string) (Equatable, error) {
	tokens := Tokenize(input)
	if len(tokens) == 1 {
		if tokens[0].kind == TVariable {
			return NewVariable(tokens[0].val),nil
		} else if tokens[0].kind == TNumber {
			return NewConstant(tokens[0].val),nil
		}
	}
	tokens,err := Validate(tokens)
	if err != nil {
		return nil, err
	}
	eq := build(tokens)
	return eq, nil
}

func Interpreter() {
	run := true
	stdin := bufio.NewReader(os.Stdin)
	var eq *Equality
	for run {
		line,_,_ := stdin.ReadLine()
		if len(line) == 0 {
			continue
		}
		if line[0] == ':' {
			eq,_ = ParseEquation(string(line[1:]))
		} else if line[0] == '?' {
			v := Vars[line[1]]
			if v != nil {
				fmt.Printf("%s = %f\n", string(line[1:2]), v.Value())
			} else {
				fmt.Printf("Error: %s does not exist!\n", line[1:2])
			}
		} else if line[0] == '!' {
			fmt.Printf("Solving for %s in:\n\t%s\n", string(line[1:2]), eq.Print())
			ans, er := eq.SolveFor(line[1])
			fmt.Printf("%s = %f\n", string(line[1:2]), ans)
			fmt.Printf("\terror: %f\n",er)
		} else {
			if string(line) == "quit" {
				return
			}
			if strings.Contains(string(line),"=") {
				v := Vars[line[0]]
				if v != nil {
					v.val,_ = strconv.ParseFloat(string(line[2:]),64)
				} else {
					fmt.Printf("Error: %s does not exist!\n", line[0:1])
				}
			}
		}
	}
	eq.SolveFor(' ')
}

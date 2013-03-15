package main

import (
	"strconv"
	"fmt"
)

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

func (c *Constant) ContainsVar(v uint8) bool {
	return false
}

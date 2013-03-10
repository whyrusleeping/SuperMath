package main

type TokStack struct {
	items []*Token
	clen int
}

func (ts *TokStack) Push(t *Token) {
	if t == nil {
		panic("nil token cant be pushed!!")
	}
	ts.items[ts.clen] = t
	ts.clen++
}

func (ts *TokStack) Pop() *Token {
	t := ts.items[ts.clen - 1]
	ts.items[ts.clen - 1] = nil
	ts.clen--
	return t
}

func (ts *TokStack) Peek() *Token {
	return ts.items[ts.clen - 1]
}

func (es *TokStack) GetSlice() []*Token {
	return es.items[:es.clen]
}

func (ts *TokStack) Size() int {
	return ts.clen
}

func (ts *TokStack) Clear() {
	for ts.clen > 0 {
		ts.Pop()
	}
}

func NewTokStack(size int) *TokStack {
	ts := TokStack{make([]*Token,size), 0}
	return &ts
}

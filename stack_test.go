package main

import "testing"

func TestStack(t *testing.T) {
	s := NewTokStack(32)
	for i := 0; i < 10; i++ {
		s.Push(&Token{byte(i),""})
	}
	for i := 0; i < 10; i++ {
		tk := s.Pop()
		if tk == nil {
			t.Fatalf("Not supposed to be nil... %d",i)
		}
	}
}

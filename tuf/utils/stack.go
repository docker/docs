package utils

import (
	"fmt"
)

type ErrEmptyStack struct {
	action string
}

func (err ErrEmptyStack) Error() string {
	return fmt.Sprintf("attempted to %s with empty stack", err.action)
}

type ErrBadTypeCast struct{}

func (err ErrBadTypeCast) Error() string {
	return "attempted to do a typed pop and item was not of type"
}

type Stack []interface{}

func NewStack() *Stack {
	s := make(Stack, 0)
	return &s
}

func (s *Stack) Push(item interface{}) {
	*s = append(*s, item)
}

func (s *Stack) Pop() (interface{}, error) {
	l := len(*s)
	if l > 0 {
		item := (*s)[l-1]
		*s = (*s)[:l-1]
		return item, nil
	}
	return nil, ErrEmptyStack{action: "pop"}
}

func (s *Stack) PopString() (string, error) {
	l := len(*s)
	if l > 0 {
		item := (*s)[l-1]
		if item, ok := item.(string); ok {
			*s = (*s)[:l-1]
			return item, nil
		}
		return "", ErrBadTypeCast{}
	}
	return "", ErrEmptyStack{action: "pop"}
}

func (s Stack) Empty() bool {
	return len(s) == 0
}

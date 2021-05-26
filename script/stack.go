package script

import (
	"container/list"	
	"fmt"
)


type Stack struct {
	stack *list.List
}

func (s *Stack) Push(opcode []byte) {
	s.stack.PushFront(opcode)
}

func (s *Stack) Pop() error {
	if s.stack.Len() > 0 {
		opcode := s.stack.Front()
		s.stack.Remove(opcode)
	}
	return fmt.Errorf("Pop Error: Queue is empty")
}

func (s *Stack) Front() ([]byte, error) {
	if s.stack.Len() > 0 {
		if opcode, ok := s.stack.Front().Value.([]byte); ok {
			return opcode, nil
		}
		return nil, fmt.Errorf("Peep Error: Queue Datatype is incorrect")
	}
	return nil, fmt.Errorf("Peep Error: Queue is empty")
}

func (s *Stack) Size() int {
	return s.stack.Len()
}

func (s *Stack) Empty() bool {
	return s.stack.Len() == 0
}

package script

import (
    "fmt"
    "sync"
)

type Stack struct {
    stack [][]byte
    lock  sync.RWMutex
}

func (s *Stack) Push(opcode []byte) {
    s.lock.Lock()
    defer s.lock.Unlock()
    s.stack = append(s.stack, opcode)
}

func (s *Stack) Pop() error {
    len := len(s.stack)
    if len > 0 {
        s.lock.Lock()
        defer s.lock.Unlock()
        s.stack = s.stack[:len-1]
        return nil
    }
    return fmt.Errorf("Pop Error: Queue is empty")
}

func (s *Stack) Front() ([]byte, error) {
    len := len(s.stack)
    if len > 0 {
        s.lock.Lock()
        defer s.lock.Unlock()
        return s.stack[len-1], nil
    }
    return nil, fmt.Errorf("Peep Error: Queue is empty")
}

func (s *Stack) Size() int {
    return len(s.stack)
}

func (s *Stack) Empty() bool {
    return len(s.stack) == 0
}
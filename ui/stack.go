package ui

import (
	"sync"

	"github.com/ebitenui/ebitenui/widget"
)

type UIStackElement struct {
	Widget     widget.HasWidget
	RemoveFunc widget.RemoveChildFunc
}

type UIStack struct {
	stack []UIStackElement
	mu    sync.Mutex
}

// Push adds a widget and its remove function to the UIStack.
// Push is concurrent-safe.
func (s *UIStack) Push(w widget.HasWidget, removeFunc widget.RemoveChildFunc) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.stack = append(s.stack, UIStackElement{
		Widget: w, RemoveFunc: removeFunc,
	})
}

// Pop removes the last added element from the UIStack (and from ebiten UI),
// and returns it if it exists, and whether it actually does or not.
// Pop is concurrent-safe.
func (s *UIStack) Pop() (el UIStackElement, ok bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.stack) == 0 {
		return UIStackElement{}, false
	}
	poppedElement := s.stack[len(s.stack)-1]
	poppedElement.RemoveFunc()
	s.stack = s.stack[:len(s.stack)-1]
	return poppedElement, true
}

// Peek returns the element last added to the UIStack,
// or nil if it is empty. Peek is concurrent-safe.
func (s *UIStack) Peek() (el *UIStackElement) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.stack) == 0 {
		return nil
	}
	return &(s.stack[len(s.stack)-1])
}

func NewUIStack() *UIStack {
	return &UIStack{
		mu:    sync.Mutex{},
		stack: make([]UIStackElement, 0, 10),
	}
}

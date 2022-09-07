package ws

import "github.com/google/uuid"

type handler[T any] struct {
	id string
	fn func(T)
}

// Emitter is an event emitter for the event type T.
type Emitter[T any] struct {
	handlers []handler[T]
}

// NewEmitter creates a new event emitter for the event type T.
func NewEmitter[T any]() *Emitter[T] {
	return &Emitter[T]{
		handlers: make([]handler[T], 0),
	}
}

// On registers a new event handler for the event type T.
// The returned function can be used to remove the handler.
func (e *Emitter[T]) On(fn func(T)) func() {
	id := uuid.New().String()
	e.handlers = append(e.handlers, handler[T]{id, fn})
	cancel := func() {
		e.off(id)
	}
	return cancel
}

// off removes an event handler by id.
func (e *Emitter[T]) off(id string) {
	for i, h := range e.handlers {
		if h.id == id {
			e.handlers = append(e.handlers[:i], e.handlers[i+1:]...)
			return
		}
	}
}

// Emit triggers an event of type T.
func (e *Emitter[T]) Emit(data T) {
	for _, h := range e.handlers {
		h.fn(data)
	}
}

func (e *Emitter[T]) Clear() {
	e.handlers = make([]handler[T], 0)
}

// Example usage:
//
// package main
//
// import (
//
//	"fmt"
//	"wsd/ws"
//
// )
//
//  type EventA struct {
//		Message string
//  }
//
//	func main() {
//		emitter := ws.NewEmitter[ws.EventA]()
//		id := emitter.On(func(data ws.EventA) {
//			fmt.Println(data.Message)
//		})
//		emitter.Emit(ws.EventA{"hello"})
//		emitter.Off(id)
//		emitter.Emit(ws.EventA{"world"})
//	}
//
// Output:
// hello

// Example usage:
//
//  type EventA struct {
//		Message string
//  }
//  type Program struct {
//		EventA ws.Emitter[ws.EventA]
//  }
//
//	func main() {
//		program := Program{}
//		remove := program.EventA.On(func(data ws.EventA) {
//			fmt.Println(data.Message)
//		})
// 		defer remove()
//		program.EventA.Emit(ws.EventA{"hello"})
//	}
//
// Output:
// hello

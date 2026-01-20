package main

import (
	"sync"
)

// Resetter — это интерфейс .
type Resetter interface {
	Reset()
}

// Pool — это обертка над sync.Pool.
type Pool[T Resetter] struct {
	internal sync.Pool
}

// New конструктор Pool
func New[T Resetter](creator func() T) *Pool[T] {
	return &Pool[T]{
		internal: sync.Pool{
			New: func() any {
				return creator()
			},
		},
	}
}

// Get получает объект из пула.
func (p *Pool[T]) Get() T {
	return p.internal.Get().(T)
}

// Put возвращает объект в пул, предварительно сбрасывая его состояние.
func (p *Pool[T]) Put(x T) {
	if any(x) == nil {
		return
	}
	x.Reset()
	p.internal.Put(x)
}

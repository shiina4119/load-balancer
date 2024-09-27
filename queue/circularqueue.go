package queue

import "fmt"

type CircularQueue[T any] struct {
	Data    []T
	IsFull  bool
	IsEmpty bool
	start   int
	end     int
}

func NewCircularQueue[T any](capacity int) *CircularQueue[T] {
	return &CircularQueue[T]{
		Data:    make([]T, capacity),
		IsFull:  false,
		IsEmpty: true,
		start:   0,
		end:     0,
	}
}

func (r *CircularQueue[T]) Push(elem T) error {
	if r.IsFull {
		return fmt.Errorf("Out of bounds push, queue is full")
	}

	r.Data[r.end] = elem
	r.end = (r.end + 1) % len(r.Data)
	r.IsFull = r.end == r.start
	r.IsEmpty = false

	return nil
}

func (r *CircularQueue[T]) Pop() (T, error) {
	var res T
	if r.IsEmpty {
		return res, fmt.Errorf("Cannot pop, queue is empty")
	}

	res = r.Data[r.start]
	r.start = (r.start + 1) % len(r.Data)
	r.IsFull = false
	r.IsEmpty = r.end == r.start

	return res, nil
}

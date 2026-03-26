package list

import (
	"github.com/samber/mo"
)

// RingBuffer is a fixed-capacity FIFO circular buffer.
// When full, Push overwrites the oldest item and returns it.
type RingBuffer[T any] struct {
	buf  []T
	head int
	size int
}

// NewRingBuffer creates a ring buffer with fixed capacity.
// capacity <= 0 creates an empty non-writable buffer.
func NewRingBuffer[T any](capacity int) *RingBuffer[T] {
	if capacity < 0 {
		capacity = 0
	}
	return &RingBuffer[T]{
		buf: make([]T, capacity),
	}
}

// Capacity returns max item capacity.
func (r *RingBuffer[T]) Capacity() int {
	if r == nil {
		return 0
	}
	return len(r.buf)
}

// Len returns current item count.
func (r *RingBuffer[T]) Len() int {
	if r == nil {
		return 0
	}
	return r.size
}

// IsEmpty reports whether buffer has no items.
func (r *RingBuffer[T]) IsEmpty() bool {
	return r.Len() == 0
}

// IsFull reports whether buffer reached capacity.
func (r *RingBuffer[T]) IsFull() bool {
	if r == nil {
		return false
	}
	return r.size == len(r.buf) && len(r.buf) > 0
}

// Push writes value at tail.
// If full, oldest value is evicted and returned as mo.Option.
func (r *RingBuffer[T]) Push(value T) mo.Option[T] {
	if r == nil || len(r.buf) == 0 {
		return mo.None[T]()
	}

	if r.size < len(r.buf) {
		tail := (r.head + r.size) % len(r.buf)
		r.buf[tail] = value
		r.size++
		return mo.None[T]()
	}

	evicted := r.buf[r.head]
	r.buf[r.head] = value
	r.head = (r.head + 1) % len(r.buf)
	return mo.Some(evicted)
}

// Pop removes and returns oldest value.
func (r *RingBuffer[T]) Pop() (T, bool) {
	var zero T
	if r == nil || r.size == 0 {
		return zero, false
	}
	value := r.buf[r.head]
	r.buf[r.head] = zero
	r.head = (r.head + 1) % len(r.buf)
	r.size--
	if r.size == 0 {
		r.head = 0
	}
	return value, true
}

// Peek returns oldest value without removing it.
func (r *RingBuffer[T]) Peek() (T, bool) {
	var zero T
	if r == nil || r.size == 0 {
		return zero, false
	}
	return r.buf[r.head], true
}

// Values returns items from oldest to newest.
func (r *RingBuffer[T]) Values() []T {
	if r == nil || r.size == 0 {
		return nil
	}
	out := make([]T, r.size)
	for i := range r.size {
		out[i] = r.buf[(r.head+i)%len(r.buf)]
	}
	return out
}

// Clear removes all values.
func (r *RingBuffer[T]) Clear() {
	if r == nil {
		return
	}
	var zero T
	for i := range r.size {
		r.buf[(r.head+i)%len(r.buf)] = zero
	}
	r.head = 0
	r.size = 0
}

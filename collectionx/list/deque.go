package list

// Deque is a growable double-ended queue backed by a ring buffer.
// Zero value is ready to use.
type Deque[T any] struct {
	buf  []T
	head int
	size int
}

// NewDeque creates an empty deque with optional initial items.
func NewDeque[T any](items ...T) *Deque[T] {
	d := &Deque[T]{}
	d.PushBack(items...)
	return d
}

// PushFront inserts items at front in argument order.
// Example: PushFront(1,2) -> front sequence becomes [1,2,...].
func (d *Deque[T]) PushFront(items ...T) {
	if len(items) == 0 {
		return
	}
	// Keep argument order near front.
	for _, item := range items {
		d.pushFrontOne(item)
	}
}

// PushBack appends items at back.
func (d *Deque[T]) PushBack(items ...T) {
	if len(items) == 0 {
		return
	}
	for _, item := range items {
		d.pushBackOne(item)
	}
}

// PopFront removes and returns front item.
func (d *Deque[T]) PopFront() (T, bool) {
	var zero T
	if d.size == 0 {
		return zero, false
	}
	value := d.buf[d.head]
	d.buf[d.head] = zero
	d.head = (d.head + 1) % len(d.buf)
	d.size--
	if d.size == 0 {
		d.head = 0
	}
	return value, true
}

// PopBack removes and returns back item.
func (d *Deque[T]) PopBack() (T, bool) {
	var zero T
	if d.size == 0 {
		return zero, false
	}
	idx := d.physicalIndex(d.size - 1)
	value := d.buf[idx]
	d.buf[idx] = zero
	d.size--
	if d.size == 0 {
		d.head = 0
	}
	return value, true
}

// Front returns front item without removing it.
func (d *Deque[T]) Front() (T, bool) {
	var zero T
	if d.size == 0 {
		return zero, false
	}
	return d.buf[d.head], true
}

// Back returns back item without removing it.
func (d *Deque[T]) Back() (T, bool) {
	var zero T
	if d.size == 0 {
		return zero, false
	}
	return d.buf[d.physicalIndex(d.size-1)], true
}

// Get returns item at logical index from front.
func (d *Deque[T]) Get(index int) (T, bool) {
	var zero T
	if index < 0 || index >= d.size {
		return zero, false
	}
	return d.buf[d.physicalIndex(index)], true
}

// Len returns item count.
func (d *Deque[T]) Len() int {
	return d.size
}

// IsEmpty reports whether deque has no items.
func (d *Deque[T]) IsEmpty() bool {
	return d.Len() == 0
}

// Clear removes all items.
func (d *Deque[T]) Clear() {
	var zero T
	for i := range d.size {
		d.buf[d.physicalIndex(i)] = zero
	}
	d.head = 0
	d.size = 0
}

// Values returns items from front to back.
func (d *Deque[T]) Values() []T {
	if d.size == 0 {
		return nil
	}
	out := make([]T, d.size)
	for i := range d.size {
		out[i] = d.buf[d.physicalIndex(i)]
	}
	return out
}

// Range iterates items from front to back until fn returns false.
func (d *Deque[T]) Range(fn func(index int, item T) bool) {
	if fn == nil {
		return
	}
	for i := range d.size {
		if !fn(i, d.buf[d.physicalIndex(i)]) {
			return
		}
	}
}

func (d *Deque[T]) pushFrontOne(item T) {
	d.ensureCapacity(1)
	if d.size == 0 {
		d.head = 0
		d.buf[0] = item
		d.size = 1
		return
	}
	d.head = (d.head - 1 + len(d.buf)) % len(d.buf)
	d.buf[d.head] = item
	d.size++
}

func (d *Deque[T]) pushBackOne(item T) {
	d.ensureCapacity(1)
	idx := d.physicalIndex(d.size)
	d.buf[idx] = item
	d.size++
}

func (d *Deque[T]) ensureCapacity(extra int) {
	if d.buf == nil {
		capacity := 4
		for capacity < extra {
			capacity *= 2
		}
		d.buf = make([]T, capacity)
		return
	}

	need := d.size + extra
	if need <= len(d.buf) {
		return
	}

	newCap := len(d.buf) * 2
	newCap = max(newCap, 4)
	for newCap < need {
		newCap *= 2
	}

	newBuf := make([]T, newCap)
	for i := range d.size {
		newBuf[i] = d.buf[d.physicalIndex(i)]
	}
	d.buf = newBuf
	d.head = 0
}

func (d *Deque[T]) physicalIndex(logicalIndex int) int {
	return (d.head + logicalIndex) % len(d.buf)
}

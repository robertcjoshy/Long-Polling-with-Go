package handler

import "sync"

type Cappedqueue[T any] struct {
	items    []T
	lock     *sync.RWMutex
	capacity int
}

func Newcappedqueue[T any](capacity int) *Cappedqueue[T] {
	return &Cappedqueue[T]{
		items:    make([]T, 0, capacity),
		lock:     new(sync.RWMutex),
		capacity: capacity,
	}
}

func (s *Cappedqueue[T]) Append(item T) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if l := len(s.items); l == 0 {
		s.items = append(s.items, item)
	} else {
		to := s.capacity - 1
		if l < s.capacity {
			to = l
		}
		s.items = append([]T{item}, s.items[:to]...)
	}

}

func (s *Cappedqueue[T]) Copy() []T {
	s.lock.RLock()
	defer s.lock.RUnlock()

	copied := make([]T, len(s.items))

	copy(copied, s.items)

	return copied

}

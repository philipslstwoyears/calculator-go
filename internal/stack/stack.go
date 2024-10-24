package stack

type element interface {
	int | float64 | string
}

type Stack[T element] struct {
	items    []T
	topIndex int
}

func New[T element]() *Stack[T] {
	return &Stack[T]{
		items:    make([]T, 0),
		topIndex: -1,
	}
}

func (s *Stack[T]) IsEmpty() bool {
	return s.topIndex == -1
}

func (s *Stack[T]) Push(item T) {
	s.topIndex++
	s.items = append(s.items, item)
}

func (s *Stack[T]) Pop() T {
	popedElement := s.items[s.topIndex]
	s.topIndex--
	s.items = s.items[:s.topIndex+1]
	return popedElement
}

func (s *Stack[T]) Peek() T {
	return s.items[s.topIndex]
}

func (s *Stack[T]) Size() int {
	return len(s.items)
}

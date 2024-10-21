package stack

type Stack struct {
	elements []int
	topIndex int
}

func New() *Stack {
	return &Stack{elements: []int{}, topIndex: -1}
}
func (s *Stack) Push(v int) {
	s.elements = append(s.elements, v)
	s.topIndex = s.topIndex + 1
}
func (s *Stack) Pop() int {
	element := s.elements[s.topIndex]
	s.elements = s.elements[:s.topIndex]
	s.topIndex = s.topIndex - 1
	return element

}

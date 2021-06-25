package stack

type Stack []interface{}

func New() *Stack {
	return new(Stack)
}

func (s *Stack) Clear() {
	*s = (*s)[0:0]
}

func (s *Stack) IsEmpty() bool {
	return len(*s) == 0
}

func (s *Stack) Push(e interface{}) {
	*s = append(*s, e)
}

func (s *Stack) PushAll(e []interface{}) {
	for _, v := range e {
		s.Push(v)
	}
}

func (s *Stack) Pop() (interface{}, bool) {
	if s.IsEmpty() {
		return nil, false
	} else {
		index := len(*s) - 1
		element := (*s)[index]
		*s = (*s)[:index]
		return element, true
	}
}

func (s *Stack) Swap(i int, j int) {
	(*s)[i], (*s)[j] = (*s)[j], (*s)[i]
}

func (s *Stack) Size() int {
	return len(*s)
}

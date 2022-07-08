package lib

type SliceIter struct {
	Items []interface{}
	index int
}

func (s *SliceIter) Next() bool {
	if s.index == len(s.Items) {
		return false
	}
	s.index += 1

	return true
}

func (s SliceIter) Current() interface{} {
	return s.Items[s.index-1]
}

func (s SliceIter) Err() error {
	return nil
}

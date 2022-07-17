package lib

type SliceIter struct {
	Items []interface{}
	count int
}

func (s *SliceIter) Next() bool {
	if s.lastItem() {
		return false
	}
	s.count += 1

	return true
}

func (s *SliceIter) Current() interface{} {
	return s.Items[s.count-1]
}

func (s *SliceIter) Err() error {
	return nil
}

func (s *SliceIter) lastItem() bool {
	return s.count == len(s.Items)
}

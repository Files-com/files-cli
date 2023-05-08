package lib

type CellWrapper struct {
	cell     string
	data     interface{}
	Iterable bool
}

func (c CellWrapper) String() string {
	return c.cell
}

func (c CellWrapper) Data() interface{} {
	return c.Data
}

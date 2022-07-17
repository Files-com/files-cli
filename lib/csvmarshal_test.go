package lib

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Person struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Age       int    `json:"age"`
}

func TestCSVMarshal(t *testing.T) {
	a := assert.New(t)
	p1 := Person{FirstName: "Dustin", LastName: "Zeisler", Age: 100}
	buf := bytes.NewBufferString("")

	CSVMarshal(p1, "", buf)

	a.Equal(`first_name,last_name,age
Dustin,Zeisler,100
`, buf.String())
}

type PersonNil struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Age         int    `json:"age"`
	DontShowNil *bool  `json:"dont-show-nil"`
}

func TestCSVMarshalNil(t *testing.T) {
	a := assert.New(t)
	p1 := PersonNil{FirstName: "Dustin", LastName: "Zeisler", Age: 100}
	buf := bytes.NewBufferString("")

	CSVMarshal(p1, "", buf)

	a.Equal(`first_name,last_name,age,dont-show-nil
Dustin,Zeisler,100,
`, buf.String())
}

func TestCSVMarshal_Fields(t *testing.T) {
	a := assert.New(t)
	p1 := Person{FirstName: "Dustin", LastName: "Zeisler", Age: 100}
	buf := bytes.NewBufferString("")

	CSVMarshal(p1, "first_name,last_name", buf)

	a.Equal(`first_name,last_name
Dustin,Zeisler
`, buf.String())
}

type MockIter struct {
	SliceIter
	eofPage func(*MockIter) bool
}

func (m *MockIter) EOFPage() bool {
	if m.eofPage == nil {
		return m.lastItem()
	} else {
		return m.eofPage(m)
	}
}

func (m *MockIter) NextPage() bool {
	return true
}

func (m *MockIter) GetPage() bool {
	return true
}

func TestCSVMarshalIter(t *testing.T) {
	a := assert.New(t)
	p1 := Person{FirstName: "Dustin", LastName: "Zeisler", Age: 100}
	p2 := Person{FirstName: "Tom", LastName: "Smith", Age: 99}
	it := MockIter{SliceIter: SliceIter{Items: []interface{}{p1, p2}}}
	buf := bytes.NewBufferString("")

	CSVMarshalIter(&it, "", nil, buf)

	a.Equal(`first_name,last_name,age
Dustin,Zeisler,100
Tom,Smith,99
`, buf.String())
}
func TestCSVMarshalIter_FilterIter(t *testing.T) {
	a := assert.New(t)
	p1 := Person{FirstName: "Dustin", LastName: "Zeisler", Age: 100}
	p2 := Person{FirstName: "Tom", LastName: "Smith", Age: 99}
	it := MockIter{SliceIter: SliceIter{Items: []interface{}{p1, p2}}}
	buf := bytes.NewBufferString("")

	CSVMarshalIter(&it, "", func(i interface{}) bool {
		return i.(Person).FirstName == "Dustin"
	}, buf)

	a.Equal(`first_name,last_name,age
Dustin,Zeisler,100
`, buf.String())
}

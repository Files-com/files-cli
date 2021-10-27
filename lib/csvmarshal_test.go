package lib

import (
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
	output := CaptureOutput(func() {
		CSVMarshal(p1, "")
	})

	a.Equal(`first_name,last_name,age
Dustin,Zeisler,100
`, output)
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
	output := CaptureOutput(func() {
		CSVMarshal(p1, "")
	})

	a.Equal(`first_name,last_name,age,dont-show-nil
Dustin,Zeisler,100,
`, output)
}

func TestCSVMarshal_Fields(t *testing.T) {
	a := assert.New(t)
	p1 := Person{FirstName: "Dustin", LastName: "Zeisler", Age: 100}
	output := CaptureOutput(func() {
		CSVMarshal(p1, "first_name,last_name")
	})

	a.Equal(`first_name,last_name
Dustin,Zeisler
`, output)
}

type MockIter struct {
	People  []Person
	index   int
	eofPage bool
}

func (m *MockIter) Next() bool {
	if m.index == len(m.People) {
		return false
	}
	m.index += 1

	return true
}

func (m MockIter) Current() interface{} {
	return m.People[m.index-1]
}

func (m MockIter) Err() error {
	return nil
}

func (m MockIter) EOFPage() bool {
	return m.eofPage
}

func TestCSVMarshalIter(t *testing.T) {
	a := assert.New(t)
	p1 := Person{FirstName: "Dustin", LastName: "Zeisler", Age: 100}
	p2 := Person{FirstName: "Tom", LastName: "Smith", Age: 99}
	it := MockIter{People: []Person{p1, p2}}
	output := CaptureOutput(func() {
		CSVMarshalIter(&it, "", nil)
	})

	a.Equal(`first_name,last_name,age
Dustin,Zeisler,100
Tom,Smith,99
`, output)
}
func TestCSVMarshalIter_FilterIter(t *testing.T) {
	a := assert.New(t)
	p1 := Person{FirstName: "Dustin", LastName: "Zeisler", Age: 100}
	p2 := Person{FirstName: "Tom", LastName: "Smith", Age: 99}
	it := MockIter{People: []Person{p1, p2}}
	output := CaptureOutput(func() {
		CSVMarshalIter(&it, "", func(i interface{}) bool {
			return i.(Person).FirstName == "Dustin"
		})
	})

	a.Equal(`first_name,last_name,age
Dustin,Zeisler,100
`, output)
}

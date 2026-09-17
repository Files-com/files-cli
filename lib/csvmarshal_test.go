package lib

import (
	"bytes"
	"os"
	"testing"

	"github.com/Files-com/files-cli/lib/ptytest"
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

	CSVMarshal(p1, []string{}, buf, "")

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

	CSVMarshal(p1, []string{}, buf, "")

	a.Equal(`first_name,last_name,age,dont-show-nil
Dustin,Zeisler,100,
`, buf.String())
}

func TestCSVMarshal_Fields(t *testing.T) {
	a := assert.New(t)
	p1 := Person{FirstName: "Dustin", LastName: "Zeisler", Age: 100}
	buf := bytes.NewBufferString("")

	CSVMarshal(p1, []string{"first_name", "last_name"}, buf, "")

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

	CSVMarshalIter(&it, []string{}, nil, buf, "")

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

	CSVMarshalIter(&it, []string{}, func(i interface{}) (interface{}, bool, error) {
		return i, i.(Person).FirstName == "Dustin", nil
	}, buf, "")

	a.Equal(`first_name,last_name,age
Dustin,Zeisler,100
`, buf.String())
}

func TestCSVMarshal_TerminalOutputEscapesControls(t *testing.T) {
	remote := RemoteFile{Path: oscClipboardST, DisplayName: c1AndDEL}

	got := ptytest.Capture(t, func(slave *os.File) error {
		return CSVMarshal(remote, []string{}, slave, "")
	})

	assertNoTerminalControls(t, got)
	assert.Contains(t, got, "path,display_name")
	assert.Contains(t, got, oscClipboardSTEscaped+","+c1AndDELEscaped)
}

func TestCSVMarshalIter_RedirectedOutputKeepsExactValues(t *testing.T) {
	it := &SliceIter{Items: []interface{}{
		RemoteFile{Path: oscClipboardBEL, DisplayName: unicodeName},
		RemoteFile{Path: "notes.txt", DisplayName: "multi\nline\ttab"},
	}}

	got := capturePipeOutput(t, func(w *os.File) error {
		return CSVMarshalIter(it, []string{}, nil, w, "")
	})

	assert.Equal(t, "path,display_name\n"+
		oscClipboardBEL+","+unicodeName+"\n"+
		"notes.txt,\"multi\nline\ttab\"\n", got)
}

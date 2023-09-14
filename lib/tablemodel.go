package lib

import (
	"context"
	"fmt"
	"io"
	"strings"
	"sync"
	"time"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
	"github.com/olekukonko/ts"
	"github.com/samber/lo"
)

type tableModel struct {
	table.Model
	*sync.Once
	tableLoader
	*tea.Program
	fields          []string
	parentResources []interface{}
	filterTextInput textinput.Model
	out             io.Writer
	maxColumnWidth  map[string]int
	FilterIter
}

var tableBoarder = table.Border{
	Top:    "─",
	Left:   "│",
	Right:  "│",
	Bottom: "─",

	TopRight:    "┐",
	TopLeft:     "┌",
	BottomRight: "┘",
	BottomLeft:  "└",

	TopJunction:    "┬",
	LeftJunction:   "├",
	RightJunction:  "┤",
	BottomJunction: "┴",
	InnerJunction:  "┼",

	InnerDivider: "│",
}

type Id interface {
	Id() interface{}
}

func (t *tableModel) Init() tea.Cmd {
	t.maxColumnWidth = make(map[string]int)
	return nil
}

func (t *tableModel) View() string {
	if t.TotalRows() == 0 && t.Context().Err() == nil && !t.filterTextInput.Focused() {
		return t.Spinner().View()
	}

	if t.tableLoader.Err() != nil {
		return t.tableLoader.Err().Error()
	}

	body := strings.Builder{}
	if t.TotalRows() == 0 && t.Context().Err() != nil {
		body.WriteString(fmt.Sprintf("%v (no rows)", t.currentResource()))
	} else {
		body.WriteString(t.Model.View())
	}

	body.WriteString("\nquit [q] " + t.ControlsFooter())
	body.WriteString(" filter [/]")
	if t.filterTextInput.Focused() {
		body.WriteString(fmt.Sprintf("\n%v\n", t.filterTextInput.View()))
	}

	return body.String()
}

func (t *tableModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			t.Cancel()
			return t, tea.Quit
		}

		if t.filterTextInput.Focused() {
			switch msg.String() {
			case "enter":
				t.filterTextInput.Blur()
			case "esc":
				t.filterTextInput.Reset()
				t.filterTextInput.Blur()
			default:
				t.filterTextInput, cmd = t.filterTextInput.Update(msg)
			}
			t.Model = t.WithFilterInput(t.filterTextInput)
			return t, cmd
		}

		switch msg.String() {
		case "/":
			t.filterTextInput.Focus()
		case "esc":
			t.Cancel()
		case "backspace":
			if t.GetIsFilterInputFocused() {
				t.Model, cmd = t.Model.Update(msg)
				return t, cmd
			}
		case "q", "ctrl+c":
			t.Cancel()
			return t, tea.Quit
		}
	}

	t.tableLoader, cmd = t.tableLoader.Update(t, msg)
	if cmd != nil {
		return t, cmd
	}

	t.Model, cmd = t.Model.Update(msg)
	return t, cmd
}

func (t *tableModel) dimensions() []int {
	size, err := ts.GetSize()
	var width int
	var height int
	if err == nil {
		width = size.Col()
		height = size.Row()
		if height > 9 {
			height = height - 9
		}
	} else {
		width = 400
		height = 200
	}

	return []int{width, height}
}

func (t *tableModel) LoadFirstPage() error {
	t.Once = &sync.Once{}
	err := t.tableLoader.LoadFirstPage(t)
	return err
}

func (t *tableModel) LoadRest() {
	t.tableLoader.LoadRest(t)
}

func (t *tableModel) LoadAndBuild(ctx context.Context) (*tea.Program, error) {
	t.Program = tea.NewProgram(t, tea.WithContext(ctx), tea.WithOutput(t.out))
	t.SetLoader(t.tableLoader)
	return t.Load()
}

func (t *tableModel) SetLoader(loader tableLoader) {
	t.filterTextInput = textinput.New()
	t.Model = table.Model{}
	t.tableLoader = loader
}

func (t *tableModel) Load() (*tea.Program, error) {
	t.tableLoader.Load()
	go func() {
		tick := time.Tick(time.Millisecond * 250)
		for {
			select {
			case <-tick:
				t.updateFooter()
			case <-t.Context().Done():
				t.updateFooter()
				return
			}
		}
	}()
	err := t.LoadFirstPage()
	if err != nil {
		t.Cancel()
		return t.Program, nil
	}
	t.LoadRest()
	return t.Program, nil
}

func (t *tableModel) addRow(result interface{}, rows []table.Row) ([]table.Row, error) {
	filter := true
	if t.FilterIter != nil {
		var err error
		result, filter, err = t.FilterIter(result)
		if err != nil {
			return rows, err
		}
	}
	if !filter {
		return rows, nil
	}

	var columns []table.Column
	record, orderedKeys, err := OnlyFields(t.fields, result)
	if err != nil {
		return rows, err
	}
	rowData := make(table.RowData)
	idResult, ok := result.(files_sdk.Identifier)
	var id interface{}
	if ok {
		id = idResult.Identifier()
	}
	iteratable := true
	iteratableResult, okIteratable := result.(files_sdk.Iterable)
	if okIteratable {
		iteratable = iteratableResult.Iterable()
	}

	for i, key := range orderedKeys {
		cell := fmt.Sprintf("%v", formatValuePretty(key, record[key]))

		if i == 0 && ok {
			rowData[key] = CellWrapper{cell: cell, data: id, Iterable: iteratable}
		} else {
			rowData[key] = cell
		}

		currentWidth := lo.Min[int]([]int{
			lo.Max[int](
				[]int{
					len(cell),
					len(key),
				},
			),
			t.dimensions()[0] / 4,
		})

		if t.maxColumnWidth[fmt.Sprintf("%v", key)] < currentWidth {
			t.maxColumnWidth[fmt.Sprintf("%v", key)] = currentWidth
		}
		columns = append(
			columns,
			table.NewColumn(
				fmt.Sprintf("%v", key),
				fmt.Sprintf("%v", key),
				t.maxColumnWidth[fmt.Sprintf("%v", key)],
			).WithFiltered(true),
		)
	}

	row := table.NewRow(rowData)
	rows = append(rows, row)
	t.updateTable(columns, rows)
	return rows, nil
}

func (t *tableModel) updateTable(columns []table.Column, rows []table.Row) {
	t.Once.Do(func() {
		t.buildTable(columns)
	})

	t.Model = t.Model.
		WithColumns(columns).
		WithRows(rows)
}

func (t *tableModel) updateFooter() {
	footer := fmt.Sprintf("Page %v of %v | Total: %v%v%v", t.Model.CurrentPage(), t.Model.MaxPages(), t.Model.TotalRows(), t.Loading(), t.currentResource())
	if t.Err() != nil {
		footer = fmt.Sprintf("%v - Error: %v", footer, t.Err().Error())
	}
	t.Model = t.Model.
		WithStaticFooter(footer)
}

func (t *tableModel) currentResource() interface{} {
	var resource interface{}
	if len(t.parentResources) > 0 {
		resource = t.parentResources[len(t.parentResources)-1]
	}
	if resource == nil {
		resource = ""
	}
	return resource
}

func (t *tableModel) buildTable(columns []table.Column) {
	t.Model = table.New(columns).
		WithKeyMap(table.KeyMap{
			RowDown: key.NewBinding(
				key.WithKeys("down", "j"),
			),
			RowUp: key.NewBinding(
				key.WithKeys("up", "k"),
			),
			RowSelectToggle: key.NewBinding(
				key.WithKeys(" ", "enter"),
			),
			PageDown: key.NewBinding(
				key.WithKeys("l", "pgdown"),
			),
			PageUp: key.NewBinding(
				key.WithKeys("h", "pgup"),
			),
			PageFirst: key.NewBinding(
				key.WithKeys("home", "g"),
			),
			PageLast: key.NewBinding(
				key.WithKeys("end", "G"),
			),
			Filter: key.NewBinding(
				key.WithKeys("/"),
			),
			FilterBlur: key.NewBinding(
				key.WithKeys("enter", "esc"),
			),
			FilterClear: key.NewBinding(
				key.WithKeys("esc"),
			),
			ScrollRight: key.NewBinding(
				key.WithKeys("right"),
			),
			ScrollLeft: key.NewBinding(
				key.WithKeys("left"),
			),
		},
		).
		Focused(true).
		WithMaxTotalWidth(t.dimensions()[0]).
		WithHorizontalFreezeColumnCount(1).
		WithHeaderVisibility(true).
		WithPaginationWrapping(true).
		WithPageSize(t.dimensions()[1]).
		WithBaseStyle(lipgloss.NewStyle().Align(lipgloss.Left)).
		Filtered(true).
		HeaderStyle(lipgloss.NewStyle().Bold(true)).
		Border(tableBoarder)
	t.Model = t.tableLoader.TableInit(t, t.Model)
}

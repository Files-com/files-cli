package lib

import (
	"context"
	"fmt"
	"strings"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
	"github.com/samber/lo"
)

type tableResource struct {
	resource interface{}
	context  context.Context
	rows     []table.Row
	error
	spinner    TitledSpinner
	title      string
	backLoader tableLoader
}

func (t *tableResource) Init(ctx context.Context, resource interface{}, title interface{}) *tableResource {
	t.spinner = TitledSpinner{spinner.New(spinner.WithSpinner(spinner.Points)), title}
	t.context = ctx
	t.resource = resource
	return t
}

func (t *tableResource) Load() {

}

func (t *tableResource) LoadFirstPage(model *tableModel) error {
	t.error = t.addRow(model, t.resource)
	if t.error != nil {
		return t.error
	}
	return nil
}

func (t *tableResource) addRow(model *tableModel, result interface{}) error {
	var columns []table.Column
	record, orderedKeys, err := OnlyFields(model.fields, result)
	if err != nil {
		return err
	}

	maxColumnLen := lo.Max[int](lo.Map[string, int](orderedKeys, func(item string, index int) int {
		return len(item)
	}))

	var resourceTitle string

	resourceWithPackage := strings.Split(fmt.Sprintf("%T", result), ".")
	if len(resourceWithPackage) > 1 {
		resourceTitle = resourceWithPackage[1]
	} else {
		resourceTitle = "" // might just be a map
	}

	columnsTitle := resourceTitle
	valuesTitle := ""

	if maxColumnLen < len(resourceTitle) {
		columnsTitle = ""
		valuesTitle = resourceTitle
	}

	columns = append(
		columns,
		table.NewColumn(
			"column",
			columnsTitle,
			maxColumnLen,
		).WithFiltered(true).WithStyle(lipgloss.NewStyle().Bold(true)),
	)

	columns = append(
		columns,
		table.NewFlexColumn(
			"value",
			valuesTitle,
			1,
		).WithStyle(lipgloss.NewStyle().Align(lipgloss.Left)),
	)

	for _, key := range orderedKeys {
		rowData := make(table.RowData)
		rowData["column"] = fmt.Sprintf("%v", key)
		rowData["value"] = fmt.Sprintf("%v", formatValuePretty(key, record[key]))
		t.rows = append(t.rows, table.NewRow(rowData))
	}

	model.updateTable(columns, t.rows)
	return nil
}

func (t *tableResource) LoadRest(_model *tableModel) {
}

func (t *tableResource) Loading() string {
	return ""
}

func (t *tableResource) Context() context.Context {
	return t.context
}

func (t *tableResource) Cancel() {

}

func (t *tableResource) Spinner() tea.Model {
	return t.spinner
}

func (t *tableResource) Err() error {
	return nil
}

func (t *tableResource) Update(model *tableModel, msg tea.Msg) (tableLoader, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "c":
			var cmd tea.Cmd
			row := model.HighlightedRow()
			text := fmt.Sprintf("%v", row.Data["value"])

			err := clipboard.WriteAll(text)
			if err != nil {
				return t, tea.Batch(
					tea.Printf(err.Error()),
				)
			}
			return t, cmd
		case "j":
			fmt.Println("")
			err := Format(t.Context(), t.resource, []string{"json"}, model.fields, false)
			if err != nil {
				return t, tea.Batch(
					tea.Printf(err.Error()),
					tea.Quit,
				)
			}
			return t, tea.Quit
		case "backspace":
			return t.backLoader.Update(model, msg)
		}
	}
	return t, cmd
}

func (t *tableResource) ResourceIterator() (files_sdk.ResourceIterator, bool) {
	return nil, false
}

func (t *tableResource) ControlsFooter() string {
	base := "copy [c] json [j]"
	if t.backLoader != nil {
		base += " back [backspace]"
	}
	return base
}

func (t *tableResource) TableInit(model *tableModel, table table.Model) table.Model {
	return table.WithFooterVisibility(false).WithTargetWidth(model.dimensions()[0])
}

func (t *tableResource) SetBackLoader(backLoader tableLoader) {
	t.backLoader = backLoader
}

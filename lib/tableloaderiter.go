package lib

import (
	"context"
	"fmt"
	"strings"

	files_sdk "github.com/Files-com/files-sdk-go/v2"

	"github.com/bradfitz/iter"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/evertras/bubble-table/table"
	"github.com/samber/lo"
)

type tableLoaderIter struct {
	Iter
	rows          []table.Row
	tableRower    chan interface{}
	tableRowerErr chan error
	context.CancelFunc
	context      context.Context
	loadingIndex int
	spinner      tea.Model
	error
}

func (t *tableLoaderIter) Init(ctx context.Context, title interface{}, iter func(context.Context) (Iter, error)) (*tableLoaderIter, error) {
	t.tableRower = make(chan interface{})
	t.tableRowerErr = make(chan error)
	t.spinner = TitledSpinner{spinner.New(spinner.WithSpinner(spinner.Points)), title}
	t.context, t.CancelFunc = context.WithCancel(ctx)
	var err error
	t.Iter, err = iter(t.context)
	return t, err
}

func (t *tableLoaderIter) Context() context.Context {
	return t.context
}

func (t *tableLoaderIter) Spinner() tea.Model {
	return t.spinner
}

func (t *tableLoaderIter) Err() error {
	return t.error
}

func (t *tableLoaderIter) Update(model *tableModel, msg tea.Msg) (tableLoader, tea.Cmd) {
	var cmd tea.Cmd
	if model.TotalRows() == 0 && t.Context().Err() == nil {
		t.spinner, cmd = t.spinner.Update(msg)
		if cmd != nil {
			return t, cmd
		}
	}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			model.Model, cmd = t.OnEnter(model)
			return model.tableLoader, cmd
		case "backspace":
			model.Model, cmd = t.OnBackspace(model)
			return model.tableLoader, cmd
		}
	}

	return t, cmd
}

func (t *tableLoaderIter) Cancel() {
	t.CancelFunc()
}

func (t *tableLoaderIter) ResourceIterator() (files_sdk.ResourceIterator, bool) {
	it, ok := t.Iter.(files_sdk.ResourceIterator)
	return it, ok
}

func (t *tableLoaderIter) ResourceLoader() (files_sdk.ResourceLoader, bool) {
	it, ok := t.Iter.(files_sdk.ResourceLoader)
	return it, ok
}

func (t *tableLoaderIter) ReloadIterator() (files_sdk.ReloadIterator, bool) {
	it, ok := t.Iter.(files_sdk.ReloadIterator)
	return it, ok
}

func (t *tableLoaderIter) Load() {
	go func() {
		defer close(t.tableRower)
		defer t.CancelFunc()
		for t.Iter.Next() {
			if t.context.Err() != nil {
				return
			}
			if t.Iter.Err() != nil {
				t.tableRowerErr <- t.Iter.Err()
			}
			if t.Iter.Current() == nil {
				return
			}
			t.tableRower <- t.Iter.Current()
		}
		if t.Iter.Err() != nil {
			t.tableRowerErr <- t.Iter.Err()
			return
		}
	}()
}

func (t *tableLoaderIter) Loading() string {
	if t.loadingIndex == len(Loading) {
		t.loadingIndex = 0
	}
	if t.context.Err() != nil {
		return " "
	}
	defer func() { t.loadingIndex += 1 }()
	return fmt.Sprintf(" %v ", Loading[t.loadingIndex])
}

func (t *tableLoaderIter) LoadFirstPage(model *tableModel) error {
	for range iter.N(model.dimensions()[1]) {
		select {
		case result := <-t.tableRower:
			var err error
			if result == nil {
				return nil
			}
			t.rows, err = model.addRow(result, t.rows)
			if err != nil {
				return err
			}
		case err := <-t.tableRowerErr:
			t.error = err
		case <-t.context.Done():
			return nil
		}
	}

	return nil
}

func (t *tableLoaderIter) LoadRest(model *tableModel) {
	go func() {
		for {
			select {
			case result := <-t.tableRower:
				var err error
				if result == nil {
					return
				}

				err = t.addRow(model, result)
				if err != nil {
					t.CancelFunc()
					return
				}
				if model.Program != nil {
					model.Program.Send(nil)
				}
			case <-t.tableRowerErr:
				return
			case <-t.context.Done():
				return
			}
		}
	}()
}

func (t *tableLoaderIter) addRow(model *tableModel, result interface{}) error {
	filter := true
	if model.FilterIter != nil {
		var err error
		result, filter, err = model.FilterIter(result)
		if err != nil {
			return err
		}
	}

	if !filter {
		return nil
	}

	var columns []table.Column
	record, orderedKeys, err := OnlyFields(model.fields, result)
	if err != nil {
		return err
	}
	rowData := make(table.RowData)
	idResult, okId := result.(files_sdk.Identifier)
	var id interface{}
	if okId {
		id = idResult.Identifier()
	}

	iteratable := true
	itResult, okIt := result.(files_sdk.Iterable)
	if okIt {
		iteratable = itResult.Iterable()
	}

	for i, key := range orderedKeys {
		cell := fmt.Sprintf("%v", formatValuePretty(key, record[key]))

		if i == 0 && okId {
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
			model.dimensions()[0] / 4,
		})

		if model.maxColumnWidth[fmt.Sprintf("%v", key)] < currentWidth {
			model.maxColumnWidth[fmt.Sprintf("%v", key)] = currentWidth
		}
		columns = append(
			columns,
			table.NewColumn(
				fmt.Sprintf("%v", key),
				fmt.Sprintf("%v", key),
				model.maxColumnWidth[fmt.Sprintf("%v", key)],
			).WithFiltered(true),
		)
	}

	row := table.NewRow(rowData)
	t.rows = append(t.rows, row)
	model.updateTable(columns, t.rows)
	return nil
}

func (t *tableLoaderIter) ControlsFooter() string {
	var footer []string
	footer = append(footer, "h-scroll [←][→]")
	if _, ok := t.ResourceIterator(); ok {
		footer = append(footer, "back [backspace]")

	}

	if t.Context().Err() == nil {
		footer = append(footer, " cancel [esc]")
	}
	return strings.Join(footer, " ")
}

func (t *tableLoaderIter) TableInit(_ *tableModel, table table.Model) table.Model {
	return table
}

func (t *tableLoaderIter) OnEnter(model *tableModel) (table.Model, tea.Cmd) {
	var cmd tea.Cmd

	row := model.HighlightedRow()
	var id interface{}
	iteratable := true
	for _, v := range row.Data {
		cell, cellOk := v.(CellWrapper)
		if cellOk {
			id = cell.data
			iteratable = cell.Iterable
			break
		}
	}
	if id == nil {
		return model.Model, cmd
	}

	itLoader, ok := t.ResourceIterator()
	if ok && iteratable {
		t.Cancel()

		model.parentResources = append(model.parentResources, id)

		loader, err := (&tableLoaderIter{}).Init(context.Background(), id, func(ctx context.Context) (Iter, error) {
			return itLoader.Iterate(id, files_sdk.WithContext(ctx))
		})
		if err != nil {
			return model.Model, tea.Batch(
				tea.Printf(err.Error()),
			)
		}
		model.SetLoader(loader)
		go model.Load()
		model.tableLoader, cmd = loader.Update(model, loader.spinner.(TitledSpinner).Tick())
		return model.Model, cmd
	}

	resourceLoader, ok := t.ResourceLoader()
	if ok {
		t.Cancel()

		model.parentResources = append(model.parentResources, id)
		it, err := resourceLoader.LoadResource(id)
		if err != nil {
			return model.Model, tea.Batch(
				tea.Printf(err.Error()),
			)
		}

		loader := (&tableResource{}).Init(context.Background(), it, id)
		loader.SetBackLoader(t)
		model.SetLoader(loader)
		model.Load()
		return model.Model, cmd
	}

	return model.Model, cmd
}

func (t *tableLoaderIter) OnBackspace(model *tableModel) (table.Model, tea.Cmd) {
	var cmd tea.Cmd
	itLoader, ok := t.ResourceIterator()
	if ok {
		t.Cancel()
		var parent interface{}
		if len(model.parentResources) >= 2 {
			parent = model.parentResources[len(model.parentResources)-2]
		}
		if len(model.parentResources) > 0 {
			model.parentResources = model.parentResources[0 : len(model.parentResources)-1]
		}

		loader, err := (&tableLoaderIter{}).Init(context.Background(), parent, func(ctx context.Context) (Iter, error) {
			return itLoader.Iterate(parent, files_sdk.WithContext(ctx))
		})
		if err != nil {
			return model.Model, tea.Batch(
				tea.Printf(err.Error()),
			)
		}
		model.SetLoader(loader)
		go model.Load()
		model.tableLoader, cmd = loader.Update(model, loader.spinner.(TitledSpinner).Tick())
		return model.Model, cmd
	}

	itReload, ok := t.ReloadIterator()
	if ok {
		t.Cancel()
		var parent interface{}
		if len(model.parentResources) >= 2 {
			parent = model.parentResources[len(model.parentResources)-2]
		}
		if len(model.parentResources) > 0 {
			model.parentResources = model.parentResources[0 : len(model.parentResources)-1]
		}

		loader, err := (&tableLoaderIter{}).Init(context.Background(), parent, func(ctx context.Context) (Iter, error) {
			return itReload.Reload(files_sdk.WithContext(ctx)), nil
		})
		if err != nil {
			return model.Model, tea.Batch(
				tea.Printf(err.Error()),
			)
		}
		model.SetLoader(loader)
		go model.Load()
		model.tableLoader, cmd = loader.Update(model, loader.spinner.(TitledSpinner).Tick())
		return model.Model, cmd
	}
	return model.Model, cmd
}

func (t *tableLoaderIter) SetBackLoader(tableLoader) {
	//	noop
}

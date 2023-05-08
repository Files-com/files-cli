package lib

import (
	"context"

	files_sdk "github.com/Files-com/files-sdk-go/v2"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/evertras/bubble-table/table"
)

type tableLoader interface {
	Load()
	LoadFirstPage(model *tableModel) error
	LoadRest(model *tableModel)
	Loading() string
	Context() context.Context
	Cancel()
	Spinner() tea.Model
	Err() error
	ResourceIterator() (files_sdk.ResourceIterator, bool)
	ControlsFooter() string
	TableInit(model *tableModel, table table.Model) table.Model
	Update(model *tableModel, msg tea.Msg) (tableLoader, tea.Cmd)
	SetBackLoader(tableLoader)
}

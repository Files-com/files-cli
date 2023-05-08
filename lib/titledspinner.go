package lib

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type TitledSpinner struct {
	spinner.Model
	Title interface{}
}

func (t TitledSpinner) View() string {
	if t.Title == nil {
		return fmt.Sprintf("%v", t.Model.View())
	} else {
		return fmt.Sprintf("%v %v", t.Model.View(), t.Title)
	}
}

func (t TitledSpinner) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	t.Model, cmd = t.Model.Update(msg)
	return t, cmd
}

func (t TitledSpinner) Init() tea.Cmd {
	return nil
}

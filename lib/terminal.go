package lib

import (
	"fmt"

	//tea "test3/go/pkg/mod/github.com/charmbracelet/bubbletea@v0.20.0"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type errMsg error

type PretextInputModel struct {
	textInput textinput.Model
	display   string
	err       error
}

func (p *PretextInputModel) new(display string, pretext string, charLimit int, width int) {
	p.err = nil
	p.display = display
	p.textInput = textinput.New()
	p.textInput.SetValue(pretext)
	p.textInput.CharLimit = charLimit
	p.textInput.Width = width
	p.textInput.Focus()
}

func (p PretextInputModel) Init() tea.Cmd {
	return textinput.Blink
}

func (p PretextInputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			p.textInput, cmd = p.textInput.Update(msg)
			return p, tea.Quit
		case tea.KeyCtrlC, tea.KeyEsc:
			return p, tea.Quit
		}
	case errMsg:
		p.err = msg
		return p, nil
	}
	p.textInput, cmd = p.textInput.Update(msg)
	return p, cmd
}

func (p PretextInputModel) View() string {
	return fmt.Sprintf(p.display, p.textInput.View()) + "\n"
}

func PromptUserWithPretext(display string, pretext string, config Config) (string, error) {
	var pretextModel PretextInputModel
	maxInputLength := 156
	var maxAdditionalCharViewSpace int
	if pretext == "" {
		maxAdditionalCharViewSpace = 40
	} else {
		maxAdditionalCharViewSpace = 10
	}
	pretextModel.new(display, pretext, maxInputLength, len(pretext)+maxAdditionalCharViewSpace)
	p := tea.NewProgram(pretextModel, tea.WithOutput(config.Out), tea.WithInput(config.In))
	newModel, err := p.StartReturningModel()
	return newModel.(PretextInputModel).textInput.Value(), err
}

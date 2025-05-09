package lib

import (
	"context"
	"errors"
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type PretextInputModel struct {
	textInput textinput.Model
	display   string
	err       error
	example   string
}

func (p *PretextInputModel) new(display string, pretext string, charLimit int, width int, example string) {
	p.err = nil
	p.display = display
	p.textInput = textinput.New()
	p.textInput.SetValue(pretext)
	p.textInput.CharLimit = charLimit
	p.textInput.Width = width
	p.textInput.Focus()
	p.example = example
}

func (p *PretextInputModel) Init() tea.Cmd {
	return textinput.Blink
}

func (p *PretextInputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			p.textInput, _ = p.textInput.Update(msg)
			return p, tea.Quit
		case tea.KeyCtrlC, tea.KeyEsc:
			// this runs in a subprocess, so there's no need to
			// use the clierr package
			p.err = errors.New("exited session logging")
			return p, tea.Quit
		}
	case error:
		p.err = msg
		return p, nil
	}
	p.textInput, cmd = p.textInput.Update(msg)
	return p, cmd
}

func (p *PretextInputModel) View() string {
	if p.example == "" {
		return fmt.Sprintf(p.display+"\n", p.textInput.View())
	}
	lightGrey := "\033[37m" // Light grey color code
	reset := "\033[0m"      // Reset code
	exampleText := fmt.Sprintf("%s%s%s", lightGrey, p.example, reset)
	return fmt.Sprintf("%v\n"+p.display+"\n", exampleText, p.textInput.View())
}

func PromptUserWithPretext(ctx context.Context, display string, pretext string, example string, profile *Profiles) (string, error) {
	pretextModel := &PretextInputModel{}
	maxInputLength := 156
	var maxAdditionalCharViewSpace int
	if pretext == "" {
		maxAdditionalCharViewSpace = 40
	} else {
		maxAdditionalCharViewSpace = 10
	}
	pretextModel.new(display, pretext, maxInputLength, len(pretext)+maxAdditionalCharViewSpace, example)
	p := tea.NewProgram(pretextModel, tea.WithOutput(profile.Out), tea.WithInput(profile.In), tea.WithContext(ctx))
	_, err := p.Run()
	if err == nil {
		err = pretextModel.err
	}
	return pretextModel.textInput.Value(), err
}

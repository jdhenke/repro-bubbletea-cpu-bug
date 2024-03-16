package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"log"
	"os"
)

func main() {
	ti := textinput.New()
	ti.Placeholder = "Placeholder..."
	ti.Focus()
	ti.Width = 50
	if _, err := tea.NewProgram(model{
		input: ti,
	}).Run(); err != nil {
		log.Fatal(err)
	}
}

var _ tea.Model = model{}

type model struct {
	input textinput.Model
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyEnter:
			text := m.input.Value()
			cmds = append(cmds, func() tea.Msg {
				// for illustrative purposes... imagine something that should be done off the UI thread
				_ = os.WriteFile("last.txt", []byte(text), 0644)
				return nil
			})
		}
	}
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Sequence(cmds...) // 👈 BUG: causes an infinite loop spiking CPU; but using tea.Batch does NOT
}

func (m model) View() string {
	return fmt.Sprintf("Check CPU usage\n\n%s", m.input.View())
}

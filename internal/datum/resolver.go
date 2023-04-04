package datum

import (
	"fmt"
	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/davfer/capi/internal/store"
	"github.com/davfer/capi/pkg/model"
	"github.com/pkg/errors"
	"os"
	"strings"
)

type Resolver interface {
	Resolve(*DataContext, *model.Collection) error
}

type FallbackResolver struct {
	resolver Resolver
	fallback Resolver
}

type DbResolver struct {
	store    store.Store
	fallback Resolver
}

func NewDbResolver(db store.Store, chain Resolver) *DbResolver {
	return &DbResolver{
		store:    db,
		fallback: chain,
	}
}

func (r *DbResolver) Resolve(data *DataContext, c *model.Collection) error {
	if data.IsComplete() {
		return nil
	}

	var missing []DataKey
	for _, k := range data.GetMissingData() {
		if k.Typ == "var" {
			got, err := r.store.Get(k.Col, k.Name)
			if err != nil {
				if errors.Is(err, store.ErrNotFoundKey) && r.fallback != nil {
					missing = append(missing, k)
					continue
				}

				return err
			}

			data.Set(k, got)
		}
	}

	if !data.IsComplete() && r.fallback != nil {
		err := r.fallback.Resolve(data, c)
		if err != nil {
			return err
		}

		for _, k := range missing {
			if v, err := data.Get(k); err == nil {
				if err := r.store.Set(k.Col, k.Name, v); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

type AskResolver struct {
}

func NewAskResolver() *AskResolver {
	return &AskResolver{}
}

func (r *AskResolver) Resolve(data *DataContext, c *model.Collection) error {
	if data.IsComplete() {
		return nil
	}

	m := initialModel(data.GetMissingData())
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Printf("could not start program: %s\n", err)
		os.Exit(1)
	}

	for _, v := range m.inputs {
		for _, k := range data.GetMissingData() {
			if k.Name == v.Placeholder {
				data.Set(k, v.Value())
			}
		}
	}

	return nil
}

var (
	focusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle         = focusedStyle.Copy()
	noStyle             = lipgloss.NewStyle()
	helpStyle           = blurredStyle.Copy()
	cursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))

	focusedButton = focusedStyle.Copy().Render("[ Submit ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
)

type askModel struct {
	focusIndex int
	inputs     []textinput.Model
	cursorMode cursor.Mode
}

func initialModel(keys []DataKey) *askModel {
	m := askModel{
		inputs: make([]textinput.Model, len(keys)),
	}

	var t textinput.Model
	for i, k := range keys {
		t = textinput.New()
		t.CursorStyle = cursorStyle
		t.CharLimit = 32

		t.Placeholder = k.Name
		t.PromptStyle = focusedStyle
		t.TextStyle = focusedStyle

		if i == 0 {
			t.Focus()
		}
		m.inputs[i] = t
	}

	return &m
}

func (m *askModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m *askModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit

		// Change cursor mode
		case "ctrl+r":
			m.cursorMode++
			if m.cursorMode > cursor.CursorHide {
				m.cursorMode = cursor.CursorBlink
			}
			cmds := make([]tea.Cmd, len(m.inputs))
			for i := range m.inputs {
				cmds[i] = m.inputs[i].Cursor.SetMode(m.cursorMode)
			}
			return m, tea.Batch(cmds...)

		// Set focus to next input
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			// Did the user press enter while the submit button was focused?
			// If so, exit.
			if s == "enter" && m.focusIndex == len(m.inputs) {
				return m, tea.Quit
			}

			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > len(m.inputs) {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.inputs)
			}

			cmds := make([]tea.Cmd, len(m.inputs))
			for i := 0; i <= len(m.inputs)-1; i++ {
				if i == m.focusIndex {
					// Set focused state
					cmds[i] = m.inputs[i].Focus()
					m.inputs[i].PromptStyle = focusedStyle
					m.inputs[i].TextStyle = focusedStyle
					continue
				}
				// Remove focused state
				m.inputs[i].Blur()
				m.inputs[i].PromptStyle = noStyle
				m.inputs[i].TextStyle = noStyle
			}

			return m, tea.Batch(cmds...)
		}
	}

	// Handle character input and blinking
	cmd := m.updateInputs(msg)

	return m, cmd
}

func (m *askModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m *askModel) View() string {
	var b strings.Builder

	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		if i < len(m.inputs)-1 {
			b.WriteRune('\n')
		}
	}

	button := &blurredButton
	if m.focusIndex == len(m.inputs) {
		button = &focusedButton
	}
	fmt.Fprintf(&b, "\n\n%s\n\n", *button)

	b.WriteString(helpStyle.Render("cursor mode is "))
	b.WriteString(cursorModeHelpStyle.Render(m.cursorMode.String()))
	b.WriteString(helpStyle.Render(" (ctrl+r to change style)"))

	return b.String()
}

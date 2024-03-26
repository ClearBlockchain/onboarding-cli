package ui

import (
	"fmt"
	"strings"

	"github.com/ClearBlockchain/onboarding-cli/pkg/utils"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

type ModelState int

const (
	StatusNormal ModelState = iota
	StatusDone
)

type Model struct {
	Status ModelState
	Renderer *lipgloss.Renderer
	Form *huh.Form
	Width int
}

func NewModel(form *huh.Form) Model {
	m := Model{Width: DefaultWidth}
	m.Renderer = lipgloss.DefaultRenderer()
	m.Form = form.
		WithWidth(DefaultWidth - 70).
		WithShowHelp(false).
		WithShowErrors(false)

	return m
}

func (m Model) Init() tea.Cmd {
	return m.Form.Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = utils.Min(msg.Width, DefaultWidth) - Base.GetHorizontalFrameSize()
	case tea.KeyMsg:
		switch msg.String() {
			case "esc", "ctrl+c", "q":
				return m, tea.Quit
		}
	}

	var cmds []tea.Cmd

	// process the form
	form, cmd := m.Form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.Form = f
		cmds = append(cmds, cmd)
	}

	if m.Form.State == huh.StateCompleted {
		// quit when the form is done
		cmds = append(cmds, tea.Quit)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	switch m.Form.State {
	case huh.StateCompleted:
		var b strings.Builder
		fmt.Fprintf(&b, "Welcome to %s! We'll take it from here 🤗\n\n", Highlight.Render("ClearX OGI"))
		fmt.Fprintf(&b, "Please follow the instructions on the browser to complete the setup.")
		return Status.Copy().Margin(0, 1).Padding(1, 2).Width(48).Render(b.String()) + "\n\n"
	default:

		var endpoints string
		// endpoints is array of strings
		if m.Form.Get("endpoints") != nil {
			items := m.Form.Get("endpoints").([]string)
			highlightedItems := make([]string, len(items))
			for i, item := range items {
				highlightedItems[i] = Highlight.Render(item)
			}

			// cast to []string and join with new line and -
			endpoints = fmt.Sprintf("Endpoints:\n- %s\n", strings.Join(highlightedItems, "\n- "))
		}

		// Form (left side)
		v := strings.TrimSuffix(m.Form.View(), "\n\n")
		form := m.Renderer.NewStyle().Margin(1, 0).Render(v)

		// Status (right side)
		var status string
		{
			var (
				buildInfo      = "(None)"
				gcpProject string
			)

			if m.Form.Get("endpoints") != nil {
				buildInfo = endpoints
			}

			if m.Form.GetString("gcpProject") != "" {
				gcpProject = fmt.Sprintf("\nGCP Project:\n- %s\n", Highlight.Render(m.Form.GetString("gcpProject")))
			}

			const statusWidth = 50
			statusMarginLeft := m.Width - statusWidth - lipgloss.Width(form) - Status.GetMarginRight()
			status = Status.Copy().
				Height(lipgloss.Height(form)).
				Width(statusWidth).
				MarginLeft(statusMarginLeft).
				Render(StatusHeader.Render("Your ClearX OGI Setup") + "\n" +
					buildInfo +
					gcpProject)
		}

		errors := m.Form.Errors()
		header := m.appBoundaryView("ClearX OGI Project Initialization")
		if len(errors) > 0 {
			header = m.appErrorBoundaryView(m.errorView())
		}
		body := lipgloss.JoinHorizontal(lipgloss.Top, form, status)

		footer := m.appBoundaryView(m.Form.Help().ShortHelpView(m.Form.KeyBinds()))
		if len(errors) > 0 {
			footer = m.appErrorBoundaryView("")
		}

		return Base.Render(header + "\n" + body + "\n\n" + footer)
	}
}

func (m Model) errorView() string {
	var s string
	for _, err := range m.Form.Errors() {
		s += err.Error()
	}
	return s
}

func (m Model) appBoundaryView(text string) string {
	return lipgloss.PlaceHorizontal(
		m.Width,
		lipgloss.Left,
		HeaderText.Render(text),
		lipgloss.WithWhitespaceChars("/"),
		lipgloss.WithWhitespaceForeground(Indigo),
	)
}

func (m Model) appErrorBoundaryView(text string) string {
	return lipgloss.PlaceHorizontal(
		m.Width,
		lipgloss.Left,
		ErrorHeaderText.Render(text),
		lipgloss.WithWhitespaceChars("/"),
		lipgloss.WithWhitespaceForeground(Red),
	)
}


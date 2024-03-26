package ui

import (
	_ "github.com/charmbracelet/bubbletea"
	_ "github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

var (
	Keyword = lipgloss.NewStyle().
		Foreground(Brand)

	Underline = lipgloss.NewStyle().
		Underline(true)

	Paragraph = lipgloss.NewStyle().
		Width(DefaultWidth).
		Padding(0, 0, 0, 2)

	Title = lipgloss.NewStyle().
		Bold(true).
		Foreground(Brand)

	Box = lipgloss.NewStyle().
		Border(lipgloss.Border{
			Top:         "─",
			Bottom:      "─",
			Left:        "│",
			Right:       "│",
			TopLeft:     "╭",
			TopRight:    "╮",
			BottomRight:  "┘",
			BottomLeft: "└",
		}).
		BorderForeground(Gray).
		Padding(1, 2).
		Width(DefaultWidth).
		Align(lipgloss.Center)

	Base = lipgloss.NewStyle().
		Padding(1, 4, 0, 1)

	HeaderText = lipgloss.NewStyle().
		Foreground(Indigo).
		Bold(true).
		Padding(0, 1, 0, 2)

	Status = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(Indigo).
		PaddingLeft(1).
		MarginTop(1)

	StatusHeader = lipgloss.NewStyle().
		Foreground(Green).
		Bold(true)

	Highlight = lipgloss.NewStyle().
		Foreground(lipgloss.Color("212"))

	ErrorHeaderText = HeaderText.
		Copy().
		Foreground(Red)

	Help = lipgloss.NewStyle().
		Foreground(lipgloss.Color("240"))
)

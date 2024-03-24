package ui

import (
	_ "github.com/charmbracelet/bubbletea"
	_ "github.com/charmbracelet/huh"

	. "github.com/charmbracelet/lipgloss"
)

var (
	Keyword = NewStyle().
		Foreground(Color("#304ffe")).
		Render

	Underline = NewStyle().
		Underline(true).
		Render

	Paragraph = NewStyle().
		Width(78).
		Padding(0, 0, 0, 2).
		Render

	Title = NewStyle().
		Bold(true).
		Foreground(Color("#304ffe")).
		Render

	Box = NewStyle().
		Border(Border{
			Top:         "─",
			Bottom:      "─",
			Left:        "│",
			Right:       "│",
			TopLeft:     "╭",
			TopRight:    "╮",
			BottomRight:  "┘",
			BottomLeft: "└",
		}).
		BorderForeground(AdaptiveColor{
			Light: "#f0f0f0",
			Dark:  "#333333",
		}).
		Padding(1, 2).
		Width(78).
		Align(Center).
		Render
)

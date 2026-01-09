package tui

import "github.com/charmbracelet/lipgloss"

var (
	BannerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("205")).
			Bold(true)

	VersionStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("245"))

	HeaderStyle          = lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
	CPUHighStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true)
	CPUMediumStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("226"))
	CPULowStyle          = lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
	MemoryStyle          = lipgloss.NewStyle().Foreground(lipgloss.Color("51"))
	PIDStyle             = lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
	NameStyle            = lipgloss.NewStyle().Foreground(lipgloss.Color("255"))
	CursorIndicatorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Bold(true)
	CheckedStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
	CheckmarkStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("42")).Bold(true)
	RowHighlightStyle    = lipgloss.NewStyle()
	SelectedCountStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("42")).Bold(true)
	HelpStyle            = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	FilterStyle          = lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Bold(true)
	CursorStyle          = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	DimStyle             = lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
	SuccessStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
	ErrorStyle           = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
)

package tui

import "github.com/charmbracelet/lipgloss"

var (
	// Primary: Cyan accent for focus states
	// Secondary: Magenta for highlights
	// Neutrals: Charcoal and slate grays for hierarchy
	primaryColor   = lipgloss.Color("#00D9FF") // Bright Cyan
	secondaryColor = lipgloss.Color("#FF006E") // Vibrant Magenta
	accentColor    = lipgloss.Color("#8338EC") // Purple accent
	textColor      = lipgloss.Color("#F0F0F0") // Off-white
	subtleColor    = lipgloss.Color("#78787A") // Medium gray
	mutedColor     = lipgloss.Color("#3C3C3C") // Dark gray
	borderColor    = lipgloss.Color("#2A2A2A") // Charcoal
	errorColor     = lipgloss.Color("#FF5555") // Softer red

	appStyle = lipgloss.NewStyle().Padding(1, 3)

	listTitleStyle = lipgloss.NewStyle().
			Foreground(textColor).
			Background(lipgloss.Color("#1A1A1A")).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(primaryColor).
			BorderBottom(true).
			Padding(1, 2).
			Bold(true).
			MarginBottom(1)

	listStatusStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true)

	formTitleStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true).
			MarginBottom(2).
			PaddingBottom(1)

	focusedPromptStyle = lipgloss.NewStyle().
				Foreground(primaryColor).
				Bold(true)

	blurredPromptStyle = lipgloss.NewStyle().
				Foreground(subtleColor)

	inputStyle = lipgloss.NewStyle().
			Foreground(textColor)

	focusedButtonStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#0A0A0A")).
				Background(primaryColor).
				Padding(0, 4).
				Bold(true).
				BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(primaryColor)

	blurredButtonStyle = lipgloss.NewStyle().
				Foreground(subtleColor).
				BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(borderColor).
				Padding(0, 4)

	errorMessageStyle = lipgloss.NewStyle().
				Foreground(errorColor).
				Bold(true).
				Padding(1, 2).
				Border(lipgloss.RoundedBorder()).
				BorderForeground(errorColor).
				Background(lipgloss.Color("#2A1A1A"))

	statusMessageStyle = lipgloss.NewStyle().
				Foreground(primaryColor).
				Bold(true)

	helpStyle = lipgloss.NewStyle().
			Foreground(subtleColor).
			MarginTop(2).
			Italic(true)
)

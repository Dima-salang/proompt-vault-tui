package tui

import "github.com/charmbracelet/lipgloss"

var (
	// Modern color palette - no backgrounds, works with any terminal theme
	primaryColor   = lipgloss.Color("#00E5FF") // Electric Cyan
	secondaryColor = lipgloss.Color("#B388FF") // Soft Purple
	accentColor    = lipgloss.Color("#69F0AE") // Mint Green
	dangerColor    = lipgloss.Color("#FF5370") // Coral Red
	textColor      = lipgloss.Color("#E4E4E7") // Cool white
	subtleColor    = lipgloss.Color("#94A3B8") // Slate
	mutedColor     = lipgloss.Color("#64748B") // Medium slate
	borderColor    = lipgloss.Color("#475569") // Visible border

	appStyle = lipgloss.NewStyle().
			Padding(1, 2)

	listTitleStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			BorderStyle(lipgloss.ThickBorder()).
			BorderForeground(primaryColor).
			BorderBottom(true).
			Padding(0, 1, 1, 1).
			Bold(true).
			MarginBottom(1)

	listStatusStyle = lipgloss.NewStyle().
			Foreground(accentColor).
			Bold(true)

	formTitleStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true).
			MarginBottom(1).
			PaddingBottom(1).
			PaddingLeft(1).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(primaryColor).
			BorderLeft(true)

	focusedPromptStyle = lipgloss.NewStyle().
				Foreground(primaryColor).
				Bold(true)

	blurredPromptStyle = lipgloss.NewStyle().
				Foreground(mutedColor)

	inputStyle = lipgloss.NewStyle().
			Foreground(textColor)

	focusedButtonStyle = lipgloss.NewStyle().
				Foreground(primaryColor).
				Padding(0, 3).
				Bold(true).
				BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(primaryColor)

	blurredButtonStyle = lipgloss.NewStyle().
				Foreground(mutedColor).
				BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(borderColor).
				Padding(0, 3)

	errorMessageStyle = lipgloss.NewStyle().
				Foreground(dangerColor).
				Bold(true).
				Padding(1, 2).
				Border(lipgloss.RoundedBorder()).
				BorderForeground(dangerColor)

	statusMessageStyle = lipgloss.NewStyle().
				Foreground(accentColor).
				Bold(true)

	helpStyle = lipgloss.NewStyle().
			Foreground(subtleColor).
			MarginTop(1).
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(borderColor)
)
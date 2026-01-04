package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/Dima-salang/proompt-vault-tui/internal/vault"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type sessionState int

const (
	stateList sessionState = iota
	stateCreate
	stateDeleteConfirm
)

type Model struct {
	state   sessionState
	service vault.PromptService
	list    list.Model

	// Form inputs
	titleInput       textinput.Model
	descriptionInput textinput.Model
	contentInput     textarea.Model
	focusIndex       int

	err    error
	width  int
	height int

	activePrompt *vault.Prompt // if nil, we are creating. if not, we are editing.
}

func NewModel(service vault.PromptService) Model {
	// Initialize inputs with clean styling
	ti := textinput.New()
	ti.Placeholder = "Enter prompt title..."
	ti.Focus()
	ti.CharLimit = 50
	ti.Width = 60
	ti.PromptStyle = focusedPromptStyle
	ti.TextStyle = inputStyle

	desc := textinput.New()
	desc.Placeholder = "Brief description..."
	desc.CharLimit = 100
	desc.Width = 60
	desc.TextStyle = inputStyle

	cont := textarea.New()
	cont.Placeholder = "Write your prompt content here..."
	cont.ShowLineNumbers = true
	cont.SetWidth(70)
	cont.SetHeight(12)
	cont.FocusedStyle.Base = lipgloss.NewStyle().Foreground(textColor)

	// Initialize List with clean delegate styling
	items := []list.Item{}
	delegate := list.NewDefaultDelegate()

	// Selected item - no background, just color and border
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.
		Foreground(primaryColor).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(primaryColor).
		BorderLeft(true).
		Padding(0, 0, 0, 2).
		Bold(true)

	delegate.Styles.SelectedDesc = delegate.Styles.SelectedDesc.
		Foreground(subtleColor).
		Padding(0, 0, 0, 2)

	// Normal items
	delegate.Styles.NormalTitle = delegate.Styles.NormalTitle.
		Foreground(textColor).
		Padding(0, 0, 0, 1)

	delegate.Styles.NormalDesc = delegate.Styles.NormalDesc.
		Foreground(mutedColor).
		Padding(0, 0, 0, 1)

	l := list.New(items, delegate, 0, 0)
	l.Title = "Prompt Vault"
	l.Styles.Title = listTitleStyle
	l.Styles.FilterPrompt = lipgloss.NewStyle().Foreground(primaryColor).Bold(true)
	l.Styles.FilterCursor = lipgloss.NewStyle().Foreground(accentColor)

	l.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			key.NewBinding(
				key.WithKeys("a"),
				key.WithHelp("a", "new"),
			),
			key.NewBinding(
				key.WithKeys("e"),
				key.WithHelp("e", "edit"),
			),
			key.NewBinding(
				key.WithKeys("d"),
				key.WithHelp("d", "delete"),
			),
			key.NewBinding(
				key.WithKeys("enter"),
				key.WithHelp("↵", "copy"),
			),
		}
	}
	l.AdditionalFullHelpKeys = l.AdditionalShortHelpKeys

	return Model{
		state:            stateList,
		service:          service,
		list:             l,
		titleInput:       ti,
		descriptionInput: desc,
		contentInput:     cont,
		focusIndex:       0,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.fetchPrompts,
		textinput.Blink,
		tea.EnableMouseCellMotion,
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		h, v := appStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
		return m, nil

	case tea.KeyMsg:
		if m.state == stateList {
			switch msg.String() {
			case "ctrl+c", "q":
				if m.list.FilterState() == list.Filtering {
					break
				}
				return m, tea.Quit
			case "a":
				if m.list.FilterState() == list.Filtering {
					break
				}
				m.state = stateCreate
				m.resetForm()
				return m, nil
			case "e":
				if m.list.FilterState() == list.Filtering {
					break
				}
				if i, ok := m.list.SelectedItem().(item); ok {
					m.state = stateCreate
					m.activePrompt = &i.prompt
					m.setForm(i.prompt)
				}
				return m, nil
			case "enter":
				if m.list.FilterState() == list.Filtering {
					break
				}
				if i, ok := m.list.SelectedItem().(item); ok {
					err := vault.CopyToClipboard(&i.prompt)
					if err != nil {
						m.err = err
						return m, nil
					}
					return m, m.list.NewStatusMessage(statusMessageStyle.Render("✓ Copied to clipboard!"))
				}
				return m, nil
			case "d":
				if m.list.FilterState() == list.Filtering {
					break
				}
				if i, ok := m.list.SelectedItem().(item); ok {
					m.activePrompt = &i.prompt
					m.state = stateDeleteConfirm
				}
				return m, nil
			}
		} else if m.state == stateDeleteConfirm {
			switch msg.String() {
			case "y", "Y", "enter":
				if m.activePrompt != nil {
					return m, m.deletePrompt
				}
			case "n", "N", "esc", "q":
				m.state = stateList
				m.activePrompt = nil
				return m, nil
			}
		} else if m.state == stateCreate {
			switch msg.String() {
			case "ctrl+c":
				return m, tea.Quit
			case "esc":
				m.state = stateList
				return m, nil
			case "tab", "shift+tab", "enter", "up", "down":
				s := msg.String()

				// Prioritize Submit if Enter is pressed on Submit button
				if s == "enter" && m.focusIndex == 3 {
					return m, m.createPrompt
				}

				// If in textarea (content input), enter should add new line unless ctrl+enter or moved away
				if m.focusIndex == 2 {
					if s == "enter" {
						break // Textarea handles enter natively
					}
					if s == "up" || s == "down" {
						break // Textarea handles up/down natively
					}
				}

				// Navigation logic
				if s == "up" || s == "shift+tab" {
					m.focusIndex--
				} else if s == "down" || s == "tab" || (s == "enter" && m.focusIndex != 2) {
					m.focusIndex++
				}

				if m.focusIndex > 3 { // 0: Title, 1: Desc, 2: Content, 3: Submit
					m.focusIndex = 0
				} else if m.focusIndex < 0 {
					m.focusIndex = 3
				}

				// Update focus
				cmds = append(cmds, m.updateFocus())
				return m, tea.Batch(cmds...)
			}
		}

	case promptsMsg:
		items := make([]list.Item, len(msg))
		for i, p := range msg {
			items[i] = item{prompt: p}
		}
		cmds = append(cmds, m.list.SetItems(items))

	case promptCreatedMsg:
		m.state = stateList
		m.resetForm()
		cmds = append(cmds, m.fetchPrompts) // Refresh list

	case promptDeletedMsg:
		m.state = stateList
		m.activePrompt = nil
		cmds = append(cmds, m.fetchPrompts)
		cmds = append(cmds, m.list.NewStatusMessage(statusMessageStyle.Render("✓ Prompt deleted")))

	case errMsg:
		m.err = msg
		return m, nil
	}

	// Update children based on state
	if m.state == stateList {
		m.list, cmd = m.list.Update(msg)
		cmds = append(cmds, cmd)
	} else {
		m.titleInput, cmd = m.titleInput.Update(msg)
		cmds = append(cmds, cmd)
		m.descriptionInput, cmd = m.descriptionInput.Update(msg)
		cmds = append(cmds, cmd)
		m.contentInput, cmd = m.contentInput.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.err != nil {
		return errorMessageStyle.Render("⚠ Error: " + m.err.Error())
	}

	if m.state == stateList {
		return appStyle.Render(m.list.View())
	}

	if m.state == stateDeleteConfirm {
		// Clean confirmation dialog
		confirmBox := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(dangerColor).
			Padding(2, 4).
			Width(60)

		titleStyle := lipgloss.NewStyle().
			Foreground(dangerColor).
			Bold(true).
			MarginBottom(1)

		promptStyle := lipgloss.NewStyle().
			Foreground(primaryColor).
			Padding(1, 0).
			MarginTop(1).
			MarginBottom(2).
			Italic(true)

		helpText := lipgloss.NewStyle().
			Foreground(subtleColor)

		content := titleStyle.Render("⚠ Delete Prompt?") + "\n\n" +
			promptStyle.Render(m.activePrompt.Title) + "\n" +
			helpText.Render("y/↵ confirm  •  n/esc cancel")

		return appStyle.Render("\n" + confirmBox.Render(content))
	}

	// Create Form View
	var b strings.Builder

	// Form header
	title := "Create New Prompt"
	if m.activePrompt != nil {
		title = "Edit Prompt"
	}
	b.WriteString(formTitleStyle.Render(title))
	b.WriteString("\n\n")

	// Form fields
	b.WriteString(m.inputView("Title", m.titleInput, m.focusIndex == 0))
	b.WriteString("\n")
	b.WriteString(m.inputView("Description", m.descriptionInput, m.focusIndex == 1))
	b.WriteString("\n")

	// Content textarea
	label := "Content"
	labelStyle := blurredPromptStyle
	if m.focusIndex == 2 {
		labelStyle = focusedPromptStyle
		label = "▸ " + label
	} else {
		label = "  " + label
	}
	b.WriteString(labelStyle.Render(label) + "\n")

	// Content box with subtle border
	contentBorder := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Padding(0, 1)

	if m.focusIndex == 2 {
		contentBorder = contentBorder.BorderForeground(primaryColor)
	}

	b.WriteString(contentBorder.Render(m.contentInput.View()))
	b.WriteString("\n\n")

	// Submit button
	btn := blurredButtonStyle.Render("  Submit  ")
	if m.focusIndex == 3 {
		btn = focusedButtonStyle.Render("▸ Submit ◂")
	}
	b.WriteString(btn)
	b.WriteString("\n")

	// Help text
	b.WriteString(helpStyle.Render("esc cancel  •  tab/shift+tab navigate  •  ↵ submit"))

	return appStyle.Render(b.String())
}

func (m Model) inputView(label string, input textinput.Model, focused bool) string {
	labelStyle := blurredPromptStyle
	if focused {
		labelStyle = focusedPromptStyle
		label = "▸ " + label
		input.PromptStyle = focusedPromptStyle
		input.TextStyle = inputStyle
	} else {
		label = "  " + label
		input.PromptStyle = blurredPromptStyle
		input.TextStyle = lipgloss.NewStyle().Foreground(mutedColor)
	}

	// Input with clean border
	inputBorder := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Padding(0, 1)

	if focused {
		inputBorder = inputBorder.BorderForeground(primaryColor)
	}

	return fmt.Sprintf("%s\n%s\n",
		labelStyle.Render(label),
		inputBorder.Render(input.View()))
}

func (m *Model) updateFocus() tea.Cmd {
	m.titleInput.Blur()
	m.descriptionInput.Blur()
	m.contentInput.Blur()

	switch m.focusIndex {
	case 0:
		return m.titleInput.Focus()
	case 1:
		return m.descriptionInput.Focus()
	case 2:
		return m.contentInput.Focus()
	}
	return nil
}

func (m *Model) resetForm() {
	m.titleInput.SetValue("")
	m.descriptionInput.SetValue("")
	m.contentInput.SetValue("")
	m.activePrompt = nil
	m.focusIndex = 0
	m.titleInput.Focus()
}

func (m *Model) setForm(p vault.Prompt) {
	m.titleInput.SetValue(p.Title)
	m.descriptionInput.SetValue(p.Description)
	m.contentInput.SetValue(p.PromptContent)
	m.focusIndex = 0
	m.titleInput.Focus()
}

// -- Commands --

type promptsMsg []vault.Prompt
type promptCreatedMsg struct{}
type promptDeletedMsg struct{}
type errMsg error

func (m Model) fetchPrompts() tea.Msg {
	prompts, err := m.service.GetAllPrompts()
	if err != nil {
		return errMsg(err)
	}
	return promptsMsg(prompts)
}

func (m Model) createPrompt() tea.Msg {
	id := 0
	var createdAt time.Time
	if m.activePrompt != nil {
		id = m.activePrompt.ID
		createdAt = m.activePrompt.CreatedAt
	}

	p := &vault.Prompt{
		ID:            id,
		CreatedAt:     createdAt,
		Title:         m.titleInput.Value(),
		Description:   m.descriptionInput.Value(),
		PromptContent: m.contentInput.Value(),
	}

	_, err := m.service.CreateOrUpdatePrompt(p)
	if err != nil {
		return errMsg(err)
	}
	return promptCreatedMsg{}
}

func (m Model) deletePrompt() tea.Msg {
	if m.activePrompt == nil {
		return nil
	}

	err := m.service.DeletePrompt(m.activePrompt.ID)
	if err != nil {
		return errMsg(err)
	}

	return promptDeletedMsg{}
}

// -- List Item Adapter --

type item struct {
	prompt vault.Prompt
}

func (i item) Title() string       { return i.prompt.Title }
func (i item) Description() string { return i.prompt.Description }
func (i item) FilterValue() string { return i.prompt.Title }

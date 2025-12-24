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
	// Initialize inputs
	ti := textinput.New()
	ti.Placeholder = "Prompt Title"
	ti.Focus()
	ti.CharLimit = 50
	ti.Width = 30

	desc := textinput.New()
	desc.Placeholder = "Description"
	desc.CharLimit = 100
	desc.Width = 50

	cont := textarea.New()
	cont.Placeholder = "Prompt Content..."
	cont.ShowLineNumbers = true
	cont.SetWidth(50)
	cont.SetHeight(10)

	// Initialize List
	items := []list.Item{}
	// Delegate will be set in Update or separate init function, but basic list config here
	// Delegate with custom styles
	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.
		Foreground(primaryColor).
		Background(lipgloss.Color("#1A1A1A")).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(primaryColor).
		BorderLeft(true).
		Padding(0, 1, 0, 2)
	delegate.Styles.SelectedDesc = delegate.Styles.SelectedDesc.
		Foreground(subtleColor).
		Background(lipgloss.Color("#1A1A1A")).
		Padding(0, 1, 0, 2)

	l := list.New(items, delegate, 0, 0)
	l.Title = "Prompts Vault"
	l.Styles.Title = listTitleStyle
	l.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			key.NewBinding(
				key.WithKeys("a"),
				key.WithHelp("a", "create prompt"),
			),
			key.NewBinding(
				key.WithKeys("e"),
				key.WithHelp("e", "edit prompt"),
			),
			key.NewBinding(
				key.WithKeys("d"),
				key.WithHelp("d", "delete prompt"),
			),
			key.NewBinding(
				key.WithKeys("enter"),
				key.WithHelp("enter", "copy to clipboard"),
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
					return m, m.list.NewStatusMessage(statusMessageStyle.Render("Copied to clipboard!"))
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
				if m.focusIndex == 2 && s == "enter" {
					// Textarea handles enter natively
					break
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
		cmds = append(cmds, m.list.NewStatusMessage(statusMessageStyle.Render("Prompt deleted!")))

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
		return errorMessageStyle.Render("Error: " + m.err.Error())
	}

	if m.state == stateList {
		return appStyle.Render(m.list.View())
	}

	if m.state == stateDeleteConfirm {
		var b strings.Builder
		b.WriteString("\n\n")
		b.WriteString(errorMessageStyle.Render("Are you sure you want to delete this prompt?"))
		b.WriteString("\n\n")
		b.WriteString(formTitleStyle.Render(m.activePrompt.Title))
		b.WriteString("\n\n")
		b.WriteString(statusMessageStyle.Render("(y/enter to confirm • n/esc to cancel)"))
		return appStyle.Render(b.String())
	}

	// Create Form View
	var b strings.Builder

	title := "✨ Create New Prompt"
	if m.activePrompt != nil {
		title = "✏️ Edit Prompt"
	}

	b.WriteString(formTitleStyle.Render(title))
	b.WriteString("\n")

	b.WriteString(m.inputView("Title", m.titleInput, m.focusIndex == 0))
	b.WriteString("\n")
	b.WriteString(m.inputView("Description", m.descriptionInput, m.focusIndex == 1))
	b.WriteString("\n")

	// Textarea special handling with improved styling
	label := "Content"
	if m.focusIndex == 2 {
		label = focusedPromptStyle.Render("▶ " + label)
	} else {
		label = blurredPromptStyle.Render("  " + label)
	}
	b.WriteString(label + "\n")
	b.WriteString(m.contentInput.View())
	b.WriteString("\n\n")

	btn := blurredButtonStyle.Render("Submit")
	if m.focusIndex == 3 {
		btn = focusedButtonStyle.Render("▶ Submit ◀")
	}
	b.WriteString(btn)

	b.WriteString(helpStyle.Render("\n\n┌─ (esc cancel • tab navigate • enter submit) ─┐"))

	return appStyle.Render(b.String())
}

func (m Model) inputView(label string, input textinput.Model, focused bool) string {
	if focused {
		label = focusedPromptStyle.Render("▶ " + label)
		input.TextStyle = focusedPromptStyle
	} else {
		label = blurredPromptStyle.Render("  " + label)
		input.TextStyle = blurredPromptStyle
	}
	return fmt.Sprintf("%s\n%s\n", label, input.View())
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

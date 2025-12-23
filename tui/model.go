package tui

import (
	"fmt"
	"strings"

	"github.com/Dima-salang/proompt-vault-tui/internal/vault"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type sessionState int

const (
	stateList sessionState = iota
	stateCreate
)

type Model struct {
	state   sessionState
	service vault.PromptService
	list    list.Model

	// Form inputs
	titleInput       textinput.Model
	descriptionInput textinput.Model
	contentInput     textarea.Model
	tagsInput        textinput.Model
	focusIndex       int

	err    error
	width  int
	height int
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

	tags := textinput.New()
	tags.Placeholder = "Tags (comma separated)"
	tags.CharLimit = 100
	tags.Width = 50

	// Initialize List
	items := []list.Item{}
	// Delegate will be set in Update or separate init function, but basic list config here
	delegate := list.NewDefaultDelegate()
	l := list.New(items, delegate, 0, 0)
	l.Title = "Prompts"
	l.SetShowHelp(true)
	l.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			key.NewBinding(
				key.WithKeys("a"),
				key.WithHelp("a", "create prompt"),
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
		tagsInput:        tags,
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
				m.state = stateCreate
				m.resetForm()
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

				// If in textarea (content input), enter should add new line unless ctrl+enter or moved away
				if m.focusIndex == 2 && s == "enter" {
					// Textarea handles enter natively
					break
				}

				// Navigation logic
				if s == "up" || s == "shift+tab" {
					m.focusIndex--
				} else if s == "down" || s == "tab" || (s == "enter" && m.focusIndex != 2) { // Skip enter for textarea
					m.focusIndex++
				}

				if m.focusIndex > 4 { // 0: Title, 1: Desc, 2: Content, 3: Tags, 4: Submit
					m.focusIndex = 0
				} else if m.focusIndex < 0 {
					m.focusIndex = 4
				}

				// Handle Submit on Enter if focused on Submit button (represented by index 4)
				if m.focusIndex == 4 && s == "enter" {
					return m, m.createPrompt
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
		m.tagsInput, cmd = m.tagsInput.Update(msg)
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

	// Create Form View
	var b strings.Builder
	b.WriteString(headerStyle.Render("Create New Prompt"))
	b.WriteString("\n\n")

	b.WriteString(m.inputView("Title", m.titleInput, m.focusIndex == 0))
	b.WriteString("\n")
	b.WriteString(m.inputView("Description", m.descriptionInput, m.focusIndex == 1))
	b.WriteString("\n")

	// Textarea special handling for label
	label := "Content"
	if m.focusIndex == 2 {
		label = statusMessageStyle.Render(label)
	}
	b.WriteString(label + "\n")
	b.WriteString(m.contentInput.View())
	b.WriteString("\n\n")

	b.WriteString(m.inputView("Tags", m.tagsInput, m.focusIndex == 3))
	b.WriteString("\n\n")

	// Submit Button
	btn := "[ Submit ]"
	if m.focusIndex == 4 {
		btn = statusMessageStyle.Render(btn)
	}
	b.WriteString(btn)

	b.WriteString("\n\n(esc to cancel, tab to navigate)")

	return appStyle.Render(b.String())
}

func (m Model) inputView(label string, input textinput.Model, focused bool) string {
	if focused {
		label = statusMessageStyle.Render(label)
	}
	return fmt.Sprintf("%s\n%s\n", label, input.View())
}

func (m *Model) updateFocus() tea.Cmd {
	m.titleInput.Blur()
	m.descriptionInput.Blur()
	m.contentInput.Blur()
	m.tagsInput.Blur()

	switch m.focusIndex {
	case 0:
		return m.titleInput.Focus()
	case 1:
		return m.descriptionInput.Focus()
	case 2:
		return m.contentInput.Focus()
	case 3:
		return m.tagsInput.Focus()
	}
	return nil
}

func (m *Model) resetForm() {
	m.titleInput.SetValue("")
	m.descriptionInput.SetValue("")
	m.contentInput.SetValue("")
	m.tagsInput.SetValue("")
	m.focusIndex = 0
	m.titleInput.Focus()
}

// -- Commands --

type promptsMsg []vault.Prompt
type promptCreatedMsg struct{}
type errMsg error

func (m Model) fetchPrompts() tea.Msg {
	prompts, err := m.service.GetAllPrompts()
	if err != nil {
		return errMsg(err)
	}
	return promptsMsg(prompts)
}

func (m Model) createPrompt() tea.Msg {
	tags := strings.Split(m.tagsInput.Value(), ",")
	for i, t := range tags {
		tags[i] = strings.TrimSpace(t)
	}

	p := &vault.Prompt{
		Title:         m.titleInput.Value(),
		Description:   m.descriptionInput.Value(),
		PromptContent: m.contentInput.Value(),
		Tags:          tags,
	}

	_, err := m.service.CreateOrUpdatePrompt(p)
	if err != nil {
		return errMsg(err)
	}
	return promptCreatedMsg{}
}

// -- List Item Adapter --

type item struct {
	prompt vault.Prompt
}

func (i item) Title() string       { return i.prompt.Title }
func (i item) Description() string { return i.prompt.Description }
func (i item) FilterValue() string { return i.prompt.Title }

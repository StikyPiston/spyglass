package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/stikypiston/spyglass/lens"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type viewState int

const (
	stateEntries viewState = iota
	stateContext
)

type model struct {
	lenses     []lens.Lens
	activeLens int

	search textinput.Model

	entries    []lens.Entry
	selected   int
	scroll     int
	state      viewState
	actions    []lens.Action
	contextFor lens.Entry

	width  int
	height int
}

func newModel() model {
	ti := textinput.New()
	ti.Placeholder = "Search..."
	ti.Focus()
	ti.CharLimit = 256

	m := model{
		lenses: Lenses,
		search: ti,
	}
	m.refresh()
	return m
}

func (m *model) refresh() {
	entries, _ := m.lenses[m.activeLens].Search(m.search.Value())
	m.entries = entries
	if m.selected >= len(entries) {
		m.selected = len(entries) - 1
	}
	if m.selected < 0 {
		m.selected = 0
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		switch msg.Type {

		case tea.KeyCtrlC:
			return m, tea.Quit

		case tea.KeyTab:
			m.activeLens = (m.activeLens + 1) % len(m.lenses)
			m.selected = 0
			m.scroll = 0
			m.refresh()

		case tea.KeyShiftTab:
			m.activeLens--
			if m.activeLens < 0 {
				m.activeLens = len(m.lenses) - 1
			}
			m.selected = 0
			m.scroll = 0
			m.refresh()

		case tea.KeyUp:
			if m.selected > 0 {
				m.selected--
			}

		case tea.KeyDown:
			if m.selected < len(m.entries)-1 {
				m.selected++
			}

		case tea.KeyEnter:
			if msg.String() == "shift+enter" {
				if len(m.entries) > 0 {
					entry := m.entries[m.selected]
					m.contextFor = entry
					m.actions = m.lenses[m.activeLens].ContextActions(entry)
					m.state = stateContext
					m.selected = 0
				}
				return m, nil
			}

			if m.state == stateEntries && len(m.entries) > 0 {
				entry := m.entries[m.selected]
				m.lenses[m.activeLens].Enter(entry)
			} else if m.state == stateContext {
				if len(m.actions) > 0 {
					m.actions[m.selected].Run(m.contextFor)
				}
				m.state = stateEntries
				m.refresh()
			}

		case tea.KeyEsc:
			m.state = stateEntries
		}
	}

	var cmd tea.Cmd
	m.search, cmd = m.search.Update(msg)
	m.refresh()
	return m, cmd
}

func (m model) View() string {
	if m.width <= 0 || m.height <= 0 {
		return ""
	}

	border := lipgloss.RoundedBorder()

	// Compute available space properly
	totalHeight := m.height

	tabHeight := 3
	descHeight := 5
	searchHeight := 3

	listHeight := totalHeight - tabHeight - descHeight - searchHeight
	if listHeight < 3 {
		listHeight = 3
	}

	contentWidth := m.width - 2

	tabStyle := lipgloss.NewStyle().
		Border(border).
		Width(contentWidth).
		Height(tabHeight-2).
		Padding(0, 1)

	listStyle := lipgloss.NewStyle().
		Border(border).
		Width(contentWidth).
		Height(listHeight-2).
		Padding(0, 1)

	descStyle := lipgloss.NewStyle().
		Border(border).
		Width(contentWidth).
		Height(descHeight-2).
		Padding(0, 1)

	searchStyle := lipgloss.NewStyle().
		Border(border).
		Width(contentWidth).
		Height(searchHeight-2).
		Padding(0, 1)

	// Tabs
	var tabs []string
	for i, l := range m.lenses {
		if i == m.activeLens {
			tabs = append(tabs, "["+l.Name()+"]")
		} else {
			tabs = append(tabs, l.Name())
		}
	}
	tabsBox := tabStyle.Render(strings.Join(tabs, " | "))

	// LIST CONTENT
	var listBuilder strings.Builder

	if m.state == stateContext {
		for i, a := range m.actions {
			cursor := "  "
			if i == m.selected {
				cursor = "> "
			}
			listBuilder.WriteString(cursor + a.Name + "\n")
		}
	} else {
		maxVisible := listHeight - 2

		if m.selected < m.scroll {
			m.scroll = m.selected
		}
		if m.selected >= m.scroll+maxVisible {
			m.scroll = m.selected - maxVisible + 1
		}

		start := m.scroll
		end := start + maxVisible
		if end > len(m.entries) {
			end = len(m.entries)
		}

		for i := start; i < end; i++ {
			cursor := "  "
			if i == m.selected {
				cursor = "> "
			}
			listBuilder.WriteString(
				fmt.Sprintf("%s%s %s\n", cursor, m.entries[i].Icon, m.entries[i].Title),
			)
		}
	}

	listBox := listStyle.Render(listBuilder.String())

	// Description
	desc := ""
	if len(m.entries) > 0 && m.selected < len(m.entries) {
		desc = m.entries[m.selected].Description
	}
	descBox := descStyle.Render(desc)

	searchBox := searchStyle.Render(m.search.View())

	return lipgloss.JoinVertical(
		lipgloss.Left,
		tabsBox,
		listBox,
		descBox,
		searchBox,
	)
}

func main() {
	p := tea.NewProgram(newModel())
	if err := p.Start(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

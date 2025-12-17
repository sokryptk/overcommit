package components

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

type scopeItem string

func (s scopeItem) FilterValue() string { return string(s) }

type scopeDelegate struct{}

func (d scopeDelegate) Height() int                             { return 1 }
func (d scopeDelegate) Spacing() int                            { return 0 }
func (d scopeDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

func (d scopeDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	s := string(item.(scopeItem))
	txt := fmt.Sprintf("  %s", s)
	if index == m.Index() {
		txt = termenv.String(fmt.Sprintf("> %s", s)).Foreground(term.Color("#8AA8F9")).Underline().String()
	} else {
		txt = termenv.String(txt).Faint().String()
	}
	fmt.Fprint(w, txt)
}

type ScopeSelectorView struct {
	view list.Model
}

func NewScopeSelector(scopes []string) ScopeSelectorView {
	items := make([]list.Item, len(scopes))
	for i, s := range scopes {
		items[i] = scopeItem(s)
	}

	height := len(items) + 4
	if height > 12 {
		height = 12
	}
	li := list.New(items, scopeDelegate{}, 40, height)
	li.Title = "Select scope (Esc to skip):"
	li.SetShowTitle(true)
	li.SetShowStatusBar(false)
	li.SetShowPagination(false)
	li.SetShowHelp(false)
	li.SetFilteringEnabled(true)
	li.Styles.Title = lipgloss.NewStyle().Foreground(lipgloss.Color("#8AA8F9")).Bold(true)
	li.Styles.TitleBar = lipgloss.NewStyle()

	return ScopeSelectorView{view: li}
}

func (s *ScopeSelectorView) Update(msg tea.Msg, v PageView) (PageView, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.view.SetSize(msg.Width, msg.Height)
		return v, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if s.view.FilterState() == list.Filtering {
				break
			}
			if len(s.view.Items()) > 0 && s.view.Index() < len(s.view.Items()) {
				v.scope = string(s.view.SelectedItem().(scopeItem))
			}
			v.Page = MSG
			return v, nil
		case "esc":
			if s.view.FilterState() != list.Filtering {
				v.scope = ""
				v.Page = MSG
				return v, nil
			}
		}
	}

	s.view, cmd = s.view.Update(msg)
	return v, cmd
}

func (s ScopeSelectorView) View() string {
	return s.view.View()
}

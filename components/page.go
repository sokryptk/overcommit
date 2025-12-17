package components

import (
	tea "github.com/charmbracelet/bubbletea"
	"me.kryptk.overcommit/utils"
)

type Page int

const (
	SELECTION = iota
	SCOPE
	MSG
)

type PageView struct {
	Page          Page
	selected      utils.Key
	scope         string
	Template      utils.Template
	Selector      *TypeSelectorView
	ScopeSelector *ScopeSelectorView
	Committer     *CommitView
	FinalMessage  string
}

func (p PageView) Init() tea.Cmd {
	return nil
}

func (p PageView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return p, tea.Quit
		}
	}

	switch p.Page {
	case SELECTION:
		return p.Selector.Update(msg, p)
	case SCOPE:
		return p.ScopeSelector.Update(msg, p)
	default:
		return p.Committer.Update(msg, p)
	}
}

func (p PageView) View() string {
	switch p.Page {
	case SELECTION:
		if p.Selector == nil {
			return ""
		}
		return p.Selector.View()
	case SCOPE:
		if p.ScopeSelector == nil {
			return ""
		}
		return p.ScopeSelector.View()
	default:
		return p.Committer.View(p)
	}
}

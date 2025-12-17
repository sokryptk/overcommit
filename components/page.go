package components

import (
	tea "github.com/charmbracelet/bubbletea"
	"me.kryptk.overcommit/utils"
)

type Page int

const (
	SELECTION = iota
	MSG
)

type PageView struct {
	Page      Page
	selected  utils.Key
	message   string
	Template  utils.Template
	Selector  *TypeSelectorView
	Committer *CommitView
}

func (p PageView) Init() tea.Cmd {
	return nil
}

func (p PageView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c":
			return p, tea.Quit
		}
	}

	if p.Page == SELECTION {
		return p.Selector.Update(msg, p)
	}

	return p.Committer.Update(msg, p)
}

func (p PageView) View() string {
	switch p.Page {
	case SELECTION:
		if p.Selector == nil {
			return ""
		}

		return p.Selector.View()
	}

	return p.Committer.View(p)
}

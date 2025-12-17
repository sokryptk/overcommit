package components

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/termenv"
	"me.kryptk.overcommit/utils"
)

type CommitView struct {
	msgInput  textinput.Model
	maxLength int
	err       string
}

func NewCommitView(maxLength int) CommitView {
	ti := textinput.New()
	ti.Prompt = ""
	ti.Placeholder = "describe your change"
	ti.Focus()

	return CommitView{
		msgInput:  ti,
		maxLength: maxLength,
	}
}

func (i *CommitView) Update(msg tea.Msg, v PageView) (PageView, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			val := i.msgInput.Value()
			if len(val) == 0 {
				i.err = "message required"
				return v, nil
			}
			if len(val) > i.maxLength {
				i.err = fmt.Sprintf("exceeds %d chars", i.maxLength)
				return v, nil
			}
			i.err = ""

			fileName := os.Args[1]
			commitMsg := utils.BuildCommitMessage(v.Template, v.selected.Prefix, v.scope, val)
			_ = utils.ReplaceHeaderFromCommit(commitMsg, fileName)
			return PageView{}, tea.Quit
		}
	}

	i.msgInput, cmd = i.msgInput.Update(msg)
	return v, cmd
}

func (i CommitView) View(v PageView) string {
	style := termenv.String().Bold().Foreground(ACCENT).Styled
	errStyle := termenv.String().Bold().Foreground(term.Color("#FF5555")).Styled

	currentLen := len(i.msgInput.Value())
	counter := fmt.Sprintf("[%d/%d]", currentLen, i.maxLength)
	if currentLen > i.maxLength {
		counter = errStyle(counter)
	}

	view := fmt.Sprintf("%s : %s - %s\n", style("[Commit Type]"), v.selected.Prefix, v.selected.Description)
	if v.scope != "" {
		view += fmt.Sprintf("%s : %s\n", style("[Scope]"), v.scope)
	}
	view += fmt.Sprintf("%s %s : %s", style("[Message]"), counter, i.msgInput.View())

	if i.err != "" {
		view += "\n" + errStyle(i.err)
	}

	return view
}

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
	msgInput textinput.Model
}

func NewCommitView() CommitView {
	ti := textinput.New()
	ti.Prompt = fmt.Sprintf("%s : ", SetTextStyle("[Enter commit message]"))
	ti.Placeholder = "pikachu"

	ti.Focus()

	return CommitView{
		msgInput: ti,
	}
}

func (i *CommitView) Update(msg tea.Msg, v PageView) (PageView, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "enter":
			fileName := os.Args[1]
			_ = utils.ReplaceHeaderFromCommit(utils.BuildPrefixWithMsg(v.Template, v.selected.Prefix, i.msgInput.Value()), fileName)

			return PageView{}, tea.Quit
		}
	}

	i.msgInput, cmd = i.msgInput.Update(msg)

	return v, cmd
}

func (i CommitView) View(v PageView) string {
	var view string

	style := termenv.String().Bold().Foreground(ACCENT).Styled

	view += fmt.Sprintf("%s : %s - %s\n", style("[Commit Type]"), v.selected.Prefix, v.selected.Description)
	view += i.msgInput.View()

	return view
}

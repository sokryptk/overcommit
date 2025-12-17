package components

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
	"me.kryptk.overcommit/utils"
)

type generatedMsg struct {
	text string
	err  error
}

type CommitView struct {
	msgInput   textinput.Model
	spinner    spinner.Model
	maxLength  int
	err        string
	llmClient  utils.LLMClient
	generating bool
}

func NewCommitView(maxLength int, llmCfg utils.LLMConfig) CommitView {
	ti := textinput.New()
	ti.Prompt = ""
	ti.Placeholder = "describe your change (g to generate)"
	ti.Focus()

	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#8AA8F9"))

	return CommitView{
		msgInput:  ti,
		spinner:   sp,
		maxLength: maxLength,
		llmClient: utils.NewLLMClient(llmCfg),
	}
}

func (c *CommitView) generate(v PageView) tea.Cmd {
	return func() tea.Msg {
		diff, err := utils.GetStagedDiff()
		if err != nil {
			return generatedMsg{err: err}
		}
		if diff == "" {
			return generatedMsg{err: fmt.Errorf("no staged changes")}
		}

		prompt := utils.BuildPrompt(v.selected.Prefix, v.scope, diff)
		text, err := c.llmClient.Generate(prompt)
		return generatedMsg{text: text, err: err}
	}
}

func (c *CommitView) Update(msg tea.Msg, v PageView) (PageView, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case generatedMsg:
		c.generating = false
		if msg.err != nil {
			c.err = msg.err.Error()
		} else {
			c.msgInput.SetValue(msg.text)
			c.err = ""
		}
		return v, nil

	case spinner.TickMsg:
		if c.generating {
			c.spinner, cmd = c.spinner.Update(msg)
			return v, cmd
		}

	case tea.KeyMsg:
		if c.generating {
			return v, nil
		}

		switch msg.String() {
		case "g":
			if c.msgInput.Value() == "" {
				c.generating = true
				c.err = ""
				return v, tea.Batch(c.spinner.Tick, c.generate(v))
			}
		case "enter":
			val := c.msgInput.Value()
			if len(val) == 0 {
				c.err = "message required"
				return v, nil
			}
			if len(val) > c.maxLength {
				c.err = fmt.Sprintf("exceeds %d chars", c.maxLength)
				return v, nil
			}
			c.err = ""

			fileName := os.Args[1]
			commitMsg := utils.BuildCommitMessage(v.Template, v.selected.Prefix, v.scope, val)
			_ = utils.ReplaceHeaderFromCommit(commitMsg, fileName)
			return PageView{}, tea.Quit
		}
	}

	c.msgInput, cmd = c.msgInput.Update(msg)
	return v, cmd
}

func (c CommitView) View(v PageView) string {
	style := termenv.String().Bold().Foreground(ACCENT).Styled
	errStyle := termenv.String().Bold().Foreground(term.Color("#FF5555")).Styled

	currentLen := len(c.msgInput.Value())
	counter := fmt.Sprintf("[%d/%d]", currentLen, c.maxLength)
	if currentLen > c.maxLength {
		counter = errStyle(counter)
	}

	view := fmt.Sprintf("%s : %s - %s\n", style("[Commit Type]"), v.selected.Prefix, v.selected.Description)
	if v.scope != "" {
		view += fmt.Sprintf("%s : %s\n", style("[Scope]"), v.scope)
	}

	if c.generating {
		view += fmt.Sprintf("%s : %s generating...", style("[Message]"), c.spinner.View())
	} else {
		view += fmt.Sprintf("%s %s : %s", style("[Message]"), counter, c.msgInput.View())
	}

	if c.err != "" {
		view += "\n" + errStyle(c.err)
	}

	return view
}

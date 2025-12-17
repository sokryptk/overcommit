package components

import (
	"fmt"
	"io"
	"strconv"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
	"me.kryptk.overcommit/utils"
)

var (
	term = termenv.TrueColor
)

func NewTypeSelector(keys []utils.Key) TypeSelectorView {
	// just a type coercion
	items := keysToItems(keys)

	li := list.New(items, listDelegate{}, 40, len(items)+4)
	li.Title = "Select commit type:"
	li.SetShowTitle(true)
	li.SetShowStatusBar(false)
	li.SetShowPagination(false)
	li.SetShowHelp(false)
	li.SetFilteringEnabled(true)
	li.Styles.Title = lipgloss.NewStyle().Foreground(lipgloss.Color("#8AA8F9")).Bold(true)
	li.Styles.TitleBar = lipgloss.NewStyle()
	return TypeSelectorView{
		view: li,
	}
}

type TypeSelectorView struct {
	view list.Model
}

type listDelegate struct{}

func (l listDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	selected := index == m.Index()

	i, ok := item.(utils.Key)
	if !ok {
		return
	}

	txt := fmt.Sprintf("(%s) - %s [%d]", i.Prefix, i.Description, index+1)

	if selected {
		txt = termenv.String(txt).Foreground(term.Color("#8AA8F9")).Underline().String()
	} else {
		txt = termenv.String(txt).Faint().String()
	}

	_, _ = fmt.Fprint(w, txt)
}

func (l listDelegate) Height() int { return 1 }

func (l listDelegate) Spacing() int {
	return 0
}

func (l listDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd {
	return nil
}

func (tsv TypeSelectorView) View() string {
	return tsv.view.View()
}

func (tsv *TypeSelectorView) Update(msg tea.Msg, v PageView) (PageView, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		tsv.view.SetSize(msg.Width, msg.Height)
		return v, nil
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "enter", tea.KeyRight.String():
			v.selected = tsv.view.SelectedItem().(utils.Key)
			v.Page = SCOPE
			return v, nil
		default:
			index, err := strconv.Atoi(keypress)
			if err != nil {
				break
			}
			if index >= 1 && index <= len(tsv.view.Items()) {
				v.selected = tsv.view.Items()[index-1].(utils.Key)
				v.Page = SCOPE
				return v, nil
			}
		}
	}

	tsv.view, cmd = tsv.view.Update(msg)

	return v, cmd
}

func keysToItems(keys []utils.Key) []list.Item {
	items := make([]list.Item, len(keys))

	for i, k := range keys {
		items[i] = k
	}

	return items
}

package main

import (
	_ "embed"
	"fmt"
	"log"
	"os"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
	"me.kryptk.overcommit/components"
	"me.kryptk.overcommit/utils"
)

//go:embed config.toml
var config string

func main() {
	hook := "$PWD/.git/hooks/prepare-commit-msg"

	_, err := os.ReadDir(os.ExpandEnv("$PWD/.git"))
	if err != nil {
		fmt.Println("not a git repository")
		return
	}

	if len(os.Args) <= 1 {
		fmt.Println("Hi, set up using -i or --init")
		return
	}

	if len(os.Args) > 1 {
		// check if initing
		isInit := os.Args[1] == "-i" || os.Args[1] == "--init"

		if isInit {
			_ = os.Rename(os.ExpandEnv(fmt.Sprintf("%s.sample", hook)), os.ExpandEnv(fmt.Sprintf("%s.bak", hook)))
			hookScript := "#!/bin/sh\novercommit \"$1\" \"$2\"\n"
			_ = os.WriteFile(os.ExpandEnv(hook), []byte(hookScript), 0755)

			exec.Command("git", "config", "alias.oc", "commit --no-edit").Run()

			fmt.Println("overcommit installed!")
			fmt.Println("use: git oc")

			return
		}
	}

	c, err := utils.LoadConfig(config)
	if err != nil {
		log.Fatal(err)
	}

	selector := components.NewTypeSelector(c.Keys)
	scopeSelector := components.NewScopeSelector(utils.GetScopes())
	committer := components.NewCommitView(c.Lint.MaxSubjectLength, c.LLM)

	m := components.PageView{
		Page:          components.SELECTION,
		Selector:      &selector,
		ScopeSelector: &scopeSelector,
		Committer:     &committer,
		Template:      c.Template,
	}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

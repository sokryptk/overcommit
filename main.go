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
	_, err := os.ReadDir(os.ExpandEnv("$PWD/.git"))
	if err != nil {
		fmt.Println("not a git repository")
		return
	}

	if len(os.Args) > 1 && (os.Args[1] == "-i" || os.Args[1] == "--init") {
		hook := "$PWD/.git/hooks/prepare-commit-msg"
		_ = os.Rename(os.ExpandEnv(fmt.Sprintf("%s.sample", hook)), os.ExpandEnv(fmt.Sprintf("%s.bak", hook)))
		_ = os.WriteFile(os.ExpandEnv(hook), []byte("#!/bin/sh\novercommit \"$1\" \"$2\"\n"), 0755)
		fmt.Println("overcommit installed! just run: overcommit")
		return
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

	finalModel, err := tea.NewProgram(m).Run()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	result := finalModel.(components.PageView)
	if result.FinalMessage == "" {
		return
	}

	if len(os.Args) <= 1 {
		cmd := exec.Command("git", "commit", "-m", result.FinalMessage)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
	}
}

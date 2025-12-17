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
		initRepo()
		return
	}

	runTUI()
}

func initRepo() {
	hookDir := ".githooks"
	os.MkdirAll(hookDir, 0755)

	hook := `#!/bin/sh
msg=$(head -1 "$1")
if ! echo "$msg" | grep -qE '^(feat|fix|docs|style|refactor|test|chore)(\(.+\))?: .+'; then
  echo "bad commit message: $msg"
  echo ""
  echo "expected: type(scope): message"
  echo "types: feat|fix|docs|style|refactor|test|chore"
  echo ""
  echo "use 'overcommit' for easy conventional commits"
  exit 1
fi
`
	os.WriteFile(hookDir+"/commit-msg", []byte(hook), 0755)
	exec.Command("git", "config", "core.hooksPath", hookDir).Run()

	fmt.Println("done. commit .githooks/ to enforce for team")
}

func runTUI() {
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

	cmd := exec.Command("git", "commit", "-m", result.FinalMessage)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

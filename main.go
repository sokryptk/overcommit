package main

import (
	_ "embed"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"log"
	"me.kryptk.overcommit/components"
	"me.kryptk.overcommit/utils"
	"os"
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
			//back up previous stuff
			_ = os.Rename(os.ExpandEnv(fmt.Sprintf("%s.sample", hook)), os.ExpandEnv(fmt.Sprintf("%s.bak", hook)))
			_ = os.WriteFile(os.ExpandEnv(hook), []byte("overcommit $1 $2"), 0755)

			fmt.Println(os.ExpandEnv("Successfully set up overcommit in $PWD!"))

			return
		}
	}

	c, err := utils.GenerateConfig(config)
	if err != nil {
		log.Fatal(err)
	}

	selector := components.NewTypeSelector(c.Keys)
	committer := components.NewCommitView()

	m := components.PageView{Page: components.SELECTION, Selector: &selector, Committer: &committer, Template: c.Template}

	if err := tea.NewProgram(m, tea.WithAltScreen(), tea.WithANSICompressor()).Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

package main

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/ParsaImi/sshelf/internal/config"
	"github.com/ParsaImi/sshelf/internal/git"
	"github.com/ParsaImi/sshelf/internal/history"
	"github.com/ParsaImi/sshelf/internal/model"
	"github.com/ParsaImi/sshelf/internal/ssh"
	"github.com/ParsaImi/sshelf/internal/ui"
	"github.com/ParsaImi/sshelf/internal/ui/screens"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "-h" {
		runHistory()
		return
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	dir := config.Dir()
	if cfg.Meta.GitAutoPull && git.IsRepo(dir) {
		if err := git.Pull(dir); err != nil {
			fmt.Fprintf(os.Stderr, "sshelf: git pull failed: %v\n", err)
		} else {
			if newCfg, loadErr := config.Load(); loadErr == nil {
				cfg = newCfg
			}
		}
	}

	p := tea.NewProgram(ui.NewApp(cfg), tea.WithAltScreen())
	finalModel, err := p.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "sshelf: %v\n", err)
		os.Exit(1)
	}

	app := finalModel.(ui.App)
	if app.SelectedHost != nil {
		connect(*app.SelectedHost)
	}
}

func runHistory() {
	hosts := history.Recent(5)
	if len(hosts) == 0 {
		fmt.Println("No connection history yet.")
		return
	}

	p := tea.NewProgram(
		screens.NewHistoryList(hosts),
		tea.WithAltScreen(),
	)
	finalModel, err := p.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "sshelf: %v\n", err)
		os.Exit(1)
	}

	hl := finalModel.(screens.HistoryList)
	if hl.Selected != nil {
		connect(*hl.Selected)
	}
}

func connect(h model.Host) {
	history.Record(h)
	if err := ssh.Exec(h); err != nil {
		fmt.Fprintf(os.Stderr, "sshelf: %v\n", err)
		os.Exit(1)
	}
}

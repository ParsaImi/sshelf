package screens

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/ParsaImi/imiterm/internal/model"
	"github.com/ParsaImi/imiterm/internal/ui/keys"
	"github.com/ParsaImi/imiterm/internal/ui/styles"
)

type HistoryList struct {
	hosts    []model.Host
	cursor   int
	Selected *model.Host
	width    int
	height   int
}

func NewHistoryList(hosts []model.Host) HistoryList {
	return HistoryList{hosts: hosts}
}

func (h HistoryList) Init() tea.Cmd { return nil }

func (h HistoryList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h.width = msg.Width
		h.height = msg.Height
		return h, nil
	}

	if msg, ok := msg.(tea.KeyMsg); ok {
		switch {
		case keys.Match(msg, keys.GroupList.Up):
			if h.cursor > 0 {
				h.cursor--
			}
		case keys.Match(msg, keys.GroupList.Down):
			if h.cursor < len(h.hosts)-1 {
				h.cursor++
			}
		case keys.Match(msg, keys.GroupList.Enter):
			host := h.hosts[h.cursor]
			h.Selected = &host
			return h, tea.Quit
		case keys.Match(msg, keys.GroupList.Quit):
			return h, tea.Quit
		}
	}
	return h, nil
}

func (h HistoryList) View() string {
	width := h.width
	if width <= 0 {
		width = 120
	}
	header := styles.RenderHeader(
		"IMITERM",
		"Recent connections",
		styles.RenderPill("saved", fmt.Sprintf("%d", len(h.hosts)), true),
	)

	var rows string
	if len(h.hosts) == 0 {
		rows = styles.RenderEmpty("No history yet", "Connect to a host and it will appear here.")
	} else {
		for i, host := range h.hosts {
			rows += renderHostRow(host, i == h.cursor, "") + "\n"
		}
	}

	content := styles.RenderCard("History", rows, min(width-8, 100), false)
	footer := styles.Pill.Render("enter reconnect") + "  " + styles.Pill.Render("q quit")

	help := styles.RenderHelp("enter", "connect", "q", "quit")

	return header + "\n" + content + "\n" + footer + "\n" + help + "\n"
}

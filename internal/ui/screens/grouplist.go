package screens

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/ParsaImi/imiterm/internal/model"
	"github.com/ParsaImi/imiterm/internal/ui/keys"
	"github.com/ParsaImi/imiterm/internal/ui/styles"
)

type GroupSelectedMsg struct {
	Index int
}

type GroupList struct {
	groups []model.Group
	cursor int
	width  int
	height int
}

func NewGroupList(groups []model.Group) GroupList {
	return GroupList{groups: groups, cursor: 0}
}

func (g GroupList) Cursor() int { return g.cursor }

func (g GroupList) Init() tea.Cmd { return nil }

func (g GroupList) Update(msg tea.Msg) (GroupList, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		g.width = msg.Width
		g.height = msg.Height
	case tea.KeyMsg:
		switch {
		case keys.Match(msg, keys.GroupList.Up):
			if g.cursor > 0 {
				g.cursor--
			}
		case keys.Match(msg, keys.GroupList.Down):
			if g.cursor < len(g.groups)-1 {
				g.cursor++
			}
		case keys.Match(msg, keys.GroupList.Enter):
			if len(g.groups) > 0 {
				return g, func() tea.Msg { return GroupSelectedMsg{Index: g.cursor} }
			}
		case keys.Match(msg, keys.GroupList.Quit):
			return g, tea.Quit
		}
	}
	return g, nil
}

func (g GroupList) View() string {
	totalHosts := 0
	for _, group := range g.groups {
		totalHosts += len(group.Hosts)
	}

	width := g.width
	if width <= 0 {
		width = 120
	}

	header := styles.RenderHeader(
		"IMITERM",
		"SSH manager",
		styles.RenderPill("groups", fmt.Sprintf("%d", len(g.groups)), true),
		styles.RenderPill("hosts", fmt.Sprintf("%d", totalHosts), false),
	)

	leftWidth := max(40, min(54, width/3))
	rightWidth := max(34, width-leftWidth-6)

	leftBody := g.renderGroupListBody()
	rightBody := g.renderGroupDetailBody()

	leftPanel := styles.RenderCard("Groups", leftBody, leftWidth, true)
	rightPanel := styles.RenderCard("Selected group", rightBody, rightWidth, false)
	content := styles.RenderDualCards(leftPanel, rightPanel, "  ")
	if width < 100 {
		content = leftPanel + "\n\n" + rightPanel
	}

	footer := styles.Pill.Render("↑/k down") + "  " + styles.Pill.Render("enter open") + "  " + styles.Pill.Render("a add") + "  " + styles.Pill.Render("e edit") + "  " + styles.Pill.Render("d delete") + "  " + styles.Pill.Render("/ search") + "  " + styles.Pill.Render("s sync") + "  " + styles.Pill.Render("q quit")

	help := styles.RenderHelp(
		"↑/k", "up", "↓/j", "down", "enter", "open",
		"a", "add", "e", "edit", "d", "delete",
		"/", "search", "s", "sync", "q", "quit",
	)

	return header + "\n" + content + "\n" + footer + "\n" + help + "\n"
}

func (g GroupList) renderGroupListBody() string {
	if len(g.groups) == 0 {
		return styles.RenderEmpty(
			"No groups yet",
			"Press a to create your first group, then add hosts and connect with Enter.",
		) + "\n\n" + styles.PanelSubtle.Render("Config: ~/.config/imiterm/config.toml")
	}

	var rows []string
	for i, group := range g.groups {
		rows = append(rows, renderGroupRow(group, i == g.cursor))
	}
	return strings.Join(rows, "\n")
}

func (g GroupList) renderGroupDetailBody() string {
	if len(g.groups) == 0 {
		return styles.RenderEmpty("Get started", "Add a group, then put hosts inside it.")
	}

	group := g.groups[g.cursor]
	var b strings.Builder
	b.WriteString(styles.DetailLabel.Render("Name") + styles.DetailValue.Render(group.Name) + "\n")
	if group.Description != "" {
		b.WriteString(styles.DetailLabel.Render("About") + styles.DetailValue.Render(group.Description) + "\n")
	} else {
		b.WriteString(styles.DetailLabel.Render("About") + styles.PanelSubtle.Render("No description") + "\n")
	}
	b.WriteString(styles.DetailLabel.Render("Hosts") + styles.DetailValue.Render(fmt.Sprintf("%d", len(group.Hosts))) + "\n")
	b.WriteString(styles.DetailLabel.Render("Color") + styles.GroupTag.Render(group.Color) + "\n\n")
	b.WriteString(styles.SectionTitle.Render("Preview") + "\n")
	if len(group.Hosts) == 0 {
		b.WriteString(styles.PanelSubtle.Render("No hosts in this group yet."))
	} else {
		limit := min(4, len(group.Hosts))
		for i := 0; i < limit; i++ {
			host := group.Hosts[i]
			b.WriteString(styles.SubtleRow.Render(fmt.Sprintf("%s  %s@%s:%d", host.Name, host.User, host.Hostname, host.Port)) + "\n")
		}
		if len(group.Hosts) > limit {
			b.WriteString(styles.PanelSubtle.Render(fmt.Sprintf("… and %d more", len(group.Hosts)-limit)))
		}
	}
	return strings.TrimRight(b.String(), "\n")
}

func renderGroupRow(group model.Group, selected bool) string {
	dot := lipgloss.NewStyle().Foreground(lipgloss.Color(group.Color)).Render("●")
	line := fmt.Sprintf("%s  %-16s  %s", dot, group.Name, styles.GroupCount.Render(fmt.Sprintf("%d hosts", len(group.Hosts))))
	if selected {
		return styles.ActiveRow.Render(styles.SelectedMarker.Render("›") + " " + line)
	}
	return styles.PassiveRow.Render("  " + line)
}

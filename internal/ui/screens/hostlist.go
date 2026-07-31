package screens

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/ParsaImi/imiterm/internal/model"
	"github.com/ParsaImi/imiterm/internal/ui/keys"
	"github.com/ParsaImi/imiterm/internal/ui/styles"
)

type HostSelectedMsg struct {
	Host model.Host
}

type BackMsg struct{}

type HostList struct {
	group  model.Group
	hosts  []model.Host
	cursor int
	width  int
	height int
}

func NewHostList(group model.Group) HostList {
	return HostList{group: group, hosts: group.Hosts, cursor: 0}
}

func (h HostList) Cursor() int { return h.cursor }

func (h HostList) Init() tea.Cmd { return nil }

func (h HostList) Update(msg tea.Msg) (HostList, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h.width = msg.Width
		h.height = msg.Height
	case tea.KeyMsg:
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
			if len(h.hosts) > 0 {
				selected := h.hosts[h.cursor]
				return h, func() tea.Msg { return HostSelectedMsg{Host: selected} }
			}
		case keys.Match(msg, keys.HostList.Back):
			return h, func() tea.Msg { return BackMsg{} }
		case keys.Match(msg, keys.GroupList.Quit):
			return h, tea.Quit
		}
	}
	return h, nil
}

func (h HostList) View() string {
	width := h.width
	if width <= 0 {
		width = 120
	}

	header := styles.RenderHeader(
		"IMITERM",
		h.group.Name,
		styles.RenderPill("hosts", fmt.Sprintf("%d", len(h.hosts)), true),
		styles.RenderPill("group", "open", false),
	)
	if h.group.Description != "" {
		header += "\n" + styles.HeaderSubtitle.Render("  " + h.group.Description)
	}

	leftWidth := max(44, min(58, width/3))
	rightWidth := max(34, width-leftWidth-6)

	leftBody := h.renderHostListBody()
	rightBody := h.renderHostDetailBody()

	leftPanel := styles.RenderCard("Hosts", leftBody, leftWidth, true)
	rightPanel := styles.RenderCard("Host details", rightBody, rightWidth, false)
	content := styles.RenderDualCards(leftPanel, rightPanel, "  ")
	if width < 100 {
		content = leftPanel + "\n\n" + rightPanel
	}

	footer := styles.Pill.Render("enter connect") + "  " + styles.Pill.Render("a add") + "  " + styles.Pill.Render("e edit") + "  " + styles.Pill.Render("d delete") + "  " + styles.Pill.Render("b back")

	help := styles.RenderHelp(
		"enter", "connect", "esc/b", "back",
		"a", "add", "e", "edit", "d", "delete", "q", "quit",
	)

	return header + "\n" + content + "\n" + footer + "\n" + help + "\n"
}

func (h HostList) renderHostListBody() string {
	if len(h.hosts) == 0 {
		return styles.RenderEmpty("No hosts here", "Press a to add one, or b to return to groups.")
	}

	var rows []string
	for i, host := range h.hosts {
		rows = append(rows, renderHostRow(host, i == h.cursor, ""))
	}
	return strings.Join(rows, "\n")
}

func (h HostList) renderHostDetailBody() string {
	if len(h.hosts) == 0 {
		return styles.RenderEmpty("Get started", "Create a host entry to connect faster.")
	}

	host := h.hosts[h.cursor]
	var b strings.Builder
	b.WriteString(styles.DetailLabel.Render("Name") + styles.DetailValue.Render(host.Name) + "\n")
	b.WriteString(styles.DetailLabel.Render("Endpoint") + styles.DetailValue.Render(fmt.Sprintf("%s@%s:%d", host.User, host.Hostname, host.Port)) + "\n")
	b.WriteString(styles.DetailLabel.Render("Auth") + styles.RenderAuthBadge(authKind(host)) + "\n")
	if host.KeyPath != "" {
		b.WriteString(styles.DetailLabel.Render("Key") + styles.DetailValue.Render(host.KeyPath) + "\n")
	} else if host.Password != "" {
		b.WriteString(styles.DetailLabel.Render("Password") + styles.PanelSubtle.Render("stored locally") + "\n")
	} else {
		b.WriteString(styles.DetailLabel.Render("Auth") + styles.PanelSubtle.Render("ssh-agent / default key") + "\n")
	}
	if len(host.Tags) > 0 {
		b.WriteString("\n" + styles.SectionTitle.Render("Tags") + "\n")
		for _, tag := range host.Tags {
			b.WriteString(styles.GroupTag.Render(tag) + " ")
		}
	}
	return strings.TrimSpace(b.String())
}

func authKind(host model.Host) string {
	switch host.AuthMethod() {
	case model.AuthKey:
		return "key"
	case model.AuthPass:
		return "pass"
	default:
		return "agent"
	}
}

func renderHostRow(host model.Host, selected bool, groupName string) string {
	name := styles.RowTitle.Render(host.Name)

	marker := styles.PanelSubtle.Render("•")
	rowStyle := styles.PassiveRow
	if selected {
		marker = styles.SelectedMarker.Render("›")
		rowStyle = styles.ActiveRow
	}

	return rowStyle.Render(marker + "  " + name)
}

package screens

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/ParsaImi/imiterm/internal/model"
	"github.com/ParsaImi/imiterm/internal/ui/keys"
	"github.com/ParsaImi/imiterm/internal/ui/styles"
)

type searchResult struct {
	host      model.Host
	groupName string
}

type Search struct {
	input   textinput.Model
	all     []searchResult
	matched []searchResult
	cursor  int
	width   int
	height  int
}

func NewSearch(groups []model.Group) Search {
	input := textinput.New()
	input.Placeholder = "type to filter hosts..."
	input.Prompt = "  / "
	input.Focus()
	input.PromptStyle = styles.HelpKey
	input.TextStyle = styles.HostAddr

	var all []searchResult
	for _, g := range groups {
		for _, h := range g.Hosts {
			all = append(all, searchResult{host: h, groupName: g.Name})
		}
	}

	return Search{input: input, all: all, matched: all}
}

func (s Search) Init() tea.Cmd { return textinput.Blink }

func (s Search) Update(msg tea.Msg) (Search, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.width = msg.Width
		s.height = msg.Height
		return s, nil
	}

	if msg, ok := msg.(tea.KeyMsg); ok {
		switch {
		case keys.Match(msg, keys.GroupList.Quit):
			return s, tea.Quit
		case msg.String() == "esc" || msg.String() == "ctrl+b":
			return s, func() tea.Msg { return BackMsg{} }
		case msg.String() == "up":
			if s.cursor > 0 {
				s.cursor--
			}
			return s, nil
		case msg.String() == "down":
			if s.cursor < len(s.matched)-1 {
				s.cursor++
			}
			return s, nil
		case keys.Match(msg, keys.GroupList.Enter):
			if len(s.matched) > 0 {
				host := s.matched[s.cursor].host
				return s, func() tea.Msg { return HostSelectedMsg{Host: host} }
			}
			return s, nil
		}
	}

	var cmd tea.Cmd
	s.input, cmd = s.input.Update(msg)
	s.filter()
	return s, cmd
}

func (s *Search) filter() {
	query := strings.ToLower(strings.TrimSpace(s.input.Value()))
	if query == "" {
		s.matched = s.all
		s.cursor = 0
		return
	}

	var matched []searchResult
	for _, r := range s.all {
		searchable := strings.ToLower(
			r.host.Name + " " + r.host.Hostname + " " + r.host.User + " " +
				r.groupName + " " + strings.Join(r.host.Tags, " "),
		)
		if strings.Contains(searchable, query) {
			matched = append(matched, r)
		}
	}

	s.matched = matched
	if s.cursor >= len(s.matched) {
		s.cursor = max(0, len(s.matched)-1)
	}
}

func (s Search) View() string {
	width := s.width
	if width <= 0 {
		width = 120
	}
	height := s.height
	if height <= 0 {
		height = 36
	}

	header := styles.RenderHeader(
		"IMITERM",
		"Search across all hosts",
		styles.RenderPill("matches", fmt.Sprintf("%d", len(s.matched)), true),
		styles.RenderPill("query", strings.TrimSpace(s.input.Value()), false),
	)

	searchBar := styles.RenderCard("Search", styles.SearchBox.Render(s.input.View()), min(width-8, 100), true)

	leftWidth := max(44, min(58, width/3))
	rightWidth := max(34, width-leftWidth-6)
	maxRows := max(3, min(8, height-12))

	var rows []string
	if len(s.matched) == 0 {
		rows = append(rows, styles.RenderEmpty("No matches", "Try host, hostname, user, group, or tags."))
	} else {
		limit := min(len(s.matched), maxRows)
		for i := 0; i < limit; i++ {
			r := s.matched[i]
			rows = append(rows, renderHostRow(r.host, i == s.cursor, r.groupName))
		}
		if len(s.matched) > limit {
			rows = append(rows, styles.PanelSubtle.Render(fmt.Sprintf("… %d more results", len(s.matched)-limit)))
		}
	}
	resultsBody := strings.Join(rows, "\n")

	content := styles.RenderCard("Results", resultsBody, leftWidth, false)
	if width >= 100 && height >= 30 {
		content = styles.RenderDualCards(
			styles.RenderCard("Results", resultsBody, leftWidth, false),
			s.renderDetail(rightWidth),
			"  ",
		)
	}
	if width < 100 {
		content = styles.RenderCard("Results", resultsBody, width-8, false)
	}
	footer := styles.Pill.Render("enter connect") + "  " + styles.Pill.Render("esc back") + "  " + styles.Pill.Render("↑ down") + "  " + styles.Pill.Render("↓ up")

	return header + "\n" + searchBar + "\n" + content + "\n" + footer + "\n"
}

func (s Search) renderDetail(width int) string {
	if len(s.matched) == 0 {
		return styles.RenderCard("Inspect", styles.RenderEmpty("No selection", "Type to narrow the host list."), width, false)
	}

	r := s.matched[s.cursor]
	var b strings.Builder
	b.WriteString(styles.DetailLabel.Render("Host") + styles.DetailValue.Render(r.host.Name) + "\n")
	b.WriteString(styles.DetailLabel.Render("Group") + styles.DetailValue.Render(r.groupName) + "\n")
	b.WriteString(styles.DetailLabel.Render("Endpoint") + styles.DetailValue.Render(fmt.Sprintf("%s@%s:%d", r.host.User, r.host.Hostname, r.host.Port)) + "\n")
	b.WriteString(styles.DetailLabel.Render("Tags") + styles.DetailValue.Render(strings.Join(r.host.Tags, ", ")))
	return styles.RenderCard("Inspect", b.String(), width, false)
}

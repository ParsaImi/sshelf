package styles

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Dracula-inspired palette
var (
	colorPink    = lipgloss.Color("#FF79C6")
	colorGreen   = lipgloss.Color("#50FA7B")
	colorPurple  = lipgloss.Color("#BD93F9")
	colorCyan    = lipgloss.Color("#8BE9FD")
	colorOrange  = lipgloss.Color("#FFB86C")
	colorRed     = lipgloss.Color("#FF5555")
	colorYellow  = lipgloss.Color("#F1FA8C")
	colorComment = lipgloss.Color("#6272A4")
	colorFg      = lipgloss.Color("#F8F8F2")
	colorBg      = lipgloss.Color("#282A36")
	colorBgAlt   = lipgloss.Color("#44475A")
)

// Header / title bar
var Header = lipgloss.NewStyle().
	Bold(true).
	Foreground(colorBg).
	Background(colorPurple).
	Padding(0, 1)

var HeaderSubtitle = lipgloss.NewStyle().
	Foreground(colorComment).
	Italic(true)

var ScreenMeta = lipgloss.NewStyle().
	Foreground(colorComment)

// Panel / card surfaces
var ContentBox = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderForeground(colorComment).
	Padding(0, 1).
	MarginLeft(0).
	MarginRight(0)

var Card = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderForeground(colorComment).
	Padding(1, 2)

var CardAccent = Card.Copy().BorderForeground(colorPurple)

var CardSoft = Card.Copy().BorderForeground(colorBgAlt)

var CardWarn = Card.Copy().BorderForeground(colorOrange)

var PanelTitle = lipgloss.NewStyle().
	Foreground(colorCyan).
	Bold(true)

var SectionTitle = lipgloss.NewStyle().
	Foreground(colorPurple).
	Bold(true)

var PanelSubtle = lipgloss.NewStyle().
	Foreground(colorComment)

var ActiveRow = lipgloss.NewStyle().
	Foreground(colorFg).
	Bold(true).
	Padding(0, 1)

var SelectedMarker = lipgloss.NewStyle().
	Foreground(colorGreen).
	Bold(true)

var PassiveRow = lipgloss.NewStyle().
	Foreground(colorFg).
	Padding(0, 1)

var RowCard = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderForeground(colorBgAlt).
	Padding(0, 1).
	MarginBottom(1)

var RowCardSelected = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderForeground(colorGreen).
	Padding(0, 1).
	MarginBottom(1).
	Background(colorBgAlt)

var RowTitle = lipgloss.NewStyle().
	Foreground(colorFg).
	Bold(true)

var RowMeta = lipgloss.NewStyle().
	Foreground(colorComment)

var RowLabel = lipgloss.NewStyle().
	Foreground(colorComment).
	Width(10)

var SubtleRow = lipgloss.NewStyle().
	Foreground(colorComment).
	Padding(0, 1)

var DetailLabel = lipgloss.NewStyle().
	Foreground(colorComment).
	Width(12)

var DetailValue = lipgloss.NewStyle().
	Foreground(colorFg)

var Pill = lipgloss.NewStyle().
	Foreground(colorBg).
	Background(colorComment).
	Padding(0, 1).
	Bold(true)

var PillAccent = lipgloss.NewStyle().
	Foreground(colorBg).
	Background(colorPurple).
	Padding(0, 1).
	Bold(true)

var PillSuccess = lipgloss.NewStyle().
	Foreground(colorBg).
	Background(colorGreen).
	Padding(0, 1).
	Bold(true)

var PillWarn = lipgloss.NewStyle().
	Foreground(colorBg).
	Background(colorOrange).
	Padding(0, 1).
	Bold(true)

var EmptyState = lipgloss.NewStyle().
	Foreground(colorComment).
	Italic(true).
	PaddingLeft(2)

// List items
var ListItem = lipgloss.NewStyle().
	PaddingLeft(2)

var ListItemSelected = lipgloss.NewStyle().
	PaddingLeft(1).
	Foreground(colorGreen).
	Bold(true)

var ListItemDim = lipgloss.NewStyle().
	Foreground(colorComment)

// Host details within a list row
var HostName = lipgloss.NewStyle().
	Foreground(colorCyan).
	Bold(true)

var HostAddr = lipgloss.NewStyle().
	Foreground(colorFg)

var AuthBadge = lipgloss.NewStyle().
	Foreground(colorBg).
	Background(colorPurple).
	Padding(0, 1).
	Bold(true)

var AuthBadgePass = lipgloss.NewStyle().
	Foreground(colorBg).
	Background(colorOrange).
	Padding(0, 1).
	Bold(true)

var AuthBadgeAgent = lipgloss.NewStyle().
	Foreground(colorBg).
	Background(colorComment).
	Padding(0, 1).
	Bold(true)

// Group styling
var GroupName = lipgloss.NewStyle().
	Foreground(colorFg).
	Bold(true)

var GroupCount = lipgloss.NewStyle().
	Foreground(colorComment)

var GroupDesc = lipgloss.NewStyle().
	Foreground(colorComment).
	Italic(true).
	PaddingLeft(6)

var GroupTag = lipgloss.NewStyle().
	Foreground(colorBg).
	Background(colorCyan).
	Padding(0, 1)

// Form styling
var FormBox = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderForeground(colorPurple).
	Padding(1, 2).
	MarginLeft(1).
	MarginRight(1)

var FormLabel = lipgloss.NewStyle().
	Foreground(colorPurple).
	Bold(true).
	Width(14).
	Align(lipgloss.Right).
	PaddingRight(1)

var FormError = lipgloss.NewStyle().
	Foreground(colorRed).
	Bold(true).
	PaddingLeft(2)

// Help bar at the bottom
var HelpBar = lipgloss.NewStyle().
	Foreground(colorComment).
	PaddingLeft(1).
	MarginTop(1)

var HelpKey = lipgloss.NewStyle().
	Foreground(colorPink).
	Bold(true)

var HelpSep = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#44475A"))

// Status bar
var StatusOk = lipgloss.NewStyle().
	Foreground(colorGreen).
	PaddingLeft(1)

var StatusErr = lipgloss.NewStyle().
	Foreground(colorRed).
	PaddingLeft(1)

// Confirm dialog
var ConfirmBox = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderForeground(colorRed).
	Padding(1, 3).
	MarginLeft(1).
	MarginTop(1)

var ConfirmPrompt = lipgloss.NewStyle().
	Foreground(colorOrange).
	Bold(true)

// Search
var SearchBox = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderForeground(colorCyan).
	Padding(0, 1).
	MarginLeft(1).
	MarginRight(1)

var SearchNoMatch = lipgloss.NewStyle().
	Foreground(colorComment).
	Italic(true).
	PaddingLeft(2)

// Helpers

func RenderHelp(pairs ...string) string {
	var parts []string
	for i := 0; i < len(pairs)-1; i += 2 {
		key := HelpKey.Render(pairs[i])
		desc := lipgloss.NewStyle().Foreground(colorComment).Render(pairs[i+1])
		parts = append(parts, key+" "+desc)
	}

	sep := HelpSep.Render(" | ")
	result := ""
	for i, part := range parts {
		if i > 0 {
			result += sep
		}
		result += part
	}
	return HelpBar.Render(result)
}

func RenderHeader(title, subtitle string, meta ...string) string {
	line := Header.Render(title)
	if subtitle != "" {
		line += "  " + HeaderSubtitle.Render(subtitle)
	}
	if len(meta) > 0 {
		line += "  " + ScreenMeta.Render(strings.Join(meta, "  "))
	}
	return line
}

func RenderPill(label, value string, accent bool) string {
	content := label + ": " + value
	if accent {
		return PillAccent.Render(content)
	}
	return Pill.Render(content)
}

func RenderEmpty(title, body string) string {
	return EmptyState.Render(title + "\n" + body)
}

func CenterBlock(view string, width int) string {
	if width <= 0 {
		return view
	}
	return lipgloss.PlaceHorizontal(width, lipgloss.Center, strings.TrimRight(view, "\n"))
}

func RenderViewport(view string, width, height int) string {
	trimmed := strings.TrimRight(view, "\n")
	if width <= 0 {
		return trimmed
	}
	if height <= 0 {
		return lipgloss.PlaceHorizontal(width, lipgloss.Center, trimmed)
	}
	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, trimmed)
}

func RenderCard(title, body string, width int, accent bool) string {
	panel := Card
	if accent {
		panel = CardAccent
	}
	if width > 0 {
		panel = panel.Width(width)
	}
	content := body
	if title != "" {
		content = SectionTitle.Render(title) + "\n" + body
	}
	return panel.Render(content)
}

func RenderDualCards(left, right string, gap string) string {
	return lipgloss.JoinHorizontal(lipgloss.Top, left, gap, right)
}

func RenderAuthBadge(kind string) string {
	switch kind {
	case "key":
		return AuthBadge.Render("KEY")
	case "pass":
		return AuthBadgePass.Render("PASS")
	default:
		return AuthBadgeAgent.Render("AGENT")
	}
}

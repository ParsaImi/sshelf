package styles

import (
	"fmt"
	"strconv"
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
	colorBgAlt   = lipgloss.Color("#3A3D53")
)

// ContentWidth is the interior width shared by every card/box so screens
// line up edge-to-edge instead of each auto-sizing to its own content.
const ContentWidth = 70

// RowWidth is the usable width inside a ContentBox row (accounts for the
// box's own 1-col padding on each side) — rows are padded to this so the
// selection highlight fills the entire line, not just the text.
const RowWidth = ContentWidth - 2

// ---------------------------------------------------------------------
// Gradient text
// ---------------------------------------------------------------------

func hexToRGB(hex string) (int, int, int) {
	hex = strings.TrimPrefix(hex, "#")
	r, _ := strconv.ParseInt(hex[0:2], 16, 0)
	g, _ := strconv.ParseInt(hex[2:4], 16, 0)
	b, _ := strconv.ParseInt(hex[4:6], 16, 0)
	return int(r), int(g), int(b)
}

func lerp(a, b int, t float64) int {
	return int(float64(a) + (float64(b)-float64(a))*t)
}

func blend(c1, c2 lipgloss.Color, t float64) lipgloss.Color {
	r1, g1, b1 := hexToRGB(string(c1))
	r2, g2, b2 := hexToRGB(string(c2))
	return lipgloss.Color(fmt.Sprintf("#%02x%02x%02x", lerp(r1, r2, t), lerp(g1, g2, t), lerp(b1, b2, t)))
}

// GradientText renders s with its foreground color interpolated across the
// given stops, left to right. Spaces are passed through unstyled.
func GradientText(s string, stops ...lipgloss.Color) string {
	runes := []rune(s)
	n := len(runes)
	if n == 0 {
		return s
	}
	if len(stops) == 1 {
		return lipgloss.NewStyle().Foreground(stops[0]).Bold(true).Render(s)
	}

	segments := len(stops) - 1
	var b strings.Builder
	for i, r := range runes {
		if r == ' ' {
			b.WriteRune(r)
			continue
		}
		pos := float64(i) / float64(max(n-1, 1)) * float64(segments)
		seg := min(int(pos), segments-1)
		localT := pos - float64(seg)
		c := blend(stops[seg], stops[seg+1], localT)
		b.WriteString(lipgloss.NewStyle().Foreground(c).Bold(true).Render(string(r)))
	}
	return b.String()
}

var logoStops = []lipgloss.Color{colorPurple, colorPink, colorOrange}

// Logo returns the gradient SSHELF wordmark used in the banner.
func Logo() string {
	return GradientText("◆ SSHELF", logoStops...)
}

// ---------------------------------------------------------------------
// Banner (top of every screen)
// ---------------------------------------------------------------------

var bannerBox = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderForeground(colorPurple).
	Padding(0, 1).
	Width(ContentWidth).
	MarginLeft(1).
	MarginBottom(1)

var bannerSubtitle = lipgloss.NewStyle().
	Foreground(colorComment).
	Italic(true)

var bannerMeta = lipgloss.NewStyle().
	Foreground(colorCyan)

// RenderBanner draws the app wordmark plus a screen subtitle, with an
// optional right-aligned meta string (stats, counts, etc.) on the logo line.
func RenderBanner(subtitle, meta string) string {
	logo := Logo()
	line := logo
	if meta != "" {
		metaStyled := bannerMeta.Render(meta)
		gap := RowWidth - lipgloss.Width(logo) - lipgloss.Width(metaStyled)
		if gap < 1 {
			gap = 1
		}
		line = logo + strings.Repeat(" ", gap) + metaStyled
	}

	body := line
	if subtitle != "" {
		// Truncate the raw text (not the styled string — cutting a styled
		// string mid-ANSI-sequence would corrupt terminal color state) so a
		// long user-entered group name can't trigger Style.Width()'s
		// word-wrap and break the box layout.
		if lipgloss.Width(subtitle) > RowWidth {
			r := []rune(subtitle)
			cut := max(RowWidth-1, 1)
			if cut > 3 {
				subtitle = string(r[:cut-3]) + "..."
			} else {
				subtitle = string(r[:cut])
			}
		}
		body += "\n" + bannerSubtitle.Render(subtitle)
	}

	return bannerBox.Render(body)
}

// ---------------------------------------------------------------------
// Content card
// ---------------------------------------------------------------------

var ContentBox = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderForeground(colorComment).
	Padding(0, 1).
	MarginLeft(1).
	MarginBottom(1)

// ---------------------------------------------------------------------
// List rows
// ---------------------------------------------------------------------

var ListItemDim = lipgloss.NewStyle().
	Foreground(colorComment)

var ColumnHeader = lipgloss.NewStyle().
	Foreground(colorComment).
	Bold(true).
	PaddingLeft(2)

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

// ---------------------------------------------------------------------
// Forms
// ---------------------------------------------------------------------

var FormBox = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderForeground(colorPurple).
	Padding(1, 2).
	Width(ContentWidth - 2).
	MarginLeft(1).
	MarginBottom(1)

var FormLabel = lipgloss.NewStyle().
	Foreground(colorPurple).
	Bold(true).
	Width(12).
	Align(lipgloss.Right).
	PaddingRight(1)

var FormLabelActive = lipgloss.NewStyle().
	Foreground(colorGreen).
	Bold(true).
	Width(12).
	Align(lipgloss.Right).
	PaddingRight(1)

var FormError = lipgloss.NewStyle().
	Foreground(colorBg).
	Background(colorRed).
	Bold(true).
	Padding(0, 1)

var FocusIndicator = lipgloss.NewStyle().Foreground(colorGreen).Bold(true)

// ---------------------------------------------------------------------
// Help / footer bar
// ---------------------------------------------------------------------

var helpBox = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderForeground(colorComment).
	Padding(0, 1).
	Width(ContentWidth).
	MarginLeft(1)

var HelpKey = lipgloss.NewStyle().
	Foreground(colorPink).
	Bold(true)

var helpDesc = lipgloss.NewStyle().
	Foreground(colorComment)

// helpKeyColors gives each key hint a color matched to what it does —
// destructive actions read as red, creation as green, navigation as cyan —
// instead of every key being the same flat pink. Falls back to colorPink.
var helpKeyColors = map[string]lipgloss.Color{
	"↑↓":        colorCyan,
	"←→":        colorCyan,
	"↵":         colorGreen,
	"tab":       colorCyan,
	"shift+tab": colorCyan,
	"esc":       colorComment,
	"a":         colorGreen,
	"e":         colorYellow,
	"d":         colorRed,
	"/":         colorPurple,
	"s":         colorCyan,
	"q":         colorComment,
}

// Deprecated compatibility shim — no longer used directly for rendering,
// kept in case other code references the exported symbol.
var HelpBar = lipgloss.NewStyle().Foreground(colorComment)

var HelpSep = lipgloss.NewStyle().Foreground(colorBgAlt)

// ---------------------------------------------------------------------
// Status bar
// ---------------------------------------------------------------------

var StatusOk = lipgloss.NewStyle().
	Foreground(colorBg).
	Background(colorGreen).
	Bold(true).
	Padding(0, 1).
	MarginLeft(1)

var StatusErr = lipgloss.NewStyle().
	Foreground(colorBg).
	Background(colorRed).
	Bold(true).
	Padding(0, 1).
	MarginLeft(1)

// ---------------------------------------------------------------------
// Confirm dialog
// ---------------------------------------------------------------------

var ConfirmBox = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderForeground(colorRed).
	Padding(1, 3).
	Width(ContentWidth - 4).
	MarginLeft(1).
	MarginTop(1)

var ConfirmPrompt = lipgloss.NewStyle().
	Foreground(colorOrange).
	Bold(true)

// ---------------------------------------------------------------------
// Search
// ---------------------------------------------------------------------

var SearchBox = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderForeground(colorCyan).
	Padding(0, 1).
	Width(ContentWidth - 2).
	MarginLeft(1).
	MarginBottom(1)

var SearchNoMatch = lipgloss.NewStyle().
	Foreground(colorComment).
	Italic(true).
	PaddingLeft(2)

// ---------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------

// RenderHelp draws the bottom key-hint bar inside a bordered box, wrapping
// onto additional lines when the hints don't fit — always at hint
// boundaries, never mid-hint (Style.Width() word-wraps on spaces, which
// would otherwise split a key from its label).
func RenderHelp(pairs ...string) string {
	const sep = "  "

	var lines []string
	var cur strings.Builder
	curWidth := 0

	for i := 0; i < len(pairs)-1; i += 2 {
		fg, ok := helpKeyColors[pairs[i]]
		if !ok {
			fg = colorPink
		}
		keyStyle := lipgloss.NewStyle().Foreground(fg).Bold(true)
		hint := keyStyle.Render(pairs[i]) + " " + helpDesc.Render(pairs[i+1])
		hintWidth := lipgloss.Width(hint)

		if curWidth > 0 && curWidth+len(sep)+hintWidth > RowWidth {
			lines = append(lines, cur.String())
			cur.Reset()
			curWidth = 0
		}
		if curWidth > 0 {
			cur.WriteString(sep)
			curWidth += len(sep)
		}
		cur.WriteString(hint)
		curWidth += hintWidth
	}
	if curWidth > 0 {
		lines = append(lines, cur.String())
	}

	return helpBox.Render(strings.Join(lines, "\n"))
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

// Divider draws a thin horizontal rule the width of a row.
func Divider() string {
	return lipgloss.NewStyle().Foreground(colorBgAlt).Render(strings.Repeat("─", RowWidth))
}

// Accent returns a single-cell colored block, used as a left-edge accent
// bar for group rows (colored per-group) and other tinted list items.
func Accent(hex string) string {
	return lipgloss.NewStyle().Foreground(lipgloss.Color(hex)).Render("▍")
}

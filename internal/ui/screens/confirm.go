package screens

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/ParsaImi/imiterm/internal/ui/styles"
)

type ConfirmResult struct {
	Confirmed bool
}

type Confirm struct {
	message string
	width   int
	height  int
}

func NewConfirm(message string) Confirm {
	return Confirm{message: message}
}

func (c Confirm) Init() tea.Cmd { return nil }

func (c Confirm) Update(msg tea.Msg) (Confirm, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		c.width = msg.Width
		c.height = msg.Height
		return c, nil
	}

	if msg, ok := msg.(tea.KeyMsg); ok {
		switch msg.String() {
		case "y", "Y":
			return c, func() tea.Msg { return ConfirmResult{Confirmed: true} }
		case "n", "N", "esc":
			return c, func() tea.Msg { return ConfirmResult{Confirmed: false} }
		}
	}
	return c, nil
}

func (c Confirm) View() string {
	width := c.width
	if width <= 0 {
		width = 92
	}
	header := styles.RenderHeader(
		"IMITERM",
		"Confirm action",
		styles.RenderPill("danger", "delete", true),
	)

	prompt := styles.ConfirmPrompt.Render(c.message)
	hint := styles.PanelSubtle.Render("Press y to confirm or n/esc to cancel")
	box := styles.RenderCard("Danger zone", prompt+"\n\n"+hint, min(width-8, 64), true)
	return header + "\n" + box + "\n"
}

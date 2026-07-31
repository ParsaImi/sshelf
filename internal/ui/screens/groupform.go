package screens

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/ParsaImi/imiterm/internal/model"
	"github.com/ParsaImi/imiterm/internal/ui/styles"
)

type GroupSavedMsg struct {
	Group   model.Group
	IsNew   bool
	EditIdx int
}

type presetColor struct {
	name string
	hex  string
}

var colorPresets = []presetColor{
	{"Red", "#FF5555"},
	{"Orange", "#FFB86C"},
	{"Yellow", "#F1FA8C"},
	{"Green", "#50FA7B"},
	{"Cyan", "#8BE9FD"},
	{"Blue", "#6272A4"},
	{"Purple", "#BD93F9"},
	{"Pink", "#FF79C6"},
}

const (
	gfName = iota
	gfDesc
	gfFieldCount
)

type GroupForm struct {
	inputs      []textinput.Model
	focused     int
	onColorRow  bool
	colorIdx    int
	isNew       bool
	editIdx     int
	err         string
	width       int
	height      int
}

func NewGroupForm() GroupForm {
	return newGroupForm(model.Group{Color: "#6272A4"}, true, -1)
}

func NewGroupFormEdit(g model.Group, idx int) GroupForm {
	return newGroupForm(g, false, idx)
}

func newGroupForm(g model.Group, isNew bool, editIdx int) GroupForm {
	inputs := make([]textinput.Model, gfFieldCount)

	inputs[gfName] = textinput.New()
	inputs[gfName].Placeholder = "Production"
	inputs[gfName].Prompt = ""
	inputs[gfName].Width = 40
	inputs[gfName].CharLimit = 256
	inputs[gfName].SetValue(g.Name)

	inputs[gfDesc] = textinput.New()
	inputs[gfDesc].Placeholder = "Live servers"
	inputs[gfDesc].Prompt = ""
	inputs[gfDesc].Width = 40
	inputs[gfDesc].CharLimit = 256
	inputs[gfDesc].SetValue(g.Description)

	inputs[gfName].Focus()

	colorIdx := 5
	for i, c := range colorPresets {
		if c.hex == g.Color {
			colorIdx = i
			break
		}
	}

	return GroupForm{
		inputs:   inputs,
		focused:  0,
		colorIdx: colorIdx,
		isNew:    isNew,
		editIdx:  editIdx,
	}
}

func (f GroupForm) Init() tea.Cmd { return textinput.Blink }

func (f GroupForm) Update(msg tea.Msg) (GroupForm, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		f.width = msg.Width
		f.height = msg.Height
		return f, nil
	}

	if msg, ok := msg.(tea.KeyMsg); ok {
		switch msg.String() {
		case "tab", "down":
			f.focusNext()
			return f, nil
		case "shift+tab", "up":
			f.focusPrev()
			return f, nil
		case "left", "h":
			if f.onColorRow {
				if f.colorIdx > 0 {
					f.colorIdx--
				}
				return f, nil
			}
		case "right", "l":
			if f.onColorRow {
				if f.colorIdx < len(colorPresets)-1 {
					f.colorIdx++
				}
				return f, nil
			}
		case "ctrl+s", "enter":
			if f.onColorRow {
				group, errMsg := f.toGroup()
				if errMsg != "" {
					f.err = errMsg
					return f, nil
				}
				return f, func() tea.Msg {
					return GroupSavedMsg{Group: group, IsNew: f.isNew, EditIdx: f.editIdx}
				}
			}
			f.focusNext()
			return f, nil
		case "esc", "ctrl+b":
			return f, func() tea.Msg { return BackMsg{} }
		}
	}

	if !f.onColorRow {
		var cmd tea.Cmd
		f.inputs[f.focused], cmd = f.inputs[f.focused].Update(msg)
		return f, cmd
	}

	return f, nil
}

func (f *GroupForm) focusNext() {
	if f.onColorRow {
		group, errMsg := f.toGroup()
		if errMsg != "" {
			f.err = errMsg
			return
		}
		_ = group
		return
	}

	f.inputs[f.focused].Blur()
	f.focused++

	if f.focused >= gfFieldCount {
		f.onColorRow = true
	} else {
		f.inputs[f.focused].Focus()
	}
}

func (f *GroupForm) focusPrev() {
	if f.onColorRow {
		f.onColorRow = false
		f.focused = gfFieldCount - 1
		f.inputs[f.focused].Focus()
		return
	}

	if f.focused > 0 {
		f.inputs[f.focused].Blur()
		f.focused--
		f.inputs[f.focused].Focus()
	}
}

func (f GroupForm) toGroup() (model.Group, string) {
	name := strings.TrimSpace(f.inputs[gfName].Value())
	if name == "" {
		return model.Group{}, "name is required"
	}

	return model.Group{
		Name:        name,
		Description: strings.TrimSpace(f.inputs[gfDesc].Value()),
		Color:       colorPresets[f.colorIdx].hex,
	}, ""
}

func (f GroupForm) View() string {
	title := "Add Group"
	if !f.isNew {
		title = "Edit Group"
	}
	width := f.width
	if width <= 0 {
		width = 120
	}
	header := styles.RenderHeader(
		"IMITERM",
		title,
		styles.RenderPill("colors", fmt.Sprintf("%d", len(colorPresets)), true),
		styles.RenderPill("mode", "group", false),
	)

	leftBody := f.renderGroupFormFields()
	rightBody := f.renderGroupPreview()
	leftWidth := max(44, min(60, width/2-4))
	rightWidth := max(28, width-leftWidth-8)
	content := styles.RenderDualCards(
		styles.RenderCard("Details", leftBody, leftWidth, true),
		styles.RenderCard("Preview", rightBody, rightWidth, false),
		"  ",
	)
	if width < 100 {
		content = styles.RenderCard("Details", leftBody, width-8, true) + "\n\n" + styles.RenderCard("Preview", rightBody, width-8, false)
	}
	footer := styles.Pill.Render("tab next") + "  " + styles.Pill.Render("shift+tab prev") + "  " + styles.Pill.Render("arrows color") + "  " + styles.Pill.Render("enter save") + "  " + styles.Pill.Render("esc cancel")

	help := styles.RenderHelp(
		"tab/↓", "next", "shift+tab/↑", "prev",
		"←/→", "color", "enter", "save", "esc", "cancel",
	)

	return header + "\n" + content + "\n" + footer + "\n" + help + "\n"
}

func (f GroupForm) renderGroupFormFields() string {
	labels := [gfFieldCount]string{"Name", "Description"}
	var b strings.Builder
	b.WriteString(styles.PanelSubtle.Render("Pick a name, add a description, then choose an accent color.") + "\n\n")
	for i := range f.inputs {
		selected := i == f.focused && !f.onColorRow
		label := styles.DetailLabel.Render(labels[i])
		line := label + f.inputs[i].View()
		if selected {
			line = styles.ActiveRow.Render(styles.SelectedMarker.Render("›") + " " + strings.TrimSpace(line))
		} else {
			line = styles.PassiveRow.Render("  " + strings.TrimSpace(line))
		}
		b.WriteString(line + "\n")
	}

	b.WriteString("\n" + styles.DetailLabel.Render("Color") + "\n")
	var chips []string
	for i, c := range colorPresets {
		chip := lipgloss.NewStyle().Foreground(lipgloss.Color(c.hex)).Render("● " + c.name)
		if i == f.colorIdx {
			chips = append(chips, lipgloss.NewStyle().Foreground(lipgloss.Color(c.hex)).Background(lipgloss.Color("#44475A")).Bold(true).Padding(0, 1).Render(chip))
		} else {
			chips = append(chips, chip)
		}
	}
	if f.onColorRow {
		b.WriteString(styles.ActiveRow.Render(strings.Join(chips, "  ")))
	} else {
		b.WriteString(styles.PassiveRow.Render(strings.Join(chips, "  ")))
	}

	if f.err != "" {
		b.WriteString("\n\n" + styles.FormError.Render("✗ "+f.err))
	}
	return strings.TrimRight(b.String(), "\n")
}

func (f GroupForm) renderGroupPreview() string {
	group, errMsg := f.toGroup()
	if errMsg != "" {
		return styles.RenderEmpty("Incomplete", "Fill the required fields to see a live preview.")
	}

	var b strings.Builder
	b.WriteString(styles.SectionTitle.Render(group.Name) + "\n")
	if group.Description != "" {
		b.WriteString(group.Description + "\n\n")
	} else {
		b.WriteString(styles.PanelSubtle.Render("No description yet") + "\n\n")
	}
	b.WriteString(styles.RenderPill("color", group.Color, true) + "\n\n")
	b.WriteString(styles.SectionTitle.Render("What this group does") + "\n")
	b.WriteString(styles.PanelSubtle.Render("Keep related hosts together, pick a memorable color, and use the list to jump into SSH fast."))
	return strings.TrimRight(b.String(), "\n")
}

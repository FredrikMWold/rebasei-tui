package ui

import (
	"fmt"
	"io"
	"strings"

	list "github.com/charmbracelet/bubbles/v2/list"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"

	"github.com/fredrikmwold/rebasei-tui/internal/ui/theme"
)

// actionOption represents an action in the modal list.
type actionOption struct{ Act action }

func (a actionOption) Title() string       { return string(a.Act) }
func (a actionOption) Description() string { return "" }
func (a actionOption) FilterValue() string { return string(a.Act) }

// openActionModal initializes and opens the action selection modal.
func (m *model) openActionModal() {
	// Build items
	opts := []list.Item{
		actionOption{Act: pick},
		actionOption{Act: squash},
		actionOption{Act: fixup},
		actionOption{Act: edit},
		actionOption{Act: drop},
	}
	// Width: longest option plus small prefix ("> ")
	maxLen := 0
	for _, it := range opts {
		if l := lipgloss.Width(it.FilterValue()); l > maxLen {
			maxLen = l
		}
	}
	width := maxLen + 2
	if width < 12 {
		width = 12
	}
	al := list.New(opts, simpleActionDelegate{}, width, len(opts))
	al.Title = ""
	// Ensure no title or extra top margin/padding in styles
	s := al.Styles
	s.Title = lipgloss.NewStyle()
	s.PaginationStyle = lipgloss.NewStyle()
	s.HelpStyle = lipgloss.NewStyle()
	al.Styles = s
	al.SetShowHelp(false)
	al.SetShowStatusBar(false)
	al.SetShowFilter(false)
	al.SetFilteringEnabled(false)
	al.SetShowPagination(false)

	// Preselect current action
	idx := m.list.Index()
	if idx >= 0 {
		if ci, ok := m.list.Items()[idx].(commitItem); ok {
			switch ci.Act {
			case pick:
				al.Select(0)
			case squash:
				al.Select(1)
			case fixup:
				al.Select(2)
			case edit:
				al.Select(3)
			case drop:
				al.Select(4)
			}
		}
	}

	m.actList = al
	m.modalOpen = true
}

// applySelectedAction applies the selected action from the modal to the current commit.
func (m *model) applySelectedAction() {
	if !m.modalOpen {
		return
	}
	idx := m.list.Index()
	if idx < 0 {
		return
	}
	sel := m.actList.Index()
	switch sel {
	case 0:
		m.setAction(pick)
	case 1:
		m.setAction(squash)
	case 2:
		m.setAction(fixup)
	case 3:
		m.setAction(edit)
	case 4:
		m.setAction(drop)
	}
}

// simpleActionDelegate renders a plain list with "> " for the selected item
type simpleActionDelegate struct{}

func (simpleActionDelegate) Height() int                               { return 1 }
func (simpleActionDelegate) Spacing() int                              { return 0 }
func (simpleActionDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (simpleActionDelegate) Render(w io.Writer, m list.Model, index int, it list.Item) {
	title := it.FilterValue()
	if index == m.Index() {
		line := "> " + title
		fmt.Fprint(w, lipgloss.NewStyle().Foreground(theme.Mauve).Render(line))
		return
	}
	fmt.Fprint(w, "  "+title)
}

// renderActionModal renders the action selection list inside a bordered box using lipgloss compositor.
func (m model) renderActionModal() string {
	// Responsive layout parameters based on terminal size
	sidePad := 1
	if m.width <= 20 {
		sidePad = 0
	}
	// inner content width available inside border and padding
	inner := m.width - 2 - (2 * sidePad)
	if inner < 8 {
		inner = 8
	}
	// Determine what to show given width/height constraints
	showDesc := inner >= 28
	showTitle := m.height >= 9
	showGap := m.height >= 11
	labelPad := 1
	if inner < 16 {
		labelPad = 0
	}

	// Box style for the modal with border
	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(0, sidePad).
		Foreground(theme.Text).
		BorderForeground(theme.Mauve)

	// Title (truncate to inner width if needed)
	titleText := "Choose rebase action"
	if lipgloss.Width(titleText) > inner {
		titleText = truncateToWidth(titleText, inner)
	}
	title := lipgloss.NewStyle().Foreground(theme.Blue).Bold(true).Render(titleText)

	// Manually render the items to avoid any header/extra spacing
	normal := lipgloss.NewStyle().Foreground(theme.Text)
	descStyle := lipgloss.NewStyle().Foreground(theme.Subtext0)
	cursor := lipgloss.NewStyle().Foreground(theme.Mauve)
	idx := m.actList.Index()
	lines := make([]string, 0, len(m.actList.Items()))
	for i, it := range m.actList.Items() {
		ao := it.(actionOption)
		// Build colored label tag for the action
		lblStyle := lipgloss.NewStyle().Padding(0, labelPad)
		switch ao.Act {
		case pick:
			lblStyle = lblStyle.Background(theme.Green).Foreground(theme.Crust)
		case squash:
			lblStyle = lblStyle.Background(theme.Peach).Foreground(theme.Crust)
		case fixup:
			lblStyle = lblStyle.Background(theme.Yellow).Foreground(theme.Crust)
		case edit:
			lblStyle = lblStyle.Background(theme.Sky).Foreground(theme.Crust)
		case drop:
			lblStyle = lblStyle.Background(theme.Red).Foreground(theme.Crust)
		}
		label := lblStyle.Render(string(ao.Act))

		// compute remaining width for description on this line
		prefixWidth := 2 // either "> " or two spaces
		labelWidth := lipgloss.Width(label)
		remaining := inner - prefixWidth - labelWidth - 2 // 2 spaces between label and desc
		var desc string
		if showDesc && remaining > 0 {
			raw := actionDescription(ao.Act)
			if lipgloss.Width(raw) > remaining {
				raw = truncateToWidth(raw, remaining)
			}
			desc = descStyle.Render(raw)
		} else {
			desc = ""
		}

		if i == idx {
			if desc == "" {
				lines = append(lines, cursor.Render(">")+" "+label)
			} else {
				lines = append(lines, cursor.Render(">")+" "+label+"  "+desc)
			}
		} else {
			if desc == "" {
				lines = append(lines, "  "+normal.Render(label))
			} else {
				lines = append(lines, "  "+normal.Render(label)+"  "+desc)
			}
		}
	}

	var content string
	if showTitle {
		if showGap {
			content = title + "\n\n" + strings.Join(lines, "\n")
		} else {
			content = title + "\n" + strings.Join(lines, "\n")
		}
	} else {
		content = strings.Join(lines, "\n")
	}
	return box.Render(content)
}

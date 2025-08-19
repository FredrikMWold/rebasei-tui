package ui

import (
	"fmt"
	"io"

	list "github.com/charmbracelet/bubbles/v2/list"
	"github.com/charmbracelet/lipgloss/v2"

	"github.com/fredrikmwold/rebasei-tui/internal/commands"
	"github.com/fredrikmwold/rebasei-tui/internal/ui/theme"
)

// commitItem holds a commit and its selected action
type commitItem struct {
	Commit commands.Commit
	Act    action
}

func (c commitItem) Title() string { return c.Commit.Subject }
func (c commitItem) Description() string {
	return fmt.Sprintf("%s • %s • %s", c.Commit.HashShort, c.Commit.Author, c.Commit.Date)
}
func (c commitItem) FilterValue() string { return c.Commit.Subject }

// commitDelegate wraps DefaultDelegate and injects an action label before the title
// while preserving default height/spacing and the built-in indicator.
type commitDelegate struct{ list.DefaultDelegate }

type wrappedItem struct {
	base  commitItem
	title string
}

func (w wrappedItem) Title() string       { return w.title }
func (w wrappedItem) Description() string { return w.base.Description() }
func (w wrappedItem) FilterValue() string { return w.base.FilterValue() }

func (d commitDelegate) Render(w io.Writer, m list.Model, index int, it list.Item) {
	if ci, ok := it.(commitItem); ok {
		// Build a colored action label tag with symmetric padding.
		lbl := string(ci.Act)
		style := lipgloss.NewStyle().Padding(0, 1)
		switch ci.Act {
		case pick:
			style = style.Background(theme.Green).Foreground(theme.Crust)
		case squash:
			style = style.Background(theme.Peach).Foreground(theme.Crust)
		case fixup:
			style = style.Background(theme.Yellow).Foreground(theme.Crust)
		case edit:
			style = style.Background(theme.Sky).Foreground(theme.Crust)
		case drop:
			style = style.Background(theme.Red).Foreground(theme.Crust)
		}
		pre := style.Render(lbl) + " "
		subj := ci.Commit.Subject
		if index == m.Index() {
			subj = lipgloss.NewStyle().Foreground(theme.Mauve).Bold(true).Render(subj)
		}
		wi := wrappedItem{base: ci, title: pre + subj}
		d.DefaultDelegate.Render(w, m, index, wi)
		return
	}
	d.DefaultDelegate.Render(w, m, index, it)
}

package ui

import (
	"os"

	key "github.com/charmbracelet/bubbles/v2/key"
	list "github.com/charmbracelet/bubbles/v2/list"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"

	"github.com/fredrikmwold/rebasei-tui/internal/commands"
	"github.com/fredrikmwold/rebasei-tui/internal/ui/theme"
)

type model struct {
	list   list.Model
	status string
	width  int
	height int
	ready  bool
	// when true, quit the TUI and run rebase with the captured actions
	doRebase bool
	actions  []commands.CommitAction

	// modal state for selecting an action via a small list
	modalOpen bool
	actList   list.Model
}

func initialModel() (model, error) {
	commits, err := commands.ListCommits(20)
	items := make([]list.Item, 0, max(0, len(commits)))
	if err == nil {
		for _, c := range commits {
			items = append(items, commitItem{Commit: c, Act: pick})
		}
	}

	// Use a wrapped-item delegate to inject the action label while preserving
	// default list indicators and alignment.
	delegate := commitDelegate{DefaultDelegate: list.NewDefaultDelegate()}
	// Keep selected description text neutral, but tint its indicator (border) Mauve.
	delegate.Styles.SelectedDesc = delegate.Styles.SelectedDesc.
		Foreground(theme.Subtext0).
		BorderForeground(theme.Mauve)
	// Tint the selection indicator to Mauve without affecting description.
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.
		BorderForeground(theme.Mauve)
	delegate.Styles.NormalTitle = delegate.Styles.NormalTitle.Foreground(theme.Text)
	delegate.Styles.NormalDesc = delegate.Styles.NormalDesc.Foreground(theme.Subtext0)

	l := list.New(items, delegate, 0, 0)
	l.Title = "Interactive Rebase"
	l.Styles.Title = lipgloss.NewStyle().Bold(true).Foreground(theme.Blue)
	// Use built-in list help line
	l.SetShowHelp(true)
	l.SetShowStatusBar(false)
	// Disable built-in filtering UI and behavior
	l.SetShowFilter(false)
	// Show built-in paginator at the bottom of the list
	l.SetShowPagination(true)
	// Subtle pagination style to fit the theme
	l.Styles.PaginationStyle = lipgloss.NewStyle().Foreground(theme.Subtext0)
	l.SetFilteringEnabled(false)

	// Provide only our desired help entries using Additional help keys
	// Short help: movement (ctrl+↑/↓), set action (enter), rebase (ctrl+r), quit
	l.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			keys.MoveUp, keys.MoveDown,
			keys.OpenAction, keys.Rebase, keys.Quit,
		}
	}
	// Full help: include movement, action change keys, rebase, quit, and arrow navigation
	l.AdditionalFullHelpKeys = func() []key.Binding {
		upNav := key.NewBinding(key.WithKeys("up"), key.WithHelp("↑", "up"))
		downNav := key.NewBinding(key.WithKeys("down"), key.WithHelp("↓", "down"))
		return []key.Binding{
			upNav, downNav,
			keys.MoveUp, keys.MoveDown,
			keys.OpenAction,
			keys.Pick, keys.Squash, keys.Fixup, keys.Edit, keys.Drop,
			keys.Rebase, keys.Quit,
		}
	}
	// Keep arrow navigation active, but hide it in short help by clearing help labels
	// Remove default arrow bindings entirely to avoid empty help entries; we'll handle arrows manually
	km := l.KeyMap
	km.CursorUp = key.Binding{}
	km.CursorDown = key.Binding{}
	l.KeyMap = km

	status := ""
	if err != nil {
		status = "No commits found or not a Git repo. Open inside a repo to begin."
	}
	return model{list: l, status: status}, nil
}

// Using default list delegate for standard selection highlighting

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.modalOpen {
			// Allow starting rebase directly from modal as well (Ctrl+Enter)
			if key.Matches(msg, keys.Rebase) {
				m.actions = m.collectActions()
				m.doRebase = true
				return m, tea.Quit
			}
			// Route keys to the action list when modal is open
			// Handle confirm/cancel explicitly
			switch msg.String() {
			case "enter":
				m.applySelectedAction()
				m.modalOpen = false
				return m, nil
			case "esc", "q":
				m.modalOpen = false
				return m, nil
			}
			var cmd tea.Cmd
			m.actList, cmd = m.actList.Update(msg)
			return m, cmd
		}
		if key.Matches(msg, keys.Quit) {
			return m, tea.Quit
		}
		if key.Matches(msg, keys.Rebase) {
			// capture current ordering and actions, quit TUI, and run rebase after p.Run returns
			m.actions = m.collectActions()
			m.doRebase = true
			return m, tea.Quit
		}
		if key.Matches(msg, keys.OpenAction) {
			m.openActionModal()
			return m, nil
		}
		if key.Matches(msg, keys.MoveDown) {
			m.moveSelected(1)
			return m, nil
		}
		if key.Matches(msg, keys.MoveUp) {
			m.moveSelected(-1)
			return m, nil
		}
		// Let default list handle plain arrow keys for selection; no reordering here
		// action keys
		if key.Matches(msg, keys.Pick) {
			return m, m.setAction(pick)
		}
		if key.Matches(msg, keys.Squash) {
			return m, m.setAction(squash)
		}
		if key.Matches(msg, keys.Fixup) {
			return m, m.setAction(fixup)
		}
		if key.Matches(msg, keys.Edit) {
			return m, m.setAction(edit)
		}
		if key.Matches(msg, keys.Drop) {
			return m, m.setAction(drop)
		}
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		// Reserve space for status line only when present
		height := msg.Height
		if m.status != "" {
			height -= 1
		}
		if height < 0 {
			height = 0
		}
		// Ensure the list takes the full terminal width for full-row backgrounds
		m.list.SetSize(msg.Width, height)
		m.ready = true
	}

	var cmd tea.Cmd
	if !m.modalOpen {
		if km, ok := msg.(tea.KeyMsg); ok {
			switch km.String() {
			case "up":
				// Move cursor up within the list (no reordering)
				m.list.CursorUp()
				return m, nil
			case "down":
				m.list.CursorDown()
				return m, nil
			}
		}
		m.list, cmd = m.list.Update(msg)
	}
	return m, cmd
}

func (m *model) moveSelected(delta int) {
	idx := m.list.Index()
	items := m.list.Items()
	if idx < 0 || idx >= len(items) {
		return
	}
	newIdx := idx + delta
	if newIdx < 0 || newIdx >= len(items) {
		return
	}
	it := items[idx]
	if delta > 0 {
		copy(items[idx:], items[idx+1:newIdx+1])
	} else {
		copy(items[newIdx+1:idx+1], items[newIdx:idx])
	}
	items[newIdx] = it
	m.list.SetItems(items)
	m.list.Select(newIdx)
}

func (m *model) setAction(a action) tea.Cmd {
	idx := m.list.Index()
	if idx < 0 {
		return nil
	}
	items := m.list.Items()
	ci, ok := items[idx].(commitItem)
	if !ok {
		return nil
	}
	// Disallow squash/fixup on the oldest (last) item to keep valid rebase order.
	if idx == len(items)-1 && (a == squash || a == fixup) {
		m.status = lipgloss.NewStyle().Foreground(theme.Red).Render("Can't squash/fixup the oldest commit. Move it above another or pick it.")
		return nil
	}
	ci.Act = a
	return m.list.SetItem(idx, ci)
}

func (m model) View() string {
	content := m.list.View()
	if m.modalOpen {
		// Build modal content
		modal := m.renderActionModal()
		// Compose base content and modal using lipgloss compositor
		base := lipgloss.NewLayer(content)
		// Center the modal
		mw := lipgloss.Width(modal)
		mh := lipgloss.Height(modal)
		x := max(0, (m.width-mw)/2)
		// Center vertically
		y := max(0, (m.height-mh)/2)
		canvas := lipgloss.NewCanvas(
			base,
			lipgloss.NewLayer(modal).X(x).Y(y).Z(10),
		)
		rendered := canvas.Render()
		if m.status == "" {
			return rendered
		}
		status := lipgloss.NewStyle().Foreground(theme.Subtext0).Background(theme.Surface0).Render(m.status)
		return rendered + "\n" + status
	}
	if m.status == "" {
		return content
	}
	status := lipgloss.NewStyle().Foreground(theme.Subtext0).Background(theme.Surface0).Render(m.status)
	return content + "\n" + status
}

func (m model) collectActions() []commands.CommitAction {
	items := m.list.Items()
	cs := make([]commands.CommitAction, 0, len(items))
	for _, it := range items {
		ci := it.(commitItem)
		cs = append(cs, commands.CommitAction{Commit: ci.Commit, Action: string(ci.Act)})
	}
	return cs
}

// Run starts the TUI program.
func Run() error {
	m, err := initialModel()
	if err != nil {
		// Should not happen now, but keep as safety.
	}
	p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithOutput(os.Stdout), tea.WithInput(os.Stdin))
	if final, err := p.Run(); err != nil {
		return err
	} else if mm, ok := final.(model); ok && mm.doRebase {
		// After exiting the TUI, run the rebase so the user regains full terminal control.
		if err := commands.RunInteractiveRebase(mm.actions); err != nil {
			return err
		}
	}
	return nil
}

// max returns the maximum of two ints.
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

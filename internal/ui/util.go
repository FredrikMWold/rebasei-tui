package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss/v2"
)

// truncateToWidth returns a string cut to maxWidth characters (approx display width) with an ellipsis when truncated.
func truncateToWidth(s string, maxWidth int) string {
	if maxWidth <= 0 {
		return ""
	}
	if lipgloss.Width(s) <= maxWidth {
		return s
	}
	// Reserve 1 for ellipsis if possible
	target := maxWidth
	ell := ""
	if maxWidth > 1 {
		target = maxWidth - 1
		ell = "…"
	}
	// naive rune-based truncation (sufficient for ASCII used here)
	w := 0
	b := make([]rune, 0, target)
	for _, r := range s {
		rw := 1
		if w+rw > target {
			break
		}
		b = append(b, r)
		w += rw
	}
	return string(b) + ell
}

// countLines returns the number of visual lines by counting newlines in s.
// Trailing newline is ignored to avoid off-by-one.
func countLines(s string) int {
	if s == "" {
		return 0
	}
	s = strings.TrimSuffix(s, "\n")
	if s == "" {
		return 0
	}
	return 1 + strings.Count(s, "\n")
}

package ui

// action represents a rebase action for a commit.
type action string

const (
	pick   action = "pick"
	squash action = "squash"
	fixup  action = "fixup"
	edit   action = "edit"
	drop   action = "drop"
)

// actionDescription returns a short explanation for each rebase action.
func actionDescription(a action) string {
	switch a {
	case pick:
		return "use commit as-is"
	case squash:
		return "combine into previous commit; edit combined message"
	case fixup:
		return "combine into previous commit; keep previous message"
	case edit:
		return "pause to edit this commit during rebase"
	case drop:
		return "remove this commit"
	default:
		return ""
	}
}

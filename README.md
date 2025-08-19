# rebasei-tui

[![Go Reference](https://pkg.go.dev/badge/github.com/fredrikmwold/rebasei-tui.svg)](https://pkg.go.dev/github.com/fredrikmwold/rebasei-tui)
[![Release](https://img.shields.io/github/v/release/FredrikMWold/rebasei-tui?sort=semver)](https://github.com/FredrikMWold/rebasei-tui/releases)

**A minimal, keyboard-first TUI for interactive Git rebase** built with [Bubble Tea](https://github.com/charmbracelet/bubbletea). Reorder commits, choose actions (pick, squash, fixup, edit, drop), and kick off an interactive rebase — all from your terminal.

![Demo](./demo.gif)

<details>
	<summary><strong>Quick keys</strong></summary>

| Context | Key | Action |
|---|---|---|
| List | `↑`/`↓` | Move selection |
| List | `Ctrl+↑`/`Ctrl+↓` | Move commit up/down |
| List | `Enter` | Choose/set action (opens modal) |
| List | `p` | Mark as pick |
| List | `s` | Mark as squash |
| List | `f` | Mark as fixup |
| List | `e` | Mark as edit |
| List | `x`/`d` | Mark as drop |
| Anywhere | `Ctrl+r` | Start rebase |
| Modal | `Enter` | Confirm selected action |
| Modal | `Esc`/`q` | Cancel and close modal |
| Anywhere | `q`/`Ctrl+C` | Quit |

> Tip: The help footer updates based on what you can do at the moment.

</details>

## Features

- ⛓️ Reorder recent commits with keyboard controls
- ✍️ One-key actions: pick, squash, fixup, edit, drop

## Install

Install with Go:

```sh
go install github.com/fredrikmwold/rebasei-tui/cmd/rebasei-tui@latest
```

Or download a prebuilt binary from the Releases page and place it on your PATH:

- https://github.com/FredrikMWold/rebasei-tui/releases
Output demo.gif

# Terminal appearance and environment
Set FontSize 24
Set FontFamily "MesloLGS NF"
Set Theme "Catppuccin Mocha"
Set Width 1800
Set Height 900
Set Padding 10
Set Margin 30
Set MarginFill "#674EFF"
Set BorderRadius 12
Set Shell zsh

# Build the app and prepare a temporary Git repo with sample commits
Hide
Type "go build -o rebasei-tui ./cmd/rebasei-tui"
Enter
Type `PROJ="$PWD"`
Enter
Type `TMP="$(mktemp -d)"`
Enter
Type `cd "$TMP"`
Enter
Type "git init -q && git config user.name Demo && git config user.email demo@example.com"
Enter
Type "printf 'one\n' > file.txt && git add file.txt && git commit -q -m 'feat: initial commit'"
Enter
Type "printf 'two\n' >> file.txt && git commit -qam 'feat: add two'"
Enter
Type "printf 'three\n' >> file.txt && git commit -qam 'feat: add three'"
Enter
Type "printf 'four\n' >> file.txt && git commit -qam 'feat: add four'"
Enter
Type "printf 'five\n' >> file.txt && git commit -qam 'fix: adjust five'"
Enter
Type `alias rebasei-tui="$PROJ/rebasei-tui" && clear`
Enter
Show

# Launch the TUI (inside the demo repo)
Type "rebasei-tui"
Enter
Sleep 2.2s

# Navigate the commit list
Down
Sleep 600ms
Down
Sleep 600ms
Up
Sleep 1s

# Open action modal and choose 'squash'
Enter
Sleep 800ms
Down
Sleep 500ms
Enter
Sleep 800ms

# Mark next commit as fixup, and another as drop
Down
Sleep 500ms
Type "f"
Sleep 900ms
Down
Sleep 500ms
Type "x"
Sleep 3s
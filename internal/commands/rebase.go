package commands

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

type CommitAction struct {
	Commit Commit
	Action string // pick, squash, fixup, edit, drop
}

func RunInteractiveRebase(list []CommitAction) error {
	if len(list) == 0 {
		return fmt.Errorf("no commits to rebase")
	}
	// Build todo in chronological order (oldest first) so squash/fixup have a previous commit.
	todo := ""
	for i := len(list) - 1; i >= 0; i-- {
		ca := list[i]
		if ca.Action == "" {
			ca.Action = "pick"
		}
		todo += fmt.Sprintf("%s %s %s\n", ca.Action, ca.Commit.Hash, ca.Commit.Subject)
	}

	tmpDir, err := os.MkdirTemp("", "rebasei-tui-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)
	todoPath := filepath.Join(tmpDir, "todo.txt")
	if err := os.WriteFile(todoPath, []byte(todo), 0o644); err != nil {
		return err
	}

	scriptPath := filepath.Join(tmpDir, "seq_editor.sh")
	script := fmt.Sprintf("#!/bin/sh\ncat '%s' > \"$1\"\n", todoPath)
	if err := os.WriteFile(scriptPath, []byte(script), 0o755); err != nil {
		return err
	}

	n := len(list)
	// Use --root when there aren't enough ancestors for HEAD~n
	countCmd := exec.Command("git", "rev-list", "--count", "HEAD")
	countCmd.Env = append(os.Environ(), "GIT_PAGER=cat")
	out, _ := countCmd.Output()
	total := 0
	if len(out) > 0 {
		if v, err := strconv.Atoi(strings.TrimSpace(string(out))); err == nil {
			total = v
		}
	}
	var cmd *exec.Cmd
	if total > 0 && n >= total {
		cmd = exec.Command("git", "-c", "sequence.editor="+scriptPath, "rebase", "-i", "--root")
	} else {
		base := fmt.Sprintf("HEAD~%d", n)
		cmd = exec.Command("git", "-c", "sequence.editor="+scriptPath, "rebase", "-i", base)
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

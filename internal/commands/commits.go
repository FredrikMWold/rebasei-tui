package commands

import (
	"bufio"
	"bytes"
	"errors"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type Commit struct {
	Hash      string
	HashShort string
	Subject   string
	Author    string
	Date      string // YYYY-MM-DD
	Tags      []string
}

func ListCommits(n int) ([]Commit, error) {
	// %D includes ref names like "HEAD -> main, tag: v1.0.0, origin/main"
	args := []string{"log", "--date=short", "--pretty=format:%h\t%H\t%s\t%an\t%ad\t%D", "-n", strconv.Itoa(n)}
	cmd := exec.Command("git", args...)
	cmd.Env = append(os.Environ(), "GIT_PAGER=cat")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	res := []Commit{}
	s := bufio.NewScanner(bytes.NewReader(out))
	for s.Scan() {
		line := s.Text()
		parts := strings.SplitN(line, "\t", 6)
		if len(parts) < 5 { // tolerate missing %D on some lines
			continue
		}
		var tags []string
		if len(parts) >= 6 {
			// Parse %D for tags: look for segments starting with "tag: "
			for _, seg := range strings.Split(parts[5], ",") {
				seg = strings.TrimSpace(seg)
				const p = "tag: "
				if strings.HasPrefix(seg, p) {
					t := strings.TrimPrefix(seg, p)
					if t != "" {
						tags = append(tags, t)
					}
				}
			}
		}
		res = append(res, Commit{
			HashShort: parts[0],
			Hash:      parts[1],
			Subject:   parts[2],
			Author:    parts[3],
			Date:      parts[4],
			Tags:      tags,
		})
	}
	if err := s.Err(); err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, errors.New("no commits found; are you in a git repo?")
	}
	return res, nil
}

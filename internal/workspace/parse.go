// Package workspace derives the names that identify a feature workspace
// (branch, worktree path, tmux session) from the user-supplied argument.
package workspace

import (
	"fmt"
	"path/filepath"
	"strings"
)

type Names struct {
	Ticket      string
	Description string
	Workspace   string
	Branch      string
	Path        string
	Session     string
}

func Parse(arg, prefix, base string) (Names, error) {
	arg = strings.TrimSpace(arg)
	if arg == "" {
		return Names{}, fmt.Errorf("empty argument")
	}
	if prefix == "" {
		return Names{}, fmt.Errorf("empty branch prefix")
	}
	if base == "" {
		return Names{}, fmt.Errorf("empty worktree base")
	}

	left, right, ok := strings.Cut(arg, "/")
	if !ok {
		return Names{}, fmt.Errorf("argument %q must be of the form TICKET/description", arg)
	}
	ticket := strings.TrimSpace(left)
	desc := strings.TrimSpace(right)
	if ticket == "" {
		return Names{}, fmt.Errorf("ticket part is empty in %q", arg)
	}
	if desc == "" {
		return Names{}, fmt.Errorf("description part is empty in %q", arg)
	}

	workspace := ticket + "-" + desc
	return Names{
		Ticket:      ticket,
		Description: desc,
		Workspace:   workspace,
		Branch:      prefix + "/" + ticket + "/" + desc,
		Path:        filepath.Join(base, workspace),
		Session:     workspace,
	}, nil
}

package cli

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/tscolari/work/internal/config"
	"github.com/tscolari/work/internal/gitx"
	"github.com/tscolari/work/internal/tmuxx"
)

// RunWorkend executes the `workend` command.
func RunWorkend(args []string, stdout, stderr io.Writer) error {
	fs := flag.NewFlagSet("workend", flag.ContinueOnError)
	fs.SetOutput(stderr)
	force := fs.Bool("force", false, "skip the unmerged-commits check and delete the branch anyway")
	dryRun := fs.Bool("dry-run", false, "print what would happen without doing anything")
	if err := fs.Parse(args); err != nil {
		return userErr("%v", err)
	}
	if fs.NArg() != 0 {
		return userErr("workend takes no positional arguments")
	}

	cfg, err := config.Load()
	if err != nil {
		return sysErr("load config: %v", err)
	}

	cwd, err := os.Getwd()
	if err != nil {
		return sysErr("getwd: %v", err)
	}
	cwdResolved, err := filepath.EvalSymlinks(cwd)
	if err != nil {
		cwdResolved = cwd
	}
	baseResolved, err := filepath.EvalSymlinks(cfg.WorktreeBase)
	if err != nil {
		baseResolved = cfg.WorktreeBase
	}

	if filepath.Dir(cwdResolved) != baseResolved {
		return userErr("workend must be run from inside a workspace under %s (cwd is %s)", cfg.WorktreeBase, cwd)
	}
	name := filepath.Base(cwdResolved)

	branch, err := gitx.BranchOfWorktree(cwd)
	if err != nil {
		return sysErr("read worktree HEAD: %v", err)
	}

	repoDir, err := gitx.CommonDir(cwd)
	if err != nil {
		return sysErr("locate main repo: %v", err)
	}

	other, err := gitx.BranchCheckedOutElsewhere(repoDir, branch, cwd)
	if err != nil {
		return sysErr("check other worktrees: %v", err)
	}
	if other != "" {
		return userErr("branch %s is checked out at %s; remove that worktree first", branch, other)
	}

	if !*force {
		unmerged, err := gitx.HasUnmergedCommits(repoDir, branch)
		if err != nil {
			return sysErr("check unmerged commits: %v", err)
		}
		if unmerged {
			return userErr("branch %s has unmerged commits vs upstream; pass --force to delete anyway", branch)
		}
	}

	if *dryRun {
		fmt.Fprintf(stdout, "would kill tmux session: %s\n", name)
		fmt.Fprintf(stdout, "would remove worktree:   %s\n", cwd)
		fmt.Fprintf(stdout, "would delete branch:     %s\n", branch)
		return nil
	}

	if tmuxx.Installed() {
		has, err := tmuxx.HasSession(name)
		if err != nil {
			return sysErr("check tmux session: %v", err)
		}
		if has {
			if err := tmuxx.KillSession(name); err != nil {
				fmt.Fprintf(stderr, "warning: kill tmux session %s: %v\n", name, err)
			}
		}
	}

	if err := gitx.WorktreeRemove(repoDir, cwd); err != nil {
		return sysErr("remove worktree: %v", err)
	}
	if err := gitx.BranchDelete(repoDir, branch); err != nil {
		return sysErr("delete branch: %v", err)
	}

	return nil
}

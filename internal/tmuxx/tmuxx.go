// Package tmuxx wraps the subset of tmux commands that work/workend need.
package tmuxx

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

func Installed() bool {
	_, err := exec.LookPath("tmux")
	return err == nil
}

func InsideTmux() bool {
	return os.Getenv("TMUX") != ""
}

func HasSession(name string) (bool, error) {
	cmd := exec.Command("tmux", "has-session", "-t="+name)
	if err := cmd.Run(); err != nil {
		var ee *exec.ExitError
		if errors.As(err, &ee) {
			return false, nil
		}
		return false, fmt.Errorf("tmux has-session: %w", err)
	}
	return true, nil
}

// NewSessionAttached replaces the current process with `tmux new-session`.
// On success, this function does not return.
func NewSessionAttached(name, cwd string) error {
	bin, err := exec.LookPath("tmux")
	if err != nil {
		return fmt.Errorf("locate tmux: %w", err)
	}
	argv := []string{"tmux", "new-session", "-s", name, "-c", cwd}
	if err := syscall.Exec(bin, argv, os.Environ()); err != nil {
		return fmt.Errorf("exec tmux: %w", err)
	}
	return nil
}

// AttachSession replaces the current process with `tmux attach-session`.
// Use when a session already exists and we are outside tmux.
// On success, this function does not return.
func AttachSession(name string) error {
	bin, err := exec.LookPath("tmux")
	if err != nil {
		return fmt.Errorf("locate tmux: %w", err)
	}
	argv := []string{"tmux", "attach-session", "-t", name}
	if err := syscall.Exec(bin, argv, os.Environ()); err != nil {
		return fmt.Errorf("exec tmux: %w", err)
	}
	return nil
}

func NewSessionDetached(name, cwd string) error {
	return run("tmux", "new-session", "-d", "-s", name, "-c", cwd)
}

func SwitchClient(name string) error {
	return run("tmux", "switch-client", "-t", name)
}

func KillSession(name string) error {
	return run("tmux", "kill-session", "-t", name)
}

func run(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s %s: %w: %s", name, strings.Join(args, " "), err, strings.TrimSpace(string(out)))
	}
	return nil
}

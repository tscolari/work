# `work` / `workend` CLI Tool Spec

## Overview

Two companion commands that manage the lifecycle of a feature workspace: git worktree + branch + tmux session, created and torn down as a unit.

Written in Go. Single binary, no external dependencies beyond git and tmux.

## Configuration

A config file at `~/.config/work/config` (or a path set via `WORK_CONFIG` env var), format:

```
worktree_base=~/worktrees
branch_prefix=tscolari
```

Both values have sensible defaults: `worktree_base` defaults to `~/worktrees`, `branch_prefix` defaults to the output of `git config user.name` (lowercased, spaces replaced with dashes), falling back to `whoami`.

---

## `work`

### Usage

```
work <TICKET/description-of-feature>
```

The argument is a single token with the form `TICKET-NUM/kebab-description`. The slash is the delimiter between ticket and description.

### Behavior

Given `work SOM-123/add-retry-logic` and the config above:

1. **Parse input** — split on the first `/` to get:
   - ticket: `SOM-123`
   - description: `add-retry-logic`
   - workspace name: `SOM-123-add-retry-logic` (slash replaced with dash)

2. **Create branch** — from the current HEAD:
   - branch name: `tscolari/SOM-123/add-retry-logic`
   - fail if the branch already exists (with a helpful message suggesting the worktree path to `cd` into)

3. **Create worktree** at `~/worktrees/SOM-123-add-retry-logic` using that branch.

4. **Create and attach a tmux session** named `SOM-123-add-retry-logic`, with the working directory set to the worktree path.
   - If not inside tmux (`TMUX` env var is unset): create and attach with `tmux new-session -s <name> -c <path>`. Use `syscall.Exec` to replace the current process.
   - If already inside tmux: create the session detached (`tmux new-session -d -s <name> -c <path>`) then switch to it (`tmux switch-client -t <name>`).

### Error cases

- Missing argument → print usage and exit 1.
- Not inside a git repo → error.
- Branch already exists → error with message: `branch already exists; cd into ~/worktrees/SOM-123-add-retry-logic or use workend to clean up first`.
- Worktree path already exists → same treatment.
- tmux not installed → skip session creation, print warning, still create worktree + branch.

---

## `workend`

### Usage

```
workend
```

No arguments. Operates based on the current working directory.

### Behavior

1. **Detect workspace name** — derive it from the current directory's basename. If the current directory is not under the configured `worktree_base`, error out.

2. **Derive the branch name** — read it from the worktree's `.git` metadata (parse the worktree's HEAD to get the checked-out branch), so it stays correct even if naming conventions drift.

3. **Kill tmux session** named after the workspace (if it exists). If running inside that tmux session, tmux will automatically move attached clients to the next available session, or detach them if none exists.

4. **Remove worktree** — `git worktree remove --force <path>`.

5. **Delete local branch** — `git branch -D <branch>`. Refuse if the branch is checked out elsewhere (another worktree). Print a warning if the branch has unmerged commits compared to its upstream, requiring a `--force` flag to proceed.

### Flags

- `--force` — skip the unmerged-commits check and delete anyway.
- `--dry-run` — print what would happen without doing anything.

### Error cases

- Not inside a worktree directory under `worktree_base` → error with message explaining this must be run from inside a workspace.
- Branch checked out in another worktree → error, suggest removing that worktree first.

---

## Binary layout

Either:
- A single binary with subcommands (`work start`, `work end`), or
- Two separate binaries (`work`, `workend`)

Implementor's choice. Both are fine.

## State management

No daemon, no background processes, no state files beyond the config. All state lives in git and tmux, which are the source of truth.

## Exit codes

- `0` — success
- `1` — user error (bad input, missing prereqs)
- `2` — unexpected failure

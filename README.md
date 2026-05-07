# WORK / WORKEND

These are 2 vibecoded helpers that I created for my day to day usage.

## The flow

They assist me with how I like to work in a codebase:

1. From inside the repository I want to work on
2. Create a new branch, prefixed with my handler, the ticket code and followed by a small description
3. Create a worktree for that branch and jump into it
4. Create a new tmux session, named after the ticket code and description
5. Jump into the session and start working

With a quick tear down
1. Delete the branch / worktree and tmux session

---

## With the tool

1. `cd ./codebase`
2. `work TKT-1234/creating-api-mocks`
3. ... work
4. `workend` (from the worktree dir)

--

## Configuration

~/.config/work/config (or set `WORK_CONFIG`)
```
worktree_base=~/inflight
branch_prefix=myname
```

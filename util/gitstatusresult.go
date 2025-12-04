package util

import (
	"fmt"
	"github.com/go-git/go-git/v6"
	"slices"
)

type GitStatusResult struct {
	Staging  StatusCode
	Worktree StatusCode
	Status   git.Status
	Path     RepoPath
}

func (s GitStatusResult) print(preifx string) {
	fmt.Printf("%-10s status to string: %v", preifx, s.Status.String())
	fmt.Printf("%-10s %-50s s:%v w:%v\n", preifx, s.Path, s.Staging, s.Worktree)
}

func (s GitStatusResult) NeedGitCommitFiles(states []git.StatusCode) (ret []RepoPath) {
	for k, v := range s.Status {
		status := v.Staging
		if slices.Contains(states, status) {
			ret = append(ret, RepoPath(k))
		}
	}
	return
}

func (s GitStatusResult) NeedGitCommit() string {
	action := ""
	for _, v := range s.Status {
		status := v.Staging
		switch status {
		case git.Added:
			action = "ADD"
		case git.Deleted:
			action = "RM"
		case git.Modified:
			action = "UPDATE"
		}
		if action != "" {
			return action
		}
	}
	return ""
}

func (s GitStatusResult) NeedGitRMFiles(work bool) (ret []RepoPath) {
	for k, v := range s.Status {
		status := v.Worktree
		if !work {
			status = v.Staging
		}
		if status == git.Deleted {
			ret = append(ret, RepoPath(k))
		}
	}
	return ret
}

func (s GitStatusResult) NeedGitAddFiles() (ret []string) {
	for k, v := range s.Status {
		status := v.Worktree
		if status == git.Modified || status == git.Untracked {
			ret = append(ret, k)
		}
	}
	return ret
}

func (s GitStatusResult) NeedGitAdd() bool {
	for _, v := range s.Status {
		status := v.Worktree
		if status == git.Modified || status == git.Untracked {
			return true
		}
	}
	return false
}

func (s GitStatusResult) NeedGitRm() bool {
	for _, v := range s.Status {
		status := v.Worktree
		if status == git.Deleted {
			return true
		}
	}
	return false
}

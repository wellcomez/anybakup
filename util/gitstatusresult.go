package util

import (
	"fmt"
	"slices"
	"strings"

	"github.com/go-git/go-git/v6"
)

type FileStatus2 struct {
	// Staging is the status of a file in the staging area
	Staging git.StatusCode
	// Worktree is the status of a file in the worktree
	Worktree git.StatusCode
	// Extra contains extra information, such as the previous name in a rename
	Extra string
}

func (s GitStatus2) String() string {
	sp := "\n"
	sp += fmt.Sprintf("%-100s %-10s %-10s %-10s\n", "Path", "Staging", "Worktree", "Extra")
	for k, v := range s {
		sp += fmt.Sprintf("%-100s %-10c %-10c %-10s\n", k, v.Staging, v.Worktree, v.Extra)
	}
	return sp
}

type GitStatus2 map[string]*FileStatus2
type GitStatusResult struct {
	Staging  StatusCode
	Worktree StatusCode
	Status   git.Status
	// StatusOrgin git.Status
	Path RepoPath
}

func (s GitStatusResult) print(preifx string) {
	// fmt.Printf("%-10s status to string: %v", preifx, s.StatusOrgin.String())
	fmt.Printf("%-10s %v", preifx, s.Status.String())
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
			if strings.HasPrefix(k, s.Path.Sting()) {
				ret = append(ret, RepoPath(k))
			}
		}
	}
	return ret
}

func (s GitStatusResult) NeedGitAddFiles() (ret []RepoPath) {
	for k, v := range s.Status {
		status := v.Worktree
		if status == git.Modified || status == git.Untracked {
			if strings.HasPrefix(k, s.Path.Sting()) {
				ret = append(ret, RepoPath(k))
			}
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

func CheckStatus(git *git.Repository) (git.Status, error) {
	w, err := git.Worktree()
	if err != nil {
		return nil, fmt.Errorf("status file worktree %v", err)
	}

	// w.StatusWithOptions(git.StatusOptions{Strategy: git.Preload})
	status, err := w.Status()
	fmt.Printf("status to string: %v", status.String())
	fmt.Println("")
	if err != nil {
		return nil, fmt.Errorf("workstate err=%v", err)
	}
	for a, k := range status {
		fmt.Printf("%-50s Staging=%c worktree=%c Extra=%s\n", a, k.Staging, k.Worktree, k.Extra)
	}
	return status, nil
}

func (repo *GitRepo) GetStateStage(file string) (StatusCode, error) {
	if r, err := repo.Status(RepoPath(repo.AbsRepo2Repo(file).Sting())); err != nil {
		return GitStatusErro, err
	} else {
		return r.Staging, nil
	}
}

func (repo *GitRepo) GetStateWorkTree(file string) (StatusCode, error) {
	// if repo == nil {
	// 	r, _ := NewGitReop()
	// 	repo = r
	// }
	if r, err := repo.Status(repo.AbsRepo2Repo(file)); err != nil {
		return GitStatusErro, err
	} else {
		return r.Worktree, nil
	}
}

func GetStatuscode(newVar git.StatusCode) StatusCode {
	switch newVar {
	case git.Unmodified:
		return GitUnmodified
	case git.Untracked:
		return GitUntracked
	case git.Modified:
		return GitModified
	case git.Added:
		return GitAdded
	case git.Deleted:
		return GitDeleted
	case git.Renamed:
		return GitRenamed
	case git.Copied:
		return GitCopied
	case git.UpdatedButUnmerged:
		return GitUpdatedButUnmerged
	default:
		return GitStatusErro
	}
}

type StatusCode string

const (
	GitUnmodified         StatusCode = "N"
	GitUntracked          StatusCode = "?"
	GitModified           StatusCode = "M"
	GitAdded              StatusCode = "A"
	GitDeleted            StatusCode = "D"
	GitRenamed            StatusCode = "R"
	GitCopied             StatusCode = "C"
	GitUpdatedButUnmerged StatusCode = "U"
	GitStatusErro         StatusCode = "E"
)

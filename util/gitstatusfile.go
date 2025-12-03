package util

import (
	"fmt"

	"github.com/go-git/go-git/v6"
	// fixtures "github.com/go-git/go-git-fixtures/v5"
)

// type GitStatusFile struct {
// 	// status   git.Status
// 	reporoot *GitRepo
// }

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

func GetStateStage(file string, repo *GitRepo) (StatusCode, error) {
	if repo == nil {
		r, _ := NewGitReop()
		repo = r
	}
	if r, err := repo.Status(RepoPath(repo.AbsRepo2Repo(file).Sting())); err != nil {
		return GitStatusErro, err
	} else {
		return r.Staging, nil
	}
}

func GetStateWorkTree(file string, repo *GitRepo) (StatusCode, error) {
	if repo == nil {
		r, _ := NewGitReop()
		repo = r
	}
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

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
	status, err := CheckStatus(repo.repo)
	if err != nil {
		return GitStatusErro, err
	}
	gitfile, err := repo.rel(file)
	if err != nil {
		return GitStatusErro, err
	}
	st := status.File(gitfile)
	newVar := st.Staging
	switch newVar {
	case git.Unmodified:
		return GitUnmodified, nil
	case git.Untracked:
		return GitUntracked, nil
	case git.Modified:
		return GitModified, nil
	case git.Added:
		return GitAdded, nil
	case git.Deleted:
		return GitDeleted, nil
	case git.Renamed:
		return GitRenamed, nil
	case git.Copied:
		return GitCopied, nil
	case git.UpdatedButUnmerged:
		return GitUpdatedButUnmerged, nil
	}
	return GitStatusErro, fmt.Errorf("workstate %v %v", st.Worktree, file)
}
func GetStateWorkTree(file string, repo *GitRepo) (StatusCode, error) {
	if repo == nil {
		r, _ := NewGitReop()
		repo = r
	}
	status, err := CheckStatus(repo.repo)
	if err != nil {
		return GitStatusErro, err
	}
	gitfile, err := repo.rel(file)
	if err != nil {
		return GitStatusErro, err
	}
	st := status.File(gitfile)
	newVar := st.Worktree
	switch newVar {
	case git.Unmodified:
		return GitUnmodified, nil
	case git.Untracked:
		return GitUntracked, nil
	case git.Modified:
		return GitModified, nil
	case git.Added:
		return GitAdded, nil
	case git.Deleted:
		return GitDeleted, nil
	case git.Renamed:
		return GitRenamed, nil
	case git.Copied:
		return GitCopied, nil
	case git.UpdatedButUnmerged:
		return GitUpdatedButUnmerged, nil
	}
	return GitStatusErro, fmt.Errorf("workstate %v %v", st.Worktree, file)
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

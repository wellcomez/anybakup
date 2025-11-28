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

//	func NewGitCheckStatus() (*GitStatusFile, error) {
//		s := GitStatusFile{}
//		if s.reporoot == nil {
//			var err error
//			s.reporoot, err = NewGitReop()
//			if err != nil {
//				return nil, fmt.Errorf("status file %v", err)
//			}
//		}
//		return &s, nil
//	}
func CheckStatus(git *git.Repository) (git.Status, error) {
	w, err := git.Worktree()
	if err != nil {
		return nil, fmt.Errorf("status file worktree %v", err)
	}

	// w.StatusWithOptions(git.StatusOptions{Strategy: git.Preload})
	status, err := w.Status()
	fmt.Printf("status file Status %v\n", status.String())
	if err != nil {
		return nil, fmt.Errorf("status file Status %v", err)
	}
	for a, k := range status {
		fmt.Printf("%-50s w:=[%c]|[%s]\n", a, k.Worktree, k.Extra)
	}
	return status, nil
}
func GetState(file string, repo *GitRepo) (StatusCode, error) {
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
	if _, ok := status[gitfile]; !ok {
		return GitUnmodified, nil
	}
	if st.Worktree == git.Unmodified {
		return GitUnmodified, nil
	}
	if st.Worktree == git.Untracked {
		return GitUntracked, nil
	}
	if st.Worktree == git.Modified {
		return GitModified, nil
	}
	if st.Worktree == git.Added {
		return GitAdded, nil
	}
	if st.Worktree == git.Deleted {
		return GitDeleted, nil
	}
	if st.Worktree == git.Renamed {
		return GitRenamed, nil
	}
	if st.Worktree == git.Copied {
		return GitCopied, nil
	}
	if st.Worktree == git.UpdatedButUnmerged {
		return GitUpdatedButUnmerged, nil
	}
	return GitStatusErro, fmt.Errorf("status file %v %v", err, file)
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

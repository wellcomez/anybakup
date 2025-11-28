package util

import (
	"fmt"

	"github.com/go-git/go-git/v5"
)

type GitStatusFile struct {
	status   git.Status
	reporoot *GitRepo
}

func NewGitCheckStatus() (*GitStatusFile, error) {
	s := GitStatusFile{}
	if s.reporoot == nil {
		var err error
		s.reporoot, err = NewGitReop()
		if err != nil {
			return nil, fmt.Errorf("status file %v", err)
		}
	}
	return &s, nil
}
func (s *GitStatusFile) CheckStatus() error {
	repo, err := s.reporoot.Open()
	if err != nil {
		return fmt.Errorf("status file Open %v", err)
	}
	w, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("status file worktree %v", err)
	}
	status, err := w.Status()
	if err != nil {
		return fmt.Errorf("status file Status %v", err)
	}
	s.status = status
	return nil
}
func (s GitStatusFile) GetState(file string) (StatusCode, error) {
	gitfile, err := s.reporoot.rel(file)
	if err != nil {
		return GitStatusErro, err
	}
	st := s.status.File(gitfile)
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

type StatusCode uint8

const (
	GitUnmodified         StatusCode = ' '
	GitUntracked          StatusCode = '?'
	GitModified           StatusCode = 'M'
	GitAdded              StatusCode = 'A'
	GitDeleted            StatusCode = 'D'
	GitRenamed            StatusCode = 'R'
	GitCopied             StatusCode = 'C'
	GitUpdatedButUnmerged StatusCode = 'U'
	GitStatusErro         StatusCode = 'E'
)

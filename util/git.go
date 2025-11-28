package util

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

type GitRepo struct {
	root string
}

func (r GitRepo) rel(s string) (string, error) {
	if !filepath.IsAbs(s) {
		f := filepath.Join(r.root, s)
		if _, err := os.Stat(f); err != nil {
			return "", fmt.Errorf("git repo %v %v", err, s)
		}
		return s, nil
	}

	rel, err := filepath.Rel(r.root, s)
	if err != nil {
		return "", fmt.Errorf("git repo %v %v", err, s)
	}
	return rel, nil
}
func (r *GitRepo) load() error {
	conf := Config{}
	if err := conf.Load(); err != nil {
		return fmt.Errorf("git repo %v", err)
	}
	st, err := os.Stat(conf.RepoDir.String())
	if err != nil {
		return fmt.Errorf("git repo %v %v", err, conf.RepoDir)
	}
	if !st.IsDir() {
		return fmt.Errorf("git repo %v is not a directory", conf.RepoDir)
	}
	r.root = conf.RepoDir.String()
	return nil
}
func NewGitReop() (*GitRepo, error) {
	ret := &GitRepo{}
	if err := ret.load(); err != nil {
		return nil, err
	}
	return ret, nil
}

func (r GitRepo) Open() (*git.Repository, error) {
	git, err := git.PlainOpen(r.root)
	if err != nil {
		return nil, fmt.Errorf("open git repo %v %v", err, r.root)
	}
	return git, nil
}

func Initgit() error {
	repo, err := NewGitReop()
	if err != nil {
		return fmt.Errorf("init git repo %v", err)
	}
	_, err = git.PlainInit(repo.root, false)
	if err != nil {
		return fmt.Errorf("init git repo %v %v", err, repo.root)
	}
	return nil
}
func (r GitRepo) GitAddFile(file string) error {
	repo, err := r.Open()
	if err != nil {
		return fmt.Errorf("git add %v", err)
	}
	w, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("git add %v %v", err, file)
	}
	gitfile, err := r.rel(file)
	if err != nil {
		return fmt.Errorf("git add %v %v", err, file)
	}
	_, err = w.Add(gitfile)
	if err != nil {
		return fmt.Errorf("git add %v %v", err, file)
	}
	msg := fmt.Sprintf("ADD %v", file)
	_, err = w.Commit(msg, &git.CommitOptions{
		Author: &object.Signature{
			Name: "anybakup",
			When: time.Now(),
		},
	})
	if err != nil {
		return fmt.Errorf("git add %v %v", err, file)
	}
	return nil
}

type GitStatusFile struct {
	status   git.Status
	cleand   bool
	reporoot *GitRepo
}

func (s *GitStatusFile) CheckStatus() error {
	if s.reporoot == nil {
		var err error
		s.reporoot, err = NewGitReop()
		if err != nil {
			return fmt.Errorf("status file %v", err)
		}
	}
	dir := s.reporoot.root
	repo, err := git.PlainOpen(dir)
	if err != nil {
		return fmt.Errorf("status file %v %v", err, dir)
	}
	w, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("status file %v %v", err, dir)
	}
	status, err := w.Status()
	if err != nil {
		return fmt.Errorf("status file %v %v", err, dir)
	}
	s.status = status
	s.cleand = status.IsClean()
	return nil
}
func GitCommitFile(file string) error {
	reporoot, err := NewGitReop()
	if err != nil {
		return fmt.Errorf("commit file %v", err)
	}
	dir := reporoot.root
	repo, err := git.PlainOpen(dir)
	if err != nil {
		return fmt.Errorf("commit file %v %v", err, dir)
	}
	w, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("commit file %v %v", err, dir)
	}
	_, err = w.Commit("commit file", &git.CommitOptions{
		Author: &object.Signature{
			Name: "anybakup",
			When: time.Now(),
		},
	})
	if err != nil {
		return fmt.Errorf("commit file %v %v", err, file)
	}
	return nil
}

func GitDiffFile(file string) error {
	reporoot, err := NewGitReop()
	if err != nil {
		return fmt.Errorf("diff file %v", err)
	}
	dir := reporoot.root
	repo, err := git.PlainOpen(dir)
	if err != nil {
		return fmt.Errorf("diff file %v %v", err, dir)
	}
	w, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("diff file %v %v", err, dir)
	}
	status, err := w.Status()
	if err != nil {
		return fmt.Errorf("diff file %v %v", err, file)
	}
	if status.IsClean() {
		return nil
	}
	return nil
}

package util

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

type reporoot struct {
	root string
}

func (r reporoot) rel(s string) (string, error) {
	if !filepath.IsAbs(s) {
		return "", fmt.Errorf("git repo %v is not absolute path", s)
	}
	rel, err := filepath.Rel(r.root, s)
	if err != nil {
		return "", fmt.Errorf("git repo %v %v", err, s)
	}
	return rel, nil
}
func (r *reporoot) load() error {
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
func new_repo() (*reporoot, error) {
	ret := &reporoot{}
	if err := ret.load(); err != nil {
		return nil, err
	}
	return ret, nil
}
func Initgit() error {
	repo, err := new_repo()
	if err != nil {
		return fmt.Errorf("init git repo %v", err)
	}
	_, err = git.PlainInit(repo.root, false)
	if err != nil {
		return fmt.Errorf("init git repo %v %v", err, repo.root)
	}
	return nil
}
func GitAddFile(file string) error {
	reporoot, err := new_repo()
	if err != nil {
		return fmt.Errorf("git add %v", err)
	}
	dir := reporoot.root
	repo, err := git.PlainOpen(dir)
	if err != nil {
		return fmt.Errorf("git add %v %v", err, dir)
	}
	w, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("git add %v %v", err, dir)
	}
	gitfile, err := reporoot.rel(file)
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
type GitStatusFile struct  {
	status git.Status
	cleand bool
}

func(s* GitStatusFile)Check() error{
	reporoot, err := new_repo()
	if err != nil {
		return fmt.Errorf("status file %v", err)
	}
	dir := reporoot.root
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
	reporoot, err := new_repo()
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
	reporoot, err := new_repo()
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

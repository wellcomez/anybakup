package util

import (
	"fmt"
	"os"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func repo_root() (string, error) {
	conf := Config{}
	if err := conf.Load(); err != nil {
		return "", fmt.Errorf("git repo %v", err)
	}
	st, err := os.Stat(conf.RepoDir)
	if err != nil {
		return "", fmt.Errorf("git repo %v %v", err, conf.RepoDir)
	}
	if !st.IsDir() {
		return "", fmt.Errorf("git repo %v is not a directory", conf.RepoDir)
	}
	return conf.RepoDir, nil
}
func Initgit() error {
	dir, err := repo_root()
	if err != nil {
		return fmt.Errorf("init git repo %v", err)
	}
	_, err = git.PlainInit(dir, false)
	if err != nil {
		return fmt.Errorf("init git repo %v %v", err, dir)
	}
	return nil
}
func GitAddFile(file string) error {
	dir, err := repo_root()
	if err != nil {
		return fmt.Errorf("add file %v", err)
	}
	repo, err := git.PlainOpen(dir)
	if err != nil {
		return fmt.Errorf("add file %v %v", err, dir)
	}
	w, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("add file %v %v", err, dir)
	}
	_, err = w.Add(file)
	if err != nil {
		return fmt.Errorf("add file %v %v", err, file)
	}
	_, err = w.Commit("add file", &git.CommitOptions{
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
func GitStatusFile(file string) error {
	dir, err := repo_root()
	if err != nil {
		return fmt.Errorf("status file %v", err)
	}
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
		return fmt.Errorf("status file %v %v", err, file)
	}
	if status.IsClean() {
		return nil
	}
	return nil
}
func GitCommitFile(file string) error {
	dir, err := repo_root()
	if err != nil {
		return fmt.Errorf("commit file %v", err)
	}
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
	dir, err := repo_root()
	if err != nil {
		return fmt.Errorf("diff file %v", err)
	}
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

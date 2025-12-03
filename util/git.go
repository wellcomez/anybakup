package util

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing"
	"github.com/go-git/go-git/v6/plumbing/object"
)

type (
	RepoPath string
	SrcPath  string
	SysPath  string
	RepoRoot string
	GitRepo  struct {
		root string
		repo *git.Repository
	}
)

type GitStatusResult struct {
	Staging  StatusCode
	Worktree StatusCode
	Status   git.Status
	Path     RepoPath
}

func (s GitStatusResult) print(preifx string) {
	fmt.Printf("%-10s status to string: %v", preifx, s.Status.String())
	fmt.Printf("%-10s add %-50s s:%v w:%v\n", preifx, s.Path, s.Staging, s.Worktree)
}

func (s GitStatusResult) NeedGitCommitFiles(states []git.StatusCode) (ret []string) {
	for k, v := range s.Status {
		status := v.Staging
		for _, state := range states {
			if status == state {
				ret = append(ret, k)
				break
			}
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

func (r *GitRepo) Status(gitfile RepoPath) (GitStatusResult, error) {
	ret := GitStatusResult{Staging: GitStatusErro, Worktree: GitStatusErro, Path: gitfile}
	if git, err := r.Open(); err != nil {
		return ret, fmt.Errorf("status file Open %v", err)
	} else {
		if w, err := git.Worktree(); err != nil {
			return ret, fmt.Errorf("status file worktree %v", err)
		} else if status, err := w.Status(); err != nil {
			return ret, fmt.Errorf("status file status %v", err)
		} else {
			st := status.File(gitfile.Sting())
			ret.Staging = GetStatuscode(st.Staging)
			ret.Worktree = GetStatuscode(st.Worktree)
			ret.Status = status
			return ret, nil
		}
	}
}

func (s SrcPath) Repo(d RepoRoot) RepoPath {
	rel, err := filepath.Rel("/", string(s))
	if err != nil {
		return ""
	}
	return RepoPath(rel)
}

func (s SrcPath) Sting() string {
	return string(s)
}

func (s RepoPath) Sting() string {
	return string(s)
}

func (s RepoPath) ToAbs(repo GitRepo) string {
	r := RepoRoot(repo.root)
	return r.With(s.Sting())
}

func (r GitRepo) Rel(s string) (string, error) {
	if !filepath.IsAbs(s) {
		// f := filepath.Join(r.root, s)
		// if _, err := os.Stat(f); err != nil {
		// 	return "", fmt.Errorf("git repo %v %v", err, s)
		// }
		return s, nil
	}

	rel, err := filepath.Rel(r.root, s)
	if err != nil {
		return "", fmt.Errorf("git repo %v %v", err, s)
	}
	return rel, nil
}

func (conf *GitRepo) CopyToRepo(src SrcPath) (RepoPath, error) {
	// Verify source exists and get its info
	srcInfo, err := os.Stat(src.Sting())
	if err != nil {
		return "", fmt.Errorf("copytorepo error stat src: %v", err)
	}

	// Create destination path by appending src path (without leading /) to repo dir
	ret := src.Repo(RepoRoot(conf.root))

	reporoot := RepoRoot(conf.root)
	dest := reporoot.With(ret.Sting())
	// Copy based on whether src is a file or directory
	if srcInfo.IsDir() {
		err = copyDir(src.Sting(), dest)
	} else {
		err = copyFile(src.Sting(), dest)
	}

	if err != nil {
		return "", fmt.Errorf("copytorepo error copying: %v", err)
	}
	return ret, nil
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

func (r GitRepo) AbsRepo2Repo(s string) RepoPath {
	rel, err := filepath.Rel(r.root, s)
	if err != nil {
		return ""
	}
	return RepoPath(rel)
}

func (r GitRepo) Src2Repo(s string) RepoPath {
	rel, err := filepath.Rel("/", s)
	if err != nil {
		return ""
	}
	return RepoPath(rel)
}

//	func (r GitRepo) Covert2Repo(s RepoPath) string {
//		b := RepoRoot(r.root)
//		return b.With(s.Sting())
//	}
// func (r GitRepo) AbsOfRepo(s RepoPath) string {
// 	b := RepoRoot(r.root)
// 	return b.With(s.Sting())
// }

func NewGitReop() (*GitRepo, error) {
	ret := &GitRepo{}
	if err := ret.load(); err != nil {
		return nil, err
	}
	if err := ret.Init(); err != nil {
		return nil, err
	}
	if r, err := ret.Open(); err != nil {
		return nil, err
	} else {
		ret.repo = r
	}
	return ret, nil
}

func (r GitRepo) Open() (*git.Repository, error) {
	if r.repo != nil {
		return r.repo, nil
	}
	git, err := git.PlainOpen(r.root)
	if err != nil {
		return nil, fmt.Errorf("open git repo %v %v", err, r.root)
	}
	return git, nil
}

func (r GitRepo) Init() error {
	if _, err := os.Stat(r.root); err != nil {
		return fmt.Errorf("git repo %v already exists", r.root)
	}
	if _, err := os.Stat(filepath.Join(r.root, ".git")); err == nil {
		return nil
	}
	_, err := git.PlainInit(r.root, false)
	if err != nil {
		return fmt.Errorf("init git repo %v %v", err, r.root)
	}
	return nil
}

func (r GitRepo) GitRmFile(real_path RepoPath) (GitResult, error) {
	add := GitResult{
		Action: GitResultTypeError,
	}
	repo, err := r.Open()
	if err != nil {
		return add, fmt.Errorf("git rm err=%v file=%v", err, real_path)
	}
	w, err := repo.Worktree()
	if err != nil {
		return add, fmt.Errorf("git rm err=%v file=%v", err, real_path)
	}

	abspath := real_path.ToAbs(r)
	os.Remove(abspath)

	// gitfile := real_path.Sting()
	state, err := r.Status(real_path)
	if err != nil {
		return add, fmt.Errorf("git rm err=%v file=%v", err, real_path)
	}
	state.print("before")
	if !state.NeedGitRm() {
		add.Action = GitResultTypeNochange
		return add, nil
	}
	_, err = w.Remove(string(real_path))
	if err != nil {
		return add, fmt.Errorf("git rm err=%v file=%v", err, real_path)
	}

	afterState, err := r.Status(real_path)
	if err != nil {
		return add, fmt.Errorf("git rm err=%v file=%v", err, real_path)
	}
	if action := afterState.NeedGitCommit(); action == "" {
		add.Action = GitResultTypeNochange
		return add, nil
	} else {
		msg := fmt.Sprintf("RM %v", real_path)
		_, err = w.Commit(msg, &git.CommitOptions{
			Author: &object.Signature{
				Name: "anybakup",
				When: time.Now(),
			},
		})
		if err != nil {
			return add, fmt.Errorf("git commit err=%v file=%v:%v", err, real_path, real_path)
		}
		add.Action = GitResultTypeRm
		return add, nil
	}
}

type GitResult struct {
	Action GitAction
	Files  []string
}
type GitAction string

const (
	GitResultTypeAdd      GitAction = "add"
	GitResultTypeRm       GitAction = "rm"
	GitResultTypeNochange GitAction = "nochange"
	GitResultTypeError    GitAction = "error"
)

func (r GitRepo) GitAddFile(gitpath RepoPath) (GitResult, error) {
	abspath := gitpath.ToAbs(r)
	// gitfile := gitpath.Sting()
	ret := GitResult{
		Action: GitResultTypeError,
	}
	repo, err := r.Open()
	if err != nil {
		return ret, fmt.Errorf("git add %v", err)
	}
	w, err := repo.Worktree()
	if err != nil {
		return ret, fmt.Errorf("git add %v %v", err, abspath)
	}
	fmt.Println("-----------------Before---------------")
	state, err := r.Status(gitpath)
	if err != nil {
		return ret, fmt.Errorf("git add %v", err)
	}
	state.print("before")
	for _, v := range state.NeedGitAddFiles() {
		fmt.Printf(">>>>>find need to change %v\n", v)
	}
	if !state.NeedGitAdd() {
		ret.Action = GitResultTypeNochange
		return ret, nil
	}
	_, err = w.Add(gitpath.Sting())
	if err != nil {
		return ret, fmt.Errorf("git add %v %v", err, abspath)
	}
	fmt.Println("-----------------after----------------")

	state, err = r.Status(gitpath)
	if err != nil {
		return ret, fmt.Errorf("git add %v", err)
	}
	state.print("after")
	action := state.NeedGitCommit()
	fmt.Printf("action %s\n", action)
	if action == "" {
		ret.Action = GitResultTypeNochange
		return ret, nil
	}
	msg := fmt.Sprintf("%v %v", action, gitpath)
	options := []git.StatusCode{git.Added, git.Modified}
	ret.Files = state.NeedGitCommitFiles(options)
	for _, k := range ret.Files {
		fmt.Printf(">>>>need to commit %v\n", k)
	}
	_, err = w.Commit(msg, &git.CommitOptions{
		Author: &object.Signature{
			Name: "anybakup",
			When: time.Now(),
		},
	})
	if err != nil {
		return ret, fmt.Errorf("git commit %v %v", err, abspath)
	}
	ret.Action = GitResultTypeAdd
	return ret, nil
}

func (r GitRepo) GitViewFile(gitpath RepoPath, commitHash string, outpath string) (string, error) {
	repo, err := r.Open()
	if err != nil {
		return "", fmt.Errorf("git view file: failed to open repo: %v", err)
	}

	gitfile := gitpath.Sting()

	// Parse commit hash
	hash := plumbing.NewHash(commitHash)

	// Get commit object
	commit, err := repo.CommitObject(hash)
	if err != nil {
		return "", fmt.Errorf("git view file: failed to get commit %s: %v", commitHash, err)
	}

	// Get tree from commit
	tree, err := commit.Tree()
	if err != nil {
		return "", fmt.Errorf("git view file: failed to get tree: %v", err)
	}

	// Get file from tree
	file, err := tree.File(gitfile)
	if err != nil {
		return "", fmt.Errorf("git view file: file %s not found in commit %s: %v", gitfile, commitHash, err)
	}

	// Read file contents
	contents, err := file.Contents()
	if err != nil {
		return "", fmt.Errorf("git view file: failed to read file contents: %v", err)
	}

	// Create output directory if it doesn't exist
	outDir := filepath.Dir(outpath)
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return "", fmt.Errorf("git view file: failed to create output directory: %v", err)
	}

	// Write contents to output file
	if err := os.WriteFile(outpath, []byte(contents), 0o644); err != nil {
		return "", fmt.Errorf("git view file: failed to write output file: %v", err)
	}

	return outpath, nil
}

// GitDiffFile compares a file between working directory and HEAD commit
// Returns the diff as a string, or empty string if no changes
func (r GitRepo) GitDiffFile(file string) (string, error) {
	repo, err := r.Open()
	if err != nil {
		return "", fmt.Errorf("git diff %v", err)
	}

	gitfile, err := r.Rel(file)
	if err != nil {
		return "", fmt.Errorf("git diff %v %v", err, file)
	}

	// Get HEAD commit
	ref, err := repo.Head()
	if err != nil {
		return "", fmt.Errorf("git diff: failed to get HEAD: %v", err)
	}

	commit, err := repo.CommitObject(ref.Hash())
	if err != nil {
		return "", fmt.Errorf("git diff: failed to get commit: %v", err)
	}

	// Get the tree from HEAD commit
	headTree, err := commit.Tree()
	if err != nil {
		return "", fmt.Errorf("git diff: failed to get tree: %v", err)
	}

	// Get worktree
	w, err := repo.Worktree()
	if err != nil {
		return "", fmt.Errorf("git diff: failed to get worktree: %v", err)
	}

	// Get current worktree status
	status, err := w.Status()
	if err != nil {
		return "", fmt.Errorf("git diff: failed to get status: %v", err)
	}

	// Check if file has changes
	fileStatus := status.File(gitfile)
	if fileStatus.Worktree == git.Unmodified && fileStatus.Staging == git.Unmodified {
		return "", nil // No changes
	}

	// Get file content from HEAD
	headFile, err := headTree.File(gitfile)
	var headContent string
	if err == nil {
		headContent, err = headFile.Contents()
		if err != nil {
			return "", fmt.Errorf("git diff: failed to read HEAD content: %v", err)
		}
	} else {
		// File doesn't exist in HEAD (new file)
		headContent = ""
	}

	// Get file content from working directory
	absPath := filepath.Join(r.root, gitfile)
	workingContent, err := os.ReadFile(absPath)
	if err != nil {
		return "", fmt.Errorf("git diff: failed to read working file: %v", err)
	}

	// Simple diff representation
	diff := ""
	if headContent != string(workingContent) {
		diff += fmt.Sprintf("--- a/%s\n+++ b/%s\n", gitfile, gitfile)
		diff += "@@ File changed @@\n"
		diff += fmt.Sprintf("- HEAD: %d bytes\n", len(headContent))
		diff += fmt.Sprintf("+ Working: %d bytes\n", len(workingContent))
	}

	return diff, nil
}

type GitChanges struct {
	Commit  string
	Author  string
	Date    string
	Message string
}

// GitLogFile retrieves the commit history for a specific file
// Returns a formatted string with commit logs
func (r GitRepo) GitLogFile(repoRelPath RepoPath) ([]GitChanges, error) {
	repo := r.repo
	gitfile := repoRelPath.Sting()

	// Get commit log with file path filter
	commitIter, err := repo.Log(&git.LogOptions{
		FileName: &gitfile,
	})
	if err != nil {
		return nil, fmt.Errorf("git changes: failed to get log: %v", err)
	}

	commitCount := 0
	ret := []GitChanges{}
	// Iterate through commits
	err = commitIter.ForEach(func(c *object.Commit) error {
		commitCount++
		r := GitChanges{
			Commit:  c.Hash.String(),
			Author:  c.Author.Name,
			Date:    c.Author.When.Format("2006-01-02 15:04:05"),
			Message: c.Message,
		}
		ret = append(ret, r)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("git changes: failed to iterate commits: %v", err)
	}

	if commitCount == 0 {
		return nil, fmt.Errorf("no commits found for file: %s", gitfile)
	}

	return ret, nil
}

package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"anybakup/util"
)

type GitCmd struct {
	C *util.Config
}

func NewGitCmd(profilname string) *GitCmd {
	c := util.Config{}
	c.Load()
	if config := c.GetProfile(profilname); config != nil {
		return &GitCmd{
			C: config,
		}
	} else {
		return &GitCmd{C: &c}
	}
}

// GetFileLogAbs returns the git log for a specific file
func (g GitCmd) GetFileLogAbs(filePath string) ([]util.GitChanges, error) {
	absFilePath, err := filepath.Abs(filePath)
	if err != nil {
		return nil, err
	}
	repo, err := util.NewGitReop(g.C)
	if err != nil {
		return nil, err
	}
	return g.GetFileLog(repo.Src2Repo(absFilePath))
}

func (g GitCmd) GetFileLog(filePath util.RepoPath) ([]util.GitChanges, error) {
	repo, err := util.NewGitReop(g.C)
	if err != nil {
		return nil, err
	}
	logs, err := repo.GitLogFile(filePath)
	if err != nil {
		return nil, err
	}
	return logs, nil
}

type Result_git_add struct {
	Dest   util.RepoPath
	Err    error
	Result util.GitAction
}

// AddFile adds a file to the git repository
func (g GitCmd) AddFile(arg string) (ret Result_git_add) {
	file, err := filepath.Abs(arg)
	if err != nil {
		ret.Err = err
		return
	}
	ret = Result_git_add{
		Dest:   "",
		Err:    nil,
		Result: util.GitResultTypeError,
	}
	repo, err := util.NewGitReop(g.C)
	if err != nil {
		ret.Err = err
		return
	}
	dest, err := repo.CopyToRepo(util.SrcPath(file))
	if err != nil {
		ret.Err = err
		return
	}
	if yes, err := repo.GitAddFile(dest); err != nil {
		ret.Err = err
	} else {
		ret.Result = yes.Action
		fmt.Printf(">>> result %v\n", ret.Result)
		switch yes.Action {
		case util.GitResultTypeAdd:
			ret.Dest = dest
		case util.GitResultTypeNochange:
			ret.Dest = dest
		default:
			ret.Err = fmt.Errorf("add unexpected result %v", yes)
			return
		}
		isfile, err := IsFile(file)
		fmt.Printf("isfile %v %v\n", isfile, err)
		if err != nil {
			ret.Err = err
			return
		}
		if err := BakupOptAdd(file, ret.Dest, isfile, false, g); err != nil {
			fmt.Printf("failed to add sql backup record %v", err)
		}
		if !isfile {
			for _, f := range yes.Files {
				if err := BakupOptAdd(fmt.Sprintf("/%v", f), f, true, true, g); err != nil {
					fmt.Printf("failed to add sql backup record %v", err)
				} else {
					fmt.Printf(">>> added to sql %v\n", f)
				}
			}
		}
	}
	return
}

// RmFileAbs removes a file from the git repository using an absolute path
func (g GitCmd) RmFileAbs(arg string) error {
	file, err := filepath.Abs(arg)
	if err != nil {
		return err
	}
	repo, err := util.NewGitReop(g.C)
	if err != nil {
		return err
	}
	gitPath := repo.Src2Repo(file)
	return g.RmFile(gitPath)
}

// RmFile removes a file from the git repository using a repository path
func (g GitCmd) RmFile(gitPath util.RepoPath) error {
	repo, err := util.NewGitReop(g.C)
	if err != nil {
		return err
	}
	if yes, err := repo.GitRmFile(gitPath); err != nil {
		return err
	} else {
		switch yes.Action {
		case util.GitResultTypeRm:
			break
		case util.GitResultTypeNochange:
			break
		default:
			return fmt.Errorf("rm unexpected result %v", yes)
		}
		if err := BakupOptRm(gitPath, g.C); err != nil {
			fmt.Println(err, gitPath)
		}
		for _, v := range yes.Files {
			if err := BakupOptRm(v, g.C); err != nil {
				fmt.Println(err, v)
			}
		}
		for _, v := range yes.Dirs {
			if err := BakupOptRm(v, g.C); err != nil {
				fmt.Println(err, v)
			}
		}
		return nil
	}
}

func IsFile(file string) (bool, error) {
	if st, err := os.Stat(file); err != nil {
		return true, err
	} else {
		return st.Mode().IsRegular(), nil
	}
}

// GetFile retrieves a file from a specific commit
func (g GitCmd) GetFile(filePath string, commit string, target string) error {
	var err error
	absFilePath, err := filepath.Abs(filePath)
	if err != nil {
		return err
	}
	repo, err := util.NewGitReop(g.C)
	if err != nil {
		return err
	}

	// If commit is empty, use HEAD
	if commit == "" {
		r, err := repo.Open()
		if err != nil {
			return fmt.Errorf("failed to open repo: %v", err)
		}
		ref, err := r.Head()
		if err != nil {
			return fmt.Errorf("failed to get HEAD: %v", err)
		}
		commit = ref.Hash().String()
	}

	_, err = repo.GitViewFile(repo.Src2Repo(absFilePath), commit, target)
	if err != nil {
		return err
	}
	return nil
}

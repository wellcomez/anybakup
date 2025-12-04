package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"anybakup/util"
)

// GetFileLog returns the git log for a specific file
func GetFileLog(filePath string) ([]util.GitChanges, error) {
	absFilePath, err := filepath.Abs(filePath)
	if err != nil {
		return nil, err
	}
	repo, err := util.NewGitReop()
	if err != nil {
		return nil, err
	}
	logs, err := repo.GitLogFile(repo.Src2Repo(absFilePath))
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
func AddFile(arg string) (ret Result_git_add) {
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
	repo, err := util.NewGitReop()
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
			ret.Err = fmt.Errorf("add  unexpected result %v", yes)
			return
		}
		isfile, err := IsFile(file)
		fmt.Printf("isfile %v %v\n", isfile, err)
		if err != nil {
			ret.Err = err
			return
		}
		if err := BakupOptAdd(file, ret.Dest, isfile, false); err != nil {
			fmt.Printf("failed to add sql backup record %v", err)
		}
		if !isfile {
			for _, f := range yes.Files {
				if err := BakupOptAdd(fmt.Sprintf("/%v", f), f, true, true); err != nil {
					fmt.Printf("failed to add sql backup record %v", err)
				} else {
					fmt.Printf(">>> added to sql %v\n", f)
				}
			}
		}
	}
	return
}

func RmFileAbs(arg string) error {
	file, err := filepath.Abs(arg)
	if err != nil {
		return err
	}
	repo, err := util.NewGitReop()
	if err != nil {
		return err
	}
	gitPath := repo.Src2Repo(file)
	return RmFile(gitPath)
}
func RmFile(gitPath util.RepoPath) error {
	repo, err := util.NewGitReop()
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
			return fmt.Errorf("rm unexpected  result %v", yes)
		}
		if err := BakupOptRm(gitPath); err != nil {
			fmt.Println(err)
		}
		for _, v := range yes.Files {
			if err := BakupOptRm(v); err != nil {
				fmt.Println(err)
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
func GetFile(filePath string, commit string, target string) error {
	var err error
	absFilePath, err := filepath.Abs(filePath)
	if err != nil {
		return err
	}
	repo, err := util.NewGitReop()
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

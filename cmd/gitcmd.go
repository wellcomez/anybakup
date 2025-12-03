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
	Dest   string
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
		switch yes.Action {
		case util.GitResultTypeAdd:
			ret.Dest = dest.Sting()
		case util.GitResultTypeNochange:
			ret.Dest = dest.Sting()
		default:
			ret.Err = fmt.Errorf("add  unexpected result %v", yes)
			return
		}
		isfile, err := IsFile(file)
		if err != nil {
			ret.Err = err
			return
		}
		if err := BakupOptAdd(file, ret.Dest, isfile, false); err != nil {
			fmt.Printf("failed to add sql backup record %v", err)
		}
		if !isfile {
			for _, f := range yes.Files {
				fmt.Printf("added %v\n", f)
				if err := BakupOptAdd(fmt.Sprintf("/%v", f), f, true, true); err != nil {
					fmt.Printf("failed to add sql backup record %v", err)
				}
			}
		}
	}
	return
}

func RmFile(arg string) error {
	file, err := filepath.Abs(arg)
	if err != nil {
		return err
	}
	repo, err := util.NewGitReop()
	if err != nil {
		return err
	}
	if yes, err := repo.GitRmFile(repo.Src2Repo(file)); err != nil {
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
	}
	if err := BakupOptRm(file); err != nil {
		fmt.Println(err)
	}
	return nil
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

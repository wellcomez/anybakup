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
	Result util.GitResult
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
		Result: util.GitResultError,
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
		ret.Result = yes
		switch yes {
		case util.GitResultAdd:
			ret.Dest = dest.Sting()
		case util.GitResultNochange:
			ret.Dest = dest.Sting()
		default:
			ret.Err = fmt.Errorf("unknown result %v", yes)
		}
	}
	isfile, err := IsFile(file)
	BackupOptAdd(file, ret.Dest, isfile)
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
		switch yes {
		case util.GitResultRm:
		case util.GitResultNochange:
			if err := BackupOptRm(file); err != nil {
				fmt.Println(err)
			}
			return nil
		default:
			return fmt.Errorf("unknown result %v", yes)
		}
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

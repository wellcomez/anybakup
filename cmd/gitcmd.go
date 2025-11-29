package cmd

import (
	"anybakup/util"
	"fmt"
	"path/filepath"
)

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

type result_git_add struct {
	dest   string
	err    error
	resutl util.GitResult
}

func AddFile(arg string) (ret result_git_add) {
	file, err := filepath.Abs(arg)
	if err != nil {
		ret.err = err
		return
	}
	ret = result_git_add{
		dest:   "",
		err:    nil,
		resutl: util.GitResultError,
	}
	repo, err := util.NewGitReop()
	if err != nil {
		ret.err = err
		return
	}
	dest, err := repo.CopyToRepo(util.SrcPath(file))
	if err != nil {
		ret.err = err
		return
	}
	if yes, err := repo.GitAddFile(dest); err != nil {
		ret.err = err
	} else {
		ret.resutl = yes
		switch yes {
		case util.GitResultAdd:
			ret.dest = dest.Sting()
		case util.GitResultNochange:
			ret.dest = dest.Sting()
		default:
			ret.err = fmt.Errorf("unknown result %v", yes)
		}
	}
	return
}
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
	_, err = repo.GitViewFile(repo.Src2Repo(absFilePath), commit, target)
	if err != nil {
		return err
	}
	return nil
}

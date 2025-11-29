package cmd

import (
	"anybakup/util"
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

package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"anybakup/util"

	"github.com/spf13/cobra"
)

type result_git_add struct {
	dest   string
	err    error
	resutl util.GitResult
}

func add_file(file string) (ret result_git_add) {
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

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add [file]",
	Short: "Add a file to the repository",
	Long:  `Add a file to the repository. This copies the file to the configured repository directory.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filePath := args[0]
		absFilePath, err := filepath.Abs(filePath)
		if err != nil {
			fmt.Printf("Error add file %v: [%v]\n", filePath, err)
			os.Exit(1)
		}
		if ret := add_file(absFilePath); ret.err != nil {
			fmt.Printf("Error add file %v: [%v]\n", filePath, ret.err)
			os.Exit(1)
		} else {
			fmt.Printf("add %s to %s\n", filePath, ret.dest)
		}
	},
}

var logCmd = &cobra.Command{
	Use:   "log [file]",
	Short: "log a file to the repository",
	Long:  `log a file to the repository. This copies the file to the configured repository directory.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filePath := args[0]
		if logs, err := GetFileLog(filePath); err != nil {
			fmt.Printf("Error log file %v: [%v]\n", filePath, err)
			os.Exit(1)
		} else {
			for _, l := range logs {
				fmt.Printf("%-10s %-10s %-10s\n", l.Commit, l.Author, l.Date)
			}
		}
	},
}
var viewCmd = &cobra.Command{
	Use:   "view [file]",
	Short: "view a file to the repository",
	Long:  `view a file to the repository. This copies the file to the configured repository directory.`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		filePath := args[0]
		absFilePath, err := filepath.Abs(filePath)
		if err != nil {
			fmt.Printf("Error log file %v: [%v]\n", filePath, err)
			os.Exit(1)
		}
		repo, err := util.NewGitReop()
		if err != nil {
			fmt.Printf("Error log file %v: [%v]\n", filePath, err)
			os.Exit(1)
		}
		logs, err := repo.GitLogFile(repo.Src2Repo(absFilePath))
		if err != nil {
			fmt.Printf("Error %v", err)
		}
		for _, l := range logs {
			fmt.Printf("%-10s %-10s %-10s\n", l.Commit, l.Author, l.Date)
		}
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(logCmd)
	rootCmd.AddCommand(viewCmd)
}

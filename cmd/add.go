package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"anybakup/util"

	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add [file]",
	Short: "Add a file to the repository",
	Long:  `Add a file to the repository. This copies the file to the configured repository directory.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if ret := AddFile(args[0]); ret.err != nil {
			fmt.Printf("Error add file %v: [%v]\n", args[0], ret.err)
			os.Exit(1)
		} else {
			fmt.Printf("add %s to %s\n", args[0], ret.dest)
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

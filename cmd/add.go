package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"anybakup/util"

	"github.com/spf13/cobra"
)

func add_file(file string) (string, error) {
	repo, err := util.NewGitReop()
	if err != nil {
		return "", err
	}
	dest, err := repo.CopyToRepo(file)
	if err != nil {
		return "", err
	}
	if yes, err := repo.GitAddFile(dest); err != nil {
		return "", err
	} else if !yes {
		return "", fmt.Errorf("no change")
	}
	return dest.Sting(), nil
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
		if dest, err := add_file(absFilePath); err != nil {
			fmt.Printf("Error add file %v: [%v]\n", filePath, err)
			os.Exit(1)
		} else {
			fmt.Printf("add %s to %s\n", filePath, dest)
		}
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}

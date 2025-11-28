package cmd

import (
	"anybakup/util"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// rmCmd represents the rm command
var rmCmd = &cobra.Command{
	Use:   "rm [file]",
	Short: "Remove a file from the repository",
	Long:  `Remove a file from the repository. This does not delete the original file, only the backup.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filePath := args[0]
		absFilePath, err := filepath.Abs(filePath)
		if err != nil {
			fmt.Printf("Error add file %v: [%v]\n", filePath, err)
			os.Exit(1)
		}
		if repo, err := util.NewGitReop(); err != nil {
			fmt.Printf("Error add file %v: [%v]\n", filePath, err)
			os.Exit(1)
		} else {
			yes, err := repo.GitRmFile(repo.PathOfRepo(absFilePath))
			if err != nil {
				fmt.Printf("Error add file %v: [%v]\n", filePath, err)
				os.Exit(1)
			}
			if yes {
				fmt.Printf("rm %s from %s\n", filePath, absFilePath)
			}
		}

	},
}

func init() {
	rootCmd.AddCommand(rmCmd)
}

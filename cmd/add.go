package cmd

import (
	"anybakup/util"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func add_file(file string) (dest string, err error) {
	dest, err = util.CopyToRepo(file)
	if err != nil {
		return
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
			fmt.Printf("Error add file %v: %v\n", filePath, err)
			os.Exit(1)
		}
		if dest, err := add_file(absFilePath); err != nil {
			fmt.Printf("Error add file %v: %v\n", filePath, err)
			os.Exit(1)
		} else {
			fmt.Printf("add %s to %s\n", filePath, dest)
		}
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}

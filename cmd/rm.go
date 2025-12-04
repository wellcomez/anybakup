package cmd

import (
	"fmt"

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
		g:=GitCmd{}
		err := g.RmFileAbs(filePath)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("rm %s OK\n", filePath)
		}

	},
}

func init() {
	rootCmd.AddCommand(rmCmd)
}

package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rmCmd represents the rm command
var rmCmd = &cobra.Command{
	Use:   "rm [file]",
	Short: "Remove a file from the repository",
	Long:  `Remove a file from the repository. This does not delete the original file, only the backup.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filePath := args[0]
		fileName := filepath.Base(filePath)

		repoDir := viper.GetString("repodir")
		if repoDir == "" {
			fmt.Println("Repository directory not configured. Run 'anybakup init <dir>' first.")
			os.Exit(1)
		}

		destPath := filepath.Join(repoDir, fileName)

		if err := os.Remove(destPath); err != nil {
			if os.IsNotExist(err) {
				fmt.Printf("File '%s' is not in the repository.\n", fileName)
			} else {
				fmt.Printf("Error removing file: %v\n", err)
				os.Exit(1)
			}
		} else {
			fmt.Printf("Removed '%s' from repository.\n", fileName)
		}
	},
}

func init() {
	rootCmd.AddCommand(rmCmd)
}

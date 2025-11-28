package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status [file]",
	Short: "Check the status of a file",
	Long:  `Check if a file is tracked in the repository.`,
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

		if _, err := os.Stat(destPath); err == nil {
			fmt.Printf("File '%s' is tracked in repository.\n", fileName)
		} else if os.IsNotExist(err) {
			fmt.Printf("File '%s' is NOT tracked in repository.\n", fileName)
		} else {
			fmt.Printf("Error checking status: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}

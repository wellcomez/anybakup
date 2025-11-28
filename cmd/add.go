package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

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
			fmt.Printf("Error getting absolute path: %v\n", err)
			os.Exit(1)
		}

		repoDir := viper.GetString("repodir")
		if repoDir == "" {
			fmt.Println("Repository directory not configured. Run 'anybakup init <dir>' first.")
			os.Exit(1)
		}

		// Check if source file exists
		fileInfo, err := os.Stat(absFilePath)
		if os.IsNotExist(err) {
			fmt.Printf("File does not exist: %s\n", absFilePath)
			os.Exit(1)
		}
		if fileInfo.IsDir() {
			fmt.Printf("Cannot add a directory: %s\n", absFilePath)
			os.Exit(1)
		}

		// Destination path
		fileName := filepath.Base(absFilePath)
		destPath := filepath.Join(repoDir, fileName)

		// Copy file
		if err := copyFile(absFilePath, destPath); err != nil {
			fmt.Printf("Error copying file: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Added %s to repository\n", fileName)
	},
}

func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destinationFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	if _, err := io.Copy(destinationFile, sourceFile); err != nil {
		return err
	}
	return nil
}

func init() {
	rootCmd.AddCommand(addCmd)
}

package cmd

import (
	"anybakup/util"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func GitInitProfile(profile, repoPath string) error {
	p := util.Profile{
		RepoDir: util.RepoRoot(repoPath),
	}
	c := util.NewConfig()
	if err := c.SetProfile(profile, p); err != nil {
		return fmt.Errorf("error saving config: %v", err)
	}
	// Initialize git repository
	if _, err := util.NewGitReop(c); err != nil {
		return fmt.Errorf("error initializing git repository: %v", err)
	}
	return nil
}

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init [directory]",
	Short: "Initialize a new repository",
	Long:  `Initialize a new repository at the specified directory.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		repoPath := args[0]
		absPath, err := filepath.Abs(repoPath)
		if err != nil {
			fmt.Printf("Error getting absolute path: %v\n", err)
			os.Exit(1)
		}

		// // 1. Create and ensure dir a
		// if err := os.MkdirAll(absPath, 0755); err != nil {
		// 	fmt.Printf("Error creating directory: %v\n", err)
		// 	os.Exit(1)
		// }
		// fmt.Printf("Initialized empty repository in %s\n", absPath)

		// p := util.Profile{
		// 	RepoDir: util.RepoRoot(absPath),
		// }
		// c := util.NewConfig()
		// if err := c.SetProfile("", p); err != nil {
		// 	fmt.Printf("Error saving config: %v\n", err)
		// 	os.Exit(1)
		// }
		// // Initialize git repository
		// if _, err := util.NewGitReop(c); err != nil {
		// 	fmt.Printf("Error initializing git repository: %v\n", err)
		// 	os.Exit(1)
		// }

		if err := GitInitProfile("", absPath); err != nil {
			fmt.Printf("Error initializing profile: %v\n", err)
			os.Exit(1)
		}

	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

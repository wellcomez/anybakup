package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"anybakup/util"

	"github.com/spf13/cobra"
)

func GitInitProfile(profile, repoPath string) (*util.GitRepo, error) {
	if err := os.MkdirAll(repoPath, 0o755); err != nil {
		fmt.Printf("Error creating directory: %v\n", err)
		return nil, err
	}
	p := util.Profile{
		RepoDir: util.RepoRoot(repoPath),
	}
	c := util.NewConfig()
	if err := c.SetProfile(profile, p); err != nil {
		c.Print()
		return nil, fmt.Errorf("error saving config: %v", err)
	}
	c.Print()
	// Initialize git repository
	if ret, err := util.NewGitReop(c); err != nil {
		return ret, fmt.Errorf("error initializing git repository: %v", err)
	} else {
		return ret, nil
	}
}

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init [directory] [profile]",
	Short: "Initialize a new repository",
	Long:  `Initialize a new repository at the specified directory with an optional profile.`,
	Args:  cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		repoPath := args[0]
		absPath, err := filepath.Abs(repoPath)
		if err != nil {
			fmt.Printf("Error getting absolute path: %v\n", err)
			os.Exit(1)
		}

		profile := ""
		if len(args) > 1 {
			profile = args[1]
		}

		if ret, err := GitInitProfile(profile, absPath); err != nil {
			fmt.Printf("Error initializing profile: %v\n", err)
		} else if ret != nil {
		}
		os.Exit(1)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

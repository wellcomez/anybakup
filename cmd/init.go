package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

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

		// 1. Create and ensure dir a
		if err := os.MkdirAll(absPath, 0755); err != nil {
			fmt.Printf("Error creating directory: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Initialized empty repository in %s\n", absPath)

		// 2. Create config file ~/.config/anybakup/config.yaml
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Printf("Error getting home directory: %v\n", err)
			os.Exit(1)
		}

		configDir := filepath.Join(home, ".config", "anybakup")
		if err := os.MkdirAll(configDir, 0755); err != nil {
			fmt.Printf("Error creating config directory: %v\n", err)
			os.Exit(1)
		}

		viper.Set("repodir", absPath)
		configFilePath := filepath.Join(configDir, "config.yaml")
		if err := viper.WriteConfigAs(configFilePath); err != nil {
			fmt.Printf("Error writing config file: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Configuration saved to %s\n", configFilePath)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

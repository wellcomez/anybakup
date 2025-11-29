package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add [file]",
	Short: "Add a file to the repository",
	Long:  `Add a file to the repository. This copies the file to the configured repository directory.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if ret := AddFile(args[0]); ret.err != nil {
			fmt.Printf("Error add file %v: [%v]\n", args[0], ret.err)
			os.Exit(1)
		} else {
			fmt.Printf("add %s to %s\n", args[0], ret.dest)
		}
	},
}
func run_list_file(filePath string) {
	if logs, err := GetFileLog(filePath); err != nil {
		fmt.Printf("Error log file %v: [%v]\n", filePath, err)
		os.Exit(1)
	} else {
		for i, l := range logs {
			fmt.Printf("%3d: %-10s %-10s %-10s\n", i, l.Commit, l.Author, l.Date)
		}
	}
}
var logCmd = &cobra.Command{
	Use:   "log [file]",
	Short: "log a file to the repository",
	Long:  `log a file to the repository. This copies the file to the configured repository directory.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filePath := args[0]
		run_list_file(filePath)
	},
}
var getCmd = &cobra.Command{
	Use:   "get [file] [target] [commit]",
	Short: "get a file from the repository",
	Long:  `get a file from the repository. This copies the file to the configured repository directory.`,
	Args:  cobra.RangeArgs(1, 3),
	Run: func(cmd *cobra.Command, args []string) {
		filePath := args[0]
		target := ""
		commit := ""
		if len(args) == 3 {
			target = args[1]
			commit = args[2]
		}else{
			run_list_file(filePath)	
			os.Exit(1)
		}
		if err := GetFile(filePath, commit, target); err != nil {
			fmt.Printf("Error get file %v: [%v]\n", filePath, err)
			os.Exit(1)
		} else {
			if commit == "" {
				fmt.Printf("get %s from HEAD to %s\n", filePath, target)
			} else {
				fmt.Printf("get %s from %s to %s\n", filePath, commit, target)
			}
			os.Exit(0)
		}
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(logCmd)
	rootCmd.AddCommand(getCmd)
}

package cmd

import (
	"anybakup/util"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

)
func add_file(file string)error{
	dest,err:=util.CopyToRepo(file)
	if err!=nil{
		return err
	}
	fmt.Printf("Added %s -> %s to repository\n",file, dest)
	return nil
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
			fmt.Printf("Error getting absolute path: %v\n", err)
			os.Exit(1)
		}
		add_file(absFilePath)
	},
}



func init() {
	rootCmd.AddCommand(addCmd)
}

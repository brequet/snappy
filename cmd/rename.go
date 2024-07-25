package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var renameCmd = &cobra.Command{
	Use:   "rename <snapshot-name> <new-snapshot-name>",
	Short: "rename a snapshot",
	Long:  `rename a snapshot`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		err := snapshotService.RenameSnapshot(args[0], args[1])
		if err != nil {
			fmt.Println("Error renaming snapshot:", err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(renameCmd)
}

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "remove <snapshot-name>",
	Short: "remove a snapshot",
	Long:  `remove a snapshot`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := snapshotService.RemoveSnapshot(args[0])
		if err != nil {
			fmt.Println("Error removing snapshot:", err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
}

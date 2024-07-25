package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var restoreCmd = &cobra.Command{
	Use:   "restore <snapshot-name>",
	Short: "restore a snapshot",
	Long:  `restore a snapshot`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := snapshotService.RestoreSnapshot(args[0])
		if err != nil {
			fmt.Println("Error restoring snapshot:", err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(restoreCmd)
}

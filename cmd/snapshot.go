package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var snapshotCmd = &cobra.Command{
	Use:   "snapshot <db> <snapshot-name>",
	Short: "snapshot a database",
	Long:  `snapshot a database`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		dbName := args[0]
		snapshotName := args[1]

		err := snapshotService.CreateSnapshot(dbName, snapshotName)
		if err != nil {
			fmt.Println("Error creating snapshot:", err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(snapshotCmd)
}

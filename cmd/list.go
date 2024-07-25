package cmd

import (
	"fmt"
	"time"

	"github.com/brequet/snappy/entity"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list snapshots",
	Long:  `list snapshots`,
	Run: func(cmd *cobra.Command, args []string) {
		names, err := snapshotService.ListSnapshots()
		if err != nil {
			fmt.Println("Error listing snapshots:", err)
			return
		}

		dbToSnapshots := make(map[string][]entity.Snapshot)
		for _, snapshot := range names {
			dbToSnapshots[snapshot.ReferenceDb] = append(dbToSnapshots[snapshot.ReferenceDb], snapshot)
		}

		for db, snapshots := range dbToSnapshots {
			fmt.Printf("Snapshots for %s:\n", db)
			for _, snapshot := range snapshots {
				fmt.Printf("  - %-20s (last updated: %s)\n", snapshot.Name, snapshot.UpdatedAt.Format(time.RFC3339))
			}
			fmt.Println()
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}

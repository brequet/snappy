package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/brequet/snappy/config"
	"github.com/brequet/snappy/database"
	"github.com/brequet/snappy/repository"
	"github.com/brequet/snappy/snapshot"
	"github.com/spf13/cobra"
)

var (
	snapshotService *snapshot.SnapshotService
	userFlag        string
	hostFlag        string
	portFlag        string
)

var rootCmd = &cobra.Command{
	Use:   config.PROGRAM_NAME,
	Short: "A CLI tool to handle snapshot for postgres databases.",
	Long:  `A CLI tool to handle snapshot for postgres databases.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		initialize()
	},
	Version: config.Version,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// disabling default -h flag
	rootCmd.SetHelpCommand(&cobra.Command{Hidden: true})
	rootCmd.PersistentFlags().BoolP("help", "", false, "Display help")

	rootCmd.PersistentFlags().StringVarP(&userFlag, "username", "U", "", "PostgreSQL user (overrides PGUSER environment variable)")
	rootCmd.PersistentFlags().StringVarP(&hostFlag, "host", "h", "", "PostgreSQL host")
	rootCmd.PersistentFlags().StringVarP(&portFlag, "port", "p", "", "PostgreSQL port")
}

func initialize() {
	pgConfig := config.GetPostgresConfig()

	if userFlag != "" {
		pgConfig.User = userFlag
	}
	if hostFlag != "" {
		pgConfig.Host = hostFlag
	}
	if portFlag != "" {
		pgConfig.Port = portFlag
	}

	postgres, err := database.NewPostgres(pgConfig)
	if err != nil {
		fmt.Println(err)
		fmt.Println()

		if strings.Contains(err.Error(), "password authentication failed") {
			fmt.Println("You can use the environment variable 'PGPASSWORD' to set the password.")
		}

		fmt.Println()
		os.Exit(1)
	}

	snappyRepository := repository.NewSnappyRepository(postgres)
	postgresRepository := repository.NewPostgresRepository(postgres)

	snapshotService = snapshot.NewSnapshotService(snappyRepository, postgresRepository)
}

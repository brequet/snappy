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

	rootCmd.PersistentFlags().StringVarP(&userFlag, "username", "U", "", "PostgreSQL user")
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
			if pgConfig.User == "postgres" {
				fmt.Println("If you do not want to use the default user 'postgres', please set the environment variable 'PGUSER' to the correct user or use the -U flag to specify the user.")
			}
			// if !config.IsPgPasswordSet() {
			// 	fmt.Printf("ariba\n")
			// }
			fmt.Println("To use a password, please set the environment variable 'PGPASSWORD' to the correct password.")
		}

		fmt.Println()
		os.Exit(1)
	}

	snappyRepository := repository.NewSnappyRepository(postgres)
	postgresRepository := repository.NewPostgresRepository(postgres)

	snapshotService = snapshot.NewSnapshotService(snappyRepository, postgresRepository)
}

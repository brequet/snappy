package cmd

import (
	"fmt"
	"os"

	"github.com/brequet/snappy/config"
	"github.com/brequet/snappy/database"
	"github.com/brequet/snappy/repository"
	"github.com/brequet/snappy/snapshot"
	"github.com/spf13/cobra"
)

var (
	snapshotService *snapshot.SnapshotService
)

var rootCmd = &cobra.Command{
	Use:   config.PROGRAM_NAME,
	Short: "A CLI tool to handle snapshot for postgres databases.",
	Long:  `A CLI tool to handle snapshot for postgres databases.`,
}

// TODO: handle when password or psql access not trusted

func Execute() {
	initialize()

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func initialize() {
	postgres, err := database.NewPostgres(config.GetPostgresConfig())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	snappyRepository := repository.NewSnappyRepository(postgres)
	postgresRepository := repository.NewPostgresRepository(postgres)

	snapshotService = snapshot.NewSnapshotService(snappyRepository, postgresRepository)
}

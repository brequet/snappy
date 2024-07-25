package config

import (
	"os"
)

var Version string

const (
	PROGRAM_NAME = "snappy"
)

func IsPgPasswordSet() bool {
	return os.Getenv("PGPASSWORD") != ""
}

func GetPgEnvVars() (string, string) {
	pgUser := os.Getenv("PGUSER")
	if pgUser == "" {
		pgUser = os.Getenv("USERNAME")
	}

	pgPassword := os.Getenv("PGPASSWORD")
	if pgPassword == "" {
		pgPassword = ""
	}

	return pgUser, pgPassword
}

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

func GetPostgresConfig() PostgresConfig {
	pgUser, pgPassword := GetPgEnvVars()

	return PostgresConfig{
		Host:     "localhost",
		Port:     "5432",
		User:     pgUser,
		Password: pgPassword,
		Database: PROGRAM_NAME,
	}
}

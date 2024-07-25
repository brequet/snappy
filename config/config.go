package config

import (
	"os"
)

const (
	PROGRAM_NAME = "snappy"
)

func GetPgEnvVars() (string, string) {
	pgUser := os.Getenv("PGUSER")
	if pgUser == "" {
		pgUser = "postgres"
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

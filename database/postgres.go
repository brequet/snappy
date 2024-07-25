package database

import (
	"context"
	"fmt"

	"github.com/brequet/snappy/config"
	"github.com/jackc/pgx/v4"
)

type Postgres struct {
	Conn *pgx.Conn
}

func NewPostgres(pgConfig config.PostgresConfig) (*Postgres, error) {
	err := createSnappyDatabase(pgConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize snappy database: %w", err)
	}

	connString := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable", pgConfig.User, pgConfig.Password, pgConfig.Host, pgConfig.Port, pgConfig.Database)
	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}

	pg := &Postgres{Conn: conn}

	err = pg.initSnappyDatabase()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize snappy database: %w", err)
	}

	return pg, nil
}

func createSnappyDatabase(pgConfig config.PostgresConfig) error {
	defaultConnString := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=postgres sslmode=disable", pgConfig.User, pgConfig.Password, pgConfig.Host, pgConfig.Port)
	defaultConn, err := pgx.Connect(context.Background(), defaultConnString)
	if err != nil {
		return fmt.Errorf("failed to connect to postgres: %w", err)
	}
	defer defaultConn.Close(context.Background())

	snappyExists, err := databaseExists(defaultConn, pgConfig.Database)
	if err != nil {
		return fmt.Errorf("failed to check if snappy database exists: %w", err)
	}

	if !snappyExists {
		_, err := defaultConn.Exec(context.Background(), fmt.Sprintf("CREATE DATABASE %s", pgConfig.Database))
		if err != nil {
			return fmt.Errorf("failed to create snappy database: %w", err)
		}
	}

	return nil
}

func (s *Postgres) initSnappyDatabase() error {
	_, err := s.Conn.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS snapshots (
			id SERIAL PRIMARY KEY,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			name TEXT UNIQUE,
			reference_db TEXT,
			snapshot_db TEXT
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create snapshots table: %w", err)
	}

	return nil
}

func databaseExists(conn *pgx.Conn, name string) (bool, error) {
	res, err := conn.Query(context.Background(), fmt.Sprintf("SELECT 1 FROM pg_database WHERE datname = '%s'", name))
	if err != nil {
		return false, fmt.Errorf("failed to check if database exists: %w", err)
	}
	defer res.Close()

	return res.Next(), nil
}

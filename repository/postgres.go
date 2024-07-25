package repository

import (
	"context"
	"fmt"

	"github.com/brequet/snappy/database"
)

type PostgresRepository struct {
	db *database.Postgres
}

func NewPostgresRepository(db *database.Postgres) *PostgresRepository {
	return &PostgresRepository{
		db: db,
	}
}

func (r *PostgresRepository) ListDatabases() ([]string, error) {
	rows, err := r.db.Conn.Query(context.Background(), "SELECT datname FROM pg_database WHERE datistemplate = false;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var databases []string
	for rows.Next() {
		var dbName string
		if err := rows.Scan(&dbName); err != nil {
			return nil, err
		}
		databases = append(databases, dbName)
	}

	return databases, nil
}

func (r *PostgresRepository) CreateDatabase(dbName, templateName string) error {
	_, err := r.db.Conn.Exec(context.Background(), fmt.Sprintf("CREATE DATABASE %s TEMPLATE %s;", templateName, dbName))
	return err
}

func (r *PostgresRepository) DropDatabase(dbName string) error {
	_, err := r.db.Conn.Exec(context.Background(), fmt.Sprintf("DROP DATABASE IF EXISTS %s;", dbName))
	return err
}

func (r *PostgresRepository) RenameDatabase(oldName, newName string) error {
	_, err := r.db.Conn.Exec(context.Background(), fmt.Sprintf("ALTER DATABASE %s RENAME TO %s;", oldName, newName))
	return err
}

func (r *PostgresRepository) IsDatabaseInUse(dbName string) (bool, error) {
	var count int
	err := r.db.Conn.QueryRow(context.Background(), fmt.Sprintf("SELECT count(*) FROM pg_stat_activity WHERE datname = '%s' AND pid <> pg_backend_pid();", dbName)).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *PostgresRepository) TerminateAllConnections(dbName string) error {
	_, err := r.db.Conn.Exec(context.Background(), fmt.Sprintf("SELECT pg_terminate_backend(pid) FROM pg_stat_activity WHERE datname = '%s' AND pid <> pg_backend_pid();", dbName))
	return err
}

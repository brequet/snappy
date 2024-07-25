package repository

import (
	"context"
	"time"

	"github.com/brequet/snappy/database"
	"github.com/brequet/snappy/entity"
)

type SnappyRepository struct {
	db *database.Postgres
}

func NewSnappyRepository(db *database.Postgres) *SnappyRepository {
	return &SnappyRepository{
		db: db,
	}
}

func (s *SnappyRepository) CreateSnapshot(name, referenceDb, snapshotDb string) error {
	_, err := s.db.Conn.Exec(context.Background(),
		"INSERT INTO snapshots (name, reference_db, snapshot_db) VALUES ($1, $2, $3)",
		name, referenceDb, snapshotDb)
	return err
}

func (s *SnappyRepository) GetAllSnapshots() ([]entity.Snapshot, error) {
	rows, err := s.db.Conn.Query(context.Background(),
		"SELECT id, created_at, updated_at, name, reference_db, snapshot_db FROM snapshots")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var snapshots []entity.Snapshot
	for rows.Next() {
		var snapshot entity.Snapshot
		err := rows.Scan(&snapshot.ID, &snapshot.CreatedAt, &snapshot.UpdatedAt, &snapshot.Name, &snapshot.ReferenceDb, &snapshot.SnapshotDb)
		if err != nil {
			return nil, err
		}
		snapshots = append(snapshots, snapshot)
	}
	return snapshots, nil
}

func (s *SnappyRepository) GetSnapshotByName(name string) (*entity.Snapshot, error) {
	var snapshot entity.Snapshot
	err := s.db.Conn.QueryRow(context.Background(),
		"SELECT id, created_at, updated_at, name, reference_db, snapshot_db FROM snapshots WHERE name = $1",
		name).Scan(&snapshot.ID, &snapshot.CreatedAt, &snapshot.UpdatedAt, &snapshot.Name, &snapshot.ReferenceDb, &snapshot.SnapshotDb)
	if err != nil {
		return nil, err
	}
	return &snapshot, nil
}

func (s *SnappyRepository) DeleteSnapshot(name string) error {
	_, err := s.db.Conn.Exec(context.Background(), "DELETE FROM snapshots WHERE name = $1", name)
	return err
}

func (s *SnappyRepository) RenameSnapshot(oldName, newName string) error {
	_, err := s.db.Conn.Exec(context.Background(),
		"UPDATE snapshots SET name = $1, updated_at = $2 WHERE name = $3",
		newName, time.Now(), oldName)
	return err
}

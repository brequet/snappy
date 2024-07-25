package snapshot

import (
	"fmt"

	"github.com/brequet/snappy/entity"
	"github.com/brequet/snappy/repository"
	"github.com/oklog/ulid/v2"
)

type SnapshotService struct {
	snappyRepository   *repository.SnappyRepository
	postgresRepository *repository.PostgresRepository
}

func NewSnapshotService(
	snappyRepository *repository.SnappyRepository,
	postgresRepository *repository.PostgresRepository,
) *SnapshotService {
	return &SnapshotService{
		snappyRepository:   snappyRepository,
		postgresRepository: postgresRepository,
	}
}

func (s *SnapshotService) ListSnapshots() ([]entity.Snapshot, error) {
	return s.snappyRepository.GetAllSnapshots()
}

func (s *SnapshotService) CreateSnapshot(sourceDB, snapshotName string) error {
	if sourceDB == "" {
		return fmt.Errorf("source database cannot be empty")
	}

	if snapshotName == "" {
		return fmt.Errorf("snapshot name cannot be empty")
	}

	err := s.promptToStopConnectionsIfNeeded(sourceDB)
	if err != nil {
		return fmt.Errorf("failed to stop connections: %w", err)
	}

	generatedDbName := generateRandomSnapshotName()
	err = s.postgresRepository.CreateDatabase(sourceDB, generatedDbName)
	if err != nil {
		return fmt.Errorf("failed to copy database: %w", err)
	}

	err = s.snappyRepository.CreateSnapshot(snapshotName, sourceDB, generatedDbName)
	if err != nil {
		return fmt.Errorf("failed to save snapshot metadata: %w", err)
	}

	fmt.Printf("Snapshot '%s' created successfully from database '%s'\n", snapshotName, sourceDB)
	return nil
}

func (s *SnapshotService) RestoreSnapshot(snapshotName string) error {
	snapshot, err := s.snappyRepository.GetSnapshotByName(snapshotName)
	if err != nil {
		// TODO: check if record not found: change message error, tell snapshot does not exist
		return fmt.Errorf("failed to get snapshot: %w", err)
	}

	err = s.promptToStopConnectionsIfNeeded(snapshot.ReferenceDb)
	if err != nil {
		return fmt.Errorf("failed to stop connections: %w", err)
	}

	err = s.postgresRepository.DropDatabase(snapshot.ReferenceDb)
	if err != nil {
		return fmt.Errorf("failed to drop source db: %w", err)
	}

	err = s.promptToStopConnectionsIfNeeded(snapshotName)
	if err != nil {
		return fmt.Errorf("failed to stop connections: %w", err)
	}

	err = s.postgresRepository.CreateDatabase(snapshot.SnapshotDb, snapshot.ReferenceDb)
	if err != nil {
		return fmt.Errorf("failed to restore snapshot: %w", err)
	}

	fmt.Printf("Snapshot '%s' restored successfully into database '%s'\n", snapshotName, snapshot.SnapshotDb)
	return nil
}

func (s *SnapshotService) RemoveSnapshot(snapshotName string) error {
	err := s.promptToStopConnectionsIfNeeded(snapshotName)
	if err != nil {
		return fmt.Errorf("failed to stop connections: %w", err)
	}

	err = s.postgresRepository.DropDatabase(snapshotName)
	if err != nil {
		return fmt.Errorf("failed to drop snapshot: %w", err)
	}

	err = s.snappyRepository.DeleteSnapshot(snapshotName)
	if err != nil {
		return fmt.Errorf("failed to delete snapshot metadata: %w", err)
	}

	fmt.Printf("Snapshot '%s' removed successfully\n", snapshotName)

	return nil
}

func (s *SnapshotService) RenameSnapshot(oldName, newName string) error {
	err := s.promptToStopConnectionsIfNeeded(oldName)
	if err != nil {
		return fmt.Errorf("failed to stop connections: %w", err)
	}

	err = s.snappyRepository.RenameSnapshot(oldName, newName)
	if err != nil {
		return fmt.Errorf("failed to rename snapshot metadata: %w", err)
	}

	fmt.Printf("Snapshot '%s' renamed to '%s' with success\n", oldName, newName)
	return nil
}

func (s *SnapshotService) promptToStopConnectionsIfNeeded(dbName string) error {
	isUsed, err := s.postgresRepository.IsDatabaseInUse(dbName)
	if err != nil {
		return fmt.Errorf("failed to check if db is used: %w", err)
	}
	if !isUsed {
		return nil
	}

	fmt.Printf("Database '%s' is in use, do you wish to terminate all connections? [y/N] ", dbName)
	var response string
	fmt.Scanln(&response)

	if response != "y" && response != "Y" {
		return nil
	}

	fmt.Println("Stopping all connections...")
	err = s.postgresRepository.TerminateAllConnections(dbName)
	if err != nil {
		return fmt.Errorf("failed to terminate all connections: %w", err)
	}
	fmt.Println("All connections stopped successfully")

	return nil
}

func generateRandomSnapshotName() string {
	return fmt.Sprintf("snappy_%s", ulid.Make())
}

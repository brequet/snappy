package entity

import "time"

type Snapshot struct {
	ID          int64
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Name        string
	ReferenceDb string
	SnapshotDb  string
}

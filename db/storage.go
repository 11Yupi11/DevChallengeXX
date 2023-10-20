//go:generate mockgen -source=./storage.go -destination=./mock/storage_mock.go

package db

import (
	"context"
	"database/sql"
	"fmt"

	"dev-challenge/internal/models"
)

type Storage interface {
	GetCellInput(ctx context.Context, sheetID, cellID string) (*models.Data, error)
	AddCellInput(ctx context.Context, tx *sql.Tx, data Input) (resp *models.Data, wasUpdated bool, err error)
	GetSheetInput(ctx context.Context, sheetID string) (map[string]models.Data, error)
	GetCellInputBatch(ctx context.Context, tx *sql.Tx, sheetID string, cells []string) (map[string]string, error)
	GetIDList(ctx context.Context, tx *sql.Tx, cellID string) ([]int, error)
	GetInputBatchByIDs(ctx context.Context, tx *sql.Tx, IDs []int) (*[]Input, error)
	BeginTransaction(ctx context.Context) (*sql.Tx, error)
}

func NewStorage(ext *sql.DB) Storage {
	return &storage{
		ext: ext,
	}
}

type storage struct {
	ext *sql.DB
}

func (s *storage) BeginTransaction(ctx context.Context) (*sql.Tx, error) {
	// Start a new transaction
	tx, err := s.ext.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	return tx, nil
}

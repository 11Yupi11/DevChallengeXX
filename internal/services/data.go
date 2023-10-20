//go:generate mockgen -source=./data.go -destination=./mock/data_mock.go

package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"dev-challenge/db"
	"dev-challenge/internal/models"
)

type ExcelLikeService interface {
	GetCellInput(ctx context.Context, sheetID, cellID string) (*models.Data, error)
	AddCellInputTX(ctx context.Context, sheetID, cellID string, inputData *models.Data) (*models.Data, error)
	AddCellInput(ctx context.Context, tx *sql.Tx, sheetID, cellID string, inputData *models.Data) (*models.Data, error)
	GetSheetInput(ctx context.Context, sheetID string) (map[string]models.Data, error)
}

type excelLikeService struct {
	storage db.Storage
}

func NewExcelLikeService(storage db.Storage) ExcelLikeService {
	return &excelLikeService{
		storage: storage,
	}
}

func (s *excelLikeService) GetCellInput(ctx context.Context, sheetID, cellID string) (*models.Data, error) {
	return s.storage.GetCellInput(ctx, sheetID, cellID)
}

func (s *excelLikeService) AddCellInputTX(ctx context.Context, sheetID, cellID string, inputData *models.Data) (*models.Data, error) {
	// Start a transaction
	tx, err := s.storage.BeginTransaction(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			tx.Rollback() // Rollback if there's an error
		}
	}()

	resp, err := s.AddCellInput(ctx, tx, sheetID, cellID, inputData)
	if err != nil {
		return nil, err
	}

	// Commit the transaction
	if txErr := tx.Commit(); txErr != nil {
		return nil, txErr
	}

	return resp, nil
}

func (s *excelLikeService) AddCellInput(ctx context.Context, tx *sql.Tx, sheetID, cellID string, inputData *models.Data) (*models.Data, error) {
	if !isValid(inputData.Value) {
		return nil, errors.New("input value is not correct")
	}
	value := strings.ToLower(inputData.Value)

	var m map[string]string
	cellsToGet := extractParams(value)
	if len(cellsToGet) > 0 {
		if contains(cellsToGet, cellID) {
			return nil, errors.New("cell can't link to itself")
		}
		paramsToValues, err := s.storage.GetCellInputBatch(ctx, tx, sheetID, cellsToGet)
		if err != nil {
			return nil, err
		}
		m = paramsToValues
	}

	var newExpression string
	if m != nil {
		newExpression = replaceParamsWithValues(value, m)
	} else {
		newExpression = value
	}

	result, err := calculate(newExpression)
	if err != nil {
		return nil, err
	}

	input := db.Input{
		SheetID:    sheetID,
		CellID:     cellID,
		Value:      value,
		Result:     result,
		UsedParams: cellsToGet,
	}

	resp, wasUpdated, err := s.storage.AddCellInput(ctx, tx, input)
	if err != nil {
		return nil, err
	}

	if wasUpdated {
		if dependentErr := s.updateDependentCells(ctx, tx, input.CellID); dependentErr != nil {
			return nil, dependentErr
		}
	}

	return resp, nil
}

func (s *excelLikeService) GetSheetInput(ctx context.Context, sheetID string) (map[string]models.Data, error) {
	return s.storage.GetSheetInput(ctx, sheetID)
}

func (s *excelLikeService) updateDependentCells(ctx context.Context, tx *sql.Tx, cellID string) error {
	// 1) select distinct ID
	// 2) select all inputs by ID
	// 3) in cyclo for all inputs start AddCellInput

	needToBeChanged, err := s.storage.GetIDList(ctx, tx, cellID)
	if err != nil {
		return err
	}

	allInputs, err := s.storage.GetInputBatchByIDs(ctx, tx, needToBeChanged)
	if err != nil {
		return err
	}

	for _, input := range *allInputs {
		_, err := s.AddCellInput(ctx, tx, input.SheetID, input.CellID, &models.Data{
			Value:  input.Value,
			Result: fmt.Sprintf("%f", input.Result),
		})
		if err != nil {
			return err
		}
	}
	return nil
}

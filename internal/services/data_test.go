package services

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"dev-challenge/db"
	mock_db "dev-challenge/db/mock"
	"dev-challenge/internal/models"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestExcelLikeService_GetCellInput(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	storage := mock_db.NewMockStorage(ctrl)

	s := &excelLikeService{
		storage: storage,
	}

	tests := []struct {
		name         string
		sheetID      string
		cellID       string
		mockBehavior func()
		expectedData *models.Data
		expectedErr  error
	}{
		{
			name:    "Success",
			sheetID: "sheet1",
			cellID:  "cellA1",
			mockBehavior: func() {
				storage.EXPECT().GetCellInput(gomock.Any(), "sheet1", "cellA1").Return(&models.Data{
					Value:  "SampleValue",
					Result: "SomeResult",
				}, nil)
			},
			expectedData: &models.Data{
				Value:  "SampleValue",
				Result: "SomeResult",
			},
			expectedErr: nil,
		},
		{
			name:    "Storage Error",
			sheetID: "sheet1",
			cellID:  "cellA2",
			mockBehavior: func() {
				storage.EXPECT().GetCellInput(gomock.Any(), "sheet1", "cellA2").Return(nil, errors.New("Some storage error"))
			},
			expectedData: nil,
			expectedErr:  errors.New("Some storage error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior()

			data, err := s.GetCellInput(context.TODO(), tt.sheetID, tt.cellID)

			assert.Equal(t, tt.expectedData, data)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestExcelLikeService_AddCellInput(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	storage := mock_db.NewMockStorage(ctrl)
	tx := &sql.Tx{}

	s := &excelLikeService{
		storage: storage,
	}

	tests := []struct {
		name          string
		sheetID       string
		cellID        string
		inputData     *models.Data
		mockBehavior  func()
		expectedData  *models.Data
		expectedError string
	}{
		{
			name:    "Invalid Input Value",
			sheetID: "sheet1",
			cellID:  "cellA1",
			inputData: &models.Data{
				Value: "InvalidValue",
			},
			mockBehavior:  func() {},
			expectedData:  nil,
			expectedError: "input value is not correct",
		},
		{
			name:    "Self-referencing cell",
			sheetID: "sheet1",
			cellID:  "cella1",
			inputData: &models.Data{
				Value: "=cella1+5",
			},
			mockBehavior:  func() {},
			expectedData:  nil,
			expectedError: "cell can't link to itself",
		},
		{
			name:    "Storage error when fetching cells",
			sheetID: "sheet1",
			cellID:  "cellA1",
			inputData: &models.Data{
				Value: "=cellB1+5",
			},
			mockBehavior: func() {
				storage.EXPECT().GetCellInputBatch(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("some DB error"))
			},
			expectedData:  nil,
			expectedError: "some DB error",
		},
		{
			name:    "Calculate error division by 0",
			sheetID: "sheet1",
			cellID:  "cellA1",
			inputData: &models.Data{
				Value: "=5/0",
			},
			mockBehavior:  func() {},
			expectedData:  nil,
			expectedError: "division by zero",
		},
		{
			name:    "Calculate error unsupported operation",
			sheetID: "sheet1",
			cellID:  "cellA1",
			inputData: &models.Data{
				Value: "=5^2",
			},
			mockBehavior:  func() {},
			expectedData:  nil,
			expectedError: "unsupported binary operator: ^",
		},
		{
			name:    "AddCellInput storage error",
			sheetID: "sheet2",
			cellID:  "cellA2",
			inputData: &models.Data{
				Value: "5",
			},
			mockBehavior: func() {
				storage.EXPECT().AddCellInput(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, false, errors.New("Insert error"))
			},
			expectedData:  nil,
			expectedError: "Insert error",
		},
		{
			name:    "Update dependent cells error (failed to get ids)",
			sheetID: "sheet1",
			cellID:  "cellA1",
			inputData: &models.Data{
				Value: "5",
			},
			mockBehavior: func() {
				storage.EXPECT().AddCellInput(gomock.Any(), gomock.Any(), gomock.Any()).Return(&models.Data{}, true, nil) // Indicate wasUpdated=true
				storage.EXPECT().GetIDList(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("failed to get dependent inputs"))
			},
			expectedData:  nil,
			expectedError: "failed to get dependent inputs",
		},
		{
			name:    "Update dependent cells error (failed to input batch)",
			sheetID: "sheet1",
			cellID:  "cellA1",
			inputData: &models.Data{
				Value: "5",
			},
			mockBehavior: func() {
				storage.EXPECT().AddCellInput(gomock.Any(), gomock.Any(), gomock.Any()).Return(&models.Data{}, true, nil)
				storage.EXPECT().GetIDList(gomock.Any(), gomock.Any(), gomock.Any()).Return([]int{1, 2}, nil)
				storage.EXPECT().GetInputBatchByIDs(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("failed to get inputs batch"))
			},
			expectedData:  nil,
			expectedError: "failed to get inputs batch",
		},
		{
			name:    "Update dependent cells error (failed to update dependent)",
			sheetID: "sheet1",
			cellID:  "cellA1",
			inputData: &models.Data{
				Value: "5",
			},
			mockBehavior: func() {
				input := []db.Input{
					{
						SheetID:    "1",
						CellID:     "2",
						Value:      "3",
						Result:     3,
						UsedParams: nil,
					},
				}
				storage.EXPECT().AddCellInput(gomock.Any(), gomock.Any(), gomock.Any()).Return(&models.Data{}, true, nil)
				storage.EXPECT().GetIDList(gomock.Any(), gomock.Any(), gomock.Any()).Return([]int{1, 2}, nil)
				storage.EXPECT().GetInputBatchByIDs(gomock.Any(), gomock.Any(), gomock.Any()).Return(&input, nil)
				storage.EXPECT().AddCellInput(gomock.Any(), gomock.Any(), gomock.Any()).Return(&models.Data{}, false, errors.New("failed to update"))
			},
			expectedData:  nil,
			expectedError: "failed to update",
		},
		{
			name:    "Success scenario",
			sheetID: "sheet4",
			cellID:  "cellA4",
			inputData: &models.Data{
				Value: "=5",
			},
			mockBehavior: func() {
				input := []db.Input{
					{
						SheetID:    "1",
						CellID:     "2",
						Value:      "3",
						Result:     3,
						UsedParams: nil,
					},
				}
				storage.EXPECT().AddCellInput(gomock.Any(), gomock.Any(), gomock.Any()).Return(&models.Data{
					Value:  "=cellB1",
					Result: "10",
				}, true, nil)
				storage.EXPECT().GetIDList(gomock.Any(), gomock.Any(), gomock.Any()).Return([]int{1, 2}, nil)
				storage.EXPECT().GetInputBatchByIDs(gomock.Any(), gomock.Any(), gomock.Any()).Return(&input, nil)
				storage.EXPECT().AddCellInput(gomock.Any(), gomock.Any(), gomock.Any()).Return(&models.Data{}, false, nil)
			},
			expectedData:  &models.Data{Value: "=cellB1", Result: "10"},
			expectedError: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior()

			data, err := s.AddCellInput(context.TODO(), tx, tt.sheetID, tt.cellID, tt.inputData)

			if err != nil {
				assert.Contains(t, err.Error(), tt.expectedError)
			}
			assert.Equal(t, tt.expectedData, data)
		})
	}
}

func TestExcelLikeService_GetSheetInput(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	storage := mock_db.NewMockStorage(ctrl)

	s := &excelLikeService{
		storage: storage,
	}

	out := map[string]models.Data{"cell1": {
		Value:  "SampleValue",
		Result: "SomeResult",
	},
	}

	tests := []struct {
		name         string
		sheetID      string
		cellID       string
		mockBehavior func()
		expectedData map[string]models.Data
		expectedErr  error
	}{
		{
			name:    "Success",
			sheetID: "sheet1",
			cellID:  "cellA1",
			mockBehavior: func() {
				storage.EXPECT().GetSheetInput(gomock.Any(), "sheet1").Return(out, nil)
			},
			expectedData: out,
			expectedErr:  nil,
		},
		{
			name:    "Storage Error",
			sheetID: "sheet1",
			cellID:  "cellA2",
			mockBehavior: func() {
				storage.EXPECT().GetSheetInput(gomock.Any(), "sheet1").Return(nil, errors.New("some storage error"))
			},
			expectedData: nil,
			expectedErr:  errors.New("some storage error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior()

			data, err := s.GetSheetInput(context.TODO(), tt.sheetID)

			assert.Equal(t, tt.expectedData, data)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

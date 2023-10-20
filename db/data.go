package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"

	"dev-challenge/internal/models"
)

type Input struct {
	SheetID    string   `db:"sheet_id"`
	CellID     string   `db:"cell_id"`
	Value      string   `db:"value"`
	Result     float64  `db:"result"`
	UsedParams []string `db:"used_params"`
}

func (s *storage) GetCellInput(ctx context.Context, sheetID, cellID string) (resp *models.Data, err error) {
	rows, err := s.ext.QueryContext(ctx, "SELECT cell_value, cell_result FROM dev_challenge WHERE sheet_id=$1 AND cell_id=$2", sheetID, cellID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var value, result string
	for rows.Next() {
		err = rows.Scan(&value, &result)
		if err != nil {
			return nil, err
		}
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}
	resp = &models.Data{
		Value:  value,
		Result: result,
	}
	return resp, nil
}

func (s *storage) AddCellInput(ctx context.Context, tx *sql.Tx, data Input) (resp *models.Data, wasUpdated bool, err error) {
	var (
		maxIDBefore, maxIDAfter int
		wasItUpdate             bool
	)

	err = tx.QueryRowContext(ctx, "SELECT MAX(id) FROM dev_challenge").Scan(&maxIDBefore)
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.PrepareContext(ctx, "INSERT INTO dev_challenge(sheet_id, cell_id, cell_value, cell_result) VALUES($1,$2,$3,$4) "+
		"ON CONFLICT(sheet_id, cell_id) DO UPDATE SET cell_value = EXCLUDED.cell_value, cell_result = EXCLUDED.cell_result")
	if err != nil {
		return nil, wasItUpdate, err
	}
	defer stmt.Close()
	_, err = stmt.Exec(data.SheetID, data.CellID, data.Value, data.Result)
	if err != nil {
		return nil, wasItUpdate, err
	}

	err = tx.QueryRowContext(ctx, "SELECT MAX(id) FROM dev_challenge").Scan(&maxIDAfter)
	if err != nil {
		return nil, wasItUpdate, err
	}

	if maxIDAfter > maxIDBefore {
		for _, strValue := range data.UsedParams {
			_, err = tx.ExecContext(ctx, "INSERT INTO string_array(dev_challenge_id, string_value) VALUES($1,$2)", maxIDAfter, strValue)
			if err != nil {
				return nil, wasItUpdate, err
			}
		}
	} else {
		wasItUpdate = true
		// Delete the old strings
		_, err = tx.ExecContext(ctx, "DELETE FROM string_array WHERE dev_challenge_id = (SELECT id FROM dev_challenge WHERE sheet_id = $1 AND cell_id = $2)", data.SheetID, data.CellID)
		if err != nil {
			return nil, wasItUpdate, err
		}

		// Then insert the new strings
		for _, strValue := range data.UsedParams {
			_, err = tx.ExecContext(ctx, "INSERT INTO string_array(dev_challenge_id, string_value) VALUES((SELECT id FROM dev_challenge WHERE sheet_id = $1 AND cell_id = $2),$3)", data.SheetID, data.CellID, strValue)
			if err != nil {
				return nil, wasItUpdate, err
			}
		}
	}

	resp = &models.Data{Value: data.Value, Result: fmt.Sprintf("%f", data.Result)}
	return resp, wasItUpdate, nil
}

func (s *storage) GetSheetInput(ctx context.Context, sheetID string) (map[string]models.Data, error) {
	rows, err := s.ext.QueryContext(ctx, "SELECT  cell_id, cell_value, cell_result FROM dev_challenge WHERE sheet_id = $1", sheetID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := make(map[string]models.Data)
	for rows.Next() {
		var data Input
		if err := rows.Scan(&data.CellID, &data.Value, &data.Result); err != nil {
			return nil, err
		}
		res[data.CellID] = models.Data{
			Value:  data.Value,
			Result: fmt.Sprintf("%f", data.Result),
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(res) < 1 {
		return nil, errors.New("sheet not found")
	}

	return res, nil
}

func (s *storage) GetCellInputBatch(ctx context.Context, tx *sql.Tx, sheetID string, cells []string) (map[string]string, error) {
	placeholders := make([]string, len(cells))
	for i := range cells {
		placeholders[i] = fmt.Sprintf("$%d", i+2) // starting from $2 because $1 is used for sheetID
	}

	query := fmt.Sprintf("SELECT cell_id, cell_result FROM dev_challenge WHERE sheet_id = $1 AND cell_id IN (%s)", strings.Join(placeholders, ", "))

	args := make([]interface{}, len(cells)+1)
	args[0] = sheetID
	for i, cell := range cells {
		args[i+1] = cell
	}

	rows, err := tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var datas []Input
	for rows.Next() {
		var data Input
		if err := rows.Scan(&data.CellID, &data.Result); err != nil {
			return nil, err
		}
		datas = append(datas, data)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	resp := make(map[string]string)
	for _, val := range datas {
		resp[val.CellID] = fmt.Sprintf("%f", val.Result)
	}
	return resp, nil
}

func (s *storage) GetIDList(ctx context.Context, tx *sql.Tx, cellID string) ([]int, error) {
	query := "SELECT DISTINCT (dev_challenge_id) FROM string_array WHERE string_value=$1"

	rows, err := tx.QueryContext(ctx, query, cellID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var IDs []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		IDs = append(IDs, id)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return IDs, nil
}

func (s *storage) GetInputBatchByIDs(ctx context.Context, tx *sql.Tx, ids []int) (*[]Input, error) {
	placeholders := make([]string, len(ids))
	for i := range ids {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
	}

	query := fmt.Sprintf("SELECT sheet_id, cell_id, cell_value, cell_result FROM dev_challenge WHERE id IN (%s)", strings.Join(placeholders, ", "))

	args := make([]interface{}, len(ids))
	for i, id := range ids {
		args[i] = id
	}

	rows, err := tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var datas []Input
	for rows.Next() {
		var data Input
		if err := rows.Scan(&data.SheetID, &data.CellID, &data.Value, &data.Result); err != nil {
			return nil, err
		}
		datas = append(datas, data)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &datas, nil
}

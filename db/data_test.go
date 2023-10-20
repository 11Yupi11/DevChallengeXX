package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStorage_GetCellInput(t *testing.T) {
	defer cleanup()

	store := NewStorage(conn)
	data, err := store.GetCellInput(context.TODO(), "sheet0", "cell0")
	require.NoError(t, err)
	require.Equal(t, "0", data.Value)
	require.Equal(t, "0", data.Result)
}

func TestStorage_AddCellInput(t *testing.T) {
	defer cleanup()

	store := NewStorage(conn)
	tx, err := store.BeginTransaction(context.TODO())
	require.NoError(t, err)

	// new raw
	data, flag, err := store.AddCellInput(context.TODO(), tx, Input{
		SheetID:    "sheet1",
		CellID:     "cell2",
		Value:      "2",
		Result:     2,
		UsedParams: nil,
	})
	require.NoError(t, err)
	require.False(t, flag)
	require.Equal(t, data.Value, "2")
	require.Equal(t, data.Result, "2.000000")

	// updating existing
	data, flag, err = store.AddCellInput(context.TODO(), tx, Input{
		SheetID:    "sheet1",
		CellID:     "cell2",
		Value:      "3",
		Result:     3,
		UsedParams: nil,
	})
	require.NoError(t, err)
	require.True(t, flag)
	require.Equal(t, data.Value, "3")
	require.Equal(t, data.Result, "3.000000")

	// add data to 2 tables
	data, flag, err = store.AddCellInput(context.TODO(), tx, Input{
		SheetID:    "sheet1",
		CellID:     "cell3",
		Value:      "=cell1+cell2",
		Result:     4,
		UsedParams: []string{"cell1", "cell2"},
	})
	require.NoError(t, err)
	require.False(t, flag)
	require.Equal(t, data.Value, "=cell1+cell2")
	require.Equal(t, data.Result, "4.000000")

	// update data in 2 tables
	data, flag, err = store.AddCellInput(context.TODO(), tx, Input{
		SheetID:    "sheet1",
		CellID:     "cell3",
		Value:      "=cell1+cell2+45",
		Result:     4,
		UsedParams: []string{"cell1", "cell2"},
	})
	require.NoError(t, err)
	require.True(t, flag)
	require.Equal(t, data.Value, "=cell1+cell2+45")
	require.Equal(t, data.Result, "4.000000")

	err = tx.Commit()
	require.NoError(t, err)
}

func TestStorage_GetSheetInput(t *testing.T) {
	defer cleanup()

	store := NewStorage(conn)
	tx, err := store.BeginTransaction(context.TODO())
	require.NoError(t, err)

	// new 3 raws
	_, flag, err := store.AddCellInput(context.TODO(), tx, Input{
		SheetID:    "sheet1",
		CellID:     "cell1",
		Value:      "1",
		Result:     1,
		UsedParams: nil,
	})
	require.NoError(t, err)
	require.False(t, flag)
	_, flag, err = store.AddCellInput(context.TODO(), tx, Input{
		SheetID:    "sheet1",
		CellID:     "cell2",
		Value:      "2",
		Result:     2,
		UsedParams: nil,
	})
	require.NoError(t, err)
	require.False(t, flag)
	_, flag, err = store.AddCellInput(context.TODO(), tx, Input{
		SheetID:    "sheet1",
		CellID:     "cell3",
		Value:      "3",
		Result:     3,
		UsedParams: nil,
	})
	require.NoError(t, err)
	require.False(t, flag)

	err = tx.Commit()
	require.NoError(t, err)

	res, err := store.GetSheetInput(context.TODO(), "sheet1")
	require.NoError(t, err)
	require.Equal(t, 3, len(res))
}

func TestStorage_GetCellInputBatch(t *testing.T) {
	defer cleanup()

	store := NewStorage(conn)
	tx, err := store.BeginTransaction(context.TODO())
	require.NoError(t, err)

	// new 3 raws
	_, flag, err := store.AddCellInput(context.TODO(), tx, Input{
		SheetID:    "sheet1",
		CellID:     "cell1",
		Value:      "1",
		Result:     1,
		UsedParams: nil,
	})
	require.NoError(t, err)
	require.False(t, flag)
	_, flag, err = store.AddCellInput(context.TODO(), tx, Input{
		SheetID:    "sheet1",
		CellID:     "cell2",
		Value:      "2",
		Result:     2,
		UsedParams: nil,
	})
	require.NoError(t, err)
	require.False(t, flag)
	_, flag, err = store.AddCellInput(context.TODO(), tx, Input{
		SheetID:    "sheet1",
		CellID:     "cell3",
		Value:      "3",
		Result:     3,
		UsedParams: nil,
	})
	require.NoError(t, err)
	require.False(t, flag)

	err = tx.Commit()
	require.NoError(t, err)

	tx2, err := store.BeginTransaction(context.TODO())
	require.NoError(t, err)

	res, err := store.GetCellInputBatch(context.TODO(), tx2, "sheet1", []string{"cell1", "cell2", "cell3"})
	require.NoError(t, err)
	err = tx2.Commit()
	require.NoError(t, err)

	require.Equal(t, 3, len(res))
	require.Equal(t, "1.000000", res["cell1"])
	require.Equal(t, "2.000000", res["cell2"])
	require.Equal(t, "3.000000", res["cell3"])
}

func TestStorage_GetIDList(t *testing.T) {
	defer cleanup()

	store := NewStorage(conn)
	tx, err := store.BeginTransaction(context.TODO())
	require.NoError(t, err)

	// add raws with named params
	data, flag, err := store.AddCellInput(context.TODO(), tx, Input{
		SheetID:    "sheet1",
		CellID:     "cell1",
		Value:      "1",
		Result:     1,
		UsedParams: nil,
	})
	require.NoError(t, err)
	require.False(t, flag)
	require.Equal(t, data.Value, "1")
	require.Equal(t, data.Result, "1.000000")

	data, flag, err = store.AddCellInput(context.TODO(), tx, Input{
		SheetID:    "sheet1",
		CellID:     "cell2",
		Value:      "=cell1+2",
		Result:     3,
		UsedParams: []string{"cell1"},
	})
	require.NoError(t, err)
	require.False(t, flag)
	require.Equal(t, data.Value, "=cell1+2")
	require.Equal(t, data.Result, "3.000000")

	data, flag, err = store.AddCellInput(context.TODO(), tx, Input{
		SheetID:    "sheet1",
		CellID:     "cell3",
		Value:      "=cell1+cell2",
		Result:     4,
		UsedParams: []string{"cell1", "cell2"},
	})
	require.NoError(t, err)
	require.False(t, flag)
	require.Equal(t, data.Value, "=cell1+cell2")
	require.Equal(t, data.Result, "4.000000")

	err = tx.Commit()
	require.NoError(t, err)

	tx2, err := store.BeginTransaction(context.TODO())
	require.NoError(t, err)

	res, err := store.GetIDList(context.TODO(), tx2, "cell1")
	require.NoError(t, err)

	err = tx2.Commit()
	require.NoError(t, err)
	require.Equal(t, 2, len(res))
}

func TestStorage_GetInputBatchByIDs(t *testing.T) {
	defer cleanup()

	store := NewStorage(conn)
	tx, err := store.BeginTransaction(context.TODO())
	require.NoError(t, err)

	// add raws with named params
	data, flag, err := store.AddCellInput(context.TODO(), tx, Input{
		SheetID:    "sheet1",
		CellID:     "cell1",
		Value:      "1",
		Result:     1,
		UsedParams: nil,
	})
	require.NoError(t, err)
	require.False(t, flag)
	require.Equal(t, data.Value, "1")
	require.Equal(t, data.Result, "1.000000")

	data, flag, err = store.AddCellInput(context.TODO(), tx, Input{
		SheetID:    "sheet1",
		CellID:     "cell2",
		Value:      "=cell1+2",
		Result:     3,
		UsedParams: []string{"cell1"},
	})
	require.NoError(t, err)
	require.False(t, flag)
	require.Equal(t, data.Value, "=cell1+2")
	require.Equal(t, data.Result, "3.000000")

	data, flag, err = store.AddCellInput(context.TODO(), tx, Input{
		SheetID:    "sheet1",
		CellID:     "cell3",
		Value:      "=cell1+cell2",
		Result:     4,
		UsedParams: []string{"cell1", "cell2"},
	})
	require.NoError(t, err)
	require.False(t, flag)
	require.Equal(t, data.Value, "=cell1+cell2")
	require.Equal(t, data.Result, "4.000000")

	err = tx.Commit()
	require.NoError(t, err)

	tx2, err := store.BeginTransaction(context.TODO())
	require.NoError(t, err)

	res, err := store.GetInputBatchByIDs(context.TODO(), tx2, []int{2, 3}) // 1 is set during applying migration
	require.NoError(t, err)

	err = tx2.Commit()
	require.NoError(t, err)
	require.Equal(t, 2, len(*res))
	for id, val := range *res {
		require.Equal(t, "sheet1", val.SheetID)
		require.Equal(t, fmt.Sprintf("cell%d", id+1), val.CellID)
	}
}

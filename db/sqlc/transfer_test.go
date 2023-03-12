package db

import (
	"backendmaster/utils/random"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createRandomTransfer(t *testing.T) Transfer {
	dummyAccount1 := createRandomAccount(t)
	dummyAccount2 := createRandomAccount(t)

	var dummyTransferData CreateTransferParams = CreateTransferParams{
		FromAccountID: dummyAccount1.ID,
		ToAccountID:   dummyAccount2.ID,
		Amount:        random.RandomMoney(),
		Description:   random.RandomDescription(),
	}

	TransferFromDB, err := testQueries.CreateTransfer(context.Background(), dummyTransferData)
	require.NoError(t, err)
	require.NotEmpty(t, TransferFromDB)

	// cek pengirim
	require.Equal(t, dummyAccount1.ID, TransferFromDB.FromAccountID)

	// cek penerima
	require.Equal(t, dummyAccount2.ID, TransferFromDB.ToAccountID)

	// cek jumlah uang
	require.Equal(t, dummyTransferData.Amount, TransferFromDB.Amount)

	// cek deskripsi
	require.Equal(t, dummyTransferData.Description, TransferFromDB.Description)

	require.NotZero(t, TransferFromDB.ID)
	require.NotZero(t, TransferFromDB.CreatedAt)

	return TransferFromDB
}

func TestCreateTransfer(t *testing.T) {
	createRandomTransfer(t)
}

func TestGetTransfer(t *testing.T) {
	var createdTransfer Transfer = createRandomTransfer(t)

	transferFromDB, err := testQueries.GetTransfer(context.Background(), createdTransfer.ID)
	require.NoError(t, err)
	require.NotEmpty(t, transferFromDB)

	require.Equal(t, createdTransfer.ID, transferFromDB.ID)
	require.Equal(t, createdTransfer.FromAccountID, transferFromDB.FromAccountID)
	require.Equal(t, createdTransfer.ToAccountID, transferFromDB.ToAccountID)
	require.Equal(t, createdTransfer.Amount, transferFromDB.Amount)
	require.Equal(t, createdTransfer.Description, transferFromDB.Description)

	require.WithinDuration(t, createdTransfer.CreatedAt, transferFromDB.CreatedAt, time.Second)

}

func TestListTransfer(t *testing.T) {
	dummyAccount1 := createRandomAccount(t)
	dummyAccount2 := createRandomAccount(t)

	var dummyTransferData CreateTransferParams = CreateTransferParams{
		FromAccountID: dummyAccount1.ID,
		ToAccountID:   dummyAccount2.ID,
		Amount:        random.RandomMoney(),
		Description:   random.RandomDescription(),
	}

	transferFromDB, etErr := testQueries.CreateTransfer(context.Background(), dummyTransferData)
	require.NoError(t, etErr)
	require.NotEmpty(t, transferFromDB)

	for i := 0; i < 7; i++ {
		testQueries.CreateTransfer(context.Background(), dummyTransferData)
	}

	tfArg := ListTransferParams{
		FromAccountID: dummyAccount1.ID,
		ToAccountID:   dummyAccount2.ID,
		Limit:         5,
		Offset:        0,
	}

	listOfTransfer, err := testQueries.ListTransfer(context.Background(), tfArg)
	require.NoError(t, err)
	// t.Logf("data : %v", data)

	for _, transfer := range listOfTransfer {
		require.Equal(t, transferFromDB.FromAccountID, transfer.FromAccountID)
		require.Equal(t, transferFromDB.ToAccountID, transfer.ToAccountID)
	}

}

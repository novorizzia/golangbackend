package db

import (
	"backendmaster/utils/random"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createRandomEntry(t *testing.T) Entry {
	dummyAccount := createRandomAccount(t)

	var dummyEntryData CreateEntryParams = CreateEntryParams{
		AccountID:   dummyAccount.ID,
		Amount:      random.RandomMoney(),
		Description: random.RandomDescription(),
	}

	entryFromDB, err := testQueries.CreateEntry(context.Background(), dummyEntryData)
	require.NoError(t, err)
	require.NotEmpty(t, entryFromDB)

	require.Equal(t, dummyEntryData.AccountID, entryFromDB.AccountID)
	require.Equal(t, dummyEntryData.Amount, entryFromDB.Amount)
	require.Equal(t, dummyEntryData.Description, entryFromDB.Description)

	require.NotZero(t, entryFromDB.AccountID)
	require.NotZero(t, entryFromDB.CreatedAt)

	return entryFromDB
}

func TestCreateEntry(t *testing.T) {
	createRandomEntry(t)
}

func TestGetEntry(t *testing.T) {
	var createdEntry Entry = createRandomEntry(t)

	entryFromDB, err := testQueries.GetEntry(context.Background(), createdEntry.ID)

	require.NoError(t, err)
	require.NotEmpty(t, entryFromDB)

	require.Equal(t, createdEntry.ID, entryFromDB.ID)
	require.Equal(t, createdEntry.AccountID, entryFromDB.AccountID)
	require.Equal(t, createdEntry.Amount, entryFromDB.Amount)
	require.Equal(t, createdEntry.Description, entryFromDB.Description)

	require.WithinDuration(t, createdEntry.CreatedAt, entryFromDB.CreatedAt, time.Second)

}

func TestListEntries(t *testing.T) {
	dummyAccount := createRandomAccount(t)

	var dummyEntryData CreateEntryParams = CreateEntryParams{
		AccountID:   dummyAccount.ID,
		Amount:      random.RandomMoney(),
		Description: random.RandomDescription(),
	}

	entryFromDB, etErr := testQueries.CreateEntry(context.Background(), dummyEntryData)
	require.NoError(t, etErr)
	require.NotEmpty(t, entryFromDB)

	for i := 0; i < 7; i++ {
		testQueries.CreateEntry(context.Background(), dummyEntryData)
	}

	etArg := ListEntriesParams{
		AccountID: entryFromDB.AccountID,
		Limit:     5,
		Offset:    0,
	}

	listOfEntry, err := testQueries.ListEntries(context.Background(), etArg)
	require.NoError(t, err)
	// t.Logf("data : %v", data)

	for _, entry := range listOfEntry {
		require.Equal(t, entryFromDB.AccountID, entry.AccountID)
	}

}

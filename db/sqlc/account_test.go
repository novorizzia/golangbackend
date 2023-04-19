package db // package harus sama dengan yang ada pada account.sql.go

import (
	"backendmaster/utils/random"
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) Account {
	user := createRandomUser(t)
	var arg CreateAccountParams = CreateAccountParams{
		Owner:    user.Username,
		Balance:  random.RandomMoney(),
		Currency: random.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err) // check if error is nil, if its not nil then failed the tes
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner) // cek input dengan hasil yang diharapkan
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID) // cek account.id otomatis digen oleh posgres atau tidak
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)

}

func TestGetAccount(t *testing.T) {
	var createdAccount Account = createRandomAccount(t)
	accountFromPS, err := testQueries.GetAccount(context.Background(), createdAccount.ID)

	require.NoError(t, err)
	require.NotEmpty(t, accountFromPS)

	require.Equal(t, createdAccount.ID, accountFromPS.ID)
	require.Equal(t, createdAccount.Owner, accountFromPS.Owner)
	require.Equal(t, createdAccount.Balance, accountFromPS.Balance)
	require.Equal(t, createdAccount.Currency, accountFromPS.Currency)

	require.WithinDuration(t, createdAccount.CreatedAt, accountFromPS.CreatedAt, time.Second) // cek dua waktu apakah terpisah jauh atau tidak
}

func TestUpdateAccount(t *testing.T) {
	createdAccount := createRandomAccount(t)

	updatedData := UpdateAccountParams{
		ID:      createdAccount.ID,
		Balance: random.RandomMoney(),
	}

	updatedAccount, err := testQueries.UpdateAccount(context.Background(), updatedData)

	require.NoError(t, err)
	require.NotEmpty(t, updatedAccount)

	require.Equal(t, createdAccount.ID, updatedAccount.ID)
	require.Equal(t, createdAccount.Owner, updatedAccount.Owner)
	require.NotEqual(t, createdAccount.Balance, updatedAccount.Balance)
	require.Equal(t, updatedData.Balance, updatedAccount.Balance)
	require.Equal(t, createdAccount.Currency, updatedAccount.Currency)

}

func TestDeleteAccount(t *testing.T) {
	createdAccount := createRandomAccount(t)

	deletedAccount, err := testQueries.DeleteAccount(context.Background(), createdAccount.ID)
	require.NoError(t, err)
	require.NotEmpty(t, deletedAccount)

	require.Equal(t, createdAccount.ID, deletedAccount.ID)
	require.Equal(t, createdAccount.Owner, deletedAccount.Owner)
	require.Equal(t, createdAccount.Balance, deletedAccount.Balance)
	require.Equal(t, createdAccount.Currency, deletedAccount.Currency)

	accountFromPS, err := testQueries.GetAccount(context.Background(), createdAccount.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, accountFromPS)
}

func TestListAccount(t *testing.T) {
	for i := 0; i < 3; i++ {
		createRandomAccount(t)
	}

	arg := ListAccountParams{
		Limit:  3,
		Offset: 0,
	}

	listAccount, err := testQueries.ListAccount(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, listAccount, 3)

	for _, account := range listAccount {
		require.NotEmpty(t, account)
	}
}

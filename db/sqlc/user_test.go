package db

import (
	"backendmaster/utils/password"
	"backendmaster/utils/random"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	tmpOwner := random.RandomOwner()
	// mengubah password terlebih dahulu menjadi sebuah hash sebelum disimpan didalam database
	hashedPassword,err := password.HashedPassword(random.RandomString(6))
	require.NoError(t, err)

	var arg CreateUserParams = CreateUserParams{
		Username:       tmpOwner,
		HashedPassword: hashedPassword,
		FullName:       tmpOwner,
		Email:          random.RandomEmail(tmpOwner),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err) // check if error is nil, if its not nil then failed the tes
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username) // cek input dengan hasil yang diharapkan
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)

	require.True(t, user.PasswordChangedAt.IsZero()) // cek apakah field ini kosong apa tidak
	require.NotZero(t, user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)

}

func TestGetUser(t *testing.T) {
	var createdUser User = createRandomUser(t)
	userFromPS, err := testQueries.GetUser(context.Background(), createdUser.Username)

	require.NoError(t, err)
	require.NotEmpty(t, userFromPS)

	require.Equal(t, createdUser.Username, userFromPS.Username)
	require.Equal(t, createdUser.HashedPassword, userFromPS.HashedPassword)
	require.Equal(t, createdUser.FullName, userFromPS.FullName)
	require.Equal(t, createdUser.Email, userFromPS.Email)

	require.WithinDuration(t, createdUser.CreatedAt, userFromPS.CreatedAt, time.Second) // cek dua waktu apakah terpisah jauh atau tidak
	require.WithinDuration(t, createdUser.CreatedAt, userFromPS.CreatedAt, time.Second) // cek dua waktu apakah terpisah jauh atau tidak
}

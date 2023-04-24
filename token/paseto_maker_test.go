package token

import (
	"backendmaster/utils/random"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestPasetoMaker(t *testing.T) {
	maker, err := NewPasetoMaker(random.RandomString(32))
	require.NoError(t, err)

	username := random.RandomOwner()
	duration := time.Minute // one minute

	IssuedAt := time.Now()
	expired_at := IssuedAt.Add(duration)

	token, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, IssuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expired_at, payload.ExpiredAt, time.Second)
}

func TestExpiredPasetoToken(t *testing.T) {
	maker, err := NewPasetoMaker(random.RandomString(32))
	require.NoError(t, err)

	username := random.RandomOwner()

	token, err := maker.CreateToken(username, -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrorExpiredToken.Error())
	require.Nil(t, payload)

}

// KITA TIDAK MEMBUTUHKAN TEST NONE ALGORITHM KARENA KASUS SEPERTI ITU TIDAK ADA DI PASETO 

// TUGAS : MENULIS TEST INVALID TOKEN CASE
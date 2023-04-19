package password

import (
	"backendmaster/utils/random"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)


func TestPassword(t *testing.T) {
	password := random.RandomString(6)

	// get the hashpassword
	hashPassword1,err := HashedPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashPassword1)

	err = CheckPassword(password, hashPassword1)
	require.NoError(t, err) // no error mean password is correct


	wrongPassword := random.RandomString(6)
	err = CheckPassword(wrongPassword, hashPassword1)
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())

	// get the hashpassword
	hashPassword2,err := HashedPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashPassword2)

	require.NotEqual(t, hashPassword1,hashPassword2)
}
package util

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestMakePass(t *testing.T) {
	data := Randomstring(6)

	pass,err := HashPassword(data)
	require.NoError(t,err)
	require.NotEmpty(t,pass)

	err = CheckPass("data",pass)
	require.ErrorContains(t,err,bcrypt.ErrMismatchedHashAndPassword.Error())
	err =CheckPass(data,pass)
	require.NoError(t,err)
}

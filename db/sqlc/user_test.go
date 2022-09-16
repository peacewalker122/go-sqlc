package db

import (
	"context"
	"sqlc/util"
	"testing"

	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	pass,err := util.HashPassword(util.Randomstring(6))
	require.NoError(t,err)

	arg := CreateUserParams{
		Username:       util.Randomowner(),
		HashedPassword: pass,
		FullName:       util.Randomowner(),
		Email:          util.Randomemail(),
	}
	user, err := Testqueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)

	require.True(t, user.PasswordChangedAt.IsZero())
	require.NotZero(t, user.CreatedAt)

	return user
}
func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	account1 := createRandomUser(t)
	account2, err := Testqueries.GetUser(context.Background(), account1.Username)
	require.NoError(t, err)
	require.NotEmpty(t, account2)
	require.Equal(t, account1.Username, account2.Username)
	require.Equal(t, account1.FullName, account2.FullName)
	require.Equal(t, account1.Email, account2.Email)

	//err = util.CheckPass(account1.HashedPassword,account2.HashedPassword)
	require.Equal(t, account1.HashedPassword, account2.HashedPassword)
}

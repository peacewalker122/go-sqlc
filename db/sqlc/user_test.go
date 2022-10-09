package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/peacewalker122/go-sqlc/util"

	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	pass, err := util.HashPassword(util.Randomstring(6))
	require.NoError(t, err)

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

func TestUpdateUser(t *testing.T) {
	t.Run("UpdateFullname", func(t *testing.T) {
		oldUser := createRandomUser(t)

		Newfullname := util.Randomowner()
		newUser, err := Testqueries.UpdateUser(context.Background(), UpdateUserParams{
			Username: oldUser.Username,
			FullName: sql.NullString{
				String: Newfullname,
				Valid:  true,
			},
		})

		require.NoError(t, err)
		require.NotEqual(t, oldUser.FullName, newUser.FullName)
		require.Equal(t, Newfullname, newUser.FullName)
		require.Equal(t, oldUser.Username, newUser.Username)
		require.Equal(t, oldUser.Email, newUser.Email)
		require.Equal(t, oldUser.HashedPassword, newUser.HashedPassword)
	})

	t.Run("UpdatePassword", func(t *testing.T) {
		oldUser := createRandomUser(t)

		Newpassword := util.Randomstring(6)
		newUser, err := Testqueries.UpdateUser(context.Background(), UpdateUserParams{
			Username: oldUser.Username,
			HashedPassword: sql.NullString{
				String: Newpassword,
				Valid:  true,
			},
		})

		require.NoError(t, err)
		require.NotEqual(t, oldUser.HashedPassword, newUser.HashedPassword)
		require.Equal(t, Newpassword, newUser.HashedPassword)
		require.Equal(t, oldUser.Username, newUser.Username)
		require.Equal(t, oldUser.Email, newUser.Email)
		require.Equal(t, oldUser.FullName, newUser.FullName)
	})

	t.Run("UpdateEmail", func(t *testing.T) {
		oldUser := createRandomUser(t)

		Newemail := util.Randomemail()
		newUser, err := Testqueries.UpdateUser(context.Background(), UpdateUserParams{
			Username: oldUser.Username,
			Email: sql.NullString{
				String: Newemail,
				Valid:  true,
			},
		})

		require.NoError(t, err)
		require.NotEqual(t, oldUser.Email, newUser.Email)
		require.Equal(t, Newemail, newUser.Email)
		require.Equal(t, oldUser.Username, newUser.Username)
		require.Equal(t, oldUser.HashedPassword, newUser.HashedPassword)
		require.Equal(t, oldUser.FullName, newUser.FullName)
	})

	t.Run("AllField", func(t *testing.T) {
		oldUser := createRandomUser(t)

		Newfullname := util.Randomstring(6)
		Newemail := util.Randomemail()
		Newpassword := util.Randomstring(6)
		newUser, err := Testqueries.UpdateUser(context.Background(), UpdateUserParams{
			HashedPassword: sql.NullString{
				String: Newpassword,
				Valid:  true,
			},
			FullName: sql.NullString{
				String: Newfullname,
				Valid:  true,
			},
			Email: sql.NullString{
				String: Newemail,
				Valid:  true,
			},
			Username: oldUser.Username,
		})

		require.NoError(t, err)

		require.NotEqual(t, oldUser.Email, newUser.Email)
		require.Equal(t, Newemail, newUser.Email)

		require.NotEqual(t, oldUser.HashedPassword, newUser.HashedPassword)
		require.Equal(t, Newpassword, newUser.HashedPassword)

		require.NotEqual(t, oldUser.FullName, newUser.FullName)
		require.Equal(t, Newfullname, newUser.FullName)
	})
}

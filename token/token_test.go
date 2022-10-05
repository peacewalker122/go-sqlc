package token

import (
	"fmt"
	"testing"
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/peacewalker122/go-sqlc/util"

	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/require"
)

func TestJWT(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		maker, err := NewJWTmaker(util.Randomstring(32))
		require.NoError(t, err)
		require.NotEmpty(t, maker)

		username := util.Randomowner()
		duration := time.Minute

		IssuedAt := time.Now()
		ExpiredAt := IssuedAt.Add(duration)

		token, payload, err := maker.CreateToken(username, duration)
		require.NoError(t, err)
		require.NotEmpty(t, token)
		require.NotEmpty(t, payload)

		payload, err = maker.VerifyToken(token)
		require.NoError(t, err)
		require.NotEmpty(t, payload)

		require.NotZero(t, payload.ID)
		require.Equal(t, username, payload.Username)
		require.WithinDuration(t, IssuedAt, payload.IssuedAt, time.Second)
		require.WithinDuration(t, ExpiredAt, payload.ExpiredAt, time.Second)
	})
	t.Run("ExpiredJWT_Token", func(t *testing.T) {
		maker, err := NewJWTmaker(util.Randomstring(32))
		require.NoError(t, err)
		require.NotEmpty(t, maker)

		token, payload, err := maker.CreateToken(util.Randomowner(), -time.Minute)
		require.NoError(t, err)
		require.NotEmpty(t, token)
		require.NotEmpty(t, payload)

		payload, err = maker.VerifyToken(token)
		require.Error(t, err)
		require.EqualError(t, err, ErrExpired.Error())
		require.Nil(t, payload)
	})
	t.Run("InvalidKey", func(t *testing.T) {
		maker, err := NewJWTmaker(util.Randomstring(31))
		require.Error(t, err)
		require.EqualError(t, err, fmt.Errorf("invalid Key Size must be %v length", minSecretKeySize).Error())
		require.Nil(t, maker)
	})
	t.Run("InvalidJWT_TokenAlgNone", func(t *testing.T) {
		payload, err := Newpayload(util.Randomowner(), time.Minute)
		require.NoError(t, err)

		jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
		token, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
		require.NoError(t, err)

		maker, err := NewJWTmaker(util.Randomstring(32))
		require.NoError(t, err)

		payload, err = maker.VerifyToken(token)
		require.Error(t, err)
		require.EqualError(t, err, ErrToken.Error())
		require.Nil(t, payload)
	})
}

func TestPaseto(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		maker, err := NewPasetoMaker(util.Randomstring(32))
		require.NoError(t, err)
		require.NotEmpty(t, maker)

		username := util.Randomowner()
		duration := time.Minute

		IssuedAt := time.Now()
		ExpiredAt := IssuedAt.Add(duration)

		token, payload, err := maker.CreateToken(username, duration)
		require.NoError(t, err)
		require.NotEmpty(t, token)
		require.NotEmpty(t, payload)

		payload, err = maker.VerifyToken(token)
		require.NoError(t, err)
		require.NotEmpty(t, payload)

		require.NotZero(t, payload.ID)
		require.Equal(t, username, payload.Username)
		require.WithinDuration(t, IssuedAt, payload.IssuedAt, time.Second)
		require.WithinDuration(t, ExpiredAt, payload.ExpiredAt, time.Second)
	})
	t.Run("ExpiredToken", func(t *testing.T) {
		maker, err := NewJWTmaker(util.Randomstring(32))
		require.NoError(t, err)
		require.NotEmpty(t, maker)

		token, payload, err := maker.CreateToken(util.Randomowner(), -time.Minute)
		require.NoError(t, err)
		require.NotEmpty(t, token)
		require.NotEmpty(t, payload)

		payload, err = maker.VerifyToken(token)
		require.Error(t, err)
		require.EqualError(t, err, ErrExpired.Error())
		require.Nil(t, payload)
	})
	t.Run("Invalid Token", func(t *testing.T) {
		maker, err := NewPasetoMaker(util.Randomstring(1))
		require.Error(t, err)
		require.EqualError(t, err, fmt.Errorf("invalid key size, must be equal to %v characters", chacha20poly1305.KeySize).Error())
		require.Nil(t, maker)

		//not implemented yet.

	})
}

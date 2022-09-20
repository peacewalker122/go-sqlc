package token

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
)

const minSecretKeySize = 32

type JWTmaker struct {
	secretkey string
}

// CreateToken implements Maker

func NewJWTmaker(secretkey string) (Maker, error) {
	if len(secretkey) < minSecretKeySize {
		return nil, WrongKey
	}
	return &JWTmaker{secretkey}, nil
}

func (j *JWTmaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := Newpayload(username, duration)
	if err != nil {
		return "", err
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return jwtToken.SignedString([]byte(j.secretkey))
}

// VerifyToken implements Maker
func (j *JWTmaker) VerifyToken(token string) (*Payload, error) {
	Keyfunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrToken
		}
		return []byte(j.secretkey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, Keyfunc)
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, ErrExpired) {
			return nil, ErrExpired
		}
		return nil, ErrToken
	}
	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrToken
	}
	return payload, nil
}

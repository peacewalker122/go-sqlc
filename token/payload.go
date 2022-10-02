package token

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/aead/chacha20poly1305"
	uuid "github.com/google/uuid"
)

var (
	ErrToken   = errors.New("token invalid")
	ErrExpired = errors.New("token expired")
	WrongKey   = fmt.Errorf("invalid Key Size must be %v length", minSecretKeySize)
	Pasetoerr  = fmt.Errorf("invalid key size, must be equal to %v characters", chacha20poly1305.KeySize)
)

type Payload struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

func Newpayload(username string, duration time.Duration) (*Payload, error) {
	token, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	log.Println(duration)
	payload := &Payload{
		ID:        token,
		Username:  username,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}
	return payload, nil
}

func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return ErrExpired
	}
	return nil
}

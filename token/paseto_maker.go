package token

import (
	"fmt"
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/o1egl/paseto"
)

type PasetoMaker struct {
	paseto        *paseto.V2
	symmectricKey []byte
}

func NewPasetoMaker(symmectricKey string) (Maker, error) {
	if len(symmectricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid key size, must be equal to %v characters", chacha20poly1305.KeySize)
	}
	maker := &PasetoMaker{
		paseto:        &paseto.V2{},
		symmectricKey: []byte(symmectricKey),
	}
	return maker, nil
}

func (p *PasetoMaker) CreateToken(username string, duration time.Duration) (string, *Payload, error) {
	payload, err := Newpayload(username, duration)
	if err != nil {
		return "", payload, err
	}
	token, err := p.paseto.Encrypt(p.symmectricKey, payload, nil)
	return token, payload, err
}

func (p *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}

	err := p.paseto.Decrypt(token, p.symmectricKey, payload, nil)
	if err != nil {
		return nil, ErrToken
	}

	err = payload.Valid()
	if err != nil {
		return nil, err
	}

	return payload, nil
}

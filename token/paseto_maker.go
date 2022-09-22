package token

import (
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
		return nil, Pasetoerr
	}
	maker := &PasetoMaker{
		paseto:        &paseto.V2{},
		symmectricKey: []byte(symmectricKey),
	}
	return maker, nil
}

func (p *PasetoMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := Newpayload(username, duration)
	if err != nil {
		return "", err
	}
	return p.paseto.Encrypt(p.symmectricKey, payload, nil)
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
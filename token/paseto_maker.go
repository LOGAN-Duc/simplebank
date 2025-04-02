package token

import (
	"fmt"
	"github.com/o1egl/paseto"
	"golang.org/x/crypto/chacha20poly1305"
	"time"
)

// Interface thoong bao
type PaseMaker struct {
	paseto       *paseto.V2
	symmetricKey []byte
}

// tao phien ban moi PasetoMaker
// swr dung thuat toan Chacha Poly
func NewPasetoMaker(symmetricKey string) (Maker, error) {
	if len(symmetricKey) != chacha20poly1305.KeySize {

		return nil, fmt.Errorf("symmetric key must be 32 bytes long at least 32 bytes")
	}
	maker := &PaseMaker{
		paseto:       paseto.NewV2(),
		symmetricKey: []byte(symmetricKey),
	}
	return maker, nil
}

// create a new token for a specific username and duration
func (maker *PaseMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayLoad(username, duration)
	if err != nil {
		return "", err
	}
	return maker.paseto.Encrypt(maker.symmetricKey, payload, nil)
}

// Kiem tra ma dau vao co hop la khong
func (maker *PaseMaker) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}

	err := maker.paseto.Decrypt(token, maker.symmetricKey, payload, nil)
	if err != nil {
		return nil, ErrInvalidToken
	}
	err = payload.Valid()
	if err != nil {
		return nil, err
	}
	return payload, nil
}

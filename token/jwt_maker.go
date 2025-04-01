package token

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

// gioi han cho khoa
const minSecret = 32

// trinh tao token web json, su dung khoa doi xung
type JWTMaker struct {
	secretKet string
}

// trien khai interface ma thong bao
func NewJWTMaker(secretKet string) (Maker, error) {
	if len(secretKet) < minSecret {
		return nil, fmt.Errorf("invaild key size : must be at least %d characters", minSecret)
	}
	return &JWTMaker{secretKet}, nil
}

// create a new token for a specific username and duration
func (maker *JWTMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayLoad(username, duration)
	if err != nil {
		return "", err
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return jwtToken.SignedString([]byte(maker.secretKet))
}

// Kiem tra ma dau vao co hop la khong
func (maker *JWTMaker) VerifyToken(token string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		return []byte(maker.secretKet), nil
	}
	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		var verr *jwt.ValidationError
		ok := errors.As(err, &verr)
		if ok && errors.Is(verr.Inner, ErrExpiredToken) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}
	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrInvalidToken
	}
	return payload, nil

}

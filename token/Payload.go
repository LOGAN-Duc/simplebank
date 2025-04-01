package token

import (
	"errors"
	"github.com/google/uuid"
	"time"
)

var (
	ErrExpiredToken = errors.New("expired token")
	ErrInvalidToken = errors.New("invalid token")
)

// contains the payload data of the token
type Payload struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

// tao token moi voi ten nguoi dung va thoi han cu the
func NewPayLoad(username string, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	payload := &Payload{
		ID:        tokenID,
		Username:  username,
		IssuedAt:  time.Now(),
		ExpiresAt: time.Now().Add(duration),
	}
	return payload, nil
}

// Valid checks if the token payload is valid or not
func (payload *Payload) Valid() error {
	//kiemtra token het han
	if time.Now().After(payload.ExpiresAt) {
		return ErrExpiredToken
	}
	return nil
}

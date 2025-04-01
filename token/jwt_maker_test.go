package token

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/require"
	"github.com/techschool/simplebank/Utill"
	"testing"
	"time"
)

func TestJWTMaker(t *testing.T) {
	maker, err := NewJWTMaker(Utill.RandomString(32))
	require.NoError(t, err)
	username := Utill.RandomOwner()
	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.Equal(t, payload.Username, username)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiresAt, time.Second)
}

// Kiem tra het han token
func TestExpiredToken(t *testing.T) {
	maker, err := NewJWTMaker(Utill.RandomString(32))
	require.NoError(t, err)

	username := Utill.RandomOwner()
	duration := -time.Minute

	token, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Nil(t, payload)
}

func TestInvalidTokenAlgNone(t *testing.T) {
	username := Utill.RandomOwner()
	duration := time.Minute
	payload, err := NewPayLoad(username, duration)
	require.NoError(t, err)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
	token, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	maker, err := NewJWTMaker(Utill.RandomString(32))
	require.NoError(t, err)

	payload, err = maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrInvalidToken.Error())
	require.Nil(t, payload)
}

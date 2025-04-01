package token

import (
	"time"
)

// maker is an interface for managing tokens
type Maker interface {
	//create a new token for a specific username and duration
	CreateToken(username string, duration time.Duration) (string, error)
	//Kiem tra ma dau vao co hop la khong
	VerifyToken(token string) (*Payload, error)
}

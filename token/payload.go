package token

import (
	"errors"
	"time"

	"github.com/google/uuid"
)


// berbagai error yang dikembalikan oleh VerifyToken function
var (
	ErrorExpiredToken = errors.New("token sudah basi cuy")
	ErrorInvalidToken = errors.New("token tidak valid")
)

// Payload berisi data payload yang terkandung didalam token
type Payload struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

// NewPayload membuat payload baru untuk spesifik username dan durasi
func NewPayload(username string, duration time.Duration) (*Payload, error) {
	tokenId, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	payload := &Payload{
		ID:        tokenId,
		Username:  username,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),

	}
	return payload, nil
}

// mengecek apakah payload dari token valid atau tidak
func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return ErrorExpiredToken
	}
	return nil
}

package token

import "time"

type Maker interface {
	// CreateToken create new token untuk target spesifik username dan durasi yang
	CreateToken(username string, duration time.Duration) (string, error)

	// VerifyToken melakukan verifikasi pada token. jika valid akan mengirimkan payload yang ada dalam body dari token tersebut
	VerifyToken(token string) (*Payload, error)
}

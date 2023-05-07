package token

import (
	"fmt"
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/o1egl/paseto"
)

type PasetoMaker struct {
	paseto       *paseto.V2
	symmetricKey []byte
}

func NewPasetoMaker(symmetricKey string) (Maker, error) {
	// chacha20poly1305.KeySize ukurannya adalah 32 karakter
	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("Ukuran key tidak benar : ukuran key harus tepat %d karakter", chacha20poly1305.KeySize)
	}

	maker := &PasetoMaker{
		paseto:       paseto.NewV2(),
		symmetricKey: []byte(symmetricKey),
	}

	return maker, nil
}

func (maker *PasetoMaker) CreateToken(username string, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", payload, err
	}

	token, err := maker.paseto.Encrypt(maker.symmetricKey, payload, nil)
	return token, payload, err
}

// VerifyToken melakukan verifikasi pada token. jika valid akan mengirimkan payload yang ada dalam body dari token tersebut
func (maker *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{} // kita buat empty payload object untuk nantinya digunakan menyimpan decrypted data

	err := maker.paseto.Decrypt(token, maker.symmetricKey, payload, nil)
	if err != nil {
		return nil, ErrorInvalidToken
	}

	err = payload.Valid()
	if err != nil {
		return nil, err
	}

	return payload, nil
}

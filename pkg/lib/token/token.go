package token

import (
	"os"
)

type token struct {
	privateKey []byte
	publicKey  []byte
}

func Init() (TokenItf, error) {
	prvKey, err := os.ReadFile(os.Getenv("PRIVATE_KEY"))
	if err != nil {
		return nil, err
	}

	pubKey, err := os.ReadFile(os.Getenv("PUBLIC_KEY"))
	if err != nil {
		return nil, err
	}

	return token{
		privateKey: prvKey,
		publicKey:  pubKey,
	}, nil
}

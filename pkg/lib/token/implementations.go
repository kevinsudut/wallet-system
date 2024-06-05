package token

import (
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
)

func (t token) Create(ttl time.Duration, content interface{}) (string, error) {
	key, err := jwt.ParseRSAPrivateKeyFromPEM(t.privateKey)
	if err != nil {
		return "", err
	}

	now := time.Now().UTC()

	claims := make(jwt.MapClaims)
	claims["dat"] = content
	claims["exp"] = now.Add(ttl).Unix()
	claims["iat"] = now.Unix()
	claims["nbf"] = now.Unix()

	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(key)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (t token) Validate(token string) (interface{}, error) {
	key, err := jwt.ParseRSAPublicKeyFromPEM(t.publicKey)
	if err != nil {
		return "", err
	}

	if strings.HasPrefix(token, "Bearer ") {
		token = strings.Split(token, "Bearer ")[1]
	}

	tok, err := jwt.Parse(token, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected method: %s", jwtToken.Header["alg"])
		}

		return key, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := tok.Claims.(jwt.MapClaims)
	if !ok || !tok.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims["dat"], nil
}

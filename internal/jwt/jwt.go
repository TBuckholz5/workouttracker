package jwt

import (
	"crypto/rand"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJwt(userID int64) (string, error) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		panic(err)
	}
	claims := jwt.MapClaims{
		"exp": time.Now().Add(time.Hour * 24).Unix(),
		"sub": userID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	signedToken, err := token.SignedString(key)
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

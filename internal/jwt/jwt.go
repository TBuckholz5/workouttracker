package jwt

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const ISSUER = "workout-tracker"

type JwtService interface {
	GenerateJwt(userID int64) (string, error)
	ValidateJwt(ctx *gin.Context, tokenString string) error
}

type Jwt struct {
	secret []byte
}

func NewJwtService(jwtSecret []byte) *Jwt {
	return &Jwt{
		secret: jwtSecret,
	}
}

func (j *Jwt) GenerateJwt(userID int64) (string, error) {
	claims := jwt.MapClaims{
		"exp": time.Now().Add(time.Hour * 24).Unix(),
		"sub": userID,
		"iss": ISSUER,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(j.secret)
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func (j *Jwt) ValidateJwt(ctx *gin.Context, tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return j.secret, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		return err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return err
	}
	if claims["iss"] != ISSUER {
		return fmt.Errorf("invalid token issuer")
	}
	if time.Now().Unix() > int64(claims["exp"].(float64)) {
		return fmt.Errorf("token has expired")
	}
	ctx.Set("userID", int64(claims["sub"].(float64)))
	return nil
}

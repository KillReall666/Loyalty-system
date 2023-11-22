package authentication

import (
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type claims struct {
	UserID string `json:"userid"`
	jwt.RegisteredClaims
}

// TODO: засунуть ключ в переменную окружения
const (
	TOKEN_EXP  = time.Hour * 3
	SECRET_KEY = "supersecretkey"
)

func BuildJWTString(id string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TOKEN_EXP)),
		},
		UserID: id,
	})

	tokenString, err := token.SignedString([]byte(SECRET_KEY))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

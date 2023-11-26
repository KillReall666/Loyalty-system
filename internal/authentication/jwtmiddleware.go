package authentication

import (
	"context"
	"fmt"
	"github.com/KillReall666/Loyalty-system/internal/logger"
	"github.com/KillReall666/Loyalty-system/internal/storage/redis"
	"github.com/golang-jwt/jwt/v4"
	"github.com/golang-jwt/jwt/v4/request"
	"net/http"
)

type JWTMiddleware struct {
	RedisClient *redis.RedisClient
	Log         *logger.Logger
}

func (j *JWTMiddleware) JWTMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		jwtCheck := func(w http.ResponseWriter, r *http.Request) {
			//получаем jwt из заголовка
			extractor := request.AuthorizationHeaderExtractor
			extToken, err := extractor.ExtractToken(r)
			if err != nil {
				j.Log.LogWarning("err when extract token in jwt middleware: ", err)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			claim := &claims{}
			//Проверяем подлинность jwt
			_, err = jwt.ParseWithClaims(extToken, claim, func(t *jwt.Token) (interface{}, error) {
				//Проверка заголовка алгоритма токена. Заголовок должен совпадать с тем, который использует сервер для подписи и проверки токенов.
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					w.WriteHeader(http.StatusUnauthorized)
					return nil, fmt.Errorf("unexpected signing token method: %v", t.Header["alg"])
				}
				return []byte(SECRET_KEY), nil
			})
			if err != nil {
				j.Log.LogWarning("err when parse jwt: %v", err)
				fmt.Fprintf(w, "token lifetime has expired, log in")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			userID := claim.UserID

			_, err = j.RedisClient.Get(userID)
			if err != nil {
				j.Log.LogWarning("err when get token in middleware:", err)
				fmt.Fprintf(w, "token not valid")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), "UserID", userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		}

		return http.HandlerFunc(jwtCheck)
	}
}

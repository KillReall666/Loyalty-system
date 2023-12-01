package authorization

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/KillReall666/Loyalty-system/internal/authentication"
	"github.com/KillReall666/Loyalty-system/internal/credentials"
	"github.com/KillReall666/Loyalty-system/internal/logger"
	"github.com/KillReall666/Loyalty-system/internal/storage/redis"
	"net/http"
	"time"
)

type AuthHandler struct {
	checkUser   CredentialsChecker
	RedisClient *redis.RedisClient
	log         *logger.Logger
}

type CredentialsChecker interface {
	CredentialsGetter(ctx context.Context, user string) (string, string, error)
}

func NewAuthorizationHandler(ch CredentialsChecker, redis *redis.RedisClient, log *logger.Logger) *AuthHandler {
	return &AuthHandler{
		checkUser:   ch,
		RedisClient: redis,
		log:         log,
	}
}

// AuthenticationHandler TODO: Скорее всего это должен быть мидлтварь
func (a *AuthHandler) AuthorizationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "only POST requests support!", http.StatusNotFound)
		return
	}

	var buf bytes.Buffer
	var user credentials.User
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	hashPasswordFromDB, id, err := a.checkUser.CredentialsGetter(ctx, user.Username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	var token string

	if user.ComparePassword(hashPasswordFromDB, user.PasswordHash) {
		token, err = authentication.BuildJWTString(id)
		if err != nil {
			a.log.LogWarning("err when get JWT token when authorization", err)
			return
		}

		w.Header().Set("Authorization", token)
		err = a.RedisClient.Set(id, token)
		if err != nil {
			a.log.LogWarning("err when set value to redis in auth handler:", err)
			return
		}
		w.WriteHeader(http.StatusOK)
		//Надо ли?
		fmt.Fprintf(w, "You have successfully authorized")
		a.log.LogInfo("user", id, "successfully authorized")
	} else {
		err = errors.New("incorrect password, please try again")
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

}

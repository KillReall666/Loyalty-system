package authentication

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/KillReall666/Loyalty-system/internal/authentication"
	"github.com/KillReall666/Loyalty-system/internal/credentials"
	"net/http"
	"time"
)

type AuthHandler struct {
	checkUser CredentialsChecker
}

type CredentialsChecker interface {
	CheckCredentials(ctx context.Context, user string) (error, string, int)
}

func NewAuthenticationHandler(ch CredentialsChecker) *AuthHandler {
	return &AuthHandler{
		checkUser: ch,
	}
}

// AuthenticationHandler TODO: Скорее всего это должен быть мидлтварь
func (a *AuthHandler) AuthenticationHandler(w http.ResponseWriter, r *http.Request) {
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

	err, hashPasswordFromDb, id := a.checkUser.CheckCredentials(ctx, user.Username)
	if err != nil {
		fmt.Println(err)
	}

	var token string
	if user.ComparePassword(hashPasswordFromDb, user.PasswordHash) {
		fmt.Println("Correct password")
		token, err = authentication.BuildJWTString(id)
		w.Header().Set("Authorization", token)
		fmt.Fprintf(w, "You have successfully authorized")
	} else {
		err = errors.New("incorrect password, please try again")
		http.Error(w, err.Error(), http.StatusUnauthorized)
	}

}

package getbalance

import (
	"context"
	"encoding/json"
	"github.com/KillReall666/Loyalty-system/internal/dto"
	"github.com/KillReall666/Loyalty-system/internal/logger"
	"net/http"
)

type GetBalanceHandler struct {
	log           *logger.Logger
	BalanceGetter BalanceGetter
}

type BalanceGetter interface {
	GetUserBalance(ctx context.Context, userId string) (*dto.UserBalance, error)
}

func NewGetBalanceHandler(balance BalanceGetter, log *logger.Logger) *GetBalanceHandler {
	return &GetBalanceHandler{
		log:           log,
		BalanceGetter: balance,
	}
}

func (g *GetBalanceHandler) GetUserBalanceHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "only GET requests support!", http.StatusNotFound)
		return
	}
	userId := r.Context().Value("UserID").(string)
	balance, err := g.BalanceGetter.GetUserBalance(context.Background(), userId)
	if err != nil {
		g.log.LogWarning("err when getting user balance: ", err)
	}

	jsonData, err := json.Marshal(balance)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)

}

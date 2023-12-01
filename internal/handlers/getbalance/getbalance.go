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
	GetUserBalance(ctx context.Context, userID string) (*dto.UserBalance, error)
}

func NewGetBalanceHandler(balance BalanceGetter, log *logger.Logger) *GetBalanceHandler {
	return &GetBalanceHandler{
		log:           log,
		BalanceGetter: balance,
	}
}

// GetUserBalanceHandler TODO: что делать когда баланс отсутствует или нулевой? Пока вернул ошибку.
func (g *GetBalanceHandler) GetUserBalanceHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "only GET requests support!", http.StatusNotFound)
		return
	}
	userID := r.Context().Value("UserID").(string)
	balance, err := g.BalanceGetter.GetUserBalance(context.Background(), userID)
	if err != nil {
		g.log.LogWarning("err when getting user balance: ", err)
		http.Error(w, "zero balance or no information about charges", http.StatusPaymentRequired)
		return
	}

	jsonData, err := json.Marshal(balance)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)

}

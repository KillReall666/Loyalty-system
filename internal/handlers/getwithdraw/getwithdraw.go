package getwithdraw

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/KillReall666/Loyalty-system/internal/dto"
	"github.com/KillReall666/Loyalty-system/internal/logger"
	"net/http"
	"time"
)

type GetWithdrawHandler struct {
	GetWithDraw GetWithDrawer
	Log         *logger.Logger
}

type GetWithDrawer interface {
	GetWithdrawals(ctx context.Context, userId string) ([]*dto.Billing, error)
}

func NewGetWithdrawHandler(getwd GetWithDrawer, log *logger.Logger) *GetWithdrawHandler {
	return &GetWithdrawHandler{
		GetWithDraw: getwd,
		Log:         log,
	}
}

func (g *GetWithdrawHandler) GetWithdrawHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "only GET requests support!", http.StatusNotFound)
	}
	userId := r.Context().Value("UserID").(string)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	orders, err := g.GetWithDraw.GetWithdrawals(ctx, userId)
	if err != nil {
		g.Log.LogWarning("err when getting data from db:", err)
		fmt.Fprintf(w, "") //идет двойная запись в заголовок, в запросе к бд getwithdrawals errors.New()
		http.Error(w, err.Error(), http.StatusNoContent)
		return
	}

	jsonData, err := json.Marshal(orders)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)

}

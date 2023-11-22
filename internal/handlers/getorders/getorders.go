package getorders

import (
	"context"
	"encoding/json"
	"github.com/KillReall666/Loyalty-system/internal/dto"
	"github.com/KillReall666/Loyalty-system/internal/interrogator"
	"github.com/KillReall666/Loyalty-system/internal/logger"
	"net/http"
	"time"
)

type GetOrdersHandler struct {
	log          *logger.Logger
	OrdersGetter OrdersGetter
	interrogator *interrogator.Interrogator
}

type OrdersGetter interface {
	GetOrders(ctx context.Context, userId string) ([]dto.FullOrder, error)
}

func NewGetOrdersHandler(ord OrdersGetter, log *logger.Logger, interrogator *interrogator.Interrogator) *GetOrdersHandler {
	return &GetOrdersHandler{
		log:          log,
		OrdersGetter: ord,
		interrogator: interrogator,
	}
}

func (g *GetOrdersHandler) GetOrdersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "only GET requests support!", http.StatusNotFound)
	}
	userId := r.Context().Value("UserID").(string)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	orders, err := g.OrdersGetter.GetOrders(ctx, userId)
	if err != nil {
		g.log.LogWarning("err when getting orders list", err)
	}

	jsonData, err := json.Marshal(orders)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

package charge

import (
	"context"
	"encoding/json"
	"github.com/KillReall666/Loyalty-system/internal/dto"
	"github.com/KillReall666/Loyalty-system/internal/logger"
	"io"
	"net/http"
)

type ChargeHandler struct {
	Charge Charger
	Log    *logger.Logger
}

type Charger interface {
	ProcessOrder(ctx context.Context, order, userId string, sum float32) error
}

func NewChargeHandler(charge Charger, log *logger.Logger) *ChargeHandler {
	return &ChargeHandler{
		Charge: charge,
		Log:    log,
	}
}

func (c *ChargeHandler) ChargeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "only POST requests support!", http.StatusNotFound)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "couldn't read data from request body", http.StatusBadRequest)
		return
	}

	var orderData dto.WithdrawOrder
	err = json.Unmarshal(body, &orderData)
	if err != nil {
		http.Error(w, "Failed tro decode JSON data", http.StatusBadRequest)
		return
	}

	userId := r.Context().Value("UserID").(string)

	err = c.Charge.ProcessOrder(context.Background(), orderData.Order, userId, orderData.Sum)
	if err != nil {
		c.Log.LogWarning("in handler: ", err)
		return
	}

}

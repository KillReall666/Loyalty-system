package addorder

import (
	"context"
	"fmt"
	"github.com/KillReall666/Loyalty-system/internal/logger"
	"github.com/ShiraazMoollatjie/goluhn"
	"io"
	"net/http"

	"github.com/KillReall666/Loyalty-system/internal/storage/redis"
)

type AddOrderHandler struct {
	addOrder    AddOrder
	RedisClient *redis.RedisClient
	Log         *logger.Logger
}

type AddOrder interface {
	OrderSetter(ctx context.Context, userId, orderNumber string) error
}

func NewPutOrderNumberHandler(order AddOrder, redis *redis.RedisClient, log *logger.Logger) *AddOrderHandler {
	return &AddOrderHandler{
		addOrder:    order,
		RedisClient: redis,
		Log:         log,
	}
}

func (a *AddOrderHandler) AddOrderNumberHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "only POST requests support!", http.StatusNotFound)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "couldn't read data from request body", http.StatusBadRequest)
		return
	}

	orderNumber := string(body)
	err = goluhn.Validate(orderNumber)
	if err != nil {
		http.Error(w, "invalid format of order number", http.StatusUnprocessableEntity)
		return
	}

	userId := r.Context().Value("UserID").(string)

	err = a.addOrder.OrderSetter(context.Background(), userId, orderNumber)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		a.Log.LogWarning("error when add order handler: ", err)
		return
	} else {
		w.WriteHeader(http.StatusAccepted)
		fmt.Fprint(w, "order added")
	}

}

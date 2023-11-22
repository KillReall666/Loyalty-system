package putordernumber

import (
	"context"
	"fmt"
	"net/http"
)

type AddOrderHandler struct {
	addOrder AddOrder
}

type AddOrder interface {
	OrderSetter(ctx context.Context, userId, orderNumber string) error
}

func NewPutOrderNumberHandler(order AddOrder) *AddOrderHandler {
	return &AddOrderHandler{
		addOrder: order,
	}
}

func (a *AddOrderHandler) AddOrderNumberHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "only POST requests support!", http.StatusNotFound)
		return
	}
	order := w.Header().Get("text/plain")
	err := a.addOrder.OrderSetter(context.Background(), order, "1")
	if err != nil {
		fmt.Println("error in add order handler: ", err)
	}
}

package interrogator

import (
	"context"
	"encoding/json"
	"github.com/KillReall666/Loyalty-system/internal/dto"
	"github.com/KillReall666/Loyalty-system/internal/logger"
	"github.com/KillReall666/Loyalty-system/internal/storage/postgres"
	"io"
	"net/http"
	"time"
)

type Interrogator struct {
	db  *postgres.Database
	log *logger.Logger
}

func NewInterrogator(db *postgres.Database, log *logger.Logger) *Interrogator {
	return &Interrogator{
		db:  db,
		log: log,
	}
}

func (i *Interrogator) OrderStatusWorker() {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	orders, err := i.db.GetOrderNumbers(ctx)
	if err != nil {
		i.log.LogWarning("err when getting orders list from db: ", err)
	}

	for j := 0; j < len(orders); j++ {
		status, accrual, err := i.GetOrderStatusFromACCRUAL(orders[j])
		if err != nil {
			i.log.LogWarning("Error retrieving order status from ACCRUAL: ", err)
			return
		}

		switch status {
		case "PROCESSED":
			// Переместить заказ в базу данных с новым статусом (PROCESSED)
			userId := i.UpdateOrderStatusInDB(orders[j], "PROCESSED", accrual)
			err = i.db.IncrementCurrent(ctx, userId, accrual)
			if err != nil {
				i.log.LogWarning("err when add user balance: ", err)
			}
		case "INVALID":
			// Переместить заказ в базу данных с новым статусом (INVALID)
			userId := i.UpdateOrderStatusInDB(orders[j], "INVALID", accrual)
			err = i.db.IncrementCurrent(ctx, userId, accrual)
			if err != nil {
				i.log.LogWarning("err when add user balance: ", err)
			}
		default:
			// Пока заказ имеет статус отличный от PROCESSED и INVALID,
			i.log.LogInfo("Order is still being processed in ACCRUAL")
		}
	}
}

func (i *Interrogator) GetOrderStatusFromACCRUAL(orderNumber string) (string, float32, error) {
	req, err := http.NewRequest("GET", "http://localhost:8888/api/orders/"+orderNumber, nil)
	if err != nil {
		i.log.LogWarning("err when create GET request: ", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		i.log.LogWarning("err when make GET request: ", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		i.log.LogWarning("err when read response body: ", err)
	}
	var order dto.FullOrder
	err = json.Unmarshal(body, &order)
	if err != nil {
		i.log.LogWarning("err when parse JSON in Interrogator:", err)
		return "", 0, err
	}
	return order.OrderStatus, order.Accrual, nil
}

// UpdateOrderStatusInDB TODO: что делать с контекстом?
func (i *Interrogator) UpdateOrderStatusInDB(orderNumber string, newStatus string, accrual float32) string {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	userId, err := i.db.StatusSetter(ctx, orderNumber, newStatus, accrual)
	if err != nil {
		i.log.LogWarning("err when trying update order status", err)
	}
	i.log.LogInfo("Order", orderNumber, "updated in the database with status", newStatus)
	return userId
}

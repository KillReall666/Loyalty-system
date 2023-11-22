package _interrogator

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/KillReall666/Loyalty-system/internal/dto"
	"github.com/KillReall666/Loyalty-system/internal/logger"
	"github.com/KillReall666/Loyalty-system/internal/storage/postgres"
	"io"
	"net/http"
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

func (i *Interrogator) OrderStatusWorker(orderNumber string) {
	status, err := i.GetOrderStatusFromACCRUAL(orderNumber)
	if err != nil {
		i.log.LogWarning("Error retrieving order status from ACCRUAL: %v\n", err)
		return
	}

	switch status {
	case "PROCESSED":
		// Переместить заказ в базу данных с новым статусом (например, PROCESSED)
		i.UpdateOrderStatusInDB(orderNumber, "PROCESSED")
	case "INVALID":
		// Переместить заказ в базу данных с новым статусом (например, INVALID)
		i.UpdateOrderStatusInDB(orderNumber, "INVALID")
	default:
		// Пока заказ имеет статус отличный от PROCESSED и INVALID,
		// считаем его обработаным в ACCRUAL, и ничего не делаем
		fmt.Println("Order is still being processed in ACCRUAL")
	}
}

func (i *Interrogator) GetOrderStatusFromACCRUAL(orderNumber string) (string, error) {
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
		i.log.LogWarning("err when parse JSON:", err)
		return "", err
	}
	return order.OrderStatus, nil
}

func (i *Interrogator) UpdateOrderStatusInDB(orderNumber string, newStatus string) {
	err := i.db.StatusSetter(context.Background(), orderNumber, newStatus)
	if err != nil {
		i.log.LogWarning("err when trying update order status", err)
	}
	fmt.Printf("Order %s updated in the database with status %s\n", orderNumber, newStatus)
}

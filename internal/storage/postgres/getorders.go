package postgres

import (
	"context"
	"github.com/KillReall666/Loyalty-system/internal/dto"
)

func (d *Database) GetOrders(ctx context.Context, userID string) ([]dto.FullOrder, error) {
	getOrdersQuery := `
		SELECT  ordernumber, status, COALESCE(accrual, 0) AS accrual, orderdate
		FROM user_orders 
		WHERE userid = $1
		ORDER BY orderdate DESC 
`

	rows, err := d.db.Query(ctx, getOrdersQuery, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders := []dto.FullOrder{}
	for rows.Next() {
		var order dto.FullOrder
		err = rows.Scan(&order.OrderNumber, &order.OrderStatus, &order.Accrual, &order.OrderDate)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	return orders, nil
}

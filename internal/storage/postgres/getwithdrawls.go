package postgres

import (
	"context"
	"errors"
	"github.com/KillReall666/Loyalty-system/internal/dto"
)

func (d *Database) GetWithdrawals(ctx context.Context, userId string) (*dto.Billing, error) {
	getOrdersQuery := `
		SELECT  ordernumber, sum, processed_at
		FROM billing
		WHERE userid = $1
		ORDER BY processed_at ASC
`
	rows, err := d.db.Query(ctx, getOrdersQuery, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders := dto.Billing{}
	for rows.Next() {
		err = rows.Scan(&orders.Order, &orders.Sum, &orders.ProcessedAt)
		if err != nil {
			return nil, err
		}
	}

	if len(orders.Order) == 0 {
		return nil, errors.New("no withdrawals found")
	}

	return &orders, nil
}

package postgres

import (
	"context"
	"errors"
	"github.com/KillReall666/Loyalty-system/internal/dto"
)

func (d *Database) GetWithdrawals(ctx context.Context, userID string) ([]*dto.Billing, error) {
	getOrdersQuery := `
		SELECT  ordernumber, sum, processed_at
		FROM billing
		WHERE userid = $1
		ORDER BY processed_at ASC
`
	rows, err := d.db.Query(ctx, getOrdersQuery, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders := []*dto.Billing{}
	for rows.Next() {
		billing := dto.Billing{}
		err = rows.Scan(&billing.Order, &billing.Sum, &billing.ProcessedAt)
		if err != nil {
			return nil, err
		}
		//billing.ProcessedAt.Format(time.RFC3339)
		orders = append(orders, &billing)
	}

	if len(orders) == 0 {
		return nil, errors.New("no withdrawals found")
	}

	return orders, nil

}

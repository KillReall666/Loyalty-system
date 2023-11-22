package postgres

import "context"

func (d *Database) GetOrderNumbers(ctx context.Context) ([]string, error) {
	query := `
        SELECT ordernumber
        FROM user_orders
        WHERE status <> 'PROCESSED' AND status <> 'INVALID'
    `

	rows, err := d.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orderNumbers []string
	for rows.Next() {
		var orderNumber string
		err = rows.Scan(&orderNumber)
		if err != nil {
			return nil, err
		}
		orderNumbers = append(orderNumbers, orderNumber)
	}

	return orderNumbers, nil
}
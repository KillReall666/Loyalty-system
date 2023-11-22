package postgres

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgconn"
)

func (d *Database) StatusSetter(ctx context.Context, orderNumber, orderStatus string, accrual float32) (string, error) {
	insertQuery := `
                UPDATE user_orders 
                SET status = $2, accrual = $3 
                WHERE orderNumber = $1
				RETURNING userid
            `
	var userId string
	err := d.db.QueryRow(ctx, insertQuery, orderNumber, orderStatus, accrual).Scan(&userId)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == "23505" {
				return "", errors.New("this order already exists, please try another one")
			} else {
				return "", err
			}
		}
	}

	return userId, nil
}

package postgres

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgconn"
)

func (d *Database) OrderSetter(ctx context.Context, userId, orderNumber string) error {
	insertQuery := `
                INSERT INTO user_orders (userid, orderNumber, status)
				VALUES ($1, $2, $3)
            `

	_, err := d.db.Exec(ctx, insertQuery, userId, orderNumber, "NEW")
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == "23505" {
				return errors.New("this order already exists, please try another one")
			} else {
				return err
			}
		}
	}

	return nil
}

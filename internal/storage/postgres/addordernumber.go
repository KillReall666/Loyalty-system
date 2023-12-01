package postgres

import (
	"context"
	"github.com/KillReall666/Loyalty-system/internal/handlers/addorder"
	"github.com/jackc/pgx/v5"
)

func (d *Database) OrderSetter(ctx context.Context, userId, orderNumber string) error {
	var existingUserID string
	selectQuery := `SELECT userid FROM user_orders WHERE orderNumber = $1 LIMIT 1`
	err := d.db.QueryRow(ctx, selectQuery, orderNumber).Scan(&existingUserID)
	if err != nil {
		if err == pgx.ErrNoRows {
			insertQuery := `INSERT INTO user_orders (userid, orderNumber, status) VALUES ($1, $2, $3)`
			_, err = d.db.Exec(ctx, insertQuery, userId, orderNumber, "NEW")
			return err
		}
		return err
	}

	if existingUserID != userId {
		return addorder.ErrDifferentUser
	}

	return addorder.ErrOrderExists
}

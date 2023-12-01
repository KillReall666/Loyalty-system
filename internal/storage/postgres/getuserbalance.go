package postgres

import (
	"context"
	"fmt"
	"github.com/KillReall666/Loyalty-system/internal/dto"
	"github.com/jackc/pgx/v5"
)

func (d *Database) GetUserBalance(ctx context.Context, userID string) (*dto.UserBalance, error) {
	query := `
        SELECT COALESCE(current, 0), COALESCE(withdrawn, 0)
        FROM user_balance
        WHERE userid = $1
    `

	rows := d.db.QueryRow(ctx, query, userID)

	var userBalance dto.UserBalance
	err := rows.Scan(&userBalance.Current, &userBalance.Withdrawn)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user balance not found for UserID: %s", userID)
		}
		return nil, err
	}

	return &userBalance, nil
}

package postgres

import "context"

func (d *Database) IncrementCurrent(ctx context.Context, userID string, value float32) error {
	insertQuery := `
        INSERT INTO user_balance (userID, current)
        VALUES ($1, $2)
        ON CONFLICT (userId)
        DO UPDATE SET current = user_balance.current + $2
    `

	_, err := d.db.Exec(ctx, insertQuery, userID, value)
	if err != nil {
		return err
	}

	return nil
}

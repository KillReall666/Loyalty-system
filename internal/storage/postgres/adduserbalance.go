package postgres

import "context"

func (d *Database) IncrementCurrent(ctx context.Context, userId string, value float32) error {
	insertQuery := `
        INSERT INTO user_balance (userId, current)
        VALUES ($1, $2)
        ON CONFLICT (userId)
        DO UPDATE SET current = user_balance.current + $2
    `

	_, err := d.db.Exec(ctx, insertQuery, userId, value)
	if err != nil {
		return err
	}

	return nil
}

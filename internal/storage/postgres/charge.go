package postgres

import (
	"context"
	"fmt"
)

func (d *Database) ProcessOrder(ctx context.Context, order, userId string, sum float32) error {
	// Проверяем баланс в таблице user_balance
	checkBalanceQuery := `
        SELECT current
        FROM user_balance
        WHERE current >= $1
    `
	var currentBalance float32
	err := d.db.QueryRow(ctx, checkBalanceQuery, sum).Scan(&currentBalance)
	if err != nil {
		return fmt.Errorf("failed to check balance: %s", err)
	}

	// Вычитаем сумму из баланса и обновляем значение в withdrawn
	updateBalanceQuery := `
        UPDATE user_balance
        SET current = current - $1, withdrawn = COALESCE(withdrawn, 0) + $1
    `
	_, err = d.db.Exec(ctx, updateBalanceQuery, sum)
	if err != nil {
		return fmt.Errorf("failed to update balance: %s", err)
	}

	// Вставляем информацию о заказе в таблицу billing
	insertBillingQuery := `
        INSERT INTO billing (userId, ordernumber, sum, processed_at)
        VALUES ($1, $2, $3, NOW())
    `
	_, err = d.db.Exec(ctx, insertBillingQuery, userId, order, sum)
	if err != nil {
		return fmt.Errorf("failed to insert billing record: %s", err)
	}

	return nil
}

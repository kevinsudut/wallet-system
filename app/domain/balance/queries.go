package domainbalance

const (
	queryGetBalanceByUsername = `
		SELECT username, amount FROM balances WHERE username = $1;
	`

	queryGrantBalanceByUsername = `
		INSERT INTO balances (username, amount) VALUES ($1, $2)
		ON CONFLICT (username)
		DO UPDATE SET
			amount = balances.amount + EXCLUDED.amount,
			updated_at = NOW();
	`

	queryDeductBalanceByUsername = `
		UPDATE balances SET
			amount = amount - $1,
			updated_at = NOW()
		WHERE username = $2 AND amount - $1 >= 0; 
	`

	queryInsertHistory = `
		INSERT INTO histories (id, username, target_username, amount, type, notes) VALUES ($1, $2, $3, $4, $5, $6);
	`

	queryUpdateHistorySummaryById = `
		INSERT INTO history_summaries (id, username, target_username, amount, type) VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (id)
		DO UPDATE SET
			amount = history_summaries.amount + EXCLUDED.amount,
			updated_at = NOW();
	`

	queryGetLatestHistoryByUsername = `
		SELECT id, username, target_username, amount, type, notes FROM histories WHERE username = $1 ORDER BY created_at DESC LIMIT $2;
	`

	queryGetHistorySummaryByUsernameAndType = `
		SELECT username, target_username, amount, type FROM history_summaries WHERE username = $1 AND type = $2 ORDER BY amount DESC LIMIT $3; 
	`
)

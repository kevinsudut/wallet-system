package domainbalance

const (
	queryGetBalanceByUserId = `
		SELECT
			user_id,
			amount
		FROM
			balances
		WHERE
			user_id = $1;
	`

	queryGrantBalanceByUserId = `
		INSERT INTO balances (user_id, amount) VALUES ($1, $2)
		ON CONFLICT (user_id)
		DO UPDATE SET
			amount = balances.amount + EXCLUDED.amount,
			updated_at = NOW();
	`

	queryDeductBalanceByUserId = `
		UPDATE balances SET
			amount = amount - $1,
			updated_at = NOW()
		WHERE user_id = $2 AND amount - $1 >= 0; 
	`

	queryInsertHistory = `
		INSERT INTO histories (id, user_id, target_user_id, amount, type, notes) VALUES ($1, $2, $3, $4, $5, $6);
	`

	queryUpdateHistorySummaryById = `
		INSERT INTO history_summaries (id, user_id, target_user_id, amount, type) VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (id)
		DO UPDATE SET
			amount = history_summaries.amount + EXCLUDED.amount,
			updated_at = NOW();
	`

	queryGetLatestHistoryByUserId = `
		SELECT
			id,
			user_id,
			target_user_id,
			amount,
			type,
			notes
		FROM
			histories
		WHERE
			user_id = $1
		ORDER BY created_at DESC
		LIMIT 10;
	`

	queryGetHistorySummaryByUserIdAndType = `
		SELECT
			user_id,
			target_user_id,
			amount,
			type
		FROM
			history_summaries
		WHERE
			user_id = $1 AND type = $2
		ORDER BY amount DESC
		LIMIT 10; 
	`
)

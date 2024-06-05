package domainbalance

const (
	queryGetBalanceByUsername = `
		SELECT username, amount FROM balances WHERE username = $1;
	`

	queryGrantBalanceByUsername = `
		INSERT INTO balances (username, amount) VALUES ($1, $2)
		ON CONFLICT (username)
		DO UPDATE SET
			amount = EXCLUDED.amount + $2
	`
)

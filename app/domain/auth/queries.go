package domainauth

const (
	queryInsertUser = `
		INSERT INTO users (id, username) VALUES ($1, $2);
	`

	queryGetUserByUsername = `
		SELECT
			id,
			username
		FROM
			users
		WHERE
			username = $1;
	`
)

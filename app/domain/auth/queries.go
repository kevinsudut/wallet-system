package domainauth

const (
	queryInsertUser = `
		INSERT INTO users (id, username) VALUES ($1, $2);
	`

	queryGetUserById = `
		SELECT
			id,
			username
		FROM
			users
		WHERE
			id = $1;
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

package domainauth

const (
	queryInsertUser = `
		INSERT INTO users (username) VALUES ($1);
	`

	queryGetUserByUsername = `
		SELECT username FROM users WHERE username = $1;
	`
)

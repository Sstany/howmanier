package db

const (
	queryUser  = "SELECT * FROM users WHERE id=$1"
	insertUser = "INSERT INTO users (id, username) values ($1, $2)"
)

const (
	insertList = "INSERT INTO fridge (id, owner_id) values ($1, $2)"
)

const (
	queryInitUsers = `CREATE TABLE IF NOT EXISTS users (
		id bigint NOT NULL,
		username text NOT NULL,
		PRIMARY KEY (id)
	  )`
	queryInitList = `CREATE TABLE IF NOT EXISTS list(
		id bigint NOT NULL,
		PRIMARY KEY (id)
		name text NOT NULL,
		count int NOT NULL,
	)`
)

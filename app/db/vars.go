package db

const (
	queryUser  = "SELECT * FROM users WHERE id=$1"
	insertUser = "INSERT INTO users (id, username) values ($1, $2)"
)

const (
	insertFridge = "INSERT INTO fridge (user_id, name, count) values ($1, $2, $3)"
	deleteFridge = "UPDATE fridge SET count=count-$2 WHERE name=$1"
	queryFridge  = "SELECT * FROM fridge WHERE name=$2"
)

const (
	queryInitUsers = `CREATE TABLE IF NOT EXISTS users (
		id bigint NOT NULL,
		username text NOT NULL,
		PRIMARY KEY (id)
	  )`
	queryInitFridge = `CREATE TABLE IF NOT EXISTS fridge(
		id serial NOT NULL,
		user_id bigint NOT NULL,
		name text NOT NULL,
		count int NOT NULL,
		PRIMARY KEY (id)
	)`
)

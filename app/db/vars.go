package db

const (
	queryUser  = "SELECT * FROM users WHERE id=$1"
	insertUser = "INSERT INTO users (id, username) values ($1, $2)"
)

const (
	insertFridge = `INSERT INTO fridge (user_id, name, count) values ($1, $2, $3)	
	 ON CONFLICT (name, user_id) DO UPDATE SET count=fridge.count+EXCLUDED.count`
	deleteFridge = "UPDATE fridge SET count=count-$3 WHERE user_id=$1 AND name=$2"
	queryFridge  = "SELECT * FROM fridge WHERE name=$2"
	listFridge   = `SELECT name,count FROM fridge WHERE user_id=$1`
)
const (
	insertRecipe = `INSERT INTO recipes (user_id, recipe_name, name, count) values ($1,$2,$3,$4) 
	ON CONFLICT (recipe_name, user_id) DO NOTHING`
	listRecipes = `SELECT recipe_name FROM recipes WHERE user_id = $1`
)

const (
	queryInitUsers = `CREATE TABLE IF NOT EXISTS users (
		id bigint NOT NULL,
		username text NOT NULL,
		PRIMARY KEY (id)
	  )`
	queryInitFridge = `CREATE TABLE IF NOT EXISTS fridge(
		id serial UNIQUE,
		user_id bigint NOT NULL,
		name text NOT NULL,
		count int NOT NULL,
		PRIMARY KEY (name, user_id)
		
	)`
	queryInitRecipes = `CREATE TABLE IF NOT EXISTS recipes(
		id serial UNIQUE,
		user_id bigint NOT NULL,
		recipe_name text,
		name text NOT NULL,
		count int NOT NULL,
		PRIMARY KEY (recipe_name, user_id)
	)`
)

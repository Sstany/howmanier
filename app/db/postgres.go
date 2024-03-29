package db

import (
	"context"
	"database/sql"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

type PostgresClient struct {
	conn   *sql.Conn
	logger *zap.Logger
}

type User struct {
	ID       int64
	Username string
	ListID   int64
}

type Product struct {
	ID     int64
	UserID int64
	Name   string
	Count  int
}
type Recipe struct {
	ID         int64
	UserID     int64
	RecipeName string
	Name       string
	Count      int
}

func NewPostgresClient(ctx context.Context, connStr string) *PostgresClient {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	conn, err := db.Conn(ctx)
	if err != nil {
		panic(err)
	}

	_, err = conn.ExecContext(ctx, queryInitUsers)
	if err != nil {
		panic(err)
	}

	_, err = conn.ExecContext(ctx, queryInitFridge)
	if err != nil {
		panic(err)
	}

	_, err = conn.ExecContext(ctx, queryInitRecipes)
	if err != nil {
		panic(err)
	}

	return &PostgresClient{
		conn: conn,
	}
}

func (r *PostgresClient) CreateUser(ctx context.Context, user *User) error {
	_, err := r.conn.ExecContext(ctx, insertUser, user.ID, user.Username)
	if err != nil {
		return err
	}

	return nil
}

func (r *PostgresClient) FetchUser(ctx context.Context, user *User) error {
	row := r.conn.QueryRowContext(ctx, queryUser, user.ID)

	if err := row.Scan(&user.ID, &user.Username); err != nil {
		return err
	}

	return nil
}

func (r *PostgresClient) AddProduct(ctx context.Context, product *Product) error {
	_, err := r.conn.ExecContext(ctx, insertFridge, product.UserID, product.Name, product.Count)
	if err != nil {
		return err
	}

	return nil
}

func (r *PostgresClient) DeleteProduct(ctx context.Context, product *Product) error {
	_, err := r.conn.ExecContext(ctx, deleteFridge, product.UserID, product.Name, product.Count)
	if err != nil {
		return err
	}

	return nil
}

func (r *PostgresClient) ListFridge(ctx context.Context, user *User) ([]Product, error) {
	rows, err := r.conn.QueryContext(ctx, listFridge, user.ID)
	if err != nil {
		return nil, err
	}

	var products []Product
	var tempProduct Product

	for rows.Next() {
		if err := rows.Scan(&tempProduct.Name, &tempProduct.Count); err != nil {
			return nil, err
		}

		products = append(products, tempProduct)

	}

	return products, nil
}
func (r *PostgresClient) AddRecipe(ctx context.Context, recipe *Recipe) error {
	_, err := r.conn.ExecContext(ctx, insertRecipe, recipe.UserID, recipe.RecipeName, recipe.Name, recipe.Count)
	if err != nil {
		return err
	}

	return nil
}

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

type Database struct {
	db     *sql.DB
	logger *zap.Logger
}

type User struct {
	ID       int64
	Username string
	ListID   int64
}

type List struct {
	ID      int64
	OwnerID int64
	Name    string
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

func (r *PostgresClient) List(ctx context.Context, wishlist *List) error {
	_, err := r.conn.ExecContext(ctx, insertList, wishlist.ID, wishlist.OwnerID, wishlist.Name)
	if err != nil {
		return err
	}

	return nil
}
package storage

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthStorage struct {
	pool *pgxpool.Pool
}

func NewAuthStorage(pool *pgxpool.Pool) *AuthStorage {
	return &AuthStorage{pool: pool}
}

var (
	UserTable    = "users"
	RefreshTable = "refresh_tokens"
)

func (a *AuthStorage) CreateUser(ctx context.Context, Email, Password string) (int, error) {
	var ID int
	m := make(map[string]interface{})
	m["email"] = Email
	m["password"] = Password

	query, args, err := squirrel.Insert(UserTable).SetMap(m).PlaceholderFormat(squirrel.Dollar).Suffix("RETURNING ID").ToSql()
	if err != nil {
		return 0, err
	}
	Row := a.pool.QueryRow(ctx, query, args...)
	if err != nil {
		return 0, err
	}

	err = Row.Scan(&ID)
	if err != nil {
		return 0, err
	}

	return ID, nil
}

func (a *AuthStorage) AuthUser(ctx context.Context, Email, PasswordHash string) (int, error) {
	var ID int
	query, args, err := squirrel.Select("id").PlaceholderFormat(squirrel.Dollar).From(UserTable).Where(squirrel.Eq{"email": Email, "password_hash": PasswordHash}).ToSql()
	if err != nil {
		return 0, err
	}

	row := a.pool.QueryRow(ctx, query, args...)

	err = row.Scan(&ID)
	if err != nil {
		return 0, err
	}

	return ID, nil
}

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

func (a *AuthStorage) AddRefresh(ctx context.Context, AccessCode, RefreshToken string) error {
	m := make(map[string]interface{})
	m["access_code"] = AccessCode
	m["refresh_token"] = RefreshToken
	query, args, err := squirrel.Insert(RefreshTable).
		Columns("refresh_token", "access_code").
		Values(RefreshToken, AccessCode).
		PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return err
	}
	_, err = a.pool.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (a *AuthStorage) CheckRefresh(ctx context.Context, AccessCode string) (string, error) {
	var Refresh string
	query, args, err := squirrel.Select("refresh_token").
		Where(squirrel.Eq{"access_code": AccessCode}).PlaceholderFormat(squirrel.Dollar).
		From(RefreshTable).ToSql()
	if err != nil {
		return "", err
	}

	row := a.pool.QueryRow(ctx, query, args...)

	err = row.Scan(&Refresh)
	if err != nil {
		return "", err
	}
	return Refresh, nil
}

func (a *AuthStorage) DeleteRefresh(ctx context.Context, AccessCode string) error {
	query, args, err := squirrel.Delete(UserTable).Where(squirrel.Eq{"access_code": AccessCode}).PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return err
	}

	_, err = a.pool.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
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

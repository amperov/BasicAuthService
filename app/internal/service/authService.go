package service

import "context"

type AuthStorage interface {
	CreateUser(ctx context.Context, Email, Password string) (int, error)
	AddRefresh(ctx context.Context, AccessCode, RefreshToken string) error
	CheckRefresh(ctx context.Context, AccessCode string) (string, error)
	DeleteRefresh(ctx context.Context, AccessCode string) error

	AuthUser(ctx context.Context, Email, PasswordHash string) (int, error)
}

type AuthRedis interface {
	InsertAccessToken(ctx context.Context, AccessCode, AccessToken string) error
	GetAccessToken(ctx context.Context, AccessCode string) (string, error)
}

type Hasher interface {
	HashPassword(Password string) string
	GenerateAccessCode(AccessToken, RefreshToken string) (string, error)
}

type TokenManager interface {
	GenerateToken(ctx context.Context, UserID int) (string, string)
	ValidateToken(ctx context.Context, Token string) (int, error)
}

type AuthService struct {
	AuthStorage  AuthStorage
	RedisClient  AuthRedis
	TokenManager TokenManager
	Hasher       Hasher
}

func NewAuthService(authStorage AuthStorage, redisClient AuthRedis, tokenManager TokenManager, hasher Hasher) *AuthService {
	return &AuthService{AuthStorage: authStorage, RedisClient: redisClient, TokenManager: tokenManager, Hasher: hasher}
}

func (a *AuthService) SignUp(ctx context.Context, email, password string) (int, string, string, error) {
	UserID, err := a.AuthStorage.CreateUser(ctx, email, a.Hasher.HashPassword(password))
	if err != nil {
		return 0, "", err.Error(), err
	}

	AccessToken, RefreshToken := a.TokenManager.GenerateToken(ctx, UserID)

	AccessCode, err := a.Hasher.GenerateAccessCode(AccessToken, RefreshToken)
	if err != nil {
		return 0, "", err.Error(), err
	}

	err = a.RedisClient.InsertAccessToken(ctx, AccessCode, AccessToken)
	if err != nil {
		return 0, "", err.Error(), err
	}

	err = a.AuthStorage.AddRefresh(ctx, AccessCode, RefreshToken)
	if err != nil {
		return 0, "", err.Error(), err
	}

	return UserID, AccessCode, "", nil
}

func (a *AuthService) SignIn(ctx context.Context, email, password string) (int, string, string, error) {
	UserID, err := a.AuthStorage.AuthUser(ctx, email, a.Hasher.HashPassword(password))
	if err != nil {
		return 0, "", "", err
	}
	AccessToken, RefreshToken := a.TokenManager.GenerateToken(ctx, UserID)

	AccessCode, err := a.Hasher.GenerateAccessCode(AccessToken, RefreshToken)
	if err != nil {
		return 0, "", err.Error(), err
	}

	err = a.RedisClient.InsertAccessToken(ctx, AccessCode, AccessToken)
	if err != nil {
		return 0, "", err.Error(), err
	}

	err = a.AuthStorage.AddRefresh(ctx, AccessCode, RefreshToken)
	if err != nil {
		return 0, "", err.Error(), err
	}

	return UserID, AccessCode, "", nil
}

func (a *AuthService) Identify(ctx context.Context, AccessCode string) (int, string, error) {
	AccessToken, err := a.RedisClient.GetAccessToken(ctx, AccessCode)
	if err != nil {
		return 0, err.Error(), err
	}

	UserID, err := a.TokenManager.ValidateToken(ctx, AccessToken)
	if err != nil {
		RefreshToken, err := a.AuthStorage.CheckRefresh(ctx, AccessCode)
		if err != nil {
			return 0, err.Error(), err
		}

		UserID, err := a.TokenManager.ValidateToken(ctx, RefreshToken)
		if err != nil {
			return 0, err.Error(), err
		}

		AccessToken, RefreshToken = a.TokenManager.GenerateToken(ctx, UserID)
		AccessCode, err := a.Hasher.GenerateAccessCode(AccessToken, RefreshToken)

		err = a.RedisClient.InsertAccessToken(ctx, AccessCode, AccessToken)
		if err != nil {
			return 0, "", err
		}

		err = a.AuthStorage.DeleteRefresh(ctx, AccessCode)
		if err != nil {
			return 0, "", err
		}

		err = a.AuthStorage.AddRefresh(ctx, AccessCode, RefreshToken)
		if err != nil {
			return 0, "", err
		}

		return UserID, AccessCode, err
	}

	return UserID, AccessCode, nil
}

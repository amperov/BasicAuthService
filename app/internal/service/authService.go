package service

import "context"

type AuthStorage interface {
	CreateUser(ctx context.Context, Email, Password string) (int, error)

	AuthUser(ctx context.Context, Email, PasswordHash string) (int, error)
}

type Hasher interface {
	HashPassword(Password string) string
}

type TokenManager interface {
	GenerateToken(ctx context.Context, UserID int) (string, string)
	ValidateToken(ctx context.Context, Token string) (int, error)
}

type AuthService struct {
	AuthStorage AuthStorage

	TokenManager TokenManager
	Hasher       Hasher
}

func NewAuthService(authStorage AuthStorage, tokenManager TokenManager, hasher Hasher) *AuthService {
	return &AuthService{AuthStorage: authStorage, TokenManager: tokenManager, Hasher: hasher}
}

func (a *AuthService) SignUp(ctx context.Context, email, password string) (int, string, string, string, error) {
	UserID, err := a.AuthStorage.CreateUser(ctx, email, a.Hasher.HashPassword(password))
	if err != nil {
		return 0, "", "", err.Error(), err
	}

	AccessToken, RefreshToken := a.TokenManager.GenerateToken(ctx, UserID)

	return UserID, AccessToken, RefreshToken, "success", nil
}

func (a *AuthService) SignIn(ctx context.Context, email, password string) (int, string, string, string, error) {
	UserID, err := a.AuthStorage.AuthUser(ctx, email, a.Hasher.HashPassword(password))
	if err != nil {
		return 0, "", "", err.Error(), err
	}
	AccessToken, RefreshToken := a.TokenManager.GenerateToken(ctx, UserID)

	return UserID, AccessToken, RefreshToken, "", nil
}

func (a *AuthService) Identify(ctx context.Context, AccessToken string, RefreshToken string) (int, string, string, error) {
	UserID, err := a.TokenManager.ValidateToken(ctx, AccessToken)
	if err != nil {

		UserID, err := a.TokenManager.ValidateToken(ctx, RefreshToken)
		if err != nil {
			return 0, "", err.Error(), err
		}

		AccessToken, RefreshToken = a.TokenManager.GenerateToken(ctx, UserID)

		return UserID, AccessToken, RefreshToken, err
	}

	return UserID, AccessToken, RefreshToken, nil
}

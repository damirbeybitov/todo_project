package auth

import (
	"context"

	"github.com/damirbeybitov/todo_project/internal/auth/repository"
	"github.com/damirbeybitov/todo_project/internal/log"
	token "github.com/damirbeybitov/todo_project/internal/token"
	authPB "github.com/damirbeybitov/todo_project/proto/auth"
)

// AuthService представляет сервис аутентификации.
type AuthService struct {
	repo *repository.Repository
	authPB.UnimplementedAuthServiceServer
}

// NewAuthService создает новый экземпляр AuthService.
func NewAuthService(repo *repository.Repository) authPB.AuthServiceServer {
	return &AuthService{repo: repo}
}

// Authenticate реализует метод аутентификации в рамках интерфейса AuthServiceServer.
func (s *AuthService) Authenticate(ctx context.Context, req *authPB.AuthenticateRequest) (*authPB.AuthenticateResponse, error) {
	log.InfoLogger.Printf("Authenticating user with username: %s", req.Username)

	// Реализация аутентификации пользователя
	if err := s.repo.CheckPassword(req.Username, req.Password); err != nil {
		return nil, err
	}

	accessToken, refreshToken, err := s.repo.GenerateTokens(req.Username)
	if err != nil {
		return nil, err
	}

	// В данном примере просто возвращается фиктивный access token и refresh token.
	return &authPB.AuthenticateResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// RefreshToken реализует метод обновления токена в рамках интерфейса AuthServiceServer.
func (s *AuthService) RefreshToken(ctx context.Context, req *authPB.RefreshTokenRequest) (*authPB.RefreshTokenResponse, error) {
	log.InfoLogger.Println("Refreshing token")

	// Реализация обновления токена
	accessToken, err := token.RefreshToken(req.RefreshToken)
	if err != nil {
		return nil, err
	}

	// В данном примере просто возвращается фиктивный access token.
	return &authPB.RefreshTokenResponse{
		AccessToken: accessToken,
	}, nil
}

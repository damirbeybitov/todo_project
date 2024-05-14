package auth

import (
	"context"
	"database/sql"

	"github.com/damirbeybitov/todo_project/internal/log"
	authPB "github.com/damirbeybitov/todo_project/proto/auth"
)

// AuthService представляет сервис аутентификации.
type AuthService struct {
	db *sql.DB
	authPB.UnimplementedAuthServiceServer
}

// NewAuthService создает новый экземпляр AuthService.
func NewAuthService(db *sql.DB) authPB.AuthServiceServer {
	return &AuthService{db: db}
}

// Authenticate реализует метод аутентификации в рамках интерфейса AuthServiceServer.
func (s *AuthService) Authenticate(ctx context.Context, req *authPB.AuthenticateRequest) (*authPB.AuthenticateResponse, error) {
	log.InfoLogger.Printf("Authenticating user with username: %s", req.Username)

	// Реализация аутентификации пользователя

	// В данном примере просто возвращается фиктивный access token и refresh token.
	return &authPB.AuthenticateResponse{
		AccessToken:  "fake_access_token",
		RefreshToken: "fake_refresh_token",
	}, nil
}

// RefreshToken реализует метод обновления токена в рамках интерфейса AuthServiceServer.
func (s *AuthService) RefreshToken(ctx context.Context, req *authPB.RefreshTokenRequest) (*authPB.RefreshTokenResponse, error) {
	log.InfoLogger.Println("Refreshing token")

	// Реализация обновления токена

	// В данном примере просто возвращается фиктивный access token.
	return &authPB.RefreshTokenResponse{
		AccessToken: "fake_access_token",
	}, nil
}

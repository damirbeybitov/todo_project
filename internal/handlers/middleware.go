package handlers

import (
	"net/http"
	"strings"

	"github.com/damirbeybitov/todo_project/internal/log"
	token "github.com/damirbeybitov/todo_project/internal/token"
)

func (h *Handler) UserIdentity(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userToken := r.Header.Get("Authorization")
		if userToken == "" {
			log.ErrorLogger.Print("Token is missing in Header")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		userToken = strings.TrimPrefix(userToken, "Bearer ")

		_, err := token.VerifyToken(userToken)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
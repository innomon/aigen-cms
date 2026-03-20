package api

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/formcms/formcms-go/core/services"
)

type AuthApi struct {
	authService services.IAuthService
}

func NewAuthApi(authService services.IAuthService) *AuthApi {
	return &AuthApi{authService: authService}
}

func (a *AuthApi) Register(r chi.Router) {
	r.Route("/api/auth", func(r chi.Router) {
		r.Post("/register", a.DoRegister)
		r.Post("/login", a.DoLogin)
		r.With(a.JWTMiddleware).Get("/me", a.GetMe)
	})
}

type authRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (a *AuthApi) DoRegister(w http.ResponseWriter, r *http.Request) {
	var req authRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := a.authService.Register(r.Context(), req.Email, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(user)
}

func (a *AuthApi) DoLogin(w http.ResponseWriter, r *http.Request) {
	var req authRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	token, err := a.authService.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func (a *AuthApi) GetMe(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("userId").(int64)
	user, err := a.authService.Me(r.Context(), userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(user)
}

func (a *AuthApi) JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "missing authorization header", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "invalid authorization header", http.StatusUnauthorized)
			return
		}

		userId, role, err := a.authService.ValidateToken(parts[1])
		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "userId", userId)
		ctx = context.WithValue(ctx, "role", role)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

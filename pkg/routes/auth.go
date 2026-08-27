package routes

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	util "github.com/tade3910/recipe_server/pkg"
	"github.com/tade3910/recipe_server/pkg/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type authHandler struct {
	db *gorm.DB
}

func NewAuthHandler(db *gorm.DB) *authHandler {
	return &authHandler{db: db}
}

type authPayload struct {
	Email    string
	Password string
}

func (h *authHandler) respondWithToken(email string, w http.ResponseWriter) {
	// Generate tokens
	accessToken := util.GenerateRandomToken(32)
	refreshToken := util.GenerateRandomToken(64)

	now := time.Now()
	accessExpiry := now.Add(15 * time.Minute)
	refreshExpiry := now.Add(7 * 24 * time.Hour)

	h.db.Create(&models.AuthToken{
		UserEmail: email,
		Token:     accessToken,
		TokenType: "access",
		CreatedAt: now,
		ExpiresAt: accessExpiry,
	})
	h.db.Create(&models.AuthToken{
		UserEmail: email,
		Token:     refreshToken,
		TokenType: "refresh",
		CreatedAt: now,
		ExpiresAt: refreshExpiry,
	})

	util.RespondWithJSON(w, http.StatusOK, map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

// SignUp registers a new user
func (h *authHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	payload := &authPayload{}
	err := util.GetBody(r.Body, payload)
	if err != nil {
		util.RespondWithError(w, http.StatusBadRequest, "Invalid payload")
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	user := models.User{Email: payload.Email, Password: string(hashedPassword)}
	if err := h.db.Create(&user).Error; err != nil {
		util.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	h.respondWithToken(payload.Email, w)
}

// SignIn generates access + refresh tokens
func (h *authHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	payload := &authPayload{}
	err := util.GetBody(r.Body, payload)
	if err != nil {
		util.RespondWithError(w, http.StatusBadRequest, "Invalid payload")
		return
	}

	user := &models.User{}
	if err := h.db.First(user, "email = ?", payload.Email).Error; err != nil {
		util.RespondWithError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password)); err != nil {
		util.RespondWithError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}
	h.respondWithToken(payload.Email, w)
}

// RefreshToken issues a new access token using a valid refresh token
func (h *authHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		RefreshToken string `json:"refresh_token"`
	}
	err := util.GetBody(r.Body, &payload)
	if err != nil {
		util.RespondWithError(w, http.StatusBadRequest, "Invalid payload")
		return
	}

	token := &models.AuthToken{}
	if err := h.db.First(token, "token = ? AND token_type = ?", payload.RefreshToken, "refresh").Error; err != nil {
		util.RespondWithError(w, http.StatusUnauthorized, "Invalid refresh token")
		return
	}

	if token.ExpiresAt.Before(time.Now()) {
		util.RespondWithError(w, http.StatusUnauthorized, "Refresh token expired")
		return
	}

	// Create new access token
	newAccessToken := util.GenerateRandomToken(32)
	now := time.Now()
	accessExpiry := now.Add(15 * time.Minute)
	h.db.Create(&models.AuthToken{
		UserEmail: token.UserEmail,
		Token:     newAccessToken,
		TokenType: "access",
		CreatedAt: now,
		ExpiresAt: accessExpiry,
	})

	util.RespondWithJSON(w, http.StatusOK, map[string]string{
		"access_token": newAccessToken,
	})
}

func (h *authHandler) Logout(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token == "" {
		util.RespondWithError(w, http.StatusUnauthorized, "Missing token")
		return
	}

	h.db.Where("token = ?", token).Delete(&models.AuthToken{})
	util.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Logged out"})
}

// AuthMiddleware protects routes using access tokens
func (h *authHandler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		token = strings.TrimPrefix(token, "Bearer ")
		token = strings.TrimSpace(token)

		authToken := &models.AuthToken{}
		if err := h.db.First(authToken, "token = ? AND token_type = ?", token, "access").Error; err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if authToken.ExpiresAt.Before(time.Now()) {
			http.Error(w, "Access token expired", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), util.UserIDKey, authToken.UserEmail)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

var ErrUnauthorized = errors.New("unauthorized: missing or invalid user context")

// GetUserEmail extracts the user email from the request context.
func GetUserEmail(ctx context.Context) (string, error) {
	userEmail, ok := ctx.Value(util.UserIDKey).(string)
	if !ok || userEmail == "" {
		return "", ErrUnauthorized
	}
	return userEmail, nil
}

// RequireUserEmail extracts the user email or automatically writes a 401 response.
// Returns (email, true) if valid, or ("", false) if an error response was written.
func RequireUserEmail(w http.ResponseWriter, r *http.Request) (string, bool) {
	email, err := GetUserEmail(r.Context())
	if err != nil {
		util.RespondWithError(w, http.StatusUnauthorized, "Unauthorized: missing or invalid user context")
		return "", false
	}
	return email, true
}

package routes

import (
	"context"
	"net/http"
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

	util.RespondWithJSON(w, http.StatusCreated, map[string]string{
		"email": user.Email,
		"name":  user.Name,
	})
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

	// Generate tokens
	accessToken := util.GenerateRandomToken(32)
	refreshToken := util.GenerateRandomToken(64)

	now := time.Now()
	accessExpiry := now.Add(15 * time.Minute)
	refreshExpiry := now.Add(7 * 24 * time.Hour)

	h.db.Create(&models.AuthToken{
		UserEmail: user.Email,
		Token:     accessToken,
		TokenType: "access",
		CreatedAt: now,
		ExpiresAt: accessExpiry,
	})
	h.db.Create(&models.AuthToken{
		UserEmail: user.Email,
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

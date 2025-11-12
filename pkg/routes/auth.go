package routes

import (
	"context"
	"encoding/json"
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
	return &authHandler{
		db: db,
	}
}

func (h *authHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Email    string
		Password string
	}
	_ = json.NewDecoder(r.Body).Decode(&payload)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	user := models.User{Email: payload.Email, Password: string(hashedPassword)}
	if err := h.db.Create(&user).Error; err != nil {
		util.RespondWithError(w, http.StatusBadRequest, err)
		return
	}
	util.RespondWithJSON(w, http.StatusCreated, map[string]string{
		"email": user.Email,
		"name":  user.Name})
}

func (h *authHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Email    string
		Password string
	}
	_ = json.NewDecoder(r.Body).Decode(&payload)

	user := &models.User{
		Email: payload.Email,
	}
	if err := h.db.First(user).Error; err != nil {
		util.RespondWithError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password)); err != nil {
		util.RespondWithError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	// Generate token (random string or JWT)
	token := util.GenerateRandomToken(32) // implement this function
	authToken := &models.AuthToken{UserEmail: user.Email, Token: token, CreatedAt: time.Now()}
	h.db.Create(authToken)

	util.RespondWithJSON(w, http.StatusOK, map[string]string{"token": token})
}

func (h *authHandler) Logout(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token == "" {
		util.RespondWithError(w, http.StatusUnauthorized, "Missing token")
		return
	}

	auth_token := &models.AuthToken{
		Token: token,
	}

	h.db.Delete(auth_token)
	util.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Logged out"})
}

func (h *authHandler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		auth_token := &models.AuthToken{
			Token: token,
		}
		if err := h.db.First(auth_token).Error; err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), util.UserIDKey, auth_token.UserEmail)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

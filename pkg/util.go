package util

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/joho/godotenv"
)

func RespondWithJSON(w http.ResponseWriter, code int, payload any) error {
	response, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(code)
	w.Write(response)
	return nil
}

func RespondWithError(w http.ResponseWriter, code int, payload any) error {
	switch v := payload.(type) {
	case string:
		return RespondWithJSON(w, code, map[string]string{"error": v})
	case error:
		return RespondWithJSON(w, code, map[string]string{"error": v.Error()})
	default:
		return RespondWithJSON(w, code, payload)
	}
}

func GetBody[T any](Body io.ReadCloser, bodyStruct *T) error {
	defer Body.Close()
	body, err := io.ReadAll(Body)
	if err != nil {
		return fmt.Errorf("could not read body")
	}
	err = json.Unmarshal(body, bodyStruct)
	if err != nil {
		return fmt.Errorf("invalid body")
	}
	return nil
}

type envs struct {
	Port           string
	DbUrl          string
	TestUrl        string
	AllowedOrigins []string
}

var loadedEnv *envs

func LoadEnvs() *envs {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Join(filepath.Dir(b), "../.env")
	godotenv.Load(basepath)

	if loadedEnv != nil {
		return loadedEnv
	}
	port := os.Getenv("PORT")
	dbUrl := os.Getenv("DATABASE_URL")
	testUrl := os.Getenv("TEST_DATABASE_DSN")
	allowedOrigins := strings.Split(os.Getenv("FRONTEND_URLS"), ",")
	return &envs{
		Port:           port,
		DbUrl:          dbUrl,
		TestUrl:        testUrl,
		AllowedOrigins: allowedOrigins,
	}
}

func IsUrl(str string) bool {
	// u, err := url.Parse(str)
	// return err == nil && u.Scheme != "" && u.Host != ""
	return true
}

func GenerateRandomToken(length int) string {
	if length <= 0 {
		length = 32 // default length
	}

	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		// fallback in case crypto/rand fails (very unlikely)
		panic(fmt.Sprintf("failed to generate random token: %v", err))
	}

	return hex.EncodeToString(bytes)
}

type contextKey string

// Package-level constant
const UserIDKey contextKey = "userID"

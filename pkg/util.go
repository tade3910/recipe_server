package util

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
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

func GetBody[T any](body io.ReadCloser, bodyStruct *T) error {
	defer body.Close()

	decoder := json.NewDecoder(body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(bodyStruct); err != nil {
		if errors.Is(err, io.EOF) {
			return fmt.Errorf("request body cannot be empty")
		}
		if strings.HasPrefix(err.Error(), "json: unknown field") {
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("unknown field %s is not allowed", fieldName)
		}
		return fmt.Errorf("invalid body: %w", err)
	}

	if decoder.More() {
		return fmt.Errorf("body must only contain a single JSON object")
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

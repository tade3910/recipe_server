package util

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	"github.com/joho/godotenv"
)

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) error {
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

func RespondWithError(w http.ResponseWriter, code int, msg string) error {
	return RespondWithJSON(w, code, map[string]string{"error": msg})
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
	Port    string
	DbUrl   string
	TestUrl string
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
	if port == "" {
		log.Fatal("Could not read port from .env file")
	}
	if dbUrl == "" {
		log.Fatal("Could not read dbUrl from .env file")
	}
	return &envs{
		Port:    port,
		DbUrl:   dbUrl,
		TestUrl: testUrl,
	}
}

func IsUrl(str string) bool {
	// u, err := url.Parse(str)
	// return err == nil && u.Scheme != "" && u.Host != ""
	return true
}

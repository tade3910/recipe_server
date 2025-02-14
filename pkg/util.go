package util

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

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

func GetBody[T interface{}](r *http.Request, bodyStruct *T) error {
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
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
	Port  string
	DbUrl string
}

var loadedEnv *envs

func LoadEnvs() *envs {
	godotenv.Load()

	if loadedEnv != nil {
		return loadedEnv
	}
	port := os.Getenv("PORT")
	dbUrl := os.Getenv("DATABASE_URL")
	if port == "" {
		log.Fatal("Could not read port from .env file")
	}
	if dbUrl == "" {
		log.Fatal("Could not read dbUrl from .env file")
	}
	return &envs{
		Port:  port,
		DbUrl: dbUrl,
	}
}

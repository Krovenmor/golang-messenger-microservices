package config

import (
	"errors"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func Init() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("config::Init(): Trouble with loading env file/files")
	}
}

func GetEnvVar(key string) (string, error) {
	val, isExists := os.LookupEnv(key)
	if !isExists {
		log.Printf("config::GetEnvVar(%q): Not found", key)
		return val, errors.New("Not found")
	}
	return val, nil
}

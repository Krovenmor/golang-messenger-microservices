package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

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
		return val, fmt.Errorf("EnvVar: %q not found", key)
	}
	return val, nil
}

func GetEnvVarBytes(key string) ([]byte, error) {
	val, err := GetEnvVar(key)
	if err != nil {
		return []byte{}, err
	}
	keyAr, err := os.ReadFile(val)
	if err != nil {
		return []byte{}, err
	}
	return keyAr, nil
}

func GetEnvVarInt(key string) (int, error) {
	val, err := GetEnvVar(key)
	if err != nil {
		return -1, err
	}
	iVal, err := strconv.Atoi(val)
	if err != nil {
		return -1, err
	}
	return iVal, nil
}

func GetEnvVarDuration(key string) (time.Duration, error) {
	val, err := GetEnvVar(key)
	if err != nil {
		return -1, err
	}
	dur, err := time.ParseDuration(val)
	if err != nil {
		return -1, err
	}
	return dur, nil
}

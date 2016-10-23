package config

import (
	"errors"
	"github.com/joho/godotenv"
	"os"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		panic("System Init Error: could not load .env!")
	}
}

func Load(path string) error {
	err := godotenv.Load(path)
	if err != nil {
		return errors.New("could not load .env: " + err.Error())
	}
	return nil
}

func Env(key string) string {
	return os.Getenv(key)
}

package config

import (
	"errors"
	"github.com/joho/godotenv"
	"os"
)

func Load(path string) error {
	err := godotenv.Load(path)
	if err != nil {
		return errors.New("could not load " + path + ": " + err.Error())
	}
	return nil
}

func Env(key string) string {
	return os.Getenv(key)
}

package config

import (
	"errors"
	"github.com/joho/godotenv"
	"golang.org/x/net/context"
	"log"
	"os"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		panic("System Init Error: could not load .env!")
	}
}

type ConfServer struct{}

func (serv *ConfServer) Env(ctx context.Context, in *ConfigNameRequest) (*ConfigValueReply, error) {
	reply := &ConfigValueReply{Value: os.Getenv(in.Name)}
	log.Printf("Get Config For %s Result: %s\n", in.Name, reply.Value)
	if reply.Value == "" {
		return reply, errors.New("Config Empty")
	}
	return reply, nil
}

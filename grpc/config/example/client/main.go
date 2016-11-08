package main

import (
	"github.com/yinhui87/go-component/config"
	gConfig "github.com/yinhui87/go-component/grpc/config"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
)

func main() {
	conn, err := grpc.Dial("localhost"+config.Env("PORT"), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s\n", err.Error())
	}
	log.Print("Dial Success")
	defer conn.Close()
	client := gConfig.NewConfigClient(conn)
	log.Println("NewConfigClient Success")
	key := "TEST"
	ret, err := client.Env(context.Background(), &gConfig.ConfigNameRequest{Name: key})
	if err != nil {
		log.Printf("gEnv Error: %s\n", err.Error())
	} else {
		log.Printf("gEnv Success: %s\n", ret.Value)
	}
}

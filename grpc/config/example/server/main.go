package main

import (
	"github.com/yinhui87/go-component/config"
	gConfig "github.com/yinhui87/go-component/grpc/config"

	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	lis, err := net.Listen("tcp", config.Env("PORT"))
	if err != nil {
		log.Fatalf("failed to listen: %v\n", err)
	}
	s := grpc.NewServer()
	gConfig.RegisterConfigServer(s, &gConfig.ConfServer{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v\n", err)
	}
}

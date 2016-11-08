package main

import (
	"github.com/yinhui87/go-component/config"
	"github.com/yinhui87/go-component/grpc/incrementid"
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
	incrementid.RegisterIncrementIdServer(s, &incrementid.IncrServer{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v\n", err)
	}
}

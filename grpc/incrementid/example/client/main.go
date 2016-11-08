package main

import (
	"github.com/yinhui87/go-component/config"
	"github.com/yinhui87/go-component/incrementid"
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
	client := incrementid.NewIncrementIdClient(conn)
	log.Print("NewIncrementIdClient Success")
	// 先测试CheckKeyExist功能
	key := "test"
	initial_value := uint64(7)
	checkRet, err := client.CheckIncrKeyExist(context.Background(), &incrementid.GetIncrIdRequest{Name: key})
	if err != nil {
		log.Fatalf("CheckIncrKeyExist Error: %s\n", err.Error())
	}
	log.Printf("key test CheckExist Rerult: %v\n", checkRet.Exist)
	if !checkRet.Exist {
		_, err = client.CreateIncrKey(context.Background(),
			&incrementid.CreateIncrKeyRequest{Name: key, InitialValue: initial_value})
		if err != nil {
			log.Fatalf("CreateIncrKeyRequest Error: %s\n", err.Error())
		}
		log.Print("CreateIncrKeyRequest Success")
	}
	incrRet, err := client.GetIncrId(context.Background(), &incrementid.GetIncrIdRequest{Name: key})
	if err != nil {
		log.Fatalf("GetIncrId Error: %s\n", err.Error())
	}
	log.Printf("GetIncrId Result: %v\n", incrRet.IncId)
}

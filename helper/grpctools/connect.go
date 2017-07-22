package grpctools

import (
	"errors"
	"github.com/silenceper/pool"
	"google.golang.org/grpc"
	"strings"
	"time"
)

func getGrpcConnectWithInsecure(host, port string) (conn *grpc.ClientConn, err error) {
	if !strings.HasPrefix(port, ":") {
		port = ":" + port
	}
	conn, err = grpc.Dial(host+port, grpc.WithInsecure())
	if err != nil {
		return nil, errors.New("grpc.Dial error: " + err.Error())
	}
	return conn, nil
}

func closeGrpcConn(v interface{}) error {
	return v.(*grpc.ClientConn).Close()
}

func GetGrpcConnPool(host, port string, initCap, maxCap int) (p pool.Pool, err error) {
	factory := func() (interface{}, error) { return getGrpcConnectWithInsecure(host, port) }
	close := closeGrpcConn

	//创建一个连接池： 初始化5，最大链接30
	poolConfig := &pool.PoolConfig{
		InitialCap: initCap,
		MaxCap:     maxCap,
		Factory:    factory,
		Close:      close,
		//链接最大空闲时间，超过该时间的链接 将会关闭，可避免空闲时链接EOF，自动失效的问题
		IdleTimeout: 15 * time.Second,
	}
	p, err = pool.NewChannelPool(poolConfig)
	if err != nil {
		return nil, errors.New("pool.NewChannelPool error: " + err.Error())
	}
	return p, nil
}

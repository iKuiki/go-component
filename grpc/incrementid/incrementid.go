package incrementid

import (
	"errors"
	"golang.org/x/net/context"
	"log"
	"sync/atomic"
)

var increments map[string]*uint64

func init() {
	increments = make(map[string]*uint64)
}

type IncrServer struct{}

func (serv *IncrServer) GetIncrId(ctx context.Context, in *IncrIdNameRequest) (*IncrIdReply, error) {
	return serv.GetIncrIdByAmount(ctx, &IncrIdNameValueRequest{Name: in.Name, Value: 1})
}

func (serv *IncrServer) GetIncrIdByAmount(ctx context.Context, in *IncrIdNameValueRequest) (*IncrIdReply, error) {
	v, ok := increments[in.Name]
	if !ok {
		return &IncrIdReply{}, errors.New("Key Not Exist")
	}
	nId := atomic.AddUint64(v, in.Value)
	log.Printf("Get Increment Id for %s with amount %d, new Id is: %d\n", in.Name, in.Value, nId)
	return &IncrIdReply{Id: nId}, nil
}

func (serv *IncrServer) CheckIncrKeyExist(ctx context.Context, in *IncrIdNameRequest) (*IncrBoolReply, error) {
	_, ok := increments[in.Name]
	if !ok {
		return &IncrBoolReply{Ret: false}, nil
	}
	return &IncrBoolReply{Ret: true}, nil
}

func (serv *IncrServer) CreateIncrKey(ctx context.Context, in *IncrIdNameValueRequest) (*NoneReply, error) {
	// 其实此方法应加读写锁，前2个方法属于读方法，此方法属于写方法
	_, ok := increments[in.Name]
	if ok {
		return &NoneReply{}, errors.New("Key Already Exist")
	}
	nId := in.Value
	increments[in.Name] = &nId
	log.Printf("Create Increment Key %s Success, Value: %d\n", in.Name, in.Value)
	return &NoneReply{}, nil
}

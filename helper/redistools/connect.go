package redistools

import (
	"errors"
	"github.com/go-redis/redis"
	"github.com/silenceper/pool"
	"time"
)

func getRedisConn(redisinfo RedisInfo) (client *redis.Client, err error) {
	client = redis.NewClient(&redis.Options{
		Addr:     redisinfo.Host + ":" + redisinfo.Port,
		Password: redisinfo.Password,
		DB:       redisinfo.Dbno,
	})
	_, err = client.Ping().Result()
	if err != nil {
		return nil, err
	}
	return client, nil
}

func closeRedisConn(v interface{}) error {
	return v.(*redis.Client).Close()
}

func GetRedisConnPool(redisinfo RedisInfo, initCap, maxCap int) (p pool.Pool, err error) {
	factory := func() (interface{}, error) { return getRedisConn(redisinfo) }
	close := closeRedisConn

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

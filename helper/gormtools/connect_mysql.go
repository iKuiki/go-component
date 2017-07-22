package gormtools

import (
	"errors"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/silenceper/pool"
	"time"
)

func getMysqlConn(dbinfo DbInfo) (db *gorm.DB, err error) {
	db, err = gorm.Open(
		"mysql",
		dbinfo.User+":"+dbinfo.Password+"@tcp("+dbinfo.Host+":"+dbinfo.Port+")/"+dbinfo.Dbname+
			"?charset=utf8mb4,utf8&collation=utf8mb4_unicode_ci&parseTime=True&loc=Local")
	return db, err
}

func closeMysqlConn(v interface{}) error {
	return v.(*gorm.DB).Close()
}

func GetMysqlConnPool(dbinfo DbInfo, initCap, maxCap int) (p pool.Pool, err error) {
	factory := func() (interface{}, error) { return getMysqlConn(dbinfo) }
	close := closeMysqlConn

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

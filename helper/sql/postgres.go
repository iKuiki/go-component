package sql

import (
	"errors"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func PostgresInit(dbName string, dbinfo *DbInfo) error {
	var err error
	// 配置数据库
	err = orm.RegisterDataBase(dbName, "postgres", "postgres://"+dbinfo.User+":"+dbinfo.Password+"@"+dbinfo.Host+"/"+dbinfo.Dbname+"?port="+dbinfo.Port)
	if err != nil {
		return errors.New("Database Open Fail: " + err.Error())
	}
	if err := checkConnect(dbName); err != nil {
		return err
	}
	return nil
}

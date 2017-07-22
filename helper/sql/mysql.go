package sql

import (
	"errors"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func MysqlInit(dbName string, dbinfo *DbInfo) error {
	var err error
	// 配置数据库
	err = orm.RegisterDataBase(dbName, "mysql", dbinfo.User+":"+dbinfo.Password+"@tcp("+dbinfo.Host+":"+dbinfo.Port+")/"+dbinfo.Dbname+"?charset=utf8mb4,utf8&collation=utf8mb4_unicode_ci&interpolateParams=true")
	// DB, err = sql.Open("mysql", dbuser+":"+dbpassword+"@tcp("+dbhost+":"+dbport+")/"+dbname+"?charset=utf8")
	if err != nil {
		return errors.New("Database Open Fail: " + err.Error())
	}
	if err := checkConnect(dbName); err != nil {
		return err
	}
	return nil
}

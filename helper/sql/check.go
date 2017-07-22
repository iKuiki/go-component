package sql

import (
	"errors"
	"github.com/astaxie/beego/orm"
)

func checkConnect(dbName string) error {

	DB, err := orm.GetDB(dbName)
	if err != nil {
		return errors.New("Database server not ready: " + err.Error())
	}
	err = DB.Ping()
	if err != nil {
		return errors.New("Database ping Fail: " + err.Error())
	}
	return nil
}

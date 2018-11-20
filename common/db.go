package common

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

func OpenDb() error {
	var err error
	DB, err = xorm.NewEngine("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", Config.Db.User, Config.Db.PassWord, Config.Db.Host, Config.Db.Port, Config.Db.Name))
	if err!=nil{
		panic(err)
	}
	err=DB.Ping()
	if err!=nil{
		panic(err)
	}
	fmt.Println("err:", err)
	DB.ShowSQL(Config.Db.ShowSQL)
	DB.ShowExecTime(Config.Db.ShowSQL)
	return err
}

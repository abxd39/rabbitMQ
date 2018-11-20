package common

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

func OpenDb() error {
	var err error
	str:= fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", Config.Db.User, Config.Db.PassWord, Config.Db.Host, Config.Db.Port, Config.Db.Name)
	fmt.Println("数据库链接",str)
	DB, err = xorm.NewEngine("mysql",str)
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

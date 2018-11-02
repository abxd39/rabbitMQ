package common

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
	"github.com/pkg/errors"
)

func OpenDb() error {
	var err error
	master, err := xorm.NewEngine("mysql",
		fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
			Config.Db.User,
			Config.Db.Password,
			Config.Db.Host,
			Config.Db.Port,
			Config.Db.Name))
	if err != nil {
		return err
	}
	slaves := make([]*xorm.Engine, len(Config.Db.Slaves))
	for i, salve := range Config.Db.Slaves {
		slaves[i], err = xorm.NewEngine("mysql",
			fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
				Config.Db.SlaveConfig.User,
				Config.Db.SlaveConfig.Password,
				salve.Host,
				salve.Port,
				salve.Name))
		if err != nil {
			return err
		}
	}
	if len(slaves) <= 0 {
		slaves = append(slaves, master)
	}
	DB, err = xorm.NewEngineGroup(master, slaves)
	DB.SetMaxOpenConns(Config.Db.MaxOpenConns)
	DB.ShowSQL(Config.Db.ShowSQL)
	DB.ShowExecTime(Config.Db.ShowSQL)
	DB.SetLogger(&XOrmLogger{logger: Log})
	return err
}

type XOrmLogger struct {
	logger  *Logger
	level   core.LogLevel
	showSQL bool
}

// Error implement core.ILogger
func (x *XOrmLogger) Error(v ...interface{}) {
	if x.level <= core.LOG_ERR {
		Log.Errorln(v...)
	}
	return
}

// Errorf implement core.ILogger
func (x *XOrmLogger) Errorf(format string, v ...interface{}) {
	if x.level <= core.LOG_ERR {
		Log.Errorf(format, v...)
	}
	return
}

// Debug implement core.ILogger
func (x *XOrmLogger) Debug(v ...interface{}) {
	if x.level <= core.LOG_DEBUG {
		Log.Traceln(v...)
	}
	return
}

// Debugf implement core.ILogger
func (x *XOrmLogger) Debugf(format string, v ...interface{}) {
	if x.level <= core.LOG_DEBUG {
		Log.Tracef(format, v...)
	}
	return
}

// Info implement core.ILogger
func (x *XOrmLogger) Info(v ...interface{}) {
	if x.level <= core.LOG_INFO {
		Log.Infoln(v...)
	}
	return
}

// Infof implement core.ILogger
func (x *XOrmLogger) Infof(format string, v ...interface{}) {
	if x.level <= core.LOG_INFO {
		Log.Infof(format, v...)
	}
	return
}

// Warn implement core.ILogger
func (x *XOrmLogger) Warn(v ...interface{}) {
	if x.level <= core.LOG_WARNING {
		Log.Warnln(v...)
	}
	return
}

// Warnf implement core.ILogger
func (x *XOrmLogger) Warnf(format string, v ...interface{}) {
	if x.level <= core.LOG_WARNING {
		Log.Warnf(format, v...)
	}
	return
}

// Level implement core.ILogger
func (x *XOrmLogger) Level() core.LogLevel {
	return x.level
}

// SetLevel implement core.ILogger
func (x *XOrmLogger) SetLevel(l core.LogLevel) {
	x.level = l
	return
}

// ShowSQL implement core.ILogger
func (x *XOrmLogger) ShowSQL(show ...bool) {
	if len(show) == 0 {
		x.showSQL = true
		return
	}
	x.showSQL = show[0]
}

// IsShowSQL implement core.ILogger
func (x *XOrmLogger) IsShowSQL() bool {
	return x.showSQL
}

// 新增批量事务处理

//执行方法
type Action func(*xorm.Session) (interface{}, error)

//1. 匿名函数的作用域内可以访问上级函数的变量无需通过参数传递
//2. 返回参数必须是命名的才能获取到错误信息
func Transaction(action Action) (obj interface{}, err error) {

	if action == nil {
		return nil, errors.New("action can not be nil")
	}

	obj = nil

	//开启事务
	tran := DB.NewSession()
	defer tran.Close()

	err = tran.Begin()
	if err != nil {
		tran.Rollback()
		fmt.Println()
		return nil, err
	}

	//如果action发生panic错误时执行
	defer func() {
		if r := recover(); r != nil {
			tran.Rollback()
			switch x := r.(type) {
			case string:
				err = errors.New(x)
			case error:
				err = x
			default:
				err = errors.New("unknown panic")
			}
		}
	}()

	if obj, err = action(tran); err != nil {
		tran.Rollback()
		//方法执行失败
		Log.Errorln("方法执行失败" + err.Error())
		return nil, err
	}

	tran.Commit()

	return obj, nil
}

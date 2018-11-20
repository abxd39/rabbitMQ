package common

import (
	"github.com/Unknwon/goconfig"
	"github.com/koding/multiconfig"
)


var CPath *goconfig.ConfigFile

func init() {
	var err error
	CPath, err = goconfig.LoadConfigFile("E:/WorkSpace/src/sctek.com/typhoon/th-platform-gateway/configPath.ini")
	if err != nil {
		panic("load config err is " + err.Error())
	}
}

//type FlagConfig struct {
//	ConfigFile string `default:"config.json"`
//}

type LoggerConfig struct {
	Enabled    bool `default:"true"`
	LogFile    string
	TraceLevel int `default:"3"`
}

type ServerConfig struct {
	Listen       string `default:":5000"`
	RuntimePath  string `default:"runtime"`
	MaxWork      int    `json:"max_work"`
	MaxQueueSize int    `json:"max_queue_size"`
	Db           struct {
		Host        string
		Port        string `default:"3306"`
		Name        string
		User        string
		Password    string
		SlaveConfig struct {
			User     string
			Password string
		}
		Slaves []struct {
			Host string
			Port string
			Name string
		}
		ShowSQL      bool `default:"false"`
		MaxOpenConns int  `default:"100"`
	}
	Redis struct {
		Address  string
		Database int `default:"0"`
		Password string
	}
	Log struct {
		LogFile    string `default:""`
		TraceLevel int    `default:"3"`
		Logger     struct {
			Trace LoggerConfig
			Info  LoggerConfig
			Warn  LoggerConfig
			Error LoggerConfig
		}
	}
	GateWay struct {
		Token string `default:"test"`
		Host  string
		Test  bool `default:"false"`
	}
	VisitInterval string
	IsDev         bool   `default:"false"`
	Url           string `json:"url"`
}


func (c *ServerConfig) load() error {
	t := &multiconfig.TagLoader{}
	j := &multiconfig.JSONLoader{Path:CPath.MustValue("dbPath","path","config.json") }
	m := multiconfig.MultiLoader(t, j)
	err := m.Load(c)
	return err
}

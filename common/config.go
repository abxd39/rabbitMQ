package common

import (
	"github.com/Unknwon/goconfig"
	"github.com/koding/multiconfig"
)

var CPath *goconfig.ConfigFile

func init() {
	var err error
	CPath, err = goconfig.LoadConfigFile("./configPath.ini")
	if err != nil {
		panic("load config err is " + err.Error())
	}
}

type ServerConfig struct {
	Listen       string `json:"listen" default:":5000"`
	RuntimePath  string `json:"runtime_path" default:"runtime"`
	MaxWork      int    `json:"max_work"`
	MaxQueueSize int    `json:"max_queue_size"`
	JPush struct{
		 AppKey string `json:"app_key"`
		 MasterSecret string `json:"master_secret"`
	}
	Db           struct {
		Host        string `json:"host"`
		Port        string `json:"port" default:"3306"`
		Name        string `json:"name"`
		User        string `json:"user"`
		PassWord    string `json:"pass_word"`
		SlaveConfig struct {
			User     string `json:"user"`
			PassWord string `json:"pass_word"`
		}
		Slaves []struct {
			Host string `json:"host"`
			Port string `json:"port"`
			Name string `json:"name"`
		}
		ShowSQL      bool `json:"show_sql" default:"false"`
		MaxOpenConns int  `json:"max_open_conns" default:"100"`
	}
	Redis struct {
		Address  string `json:"address"`
		Database int    `json:"database" default:"0"`
		PassWord string `json:"pass_word"`
		Key string `json:"key"`
	}
	Log struct {
		LogFileName string `json:"log_file_name"`
		LogFileDir  string `json:"log_file_dir"`
		//TraceLevel  string `json:"trace_level" default:"info"`
	}
	GateWay struct {
		Token string `default:"test"`
		Host  string
		Test  bool `default:"false"`
	}
	VisitInterval string
	IsDev         bool   `default:"false"`
	Url           string `json:"url"`
	PollingTime   string `json:"polling_time"`
}

func (c *ServerConfig) load() error {
	t := &multiconfig.TagLoader{}
	j := &multiconfig.JSONLoader{Path: CPath.MustValue("dbPath", "path", "config.json")}
	m := multiconfig.MultiLoader(t, j)
	err := m.Load(c)
	return err
}

package module

import (
	"time"
)

type MonitorLog struct {
	Id            int       `xorm:"not null pk autoincr unique INT(11)"`
	CorpId        int       `xorm:"not null INT(11)"`
	MallId        int       `xorm:"not null index INT(11)"`
	ShopId        int       `xorm:"not null comment('商家id') index INT(11)"`
	DeviceId      int       `xorm:"not null comment('设备id') index INT(11)"`
	MallName      string    `xorm:"not null comment('商圈简称') VARCHAR(45)"`
	ShopName      string    `xorm:"not null comment('商店简称') VARCHAR(30)"`
	DeviceSerials string    `xorm:"not null comment('设备序列号') VARCHAR(50)"`
	DeviceType    int       `xorm:"not null comment('设备类型(1:插件 2:硬件)') TINYINT(4)"`
	Event         int       `xorm:"not null comment('事件,1:设备断网超过10分钟 2:设备故障') INT(11)"`
	Content       string    `xorm:"not null comment('报警内容') VARCHAR(255)"`
	Created       time.Time `xorm:"not null comment('创建时间') DATETIME"`
}

func (a MonitorLog) TableName() string {
	return "monitor_log"
}

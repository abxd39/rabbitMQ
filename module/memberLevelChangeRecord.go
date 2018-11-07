package module

import (
	"time"
)

type MemberLevelChangeRecord struct {
	Id              int       `xorm:"not null pk autoincr INT(11)"`
	CorpId          int       `xorm:"not null default 0 comment('机构ID') INT(11)"`
	MallId          int       `xorm:"not null default 0 comment('分店id') INT(11)"`
	MemberId        int       `xorm:"not null default 0 comment('会员ID') INT(11)"`
	UserId          int       `xorm:"not null default 0 comment('操作人ID') INT(11)"`
	UserName        string    `xorm:"not null default '' comment('操作人名字') VARCHAR(255)"`
	BeforeLevel     int       `xorm:"not null default 0 comment('调整前的级别') INT(11)"`
	BeforeLevelName string    `xorm:"not null default '' comment('调整前等级名称') VARCHAR(255)"`
	AfterLevel      int       `xorm:"not null default 0 comment('调整后的级别') INT(11)"`
	AfterLevelName  string    `xorm:"not null default '' comment('调整后等级名称') VARCHAR(255)"`
	Created         time.Time `xorm:"DATETIME"`
}

func (a MemberLevelChangeRecord) TableName() string {
	return "member_level_change_record"
}

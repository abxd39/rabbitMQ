package module

import (
	"time"
)

type MemberRight struct {
	Id           int       `xorm:"not null pk autoincr INT(11)"`
	CorpId       int       `xorm:"not null default 0 INT(11)"`
	Name         string    `xorm:"not null default '' comment('权益名称') VARCHAR(255)"`
	Desc         string    `xorm:"not null comment('权益详情') TEXT"`
	MemberLevels string    `xorm:"not null default '' comment('享有此权益的会员等级列表') VARCHAR(255)"`
	Created      time.Time `xorm:"not null DATETIME"`
	Updated      time.Time `xorm:"not null DATETIME"`
	IsDeleted    int       `xorm:"not null default 0 comment('0：正常，1：删除') TINYINT(1)"`
}

func (a MemberRight) TableName() string {
	return "member_right"
}

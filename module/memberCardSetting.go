package module

import (
	"time"
)

type MemberCardSetting struct {
	Id       int       `xorm:"not null pk autoincr INT(11)"`
	CorpId   int       `xorm:"not null default 0 comment('机构ID') INT(11)"`
	Notice   string    `xorm:"not null comment('须知') TEXT"`
	Telphone string    `xorm:"not null default '' comment('电话') VARCHAR(50)"`
	Status   int       `xorm:"not null default 0 comment('状态（0：关闭：1：开启）') TINYINT(1)"`
	Created  time.Time `xorm:"not null DATETIME"`
}

func (a MemberCardSetting) TableName() string {
	return "member_card_setting"
}

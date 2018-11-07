package module

import (
	"time"
)

type MemberInfo struct {
	MemberId int       `xorm:"not null pk comment('会员ID') INT(11)"`
	CorpId   int       `xorm:"not null default 0 INT(11)"`
	CardNo   string    `xorm:"not null default '' comment('会员卡号') index VARCHAR(50)"`
	Name     string    `xorm:"not null default '' comment('会员姓名') VARCHAR(255)"`
	Mobile   string    `xorm:"not null default '' comment(' 手机号') index VARCHAR(20)"`
	Sex      int       `xorm:"not null default 0 comment('性别（0：未知，1：男，2：女）') TINYINT(1)"`
	Birthday time.Time `xorm:"not null comment(' 生日') DATE"`
	ProvId   int       `xorm:"not null default 0 INT(11)"`
	CityId   int       `xorm:"not null default 0 INT(11)"`
	AreaId   int       `xorm:"not null default 0 INT(11)"`
	Address  string    `xorm:"not null comment('详细地址') TEXT"`
	Email    string    `xorm:"not null default '' comment('邮箱') VARCHAR(50)"`
}

func (a MemberInfo) TableName() string {
	return "member_info"
}

package module

import (
	"time"
)

type MemberPointsTactics struct {
	Id            int       `xorm:"not null pk autoincr INT(11)"`
	CorpId        int       `xorm:"not null default 0 comment('机构id') INT(11)"`
	Title         string    `xorm:"not null default '' comment('策略名称') VARCHAR(100)"`
	Type          int       `xorm:"not null default 0 comment('类型：1,通用 2,业态') TINYINT(4)"`
	EffectiveDate time.Time `xorm:"not null comment('生效日期') DATE"`
	IndustryId    int       `xorm:"not null default 0 comment('行业id') INT(11)"`
	PaymentFee    int       `xorm:"not null default 0 comment('消费金额(单位为分)') INT(11)"`
	Points        int       `xorm:"not null default 0 comment('积分') INT(11)"`
	UserId        int       `xorm:"not null default 0 comment('操作人id') INT(11)"`
	PointsLimit   int       `xorm:"not null default 0 comment('积分上限') INT(11)"`
	Status        int       `xorm:"not null default 0 comment('状态：0,待启用 1,开启 2.关闭') TINYINT(4)"`
	Created       time.Time `xorm:"not null comment('创建时间') DATETIME"`
	Updated       time.Time `xorm:"not null comment('更新时间') DATETIME"`
	MallId        int       `xorm:"not null default 0 comment('分店ID') INT(11)"`
}

func (a MemberPointsTactics) TableName() string {
	return "member_points_tactics"
}

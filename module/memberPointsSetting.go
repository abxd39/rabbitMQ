package module

import (
	"time"
)

type MemberPointsSetting struct {
	Id                int       `xorm:"not null pk autoincr INT(11)"`
	CorpId            int       `xorm:"not null INT(11)"`
	PointsLimit       int       `xorm:"not null default 0 comment('每日总积分上限') INT(11)"`
	ScanValidity      int       `xorm:"not null default 0 comment('扫描有效期') TINYINT(4)"`
	RefundValidity    int       `xorm:"not null default 0 comment('可退货时间') TINYINT(4)"`
	RefundPointsLimit int       `xorm:"not null default 0 comment('退积分限制（时间戳）') INT(11)"`
	SweepPoints       int       `xorm:"not null default 1 comment('扫码积分，0关闭，1开启') TINYINT(4)"`
	PointsShop        int       `xorm:"not null default 1 comment('积分商城，0为关闭，1为开启') TINYINT(4)"`
	Type              int       `xorm:"not null default 0 comment('类型：1,定期清零 2,滚动清零') TINYINT(4)"`
	ResetCondition    string    `xorm:"not null default '' comment('积分清零条件') VARCHAR(255)"`
	Created           time.Time `xorm:"not null comment('创建时间') DATETIME"`
	Updated           time.Time `xorm:"not null comment('更新时间') DATETIME"`
}

func (a MemberPointsSetting) TableName() string {
	return "member_points_setting"
}

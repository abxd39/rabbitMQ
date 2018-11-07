package module

import (
	"time"
)

type MemberPointsDetail struct {
	Id                int       `xorm:"not null pk autoincr INT(11)"`
	CorpId            int       `xorm:"not null INT(11)"`
	MallId            int       `xorm:"not null default 0 comment('分店id') INT(11)"`
	ShopId            int       `xorm:"not null default 0 comment('商家id') INT(11)"`
	MemberId          int       `xorm:"not null default 0 comment('会员id') index INT(11)"`
	MemberLevel       int       `xorm:"not null default 1 comment('会员等级') TINYINT(4)"`
	PaymentFee        int       `xorm:"not null default 0 comment('消费金额(单位为分)') INT(11)"`
	Code              string    `xorm:"not null default '' comment('订单号') VARCHAR(255)"`
	Points            int       `xorm:"not null default 0 comment('积分') INT(11)"`
	Type              int       `xorm:"comment('积分类型：1:赠送 2:消耗 3:过期 4:退款退积分') TINYINT(4)"`
	Status            int       `xorm:"not null default 0 comment('是否已解冻，0：否  1：是') TINYINT(4)"`
	RelateId          string    `xorm:"not null default '' comment('关联id') VARCHAR(255)"`
	RelateType        int       `xorm:"not null comment('关联类型：1：小票码 2：数方支付 3：活动 ') TINYINT(4)"`
	ReturnPointStatus int       `xorm:"not null default 0 comment('积分退还状态（0：待申请，1：待审核，2：已完成，3：已拒绝）') TINYINT(1)"`
	OrderTime         time.Time `xorm:"not null comment('下单时间') DATETIME"`
	Sdate             time.Time `xorm:"not null comment('创建日期') DATE"`
	Created           time.Time `xorm:"not null comment('创建时间') index DATETIME"`
}

func (a MemberPointsDetail) TableName() string {
	return "member_points_detail"
}

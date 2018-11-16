package module

import (
	"sctek.com/typhoon/th-platform-gateway/common"
	"time"
)

type Member struct {
	Id                int       `xorm:"not null pk autoincr INT(11)"`
	CorpId            int       `xorm:"not null default 0 comment('机构ID') INT(11)"`
	MallId            int       `xorm:"not null default 0 comment('分店id') INT(11)"`
	ShopId            int       `xorm:"not null INT(11)"`
	MemberCardId      int       `xorm:"not null default 0 comment('会员卡ID') INT(11)"`
	TotalFee          int64     `xorm:"not null default 0 comment('消费金额（单位为分）') BIGINT(20)"`
	LevelFee          int64     `xorm:"not null default 0 comment('当前等级累计消费') BIGINT(20)"`
	TotalPoints       int64     `xorm:"not null default 0 comment('会员总积分') BIGINT(20)"`
	UsablePoints      int64     `xorm:"not null default 0 comment('可用积分') BIGINT(20)"`
	FrozenPoints      int64     `xorm:"not null default 0 comment('冻结积分') BIGINT(20)"`
	PayTimes          int       `xorm:"not null default 0 comment('消费次数') INT(11)"`
	Status            int       `xorm:"not null default 1 comment('状态（0：冻结 ,1：正常）') TINYINT(4)"`
	YAmount           int       `xorm:"not null default 0 comment('截至昨日的累计消费（分）') INT(11)"`
	YCount            int       `xorm:"not null default 0 comment('截至到昨日的累计消费次数') INT(11)"`
	YStatisDate       time.Time `xorm:"not null comment('统计的时间（避免重复统计）') DATETIME"`
	IsNew             int       `xorm:"not null default 1 comment('是否为新会员（0：否，1：是）') TINYINT(4)"`
	FirstConsumeDate  time.Time `xorm:"not null comment('首次消费时间') DATETIME"`
	ToOldDate         time.Time `xorm:"not null comment('转为老会员时间') DATETIME"`
	LastVisitTime     time.Time `xorm:"not null comment('最后访问时间') DATETIME"`
	LastDownlevelTime time.Time `xorm:"not null comment('上次降级时间') DATETIME"`
	Created           time.Time `xorm:"not null comment('创建时间') DATETIME"`
}

func (m *Member) TableName() string {
	return "member"
}

func (m *Member) GetMallId(id int) int {
	engine := common.DB
	has, err := engine.Where("id=?", id).Get(m)
	if err != nil {
		common.Log.Errorln(err)
		return 0
	}
	if !has {
		common.Log.Errorf("会员%d不存在！！\r\n", id)
		return 0
	}
	return m.MallId
}

func (m *Member) GetMemberId(idList []int) ([]int, error) {
	engine := common.DB
	list := make([]int, 0)
	err := engine.Table("member").Cols("id").In("member_card_id", idList).Find(&list)
	if err != nil {
		return nil, err
	}
	return list, nil
}

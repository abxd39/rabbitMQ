package module

import (
	"sctek.com/typhoon/th-platform-gateway/common"
	"time"
)

type MemberCard struct {
	Id                int       `xorm:"not null pk autoincr INT(11)"`
	CorpId            int       `xorm:"not null default 0 comment('机构ID') INT(11)"`
	Level             int       `xorm:"not null default 0 comment('会员卡等级') TINYINT(4)"`
	LevelName         string    `xorm:"not null default '' comment('等级名称') VARCHAR(255)"`
	CoverImg          string    `xorm:"not null default '' comment('会员卡卡面') VARCHAR(255)"`
	CoverImgWxUrl     string    `xorm:"default '' comment('微信cdn地址') VARCHAR(255)"`
	UpCondition       int64     `xorm:"not null default 0 comment('升至本等级需要消费的累计金额') BIGINT(20)"`
	DownConditionTime int       `xorm:"not null default 0 comment('自开卡时间起降级的时间纬度') INT(11)"`
	DownCondition     int64     `xorm:"not null default 0 comment('降级条件') BIGINT(20)"`
	IsPrimary         int       `xorm:"not null default 0 comment('主要的，不能删除的') TINYINT(4)"`
	Created           time.Time `xorm:"not null DATETIME"`
	Updated           time.Time `xorm:"not null DATETIME"`
	Deleted           int       `xorm:"not null default 0 comment('删除状态（0：正常，1：删除）') TINYINT(1)"`
	IsTimeDelete      int       `xorm:"not null default 0 comment('是否定时删除（0：否，1：是）') TINYINT(1)"`
}

func (m *MemberCard) TableName() string {
	return "member_card"
}

type mobile struct {
	Mobile   string `json:"mobile"`
	Level    int    `json:"level"`
	CardNo   string `json:"card_no"`
	MemberId int    `json:"member_id"`
	CorpId   int    `json:"corp_id"`
}

func (m *mobile) TableName() string {
	return "member_card"
}

//根据会员等级发送
//func (m *MemberCard) SendMessageForGrade(manageId int, grade, message string) error {
//	common.Log.Infoln("根据会员等级把消息压入mq队列")
//	engine := common.DB
//
//	query := engine.Join("left", "member_info", "card_no==id")
//	query = query.In("level", grade)
//
//	list := make([]mobile, 0)
//	err := query.Find(&list)
//	if err != nil {
//		common.Log.Infoln(err)
//		log.Print(err.Error())
//		return err
//	}
//
//	for _, v := range list {
//		sendLog := new(TemplateSmsLog)
//		sendLog.TemplateManageId = manageId
//		sendLog.MemberId = v.MemberId
//		sendLog.CorpId = v.CorpId
//		sendLog.Mobile = v.Mobile
//		sendLog.MallId = new(Member).GetMallId(v.MemberId)
//		if sendLog.MallId == 0 {
//			continue
//		}
//		result, err := sendLog.marshalJson(message)
//		if err != nil {
//			common.Log.Errorln(err)
//			continue
//		}
//		Push("myPusher", "rmq_test", result)
//	}
//	return nil
//}

func (m*MemberCard) GetMessageOfGrade(grade string)([]mobile,error){
	engine := common.DB

	query := engine.Join("left", "member_info", "card_no==id")
	query = query.In("level", grade)

	list := make([]mobile, 0)
	err := query.Find(&list)
	if err != nil {
		return nil ,err
	}
	return list, nil
}
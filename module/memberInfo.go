package module

import (
	"fmt"
	"sctek.com/typhoon/th-platform-gateway/common"
	"sctek.com/typhoon/th-platform-gateway/sms"
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

func (m* MemberInfo) TableName() string {
	return "member_info"
}


func (m* MemberInfo) SendMassageForSex(sex,message string)error{
	engine:=common.DB
	list:=make([]MemberInfo,0)
	fmt.Println("sex",sex,"message",message)
	err:=engine.Where("sex=?",sex).Find(&list)
	if err!=nil{
		common.Log.Errorln(err)
		return err
	}
	//发送短信
	for _,value:=range list{
		fmt.Println("phone=",value.Mobile)
		new(sms.SMSMessage).SendMobileMessage(value.Mobile,message)
	}

	return nil
}
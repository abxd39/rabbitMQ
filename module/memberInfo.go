package module

import (
	"fmt"
	"sctek.com/typhoon/th-platform-gateway/common"
	"sctek.com/typhoon/th-platform-gateway/manageMq"
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

func (m *MemberInfo) TableName() string {
	return "member_info"
}

func (m *MemberInfo) SendMessageForSex(sex, message string) error {
	common.Log.Infoln("根据性别发送短息")
	engine := common.DB
	list := make([]MemberInfo, 0)
	fmt.Println("sex", sex, "message", message)
	err := engine.Where("sex=?", sex).Find(&list)
	if err != nil {
		common.Log.Errorln(err)
		return err
	}
	//发送短信
	for _, value := range list {
		manageMq.ExampleLoggerOutput("phone=" + value.Mobile)
		message = fmt.Sprintf("{\"phone\":\"%q\",\"message\":\"%q\"}", value.Mobile, message)
		manageMq.GlobalMq.Publish("fanout", message)
		//new(sms.SMSMessage).SendMobileMessage(value.Mobile,message)
	}

	return nil
}

//根据会员生日
func (m *MemberInfo) SendMessageForBirthDay(date, message string) error {
	common.Log.Infoln("根据会员生日发送信息")
	phone := "15920038315"
	msg := "are you sure??"
	engine := common.DB
	list := make([]MemberInfo, 0)
	err := engine.Select("select * from member_info").Find(&list)
	if err != nil {
		common.Log.Errorln(err)
		manageMq.ExampleLoggerOutput(err.Error())
		return err
	}
	for _, v := range list {
		_ = v
		message = fmt.Sprintf("{\"phone\":\"%q\",\"message\":\"%q\"}", phone, msg)
		manageMq.GlobalMq.Publish("fanout", message)
	}

	return nil
}

//全员发送
func(m*MemberInfo)SendMessageEveryOne (message string)error{
	common.Log.Infoln("即时全员发送")
	engine:=common.DB
	list :=make([]MemberInfo,0)
	err:=engine.Select("select * from member_info").Find(&list)
	if err!=nil{
		common.Log.Errorln(err)
		manageMq.ExampleLoggerOutput(err.Error())
		return err
	}
	for _,v:=range list{
		message = fmt.Sprintf("{\"phone\":\"%q\",\"message\":\"%q\"}",v.Mobile,message)
		manageMq.GlobalMq.Publish("fanout",message)
	}
	return nil
}
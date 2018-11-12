package module

import (
	"fmt"
	"log"
	"sctek.com/typhoon/th-platform-gateway/common"
	"strings"
)

type MemberInfo struct {
	MemberId int    `xorm:"not null pk comment('会员ID') INT(11)"`
	CorpId   int    `xorm:"not null default 0 INT(11)"`
	CardNo   string `xorm:"not null default '' comment('会员卡号') index VARCHAR(50)"`
	Name     string `xorm:"not null default '' comment('会员姓名') VARCHAR(255)"`
	Mobile   string `xorm:"not null default '' comment(' 手机号') index VARCHAR(20)"`
	Sex      int    `xorm:"not null default 0 comment('性别（0：未知，1：男，2：女）') TINYINT(1)"`
	Birthday string `xorm:"not null comment(' 生日') DATE"`
	ProvId   int    `xorm:"not null default 0 INT(11)"`
	CityId   int    `xorm:"not null default 0 INT(11)"`
	AreaId   int    `xorm:"not null default 0 INT(11)"`
	Address  string `xorm:"not null comment('详细地址') TEXT"`
	Email    string `xorm:"not null default '' comment('邮箱') VARCHAR(50)"`
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
	for _, v := range list {

		result,err:=marshalJson(v.Mobile,message)
		if err!=nil{
			common.Log.Errorln(err)
			continue
		}
		Push("myPusher", "rmq_test", result)
	}

	return nil
}

//全员发送
func (m *MemberInfo) SendMessageEveryOne(message string) error {
	common.Log.Infoln("即时全员发送")
	engine := common.DB
	list := make([]MemberInfo, 0)
	err := engine.Select(" * ").Find(&list)
	if err != nil {
		common.Log.Warnln(err)
		log.Print(err.Error())
		return err
	}
	for _, v := range list {
		result,err:=marshalJson(v.Mobile,message)
		if err!=nil{
			common.Log.Errorln(err)
			continue
		}
		Push("myPusher", "rmq_test", result)
	}
	return nil
}

//按照会员生日发送
func (m *MemberInfo) SendMessageOfBirthDay(birthDat, message string) error {
	common.Log.Infoln("按照生日月分发送")
	subList := strings.Split(birthDat, ",")
	engine := common.DB
	list := make([]MemberInfo, 0)
	err := engine.Find(&list)
	if err != nil {
		common.Log.Errorln(err)
		return err
	}
	for _, v := range list {
		if len(v.Birthday) != len("1991-06-14") {
			continue
		}
		month := v.Birthday[5:7]
		for _, m := range subList {
			if strings.Compare(month, m) == 0 || strings.Compare(month, "0"+m) == 0 {
				result,err:=marshalJson(v.Mobile,message)
				if err !=nil{
					common.Log.Errorln(err)
					continue
				}
				Push("myPusher","rmq_test",result)
				continue
			}
		}
	}
	return nil
}

//指定电话号码发送短息
func (m *MemberInfo) SendMessageOfPhone(Phone, message string) error {
	common.Log.Infoln("指定电话号码发送短息")
	engine := common.DB
	has, err := engine.Where("mobile=?", Phone).Get(m)
	if err != nil {
		common.Log.Errorln(err)
		return err
	}
	if !has {
		err = fmt.Errorf("手机号码为：%v的用户不存在！！", Phone)
		common.Log.Errorln(err)
		return err
	}
	result,err:=marshalJson(Phone,message)
	if err!=nil{
		common.Log.Errorln(err)
		return err
	}
	Push("myPusher","rmq_test",result)
	return nil
}

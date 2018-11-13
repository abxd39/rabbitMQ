package module

import (
	"sctek.com/typhoon/th-platform-gateway/common"
)

type TemplateSmsType struct {
	Id               int    `xorm:"not null pk autoincr INT(10)"`
	Type             int    `xorm:"default 1 comment('1-性别；2-会员等级；3-会员生日月') TINYINT(1)"`
	TypeData         string `xorm:"default '' comment('type字段对应的填充数据') VARCHAR(25)"`
	TemplateManageId int    `xorm:"default 0 comment('template_sms_manage表的主键id') INT(10)" json:"template_manage_id"`
}

//指定会员时的条件筛选
//func (t *TemplateSmsType) SearchOfManageId(mId, templateMessageId int) error {
//	engine := common.DB
//	//获取短息模板
//	message, err := new(TemplateSms).GetText(templateMessageId)
//	if err != nil {
//		common.Log.Infoln(err)
//		return err
//	}
//	list := make([]TemplateSmsType, 0)
//	err = engine.Where("template_manage_id=?", mId).Find(&list)
//	if err != nil {
//		common.Log.Errorln(err)
//		return err
//	}
//
//	for _, v := range list {
//		fmt.Println("指定发送给谁",v.Type)
//		if v.Type == 1 { //姓别
//			new(MemberInfo).SendMessageForSex(mId,v.TypeData, message)
//		} else if v.Type == 2 { //会员等级
//			new(MemberCard).SendMessageForGrade(mId,v.TypeData,message)
//		} else if v.Type == 3 { //会员生日
//			new(MemberInfo).SendMessageOfBirthDay(mId,v.TypeData,message)
//		}
//	}
//	return nil
//}


func (t*TemplateSmsType)GetSmsTypeOfId(id int)([]TemplateSmsType,error){
	engine := common.DB
	list := make([]TemplateSmsType, 0)
	err := engine.Where("template_manage_id=?", id).Find(&list)
	if err != nil {
		return nil,err
	}
	return list,nil
}
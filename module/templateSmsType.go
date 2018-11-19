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


func (t*TemplateSmsType)GetSmsTypeOfId(id int)([]TemplateSmsType,error){
	engine := common.DB
	list := make([]TemplateSmsType, 0)
	err := engine.Where("template_manage_id=?", id).Find(&list)
	if err != nil {
		return nil,err
	}
	return list,nil
}
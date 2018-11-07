package module

import (
	"fmt"
	"log"
	"sctek.com/typhoon/th-platform-gateway/common"
	"time"
)

type TemplateSms struct {
	Id           int       `xorm:"not null pk autoincr INT(10)"`
	CorpId       int       `xorm:"default 0 INT(10)"`
	MallId       int       `xorm:"default 0 INT(10)"`
	TemplateType int       `xorm:"default 0 comment('模板类型：1-推广活动；2-商场优惠；3-会员关怀') TINYINT(1)"`
	TemplateName string    `xorm:"default '' comment('模板名称') VARCHAR(30)"`
	TemplateId   string    `xorm:"default '' comment('模板id') VARCHAR(15)"`
	Content      string    `xorm:"comment('模板内容') TEXT"`
	Created      time.Time `xorm:"comment('时间') DATETIME"`
}


func (t* TemplateSms) GetText(TemplateId int) (string,error){
	engine :=common.DB
	has,err:=engine.Where("template_type=?",TemplateId).Get(t)
	if err!=nil{
		log.Printf("短息模板查询失败=%s",err.Error())
		return "",err
	}
	if !has{
		return "",fmt.Errorf("短息模板id为=%d不存在",TemplateId)
	}
	return t.Content,nil
}
package module

import (
	"encoding/json"
	"sctek.com/typhoon/th-platform-gateway/common"
	"sctek.com/typhoon/th-platform-gateway/sms"
	"time"
)

type TemplateSmsLog struct {
	Id               int       `xorm:"not null pk autoincr INT(10)"`
	MemberId         int       `xorm:"default 0 INT(10)" json:"member_id"`
	CorpId           int       `xorm:"default 0 INT(10)" json:"corp_id"`
	MallId           int       `xorm:"default 0 INT(10)" json:"mall_id"`
	TemplateManageId int       `xorm:"default 0 comment('template_sms_manage表的主键id') INT(10)" json:"template_manage_id"`
	Mobile           string    `xorm:"not null default '' comment(' 手机号') VARCHAR(20)" json:"mobile"`
	Status           int       `xorm:"default 0 comment('发送状态：1-成功；2-失败') TINYINT(1)"`
	Created          time.Time `xorm:"comment('备注时间') DATETIME"`
}

func (t *TemplateSmsLog) marshalJson(message string) ([]byte, error) {
	body := make(map[string]interface{})
	body["mobile"] = t.Mobile
	body["message"] = message
	body["member_id"] = t.MemberId
	body["corp_id"] = t.CorpId
	body["mall_id"] = t.MallId
	body["template_manage_id"] = t.TemplateManageId
	return json.Marshal(body)
}

type Result struct {
	TemplateSmsLog `xorm:"extends"`
	Message        string `json:"message"`
}

//
func (t *TemplateSmsLog) SendMobileMessage(body []byte) {
	common.Log.Traceln("开始发送短息")
	re := &Result{}
	err := json.Unmarshal(body, re)
	if err != nil {
		common.Log.Errorln(err)
		return
	}
	err = new(sms.SMSMessage).SendMobileMessage(re.Mobile, re.Message)
	if err != nil {
		common.Log.Errorln(err)
		t.Status =2
	}
	t.Mobile = re.Mobile
	t.MemberId = re.MemberId
	t.CorpId = re.CorpId
	t.TemplateManageId =re.TemplateManageId
	t.MallId =re.MallId
	t.Status =1
	t.Created = time.Now()//.Format("2006-01-02 15:04:05")
	_,err=common.DB.InsertOne(t)
	if err!=nil{
		common.Log.Errorln(err)
	}
	return
}

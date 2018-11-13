package service

import (
	"encoding/json"
	"sctek.com/typhoon/th-platform-gateway/common"
	"sctek.com/typhoon/th-platform-gateway/module"
	"sctek.com/typhoon/th-platform-gateway/sms"
	"time"
)

type MarshalJson struct {
	MemberId         int    `json:"member_id"`
	CorpId           int    `json:"corp_id"`
	MallId           int    `json:"mall_id"`
	TemplateManageId int    `json:"template_manage_id"`
	Mobile           string `json:"mobile"`
	Msg              string `json:"msg"`
}

func (m *MarshalJson) marshalJson(message string) ([]byte, error) {
	body := make(map[string]interface{})
	body["mobile"] = m.Mobile
	body["msg"] = message
	body["member_id"] = m.MemberId
	body["corp_id"] = m.CorpId
	body["mall_id"] = m.MallId
	body["template_manage_id"] = m.TemplateManageId
	return json.Marshal(body)
}

func (m *MarshalJson) UnmarshalJson(body []byte) {
	common.Log.Traceln("开始发送短息")
	err := json.Unmarshal(body, m)
	if err != nil {
		common.Log.Errorln(err)
		return
	}
	if len(m.Mobile) <= 0 {
		common.Log.Infoln(m)
		common.Log.Infof("发送短息的电话号码为空")
		return
	}
	if len(m.Msg) <= 0 {
		common.Log.Infof("发送的内容为空")
		return
	}
	ob:=new(module.TemplateSmsLog)
	ob.Status = 1
	ob.Mobile = m.Mobile
	ob.MemberId = m.MemberId
	ob.CorpId = m.CorpId
	ob.TemplateManageId = m.TemplateManageId
	ob.MallId = m.MallId
	err = new(sms.SMSMessage).SendMobileMessage(m.Mobile, m.Msg)
	if err != nil {
		common.Log.Errorln(err)
		ob.Status = 2
	}

	ob.Created = time.Now() //.Format("2006-01-02 15:04:05")
	err = ob.InsertDb()
	if err != nil {
		common.Log.Errorln(err)
	}
	return
}

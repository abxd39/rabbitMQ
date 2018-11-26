package service

import (
	"encoding/json"
	Log "github.com/sirupsen/logrus"
	"sctek.com/typhoon/th-platform-gateway/common"
	"sctek.com/typhoon/th-platform-gateway/common/worker"
	"sctek.com/typhoon/th-platform-gateway/module"
	"sctek.com/typhoon/th-platform-gateway/sms"
	"time"
)

var pool *worker.Pool

func InitPool() {
	pool = worker.NewPool(common.Config.MaxQueueSize)
	pool.Run(common.Config.MaxWork)
	Log.Infof("goroutine的个数为%v,最大任务数为%v", common.Config.MaxWork, common.Config.MaxQueueSize)
}

func ClosePool() {
	pool.Shutdown()
}

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

//负责发送短息和写发送记录到数据库
func (m *MarshalJson) Run() error {
	ob := new(module.TemplateSmsLog)

	ob.Status = 1
	ob.Mobile = m.Mobile
	ob.MemberId = m.MemberId
	ob.CorpId = m.CorpId
	ob.TemplateManageId = m.TemplateManageId
	ob.MallId = m.MallId
	if ob.MemberId == -1 {
		ob.Status = 2
	} else {
		err := new(sms.SMSMessage).SendMobileMessage(m.Mobile, m.Msg)
		if err != nil {
			Log.Errorf("发送的信息为%#v",*m)
			Log.Errorln(err)
			ob.Status = 2
		}
	}
	ob.Created = time.Now() //.Format("2006-01-02 15:04:05")
	err := ob.InsertDb()
	if err != nil {
		Log.Errorln(err)
	}
	return nil
}

func (m *MarshalJson) UnmarshalJson(body []byte) {
	//common.Log.Traceln("添加到工作池")
	err := json.Unmarshal(body, m)
	if err != nil {
		Log.Errorln(err)
		return
	}
	pool.Add(m)
	//common.Log.Infoln("添加到工作池成功！！")
	return
}

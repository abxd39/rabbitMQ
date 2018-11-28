package service

import (
	"encoding/json"
	"fmt"
	Log "github.com/sirupsen/logrus"
	"sctek.com/typhoon/th-platform-gateway/weChat"
	"strconv"
)

func SmsCallback(d MSG) {
	Log.Infof("consumer-name=%v", d.Poper)
	//发送短息
	new(MarshalJson).UnmarshalJson(d.Body)
}

func weChatServiceAccountTemplateMessageCallback(d MSG){
	Log.Infof("consumer-name=%v",d.Poper)
	ob:=new(weChat.WeChatMp)
	ob.Body =d.Body
	ob.Mark =1
	ob.ReceiveMqWeChatMessage()
}

func errCallback(d MSG) {
	Log.Infoln("errServerQueue consumer")
	fmt.Println(string(d.Body))
}

func miniProgramTemplateMessageCallback(d MSG) {
	Log.Infof("consumer-name=%v",d.Poper)
	ob:=new(weChat.WeChatMp)
	ob.Body =d.Body
	ob.Mark =2
	ob.ReceiveMqWeChatMessage()
}

func DbIdCallback(d MSG) {
	Log.Infof("consumer-name=%v ", d.Poper)
	fmt.Printf("mq中读到的数据为：%q\r\n", string(d.Body))
	UnmarshalMQBody(d.Body)
}

func UnmarshalMQBody(body []byte) error {
	Log.Infoln("数据库条目 id 解码")
	result := &struct {
		Id string `json:"id"`
	}{}
	err := json.Unmarshal(body, result)
	if err != nil {
		Log.Infoln(err)
		return err
	}
	id, err := strconv.Atoi(result.Id)
	if err != nil {
		Log.Errorln(err)
		return err
	}
	new(LogicService).AboutIdInfo(id)
	return nil
}

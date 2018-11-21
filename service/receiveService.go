package service

import (
	"encoding/json"
	"fmt"
	"sctek.com/typhoon/th-platform-gateway/common"
	"strconv"
)

func callback(d MSG) {
	common.Log.Infof("consumer-name=%v",d.Poper)
	//发送短息
	new(MarshalJson).UnmarshalJson(d.Body)
}

func errCallback(d MSG) {
	common.Log.Infoln("errServerQueue consumer")
	fmt.Println(string(d.Body))
}

func dlxCallback(d MSG) {
	common.Log.Infoln("dlxQueue consumer")
	fmt.Println(string(d.Body))
}

func otherCallback(d MSG) {
	common.Log.Infof("consumer-name=%v ",d.Poper)
	fmt.Printf("mq中读到的数据为：%q\r\n",string(d.Body))
	UnmarshalMQBody(d.Body)
}

func  UnmarshalMQBody(body []byte) error {
	common.Log.Infoln("数据库条目 id 解码")
	result := &struct {
		Id string `json:"id"`
	}{}
	err := json.Unmarshal(body, result)
	if err != nil {
		common.Log.Infoln(err)
		return err
	}
	id,err:=strconv.Atoi(result.Id)
	if err!=nil{
		common.Log.Errorln(err)
		return err
	}
	new(LogicService).AboutIdInfo(id)
	return nil
}


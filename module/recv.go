package module

import (
	"encoding/json"
	"fmt"
	"sctek.com/typhoon/th-platform-gateway/common"
	"strconv"
)

func callback(d MSG) {
	common.Log.Infoln("yf_manage_message  consumer")
	fmt.Println(string(d.Body))
	//发送短息
	new(TemplateSmsLog).SendMobileMessage(d.Body)
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
	common.Log.Infoln("yf_sms_send consumer ")
	fmt.Println(string(d.Body))
	UnmarshalMQBody(d.Body)
}

func  UnmarshalMQBody(body []byte) error {
	//log.Println(string(body))
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
	//测试
	//指定会员性别即时发送
	temp:=new(TemplateSmsManage)
	temp.AboutIdInfo(34)
	temp.AboutIdInfo(36)
	//指定会员性别定时发送
	//指定会员等级即时发送
	//指定会员等级定时发送
	//指定会员定生日即时发送
	//指定会员生日定时发送
	//指定手机号即时发送
	//指定手机号定时发送
	//全员定时发送
	//全员即时发送
	return new(TemplateSmsManage).AboutIdInfo(id)
}

//func main() {
//	if err := rmq.Init("E:/WorkSpace/src/go-rabbitmq/example/rmq.json"); err != nil {
//		fmt.Println(err)
//		return
//	}
//
//	if err := rmq.Pop("myPoper", callback); err != nil {
//		fmt.Println(err)
//		return
//	}
//
//	if err := rmq.Pop("errPoper", errCallback); err != nil {
//		fmt.Println(err)
//		return
//	}
//
//	if err := rmq.Pop("dlxPoper", dlxCallback); err != nil {
//		fmt.Println(err)
//		return
//	}
//	if err:=rmq.Pop("first",otherCallback);err!=nil{
//		fmt.Println(err)
//		return
//	}
//
//	//time.Sleep(time.Duration(1000) * time.Second)
//	//Wait for interrupt signal to gracefully shutdown the server with
//	//a timeout of 30 seconds.
//	quit := make(chan os.Signal)
//	signal.Notify(quit, os.Interrupt)
//	<-quit
//	if err := rmq.Fini(); err != nil {
//		fmt.Println(err)
//	}
//}

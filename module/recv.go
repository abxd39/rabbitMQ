package module

import (
	"encoding/json"
	"fmt"
	"sctek.com/pingtai/consumer/common"
)

func callback(d MSG) {
	fmt.Println("yf_manage_message  consumer")
	fmt.Println(string(d.Body))
	//发送短息
}

func errCallback(d MSG) {
	fmt.Println("errServerQueue consumer")
	fmt.Println(string(d.Body))
}

func dlxCallback(d MSG) {
	fmt.Println("dlxQueue consumer")
	fmt.Println(string(d.Body))
}

func otherCallback(d MSG) {
	fmt.Println("yf_sms_send consumer ")
	fmt.Println(string(d.Body))
	UnmarshalMQBody(d.Body)
}

func  UnmarshalMQBody(body []byte) error {
	//log.Println(string(body))
	result := &struct {
		Id int `json:"id"`
	}{}
	err := json.Unmarshal(body, result)
	if err != nil {
		common.Log.Infoln(err)
		return err
	}
	return new(TemplateSmsManage).AboutIdInfo(result.Id)
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

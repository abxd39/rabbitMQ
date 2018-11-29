package Jpush

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	Log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"sctek.com/typhoon/th-platform-gateway/common"
)

type JPush struct {
	Context []byte
}

const Url = "https://api.jpush.cn/v3/push"

func (j*JPush)Run()error  {
	err:=j.JMessage()
	if err!=nil{
		return err
	}
	//再此对数据库进行操作
	return nil
}

func(j*JPush)ReceiveJPushMessageFromMQ(){
	common.Pool.Add(j)
}

func (j *JPush) JMessage() error{
	//base64(appKey:masterSecret)
	Authorization := fmt.Sprintf("%s:%s", common.Config.JPush.AppKey, common.Config.JPush.MasterSecret)
	encode := base64.StdEncoding.EncodeToString([]byte(Authorization))
	params := make(map[string]interface{})
	params["test"] = "test"
	bytsData, err := json.Marshal(params)
	if err != nil {
		Log.Errorln(err)
		return err
	}
	reader:=bytes.NewReader(bytsData)
	request ,err:=http.NewRequest("POST",Url,reader)
	request.Header.Set("Content-Type","application/json;charset=UTF-8")
	request.Header.Add("Authorization","Basic "+encode)
	client:=http.Client{}
	response,err:=client.Do(request)
	if err!=nil{
		Log.Errorln(err)
		return err
	}
	result :=&struct {
		SendNo string `json:"sendno"`
		MsgId string `json:"msg_id"`
		Error struct{
			Code uint32 `json:"code"`
			Message string `json:"message"`
		}
	}{}
	body ,err:=ioutil.ReadAll(response.Body)
	if err!=nil{
		Log.Errorln(err)
		return err
	}
	if len(body)<=0{
		err=fmt.Errorf("极光消息推送请求的返回内容为空！！")
		Log.Errorln(err)
		return err
	}
	err = json.Unmarshal(body,result)
	if err!=nil{
		Log.Errorln(err)
		return err
	}
	Log.Errorln("JPush successful ^/~\\^")
	return nil
}

package weChat

import (
	"bytes"
	"encoding/json"
	"fmt"
	Log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"sctek.com/typhoon/th-platform-gateway/common"
	"strings"
)

type WeChatMp struct {
	token string
	Body []byte
}

func (w *WeChatMp)ReceiveMqWeChatMessage(){
	err:=w.GetTokenWeChat()
	if err!=nil{
		Log.Errorln(err)
		return
	}
	Log.Infof(w.token)
	return
	common.Pool.Add(w)
}

func (w*WeChatMp)Run()error{
	w.SendTemplateMessages()
	return nil
}

func(w*WeChatMp) GetTokenWeChat()error{
	reader:= strings.NewReader("------WebKitFormBoundary7MA4YWxkTrZu0gW\r\nContent-Disposition: form-data; name=\"mall_id\"\r\n\r\n97\r\n------WebKitFormBoundary7MA4YWxkTrZu0gW\r\nContent-Disposition: form-data; name=\"auth_type\"\r\n\r\n1\r\n------WebKitFormBoundary7MA4YWxkTrZu0gW--")
	url := "http://boss-api.dev-ibc.snsshop.net/wechat/tpop/get-authorizer-access-token?mall_id=97&auth_type=1"
	request, err := http.NewRequest("POST", url,reader)
	if err != nil {
		Log.Errorln(err)
		return err
	}
	request.Header.Set("content-type", "multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW")
	client := http.Client{}
	result, err := client.Do(request)
	if err != nil {
		Log.Errorln(err)
		return err
	}
	body, err := ioutil.ReadAll(result.Body)
	if err != nil {
		Log.Errorln(err)
		return err
	}
	value:= &struct {
		ErrCode int `json:"errcode"`
		ErrMsg string `json:"errmsg"`
		Token string `json:"token"`
	}{}
	Log.Infof("返回值为：%s",string(body))
	err=json.Unmarshal(body,value)
	if err!=nil{
		Log.Errorln(err)
		return err
	}
	if value.ErrCode !=0{
		err:=fmt.Errorf("errcode:%v,errmsg:%v",value.ErrCode,value.ErrMsg)
		Log.Errorln(err)
		return err
	}
	w.token = value.Token
	return nil
}


//Setting Associated Industries
func (w WeChatMp) Industries() error {
	params := make(map[string]interface{})
	params["access_token"] = w.token
	params["industry_id1"] = "1"
	params["industry_id2"] = "4"
	bytesData, err := json.Marshal(params)
	if err != nil {
		return err
	}
	reader := bytes.NewReader(bytesData)
	url := "https://api.weixin.qq.com/cgi-bin/template/api_set_industry"
	request, err := http.NewRequest("GET", url, reader)
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	client := http.Client{}
	result, err := client.Do(request)
	if err != nil {
		Log.Errorln(err)
		return err
	}
	body, err := ioutil.ReadAll(result.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(body))
	return nil
}

//公众号发送模板消息
func (w WeChatMp)SendTemplateMessages()error  {
	reader := bytes.NewReader(w.Body)
	str:= fmt.Sprintf("access_toke=%v",w.token)
	url := "https://api.weixin.qq.com/cgi-bin/message/template/send?"+str
	request, err := http.NewRequest("POST", url, reader)
	if err != nil {
		Log.Errorln(err)
		return err
	}
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	client := http.Client{}
	rsp, err := client.Do(request)
	if err != nil {
		Log.Errorln(err)
		return err
	}
	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		Log.Errorln(err)
		return err
	}
	result :=&struct {
		ErrMsg string `json:"errmsg"`
		ErrCode int `json:"errcode"`
		ErrId uint32 `json:"errid"`
	}{}
	err=json.Unmarshal(body,result)
	if err!=nil{
		Log.Errorln(err)
		return err
	}
	if result.ErrCode!=0{
		err=fmt.Errorf("ErrCode:%v,ErrMsg:%v",result.ErrCode,result.ErrMsg)
		Log.Errorln(err)
		return err
	}
	return nil
}


//小程序发送模板消息
func (w WeChatMp)SendTemplateMessage1()error{
	reader := bytes.NewReader(w.Body)
	str:= fmt.Sprintf("access_toke=%v",w.token)
	url:="https://api.weixin.qq.com/cgi-bin/message/wxopen/template/send"+str
	request, err := http.NewRequest("POST", url, reader)
	if err != nil {
		Log.Errorln(err)
		return err
	}
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	client := http.Client{}
	rsp, err := client.Do(request)
	if err != nil {
		Log.Errorln(err)
		return err
	}
	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		Log.Errorln(err)
		return err
	}
	result :=&struct {
		ErrMsg string `json:"errmsg"`
		ErrCode int `json:"errcode"`
		TemplateId uint32 `json:"template_id"`
	}{}
	err=json.Unmarshal(body,result)
	if err!=nil{
		Log.Errorln(err)
		return err
	}
	if result.ErrCode!=0{
		err=fmt.Errorf("ErrCode:%v,ErrMsg:%v",result.ErrCode,result.ErrMsg)
		Log.Errorln(err)
		return err
	}
	return nil
}
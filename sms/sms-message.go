package sms

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sctek.com/pingtai/consumer/common"
	"strconv"
	"time"
)

type FailMessageStruct struct {
	Params     string
	Results    string
	CreateTime string
	ModifyTime string
}

/**
 * 队列中取出来的数据，需要转化取出回调的url
 */
type PostDataStruct struct {
	Callback string
}

func RecordFailMsg(params FailMessageStruct) {
	sql := "insert into sms_fail_message (params,results,create_time,modify_time) values ('" + params.Params + "','" + params.Results + "','" + params.CreateTime + "','" + params.ModifyTime + "')"
	common.DB.Exec(sql)
}

/**
 * 发送消息，调php接口
 */
func Send(queueName string, msg string) {
	var m PostDataStruct
	err := json.Unmarshal([]byte(msg), &m)
	if err != nil {
		common.Log.Infoln(err)
	}
	url := m.Callback

	//调curl把数据传到php执行
	result := execTask(url, msg)
	if result == "ok" {
		return
	}

	common.Log.Infoln(result)
	//如果失败，记录返回结果值到数据表中
	p := FailMessageStruct{
		Params:     msg,
		Results:    result,
		CreateTime: strconv.FormatInt(time.Now().Unix(), 10),
		ModifyTime: strconv.FormatInt(time.Now().Unix(), 10),
	}
	RecordFailMsg(p)
}

/**
 * 调http post到php页面, application/json
 */
func execTask(url string, msg string) string {
	var jsonStr = []byte(msg)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	defer req.Body.Close()

	if err != nil {
		return "接口对象创建失败，地址：" + url
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return "接口请求失败，地址：" + url
	}
	if resp.StatusCode == 200 {
		return "ok"
	}
	body, _ := ioutil.ReadAll(resp.Body)
	respMsg := string(body)
	return respMsg
}


//发送短息
type SMSMessage struct{}

func (s *SMSMessage) SendMobileMessage(phone, message string) error {
	params := make(map[string]interface{})
	params["mobile"] = phone
	params["msg"] = message
	params["send_type"] = "ibc_mall_sign"
	bytesData, err := json.Marshal(params)
	if err != nil {
		return err
	}
	reader := bytes.NewReader(bytesData)
	url := "http://dev-ibc.snsshop.net/ec_crm/sms/qcloud_send?"
	request, err := http.NewRequest("POST", url, reader)
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	client := http.Client{}
	result, err := client.Do(request)
	if err != nil {
		return err
	}
	rsp := &struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}{}
	body, err := ioutil.ReadAll(result.Body)
	if err != nil {
		return err
	}
	log.Printf("发送短息的url:=%s", url)
	log.Printf("返回值=%s",string(body))
	err = json.Unmarshal(body, rsp)
	if err != nil {
		return err
	}
	if rsp.Code != 0 {
		return fmt.Errorf(rsp.Msg)
	}
	return nil
}

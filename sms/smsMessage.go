package sms

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sctek.com/typhoon/th-platform-gateway/common"
)


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
	url := common.Config.Url+"ec_crm/sms/qcloud-send?"
	request, err := http.NewRequest("POST", url, reader)
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	client := http.Client{}
	result, err := client.Do(request)

	if err != nil {
		fmt.Println(err)
		return err
	}
	rsp := &struct {
		Code int    `json:"code"`
		ErrMsg  string `json:"errmsg"`
	}{}
	body, err := ioutil.ReadAll(result.Body)
	if err != nil {
		return err
	}
	log.Printf("发送短息的url:=%s", url)
	common.Log.Infof("返回值=%s\r\n",string(body))
	common.Log.Infof("%q\r\n",params)
	err = json.Unmarshal(body, rsp)
	if err != nil {
		return err
	}
	if rsp.Code !=-1 {
		return fmt.Errorf("发送短息返回的错误信息为【%v】\r\n",rsp.ErrMsg)
	}
	common.Log.Infoln("短息发送成功^~^")
	return nil
}

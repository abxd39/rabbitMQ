package weChat

import (
	"bytes"
	"encoding/json"
	"fmt"
	Log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

type WeChatMp struct {
}

//Setting Associated Industries
func (w WeChatMp) Industries() error {
	params := make(map[string]interface{})
	params["access_token"] = "ACCESS_TOKEN"
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
	}
	body, err := ioutil.ReadAll(result.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(body))
	return nil
}

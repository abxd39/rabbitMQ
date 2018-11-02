package util

import (
	"encoding/json"
)

// ToJson 转换成json字符串
func ToJson(data interface{}) string {
	result, err := json.Marshal(data)
	if err != nil {
		return ""
	}
	return string(result)
}

// FromJson 转换成json对象
func FromJson(val string, a interface{}) error {
	err := json.Unmarshal([]byte(val), &a)
	return err
}

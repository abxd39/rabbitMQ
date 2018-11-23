package common

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	Log "github.com/sirupsen/logrus"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type ErrorRes struct {
	Code int `json:"code"`
	Desc interface{}
}

type WxDecrypted struct {
	PhoneNumber string `json:"phoneNumber"`
}

func GetSexIndex(sSex string) int {
	var sexIndex int
	switch sSex {
	case "男":
		sexIndex = 1
		break
	case "女":
		sexIndex = 2
		break
	default:
		sexIndex = 0
		break
	}

	return sexIndex
}

func Floor() map[string]string {
	return map[string]string{
		"-3": "B3",
		"-2": "B2",
		"-1": "B1",
		"1":  "1F",
		"2":  "2F",
		"3":  "3F",
		"4":  "4F",
		"5":  "5F",
		"6":  "6F",
		"7":  "7F",
		"8":  "8F",
		"9":  "9F",
	}
}

//分转元
func FenToYuan(i int) string {
	return fmt.Sprintf("%.2f", float64(i)/100)
}

//时间戳转字符串
func IntTimeToString(t int64) string {
	return time.Unix(t, 0).Format("2006.01.02")
}

func GetTimestamp() int64 {
	return time.Now().Unix()
}

func GetTimestampStr() string {
	ts := GetTimestamp()
	return strconv.FormatInt(ts, 10)
}

func I64ToStr(p int64) string {
	return strconv.FormatInt(p, 10)
}

//格式化时间
func FormatTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

func CheckErr(err error) {
	if err != nil {
		log.Fatal(err)
		//fmt.Println(fmt.Printf("%s %s", time.Now().Format(common.DATETIME), err.Error()))
	}
}

func GetUserId(c *gin.Context) int {
	return 293 //用于测试
	return c.GetInt("user_id")
}

func GetMallId(c *gin.Context) int {
	mallId, err := strconv.Atoi(c.Params.ByName("mall_id"))
	if err != nil {
		mallId = 0
	}
	return mallId
}

func BytesToInt(b []byte) int {
	r, _ := strconv.Atoi(string(b))
	return r
}

func RenderJSON(c *gin.Context, obj interface{}) {
	c.JSON(http.StatusOK, obj)
}

func RenderJSONWithError(c *gin.Context, err error, status ...int) {
	code := 0
	s := http.StatusInternalServerError
	cerr, ok := err.(BadRequestError)
	if ok {
		code = cerr.Code()
		s = http.StatusBadRequest
	}
	if len(status) > 0 {
		s = status[0]
	}
	c.AbortWithStatusJSON(s, gin.H{
		"status":  s,
		"code":    code,
		"message": err.Error(),
	})
}

//post 请求
func Http_Post(sUrl string, p interface{}) ([]byte, error) {
	bytesData, err := json.Marshal(p)
	if err != nil {
		Log.Errorln(err)
		return nil, err
	}
	reader := bytes.NewReader(bytesData)
	resp, err := http.Post(sUrl, "application/json", reader)
	if err != nil {
		Log.Errorln(err)
		return nil, err
	}
	defer resp.Body.Close()

	//判断返回的状态码是否是200
	if resp.StatusCode != http.StatusOK {
		return nil, NewError(resp.StatusCode)
	}

	return ioutil.ReadAll(resp.Body)
}

//post 请求
func Http_PostForm(sUrl string, uv url.Values) ([]byte, error) {
	resp, err := http.PostForm(sUrl, uv)
	if err != nil {
		Log.Errorln(err)
		return nil, err
	}
	defer resp.Body.Close()

	//判断返回的状态码是否是200
	if resp.StatusCode != http.StatusOK {
		return nil, NewError(resp.StatusCode)
	}

	return ioutil.ReadAll(resp.Body)
}

//截取字符串 start 起点下标 end 终点下标(不包括)
func Substr(str string, start int, end int) string {
	rs := []rune(str)
	length := len(rs)

	if start < 0 || start > length {
		panic("start is wrong")
	}

	if end < 0 || end > length {
		panic("end is wrong")
	}

	return string(rs[start:end])
}

func RandInt64(min, max int64) int64 {
	maxBigInt := big.NewInt(max)
	i, _ := rand.Int(rand.Reader, maxBigInt)
	if i.Int64() < min {
		RandInt64(min, max)
	}
	return i.Int64()
}

func GetWxDecryptedData(EncryptedData, sessionKey, iv string) (info *WxDecrypted, err error) {

	cipherText, err := base64.StdEncoding.DecodeString(EncryptedData)

	aesKey, err := base64.StdEncoding.DecodeString(sessionKey)
	aesIv, err := base64.StdEncoding.DecodeString(iv)

	if err != nil {
		return
	}

	raw, err := AESDecryptData(cipherText, aesKey, aesIv)

	if err != nil {
		return
	}

	if err = json.Unmarshal(raw, &info); err != nil {
		return
	}
	return
}

func AESDecryptData(cipherText []byte, aesKey []byte, iv []byte) (rawData []byte, err error) {

	const (
		BLOCK_SIZE = 32             // PKCS#7
		BLOCK_MASK = BLOCK_SIZE - 1 // BLOCK_SIZE 为 2^n 时, 可以用 mask 获取针对 BLOCK_SIZE 的余数
	)

	if len(cipherText) < BLOCK_SIZE {
		err = fmt.Errorf("the length of ciphertext too short: %d", len(cipherText))
		return
	}

	plaintext := make([]byte, len(cipherText)) // len(plaintext) >= BLOCK_SIZE

	// 解密
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		panic(err)
	}
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(plaintext, cipherText)

	// PKCS#7 去除补位
	amountToPad := int(plaintext[len(plaintext)-1])
	if amountToPad < 1 || amountToPad > BLOCK_SIZE {
		err = fmt.Errorf("the amount to pad is incorrect: %d", amountToPad)
		return
	}
	plaintext = plaintext[:len(plaintext)-amountToPad]

	// 反拼接
	// len(plaintext) == 16+4+len(rawXMLMsg)+len(appId)
	if len(plaintext) <= 20 {
		err = fmt.Errorf("plaintext too short, the length is %d", len(plaintext))
		return
	}

	rawData = plaintext

	return

}

func QueryInt(c *gin.Context, key string) int {
	val := c.Query(key)
	result, _ := strconv.Atoi(val)
	return result
}

func DefaultQueryInt(c *gin.Context, key, defaultValue string) int {
	val := c.DefaultQuery(key, defaultValue)
	result, _ := strconv.Atoi(val)
	return result
}

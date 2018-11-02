package util

import (
	"crypto/md5"
	"encoding/hex"
)

type MD5Client struct {
}

var MD5 = MD5Client{}

// Encrypt 获取字节
func (c *MD5Client) Encrypt(plantext []byte) []byte {
	result := md5.Sum(plantext)
	return result[:]
}

// GetEncryptString 获取加密字符串
func (c *MD5Client) GetEncryptString(plantext []byte) string {
	byteResult := c.Encrypt(plantext)
	result := hex.EncodeToString(byteResult)
	return result
}

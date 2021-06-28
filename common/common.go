package common

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"os"
)

// HmacSha256 加密
func HmacSha256(data string) string {
	hash := hmac.New(sha256.New, []byte(os.Getenv("KEY"))) //创建对应的sha256哈希加密算法
	hash.Write([]byte(data))
	return hex.EncodeToString(hash.Sum([]byte("")))
}

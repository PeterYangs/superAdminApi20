package common

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"os"
)

// HmacSha256 加密
func HmacSha256(data string) string {
	hash := hmac.New(sha256.New, []byte(os.Getenv("KEY"))) //创建对应的sha256哈希加密算法
	hash.Write([]byte(data))
	return hex.EncodeToString(hash.Sum([]byte("")))
}

func Paginate(tx *gorm.DB, dest interface{}, page int, size int) gin.H {

	offset := (page - 1) * size

	var count int64

	tx.Count(&count)

	tx.Offset(offset).Limit(size).Find(dest)

	return gin.H{"total": count, "data": dest, "page": page, "size": size}

}

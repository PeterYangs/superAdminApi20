package common

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"os"
	"reflect"
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

// UpdateOrCreateOne 更新或修改
func UpdateOrCreateOne(tx *gorm.DB, model interface{}, where map[string]interface{}, modelData interface{}, Omits ...string) error {

	tt := tx.Model(model)

	for s, i := range where {

		tt.Where(s+"=?", i)
	}

	re := tt.First(model)

	id := reflect.ValueOf(model).Elem().FieldByName("Id").Interface()

	if re.Error == gorm.ErrRecordNotFound {

		tx.Model(model).Create(modelData)

		return nil

	}

	if re.Error == nil {

		//fmt.Println("更新")

		fmt.Println(id)

		up := tx.Model(model).Where("id=?", id)

		if len(Omits) > 0 {

			up.Omit(Omits...)

		}

		//fmt.Println(modelData)

		up.Updates(modelData)

		//fmt.Println(rrr.Error)

		return nil
	}

	return re.Error

}

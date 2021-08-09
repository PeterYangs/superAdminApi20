package common

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"github.com/PeterYangs/tools"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"math/big"

	"crypto/rand"
	//"math/rand"
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
/**
Omits 是忽略的字段
*/
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

		//fmt.Println(id)

		up := tx.Model(model).Where("id=?", id)

		if len(Omits) > 0 {

			up.Omit(Omits...)

		}

		s := reflect.TypeOf(modelData).Elem()

		ss, b := s.FieldByName("Id")

		if b {

			fillable := ss.Tag.Get("fillable")

			if fillable != "" {

				f := tools.Explode(",", fillable)

				up = up.Select(f)
			}

		}

		up.Updates(modelData)

		return nil
	}

	return re.Error

}

func MtRand(min, max int64) int64 {

	//rand.Seed(time.Now().UnixNano())

	//return rand.Intn(max-min+1) + min

	n, _ := rand.Int(rand.Reader, big.NewInt(max-min+1))

	return n.Int64() + min
}

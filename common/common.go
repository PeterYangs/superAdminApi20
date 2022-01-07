package common

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/PeterYangs/tools"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"math/big"
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

	db := tx.WithContext(context.TODO())

	offset := (page - 1) * size

	var count int64

	db.Count(&count)

	db.Offset(offset).Limit(size).Find(dest)

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

	re := tt.First(map[string]interface{}{})

	//id := reflect.ValueOf(model).Elem().FieldByName("Id").Interface()

	if re.Error == gorm.ErrRecordNotFound {

		fmt.Println("你xx")

		cRe := tx.Create(modelData)

		if cRe.Error != nil {

			return cRe.Error
		}

		return nil

	}

	if re.Error == nil {

		up := tx.Model(model)

		for s, i := range where {

			up.Where(s+"=?", i)
		}

		if len(Omits) > 0 {

			up.Omit(Omits...)

		}

		s := reflect.TypeOf(modelData).Elem()

		ss, b := s.FieldByName("Id")

		//需要更新的字段
		if b {

			fillable := ss.Tag.Get("fillable")

			if fillable != "" {

				f := tools.Explode(",", fillable)

				up.Select(f)
			}

		}

		uRe := up.Updates(modelData)

		if uRe.Error != nil {

			return uRe.Error
		}

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

// Capitalize 字符首字母大写转换
func Capitalize(str string) string {
	var upperStr string
	vv := []rune(str) // 后文有介绍
	for i := 0; i < len(vv); i++ {
		if i == 0 {
			if vv[i] >= 97 && vv[i] <= 122 { // 后文有介绍
				vv[i] -= 32 // string的码表相差32位
				upperStr += string(vv[i])
			} else {
				fmt.Println("Not begins with lowercase letter,")
				return str
			}
		} else {
			upperStr += string(vv[i])
		}
	}
	return upperStr
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 8)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

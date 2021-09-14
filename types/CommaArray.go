package types

import (
	"database/sql/driver"
	"github.com/PeterYangs/tools"
	"github.com/spf13/cast"
)

type CommaArray []int

//获取器（读取结果）
func (m *CommaArray) Scan(value interface{}) error {

	bytes, ok := value.([]byte)

	if !ok {
		//return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}

	//err := json.Unmarshal(bytes, &r)

	s := tools.Explode(",", string(bytes))

	r := make([]int, len(s))

	for i, s2 := range s {

		r[i] = cast.ToInt(s2)

	}

	//if err != nil {
	//
	//	*m = []interface{}{}
	//
	//	return nil
	//}

	*m = r

	return nil

}

//写入数据库(修改器)
func (m CommaArray) Value() (driver.Value, error) {

	r := make([]string, len(m))

	for i, i2 := range m {

		r[i] = cast.ToString(i2)
	}

	return tools.Join(",", r), nil

}

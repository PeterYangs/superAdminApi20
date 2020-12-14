package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type JsonMap map[string]interface{}

//将数据库内容转成对应类型（获取器）
func (m *JsonMap) Scan(value interface{}) error {

	bytes, ok := value.([]byte)

	if !ok {
		//return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}

	var r map[string]interface{}

	err := json.Unmarshal(bytes, &r)

	fmt.Println()

	if err != nil {

		*m = map[string]interface{}{}

		return nil
	}

	*m = r

	return nil

}

//将插入内容转成数据库(修改器)
func (m JsonMap) Value() (driver.Value, error) {

	return json.Marshal(m)

}

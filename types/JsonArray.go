package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type JsonArray []interface{}

func (m *JsonArray) Scan(value interface{}) error {

	bytes, ok := value.([]byte)

	if !ok {
		//return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}

	var r []interface{}

	err := json.Unmarshal(bytes, &r)

	fmt.Println()

	if err != nil {

		*m = []interface{}{}

		return nil
	}

	*m = r

	return nil

}

//将插入内容转成数据库(修改器)
func (m JsonArray) Value() (driver.Value, error) {

	return json.Marshal(m)

}

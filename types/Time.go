package types

import (
	"database/sql/driver"
	"time"
)

type Time time.Time

func (t Time) Value() (driver.Value, error) {

	return time.Time(t), nil
}

func (t Time) MarshalJSON() ([]byte, error) {
	tTime := time.Time(t)
	tStr := tTime.Format("2006-01-02 15:04:05") // 设置格式

	if tStr == "0001-01-01 00:00:00" {

		return []byte("\"" + " " + "\""), nil
	}

	//// 注意 json 字符串风格要求
	return []byte("\"" + tStr + "\""), nil

}

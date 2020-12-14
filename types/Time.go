package types

import (
	"database/sql/driver"
	"time"
)

type Time time.Time

func (t Time) Value() (driver.Value, error) {

	tTime := time.Time(t)

	return tTime.Format("2006-01-02 15:04:05"), nil

}

func (t Time) MarshalJSON() ([]byte, error) {
	tTime := time.Time(t)
	tStr := tTime.Format("2006-01-02 15:04:05") // 设置格式

	// 注意 json 字符串风格要求
	return []byte("\"" + tStr + "\""), nil

}

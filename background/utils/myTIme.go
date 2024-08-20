package utils

import (
	"database/sql/driver"
	"fmt"
	"time"
)

// CustomTime 是自定义的时间类型
type CustomTime struct {
	time.Time
}

const customTimeFormat1 = "2006-01-02 15:04:05"
const customTimeFormat = "2006-01-02 15:04"

// Value 将 CustomTime 转换为数据库中可存储的格式。
// GORM 在将数据保存到数据库时会调用此方法。
// 这里将 CustomTime 格式化为字符串形式，以便数据库能够存储该字符串。
func (ct CustomTime) Value() (driver.Value, error) {
	return ct.Format(customTimeFormat), nil
}

// Scan 从数据库中读取时间字段并将其转换为 CustomTime。
// GORM 在从数据库加载数据时会调用此方法。
// 根据数据库返回的数据类型（如 time.Time 或 string），将其转换为 CustomTime。
// 如果数据库返回的是时间戳或其他格式的数据，则可以根据实际需要扩展支持更多类型。
func (ct *CustomTime) Scan(value interface{}) error {
	if value == nil {
		*ct = CustomTime{Time: time.Time{}}
		return nil
	}
	switch v := value.(type) {
	case time.Time:
		*ct = CustomTime{Time: v}
	case string:
		t, err := time.Parse(customTimeFormat, v)
		if err != nil {
			return fmt.Errorf("failed to parse time: %v", err)
		}
		*ct = CustomTime{Time: t}
	case []byte:
		t, err := time.Parse(customTimeFormat, string(v))
		if err != nil {
			return fmt.Errorf("failed to parse time: %v", err)
		}
		*ct = CustomTime{Time: t}
	default:
		return fmt.Errorf("unsupported type: %T", v)
	}
	return nil
}

// MarshalJSON 将 CustomTime 转换为 JSON 格式的字符串。
// 在将结构体转换为 JSON 时会调用此方法。
// 这里将时间格式化为自定义的时间字符串（如 "2006-01-02 15:04"），然后返回格式化后的 JSON 字符串。
func (ct CustomTime) MarshalJSON() ([]byte, error) {
	formatted := fmt.Sprintf("\"%s\"", ct.Format(customTimeFormat))
	return []byte(formatted), nil
}

// UnmarshalJSON 从 JSON 字符串解析并将其转换为 CustomTime。
// 在从 JSON 数据反序列化为结构体时会调用此方法。
// 解析自定义格式的时间字符串（如 "2006-01-02 15:04"），并将其转换为 time.Time 类型的时间对象。
func (ct *CustomTime) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}
	t, err := time.Parse(`"`+customTimeFormat+`"`, string(data))
	if err != nil {
		return fmt.Errorf("failed to parse time: %v", err)
	}
	ct.Time = t
	return nil
}

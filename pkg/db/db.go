package db

import "time"

// 时间格式模板
const (
	timeFormat = "2006-01-02 15:04:05"
)

// Time 自定义时间类型
type Time time.Time

// BaseModel 数据库中表基础结构
type BaseModel struct {
	Id        int64 // id
	CreatedAt Time  // 创建时间
	UpdateAt  Time  // 修改时间
}

// UnmarshalJSON 将 JSON 中的时间字符串解码为Time类型
func (t *Time) UnmarshalJSON(data []byte) (err error) {
	// 解析时间字符串
	now, err := time.ParseInLocation(
		`"`+timeFormat+`"`, // 时间格式
		string(data),       // json 数据
		time.Local)         // 本地时区
	// 转换为自定义Time类型，并赋值
	*t = Time(now)
	return
}

// MarshalJSON 将 Time 类型编码为 JSON 格式的字符串
func (t Time) MarshalJSON() ([]byte, error) {
	// 创建一个字节切片，长度为 timeFormat 加上前后的引号
	b := make([]byte, 0, len(timeFormat)+2)
	// 添加前引号
	b = append(b, '"')
	// 将 Time 类型的时间格式化后添加到切片中。
	b = time.Time(t).AppendFormat(b, timeFormat)
	// 添加后引号
	b = append(b, '"')
	// 返回
	return b, nil
}

// String 实现Stringer接口，提供Time类型的字符串表示
func (t Time) String() string {
	return time.Time(t).Format(timeFormat)
}

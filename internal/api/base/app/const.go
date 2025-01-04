package app

// Status app 状态
type Status int

// 定义状态常量 iota 是Go的特殊常量生成器，从 0 开始，逐行递增
const (
	// StatusDisable app 被禁用
	StatusDisable Status = iota
	// StatusEnable app被启用
	StatusEnable
)

// Int 转换Status为int
func (s Status) Int() int {
	return int(s)
}

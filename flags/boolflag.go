package flags

import (
	"strconv"
	"strings"
)

// BoolFlag 布尔类型标志结构体
// 继承BaseFlag[bool]泛型结构体,实现Flag接口
type BoolFlag struct {
	BaseFlag[bool]
}

// Type 返回标志类型
func (f *BoolFlag) Type() FlagType { return FlagTypeBool }

// Set 实现flag.Value接口,解析并设置布尔值
func (f *BoolFlag) Set(value string) error {
	boolVal, err := strconv.ParseBool(strings.ToLower(value))
	if err != nil {
		return err
	}
	return f.BaseFlag.Set(boolVal)
}

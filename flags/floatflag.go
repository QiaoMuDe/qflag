package flags

// FloatFlag 浮点型标志结构体
// 继承BaseFlag[float64]泛型结构体,实现Flag接口
type FloatFlag struct {
	BaseFlag[float64]
}

// Type 返回标志类型
func (f *FloatFlag) Type() FlagType { return FlagTypeFloat }

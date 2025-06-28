package cmd

import (
	"time"

	"gitee.com/MM-Q/qflag/flags"
)

// String 添加字符串类型标志, 返回标志对象指针
//
// 参数依次为: 长标志名、短标志、默认值、帮助说明
//
// 返回值: 字符串标志对象指针
func (c *Cmd) String(longName, shortName, defValue, usage string) *flags.StringFlag {
	f := &flags.StringFlag{}
	c.StringVar(f, longName, shortName, defValue, usage)
	return f
}

// StringVar 绑定字符串类型标志到指针并内部注册Flag对象
//
// 参数依次为: 字符串标志指针、长标志名、短标志、默认值、帮助说明
func (c *Cmd) StringVar(f *flags.StringFlag, longName, shortName, defValue, usage string) {
	// 检查指针是否为nil
	if f == nil {
		panic("StringFlag pointer cannot be nil")
	}

	// 参数校验（复用公共函数）
	if validateErr := c.validateFlag(longName, shortName); validateErr != nil {
		panic(validateErr)
	}

	// 显式初始化当前值的默认值
	currentStr := new(string)
	*currentStr = defValue

	// 初始化Flag对象
	if initErr := f.Init(longName, shortName, defValue, usage, currentStr); initErr != nil {
		panic(initErr)
	}

	// 创建FlagMeta对象
	meta := &flags.FlagMeta{
		Flag: f, // 添加标志对象 - Flag对象
	}

	// 绑定短标志
	if shortName != "" {
		c.fs.StringVar(f.GetPointer(), shortName, defValue, usage)
	}

	// 绑定长标志
	if longName != "" {
		c.fs.StringVar(f.GetPointer(), longName, defValue, usage)
	}

	// 注册Flag对象
	if registerErr := c.flagRegistry.RegisterFlag(meta); registerErr != nil {
		panic(registerErr)
	}
}

// IntVar 绑定整数类型标志到指针并内部注册Flag对象
//
// 参数依次为: 整数标志指针、长标志名、短标志、默认值、帮助说明
func (c *Cmd) IntVar(f *flags.IntFlag, longName, shortName string, defValue int, usage string) {
	// 检查指针是否为nil
	if f == nil {
		panic("IntFlag pointer cannot be nil")
	}

	// 参数校验（复用公共函数）
	if validateErr := c.validateFlag(longName, shortName); validateErr != nil {
		panic(validateErr)
	}

	// 初始化默认值
	currentInt := new(int)
	*currentInt = defValue

	// 初始化Flag对象
	if initErr := f.Init(longName, shortName, defValue, usage, currentInt); initErr != nil {
		panic(initErr)
	}

	// 创建FlagMeta对象
	meta := &flags.FlagMeta{
		Flag: f, // 添加标志对象 - Flag对象
	}

	// 绑定短标志
	if shortName != "" {
		c.fs.IntVar(f.GetPointer(), shortName, defValue, usage)
	}

	// 绑定长标志
	if longName != "" {
		c.fs.IntVar(f.GetPointer(), longName, defValue, usage)
	}

	// 注册Flag对象
	if registerErr := c.flagRegistry.RegisterFlag(meta); registerErr != nil {
		panic(registerErr)
	}
}

// Int 添加整数类型标志, 返回标志对象指针
//
// 参数依次为: 长标志名、短标志、默认值、帮助说明
// 返回值: 整数标志对象指针
func (c *Cmd) Int(longName, shortName string, defValue int, usage string) *flags.IntFlag {
	f := &flags.IntFlag{}
	c.IntVar(f, longName, shortName, defValue, usage)
	return f
}

// BoolVar 绑定布尔类型标志到指针并内部注册Flag对象
//
// 参数依次为: 布尔标志指针、长标志名、短标志、默认值、帮助说明
func (c *Cmd) BoolVar(f *flags.BoolFlag, longName, shortName string, defValue bool, usage string) {
	// 检查指针是否为nil
	if f == nil {
		panic("BoolFlag pointer cannot be nil")
	}

	// 参数校验（复用公共函数）
	if validateErr := c.validateFlag(longName, shortName); validateErr != nil {
		panic(validateErr)
	}

	// 显式初始化
	currentBool := new(bool) // 创建当前值指针
	*currentBool = defValue

	// 初始化Flag对象
	if initErr := f.Init(longName, shortName, defValue, usage, currentBool); initErr != nil {
		panic(initErr)
	}

	// 创建FlagMeta对象
	meta := &flags.FlagMeta{
		Flag: f, // 添加标志对象 - Flag对象
	}

	// 绑定短标志
	if shortName != "" {
		c.fs.BoolVar(f.GetPointer(), shortName, defValue, usage)
	}

	// 绑定长标志
	if longName != "" {
		c.fs.BoolVar(f.GetPointer(), longName, defValue, usage)
	}

	// 注册Flag对象
	if registerErr := c.flagRegistry.RegisterFlag(meta); registerErr != nil {
		panic(registerErr)
	}
}

// Bool 添加布尔类型标志, 返回标志对象指针
//
// 参数依次为: 长标志名、短标志、默认值、帮助说明
//
// 返回值: 布尔标志对象指针
func (c *Cmd) Bool(longName, shortName string, defValue bool, usage string) *flags.BoolFlag {
	f := &flags.BoolFlag{}
	c.BoolVar(f, longName, shortName, defValue, usage)
	return f
}

// Float 添加浮点型标志, 返回标志对象指针
//
// 参数依次为: 长标志名、短标志、默认值、帮助说明
//
// 返回值: 浮点型标志对象指针
func (c *Cmd) Float(longName, shortName string, defValue float64, usage string) *flags.FloatFlag {
	f := &flags.FloatFlag{}
	c.FloatVar(f, longName, shortName, defValue, usage)
	return f
}

// FloatVar 绑定浮点型标志到指针并内部注册Flag对象
//
// 参数依次为: 浮点数标志指针、长标志名、短标志、默认值、帮助说明
func (c *Cmd) FloatVar(f *flags.FloatFlag, longName, shortName string, defValue float64, usage string) {
	// 检查指针是否为空
	if f == nil {
		panic("FloatFlag pointer cannot be nil")
	}

	// 参数校验（复用公共函数）
	if validateErr := c.validateFlag(longName, shortName); validateErr != nil {
		panic(validateErr)
	}

	// 显式初始化默认值
	currentFloat := new(float64) // 显式堆分配
	*currentFloat = defValue

	// 初始化Flag对象
	if initErr := f.Init(longName, shortName, defValue, usage, currentFloat); initErr != nil {
		panic(initErr)
	}

	// 创建FlagMeta对象
	meta := &flags.FlagMeta{
		Flag: f, // 添加标志对象 - Flag对象
	}

	// 绑定短标志
	if shortName != "" {
		c.fs.Float64Var(f.GetPointer(), shortName, defValue, usage)
	}

	// 绑定长标志
	if longName != "" {
		c.fs.Float64Var(f.GetPointer(), longName, defValue, usage)
	}

	// 注册Flag对象
	if registerErr := c.flagRegistry.RegisterFlag(meta); registerErr != nil {
		panic(registerErr)
	}
}

// DurationVar 绑定时间间隔类型标志到指针并内部注册Flag对象
//
// 参数依次为: 时间间隔标志指针、长标志名、短标志、默认值、帮助说明
func (c *Cmd) DurationVar(f *flags.DurationFlag, longName, shortName string, defValue time.Duration, usage string) {
	// 检查指针是否为空
	if f == nil {
		panic("DurationFlag pointer cannot be nil")
	}

	// 参数校验
	if validateErr := c.validateFlag(longName, shortName); validateErr != nil {
		panic(validateErr)
	}

	// 初始化默认值(值类型)
	currentDuration := new(time.Duration)
	*currentDuration = defValue

	// 初始化Flag对象
	if initErr := f.Init(longName, shortName, defValue, usage, currentDuration); initErr != nil {
		panic(initErr)
	}

	// 绑定短标志
	if shortName != "" {
		c.fs.DurationVar(f.GetPointer(), shortName, defValue, usage)
	}

	// 绑定长标志
	if longName != "" {
		c.fs.DurationVar(f.GetPointer(), longName, defValue, usage)
	}

	// 创建并注册标志元数据
	meta := &flags.FlagMeta{
		Flag: f, // 添加标志对象 - Flag对象
	}

	// 注册标志元数据
	if registerErr := c.flagRegistry.RegisterFlag(meta); registerErr != nil {
		panic(registerErr)
	}
}

// Duration 添加时间间隔类型标志, 返回标志对象指针
//
// 参数依次为: 长标志名、短标志、默认值、帮助说明
//
// 返回值: 时间间隔标志对象指针
func (c *Cmd) Duration(longName, shortName string, defValue time.Duration, usage string) *flags.DurationFlag {
	f := &flags.DurationFlag{}
	c.DurationVar(f, longName, shortName, defValue, usage)
	return f
}

// Enum 添加枚举类型标志, 返回标志对象指针
//
// 参数依次为: 长标志名、短标志、默认值、帮助说明、限制该标志取值的枚举值切片
//
// 返回值: 枚举标志对象指针
func (c *Cmd) Enum(longName, shortName string, defValue string, usage string, options []string) *flags.EnumFlag {
	f := &flags.EnumFlag{}
	c.EnumVar(f, longName, shortName, defValue, usage, options)
	return f
}

// EnumVar 绑定枚举类型标志到指针并内部注册Flag对象
//
// 参数依次为: 枚举标志指针、长标志名、短标志、默认值、帮助说明、限制该标志取值的枚举值切片
func (c *Cmd) EnumVar(f *flags.EnumFlag, longName, shortName string, defValue string, usage string, options []string) {
	// 检查指针是否为空
	if f == nil {
		panic("EnumFlag pointer cannot be nil")
	}

	// 参数校验（复用公共函数）
	if validateErr := c.validateFlag(longName, shortName); validateErr != nil {
		panic(validateErr)
	}

	// 初始化枚举值
	if options == nil {
		options = make([]string, 0)
	}

	// 调用枚举专用Init方法
	if initErr := f.Init(longName, shortName, defValue, usage, options); initErr != nil {
		panic(initErr)
	}

	// 创建FlagMeta对象
	meta := &flags.FlagMeta{
		Flag: f, // 添加标志对象 - Flag对象
	}

	// 绑定短标志
	if shortName != "" {
		c.fs.StringVar(f.GetPointer(), shortName, defValue, usage)
	}

	// 绑定长标志
	if longName != "" {
		c.fs.StringVar(f.GetPointer(), longName, defValue, usage)
	}

	// 注册Flag对象
	if registerErr := c.flagRegistry.RegisterFlag(meta); registerErr != nil {
		panic(registerErr)
	}
}

// Slice 绑定字符串切片类型标志并内部注册Flag对象
//
// 参数依次为: 长标志名、短标志、默认值、帮助说明
//
// 返回值: 字符串切片标志对象指针
func (c *Cmd) Slice(longName, shortName string, defValue []string, usage string) *flags.SliceFlag {
	f := &flags.SliceFlag{}
	c.SliceVar(f, longName, shortName, defValue, usage)
	return f
}

// SliceVar 绑定字符串切片类型标志到指针并内部注册Flag对象
//
// 参数依次为: 字符串切片标志指针、长标志名、短标志、默认值、帮助说明
func (c *Cmd) SliceVar(f *flags.SliceFlag, longName, shortName string, defValue []string, usage string) {
	// 检查指针是否为空
	if f == nil {
		panic("SliceFlag pointer cannot be nil")
	}

	// 参数校验（复用公共函数）
	if validateErr := c.validateFlag(longName, shortName); validateErr != nil {
		panic(validateErr)
	}

	// 确保默认值不为空
	if defValue == nil {
		defValue = make([]string, 0)
	}

	// 初始化Flag对象字段
	if initErr := f.Init(longName, shortName, defValue, usage); initErr != nil {
		panic(initErr)
	}

	// 创建FlagMeta对象
	meta := &flags.FlagMeta{
		Flag: f, // 添加标志对象 - Flag对象
	}

	// 绑定短标志
	if shortName != "" {
		c.fs.Var(f, shortName, usage)
	}

	// 绑定长标志
	if longName != "" {
		c.fs.Var(f, longName, usage)
	}

	// 注册Flag对象
	if registerErr := c.flagRegistry.RegisterFlag(meta); registerErr != nil {
		panic(registerErr)
	}
}

// Int64Var 绑定64位整数类型标志到指针并内部注册Flag对象
//
// 参数依次为: 64位整数标志指针、长标志名、短标志、默认值、帮助说明
func (c *Cmd) Int64Var(f *flags.Int64Flag, longName, shortName string, defValue int64, usage string) {
	// 检查指针是否为nil
	if f == nil {
		panic("Int64Flag pointer cannot be nil")
	}

	// 参数校验（复用公共函数）
	if validateErr := c.validateFlag(longName, shortName); validateErr != nil {
		panic(validateErr)
	}

	// 初始化默认值
	currentInt64 := new(int64)
	*currentInt64 = defValue

	// 初始化Flag对象
	if initErr := f.Init(longName, shortName, defValue, usage, currentInt64); initErr != nil {
		panic(initErr)
	}

	// 创建FlagMeta对象
	meta := &flags.FlagMeta{
		Flag: f, // 添加标志对象 - Flag对象
	}

	// 绑定短标志
	if shortName != "" {
		c.fs.Int64Var(f.GetPointer(), shortName, defValue, usage)
	}

	// 绑定长标志
	if longName != "" {
		c.fs.Int64Var(f.GetPointer(), longName, defValue, usage)
	}

	// 注册Flag对象
	if registerErr := c.flagRegistry.RegisterFlag(meta); registerErr != nil {
		panic(registerErr)
	}
}

// Int64 添加64位整数类型标志, 返回标志对象指针
//
// 参数依次为: 长标志名、短标志、默认值、帮助说明
//
// 返回值: 64位整数标志对象指针
func (c *Cmd) Int64(longName, shortName string, defValue int64, usage string) *flags.Int64Flag {
	f := &flags.Int64Flag{}
	c.Int64Var(f, longName, shortName, defValue, usage)
	return f
}

// Uint16Var 绑定16位无符号整数类型标志到指针并内部注册Flag对象
//
// 参数依次为: 16位无符号整数标志指针、长标志名、短标志、默认值、帮助说明
func (c *Cmd) Uint16Var(f *flags.Uint16Flag, longName, shortName string, defValue uint16, usage string) {
	if f == nil {
		panic("Uint16Flag pointer cannot be nil")
	}

	if validateErr := c.validateFlag(longName, shortName); validateErr != nil {
		panic(validateErr)
	}

	currentUint16 := new(uint16)
	*currentUint16 = defValue

	if initErr := f.Init(longName, shortName, defValue, usage, currentUint16); initErr != nil {
		panic(initErr)
	}

	meta := &flags.FlagMeta{
		Flag: f,
	}

	if shortName != "" {
		c.fs.Var(f, shortName, usage)
	}

	if longName != "" {
		c.fs.Var(f, longName, usage)
	}

	if registerErr := c.flagRegistry.RegisterFlag(meta); registerErr != nil {
		panic(registerErr)
	}
}

// Uint16 添加16位无符号整数类型标志, 返回标志对象指针
//
// 参数依次为: 长标志名、短标志、默认值、帮助说明
//
// 返回值: 16位无符号整数标志对象指针
func (c *Cmd) Uint16(longName, shortName string, defValue uint16, usage string) *flags.Uint16Flag {
	f := &flags.Uint16Flag{}
	c.Uint16Var(f, longName, shortName, defValue, usage)
	return f
}

// Time 添加时间类型标志, 返回标志对象指针
//
// 参数依次为: 长标志名、短标志、默认值、帮助说明
//
// 返回值: 时间标志对象指针
func (c *Cmd) Time(longName, shortName string, defValue time.Time, usage string) *flags.TimeFlag {
	f := &flags.TimeFlag{}
	c.TimeVar(f, longName, shortName, defValue, usage)
	return f
}

// TimeVar 绑定时间类型标志到指针并内部注册Flag对象
//
// 参数依次为: 时间标志指针、长标志名、短标志、默认值、帮助说明
func (c *Cmd) TimeVar(f *flags.TimeFlag, longName, shortName string, defValue time.Time, usage string) {
	// 检查指针是否为nil
	if f == nil {
		panic("TimeFlag pointer cannot be nil")
	}

	// 参数校验（复用公共函数）
	if validateErr := c.validateFlag(longName, shortName); validateErr != nil {
		panic(validateErr)
	}

	// 初始化默认值
	currentTime := new(time.Time)
	*currentTime = defValue

	// 初始化Flag对象
	if initErr := f.Init(longName, shortName, defValue, usage, currentTime); initErr != nil {
		panic(initErr)
	}

	// 创建FlagMeta对象
	meta := &flags.FlagMeta{
		Flag: f, // 添加标志对象 - Flag对象
	}

	// 绑定短标志
	if shortName != "" {
		c.fs.Var(f, shortName, usage)
	}

	// 绑定长标志
	if longName != "" {
		c.fs.Var(f, longName, usage)
	}

	// 注册Flag对象
	if registerErr := c.flagRegistry.RegisterFlag(meta); registerErr != nil {
		panic(registerErr)
	}
}

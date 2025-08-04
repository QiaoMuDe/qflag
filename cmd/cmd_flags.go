package cmd

import (
	"time"

	"gitee.com/MM-Q/qflag/flags"
)

// =============================================================================
// 布尔类型标志
// =============================================================================

// BoolVar 绑定布尔类型标志到指针并内部注册Flag对象
//
// 参数值:
//   - f: *flags.BoolFlag - 布尔标志对象指针
//   - longName: string - 长标志名
//   - shortName: string - 短标志
//   - defValue: bool - 默认值
//   - usage: string - 帮助说明
func (c *Cmd) BoolVar(f *flags.BoolFlag, longName, shortName string, defValue bool, usage string) {
	c.rwMu.Lock()
	defer c.rwMu.Unlock()

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
	if initErr := f.Init(longName, shortName, usage, currentBool); initErr != nil {
		panic(initErr)
	}

	// 绑定长短标志
	if shortName != "" {
		c.fs.Var(f, shortName, usage)
	}
	if longName != "" {
		c.fs.Var(f, longName, usage)
	}

	// 注册Flag对象
	if registerErr := c.flagRegistry.RegisterFlag(&flags.FlagMeta{Flag: f}); registerErr != nil {
		panic(registerErr)
	}
}

// Bool 添加布尔类型标志, 返回标志对象指针
//
// 参数值:
//   - longName: string - 长标志名
//   - shortName: string - 短标志
//   - defValue: bool - 默认值
//   - usage: string - 帮助说明
//
// 返回值:
//   - *flags.BoolFlag - 布尔标志对象指针
func (c *Cmd) Bool(longName, shortName string, defValue bool, usage string) *flags.BoolFlag {
	f := &flags.BoolFlag{}
	c.BoolVar(f, longName, shortName, defValue, usage)
	return f
}

// =============================================================================
// 枚举类型标志
// =============================================================================

// Enum 添加枚举类型标志, 返回标志对象指针
//
// 参数值:
//   - longName: string - 长标志名
//   - shortName: string - 短标志
//   - defValue: string - 默认值
//   - usage: string - 帮助说明
//   - options: []string - 限制该标志取值的枚举值切片
//
// 返回值:
//   - *flags.EnumFlag - 枚举标志对象指针
func (c *Cmd) Enum(longName, shortName string, defValue string, usage string, options []string) *flags.EnumFlag {
	f := &flags.EnumFlag{}
	c.EnumVar(f, longName, shortName, defValue, usage, options)
	return f
}

// EnumVar 绑定枚举类型标志到指针并内部注册Flag对象
//
// 参数值:
//   - f: *flags.EnumFlag - 枚举标志对象指针
//   - longName: string - 长标志名
//   - shortName: string - 短标志
//   - defValue: string - 默认值
//   - usage: string - 帮助说明
//   - options: []string - 限制该标志取值的枚举值切片
func (c *Cmd) EnumVar(f *flags.EnumFlag, longName, shortName string, defValue string, usage string, options []string) {
	c.rwMu.Lock()
	defer c.rwMu.Unlock()

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

	// 绑定长短标志
	if shortName != "" {
		c.fs.Var(f, shortName, usage)
	}
	if longName != "" {
		c.fs.Var(f, longName, usage)
	}

	// 注册Flag对象
	if registerErr := c.flagRegistry.RegisterFlag(&flags.FlagMeta{Flag: f}); registerErr != nil {
		panic(registerErr)
	}
}

// =============================================================================
// 64位浮点数类型标志
// =============================================================================

// Float64 添加浮点型标志, 返回标志对象指针
//
// 参数值:
//   - longName - 长标志名
//   - shortName - 短标志
//   - defValue - 默认值
//   - usage - 帮助说明
//
// 返回值:
//   - *flags.Float64Flag - 浮点型标志对象指针
func (c *Cmd) Float64(longName, shortName string, defValue float64, usage string) *flags.Float64Flag {
	f := &flags.Float64Flag{}
	c.Float64Var(f, longName, shortName, defValue, usage)
	return f
}

// Float64Var 绑定浮点型标志到指针并内部注册Flag对象
//
// 参数值:
//   - f: *flags.Float64Flag - 浮点型标志对象指针
//   - longName: string - 长标志名
//   - shortName: string - 短标志
//   - defValue: float64 - 默认值
//   - usage: string - 帮助说明
func (c *Cmd) Float64Var(f *flags.Float64Flag, longName, shortName string, defValue float64, usage string) {
	c.rwMu.Lock()
	defer c.rwMu.Unlock()

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
	if initErr := f.Init(longName, shortName, usage, currentFloat); initErr != nil {
		panic(initErr)
	}

	// 绑定长短标志
	if shortName != "" {
		c.fs.Var(f, shortName, usage)
	}
	if longName != "" {
		c.fs.Var(f, longName, usage)
	}

	// 注册Flag对象
	if registerErr := c.flagRegistry.RegisterFlag(&flags.FlagMeta{Flag: f}); registerErr != nil {
		panic(registerErr)
	}
}

// =============================================================================
// 整数类型标志
// =============================================================================

// IntVar 绑定整数类型标志到指针并内部注册Flag对象
//
// 参数值:
//   - f: 整数标志指针
//   - longName: 长标志名
//   - shortName: 短标志名
//   - defValue: 默认值
//   - usage: 帮助说明
func (c *Cmd) IntVar(f *flags.IntFlag, longName, shortName string, defValue int, usage string) {
	c.rwMu.Lock()
	defer c.rwMu.Unlock()

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
	if initErr := f.Init(longName, shortName, usage, currentInt); initErr != nil {
		panic(initErr)
	}

	// 绑定长短标志
	if shortName != "" {
		c.fs.Var(f, shortName, usage)
	}
	if longName != "" {
		c.fs.Var(f, longName, usage)
	}

	// 注册Flag对象
	if registerErr := c.flagRegistry.RegisterFlag(&flags.FlagMeta{Flag: f}); registerErr != nil {
		panic(registerErr)
	}
}

// Int 添加整数类型标志, 返回标志对象指针
//
// 参数值:
//   - longName: 长标志名
//   - shortName: 短标志名
//   - defValue: 默认值
//   - usage: 帮助说明
//
// 返回值:
//   - *flags.IntFlag: 整数标志对象指针
func (c *Cmd) Int(longName, shortName string, defValue int, usage string) *flags.IntFlag {
	f := &flags.IntFlag{}
	c.IntVar(f, longName, shortName, defValue, usage)
	return f
}

// Int64Var 绑定64位整数类型标志到指针并内部注册Flag对象
//
// 参数值:
//   - f: 64位整数标志指针
//   - longName: 长标志名
//   - shortName: 短标志名
//   - defValue: 默认值
//   - usage: 帮助说明
func (c *Cmd) Int64Var(f *flags.Int64Flag, longName, shortName string, defValue int64, usage string) {
	c.rwMu.Lock()
	defer c.rwMu.Unlock()

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
	if initErr := f.Init(longName, shortName, usage, currentInt64); initErr != nil {
		panic(initErr)
	}

	// 绑定长短标志
	if shortName != "" {
		c.fs.Var(f, shortName, usage)
	}
	if longName != "" {
		c.fs.Var(f, longName, usage)
	}

	// 注册Flag对象
	if registerErr := c.flagRegistry.RegisterFlag(&flags.FlagMeta{Flag: f}); registerErr != nil {
		panic(registerErr)
	}
}

// Int64 添加64位整数类型标志, 返回标志对象指针
//
// 参数值:
//   - longName: 长标志名
//   - shortName: 短标志名
//   - defValue: 默认值
//   - usage: 帮助说明
//
// 返回值:
//   - *flags.Int64Flag: 64位整数标志对象指针
func (c *Cmd) Int64(longName, shortName string, defValue int64, usage string) *flags.Int64Flag {
	f := &flags.Int64Flag{}
	c.Int64Var(f, longName, shortName, defValue, usage)
	return f
}

// =============================================================================
// IP类型标志
// =============================================================================

// IP4Var 绑定IPv4地址类型标志到指针并内部注册Flag对象
//
// 参数值:
//   - f: IPv4标志指针
//   - longName: 长标志名
//   - shortName: 短标志名
//   - defValue: 默认值
//   - usage: 帮助说明
func (c *Cmd) IP4Var(f *flags.IP4Flag, longName, shortName string, defValue string, usage string) {
	c.rwMu.Lock()
	defer c.rwMu.Unlock()

	// 校验参数
	if f == nil {
		panic("IP4Flag pointer cannot be nil")
	}

	// 通用参数校验
	if validateErr := c.validateFlag(longName, shortName); validateErr != nil {
		panic(validateErr)
	}

	// 显式设置默认值
	currentIP4 := new(string)
	*currentIP4 = defValue

	// 初始化标志
	if initErr := f.Init(longName, shortName, usage, currentIP4); initErr != nil {
		panic(initErr)
	}

	// 绑定标志
	if shortName != "" {
		c.fs.Var(f, shortName, usage)
	}
	if longName != "" {
		c.fs.Var(f, longName, usage)
	}

	// 注册Flag对象
	if registerErr := c.flagRegistry.RegisterFlag(&flags.FlagMeta{Flag: f}); registerErr != nil {
		panic(registerErr)
	}
}

// IP4 添加IPv4地址类型标志, 返回标志对象指针
//
// 参数值:
//   - longName: 长标志名
//   - shortName: 短标志名
//   - defValue: 默认值
//   - usage: 帮助说明
//
// 返回值:
//   - *flags.IP4Flag: IPv4地址标志对象指针
func (c *Cmd) IP4(longName, shortName string, defValue string, usage string) *flags.IP4Flag {
	f := &flags.IP4Flag{}
	c.IP4Var(f, longName, shortName, defValue, usage)
	return f
}

// IP6Var 绑定IPv6地址类型标志到指针并内部注册Flag对象
//
// 参数值:
//   - f: IPv6标志指针
//   - longName: 长标志名
//   - shortName: 短标志名
//   - defValue: 默认值
//   - usage: 帮助说明
func (c *Cmd) IP6Var(f *flags.IP6Flag, longName, shortName string, defValue string, usage string) {
	c.rwMu.Lock()
	defer c.rwMu.Unlock()

	// 校验参数
	if f == nil {
		panic("IP6Flag pointer cannot be nil")
	}

	// 通用校验
	if validateErr := c.validateFlag(longName, shortName); validateErr != nil {
		panic(validateErr)
	}

	// 显式设置默认值
	currentIP6 := new(string)
	*currentIP6 = defValue

	// 初始化标志对象
	if initErr := f.Init(longName, shortName, usage, currentIP6); initErr != nil {
		panic(initErr)
	}

	// 绑定标志
	if shortName != "" {
		c.fs.Var(f, shortName, usage)
	}
	if longName != "" {
		c.fs.Var(f, longName, usage)
	}

	// 注册Flag对象
	if registerErr := c.flagRegistry.RegisterFlag(&flags.FlagMeta{Flag: f}); registerErr != nil {
		panic(registerErr)
	}
}

// IP6 添加IPv6地址类型标志, 返回标志对象指针
//
// 参数值:
//   - longName: 长标志名
//   - shortName: 短标志名
//   - defValue: 默认值
//   - usage: 帮助说明
//
// 返回值:
//   - *flags.IP6Flag: IPv6地址标志对象指针
func (c *Cmd) IP6(longName, shortName string, defValue string, usage string) *flags.IP6Flag {
	f := &flags.IP6Flag{}
	c.IP6Var(f, longName, shortName, defValue, usage)
	return f
}

// =============================================================================
// 键值对类型标志
// =============================================================================

// MapVar 绑定键值对类型标志到指针并内部注册Flag对象
//
// 参数值:
//   - f: 键值对标志指针
//   - longName: 长标志名
//   - shortName: 短标志名
//   - defValue: 默认值
//   - usage: 帮助说明
func (c *Cmd) MapVar(f *flags.MapFlag, longName, shortName string, defValue map[string]string, usage string) {
	c.rwMu.Lock()
	defer c.rwMu.Unlock()

	// 检查指针是否为nil
	if f == nil {
		panic("MapFlag pointer cannot be nil")
	}

	// 参数校验（复用公共函数）
	if validateErr := c.validateFlag(longName, shortName); validateErr != nil {
		panic(validateErr)
	}

	// 如果默认值为nil，则初始化为空map
	if defValue == nil {
		defValue = map[string]string{}
	}

	// 初始化值
	currentMap := new(map[string]string)
	*currentMap = defValue

	// 初始化Flag对象
	if initErr := f.Init(longName, shortName, usage, currentMap); initErr != nil {
		panic(initErr)
	}

	// 设置默认分隔符
	f.SetDelimiters(flags.FlagSplitComma, flags.FlagKVEqual)

	// 绑定短标志
	if shortName != "" {
		c.fs.Var(f, shortName, usage)
	}

	// 绑定长标志
	if longName != "" {
		c.fs.Var(f, longName, usage)
	}

	// 注册Flag对象
	if registerErr := c.flagRegistry.RegisterFlag(&flags.FlagMeta{Flag: f}); registerErr != nil {
		panic(registerErr)
	}
}

// Map 添加键值对类型标志, 返回标志对象指针
//
// 参数值:
//   - longName: 长标志名
//   - shortName: 短标志名
//   - defValue: 默认值
//   - usage: 帮助说明
//
// 返回值:
//   - *flags.MapFlag: 键值对标志对象指针
func (c *Cmd) Map(longName, shortName string, defValue map[string]string, usage string) *flags.MapFlag {
	f := &flags.MapFlag{}
	c.MapVar(f, longName, shortName, defValue, usage)
	return f
}

// =============================================================================
// 切片类型标志
// =============================================================================

// Slice 绑定字符串切片类型标志并内部注册Flag对象
//
// 参数值:
//   - longName: 长标志名
//   - shortName: 短标志名
//   - defValue: 默认值
//   - usage: 帮助说明
//
// 返回值:
//   - *flags.SliceFlag: 字符串切片标志对象指针
func (c *Cmd) Slice(longName, shortName string, defValue []string, usage string) *flags.SliceFlag {
	f := &flags.SliceFlag{}
	c.SliceVar(f, longName, shortName, defValue, usage)
	return f
}

// SliceVar 绑定字符串切片类型标志到指针并内部注册Flag对象
//
// 参数值:
//   - f: 字符串切片标志指针
//   - longName: 长标志名
//   - shortName: 短标志名
//   - defValue: 默认值
//   - usage: 帮助说明
func (c *Cmd) SliceVar(f *flags.SliceFlag, longName, shortName string, defValue []string, usage string) {
	c.rwMu.Lock()
	defer c.rwMu.Unlock()

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

	// 绑定短标志
	if shortName != "" {
		c.fs.Var(f, shortName, usage)
	}

	// 绑定长标志
	if longName != "" {
		c.fs.Var(f, longName, usage)
	}

	// 注册Flag对象
	if registerErr := c.flagRegistry.RegisterFlag(&flags.FlagMeta{Flag: f}); registerErr != nil {
		panic(registerErr)
	}
}

// =============================================================================
// 字符串类型标志
// =============================================================================

// String 添加字符串类型标志, 返回标志对象指针
//
// 参数值:
//   - longName: 长标志名
//   - shortName: 短标志名
//   - defValue: 默认值
//   - usage: 帮助说明
//
// 返回值:
//   - *flags.StringFlag: 字符串标志对象指针
func (c *Cmd) String(longName, shortName, defValue, usage string) *flags.StringFlag {
	f := &flags.StringFlag{}
	c.StringVar(f, longName, shortName, defValue, usage)
	return f
}

// StringVar 绑定字符串类型标志到指针并内部注册Flag对象
//
// 参数值:
//   - f: 字符串标志指针
//   - longName: 长标志名
//   - shortName: 短标志名
//   - defValue: 默认值
//   - usage: 帮助说明
func (c *Cmd) StringVar(f *flags.StringFlag, longName, shortName, defValue, usage string) {
	c.rwMu.Lock()
	defer c.rwMu.Unlock()

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
	if initErr := f.Init(longName, shortName, usage, currentStr); initErr != nil {
		panic(initErr)
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
	if registerErr := c.flagRegistry.RegisterFlag(&flags.FlagMeta{Flag: f}); registerErr != nil {
		panic(registerErr)
	}
}

// =============================================================================
// 时间类型标志
// =============================================================================

// Time 添加时间类型标志, 返回标志对象指针
//
// 参数值:
//   - longName: 长标志名
//   - shortName: 短标志名
//   - defValue: 默认值
//   - usage: 帮助说明
//
// 返回值:
//   - *flags.TimeFlag: 时间标志对象指针
func (c *Cmd) Time(longName, shortName string, defValue time.Time, usage string) *flags.TimeFlag {
	f := &flags.TimeFlag{}
	c.TimeVar(f, longName, shortName, defValue, usage)
	return f
}

// TimeVar 绑定时间类型标志到指针并内部注册Flag对象
//
// 参数值:
//   - f: 时间标志指针
//   - longName: 长标志名
//   - shortName: 短标志名
//   - defValue: 默认值
//   - usage: 帮助说明
func (c *Cmd) TimeVar(f *flags.TimeFlag, longName, shortName string, defValue time.Time, usage string) {
	c.rwMu.Lock()
	defer c.rwMu.Unlock()

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
	if initErr := f.Init(longName, shortName, usage, currentTime); initErr != nil {
		panic(initErr)
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
	if registerErr := c.flagRegistry.RegisterFlag(&flags.FlagMeta{Flag: f}); registerErr != nil {
		panic(registerErr)
	}
}

// DurationVar 绑定时间间隔类型标志到指针并内部注册Flag对象
//
// 参数值:
//   - f: *flags.DurationFlag - 时间间隔标志对象指针
//   - longName: string - 长标志名
//   - shortName: string - 短标志
//   - defValue: time.Duration - 默认值
//   - usage: string - 帮助说明
func (c *Cmd) DurationVar(f *flags.DurationFlag, longName, shortName string, defValue time.Duration, usage string) {
	c.rwMu.Lock()
	defer c.rwMu.Unlock()

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
	if initErr := f.Init(longName, shortName, usage, currentDuration); initErr != nil {
		panic(initErr)
	}

	// 绑定长短标志
	if shortName != "" {
		c.fs.Var(f, shortName, usage)
	}
	if longName != "" {
		c.fs.Var(f, longName, usage)
	}

	// 注册标志元数据
	if registerErr := c.flagRegistry.RegisterFlag(&flags.FlagMeta{Flag: f}); registerErr != nil {
		panic(registerErr)
	}
}

// Duration 添加时间间隔类型标志, 返回标志对象指针
//
// 参数值:
//   - longName: string - 长标志名
//   - shortName: string - 短标志
//   - defValue: time.Duration - 默认值
//   - usage: string - 帮助说明
//
// 返回值:
//   - *flags.DurationFlag - 时间间隔标志对象指针
func (c *Cmd) Duration(longName, shortName string, defValue time.Duration, usage string) *flags.DurationFlag {
	f := &flags.DurationFlag{}
	c.DurationVar(f, longName, shortName, defValue, usage)
	return f
}

// =============================================================================
// 无符号整数类型标志
// =============================================================================

// Uint16Var 绑定16位无符号整数类型标志到指针并内部注册Flag对象
//
// 参数值:
//   - f: 16位无符号整数标志指针
//   - longName: 长标志名
//   - shortName: 短标志名
//   - defValue: 默认值
//   - usage: 帮助说明
func (c *Cmd) Uint16Var(f *flags.Uint16Flag, longName, shortName string, defValue uint16, usage string) {
	c.rwMu.Lock()
	defer c.rwMu.Unlock()

	// 参数校验
	if f == nil {
		panic("Uint16Flag pointer cannot be nil")
	}

	// 通用校验
	if validateErr := c.validateFlag(longName, shortName); validateErr != nil {
		panic(validateErr)
	}

	// 显式初始化
	currentUint16 := new(uint16)
	*currentUint16 = defValue

	// 初始化Flag对象
	if initErr := f.Init(longName, shortName, usage, currentUint16); initErr != nil {
		panic(initErr)
	}

	// 注册到flagSet
	if shortName != "" {
		c.fs.Var(f, shortName, usage)
	}
	if longName != "" {
		c.fs.Var(f, longName, usage)
	}

	// 注册到flagRegistry
	if registerErr := c.flagRegistry.RegisterFlag(&flags.FlagMeta{Flag: f}); registerErr != nil {
		panic(registerErr)
	}
}

// Uint16 添加16位无符号整数类型标志, 返回标志对象指针
//
// 参数值:
//   - longName: 长标志名
//   - shortName: 短标志名
//   - defValue: 默认值
//   - usage: 帮助说明
//
// 返回值:
//   - *flags.Uint16Flag: 16位无符号整数标志对象指针
func (c *Cmd) Uint16(longName, shortName string, defValue uint16, usage string) *flags.Uint16Flag {
	f := &flags.Uint16Flag{}
	c.Uint16Var(f, longName, shortName, defValue, usage)
	return f
}

// Uint32Var 绑定32位无符号整数类型标志到指针并内部注册Flag对象
//
// 参数值:
//   - f: 32位无符号整数标志指针
//   - longName: 长标志名
//   - shortName: 短标志名
//   - defValue: 默认值
//   - usage: 帮助说明
func (c *Cmd) Uint32Var(f *flags.Uint32Flag, longName, shortName string, defValue uint32, usage string) {
	c.rwMu.Lock()
	defer c.rwMu.Unlock()

	// 校验标志指针
	if f == nil {
		panic("Uint32Flag pointer cannot be nil")
	}

	// 通用标志校验
	if validateErr := c.validateFlag(longName, shortName); validateErr != nil {
		panic(validateErr)
	}

	// 设置标志指针
	currentUint32 := new(uint32)
	*currentUint32 = defValue

	// 绑定默认值
	if initErr := f.Init(longName, shortName, usage, currentUint32); initErr != nil {
		panic(initErr)
	}

	// 绑定标志到指针
	if shortName != "" {
		c.fs.Var(f, shortName, usage)
	}
	if longName != "" {
		c.fs.Var(f, longName, usage)
	}

	if registerErr := c.flagRegistry.RegisterFlag(&flags.FlagMeta{Flag: f}); registerErr != nil {
		panic(registerErr)
	}
}

// Uint32 添加32位无符号整数类型标志, 返回标志对象指针
//
// 参数值:
//   - longName: 长标志名
//   - shortName: 短标志名
//   - defValue: 默认值
//   - usage: 帮助说明
//
// 返回值:
//   - *flags.Uint32Flag: 32位无符号整数标志对象指针
func (c *Cmd) Uint32(longName, shortName string, defValue uint32, usage string) *flags.Uint32Flag {
	f := &flags.Uint32Flag{}
	c.Uint32Var(f, longName, shortName, defValue, usage)
	return f
}

// Uint64Var 绑定64位无符号整数类型标志到指针并内部注册Flag对象
//
// 参数值:
//   - f: 64位无符号整数标志指针
//   - longName: 长标志名
//   - shortName: 短标志名
//   - defValue: 默认值
//   - usage: 帮助说明
func (c *Cmd) Uint64Var(f *flags.Uint64Flag, longName, shortName string, defValue uint64, usage string) {
	c.rwMu.Lock()
	defer c.rwMu.Unlock()

	// 检查指针是否为空
	if f == nil {
		panic("Uint64Flag pointer cannot be nil")
	}

	// 参数校验（复用公共函数）
	if validateErr := c.validateFlag(longName, shortName); validateErr != nil {
		panic(validateErr)
	}

	// 显式初始化
	currentUint64 := new(uint64)
	*currentUint64 = defValue

	// 注册flag
	if initErr := f.Init(longName, shortName, usage, currentUint64); initErr != nil {
		panic(initErr)
	}

	// 注册到flagSet
	if shortName != "" {
		c.fs.Var(f, shortName, usage)
	}
	if longName != "" {
		c.fs.Var(f, longName, usage)
	}

	// 注册到flagRegistry
	if registerErr := c.flagRegistry.RegisterFlag(&flags.FlagMeta{Flag: f}); registerErr != nil {
		panic(registerErr)
	}
}

// Uint64 添加64位无符号整数类型标志, 返回标志对象指针
//
// 参数值:
//   - longName: 长标志名
//   - shortName: 短标志名
//   - defValue: 默认值
//   - usage: 帮助说明
//
// 返回值:
//   - *flags.Uint64Flag: 64位无符号整数标志对象指针
func (c *Cmd) Uint64(longName, shortName string, defValue uint64, usage string) *flags.Uint64Flag {
	f := &flags.Uint64Flag{}
	c.Uint64Var(f, longName, shortName, defValue, usage)
	return f
}

// =============================================================================
// URL类型标志
// =============================================================================

// URLVar 绑定URL类型标志到指针并内部注册Flag对象
//
// 参数值:
//   - f: URL标志对象指针
//   - longName: 长标志名
//   - shortName: 短标志名
//   - defValue: 默认值
//   - usage: 帮助说明
func (c *Cmd) URLVar(f *flags.URLFlag, longName, shortName string, defValue string, usage string) {
	c.rwMu.Lock()
	defer c.rwMu.Unlock()

	// 检查指针是否为空
	if f == nil {
		panic("URLFlag pointer cannot be nil")
	}

	// 参数校验（复用公共函数）
	if validateErr := c.validateFlag(longName, shortName); validateErr != nil {
		panic(validateErr)
	}

	// 显式初始化标志值
	currentURL := new(string)
	*currentURL = defValue

	// 初始化Flag对象
	if initErr := f.Init(longName, shortName, usage, currentURL); initErr != nil {
		panic(initErr)
	}

	// 注册标志
	if shortName != "" {
		c.fs.Var(f, shortName, usage)
	}
	if longName != "" {
		c.fs.Var(f, longName, usage)
	}

	// 注册Flag对象
	if registerErr := c.flagRegistry.RegisterFlag(&flags.FlagMeta{Flag: f}); registerErr != nil {
		panic(registerErr)
	}
}

// URL 添加URL类型标志, 返回标志对象指针
//
// 参数值:
//   - longName: 长标志名
//   - shortName: 短标志名
//   - defValue: 默认值
//   - usage: 帮助说明
//
// 返回值:
//   - *flags.URLFlag: URL标志对象指针
func (c *Cmd) URL(longName, shortName string, defValue string, usage string) *flags.URLFlag {
	f := &flags.URLFlag{}
	c.URLVar(f, longName, shortName, defValue, usage)
	return f
}

package cmd

import (
	"fmt"

	"gitee.com/MM-Q/qflag/flags"
)

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
	c.ctx.Mutex.Lock()
	defer c.ctx.Mutex.Unlock()

	// 检查指针是否为空
	if f == nil {
		panic("FloatFlag pointer cannot be nil")
	}

	// 检查标志是否为内置标志
	if ok := c.ctx.BuiltinFlags.IsBuiltinFlag(longName); ok {
		panic(fmt.Errorf("flag long name %s is reserved", longName))
	}
	if ok := c.ctx.BuiltinFlags.IsBuiltinFlag(shortName); ok {
		panic(fmt.Errorf("flag short name %s is reserved", shortName))
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
		c.ctx.FlagSet.Var(f, shortName, usage)
	}
	if longName != "" {
		c.ctx.FlagSet.Var(f, longName, usage)
	}

	// 注册Flag对象
	if registerErr := c.ctx.FlagRegistry.RegisterFlag(&flags.FlagMeta{Flag: f}); registerErr != nil {
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
	c.ctx.Mutex.Lock()
	defer c.ctx.Mutex.Unlock()

	// 检查指针是否为nil
	if f == nil {
		panic("IntFlag pointer cannot be nil")
	}

	// 检查标志是否为内置标志
	if ok := c.ctx.BuiltinFlags.IsBuiltinFlag(longName); ok {
		panic(fmt.Errorf("flag long name %s is reserved", longName))
	}
	if ok := c.ctx.BuiltinFlags.IsBuiltinFlag(shortName); ok {
		panic(fmt.Errorf("flag short name %s is reserved", shortName))
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
		c.ctx.FlagSet.Var(f, shortName, usage)
	}
	if longName != "" {
		c.ctx.FlagSet.Var(f, longName, usage)
	}

	// 注册Flag对象
	if registerErr := c.ctx.FlagRegistry.RegisterFlag(&flags.FlagMeta{Flag: f}); registerErr != nil {
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
	c.ctx.Mutex.Lock()
	defer c.ctx.Mutex.Unlock()

	// 检查指针是否为nil
	if f == nil {
		panic("Int64Flag pointer cannot be nil")
	}

	// 检查标志是否为内置标志
	if ok := c.ctx.BuiltinFlags.IsBuiltinFlag(longName); ok {
		panic(fmt.Errorf("flag long name %s is reserved", longName))
	}
	if ok := c.ctx.BuiltinFlags.IsBuiltinFlag(shortName); ok {
		panic(fmt.Errorf("flag short name %s is reserved", shortName))
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
		c.ctx.FlagSet.Var(f, shortName, usage)
	}
	if longName != "" {
		c.ctx.FlagSet.Var(f, longName, usage)
	}

	// 注册Flag对象
	if registerErr := c.ctx.FlagRegistry.RegisterFlag(&flags.FlagMeta{Flag: f}); registerErr != nil {
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
	c.ctx.Mutex.Lock()
	defer c.ctx.Mutex.Unlock()

	// 参数校验
	if f == nil {
		panic("Uint16Flag pointer cannot be nil")
	}

	// 检查标志是否为内置标志
	if ok := c.ctx.BuiltinFlags.IsBuiltinFlag(longName); ok {
		panic(fmt.Errorf("flag long name %s is reserved", longName))
	}
	if ok := c.ctx.BuiltinFlags.IsBuiltinFlag(shortName); ok {
		panic(fmt.Errorf("flag short name %s is reserved", shortName))
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
		c.ctx.FlagSet.Var(f, shortName, usage)
	}
	if longName != "" {
		c.ctx.FlagSet.Var(f, longName, usage)
	}

	// 注册到flagRegistry
	if registerErr := c.ctx.FlagRegistry.RegisterFlag(&flags.FlagMeta{Flag: f}); registerErr != nil {
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
	c.ctx.Mutex.Lock()
	defer c.ctx.Mutex.Unlock()

	// 校验标志指针
	if f == nil {
		panic("Uint32Flag pointer cannot be nil")
	}

	// 检查标志是否为内置标志
	if ok := c.ctx.BuiltinFlags.IsBuiltinFlag(longName); ok {
		panic(fmt.Errorf("flag long name %s is reserved", longName))
	}
	if ok := c.ctx.BuiltinFlags.IsBuiltinFlag(shortName); ok {
		panic(fmt.Errorf("flag short name %s is reserved", shortName))
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
		c.ctx.FlagSet.Var(f, shortName, usage)
	}
	if longName != "" {
		c.ctx.FlagSet.Var(f, longName, usage)
	}

	if registerErr := c.ctx.FlagRegistry.RegisterFlag(&flags.FlagMeta{Flag: f}); registerErr != nil {
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
	c.ctx.Mutex.Lock()
	defer c.ctx.Mutex.Unlock()

	// 检查指针是否为空
	if f == nil {
		panic("Uint64Flag pointer cannot be nil")
	}

	// 检查标志是否为内置标志
	if ok := c.ctx.BuiltinFlags.IsBuiltinFlag(longName); ok {
		panic(fmt.Errorf("flag long name %s is reserved", longName))
	}
	if ok := c.ctx.BuiltinFlags.IsBuiltinFlag(shortName); ok {
		panic(fmt.Errorf("flag short name %s is reserved", shortName))
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
		c.ctx.FlagSet.Var(f, shortName, usage)
	}
	if longName != "" {
		c.ctx.FlagSet.Var(f, longName, usage)
	}

	// 注册到flagRegistry
	if registerErr := c.ctx.FlagRegistry.RegisterFlag(&flags.FlagMeta{Flag: f}); registerErr != nil {
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

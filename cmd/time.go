package cmd

import (
	"fmt"
	"time"

	"gitee.com/MM-Q/qflag/flags"
)

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
	c.ctx.Mutex.Lock()
	defer c.ctx.Mutex.Unlock()

	// 检查指针是否为nil
	if f == nil {
		panic("TimeFlag pointer cannot be nil")
	}

	// 检查标志是否为内置标志
	if ok := c.ctx.BuiltinFlags.IsBuiltinFlag(longName); ok {
		panic(fmt.Errorf("flag long name %s is reserved", longName))
	}
	if ok := c.ctx.BuiltinFlags.IsBuiltinFlag(shortName); ok {
		panic(fmt.Errorf("flag short name %s is reserved", shortName))
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
		c.ctx.FlagSet.Var(f, shortName, usage)
	}

	// 绑定长标志
	if longName != "" {
		c.ctx.FlagSet.Var(f, longName, usage)
	}

	// 注册Flag对象
	if registerErr := c.ctx.FlagRegistry.RegisterFlag(&flags.FlagMeta{Flag: f}); registerErr != nil {
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
	c.ctx.Mutex.Lock()
	defer c.ctx.Mutex.Unlock()

	// 检查指针是否为空
	if f == nil {
		panic("DurationFlag pointer cannot be nil")
	}

	// 检查标志是否为内置标志
	if ok := c.ctx.BuiltinFlags.IsBuiltinFlag(longName); ok {
		panic(fmt.Errorf("flag long name %s is reserved", longName))
	}
	if ok := c.ctx.BuiltinFlags.IsBuiltinFlag(shortName); ok {
		panic(fmt.Errorf("flag short name %s is reserved", shortName))
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
		c.ctx.FlagSet.Var(f, shortName, usage)
	}
	if longName != "" {
		c.ctx.FlagSet.Var(f, longName, usage)
	}

	// 注册标志元数据
	if registerErr := c.ctx.FlagRegistry.RegisterFlag(&flags.FlagMeta{Flag: f}); registerErr != nil {
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

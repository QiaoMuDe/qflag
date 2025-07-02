package cmd

import "gitee.com/MM-Q/qflag/flags"

// Uint64Var 绑定64位无符号整数类型标志到指针并内部注册Flag对象
// 参数依次为: 64位无符号整数标志指针、长标志名、短标志、默认值、帮助说明
func (c *Cmd) Uint64Var(f *flags.Uint64Flag, longName, shortName string, defValue uint64, usage string) {
	c.rwMu.Lock()
	defer c.rwMu.Unlock()

	if f == nil {
		panic("Uint64Flag pointer cannot be nil")
	}

	if validateErr := c.validateFlag(longName, shortName); validateErr != nil {
		panic(validateErr)
	}

	currentUint64 := new(uint64)
	*currentUint64 = defValue

	if initErr := f.Init(longName, shortName, usage, currentUint64); initErr != nil {
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

// Uint64 添加64位无符号整数类型标志, 返回标志对象指针
// 参数依次为: 长标志名、短标志、默认值、帮助说明
// 返回值: 64位无符号整数标志对象指针
func (c *Cmd) Uint64(longName, shortName string, defValue uint64, usage string) *flags.Uint64Flag {
	f := &flags.Uint64Flag{}
	c.Uint64Var(f, longName, shortName, defValue, usage)
	return f
}

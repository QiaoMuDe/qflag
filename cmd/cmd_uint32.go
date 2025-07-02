package cmd

import "gitee.com/MM-Q/qflag/flags"

// Uint32Var 绑定32位无符号整数类型标志到指针并内部注册Flag对象
// 参数依次为: 32位无符号整数标志指针、长标志名、短标志、默认值、帮助说明
func (c *Cmd) Uint32Var(f *flags.Uint32Flag, longName, shortName string, defValue uint32, usage string) {
	c.rwMu.Lock()
	defer c.rwMu.Unlock()

	if f == nil {
		panic("Uint32Flag pointer cannot be nil")
	}

	if validateErr := c.validateFlag(longName, shortName); validateErr != nil {
		panic(validateErr)
	}

	currentUint32 := new(uint32)
	*currentUint32 = defValue

	if initErr := f.Init(longName, shortName, usage, currentUint32); initErr != nil {
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

// Uint32 添加32位无符号整数类型标志, 返回标志对象指针
// 参数依次为: 长标志名、短标志、默认值、帮助说明
// 返回值: 32位无符号整数标志对象指针
func (c *Cmd) Uint32(longName, shortName string, defValue uint32, usage string) *flags.Uint32Flag {
	f := &flags.Uint32Flag{}
	c.Uint32Var(f, longName, shortName, defValue, usage)
	return f
}

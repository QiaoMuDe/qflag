package cmd

import "gitee.com/MM-Q/qflag/flags"

// IP4Var 绑定IPv4地址类型标志到指针并内部注册Flag对象
// 参数依次为: IPv4标志指针、长标志名、短标志、默认值、帮助说明
func (c *Cmd) IP4Var(f *flags.IP4Flag, longName, shortName string, defValue string, usage string) {
	c.rwMu.Lock()
	defer c.rwMu.Unlock()

	if f == nil {
		panic("IP4Flag pointer cannot be nil")
	}

	if validateErr := c.validateFlag(longName, shortName); validateErr != nil {
		panic(validateErr)
	}

	currentIP4 := new(string)
	*currentIP4 = defValue

	if initErr := f.Init(longName, shortName, usage, currentIP4); initErr != nil {
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

// IP4 添加IPv4地址类型标志, 返回标志对象指针
// 参数依次为: 长标志名、短标志、默认值、帮助说明
// 返回值: IPv4地址标志对象指针
func (c *Cmd) IP4(longName, shortName string, defValue string, usage string) *flags.IP4Flag {
	f := &flags.IP4Flag{}
	c.IP4Var(f, longName, shortName, defValue, usage)
	return f
}

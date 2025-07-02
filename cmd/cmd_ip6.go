package cmd

import "gitee.com/MM-Q/qflag/flags"

// IP6Var 绑定IPv6地址类型标志到指针并内部注册Flag对象
// 参数依次为: IPv6标志指针、长标志名、短标志、默认值、帮助说明
func (c *Cmd) IP6Var(f *flags.IP6Flag, longName, shortName string, defValue string, usage string) {
	c.rwMu.Lock()
	defer c.rwMu.Unlock()

	if f == nil {
		panic("IP6Flag pointer cannot be nil")
	}

	if validateErr := c.validateFlag(longName, shortName); validateErr != nil {
		panic(validateErr)
	}

	currentIP6 := new(string)
	*currentIP6 = defValue

	if initErr := f.Init(longName, shortName, usage, currentIP6); initErr != nil {
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

// IP6 添加IPv6地址类型标志, 返回标志对象指针
// 参数依次为: 长标志名、短标志、默认值、帮助说明
// 返回值: IPv6地址标志对象指针
func (c *Cmd) IP6(longName, shortName string, defValue string, usage string) *flags.IP6Flag {
	f := &flags.IP6Flag{}
	c.IP6Var(f, longName, shortName, defValue, usage)
	return f
}

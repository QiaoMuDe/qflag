package cmd

import "gitee.com/MM-Q/qflag/flags"

// URLVar 绑定URL类型标志到指针并内部注册Flag对象
// 参数依次为: URL标志指针、长标志名、短标志、默认值、帮助说明
func (c *Cmd) URLVar(f *flags.URLFlag, longName, shortName string, defValue string, usage string) {
	c.rwMu.Lock()
	defer c.rwMu.Unlock()

	if f == nil {
		panic("URLFlag pointer cannot be nil")
	}

	if validateErr := c.validateFlag(longName, shortName); validateErr != nil {
		panic(validateErr)
	}

	currentURL := new(string)
	*currentURL = defValue

	if initErr := f.Init(longName, shortName, usage, currentURL); initErr != nil {
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

// URL 添加URL类型标志, 返回标志对象指针
// 参数依次为: 长标志名、短标志、默认值、帮助说明
// 返回值: URL标志对象指针
func (c *Cmd) URL(longName, shortName string, defValue string, usage string) *flags.URLFlag {
	f := &flags.URLFlag{}
	c.URLVar(f, longName, shortName, defValue, usage)
	return f
}

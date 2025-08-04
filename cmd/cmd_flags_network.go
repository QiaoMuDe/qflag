package cmd

import "gitee.com/MM-Q/qflag/flags"

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

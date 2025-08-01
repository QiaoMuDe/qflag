package cmd

import "gitee.com/MM-Q/qflag/flags"

// PathVar 绑定路径类型标志到指针并内部注册Flag对象
//
// 参数值:
//   - f: 路径标志指针
//   - longName: 长标志名
//   - shortName: 短标志名
//   - defValue: 默认值
//   - usage: 帮助说明
func (c *Cmd) PathVar(f *flags.PathFlag, longName, shortName string, defValue string, usage string) {
	c.rwMu.Lock()
	defer c.rwMu.Unlock()

	// 检查指针是否为nil
	if f == nil {
		panic("PathFlag pointer cannot be nil")
	}

	// 参数校验（复用公共函数）
	if validateErr := c.validateFlag(longName, shortName); validateErr != nil {
		panic(validateErr)
	}

	// 初始化Flag对象
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

// Path 添加路径类型标志, 返回标志对象指针
//
// 参数值:
//   - longName: 长标志名
//   - shortName: 短标志名
//   - defValue: 默认值
//   - usage: 帮助说明
//
// 返回值:
//   - *flags.PathFlag: 路径标志对象指针
func (c *Cmd) Path(longName, shortName string, defValue string, usage string) *flags.PathFlag {
	f := &flags.PathFlag{}
	c.PathVar(f, longName, shortName, defValue, usage)
	return f
}

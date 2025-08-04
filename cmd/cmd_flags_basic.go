package cmd

import "gitee.com/MM-Q/qflag/flags"

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

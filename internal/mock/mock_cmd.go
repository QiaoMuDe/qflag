package mock

import (
	"gitee.com/MM-Q/qflag/internal/types"
)

// MockCommandBasic 简化的模拟命令实现
type MockCommandBasic struct {
	name         string
	shortName    string
	description  string
	version      string
	useChinese   bool
	envPrefix    string
	usageSyntax  string
	logoText     string
	args         []string
	parsed       bool
	runFunc      func(types.Command) error
	flagRegistry types.FlagRegistry
	cmdRegistry  types.CmdRegistry
}

// NewMockCommandBasic 创建基础模拟命令
func NewMockCommandBasic(name, shortName, description string) *MockCommandBasic {
	return &MockCommandBasic{
		name:         name,
		shortName:    shortName,
		description:  description,
		version:      "",
		useChinese:   false,
		envPrefix:    "",
		usageSyntax:  "",
		logoText:     "",
		args:         []string{},
		parsed:       false,
		runFunc:      nil,
		flagRegistry: NewMockFlagRegistry(),
		cmdRegistry:  NewMockCmdRegistry(),
	}
}

// 实现 Command 接口
func (c *MockCommandBasic) Name() string      { return c.name }
func (c *MockCommandBasic) LongName() string  { return c.name }
func (c *MockCommandBasic) ShortName() string { return c.shortName }
func (c *MockCommandBasic) Desc() string      { return c.description }

func (c *MockCommandBasic) Config() *types.CmdConfig {
	return &types.CmdConfig{
		Version:     c.version,
		UseChinese:  c.useChinese,
		EnvPrefix:   c.envPrefix,
		UsageSyntax: c.usageSyntax,
		LogoText:    c.logoText,
	}
}

func (c *MockCommandBasic) SetVersion(version string)    { c.version = version }
func (c *MockCommandBasic) SetChinese(useChinese bool)   { c.useChinese = useChinese }
func (c *MockCommandBasic) SetEnvPrefix(prefix string)   { c.envPrefix = prefix }
func (c *MockCommandBasic) SetUsageSyntax(syntax string) { c.usageSyntax = syntax }
func (c *MockCommandBasic) SetLogoText(logo string)      { c.logoText = logo }

func (c *MockCommandBasic) AddFlag(f types.Flag) error {
	return c.flagRegistry.Register(f)
}

func (c *MockCommandBasic) AddFlags(flags ...types.Flag) error {
	for _, flag := range flags {
		if err := c.AddFlag(flag); err != nil {
			return err
		}
	}
	return nil
}

func (c *MockCommandBasic) AddFlagsFrom(flags []types.Flag) error {
	return c.AddFlags(flags...)
}

func (c *MockCommandBasic) GetFlag(name string) (types.Flag, bool) {
	return c.flagRegistry.Get(name)
}

func (c *MockCommandBasic) Flags() []types.Flag {
	return c.flagRegistry.List()
}

func (c *MockCommandBasic) FlagRegistry() types.FlagRegistry {
	return c.flagRegistry
}

func (c *MockCommandBasic) AddSubCmds(cmds ...types.Command) error {
	for _, cmd := range cmds {
		if err := c.cmdRegistry.Register(cmd); err != nil {
			return err
		}
	}
	return nil
}

func (c *MockCommandBasic) AddSubCmdFrom(cmds []types.Command) error {
	return c.AddSubCmds(cmds...)
}

func (c *MockCommandBasic) GetSubCmd(name string) (types.Command, bool) {
	return c.cmdRegistry.Get(name)
}

func (c *MockCommandBasic) SubCmds() []types.Command {
	return c.cmdRegistry.List()
}

func (c *MockCommandBasic) HasSubCmd(name string) bool {
	return c.cmdRegistry.Has(name)
}

func (c *MockCommandBasic) CmdRegistry() types.CmdRegistry {
	return c.cmdRegistry
}

func (c *MockCommandBasic) IsRootCmd() bool {
	return true // 默认为根命令
}

func (c *MockCommandBasic) Path() string {
	return c.name
}

func (c *MockCommandBasic) Parse(args []string) error {
	c.args = args
	c.parsed = true
	return nil
}

func (c *MockCommandBasic) ParseAndRoute(args []string) error {
	c.args = args
	c.parsed = true
	if c.runFunc != nil {
		return c.runFunc(c)
	}
	return nil
}

func (c *MockCommandBasic) ParseOnly(args []string) error {
	c.args = args
	c.parsed = true
	return nil
}

func (c *MockCommandBasic) IsParsed() bool {
	return c.parsed
}

func (c *MockCommandBasic) SetParsed(parsed bool) {
	c.parsed = parsed
}

func (c *MockCommandBasic) Args() []string {
	return c.args
}

func (c *MockCommandBasic) Arg(index int) string {
	if index < 0 || index >= len(c.args) {
		return ""
	}
	return c.args[index]
}

func (c *MockCommandBasic) NArg() int {
	return len(c.args)
}

func (c *MockCommandBasic) SetArgs(args []string) {
	c.args = args
}

func (c *MockCommandBasic) Run() error {
	if c.runFunc != nil {
		return c.runFunc(c)
	}
	return nil
}

func (c *MockCommandBasic) SetRun(fn func(types.Command) error) {
	c.runFunc = fn
}

func (c *MockCommandBasic) HasRunFunc() bool {
	return c.runFunc != nil
}

func (c *MockCommandBasic) Help() string {
	return "Mock command help"
}

func (c *MockCommandBasic) PrintHelp() {
	// 模拟打印帮助
}

func (c *MockCommandBasic) SetParser(p types.Parser) {
	// 模拟设置解析器
}

func (c *MockCommandBasic) SetDesc(desc string) {
	c.description = desc
}

func (c *MockCommandBasic) AddExample(title, cmd string) {
	// 模拟添加示例
}

func (c *MockCommandBasic) AddExamples(examples map[string]string) {
	// 模拟添加多个示例
}

func (c *MockCommandBasic) AddNote(note string) {
	// 模拟添加注意事项
}

func (c *MockCommandBasic) AddNotes(notes []string) {
	// 模拟添加多个注意事项
}

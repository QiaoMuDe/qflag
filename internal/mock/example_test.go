package mock

import (
	"gitee.com/MM-Q/qflag/internal/types"
)

// ExampleNewTestHelper 展示如何使用测试辅助工具
func ExampleNewTestHelper() {
	// 创建测试辅助工具
	helper := NewTestHelper()

	// 创建模拟命令
	cmd := helper.CreateMockCommandWithFlags(
		"test",
		"t",
		"Test command",
		helper.CreateMockBoolFlag("verbose", "v", "Verbose output", false),
		helper.CreateMockEnumFlag("mode", "m", "Operation mode", "normal", []string{"normal", "debug", "release"}),
	)

	// 添加示例和注意事项
	cmd.AddExample("Basic usage", "test -v -m debug")
	cmd.AddNotes([]string{
		"This is a mock command for testing",
		"Use -v for verbose output",
	})

	// 创建模拟解析器
	parser := NewMockParser()

	// 测试解析
	_ = parser.Parse(cmd, []string{"-v", "--mode", "debug"})

	// 创建模拟标志注册表
	registry := NewMockFlagRegistry()

	// 注册标志
	_ = registry.Register(helper.CreateMockBoolFlag("help", "h", "Show help", false))

	// 获取标志
	flag, exists := registry.Get("help")
	if exists {
		_ = flag.Name()
	}
}

// ExampleTestHelper_CreateMockCommandTree 展示如何创建命令树
func ExampleTestHelper_CreateMockCommandTree() {
	helper := NewTestHelper()

	// 创建命令树
	root := helper.CreateMockCommandTree()

	// 获取子命令
	sub1, _ := root.GetSubCmd("sub1")
	if sub1 != nil {
		// 获取孙子命令
		sub1_1, _ := sub1.GetSubCmd("sub1-1")
		if sub1_1 != nil {
			_ = sub1_1.Name()
		}
	}
}

// ExampleTestHelper_CreateMockBoolFlag 展示如何创建不同类型的模拟标志
func ExampleTestHelper_CreateMockBoolFlag() {
	helper := NewTestHelper()

	// 创建布尔标志
	boolFlag := helper.CreateMockBoolFlag("debug", "d", "Enable debug mode", false)

	// 创建枚举标志
	enumFlag := helper.CreateMockEnumFlag("level", "l", "Log level", "info", []string{"debug", "info", "warn", "error"})

	_ = boolFlag.Name()
	_ = enumFlag.EnumValues()
}

// ExampleTestHelper_CreateMockParserWithBehavior 展示如何使用模拟解析器
func ExampleTestHelper_CreateMockParserWithBehavior() {
	helper := NewTestHelper()

	// 创建带有错误的解析器
	errorParser := helper.CreateMockParserWithBehavior(
		types.NewError("PARSE_ERROR", "parse failed", nil),
		nil,
		true,
		false,
	)

	// 创建命令
	cmd := helper.CreateMockCommandWithRunFunc(
		"test",
		"t",
		"Test command",
		func(c types.Command) error {
			_ = c.Name()
			return nil
		},
	)

	// 测试解析
	_ = errorParser.Parse(cmd, []string{"--error"})
}

// ExampleNewMockFlagRegistry 展示如何使用模拟注册表
func ExampleNewMockFlagRegistry() {
	// 创建标志注册表
	flagReg := NewMockFlagRegistry()

	// 创建命令注册表
	cmdReg := NewMockCmdRegistry()

	// 创建标志和命令
	flag := NewMockBoolFlag("help", "h", "Show help", false)
	cmd := NewMockCommand("test", "t", "Test command")

	// 注册
	_ = flagReg.Register(flag)
	_ = cmdReg.Register(cmd)

	// 检查注册
	_ = flagReg.Has("help")
	_ = cmdReg.Has("test")

	// 获取列表
	_ = flagReg.List()
	_ = cmdReg.List()

	// 获取计数
	_ = flagReg.Count()
	_ = cmdReg.Count()
}

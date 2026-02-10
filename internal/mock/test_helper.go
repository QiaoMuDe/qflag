package mock

import (
	"gitee.com/MM-Q/qflag/internal/types"
)

// TestHelper 测试辅助工具
type TestHelper struct{}

// NewTestHelper 创建测试辅助工具
func NewTestHelper() *TestHelper {
	return &TestHelper{}
}

// CreateMockCommandWithFlags 创建带有标志的模拟命令
func (h *TestHelper) CreateMockCommandWithFlags(name, shortName, description string, flags ...types.Flag) *MockCommand {
	cmd := NewMockCommand(name, shortName, description)
	for _, flag := range flags {
		if err := cmd.AddFlag(flag); err != nil {
			panic(err) // 在测试辅助工具中, 如果添加标志失败, 应该立即 panic
		}
	}
	return cmd
}

// CreateMockSubCommandWithFlags 创建带有标志的模拟子命令
func (h *TestHelper) CreateMockSubCommandWithFlags(name, shortName, description string, parent types.Command, flags ...types.Flag) *MockCommand {
	cmd := NewMockSubCommand(name, shortName, description, parent)
	for _, flag := range flags {
		if err := cmd.AddFlag(flag); err != nil {
			panic(err) // 在测试辅助工具中, 如果添加标志失败, 应该立即 panic
		}
	}
	return cmd
}

// CreateMockCommandTree 创建模拟命令树
func (h *TestHelper) CreateMockCommandTree() *MockCommand {
	// 创建根命令
	root := NewMockCommand("root", "r", "Root command")

	// 创建子命令
	sub1 := NewMockSubCommand("sub1", "s1", "Sub command 1", root)
	sub2 := NewMockSubCommand("sub2", "s2", "Sub command 2", root)

	// 添加子命令到根命令
	if err := root.AddSubCmds(sub1, sub2); err != nil {
		panic(err) // 在测试辅助工具中, 如果添加子命令失败, 应该立即 panic
	}

	// 创建孙子命令
	sub1_1 := NewMockSubCommand("sub1-1", "s11", "Sub command 1-1", sub1)
	sub1_2 := NewMockSubCommand("sub1-2", "s12", "Sub command 1-2", sub1)

	// 添加孙子命令到子命令
	if err := sub1.AddSubCmds(sub1_1, sub1_2); err != nil {
		panic(err) // 在测试辅助工具中, 如果添加子命令失败, 应该立即 panic
	}

	return root
}

// CreateMockEnumFlag 创建模拟枚举标志
func (h *TestHelper) CreateMockEnumFlag(name, short, desc string, defaultValue string, allowedValues []string) *MockFlag {
	return NewMockEnumFlag(name, short, desc, defaultValue, allowedValues)
}

// CreateMockBoolFlag 创建模拟布尔标志
func (h *TestHelper) CreateMockBoolFlag(name, short, desc string, defaultValue bool) *MockFlag {
	return NewMockBoolFlag(name, short, desc, defaultValue)
}

// CreateMockCommandWithMutexGroups 创建带有互斥组的模拟命令
func (h *TestHelper) CreateMockCommandWithMutexGroups(name, shortName, description string, mutexGroups []types.MutexGroup) *MockCommand {
	cmd := NewMockCommand(name, shortName, description)

	for _, group := range mutexGroups {
		cmd.AddMutexGroup(group.Name, group.Flags, group.AllowNone)
	}

	return cmd
}

// CreateMockCommandWithExamples 创建带有示例的模拟命令
func (h *TestHelper) CreateMockCommandWithExamples(name, shortName, description string, examples map[string]string) *MockCommand {
	cmd := NewMockCommand(name, shortName, description)
	cmd.AddExamples(examples)
	return cmd
}

// CreateMockCommandWithNotes 创建带有注意事项的模拟命令
func (h *TestHelper) CreateMockCommandWithNotes(name, shortName, description string, notes []string) *MockCommand {
	cmd := NewMockCommand(name, shortName, description)
	cmd.AddNotes(notes)
	return cmd
}

// CreateMockCommandWithRunFunc 创建带有运行函数的模拟命令
func (h *TestHelper) CreateMockCommandWithRunFunc(name, shortName, description string, runFunc func(types.Command) error) *MockCommand {
	cmd := NewMockCommand(name, shortName, description)
	cmd.SetRun(runFunc)
	return cmd
}

// AssertCommandEqual 断言两个命令相等
func (h *TestHelper) AssertCommandEqual(expected, actual types.Command) bool {
	if expected.Name() != actual.Name() {
		return false
	}
	if expected.ShortName() != actual.ShortName() {
		return false
	}
	if expected.Desc() != actual.Desc() {
		return false
	}
	return true
}

// AssertFlagEqual 断言两个标志相等
func (h *TestHelper) AssertFlagEqual(expected, actual types.Flag) bool {
	if expected.Name() != actual.Name() {
		return false
	}
	if expected.ShortName() != actual.ShortName() {
		return false
	}
	if expected.Desc() != actual.Desc() {
		return false
	}
	if expected.Type() != actual.Type() {
		return false
	}
	return true
}

// CreateMockParserWithBehavior 创建带有特定行为的模拟解析器
func (h *TestHelper) CreateMockParserWithBehavior(parseError, routeError error, shouldCallParse, shouldCallParseAndRoute bool) *MockParser {
	parser := NewMockParserWithError(parseError, routeError)

	// 如果不需要调用特定方法, 可以在这里设置
	// 这取决于 MockParser 的具体实现

	return parser
}

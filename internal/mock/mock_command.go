package mock

import (
	"gitee.com/MM-Q/qflag/internal/types"
)

// MockCommand 扩展的模拟命令实现, 支持更多测试场景
type MockCommand struct {
	*MockCommandBasic
	isRoot      bool
	parent      types.Command
	mutexGroups []types.MutexGroup
	examples    map[string]string
	notes       []string
}

// NewMockCommand 创建扩展的模拟命令
func NewMockCommand(name, shortName, description string) *MockCommand {
	return &MockCommand{
		MockCommandBasic: NewMockCommandBasic(name, shortName, description),
		isRoot:           true,
		parent:           nil,
		mutexGroups:      []types.MutexGroup{},
		examples:         make(map[string]string),
		notes:            []string{},
	}
}

// NewMockSubCommand 创建子命令
func NewMockSubCommand(name, shortName, description string, parent types.Command) *MockCommand {
	return &MockCommand{
		MockCommandBasic: NewMockCommandBasic(name, shortName, description),
		isRoot:           false,
		parent:           parent,
		mutexGroups:      []types.MutexGroup{},
		examples:         make(map[string]string),
		notes:            []string{},
	}
}

// 重写 IsRootCmd 方法
func (c *MockCommand) IsRootCmd() bool {
	return c.isRoot
}

// 设置父命令
func (c *MockCommand) SetParent(parent types.Command) {
	c.parent = parent
}

// 获取父命令
func (c *MockCommand) GetParent() types.Command {
	return c.parent
}

// 重写 Path 方法, 包含父命令路径
func (c *MockCommand) Path() string {
	if c.parent == nil {
		return c.Name()
	}
	return c.parent.Path() + " " + c.Name()
}

// 添加互斥组
func (c *MockCommand) AddMutexGroup(name string, flags []string, allowNone bool) {
	group := types.MutexGroup{
		Name:      name,
		Flags:     flags,
		AllowNone: allowNone,
	}
	c.mutexGroups = append(c.mutexGroups, group)
}

// 获取互斥组
func (c *MockCommand) GetMutexGroups() []types.MutexGroup {
	return c.mutexGroups
}

// 重写 Config 方法, 包含互斥组
func (c *MockCommand) Config() *types.CmdConfig {
	config := c.MockCommandBasic.Config()
	config.MutexGroups = c.mutexGroups
	return config
}

// 重写 AddExample 方法
func (c *MockCommand) AddExample(title, cmd string) {
	c.examples[title] = cmd
}

// 重写 AddExamples 方法
func (c *MockCommand) AddExamples(examples map[string]string) {
	for k, v := range examples {
		c.examples[k] = v
	}
}

// 重写 AddNote 方法
func (c *MockCommand) AddNote(note string) {
	c.notes = append(c.notes, note)
}

// 重写 AddNotes 方法
func (c *MockCommand) AddNotes(notes []string) {
	c.notes = append(c.notes, notes...)
}

// 获取示例
func (c *MockCommand) GetExamples() map[string]string {
	return c.examples
}

// 获取注意事项
func (c *MockCommand) GetNotes() []string {
	return c.notes
}

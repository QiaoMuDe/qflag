package mock

import (
	"gitee.com/MM-Q/qflag/internal/types"
)

// MockCommand 扩展的模拟命令实现, 支持更多测试场景
type MockCommand struct {
	*MockCommandBasic
	isRoot         bool
	parent         types.Command
	mutexGroups    []types.MutexGroup
	requiredGroups []types.RequiredGroup
	examples       map[string]string
	notes          []string
}

// NewMockCommand 创建扩展的模拟命令
func NewMockCommand(name, shortName, description string) *MockCommand {
	return &MockCommand{
		MockCommandBasic: NewMockCommandBasic(name, shortName, description),
		isRoot:           true,
		parent:           nil,
		mutexGroups:      []types.MutexGroup{},
		requiredGroups:   []types.RequiredGroup{},
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
		requiredGroups:   []types.RequiredGroup{},
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
func (c *MockCommand) AddMutexGroup(name string, flags []string, allowNone bool) error {
	for _, group := range c.mutexGroups {
		if group.Name == name {
			return types.NewError("MUTEX_GROUP_ALREADY_EXISTS",
				"mutex group already exists", nil)
		}
	}
	group := types.MutexGroup{
		Name:      name,
		Flags:     flags,
		AllowNone: allowNone,
	}
	c.mutexGroups = append(c.mutexGroups, group)
	return nil
}

// 获取互斥组
func (c *MockCommand) GetMutexGroups() []types.MutexGroup {
	return c.mutexGroups
}

// 移除互斥组
func (c *MockCommand) RemoveMutexGroup(name string) error {
	for i, group := range c.mutexGroups {
		if group.Name == name {
			c.mutexGroups = append(c.mutexGroups[:i], c.mutexGroups[i+1:]...)
			return nil
		}
	}
	return types.NewError("MUTEX_GROUP_NOT_FOUND",
		"mutex group not found", nil)
}

// 获取互斥组
func (c *MockCommand) GetMutexGroup(name string) (*types.MutexGroup, bool) {
	for i := range c.mutexGroups {
		if c.mutexGroups[i].Name == name {
			return &c.mutexGroups[i], true
		}
	}
	return nil, false
}

// 添加必需组
func (c *MockCommand) AddRequiredGroup(name string, flags []string) error {
	for _, group := range c.requiredGroups {
		if group.Name == name {
			return types.NewError("REQUIRED_GROUP_ALREADY_EXISTS",
				"required group already exists", nil)
		}
	}
	c.requiredGroups = append(c.requiredGroups, types.RequiredGroup{
		Name:  name,
		Flags: flags,
	})
	return nil
}

// 移除必需组
func (c *MockCommand) RemoveRequiredGroup(name string) error {
	for i, group := range c.requiredGroups {
		if group.Name == name {
			c.requiredGroups = append(c.requiredGroups[:i], c.requiredGroups[i+1:]...)
			return nil
		}
	}
	return types.NewError("REQUIRED_GROUP_NOT_FOUND",
		"required group not found", nil)
}

// 获取必需组
func (c *MockCommand) GetRequiredGroup(name string) (*types.RequiredGroup, bool) {
	for i := range c.requiredGroups {
		if c.requiredGroups[i].Name == name {
			return &c.requiredGroups[i], true
		}
	}
	return nil, false
}

// 获取所有必需组
func (c *MockCommand) RequiredGroups() []types.RequiredGroup {
	if len(c.requiredGroups) == 0 {
		return []types.RequiredGroup{}
	}
	result := make([]types.RequiredGroup, len(c.requiredGroups))
	copy(result, c.requiredGroups)
	return result
}

// 重写 Config 方法, 包含互斥组和必需组
func (c *MockCommand) Config() *types.CmdConfig {
	config := c.MockCommandBasic.Config()

	var mutexGroups []types.MutexGroup
	if len(c.mutexGroups) > 0 {
		mutexGroups = make([]types.MutexGroup, len(c.mutexGroups))
		copy(mutexGroups, c.mutexGroups)
	}

	var requiredGroups []types.RequiredGroup
	if len(c.requiredGroups) > 0 {
		requiredGroups = make([]types.RequiredGroup, len(c.requiredGroups))
		copy(requiredGroups, c.requiredGroups)
	}

	return &types.CmdConfig{
		Version:        config.Version,
		UseChinese:     config.UseChinese,
		EnvPrefix:      config.EnvPrefix,
		UsageSyntax:    config.UsageSyntax,
		Example:        config.Example,
		Notes:          config.Notes,
		LogoText:       config.LogoText,
		MutexGroups:    mutexGroups,
		RequiredGroups: requiredGroups,
		Completion:     config.Completion,
	}
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

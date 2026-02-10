package mock

import (
	"gitee.com/MM-Q/qflag/internal/types"
)

// MockFlagRegistry 模拟标志注册表实现
type MockFlagRegistry struct {
	flags map[string]types.Flag
}

// NewMockFlagRegistry 创建新的模拟标志注册表
func NewMockFlagRegistry() *MockFlagRegistry {
	return &MockFlagRegistry{
		flags: make(map[string]types.Flag),
	}
}

// 实现 FlagRegistry 接口
func (r *MockFlagRegistry) Register(flag types.Flag) error {
	if flag == nil {
		return types.ErrInvalidFlagType
	}

	name := flag.Name()
	shortName := flag.ShortName()

	if name == "" && shortName == "" {
		return types.NewError("INVALID_FLAG_NAME", "flag name and short name cannot both be empty", nil)
	}

	// 检查长名称是否已存在
	if name != "" {
		if _, exists := r.flags[name]; exists {
			return types.ErrFlagAlreadyExists
		}
		r.flags[name] = flag
	}

	// 检查短名称是否已存在
	if shortName != "" {
		if _, exists := r.flags[shortName]; exists {
			return types.ErrFlagAlreadyExists
		}
		r.flags[shortName] = flag
	}

	return nil
}

func (r *MockFlagRegistry) Unregister(name string) error {
	if _, exists := r.flags[name]; !exists {
		return types.ErrFlagNotFound
	}
	delete(r.flags, name)
	return nil
}

func (r *MockFlagRegistry) Get(name string) (types.Flag, bool) {
	flag, exists := r.flags[name]
	return flag, exists
}

func (r *MockFlagRegistry) List() []types.Flag {
	result := make([]types.Flag, 0, len(r.flags))
	seen := make(map[string]bool)

	for _, flag := range r.flags {
		name := flag.Name()
		if !seen[name] {
			result = append(result, flag)
			seen[name] = true
		}
	}

	return result
}

func (r *MockFlagRegistry) Has(name string) bool {
	_, exists := r.flags[name]
	return exists
}

func (r *MockFlagRegistry) Count() int {
	return len(r.flags)
}

func (r *MockFlagRegistry) Clear() {
	r.flags = make(map[string]types.Flag)
}

// MockCmdRegistry 模拟命令注册表实现
type MockCmdRegistry struct {
	commands map[string]types.Command
}

// NewMockCmdRegistry 创建新的模拟命令注册表
func NewMockCmdRegistry() *MockCmdRegistry {
	return &MockCmdRegistry{
		commands: make(map[string]types.Command),
	}
}

// 实现 CmdRegistry 接口
func (r *MockCmdRegistry) Register(cmd types.Command) error {
	if cmd == nil {
		return types.ErrInvalidFlagType
	}

	name := cmd.Name()
	shortName := cmd.ShortName()

	if name == "" && shortName == "" {
		return types.NewError("INVALID_COMMAND_NAME", "command name and short name cannot both be empty", nil)
	}

	// 检查长名称是否已存在
	if name != "" {
		if _, exists := r.commands[name]; exists {
			return types.ErrCmdAlreadyExists
		}
		r.commands[name] = cmd
	}

	// 检查短名称是否已存在
	if shortName != "" {
		if _, exists := r.commands[shortName]; exists {
			return types.ErrCmdAlreadyExists
		}
		r.commands[shortName] = cmd
	}

	return nil
}

func (r *MockCmdRegistry) Unregister(name string) error {
	if _, exists := r.commands[name]; !exists {
		return types.ErrCmdNotFound
	}
	delete(r.commands, name)
	return nil
}

func (r *MockCmdRegistry) Get(name string) (types.Command, bool) {
	cmd, exists := r.commands[name]
	return cmd, exists
}

func (r *MockCmdRegistry) List() []types.Command {
	result := make([]types.Command, 0, len(r.commands))
	seen := make(map[string]bool)

	for _, cmd := range r.commands {
		name := cmd.Name()
		if !seen[name] {
			result = append(result, cmd)
			seen[name] = true
		}
	}

	return result
}

func (r *MockCmdRegistry) Has(name string) bool {
	_, exists := r.commands[name]
	return exists
}

func (r *MockCmdRegistry) Count() int {
	return len(r.commands)
}

func (r *MockCmdRegistry) Clear() {
	r.commands = make(map[string]types.Command)
}

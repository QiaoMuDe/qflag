package builtin

import (
	"testing"

	"gitee.com/MM-Q/qflag/internal/flag"
	"gitee.com/MM-Q/qflag/internal/mock"
	"gitee.com/MM-Q/qflag/internal/types"
)

// TestNewBuiltinFlagManager 测试内置标志管理器创建
func TestNewBuiltinFlagManager(t *testing.T) {
	manager := NewBuiltinFlagManager()

	if manager == nil {
		t.Error("NewBuiltinFlagManager() should not return nil")
	}

	// 检查默认处理器是否已注册
	expectedHandlers := []types.BuiltinFlagType{
		types.HelpFlag,
		types.VersionFlag,
		types.CompletionFlag,
	}

	for _, flagType := range expectedHandlers {
		if _, exists := manager.handlers[flagType]; !exists {
			t.Errorf("Expected handler for %v to be registered", flagType)
		}
	}
}

// TestBuiltinFlagManager_RegisterHandler 测试处理器注册
func TestBuiltinFlagManager_RegisterHandler(t *testing.T) {
	manager := NewBuiltinFlagManager()

	// 测试注册新处理器
	customHandler := mock.NewCustomHandler(types.BuiltinFlagType(999))
	manager.RegisterHandler(customHandler)

	// 检查处理器是否已注册
	if _, exists := manager.handlers[types.BuiltinFlagType(999)]; !exists {
		t.Error("Custom handler should be registered")
	}
}

// TestBuiltinFlagManager_RegisterBuiltinFlags 测试内置标志注册
func TestBuiltinFlagManager_RegisterBuiltinFlags(t *testing.T) {
	helper := mock.NewTestHelper()

	tests := []struct {
		name          string
		cmd           *mock.MockCommand
		expectHelp    bool
		expectVersion bool
		expectComp    bool
	}{
		{
			name: "根命令, 有版本信息",
			cmd: func() *mock.MockCommand {
				cmd := helper.CreateMockCommandWithFlags("root", "r", "Root command")
				cmd.SetVersion("1.0.0")
				return cmd
			}(),
			expectHelp:    true,
			expectVersion: true,
			expectComp:    false,
		},
		{
			name:          "根命令, 无版本信息",
			cmd:           helper.CreateMockCommandWithFlags("root", "r", "Root command"),
			expectHelp:    true,
			expectVersion: false,
			expectComp:    false,
		},
		{
			name: "子命令, 有版本信息",
			cmd: func() *mock.MockCommand {
				root := helper.CreateMockCommandWithFlags("root", "r", "Root command")
				cmd := helper.CreateMockSubCommandWithFlags("sub", "s", "Sub command", root)
				cmd.SetVersion("1.0.0")
				return cmd
			}(),
			expectHelp:    true,
			expectVersion: false,
			expectComp:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := NewBuiltinFlagManager()
			err := manager.RegisterBuiltinFlags(tt.cmd)

			if err != nil {
				t.Errorf("RegisterBuiltinFlags() error = %v", err)
				return
			}

			// 检查帮助标志
			if tt.expectHelp {
				if _, exists := tt.cmd.GetFlag(types.HelpFlagName); !exists {
					t.Error("Help flag should be registered")
				}
			} else {
				if _, exists := tt.cmd.GetFlag(types.HelpFlagName); exists {
					t.Error("Help flag should not be registered")
				}
			}

			// 检查版本标志
			if tt.expectVersion {
				if _, exists := tt.cmd.GetFlag(types.VersionFlagName); !exists {
					t.Error("Version flag should be registered")
				}
			} else {
				if _, exists := tt.cmd.GetFlag(types.VersionFlagName); exists {
					t.Error("Version flag should not be registered")
				}
			}

			// 检查补全标志
			if tt.expectComp {
				if _, exists := tt.cmd.GetFlag(types.CompletionFlagName); !exists {
					t.Error("Completion flag should be registered")
				}
			} else {
				if _, exists := tt.cmd.GetFlag(types.CompletionFlagName); exists {
					t.Error("Completion flag should not be registered")
				}
			}
		})
	}
}

// TestBuiltinFlagManager_HandleBuiltinFlags 测试内置标志处理
func TestBuiltinFlagManager_HandleBuiltinFlags(t *testing.T) {
	helper := mock.NewTestHelper()

	tests := []struct {
		name        string
		cmd         *mock.MockCommand
		expectPanic bool
	}{
		{
			name:        "没有设置内置标志",
			cmd:         helper.CreateMockCommandWithFlags("test", "t", "Test command"),
			expectPanic: false,
		},
		{
			name: "设置了帮助标志",
			cmd: func() *mock.MockCommand {
				cmd := helper.CreateMockCommandWithFlags("test", "t", "Test command")
				helpFlag := helper.CreateMockBoolFlag(types.HelpFlagName, types.HelpFlagShortName, "Help", false)
				if err := cmd.AddFlag(helpFlag); err != nil {
					panic(err) // 在测试中, 如果添加标志失败, 应该立即 panic
				}
				// 设置标志为已设置状态
				_ = helpFlag.Set("true")
				return cmd
			}(),
			expectPanic: true, // 因为Handle方法会调用os.Exit(0)
		},
		{
			name: "设置了版本标志",
			cmd: func() *mock.MockCommand {
				cmd := helper.CreateMockCommandWithFlags("test", "t", "Test command")
				cmd.SetVersion("1.0.0")
				versionFlag := helper.CreateMockBoolFlag(types.VersionFlagName, types.VersionFlagShortName, "Version", false)
				if err := cmd.AddFlag(versionFlag); err != nil {
					panic(err) // 在测试中, 如果添加标志失败, 应该立即 panic
				}
				// 设置标志为已设置状态
				_ = versionFlag.Set("true")
				return cmd
			}(),
			expectPanic: true, // 因为Handle方法会调用os.Exit(0)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					if !tt.expectPanic {
						t.Errorf("HandleBuiltinFlags() unexpected panic: %v", r)
					}
				} else {
					if tt.expectPanic {
						t.Error("HandleBuiltinFlags() expected panic but did not panic")
					}
				}
			}()

			manager := NewBuiltinFlagManager()
			err := manager.HandleBuiltinFlags(tt.cmd)

			if err != nil && !tt.expectPanic {
				t.Errorf("HandleBuiltinFlags() error = %v", err)
			}
		})
	}
}

// TestBuiltinFlagManager_isBuiltinFlag 测试内置标志检查
func TestBuiltinFlagManager_isBuiltinFlag(t *testing.T) {
	manager := NewBuiltinFlagManager()

	tests := []struct {
		name            string
		flag            types.Flag
		expectType      types.BuiltinFlagType
		expectIsBuiltin bool
	}{
		{
			name:            "帮助标志长名称",
			flag:            flag.NewBoolFlag(types.HelpFlagName, "", "Help", false),
			expectType:      types.HelpFlag,
			expectIsBuiltin: true,
		},
		{
			name:            "帮助标志短名称",
			flag:            flag.NewBoolFlag("", types.HelpFlagShortName, "Help", false),
			expectType:      types.HelpFlag,
			expectIsBuiltin: true,
		},
		{
			name:            "版本标志长名称",
			flag:            flag.NewBoolFlag(types.VersionFlagName, "", "Version", false),
			expectType:      types.VersionFlag,
			expectIsBuiltin: true,
		},
		{
			name:            "版本标志短名称",
			flag:            flag.NewBoolFlag("", types.VersionFlagShortName, "Version", false),
			expectType:      types.VersionFlag,
			expectIsBuiltin: true,
		},
		{
			name:            "补全标志长名称",
			flag:            flag.NewEnumFlag(types.CompletionFlagName, "", "Completion", "bash", []string{"bash", "pwsh"}),
			expectType:      types.CompletionFlag,
			expectIsBuiltin: true,
		},
		{
			name:            "自定义标志",
			flag:            flag.NewBoolFlag("custom", "c", "Custom flag", false),
			expectType:      0,
			expectIsBuiltin: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flagType, isBuiltin := manager.isBuiltinFlag(tt.flag)

			if isBuiltin != tt.expectIsBuiltin {
				t.Errorf("isBuiltinFlag() isBuiltin = %v, expect %v", isBuiltin, tt.expectIsBuiltin)
			}

			if flagType != tt.expectType {
				t.Errorf("isBuiltinFlag() flagType = %v, expect %v", flagType, tt.expectType)
			}
		})
	}
}

// TestHelpHandler 测试帮助处理器
func TestHelpHandler(t *testing.T) {
	helper := mock.NewTestHelper()
	handler := &HelpHandler{}

	// 测试Type方法
	if handler.Type() != types.HelpFlag {
		t.Errorf("Expected HelpFlag, got %v", handler.Type())
	}

	// 测试ShouldRegister方法
	tests := []struct {
		name     string
		cmd      *mock.MockCommand
		expected bool
	}{
		{
			name:     "根命令",
			cmd:      helper.CreateMockCommandWithFlags("root", "r", "Root command"),
			expected: true,
		},
		{
			name: "子命令",
			cmd: func() *mock.MockCommand {
				root := helper.CreateMockCommandWithFlags("root", "r", "Root command")
				return helper.CreateMockSubCommandWithFlags("sub", "s", "Sub command", root)
			}(),
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if handler.ShouldRegister(tt.cmd) != tt.expected {
				t.Errorf("HelpHandler.ShouldRegister() = %v, expected %v", handler.ShouldRegister(tt.cmd), tt.expected)
			}
		})
	}
}

// TestVersionHandler 测试版本处理器
func TestVersionHandler(t *testing.T) {
	helper := mock.NewTestHelper()
	handler := &VersionHandler{}

	// 测试Type方法
	if handler.Type() != types.VersionFlag {
		t.Errorf("Expected VersionFlag, got %v", handler.Type())
	}

	// 测试ShouldRegister方法
	tests := []struct {
		name     string
		cmd      *mock.MockCommand
		expected bool
	}{
		{
			name: "根命令, 有版本信息",
			cmd: func() *mock.MockCommand {
				cmd := helper.CreateMockCommandWithFlags("root", "r", "Root command")
				cmd.SetVersion("1.0.0")
				return cmd
			}(),
			expected: true,
		},
		{
			name:     "根命令, 无版本信息",
			cmd:      helper.CreateMockCommandWithFlags("root", "r", "Root command"),
			expected: false,
		},
		{
			name: "子命令, 有版本信息",
			cmd: func() *mock.MockCommand {
				root := helper.CreateMockCommandWithFlags("root", "r", "Root command")
				cmd := helper.CreateMockSubCommandWithFlags("sub", "s", "Sub command", root)
				cmd.SetVersion("1.0.0")
				return cmd
			}(),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if handler.ShouldRegister(tt.cmd) != tt.expected {
				t.Errorf("VersionHandler.ShouldRegister() = %v, expected %v", handler.ShouldRegister(tt.cmd), tt.expected)
			}
		})
	}
}

// TestCompletionHandler 测试补全处理器
func TestCompletionHandler(t *testing.T) {
	helper := mock.NewTestHelper()
	handler := &CompletionHandler{}

	// 测试Type方法
	if handler.Type() != types.CompletionFlag {
		t.Errorf("Expected CompletionFlag, got %v", handler.Type())
	}

	// 测试ShouldRegister方法
	tests := []struct {
		name     string
		cmd      *mock.MockCommand
		expected bool
	}{
		{
			name:     "根命令",
			cmd:      helper.CreateMockCommandWithFlags("root", "r", "Root command"),
			expected: false,
		},
		{
			name: "子命令",
			cmd: func() *mock.MockCommand {
				root := helper.CreateMockCommandWithFlags("root", "r", "Root command")
				return helper.CreateMockSubCommandWithFlags("sub", "s", "Sub command", root)
			}(),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if handler.ShouldRegister(tt.cmd) != tt.expected {
				t.Errorf("CompletionHandler.ShouldRegister() = %v, expected %v", handler.ShouldRegister(tt.cmd), tt.expected)
			}
		})
	}
}

// TestGetShellTypeFromArgs 测试从参数获取Shell类型
func TestGetShellTypeFromArgs(t *testing.T) {
	helper := mock.NewTestHelper()

	tests := []struct {
		name     string
		cmd      *mock.MockCommand
		expected string
	}{
		{
			name: "有补全标志参数",
			cmd: func() *mock.MockCommand {
				cmd := helper.CreateMockCommandWithFlags("test", "t", "Test command")
				completionFlag := flag.NewEnumFlag(types.CompletionFlagName, "", "Completion", "pwsh", []string{"bash", "pwsh"})
				if err := cmd.AddFlag(completionFlag); err != nil {
					panic(err) // 在测试中, 如果添加标志失败, 应该立即 panic
				}
				return cmd
			}(),
			expected: "pwsh",
		},
		{
			name:     "无补全标志参数",
			cmd:      helper.CreateMockCommandWithFlags("test", "t", "Test command"),
			expected: types.CurrentShell(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getShellTypeFromArgs(tt.cmd)
			if result != tt.expected {
				t.Errorf("getShellTypeFromArgs() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

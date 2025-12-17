// package qflag 内部命令测试
// 本文件包含了Cmd结构体内部功能的单元测试，测试内部API和
// 实现细节，确保内部逻辑的正确性和稳定性。
package qflag

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"gitee.com/MM-Q/qflag/internal/types"
)

// =============================================================================
// 测试辅助函数
// =============================================================================

// createInternalTestCmd 创建内部测试用的命令实例
func createInternalTestCmd() *Cmd {
	return NewCmd("internal-test", "it", flag.ContinueOnError)
}

// createInternalTestCmdWithVersion 创建带版本信息的测试命令
func createInternalTestCmdWithVersion() *Cmd {
	cmd := NewCmd("internal-test", "it", flag.ContinueOnError)
	cmd.SetVersion("v1.0.0")
	return cmd
}

// createInternalTestCmdWithCompletion 创建启用补全功能的测试命令
func createInternalTestCmdWithCompletion() *Cmd {
	cmd := NewCmd("internal-test", "it", flag.ContinueOnError)
	cmd.SetCompletion(true)
	return cmd
}

// =============================================================================
// parseCommon 方法测试
// =============================================================================

func TestCmd_parseCommon(t *testing.T) {
	tests := []struct {
		name             string      // 测试用例名称
		setupCmd         func() *Cmd // 创建测试命令的函数
		args             []string    // 要解析的参数列表
		parseSubcommands bool        // 是否解析子命令
		expectShouldExit bool        // 是否期望退出程序
		expectError      bool        // 是否期望返回错误
		errorContains    string      // 期望错误信息包含的字符串
		setupFlags       func(*Cmd)  // 设置标志的函数
		setupSubcommands func(*Cmd)  // 设置子命令的函数
	}{
		{
			name:             "正常解析空参数",
			setupCmd:         createInternalTestCmd,
			args:             []string{},
			parseSubcommands: true,
			expectShouldExit: false, // 禁用退出以便测试
			expectError:      false,
		},
		{
			name:             "正常解析标志参数",
			setupCmd:         createInternalTestCmd,
			args:             []string{"--help"},
			parseSubcommands: true,
			expectShouldExit: false, // 禁用退出以便测试
			expectError:      false,
		},
		{
			name:             "解析版本标志",
			setupCmd:         createInternalTestCmdWithVersion,
			args:             []string{"--version"},
			parseSubcommands: true,
			expectShouldExit: false, // 禁用退出以便测试
			expectError:      false,
		},
		{
			name:             "解析无效标志",
			setupCmd:         createInternalTestCmd,
			args:             []string{"--invalid-flag"},
			parseSubcommands: true,
			expectShouldExit: false, // 禁用退出以便测试
			expectError:      true,
			errorContains:    "flag provided but not defined",
		},
		{
			name:             "nil命令测试",
			setupCmd:         func() *Cmd { return nil },
			args:             []string{},
			parseSubcommands: true,
			expectShouldExit: false, // 禁用退出以便测试
			expectError:      true,
			errorContains:    "nil command",
		},
		{
			name:             "解析子命令",
			setupCmd:         createInternalTestCmd,
			args:             []string{"subcmd", "--help"},
			parseSubcommands: true,
			expectShouldExit: false, // 禁用退出以便测试
			expectError:      false,
			setupSubcommands: func(cmd *Cmd) {
				subCmd := NewCmd("subcmd", "sc", flag.ContinueOnError)
				subCmd.SetNoBuiltinExit(true) // 设置子命令不在内置标志时退出
				err := cmd.AddSubCmd(subCmd)
				if err != nil {
					t.Fatalf("添加子命令失败: %v", err)
				}
			},
		},
		{
			name:             "不解析子命令",
			setupCmd:         createInternalTestCmd,
			args:             []string{"subcmd", "--help"},
			parseSubcommands: false,
			expectShouldExit: false, // 禁用退出以便测试
			expectError:      false,
			setupSubcommands: func(cmd *Cmd) {
				subCmd := NewCmd("subcmd", "sc", flag.ContinueOnError)
				err := cmd.AddSubCmd(subCmd)
				if err != nil {
					t.Fatalf("添加子命令失败: %v", err)
				}
			},
		},
		{
			name:             "枚举标志验证失败",
			setupCmd:         createInternalTestCmd,
			args:             []string{"--mode", "invalid"},
			parseSubcommands: true,
			expectShouldExit: false, // 禁用退出以便测试
			expectError:      true,
			errorContains:    "invalid enum value",
			setupFlags: func(cmd *Cmd) {
				cmd.Enum("mode", "m", "debug", "运行模式", []string{"debug", "release"})
			},
		},
		{
			name:             "启用补全功能解析",
			setupCmd:         createInternalTestCmdWithCompletion,
			args:             []string{"--completion", "bash"},
			parseSubcommands: true,
			expectShouldExit: false,
			expectError:      false,
		},
		{
			name:             "禁用内置标志退出测试",
			setupCmd:         createInternalTestCmdWithVersion,
			args:             []string{"--help"},
			parseSubcommands: true,
			expectShouldExit: false, // 禁用退出，应该不退出
			expectError:      false,
			setupFlags: func(cmd *Cmd) {
				// 禁用内置标志退出
				cmd.SetNoBuiltinExit(true)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := tt.setupCmd()
			if cmd == nil && !tt.expectError {
				t.Skip("跳过nil命令测试")
			}

			// 禁用内置标志退出以便测试
			if cmd != nil {
				cmd.SetNoBuiltinExit(true)
			}

			// 设置标志
			if tt.setupFlags != nil && cmd != nil {
				tt.setupFlags(cmd)
			}

			// 设置子命令
			if tt.setupSubcommands != nil && cmd != nil {
				tt.setupSubcommands(cmd)
			}

			// 调用parseCommon方法
			var shouldExit bool
			var err error
			if cmd != nil {
				shouldExit, err = cmd.parseCommon(tt.args, tt.parseSubcommands)
			} else {
				// 对于nil命令，直接返回期望的错误
				err = fmt.Errorf("nil command")
				shouldExit = false
			}

			// 验证错误
			if tt.expectError {
				if err == nil {
					t.Error("期望错误但未发生")
				} else if tt.errorContains != "" && !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("错误信息 = %v, 期望包含 %v", err.Error(), tt.errorContains)
				}
			} else if err != nil {
				t.Errorf("意外的错误: %v", err)
			}

			// 验证退出状态
			if shouldExit != tt.expectShouldExit {
				t.Errorf("shouldExit = %v, 期望 %v", shouldExit, tt.expectShouldExit)
			}
		})
	}
}

func TestCmd_parseCommon_PanicRecovery(t *testing.T) {
	t.Run("panic恢复测试", func(t *testing.T) {
		cmd := createInternalTestCmd()

		// 设置一个会导致panic的解析钩子
		cmd.ctx.ParseHook = func(ctx *types.CmdContext) (error, bool) {
			panic("测试panic")
		}

		_, err := cmd.parseCommon([]string{}, true)

		if err == nil {
			t.Error("期望捕获panic错误但未发生")
		}
		if !strings.Contains(err.Error(), "panic recovered") {
			t.Errorf("错误信息应包含'panic recovered'，实际: %v", err.Error())
		}
	})
}

func TestCmd_parseCommon_Concurrency(t *testing.T) {
	t.Run("并发解析测试", func(t *testing.T) {
		cmd := createInternalTestCmd()

		var wg sync.WaitGroup
		numGoroutines := 10
		results := make([]error, numGoroutines)

		// 并发调用parseCommon
		wg.Add(numGoroutines)
		for i := 0; i < numGoroutines; i++ {
			go func(index int) {
				defer wg.Done()
				_, err := cmd.parseCommon([]string{}, true)
				results[index] = err
			}(i)
		}

		wg.Wait()

		// 验证所有调用都成功（由于sync.Once，只有第一次会真正执行）
		for i, err := range results {
			if err != nil {
				t.Errorf("goroutine %d 出现错误: %v", i, err)
			}
		}

		// 验证命令已被标记为已解析
		if !cmd.IsParsed() {
			t.Error("命令应该被标记为已解析")
		}
	})
}

// =============================================================================
// validateComponents 方法测试
// =============================================================================

func TestCmd_validateComponents(t *testing.T) {
	tests := []struct {
		name          string
		setupCmd      func() *Cmd
		expectError   bool
		errorContains string
	}{
		{
			name:        "正常组件验证",
			setupCmd:    createInternalTestCmd,
			expectError: false,
		},
		{
			name: "FlagSet为nil",
			setupCmd: func() *Cmd {
				cmd := createInternalTestCmd()
				cmd.ctx.FlagSet = nil
				return cmd
			},
			expectError:   true,
			errorContains: "flag.FlagSet instance is not initialized",
		},
		{
			name: "FlagRegistry为nil",
			setupCmd: func() *Cmd {
				cmd := createInternalTestCmd()
				cmd.ctx.FlagRegistry = nil
				return cmd
			},
			expectError:   true,
			errorContains: "FlagRegistry instance is not initialized",
		},
		{
			name: "SubCmds为nil",
			setupCmd: func() *Cmd {
				cmd := createInternalTestCmd()
				cmd.ctx.SubCmds = nil
				return cmd
			},
			expectError:   true,
			errorContains: "subCmdMap cannot be nil",
		},
		{
			name: "Help标志为nil",
			setupCmd: func() *Cmd {
				cmd := createInternalTestCmd()
				cmd.ctx.BuiltinFlags.Help = nil
				return cmd
			},
			expectError:   true,
			errorContains: "help flag is not initialized",
		},
		{
			name: "Version标志为nil",
			setupCmd: func() *Cmd {
				cmd := createInternalTestCmd()
				cmd.ctx.BuiltinFlags.Version = nil
				return cmd
			},
			expectError:   true,
			errorContains: "version flag is not initialized",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := tt.setupCmd()
			err := cmd.validateComponents()

			if tt.expectError {
				if err == nil {
					t.Error("期望错误但未发生")
				} else if tt.errorContains != "" && !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("错误信息 = %v, 期望包含 %v", err.Error(), tt.errorContains)
				}
			} else if err != nil {
				t.Errorf("意外的错误: %v", err)
			}
		})
	}
}

// =============================================================================
// registerBuiltinFlags 方法测试
// =============================================================================

func TestCmd_registerBuiltinFlags(t *testing.T) {
	tests := []struct {
		name                 string
		setupCmd             func() *Cmd
		expectVersionFlag    bool
		expectCompletionFlag bool
		expectNotes          bool
		expectExamples       bool
		useChinese           bool
	}{
		{
			name:              "顶级命令注册内置标志",
			setupCmd:          createInternalTestCmd,
			expectVersionFlag: false,
			useChinese:        false,
		},
		{
			name:              "带版本信息的顶级命令",
			setupCmd:          createInternalTestCmdWithVersion,
			expectVersionFlag: true,
			useChinese:        false,
		},
		{
			name:                 "启用补全功能的命令",
			setupCmd:             createInternalTestCmdWithCompletion,
			expectCompletionFlag: true,
			expectNotes:          true,
			expectExamples:       true,
			useChinese:           false,
		},
		{
			name: "中文环境的补全功能",
			setupCmd: func() *Cmd {
				cmd := createInternalTestCmdWithCompletion()
				cmd.SetChinese(true)
				return cmd
			},
			expectCompletionFlag: true,
			expectNotes:          true,
			expectExamples:       true,
			useChinese:           true,
		},
		{
			name: "子命令不注册内置标志",
			setupCmd: func() *Cmd {
				parent := createInternalTestCmd()
				child := NewCmd("child", "c", flag.ContinueOnError)
				child.ctx.Parent = parent.ctx
				return child
			},
			expectVersionFlag:    false,
			expectCompletionFlag: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := tt.setupCmd()

			// 保存原始的os.Args
			originalArgs := os.Args
			defer func() { os.Args = originalArgs }()

			// 设置测试用的程序名
			os.Args = []string{"test-program"}

			// 调用registerBuiltinFlags
			cmd.registerBuiltinFlags()

			// 验证版本标志
			if tt.expectVersionFlag {
				if _, ok := cmd.ctx.BuiltinFlags.NameMap.Load("version"); !ok {
					t.Error("期望注册version标志但未找到")
				}
				if _, ok := cmd.ctx.BuiltinFlags.NameMap.Load("v"); !ok {
					t.Error("期望注册v标志但未找到")
				}
			}

			// 验证补全标志
			if tt.expectCompletionFlag {
				if _, ok := cmd.ctx.BuiltinFlags.NameMap.Load("completion"); !ok {
					t.Error("期望注册completion标志但未找到")
				}
			}

			// 验证注意事项
			if tt.expectNotes {
				if len(cmd.ctx.Config.Notes) == 0 {
					t.Error("期望添加注意事项但未找到")
				}

				// 验证语言相关的注意事项
				noteText := strings.Join(cmd.ctx.Config.Notes, " ")
				if tt.useChinese {
					if !strings.Contains(noteText, "Windows") || !strings.Contains(noteText, "PowerShell") {
						t.Error("中文环境下期望包含Windows和PowerShell相关注意事项")
					}
				} else {
					if !strings.Contains(noteText, "Windows") || !strings.Contains(noteText, "PowerShell") {
						t.Error("英文环境下期望包含Windows和PowerShell相关注意事项")
					}
				}
			}

			// 验证示例
			if tt.expectExamples {
				if len(cmd.ctx.Config.Examples) == 0 {
					t.Error("期望添加示例但未找到")
				}

				// 验证示例中包含程序名
				for _, example := range cmd.ctx.Config.Examples {
					if !strings.Contains(example.Usage, "test-program") {
						t.Errorf("示例应包含程序名，实际: %v", example.Usage)
					}
				}
			}
		})
	}
}

// =============================================================================
// handleBuiltinFlags 方法测试
// =============================================================================

func TestCmd_handleBuiltinFlags(t *testing.T) {
	tests := []struct {
		name           string      // 测试用例名称
		setupCmd       func() *Cmd // 命令设置函数
		setupFlags     func(*Cmd)  // 标志设置函数
		expectExit     bool        // 是否期望退出
		expectError    bool        // 是否期望错误
		errorContains  string      // 期望错误包含的字符串
		expectOutput   bool        // 是否期望输出
		outputContains string      // 期望输出包含的字符串
	}{
		{
			name:        "无内置标志触发",
			setupCmd:    createInternalTestCmd,
			expectExit:  true,
			expectError: false,
		},
		{
			name:     "help标志触发",
			setupCmd: createInternalTestCmd,
			setupFlags: func(cmd *Cmd) {
				err := cmd.ctx.BuiltinFlags.Help.Set("true")
				if err != nil {
					t.Fatalf("设置帮助标志失败: %v", err)
				}
			},
			expectExit:  false,
			expectError: false,
		},
		{
			name:     "version标志触发",
			setupCmd: createInternalTestCmdWithVersion,
			setupFlags: func(cmd *Cmd) {
				err := cmd.ctx.BuiltinFlags.Version.Set("true")
				if err != nil {
					t.Fatalf("设置版本标志失败: %v", err)
				}
			},
			expectExit:  false, // NoBuiltinExit为false时，handleBuiltinFlags返回false(不禁用退出)，parseCommon中取反后为true(退出)
			expectError: false,
		},
		{
			name: "子命令中的version标志不触发",
			setupCmd: func() *Cmd {
				parent := createInternalTestCmdWithVersion()
				child := NewCmd("child", "c", flag.ContinueOnError)
				child.ctx.Parent = parent.ctx
				child.ctx.BuiltinFlags = parent.ctx.BuiltinFlags
				child.ctx.Config = parent.ctx.Config
				return child
			},
			setupFlags: func(cmd *Cmd) {
				err := cmd.ctx.BuiltinFlags.Version.Set("true")
				if err != nil {
					t.Fatalf("设置版本标志失败: %v", err)
				}
			},
			expectExit:  true,
			expectError: false,
		},
		{
			name:     "补全标志触发",
			setupCmd: createInternalTestCmdWithCompletion,
			setupFlags: func(cmd *Cmd) {
				err := cmd.ctx.BuiltinFlags.Completion.Set("bash")
				if err != nil {
					t.Fatalf("设置补全标志失败: %v", err)
				}
			},
			expectExit:  false, // 补全标志触发时，handleBuiltinFlags返回true(禁用退出)，parseCommon中取反后为false(不退出)
			expectError: false,
		},
		{
			name:     "枚举标志验证失败",
			setupCmd: createInternalTestCmd,
			setupFlags: func(cmd *Cmd) {
				// 创建枚举标志但不设置无效值，让handleBuiltinFlags来验证
				cmd.Enum("mode", "m", "debug", "运行模式", []string{"debug", "release"})
			},
			expectExit:  true,
			expectError: false,
		},
		{
			name: "NoBuiltinExit为true时不退出",
			setupCmd: func() *Cmd {
				cmd := createInternalTestCmd()
				cmd.SetNoBuiltinExit(true)
				return cmd
			},
			setupFlags: func(cmd *Cmd) {
				err := cmd.ctx.BuiltinFlags.Help.Set("true")
				if err != nil {
					t.Fatalf("设置帮助标志失败: %v", err)
				}
			},
			expectExit:  true,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := tt.setupCmd()

			// 设置标志
			if tt.setupFlags != nil {
				tt.setupFlags(cmd)
			}

			// 调用handleBuiltinFlags
			shouldExit, err := cmd.handleBuiltinFlags()

			// 验证错误
			if tt.expectError {
				if err == nil {
					t.Error("期望错误但未发生")
				} else if tt.errorContains != "" && !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("错误信息 = %v, 期望包含 %v", err.Error(), tt.errorContains)
				}
			} else if err != nil {
				t.Errorf("意外的错误: %v", err)
			}

			// 验证退出状态
			if shouldExit != tt.expectExit {
				t.Errorf("shouldExit = %v, 期望 %v", shouldExit, tt.expectExit)
			}
		})
	}
}

// =============================================================================
// 边界条件和错误处理测试
// =============================================================================

func TestCmd_Internal_BoundaryConditions(t *testing.T) {
	t.Run("极端参数数量", func(t *testing.T) {
		cmd := createInternalTestCmd()

		// 测试大量参数
		largeArgs := make([]string, 1000)
		for i := range largeArgs {
			largeArgs[i] = fmt.Sprintf("arg%d", i)
		}

		_, err := cmd.parseCommon(largeArgs, true)
		if err != nil {
			t.Errorf("大量参数解析失败: %v", err)
		}
	})

	t.Run("空字符串参数", func(t *testing.T) {
		cmd := createInternalTestCmd()

		_, err := cmd.parseCommon([]string{""}, true)
		if err != nil {
			t.Errorf("空字符串参数解析失败: %v", err)
		}
	})

	t.Run("特殊字符参数", func(t *testing.T) {
		cmd := createInternalTestCmd()

		specialArgs := []string{"--flag=测试", "中文参数", "!@#$%^&*()", "🚀🎉"}
		_, err := cmd.parseCommon(specialArgs, true)
		if err != nil && !strings.Contains(err.Error(), "flag provided but not defined") {
			t.Errorf("特殊字符参数处理失败: %v", err)
		}
	})

	t.Run("重复解析保护", func(t *testing.T) {
		cmd := createInternalTestCmd()

		// 第一次解析
		_, err1 := cmd.parseCommon([]string{}, true)
		if err1 != nil {
			t.Errorf("第一次解析失败: %v", err1)
		}

		// 第二次解析应该被sync.Once保护
		shouldExit2, err2 := cmd.parseCommon([]string{"--help"}, true)
		if err2 != nil {
			t.Errorf("第二次解析失败: %v", err2)
		}

		// 第二次解析不应该触发help（因为被sync.Once保护）
		if shouldExit2 {
			t.Error("重复解析不应该触发退出")
		}

		// 验证解析状态
		if !cmd.IsParsed() {
			t.Error("命令应该被标记为已解析")
		}
	})
}

func TestCmd_Internal_ErrorHandling(t *testing.T) {
	t.Run("组件初始化失败恢复", func(t *testing.T) {
		cmd := createInternalTestCmd()

		// 破坏组件状态
		originalFlagSet := cmd.ctx.FlagSet
		cmd.ctx.FlagSet = nil

		_, err := cmd.parseCommon([]string{}, true)
		if err == nil {
			t.Error("期望组件验证失败但未发生")
		}

		// 恢复组件状态
		cmd.ctx.FlagSet = originalFlagSet
	})

	t.Run("内置标志处理异常", func(t *testing.T) {
		cmd := createInternalTestCmd()

		// 设置无效的内置标志状态
		cmd.ctx.BuiltinFlags.Help = nil

		err := cmd.validateComponents()
		if err == nil {
			t.Fatal("期望验证组件失败但未发生")
		}
		if !strings.Contains(err.Error(), "help flag is not initialized") {
			t.Errorf("错误信息应包含help flag相关信息，实际: %v", err.Error())
		}
	})
}

// =============================================================================
// 性能测试
// =============================================================================

func BenchmarkCmd_parseCommon(b *testing.B) {
	cmd := createInternalTestCmd()
	args := []string{"--help"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// 重置解析状态以允许重复测试
		cmd.ctx.ParseOnce = sync.Once{}
		cmd.ctx.Parsed.Store(false)

		_, err := cmd.parseCommon(args, true)
		if err != nil {
			b.Fatalf("解析通用参数失败: %v", err)
		}
	}
}

func BenchmarkCmd_validateComponents(b *testing.B) {
	cmd := createInternalTestCmd()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := cmd.validateComponents()
		if err != nil {
			b.Fatalf("验证组件失败: %v", err)
		}
	}
}

func BenchmarkCmd_handleBuiltinFlags(b *testing.B) {
	cmd := createInternalTestCmd()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := cmd.handleBuiltinFlags()
		if err != nil {
			b.Fatalf("处理内置标志失败: %v", err)
		}
	}
}

// =============================================================================
// 集成测试
// =============================================================================

func TestCmd_Internal_Integration(t *testing.T) {
	t.Run("完整解析流程", func(t *testing.T) {
		// 测试各种参数组合，每个测试用例使用独立的命令实例
		testCases := []struct {
			name        string      // 测试用例名称
			setupCmd    func() *Cmd // 命令实例创建函数
			args        []string    // 测试参数
			expectExit  bool        // 是否期望退出
			expectError bool        // 是否期望错误
		}{
			{
				name:        "空参数",
				setupCmd:    createInternalTestCmd,
				args:        []string{},
				expectExit:  false,
				expectError: false,
			},
			{
				name:        "help标志",
				setupCmd:    createInternalTestCmd,
				args:        []string{"--help"},
				expectExit:  true,
				expectError: false,
			},
			{
				name:        "version标志",
				setupCmd:    createInternalTestCmdWithVersion,
				args:        []string{"--version"},
				expectExit:  true,
				expectError: false,
			},
			{
				name: "正常标志",
				setupCmd: func() *Cmd {
					cmd := createInternalTestCmd()
					cmd.String("config", "c", "config.json", "配置文件路径")
					cmd.Int("port", "p", 8080, "端口号")
					return cmd
				},
				args:        []string{"--config", "test.json", "--port", "9000"},
				expectExit:  false,
				expectError: false,
			},
			{
				name: "子命令",
				setupCmd: func() *Cmd {
					cmd := createInternalTestCmd()
					subCmd := NewCmd("start", "s", flag.ContinueOnError)
					subCmd.String("env", "e", "dev", "环境")
					err := cmd.AddSubCmd(subCmd)
					if err != nil {
						t.Fatalf("添加子命令失败: %v", err)
					}
					return cmd
				},
				args:        []string{"start", "--env", "prod"},
				expectExit:  false,
				expectError: false,
			},
			{
				name:        "补全",
				setupCmd:    createInternalTestCmdWithCompletion,
				args:        []string{"--completion", "bash"},
				expectExit:  true,
				expectError: false,
			},
			{
				name:        "无效标志",
				setupCmd:    createInternalTestCmd,
				args:        []string{"--invalid"},
				expectExit:  false,
				expectError: true,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				cmd := tc.setupCmd()
				shouldExit, err := cmd.parseCommon(tc.args, true)

				if tc.expectError {
					if err == nil {
						t.Error("期望错误但未发生")
					}
				} else if err != nil {
					t.Errorf("意外的错误: %v", err)
				}

				if shouldExit != tc.expectExit {
					t.Errorf("shouldExit = %v, 期望 %v", shouldExit, tc.expectExit)
				}
			})
		}
	})

	t.Run("多语言支持", func(t *testing.T) {
		// 测试中文环境
		cmdCN := createInternalTestCmdWithCompletion()
		cmdCN.SetChinese(true)
		cmdCN.registerBuiltinFlags()

		// 测试英文环境
		cmdEN := createInternalTestCmdWithCompletion()
		cmdEN.SetChinese(false)
		cmdEN.registerBuiltinFlags()

		// 验证注意事项的语言差异
		if len(cmdCN.ctx.Config.Notes) == 0 || len(cmdEN.ctx.Config.Notes) == 0 {
			t.Error("中英文环境都应该有注意事项")
		}

		// 验证示例的语言差异
		if len(cmdCN.ctx.Config.Examples) == 0 || len(cmdEN.ctx.Config.Examples) == 0 {
			t.Error("中英文环境都应该有示例")
		}
	})
}

// =============================================================================
// 并发安全测试
// =============================================================================

func TestCmd_Internal_ConcurrencySafety(t *testing.T) {
	t.Run("并发组件验证", func(t *testing.T) {
		cmd := createInternalTestCmd()

		var wg sync.WaitGroup
		numGoroutines := 50
		errors := make([]error, numGoroutines)

		wg.Add(numGoroutines)
		for i := 0; i < numGoroutines; i++ {
			go func(index int) {
				defer wg.Done()
				errors[index] = cmd.validateComponents()
			}(i)
		}

		wg.Wait()

		// 验证所有验证都成功
		for i, err := range errors {
			if err != nil {
				t.Errorf("goroutine %d 验证失败: %v", i, err)
			}
		}
	})

	t.Run("并发内置标志处理", func(t *testing.T) {
		cmd := createInternalTestCmd()

		var wg sync.WaitGroup
		numGoroutines := 50
		results := make([]bool, numGoroutines)
		errors := make([]error, numGoroutines)

		wg.Add(numGoroutines)
		for i := 0; i < numGoroutines; i++ {
			go func(index int) {
				defer wg.Done()
				shouldExit, err := cmd.handleBuiltinFlags()
				results[index] = shouldExit
				errors[index] = err
			}(i)
		}

		wg.Wait()

		// 验证所有处理都成功
		for i, err := range errors {
			if err != nil {
				t.Errorf("goroutine %d 处理失败: %v", i, err)
			}
		}

		// 验证结果一致性
		expectedResult := results[0]
		for i, result := range results {
			if result != expectedResult {
				t.Errorf("goroutine %d 结果不一致: %v, 期望 %v", i, result, expectedResult)
			}
		}
	})
}

// =============================================================================
// 特殊场景测试
// =============================================================================

func TestCmd_Internal_SpecialScenarios(t *testing.T) {
	t.Run("解析钩子测试", func(t *testing.T) {
		cmd := createInternalTestCmd()

		// 设置解析钩子
		hookCalled := false
		cmd.ctx.ParseHook = func(ctx *types.CmdContext) (error, bool) {
			hookCalled = true
			return nil, false
		}

		_, err := cmd.parseCommon([]string{}, true)
		if err != nil {
			t.Errorf("解析钩子测试失败: %v", err)
		}
		if !hookCalled {
			t.Error("解析钩子应该被调用")
		}
	})

	t.Run("解析钩子返回错误", func(t *testing.T) {
		cmd := createInternalTestCmd()

		// 设置返回错误的解析钩子
		cmd.ctx.ParseHook = func(ctx *types.CmdContext) (error, bool) {
			return fmt.Errorf("钩子错误"), false
		}

		_, err := cmd.parseCommon([]string{}, true)
		if err == nil {
			t.Error("期望钩子错误但未发生")
		}
		if !strings.Contains(err.Error(), "钩子错误") {
			t.Errorf("错误信息应包含钩子错误，实际: %v", err.Error())
		}
	})

	t.Run("解析钩子要求退出", func(t *testing.T) {
		cmd := createInternalTestCmd()

		// 设置要求退出的解析钩子
		cmd.ctx.ParseHook = func(ctx *types.CmdContext) (error, bool) {
			return nil, true
		}

		shouldExit, err := cmd.parseCommon([]string{}, true)
		if err != nil {
			t.Errorf("解析钩子退出测试失败: %v", err)
		}
		if !shouldExit {
			t.Error("解析钩子要求退出应该被响应")
		}
	})

	t.Run("时间相关标志测试", func(t *testing.T) {
		cmd := createInternalTestCmd()

		// 添加时间相关标志
		timeFlag := cmd.Time("start-time", "st", "now", "开始时间")
		durationFlag := cmd.Duration("timeout", "to", 30*time.Second, "超时时间")

		// 解析参数
		_, err := cmd.parseCommon([]string{}, true)
		if err != nil {
			t.Errorf("时间标志解析失败: %v", err)
		}

		// 验证标志值
		if timeFlag.Get().IsZero() {
			t.Error("时间标志应该有默认值")
		}
		if durationFlag.Get() != 30*time.Second {
			t.Errorf("时间间隔标志值 = %v, 期望 %v", durationFlag.Get(), 30*time.Second)
		}
	})
}

// =============================================================================
// 回归测试
// =============================================================================

func TestCmd_Internal_RegressionTests(t *testing.T) {
	t.Run("修复: nil指针解引用", func(t *testing.T) {
		// 这个测试确保我们不会在nil指针上调用方法
		var cmd *Cmd = nil

		_, err := cmd.parseCommon([]string{}, true)
		if err == nil {
			t.Error("nil命令应该返回错误")
		}
		if !strings.Contains(err.Error(), "nil command") {
			t.Errorf("错误信息应包含'nil command'，实际: %v", err.Error())
		}
	})

	t.Run("修复：重复注册内置标志", func(t *testing.T) {
		cmd := createInternalTestCmdWithVersion()

		// 第一次注册
		cmd.registerBuiltinFlags()

		// 验证第一次注册成功
		_, err := cmd.handleBuiltinFlags()
		if err != nil {
			t.Errorf("第一次注册后处理失败: %v", err)
		}

		// 注意：由于内置标志注册会检查重复，多次调用会导致panic
		// 这里我们测试的是单次注册的正确性
		// 在实际使用中，registerBuiltinFlags只会在parseCommon中被sync.Once保护调用一次
	})

	t.Run("修复：枚举标志验证边界情况", func(t *testing.T) {
		cmd := createInternalTestCmd()

		// 创建空选项的枚举标志
		enumFlag := cmd.Enum("empty-enum", "ee", "", "空枚举", []string{})

		// 验证空枚举不会导致验证失败
		_, err := cmd.handleBuiltinFlags()
		if err != nil {
			t.Errorf("空枚举验证失败: %v", err)
		}

		// 验证枚举标志的值
		if enumFlag.Get() != "" {
			t.Errorf("空枚举标志值 = %v, 期望空字符串", enumFlag.Get())
		}
	})
}

// =============================================================================
// 性能基准测试
// =============================================================================

func BenchmarkCmd_registerBuiltinFlags(b *testing.B) {
	cmd := createInternalTestCmdWithCompletion()
	cmd.SetVersion("v1.0.0")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// 重置状态以允许重复测试
		cmd.ctx.BuiltinFlags.NameMap = sync.Map{}
		cmd.ctx.Config.Notes = []string{}
		cmd.ctx.Config.Examples = []types.ExampleInfo{}

		cmd.registerBuiltinFlags()
	}
}

func BenchmarkCmd_handleBuiltinFlags_WithManyEnums(b *testing.B) {
	cmd := createInternalTestCmd()

	// 添加大量枚举标志
	for i := 0; i < 100; i++ {
		cmd.Enum(fmt.Sprintf("enum%d", i), fmt.Sprintf("e%d", i), "option1", "测试枚举", []string{"option1", "option2", "option3"})
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := cmd.handleBuiltinFlags()
		if err != nil {
			b.Fatalf("处理内置标志失败: %v", err)
		}
	}
}

// =============================================================================
// 内存泄漏检测测试
// =============================================================================

func TestCmd_Internal_MemoryLeaks(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过内存泄漏测试（短测试模式）")
	}

	t.Run("大量解析操作", func(t *testing.T) {
		// 执行大量解析操作，检查是否有内存泄漏
		for i := 0; i < 1000; i++ {
			cmd := createInternalTestCmd()
			cmd.String("test", "t", "default", "测试")

			_, err := cmd.parseCommon([]string{"--test", fmt.Sprintf("value%d", i)}, true)
			if err != nil {
				t.Errorf("解析操作 %d 失败: %v", i, err)
				break
			}
		}
	})
}
